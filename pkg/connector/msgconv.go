// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2025 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
	"go.mau.fi/util/exmime"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/util/ffmpeg"

	"go.mau.fi/mautrix-twitter/pkg/connector/twitterfmt"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (tc *TwitterClient) convertEditToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, existing []*database.Message, data *types.MessageData) (*bridgev2.ConvertedEdit, error) {
	meta, ok := existing[0].Metadata.(*MessageMetadata)
	if !ok || meta == nil {
		meta = &MessageMetadata{}
	}
	if data.EditCount == 0 {
		data.EditCount = meta.EditCount + 1
	} else if meta.EditCount >= data.EditCount {
		return nil, fmt.Errorf("%w: db edit count %d >= remote edit count %d", bridgev2.ErrIgnoringRemoteEvent, meta.EditCount, data.EditCount)
	}
	data.Text = strings.TrimPrefix(data.Text, "Edited: ")
	editPart := tc.convertToMatrix(ctx, portal, intent, data).Parts[0].ToEditPart(existing[0])
	editPart.Part.Metadata = &MessageMetadata{EditCount: data.EditCount}
	return &bridgev2.ConvertedEdit{
		ModifiedParts: []*bridgev2.ConvertedEditPart{editPart},
	}, nil
}

func (tc *TwitterClient) convertToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, msg *types.MessageData) *bridgev2.ConvertedMessage {
	var replyTo *networkid.MessageOptionalPartID
	if msg.ReplyData.ID != "" {
		replyTo = &networkid.MessageOptionalPartID{
			MessageID: networkid.MessageID(msg.ReplyData.ID),
		}
	}

	tc.populateMentionIDs(msg)

	textPart := &bridgev2.ConvertedMessagePart{
		ID:      "",
		Type:    event.EventMessage,
		Content: twitterfmt.Parse(ctx, portal, msg),
	}

	parts := make([]*bridgev2.ConvertedMessagePart, 0)

	if msg.Attachment != nil {
		convertedAttachmentPart, indices, err := tc.twitterAttachmentToMatrix(ctx, portal, intent, msg)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to convert attachment")
			parts = append(parts, &bridgev2.ConvertedMessagePart{
				ID:   "",
				Type: event.EventMessage,
				Content: &event.MessageEventContent{
					MsgType: event.MsgNotice,
					Body:    "Failed to convert attachment from Twitter",
				},
			})
		} else {
			if msg.Attachment.Card != nil || msg.Attachment.Tweet != nil {
				textPart.Content.BeeperLinkPreviews = convertedAttachmentPart.Content.BeeperLinkPreviews
			} else {
				parts = append(parts, convertedAttachmentPart)
				removeEntityLinkFromText(textPart, indices)
			}
		}
	}

	if len(textPart.Content.Body) > 0 || len(textPart.Content.BeeperLinkPreviews) > 0 {
		parts = append(parts, textPart)
	}
	displayName := tc.getDisplayNameForUser(ctx, msg.SenderID)
	if displayName == "" {
		displayName = msg.SenderID
	}

	for _, part := range parts {
		part.DBMetadata = &MessageMetadata{
			EditCount:         msg.EditCount,
			MessageText:       msg.Text,
			SenderID:          msg.SenderID,
			SenderDisplayName: displayName,
			ReplyAttachments:  filterReplyPreviewAttachments(msg.OriginalAttachments),
		}
	}

	cm := &bridgev2.ConvertedMessage{
		ReplyTo: replyTo,
		Parts:   parts,
	}
	cm.MergeCaption()

	return cm
}

func removeEntityLinkFromText(msgPart *bridgev2.ConvertedMessagePart, indices []int) {
	if len(indices) < 2 {
		return
	}

	start, end := indices[0], indices[1]
	if start <= 0 || end <= start || end > len(msgPart.Content.Body) {
		return
	}

	msgPart.Content.Body = msgPart.Content.Body[:start-1] + msgPart.Content.Body[end:]
}

func (tc *TwitterClient) getDisplayNameForUser(ctx context.Context, userID string) string {
	if userID == "" {
		return ""
	}

	tc.userCacheLock.RLock()
	user := tc.userCache[userID]
	tc.userCacheLock.RUnlock()

	if user == nil {
		if err := tc.ensureUsersInCacheByID(ctx, []string{userID}); err != nil {
			zerolog.Ctx(ctx).Debug().
				Err(err).
				Str("user_id", userID).
				Msg("Failed to fetch user while resolving display name")
			return ""
		}
		tc.userCacheLock.RLock()
		user = tc.userCache[userID]
		tc.userCacheLock.RUnlock()
	}

	if user == nil {
		return ""
	}
	if user.Name != "" {
		return user.Name
	}
	return user.ScreenName
}

