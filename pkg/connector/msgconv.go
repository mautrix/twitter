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
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"go.mau.fi/util/exmime"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	bridgeEvt "maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/connector/twitterfmt"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (tc *TwitterClient) convertEditToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, existing []*database.Message, data *types.MessageData) (*bridgev2.ConvertedEdit, error) {
	if ec := existing[0].Metadata.(*MessageMetadata).EditCount; ec >= data.EditCount {
		return nil, fmt.Errorf("%w: db edit count %d >= remote edit count %d", bridgev2.ErrIgnoringRemoteEvent, ec, data.EditCount)
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

	textPart := &bridgev2.ConvertedMessagePart{
		ID:      "",
		Type:    bridgeEvt.EventMessage,
		Content: twitterfmt.Parse(ctx, portal, msg),
	}

	parts := make([]*bridgev2.ConvertedMessagePart, 0)

	if msg.Attachment != nil {
		convertedAttachmentPart, indices, err := tc.twitterAttachmentToMatrix(ctx, portal, intent, msg)
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Failed to convert attachment")
			parts = append(parts, &bridgev2.ConvertedMessagePart{
				ID:   "",
				Type: bridgeEvt.EventMessage,
				Content: &bridgeEvt.MessageEventContent{
					MsgType: bridgeEvt.MsgNotice,
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

	if len(textPart.Content.Body) > 0 {
		parts = append(parts, textPart)
	}
	for _, part := range parts {
		part.DBMetadata = &MessageMetadata{EditCount: msg.EditCount}
	}

	cm := &bridgev2.ConvertedMessage{
		ReplyTo: replyTo,
		Parts:   parts,
	}
	cm.MergeCaption()

	return cm
}

func removeEntityLinkFromText(msgPart *bridgev2.ConvertedMessagePart, indices []int) {
	start, end := indices[0], indices[1]
	msgPart.Content.Body = msgPart.Content.Body[:start-1] + msgPart.Content.Body[end:]
}

func (tc *TwitterClient) twitterAttachmentToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, msg *types.MessageData) (*bridgev2.ConvertedMessagePart, []int, error) {
	attachment := msg.Attachment
	var attachmentInfo *types.AttachmentInfo
	var attachmentURL string
	var mimeType string
	var msgType bridgeEvt.MessageType
	extraInfo := map[string]any{}
	if attachment.Photo != nil {
		attachmentInfo = attachment.Photo
		mimeType = "image/jpeg" // attachment doesn't include this specifically
		msgType = bridgeEvt.MsgImage
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
		msgType = bridgeEvt.MsgVideo

		highestBitRateVariant, err := attachmentInfo.VideoInfo.GetHighestBitrateVariant()
		if err != nil {
			return nil, nil, err
		}
		attachmentURL = highestBitRateVariant.URL
	} else if attachment.Card != nil {
		content := bridgeEvt.MessageEventContent{
			MsgType:            bridgeEvt.MsgText,
			BeeperLinkPreviews: []*bridgeEvt.BeeperLinkPreview{tc.attachmentCardToMatrix(ctx, attachment.Card, msg.Entities.URLs)},
		}
		return &bridgev2.ConvertedMessagePart{
			ID:      networkid.PartID(""),
			Type:    bridgeEvt.EventMessage,
			Content: &content,
		}, []int{0, 0}, nil
	} else if attachment.Tweet != nil {
		content := bridgeEvt.MessageEventContent{
			MsgType:            bridgeEvt.MsgText,
			BeeperLinkPreviews: []*bridgeEvt.BeeperLinkPreview{tc.attachmentTweetToMatrix(ctx, portal, intent, attachment.Tweet)},
		}
		return &bridgev2.ConvertedMessagePart{
			ID:      networkid.PartID(""),
			Type:    bridgeEvt.EventMessage,
			Content: &content,
		}, []int{0, 0}, nil
	} else {
		return nil, nil, fmt.Errorf("unsupported attachment type")
	}

	fileResp, err := tc.downloadFile(ctx, attachmentURL)
	if err != nil {
		return nil, nil, err
	}
	content := bridgeEvt.MessageEventContent{
		Info: &bridgeEvt.FileInfo{
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
	ext := exmime.ExtensionFromMimetype(mimeType)
	if !strings.HasSuffix(content.Body, ext) {
		content.Body += ext
	}

	content.URL, content.File, err = intent.UploadMediaStream(ctx, portal.MXID, fileResp.ContentLength, true, func(file io.Writer) (*bridgev2.FileStreamResult, error) {
		n, err := io.Copy(file, fileResp.Body)
		if err != nil {
			return nil, err
		}
		content.Info.Size = int(n)
		return &bridgev2.FileStreamResult{
			MimeType: content.Info.MimeType,
			FileName: content.Body,
		}, nil
	})

	if err != nil {
		return nil, nil, err
	}

	return &bridgev2.ConvertedMessagePart{
		ID:      networkid.PartID(""),
		Type:    bridgeEvt.EventMessage,
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

func (tc *TwitterClient) downloadFile(ctx context.Context, url string) (*http.Response, error) {
	return downloadFile(ctx, tc.client, url)
}

func (tc *TwitterClient) attachmentCardToMatrix(ctx context.Context, card *types.AttachmentCard, urls []types.URLs) *bridgeEvt.BeeperLinkPreview {
	canonicalURL := card.BindingValues.CardURL.StringValue
	for _, url := range urls {
		if url.URL == canonicalURL {
			canonicalURL = url.ExpandedURL
			break
		}
	}
	preview := &bridgeEvt.BeeperLinkPreview{
		LinkPreview: bridgeEvt.LinkPreview{
			CanonicalURL: canonicalURL,
			Title:        card.BindingValues.Title.StringValue,
			Description:  card.BindingValues.Description.StringValue,
		},
	}
	return preview
}

func (tc *TwitterClient) attachmentTweetToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, tweet *types.AttachmentTweet) *bridgeEvt.BeeperLinkPreview {
	linkPreview := bridgeEvt.LinkPreview{
		CanonicalURL: tweet.ExpandedURL,
		Title:        tweet.Status.User.Name + " on X",
		Description:  tweet.Status.FullText,
	}
	medias := tweet.Status.Entities.Media
	if len(medias) > 0 {
		media := medias[0]
		if media.Type == "photo" {
			resp, err := tc.downloadFile(ctx, media.MediaURLHTTPS)
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("failed to download tweet image")
			} else {
				linkPreview.ImageType = "image/jpeg"
				linkPreview.ImageWidth = media.OriginalInfo.Width
				linkPreview.ImageHeight = media.OriginalInfo.Height
				linkPreview.ImageSize = int(resp.ContentLength)
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
	return &bridgeEvt.BeeperLinkPreview{
		LinkPreview: linkPreview,
	}
}
