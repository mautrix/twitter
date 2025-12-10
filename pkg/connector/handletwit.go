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
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/bridgev2/simplevent"
	"maunium.net/go/mautrix/bridgev2/status"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/methods"
)

func (tc *TwitterClient) HandleTwitterEvent(rawEvt types.TwitterEvent, inbox *response.TwitterInboxData) bool {
	if rawEvt == nil && inbox != nil {
		prevState := tc.userLogin.BridgeState.GetPrevUnsent().StateEvent
		if prevState == "" || prevState == status.StateConnecting {
			tc.userLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnected})
		}
		tc.updateTwitterUserInfo(inbox)
		tc.updateTwitterReadReceipt(inbox)
		tc.userCacheLock.Lock()
		maps.Copy(tc.userCache, inbox.Users)
		tc.userCacheLock.Unlock()
		return true
	}
	isEdit := false
	if edit, ok := rawEvt.(*types.MessageEdit); ok {
		rawEvt = (*types.Message)(edit)
		isEdit = true
	}
	switch evt := rawEvt.(type) {
	case *types.PollingError:
		if evt.Error != nil {
			stateEvt := status.BridgeState{
				StateEvent: status.StateTransientDisconnect,
				Error:      "twitter-polling-error",
				Info: map[string]any{
					"go_error": evt.Error.Error(),
				},
			}
			if evt.IsAuth {
				stateEvt.StateEvent = status.StateBadCredentials
				stateEvt.Error = "twitter-invalid-credentials"
				go tc.Disconnect()
			}
			tc.userLogin.BridgeState.Send(stateEvt)
			if errors.Is(evt.Error, twittermeow.ErrCSRFMismatch) && !tc.reconnectAttempted.Swap(true) {
				tc.userLogin.Log.Info().Msg("Doing full reconnect due to 353 error")
				go tc.FullReconnect()
			}
		} else {
			tc.userLogin.BridgeState.Send(status.BridgeState{StateEvent: status.StateConnected})
		}
		return true
	case *types.Message:
		isFromMe := evt.MessageData.SenderID == string(tc.userLogin.ID)
		msgType := bridgev2.RemoteEventMessage
		if isEdit {
			msgType = bridgev2.RemoteEventEdit
		}
		conversation := inbox.GetConversationByID(evt.ConversationID)
		if isFromMe {
			tc.userLogin.Log.Debug().
				Str("sequence_id", evt.ID).
				Str("message_id", evt.MessageData.ID).
				Str("request_id", evt.RequestID).
				Msg("Received self XChat message")
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.Message[*types.MessageData]{
			EventMeta: simplevent.EventMeta{
				Type: msgType,
				LogContext: func(c zerolog.Context) zerolog.Context {
					if conversation != nil {
						c = c.
							Bool("conv_low_quality", conversation.LowQuality).
							Bool("conv_trusted", conversation.Trusted)
					} else {
						c = c.Bool("conversation_nil", true)
					}
					return c.
						Str("message_id", evt.MessageData.ID).
						Str("sender", evt.MessageData.SenderID).
						Bool("is_from_me", isFromMe)
				},
				PortalKey:    tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				CreatePortal: isFromMe || (conversation != nil && (conversation.Trusted || !conversation.LowQuality)),
				Sender:       tc.MakeEventSender(evt.MessageData.SenderID),
				StreamOrder:  methods.ParseSnowflakeInt(evt.ID),
				Timestamp:    methods.ParseSnowflake(evt.ID),
			},
			ID:            networkid.MessageID(evt.MessageData.ID),
			TransactionID: networkid.TransactionID(evt.RequestID),
			TargetMessage: networkid.MessageID(evt.MessageData.ID),
			Data:          &evt.MessageData,
			ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data *types.MessageData) (*bridgev2.ConvertedMessage, error) {
				return tc.convertToMatrix(ctx, portal, intent, data), nil
			},
			ConvertEditFunc: tc.convertEditToMatrix,
		}).Success
	case *types.ConversationRead:
		lastTarget := networkid.MessageID(evt.LastReadEventID)
		readUpTo := methods.ParseSnowflake(evt.LastReadEventID)
		if readUpTo.IsZero() {
			readUpTo = methods.ParseSnowflake(evt.ID)
		}
		readUpToStreamOrder := methods.ParseSnowflakeInt(evt.LastReadEventID)
		if readUpToStreamOrder == 0 {
			readUpToStreamOrder = methods.ParseSnowflakeInt(evt.ID)
		}
		var targets []networkid.MessageID
		if lastTarget != "" {
			targets = []networkid.MessageID{lastTarget}
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.Receipt{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventReadReceipt,
				PortalKey: tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				Sender:    tc.MakeEventSender(ParseUserLoginID(tc.userLogin.ID)),
				Timestamp: methods.ParseSnowflake(evt.ID),
				PreHandleFunc: func(ctx context.Context, portal *bridgev2.Portal) {
					if intent := tc.userLogin.User.DoublePuppet(ctx); intent != nil {
						_ = intent.EnsureJoined(ctx, portal.MXID)
					}
				},
			},
			LastTarget:          lastTarget,
			Targets:             targets,
			ReadUpTo:            readUpTo,
			ReadUpToStreamOrder: readUpToStreamOrder,
		}).Success
	case *types.MessageReactionCreate:
		portalKey := tc.makePortalKeyFromInbox(evt.ConversationID, inbox)
		wrappedEvt := tc.wrapReaction((*types.MessageReaction)(evt), portalKey, bridgev2.RemoteEventReaction)
		return tc.userLogin.QueueRemoteEvent(wrappedEvt).Success
	case *types.MessageReactionDelete:
		portalKey := tc.makePortalKeyFromInbox(evt.ConversationID, inbox)
		wrappedEvt := tc.wrapReaction((*types.MessageReaction)(evt), portalKey, bridgev2.RemoteEventReactionRemove)
		return tc.userLogin.QueueRemoteEvent(wrappedEvt).Success
	case *types.ConversationCreate:
		// honestly not sure when this is ever called... ? might be when they initialize the conversation with me?
		tc.client.Logger.Warn().Any("data", evt).Msg("Unhandled conversation create event")
		return true
	case *types.MessageDelete:
		allSuccess := true
		for _, deletedMsg := range evt.Messages {
			messageDeleteRemoteEvent := &simplevent.MessageRemove{
				EventMeta: simplevent.EventMeta{
					Type:      bridgev2.RemoteEventMessageRemove,
					PortalKey: tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
					LogContext: func(c zerolog.Context) zerolog.Context {
						return c.
							Str("message_id", deletedMsg.MessageID).
							Str("message_create_event_id", deletedMsg.MessageCreateEventID)
					},
					Timestamp:   methods.ParseSnowflake(evt.ID),
					StreamOrder: methods.ParseSnowflakeInt(evt.ID),
				},
				TargetMessage: networkid.MessageID(deletedMsg.MessageID),
			}
			allSuccess = tc.userLogin.QueueRemoteEvent(messageDeleteRemoteEvent).Success && allSuccess
		}
		return allSuccess
	case *types.ConversationNameUpdate:
		portalUpdateRemoteEvent := &simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type:   bridgev2.RemoteEventChatInfoChange,
				Sender: tc.MakeEventSender(evt.ByUserID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID).
						Str("new_name", evt.ConversationName).
						Str("changed_by_user_id", evt.ByUserID)
				},
				PortalKey:   tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				Timestamp:   methods.ParseSnowflake(evt.ID),
				StreamOrder: methods.ParseSnowflakeInt(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: &bridgev2.ChatInfo{
					Name: &evt.ConversationName,
				},
			},
		}
		return tc.userLogin.QueueRemoteEvent(portalUpdateRemoteEvent).Success
	case *types.ConversationAvatarUpdate:
		chatInfo := &bridgev2.ChatInfo{
			Avatar: makeAvatar(tc.client, evt.ConversationAvatarImageHttps),
		}
		success := tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type:        bridgev2.RemoteEventChatInfoChange,
				Sender:      tc.MakeEventSender(evt.ByUserID),
				PortalKey:   tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				StreamOrder: methods.ParseSnowflakeInt(evt.ID),
				Timestamp:   methods.ParseSnowflake(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: chatInfo,
			},
		}).Success
		return success
	case *types.ConversationMetadataUpdate:
		tc.client.Logger.Warn().Any("data", evt).Msg("Unhandled conversation metadata update event")
		return true
	case *types.ConversationJoin:
		// TODO handle
		return true
	case *types.ParticipantsJoin:
		conversation := inbox.GetConversationByID(evt.ConversationID)
		portalMembersAddedRemoteEvent := &simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatInfoChange,
				LogContext: func(c zerolog.Context) zerolog.Context {
					if conversation != nil {
						c = c.
							Bool("conv_low_quality", conversation.LowQuality).
							Bool("conv_trusted", conversation.Trusted)
					} else {
						c = c.Bool("conversation_nil", true)
					}
					return c.
						Str("conversation_id", evt.ConversationID).
						Int("total_new_members", len(evt.Participants))
				},
				PortalKey:    tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				CreatePortal: conversation != nil && conversation.Trusted,
				StreamOrder:  methods.ParseSnowflakeInt(evt.ID),
				Timestamp:    methods.ParseSnowflake(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				MemberChanges: tc.participantsToMemberList(evt.Participants, inbox),
			},
		}
		portalMembersAddedRemoteEvent.ChatInfoChange.MemberChanges.IsFull = false
		if conversation != nil {
			portalMembersAddedRemoteEvent.ChatInfoChange.MemberChanges.TotalMemberCount = len(conversation.Participants)
		}
		return tc.userLogin.QueueRemoteEvent(portalMembersAddedRemoteEvent).Success
	case *types.ParticipantsLeave:
		memberChanges := tc.participantsToMemberList(evt.Participants, inbox)
		for _, member := range memberChanges.MemberMap {
			member.Membership = event.MembershipLeave
		}
		conversation := inbox.GetConversationByID(evt.ConversationID)
		if conversation != nil {
			conversation.Participants = slices.DeleteFunc(conversation.Participants, func(pcp types.Participant) bool {
				_, remove := memberChanges.MemberMap[MakeUserID(pcp.UserID)]
				return remove
			})
			memberChanges.TotalMemberCount = len(conversation.Participants)
		}
		memberChanges.IsFull = false
		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatInfoChange,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID).
						Int("total_left_members", len(evt.Participants))
				},
				PortalKey:    tc.makePortalKeyFromInbox(evt.ConversationID, inbox),
				CreatePortal: false,
				StreamOrder:  methods.ParseSnowflakeInt(evt.ID),
				Timestamp:    methods.ParseSnowflake(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				MemberChanges: memberChanges,
			},
		}).Success
	case *types.ConversationDelete:
		portalDeleteRemoteEvent := &simplevent.ChatDelete{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventChatDelete,
				PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID)
				},
				StreamOrder: methods.ParseSnowflakeInt(evt.ID),
				Timestamp:   methods.ParseSnowflake(evt.ID),
			},
			OnlyForMe: true,
		}
		tc.client.Logger.Info().Any("data", evt).Msg("Deleted conversation")
		return tc.userLogin.QueueRemoteEvent(portalDeleteRemoteEvent).Success
	case *types.TrustConversation:
		conversation := inbox.GetConversationByID(evt.ConversationID)
		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatResync{
			EventMeta: simplevent.EventMeta{
				Type:         bridgev2.RemoteEventChatResync,
				PortalKey:    tc.MakePortalKey(conversation),
				CreatePortal: conversation != nil && conversation.Trusted,
			},
			ChatInfo: tc.conversationToChatInfo(conversation, inbox),
		}).Success
	case *types.EndAVBroadcast:
		conversation := inbox.GetConversationByID(evt.ConversationID)
		return tc.userLogin.QueueRemoteEvent(&simplevent.Message[string]{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventMessage,
				PortalKey: tc.MakePortalKey(conversation),
				Timestamp: methods.ParseSnowflake(evt.ID),
			},
			ID:   networkid.MessageID(evt.ID),
			Data: evt.CallType,
			ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, callType string) (*bridgev2.ConvertedMessage, error) {
				body := "Video"
				if callType == "AUDIO_ONLY" {
					body = "Audio"
				}
				if evt.EndReason == "HUNG_UP" {
					body += " call ended"
				} else if evt.IsCaller {
					body += " call"
				} else {
					body = "video"
					if callType == "AUDIO_ONLY" {
						body = "audio"
					}
					body = fmt.Sprintf("Missed %s call", body)
				}
				return &bridgev2.ConvertedMessage{
					Parts: []*bridgev2.ConvertedMessagePart{{
						Type: event.EventMessage,
						Content: &event.MessageEventContent{
							MsgType: event.MsgNotice,
							Body:    body,
						},
					}},
				}, nil
			},
		}).Success
	default:
		tc.client.Logger.Warn().
			Type("event_data_type", rawEvt).
			Any("event_data", rawEvt).
			Msg("Received unhandled event case from twitter library")
		return true
	}
}