func (tc *TwitterClient) populateMentionIDs(msg *types.MessageData) {
	if msg == nil || msg.Entities == nil || len(msg.Entities.UserMentions) == 0 {
		return
	}

	var textRunes []rune
	textLen := 0
	if msg.Text != "" {
		textRunes = []rune(msg.Text)
		textLen = len(textRunes)
	}

	for i := range msg.Entities.UserMentions {
		mention := &msg.Entities.UserMentions[i]
		if mention.IDStr != "" {
			continue
		}

		screenName := strings.TrimSpace(mention.ScreenName)
		screenName = strings.TrimPrefix(screenName, "@")

		if screenName == "" && len(mention.Indices) >= 2 && textLen > 0 {
			start := mention.Indices[0]
			end := mention.Indices[1]
			if start < 0 {
				start = 0
			}
			if end < start {
				end = start
			}
			if end > textLen {
				end = textLen
			}
			if start < end {
				mentionText := strings.TrimSpace(string(textRunes[start:end]))
				mentionText = strings.TrimPrefix(mentionText, "@")
				screenName = mentionText
			}
		}

		if screenName == "" {
			continue
		}
		mention.ScreenName = screenName

		idStr, id := tc.lookupUserIDByScreenName(screenName)
		if idStr == "" {
			continue
		}
		mention.IDStr = idStr
		if mention.ID == 0 && id != 0 {
			mention.ID = id
		}
	}
}

func (tc *TwitterClient) lookupUserIDByScreenName(screenName string) (string, int64) {
	normalized := strings.TrimSpace(screenName)
	normalized = strings.TrimPrefix(normalized, "@")
	if normalized == "" {
		return "", 0
	}

	if tc.userLogin != nil {
		if tc.userLogin.RemoteName != "" && strings.EqualFold(tc.userLogin.RemoteName, normalized) {
			return string(tc.userLogin.ID), 0
		}
		if tc.userLogin.RemoteProfile.Username != "" && strings.EqualFold(tc.userLogin.RemoteProfile.Username, normalized) {
			return string(tc.userLogin.ID), 0
		}
	}

	tc.userCacheLock.RLock()
	defer tc.userCacheLock.RUnlock()
	for _, user := range tc.userCache {
		if user == nil {
			continue
		}
		if strings.EqualFold(user.ScreenName, normalized) {
			if user.IDStr != "" {
				return user.IDStr, user.ID
			}
			if user.ID != 0 {
				return strconv.FormatInt(user.ID, 10), user.ID
			}
			return "", 0
		}
	}
	return "", 0
}

