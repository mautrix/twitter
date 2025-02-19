package connector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
	"go.mau.fi/util/exmime"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	bridgeEvt "maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func RemoveEntityLinkFromText(msgPart *bridgev2.ConvertedMessagePart, indices []int) {
	start, end := indices[0], indices[1]
	msgPart.Content.Body = msgPart.Content.Body[:start-1] + msgPart.Content.Body[end:]
}

func (tc *TwitterClient) ConversationToChatInfo(conv *types.Conversation) *bridgev2.ChatInfo {
	memberList := tc.ParticipantsToMemberList(conv.Participants)
	var userLocal bridgev2.UserLocalPortalInfo
	if conv.Muted {
		userLocal.MutedUntil = ptr.Ptr(bridgeEvt.MutedForever)
	} else {
		userLocal.MutedUntil = ptr.Ptr(bridgev2.Unmuted)
	}
	return &bridgev2.ChatInfo{
		Name:        &conv.Name,
		Avatar:      tc.MakeAvatar(conv.AvatarImageHttps),
		Members:     memberList,
		Type:        tc.ConversationTypeToRoomType(conv.Type),
		UserLocal:   &userLocal,
		CanBackfill: true,
	}
}

func (tc *TwitterClient) ConversationTypeToRoomType(convType types.ConversationType) *database.RoomType {
	var roomType database.RoomType
	switch convType {
	case types.ONE_TO_ONE:
		roomType = database.RoomTypeDM
	case types.GROUP_DM:
		roomType = database.RoomTypeGroupDM
	}

	return &roomType
}

func (tc *TwitterClient) UsersToMemberList(users []types.User) *bridgev2.ChatMemberList {
	selfUserID := tc.client.GetCurrentUserID()

	memberMap := map[networkid.UserID]bridgev2.ChatMember{}
	for _, user := range users {
		memberMap[networkid.UserID(user.IDStr)] = tc.UserToChatMember(user, user.IDStr == selfUserID)
	}

	return &bridgev2.ChatMemberList{
		IsFull:           true,
		TotalMemberCount: len(users),
		MemberMap:        memberMap,
	}
}

func (tc *TwitterClient) ParticipantsToMemberList(participants []types.Participant) *bridgev2.ChatMemberList {
	selfUserID := tc.client.GetCurrentUserID()
	memberMap := map[networkid.UserID]bridgev2.ChatMember{}
	for _, participant := range participants {
		memberMap[networkid.UserID(participant.UserID)] = tc.ParticipantToChatMember(participant, participant.UserID == selfUserID)
	}

	return &bridgev2.ChatMemberList{
		IsFull:           true,
		TotalMemberCount: len(participants),
		MemberMap:        memberMap,
	}
}

func (tc *TwitterClient) UserToChatMember(user types.User, isFromMe bool) bridgev2.ChatMember {
	return bridgev2.ChatMember{
		EventSender: bridgev2.EventSender{
			IsFromMe: isFromMe,
			Sender:   networkid.UserID(user.IDStr),
		},
		UserInfo: &bridgev2.UserInfo{
			Name:   &user.Name,
			Avatar: tc.MakeAvatar(user.ProfileImageURL),
		},
	}
}

func (tc *TwitterClient) ParticipantToChatMember(participant types.Participant, isFromMe bool) bridgev2.ChatMember {
	return bridgev2.ChatMember{
		EventSender: bridgev2.EventSender{
			IsFromMe: isFromMe,
			Sender:   networkid.UserID(participant.UserID),
		},
		UserInfo: tc.GetUserInfoBridge(participant.UserID),
	}
}

func (tc *TwitterClient) GetUserInfoBridge(userID string) *bridgev2.UserInfo {
	var userinfo *bridgev2.UserInfo
	if userCacheEntry, ok := tc.userCache[userID]; ok {
		userinfo = &bridgev2.UserInfo{
			Name:        ptr.Ptr(tc.connector.Config.FormatDisplayname(userCacheEntry.ScreenName, userCacheEntry.Name)),
			Avatar:      tc.MakeAvatar(userCacheEntry.ProfileImageURL),
			Identifiers: []string{fmt.Sprintf("twitter:%s", userCacheEntry.ScreenName)},
		}
	}
	return userinfo
}

func (tc *TwitterClient) TwitterAttachmentToMatrix(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, msg *types.MessageData) (*bridgev2.ConvertedMessagePart, []int, error) {
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

func (tc *TwitterClient) MakeAvatar(avatarURL string) *bridgev2.Avatar {
	return &bridgev2.Avatar{
		ID: networkid.AvatarID(avatarURL),
		Get: func(ctx context.Context) ([]byte, error) {
			resp, err := tc.downloadFile(ctx, avatarURL)
			if err != nil {
				return nil, err
			}
			data, err := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			return data, err
		},
		Remove: avatarURL == "",
	}
}

func (tc *TwitterClient) downloadFile(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}

	headers := twittermeow.BaseHeaders.Clone()
	headers.Set("Cookie", tc.client.GetCookieString())
	req.Header = headers

	getResp, err := tc.client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	return getResp, nil
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