func (tc *TwitterClient) wrapReaction(data *types.MessageReaction, portalKey networkid.PortalKey, evtType bridgev2.RemoteEventType) *simplevent.Reaction {
	return &simplevent.Reaction{
		EventMeta: simplevent.EventMeta{
			Type: evtType,
			LogContext: func(c zerolog.Context) zerolog.Context {
				return c.
					Str("message_id", data.MessageID).
					Str("sender", data.SenderID).
					Str("reaction_key", data.ReactionKey).
					Str("emoji_reaction", data.EmojiReaction)
			},
			PortalKey:   portalKey,
			Sender:      tc.MakeEventSender(data.SenderID),
			Timestamp:   methods.ParseSnowflake(data.ID),
			StreamOrder: methods.ParseSnowflakeInt(data.ID),
		},
		EmojiID:       "",
		Emoji:         data.EmojiReaction,
		TargetMessage: networkid.MessageID(data.MessageID),
	}
}

func (tc *TwitterClient) updateTwitterReadReceipt(inbox *response.TwitterInboxData) {
	for conversationID, conversation := range inbox.Conversations {
		cache := tc.participantCache[conversationID]
		for _, participant := range conversation.Participants {
			if participant.UserID == ParseUserLoginID(tc.userLogin.ID) {
				continue
			}
			if participant.LastReadEventID == "" {
				continue
			}
			var cachedParticipant *types.Participant
			for _, p := range cache {
				if p.UserID == participant.UserID {
					cachedParticipant = &p
					break
				}
			}
			if cachedParticipant == nil || cachedParticipant.LastReadEventID < participant.LastReadEventID {
				lastTarget := networkid.MessageID(participant.LastReadEventID)
				readUpTo := methods.ParseSnowflake(participant.LastReadEventID)
				readUpToStreamOrder := methods.ParseSnowflakeInt(participant.LastReadEventID)
				var targets []networkid.MessageID
				if lastTarget != "" {
					targets = []networkid.MessageID{lastTarget}
				}
				tc.userLogin.QueueRemoteEvent(&simplevent.Receipt{
					EventMeta: simplevent.EventMeta{
						Type:      bridgev2.RemoteEventReadReceipt,
						PortalKey: tc.makePortalKeyFromInbox(conversationID, inbox),
						Sender:    tc.MakeEventSender(participant.UserID),
						Timestamp: methods.ParseSnowflake(participant.LastReadEventID),
						PreHandleFunc: func(ctx context.Context, portal *bridgev2.Portal) {
							if intent := tc.userLogin.User.DoublePuppet(ctx); intent != nil {
								_ = intent.EnsureJoined(ctx, portal.MXID)
							}
						},
					},
					LastTarget:          lastTarget,
					Targets:             targets,
					ReadUpTo:            readUpTo,
					ReadUpToStreamOrder: readUpToStreamOrder,
				})
			}
		}
		tc.participantCache[conversationID] = conversation.Participants
	}
}