func (tc *TwitterClient) twitterAttachmentToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, msg *types.MessageData) (*bridgev2.ConvertedMessagePart, []int, error) {
	attachment := msg.Attachment
	var attachmentInfo *types.AttachmentInfo
	var attachmentURL string
	var mimeType string
	var msgType event.MessageType
	extraInfo := map[string]any{}
	if attachment.Photo != nil {
		attachmentInfo = attachment.Photo
		mimeType = "image/jpeg" // attachment doesn't include this specifically
		msgType = event.MsgImage
		attachmentURL = attachmentInfo.MediaURLHTTPS
	} else if attachment.Video != nil || attachment.AnimatedGif != nil {
		if attachment.AnimatedGif != nil {
			attachmentInfo = attachment.AnimatedGif
			extraInfo["fi.mau.gif"] = true
			extraInfo["fi.mau.loop"] = true
			extraInfo["fi.mau.autoplay"] = true
			extraInfo["fi.mau.hide_controls"] = true
			extraInfo["fi.mau.no_audio"] = true
		} else {
			attachmentInfo = attachment.Video
		}
		mimeType = "video/mp4"
		msgType = event.MsgVideo

		// For XChat encrypted media, skip variant lookup - we'll use MediaHashKey
		if attachmentInfo.MediaHashKey == "" {
			highestBitRateVariant, err := attachmentInfo.VideoInfo.GetHighestBitrateVariant()
			if err != nil {
				return nil, nil, err
			}
			attachmentURL = highestBitRateVariant.URL
		}
	} else if attachment.Card != nil {
		var urls []types.URLs
		if msg.Entities != nil {
			urls = msg.Entities.URLs
		}
		content := event.MessageEventContent{
			MsgType:            event.MsgText,
			BeeperLinkPreviews: []*event.BeeperLinkPreview{tc.attachmentCardToMatrix(ctx, portal, intent, attachment, urls, msg.ConversationKeyVersion)},
		}
		return &bridgev2.ConvertedMessagePart{
			ID:      networkid.PartID(""),
			Type:    event.EventMessage,
			Content: &content,
		}, []int{0, 0}, nil
	} else if attachment.Tweet != nil {
		content := event.MessageEventContent{
			MsgType:            event.MsgText,
			BeeperLinkPreviews: []*event.BeeperLinkPreview{tc.attachmentTweetToMatrix(ctx, portal, intent, attachment.Tweet)},
		}
		return &bridgev2.ConvertedMessagePart{
			ID:      networkid.PartID(""),
			Type:    event.EventMessage,
			Content: &content,
		}, []int{0, 0}, nil
	} else {
		return nil, nil, fmt.Errorf("unsupported attachment type")
	}

	content := event.MessageEventContent{
		Info: &event.FileInfo{
			MimeType: mimeType,
			Width:    attachmentInfo.OriginalInfo.Width,
			Height:   attachmentInfo.OriginalInfo.Height,
			Duration: attachmentInfo.VideoInfo.DurationMillis,
		},
		MsgType: msgType,
		Body:    attachmentInfo.IDStr,
	}
	if content.Body == "" {
		content.Body = strings.TrimPrefix(string(msgType), "m.")
	}

	audioOnly := attachment.Video != nil && attachment.Video.AudioOnly

	var err error
	// Check if this is XChat encrypted media
	if attachmentInfo.MediaHashKey != "" {
		// Download and decrypt XChat media
		conversationID := string(portal.ID)
		decryptedData, downloadErr := tc.client.DownloadXChatMedia(ctx, conversationID, attachmentInfo.MediaHashKey, msg.ConversationKeyVersion)
		if downloadErr != nil {
			return nil, nil, fmt.Errorf("failed to download XChat media: %w", downloadErr)
		}
		content.Info.Size = len(decryptedData)
		content.URL, content.File, err = intent.UploadMediaStream(ctx, portal.MXID, int64(len(decryptedData)), audioOnly, func(file io.Writer) (*bridgev2.FileStreamResult, error) {
			n, err := io.Copy(file, bytes.NewReader(decryptedData))
			if err != nil {
				return nil, err
			}
			if audioOnly && ffmpeg.Supported() {
				outFile, err := ffmpeg.ConvertPath(ctx, file.(*os.File).Name(), ".ogg", []string{}, []string{"-vn", "-c:a", "libopus"}, false)
				if err == nil {
					mimeType = "audio/ogg"
					content.Info.MimeType = mimeType
					content.Info.Width = 0
					content.Info.Height = 0
					content.MsgType = event.MsgAudio
					content.Body += ".ogg"
					return &bridgev2.FileStreamResult{
						ReplacementFile: outFile,
						MimeType:        mimeType,
						FileName:        content.Body,
					}, nil
				} else {
					zerolog.Ctx(ctx).Warn().Err(err).Msg("Failed to convert voice message to ogg")
				}
			} else {
				content.Info.Size = int(n)
			}
			ext := exmime.ExtensionFromMimetype(mimeType)
			if !strings.HasSuffix(content.Body, ext) {
				content.Body += ext
			}
			return &bridgev2.FileStreamResult{
				MimeType: content.Info.MimeType,
				FileName: content.Body,
			}, nil
		})
	} else {
		// Legacy media download
		fileResp, downloadErr := downloadFile(ctx, tc.client, attachmentURL)
		if downloadErr != nil {
			return nil, nil, downloadErr
		}
		if tc.connector.directMedia {
			content.URL, err = tc.connector.br.Matrix.GenerateContentURI(ctx, MakeMediaID(portal.Receiver, attachmentURL))
		} else {
			content.URL, content.File, err = intent.UploadMediaStream(ctx, portal.MXID, fileResp.ContentLength, audioOnly, func(file io.Writer) (*bridgev2.FileStreamResult, error) {
				n, err := io.Copy(file, fileResp.Body)
				if err != nil {
					return nil, err
				}
				if audioOnly && ffmpeg.Supported() {
					outFile, err := ffmpeg.ConvertPath(ctx, file.(*os.File).Name(), ".ogg", []string{}, []string{"-vn", "-c:a", "libopus"}, false)
					if err == nil {
						mimeType = "audio/ogg"
						content.Info.MimeType = mimeType
						content.Info.Width = 0
						content.Info.Height = 0
						content.MsgType = event.MsgAudio
						content.Body += ".ogg"
						return &bridgev2.FileStreamResult{
							ReplacementFile: outFile,
							MimeType:        mimeType,
							FileName:        content.Body,
						}, nil
					} else {
						zerolog.Ctx(ctx).Warn().Err(err).Msg("Failed to convert voice message to ogg")
					}
				} else {
					content.Info.Size = int(n)
				}
				ext := exmime.ExtensionFromMimetype(mimeType)
				if !strings.HasSuffix(content.Body, ext) {
					content.Body += ext
				}
				return &bridgev2.FileStreamResult{
					MimeType: content.Info.MimeType,
					FileName: content.Body,
				}, nil
			})
		}
	}
	if err != nil {
		return nil, nil, err
	}

	if audioOnly {
		content.MSC3245Voice = &event.MSC3245Voice{}
	}
	return &bridgev2.ConvertedMessagePart{
		ID:      networkid.PartID(""),
		Type:    event.EventMessage,
		Content: &content,
		Extra: map[string]any{
			"info": extraInfo,
		},
	}, attachmentInfo.Indices, nil
}