func (tc *TwitterClient) updateTwitterUserInfo(inbox *response.TwitterInboxData) {
	ctx := tc.userLogin.Log.With().Str("action", "update user info").Logger().WithContext(context.Background())
	for userID, user := range inbox.Users {
		cached := tc.userCache[userID]
		if cached == nil || cached.Name != user.Name || cached.ScreenName != user.ScreenName || cached.ProfileImageURLHTTPS != user.ProfileImageURLHTTPS {
			ghost, err := tc.connector.br.GetGhostByID(ctx, MakeUserID(userID))
			if err != nil {
				zerolog.Ctx(ctx).Err(err).Msg("Failed to get ghost by ID")
			} else {
				ghost.UpdateInfo(ctx, tc.connector.wrapUserInfo(tc.client, user))
			}
		}
	}
}

func (tc *TwitterClient) HandleStreamEvent(evt response.StreamEvent) {
	updateData := evt.Payload.DmUpdate
	typingData := evt.Payload.DmTyping

	if updateData != nil {
		tc.client.PollConversation(updateData.ConversationID)
	}

	if typingData != nil {
		tc.userLogin.QueueRemoteEvent(&simplevent.Typing{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventTyping,
				PortalKey: tc.MakePortalKeyFromID(typingData.ConversationID),
				Sender:    tc.MakeEventSender(typingData.UserID),
			},
			Timeout: 3 * time.Second,
		})
	}
}