func downloadFile(ctx context.Context, cli *twittermeow.Client, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	headers := twittermeow.BaseHeaders.Clone()
	headers.Set("Cookie", cli.GetCookieString())
	req.Header = headers

	getResp, err := cli.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	return getResp, nil
}

func (tc *TwitterClient) attachmentCardToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, attachment *types.Attachment, urls []types.URLs, keyVersion string) *event.BeeperLinkPreview {
	card := attachment.Card
	canonicalURL := card.BindingValues.CardURL.StringValue
	for _, url := range urls {
		if url.URL == canonicalURL {
			canonicalURL = url.ExpandedURL
			break
		}
	}
	preview := &event.BeeperLinkPreview{
		LinkPreview: event.LinkPreview{
			CanonicalURL: canonicalURL,
			Title:        card.BindingValues.Title.StringValue,
			Description:  card.BindingValues.Description.StringValue,
		},
	}

	// Download banner image if available (XChat encrypted)
	if attachment.URLBannerMediaHashKey != "" {
		conversationID := string(portal.ID)
		decryptedData, err := tc.client.DownloadXChatMedia(ctx, conversationID, attachment.URLBannerMediaHashKey, keyVersion)
		if err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Msg("Failed to download URL attachment banner image")
		} else {
			preview.ImageType = "image/jpeg"
			preview.ImageSize = event.IntOrString(len(decryptedData))
			preview.ImageURL, _, err = intent.UploadMediaStream(ctx, portal.MXID, int64(len(decryptedData)), false, func(file io.Writer) (*bridgev2.FileStreamResult, error) {
				_, err := io.Copy(file, bytes.NewReader(decryptedData))
				if err != nil {
					return nil, err
				}
				return &bridgev2.FileStreamResult{
					MimeType: "image/jpeg",
					FileName: "banner.jpeg",
				}, nil
			})
			if err != nil {
				zerolog.Ctx(ctx).Warn().Err(err).Msg("Failed to upload URL attachment banner image to Matrix")
			}
		}
	}

	return preview
}

func (tc *TwitterClient) attachmentTweetToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, tweet *types.AttachmentTweet) *event.BeeperLinkPreview {
	// Handle XChat Post attachments with empty Status (only URL available)
	if tweet.Status.FullText == "" && tweet.ExpandedURL != "" {
		return &event.BeeperLinkPreview{
			LinkPreview: event.LinkPreview{
				CanonicalURL: tweet.ExpandedURL,
				Title:        "Post on X",
			},
		}
	}

	linkPreview := event.LinkPreview{
		CanonicalURL: tweet.ExpandedURL,
		Title:        tweet.Status.User.Name + " on X",
		Description:  tweet.Status.FullText,
	}
	medias := tweet.Status.Entities.Media
	if len(medias) > 0 {
		media := medias[0]
		if media.Type == "photo" {
			resp, err := downloadFile(ctx, tc.client, media.MediaURLHTTPS)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("failed to download tweet image")
			} else {
				linkPreview.ImageType = "image/jpeg"
				linkPreview.ImageWidth = event.IntOrString(media.OriginalInfo.Width)
				linkPreview.ImageHeight = event.IntOrString(media.OriginalInfo.Height)
				linkPreview.ImageSize = event.IntOrString(resp.ContentLength)
				linkPreview.ImageURL, _, err = intent.UploadMediaStream(ctx, portal.MXID, resp.ContentLength, false, func(file io.Writer) (*bridgev2.FileStreamResult, error) {
					_, err := io.Copy(file, resp.Body)
					if err != nil {
						return nil, err
					}
					return &bridgev2.FileStreamResult{
						MimeType: linkPreview.ImageType,
						FileName: "image.jpeg",
					}, nil
				})
				if err != nil {
					zerolog.Ctx(ctx).Err(err).Msg("failed to upload tweet image to Matrix")
				}
			}
		}
	}
	return &event.BeeperLinkPreview{
		LinkPreview: linkPreview,
	}
}