// HandleXChatEvent handles events from the XChat websocket processor.
func (tc *TwitterClient) HandleXChatEvent(ctx context.Context, rawEvt types.TwitterEvent) bool {
	if rawEvt == nil {
		return true
	}

	log := tc.userLogin.Log.With().
		Str("handler", "xchat").
		Type("event_type", rawEvt).
		Logger()

	switch evt := rawEvt.(type) {
	case *types.Message:
		isFromMe := evt.MessageData.SenderID == string(tc.userLogin.ID)
		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)

		return tc.userLogin.QueueRemoteEvent(&simplevent.Message[*types.MessageData]{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventMessage,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("message_id", evt.ID).
						Str("sender", evt.MessageData.SenderID).
						Bool("is_from_me", isFromMe)
				},
				PortalKey:    portalKey,
				CreatePortal: false, // Portal should already exist from initial sync
				Sender:       tc.MakeEventSender(evt.MessageData.SenderID),
				StreamOrder:  methods.ParseInt64(evt.ID),
				Timestamp:    methods.ParseMsecTimestamp(evt.Time),
			},
			ID:            networkid.MessageID(evt.ID),
			TransactionID: networkid.TransactionID(evt.RequestID),
			TargetMessage: networkid.MessageID(evt.ID),
			Data:          &evt.MessageData,
			ConvertMessageFunc: func(ctx context.Context, portal *bridgev2.Portal, intent bridgev2.MatrixAPI, data *types.MessageData) (*bridgev2.ConvertedMessage, error) {
				return tc.convertToMatrix(ctx, portal, intent, data), nil
			},
			ConvertEditFunc: tc.convertEditToMatrix,
		}).Success

	case *types.MessageReactionCreate:
		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
		wrappedEvt := tc.wrapReaction((*types.MessageReaction)(evt), portalKey, bridgev2.RemoteEventReaction)
		return tc.userLogin.QueueRemoteEvent(wrappedEvt).Success

	case *types.MessageReactionDelete:
		portalKey := tc.MakePortalKeyFromID(evt.ConversationID)
		wrappedEvt := tc.wrapReaction((*types.MessageReaction)(evt), portalKey, bridgev2.RemoteEventReactionRemove)
		return tc.userLogin.QueueRemoteEvent(wrappedEvt).Success

	case *types.ConversationRead:
		lastTarget := networkid.MessageID(evt.LastReadEventID)
		readUpTo := methods.ParseMsecTimestamp(evt.Time)
		readUpToStreamOrder := methods.ParseInt64(evt.LastReadEventID)
		if readUpToStreamOrder == 0 {
			readUpToStreamOrder = methods.ParseInt64(evt.ID)
		}
		var targets []networkid.MessageID
		if lastTarget != "" {
			targets = []networkid.MessageID{lastTarget}
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.Receipt{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventReadReceipt,
				PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
				Sender:    tc.MakeEventSender(ParseUserLoginID(tc.userLogin.ID)),
				Timestamp: readUpTo,
				PreHandleFunc: func(ctx context.Context, portal *bridgev2.Portal) {
					if intent := tc.userLogin.User.DoublePuppet(ctx); intent != nil {
						_ = intent.EnsureJoined(ctx, portal.MXID)
					}
				},
			},
			LastTarget:          lastTarget,
			Targets:             targets,
			ReadUpTo:            readUpTo,
			ReadUpToStreamOrder: readUpToStreamOrder,
		}).Success

	case *types.MessageDelete:
		allSuccess := true
		for _, deletedMsg := range evt.Messages {
			messageDeleteRemoteEvent := &simplevent.MessageRemove{
				EventMeta: simplevent.EventMeta{
					Type:      bridgev2.RemoteEventMessageRemove,
					PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
					LogContext: func(c zerolog.Context) zerolog.Context {
						return c.
							Str("message_id", deletedMsg.MessageID).
							Str("message_create_event_id", deletedMsg.MessageCreateEventID)
					},
					Timestamp:   methods.ParseMsecTimestamp(evt.Time),
					StreamOrder: methods.ParseInt64(evt.ID),
				},
				TargetMessage: networkid.MessageID(deletedMsg.MessageID),
			}
			allSuccess = tc.userLogin.QueueRemoteEvent(messageDeleteRemoteEvent).Success && allSuccess
		}
		return allSuccess

	case *types.ConversationDelete:
		portalDeleteRemoteEvent := &simplevent.ChatDelete{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventChatDelete,
				PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.Str("conversation_id", evt.ConversationID)
				},
				StreamOrder: methods.ParseInt64(evt.ID),
				Timestamp:   methods.ParseMsecTimestamp(evt.Time),
			},
			OnlyForMe: true,
		}
		log.Info().Any("data", evt).Msg("Deleted conversation")
		return tc.userLogin.QueueRemoteEvent(portalDeleteRemoteEvent).Success

	case *types.ConversationNameUpdate:
		portalUpdateRemoteEvent := &simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type:   bridgev2.RemoteEventChatInfoChange,
				Sender: tc.MakeEventSender(evt.ByUserID),
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID).
						Str("new_name", evt.ConversationName).
						Str("changed_by_user_id", evt.ByUserID)
				},
				PortalKey:   tc.MakePortalKeyFromID(evt.ConversationID),
				Timestamp:   methods.ParseMsecTimestamp(evt.Time),
				StreamOrder: methods.ParseInt64(evt.ID),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: &bridgev2.ChatInfo{
					Name: &evt.ConversationName,
				},
			},
		}
		return tc.userLogin.QueueRemoteEvent(portalUpdateRemoteEvent).Success

	case *types.ConversationAvatarUpdate:
		chatInfo := &bridgev2.ChatInfo{
			Avatar: makeAvatar(tc.client, evt.ConversationAvatarImageHttps),
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type:        bridgev2.RemoteEventChatInfoChange,
				Sender:      tc.MakeEventSender(evt.ByUserID),
				PortalKey:   tc.MakePortalKeyFromID(evt.ConversationID),
				StreamOrder: methods.ParseInt64(evt.ID),
				Timestamp:   methods.ParseMsecTimestamp(evt.Time),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				ChatInfo: chatInfo,
			},
		}).Success

	case *types.ParticipantsJoin:
		memberChanges := &bridgev2.ChatMemberList{
			IsFull:    false,
			MemberMap: make(map[networkid.UserID]bridgev2.ChatMember),
		}
		for _, p := range evt.Participants {
			memberChanges.MemberMap[MakeUserID(p.UserID)] = bridgev2.ChatMember{
				EventSender: tc.MakeEventSender(p.UserID),
				Membership:  event.MembershipJoin,
			}
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatInfoChange,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID).
						Int("total_new_members", len(evt.Participants))
				},
				PortalKey:    tc.MakePortalKeyFromID(evt.ConversationID),
				CreatePortal: false, // Portal should already exist from initial sync
				StreamOrder:  methods.ParseInt64(evt.ID),
				Timestamp:    methods.ParseMsecTimestamp(evt.Time),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				MemberChanges: memberChanges,
			},
		}).Success

	case *types.ParticipantsLeave:
		memberChanges := &bridgev2.ChatMemberList{
			IsFull:    false,
			MemberMap: make(map[networkid.UserID]bridgev2.ChatMember),
		}
		for _, p := range evt.Participants {
			memberChanges.MemberMap[MakeUserID(p.UserID)] = bridgev2.ChatMember{
				EventSender: tc.MakeEventSender(p.UserID),
				Membership:  event.MembershipLeave,
			}
		}
		return tc.userLogin.QueueRemoteEvent(&simplevent.ChatInfoChange{
			EventMeta: simplevent.EventMeta{
				Type: bridgev2.RemoteEventChatInfoChange,
				LogContext: func(c zerolog.Context) zerolog.Context {
					return c.
						Str("conversation_id", evt.ConversationID).
						Int("total_left_members", len(evt.Participants))
				},
				PortalKey:    tc.MakePortalKeyFromID(evt.ConversationID),
				CreatePortal: false,
				StreamOrder:  methods.ParseInt64(evt.ID),
				Timestamp:    methods.ParseMsecTimestamp(evt.Time),
			},
			ChatInfoChange: &bridgev2.ChatInfoChange{
				MemberChanges: memberChanges,
			},
		}).Success

	case *types.XChatTyping:
		tc.userLogin.QueueRemoteEvent(&simplevent.Typing{
			EventMeta: simplevent.EventMeta{
				Type:      bridgev2.RemoteEventTyping,
				PortalKey: tc.MakePortalKeyFromID(evt.ConversationID),
				Sender:    tc.MakeEventSender(evt.SenderID),
			},
			Timeout: 3 * time.Second,
		})
		return true

	case *types.XChatKeyChange:
		log.Info().
			Str("conversation_id", evt.ConversationID).
			Str("new_key_version", evt.NewKeyVersion).
			Msg("Conversation key changed")
		return true

	case *types.XChatMessageFailure:
		log.Warn().
			Str("conversation_id", evt.ConversationID).
			Str("message_id", evt.MessageID).
			Int32("failure_type", int32(evt.FailureType)).
			Msg("Message failure event received")
		return true

	default:
		log.Debug().
			Type("event_data_type", rawEvt).
			Any("event_data", rawEvt).
			Msg("Received unhandled XChat event")
		return true
	}
}

// ignorePayloadImport is used to prevent the import from being removed
var _ = payload.FailureType(0)
