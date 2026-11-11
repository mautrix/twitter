package twittermeow

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.mau.fi/util/ptr"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

var (
	ErrGroupInviteLinkEnabled = errors.New("changing members isn't supported while the group invite link is on")
	ErrGroupHasPins           = errors.New("changing members isn't supported while the group has pinned messages")
	ErrGroupHasDescription    = errors.New("changing members isn't supported while the group has a description")
)

// xchatGroupState is the group state that member changes re-sign and re-encrypt under a new key.
type xchatGroupState struct {
	memberIDs     []string
	adminIDs      []string
	createdBy     string
	title         []byte
	hasTitle      bool
	avatar        []byte
	ttlMsec       *int64
	screenBlocked bool
	adminSettings *payload.GroupMemberAdminSettings
	oldKey        *crypto.ConversationKey
}

// rotatedGroupKey is a fresh conversation key with the title and avatar re-encrypted under it.
type rotatedGroupKey struct {
	key             *crypto.ConversationKey
	participantKeys []*payload.ConversationParticipantKey
	encTitle        *string
	avatarURL       *string
}

func (c *Client) fetchXChatGroupState(ctx context.Context, conversationID string) (*xchatGroupState, error) {
	resp, err := c.GetConversationData(ctx, payload.NewInboxPageConversationDataQueryVariables(conversationID, false))
	if err != nil {
		return nil, fmt.Errorf("fetch group state: %w", err)
	}
	item, err := resp.Item(conversationID)
	if err != nil {
		return nil, err
	}
	detail := item.ConversationDetail
	meta := detail.GroupMetadata
	if meta == nil || len(detail.GroupMembersResults) == 0 || len(detail.GroupAdminsResults) == 0 {
		return nil, errors.New("group state response is missing metadata, members or admins")
	} else if meta.GroupInviteDetails != nil && meta.GroupInviteDetails.InviteURL != "" {
		return nil, ErrGroupInviteLinkEnabled
	} else if len(item.PMKBundle) > 0 || len(item.LatestPinnedMessageEvents) > 0 {
		// x.com re-wraps the pinned message keys under the new conversation key, which isn't implemented here.
		return nil, ErrGroupHasPins
	} else if detail.LatestGroupDescriptionChangeEvent != "" {
		return nil, ErrGroupHasDescription
	}
	oldKey, err := c.keyManager.GetLatestConversationKey(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("get conversation key: %w", err)
	} else if oldKey == nil || len(oldKey.Key) == 0 {
		return nil, crypto.ErrKeyNotFound
	}
	state := &xchatGroupState{
		screenBlocked: meta.ScreenCaptureBlockingEnabled,
		oldKey:        oldKey,
		adminSettings: &payload.GroupMemberAdminSettings{
			EditGroupInfo:      ptr.Ptr(false),
			EditMessageTtl:     ptr.Ptr(false),
			BlockScreenCapture: ptr.Ptr(false),
			AddMember:          ptr.Ptr(false),
			SendMessage:        ptr.Ptr(false),
			StartCall:          ptr.Ptr(false),
		},
	}
	for _, r := range detail.GroupMembersResults {
		state.memberIDs = append(state.memberIDs, r.RestID)
	}
	for _, r := range detail.GroupAdminsResults {
		state.adminIDs = append(state.adminIDs, r.RestID)
	}
	if detail.GroupCreatedByUserResults != nil {
		state.createdBy = detail.GroupCreatedByUserResults.RestID
	}
	if settings := meta.AdminSettings; settings != nil {
		state.adminSettings = &payload.GroupMemberAdminSettings{
			EditGroupInfo:      ptr.Ptr(settings.EditGroupInfo),
			EditMessageTtl:     ptr.Ptr(settings.EditMessageTTL),
			BlockScreenCapture: ptr.Ptr(settings.BlockScreenCapture),
			AddMember:          ptr.Ptr(settings.AddMember),
			SendMessage:        ptr.Ptr(settings.SendMessage),
			StartCall:          ptr.Ptr(settings.StartCall),
		}
	}
	if meta.MessageTTLMsec != "" {
		ttl, err := strconv.ParseInt(meta.MessageTTLMsec, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse group message TTL: %w", err)
		}
		state.ttlMsec = &ttl
	}
	if meta.GroupName != "" {
		encrypted, err := base64.StdEncoding.DecodeString(meta.GroupName)
		if err != nil {
			return nil, fmt.Errorf("decode group title: %w", err)
		}
		state.title, err = crypto.SecretboxDecrypt(encrypted, oldKey.Key)
		if err != nil {
			return nil, fmt.Errorf("decrypt group title: %w", err)
		}
		state.hasTitle = true
	}
	if meta.GroupAvatarURL != "" {
		httpResp, err := c.FetchRaw(ctx, meta.GroupAvatarURL)
		if err != nil {
			return nil, fmt.Errorf("download group avatar: %w", err)
		}
		defer httpResp.Body.Close()
		var buf bytes.Buffer
		if _, err = crypto.SecretstreamDecrypt(httpResp.Body, oldKey.Key, &buf); err != nil {
			return nil, fmt.Errorf("decrypt group avatar: %w", err)
		}
		state.avatar = buf.Bytes()
	}
	return state, nil
}

// rotateXChatGroupKey makes a new conversation key, wraps it for every final member and
// re-encrypts the title and avatar with it, the same way x.com does for member changes.
func (c *Client) rotateXChatGroupKey(ctx context.Context, conversationID string, state *xchatGroupState, finalMembers []string) (*rotatedGroupKey, error) {
	keysResp, err := c.GetPublicKeys(ctx, finalMembers)
	if err != nil {
		return nil, fmt.Errorf("fetch member public keys: %w", err)
	}
	latest := make(map[string]response.XChatPublicKeyWithMeta, len(finalMembers))
	for _, user := range keysResp.Data.UserResultsByRestIDs {
		for _, k := range user.Result.GetPublicKeys.PublicKeysWithTokenMap {
			current, ok := latest[user.RestID]
			if !ok || compareNumericStrings(k.PublicKeyWithMetadata.Version, current.Version) > 0 {
				latest[user.RestID] = k.PublicKeyWithMetadata
			}
		}
	}

	newKey := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, newKey); err != nil {
		return nil, err
	}
	version := time.Now().UnixMilli()
	if oldVersion, err := strconv.ParseInt(state.oldKey.KeyVersion, 10, 64); err == nil && version <= oldVersion {
		version = oldVersion + 1
	}
	rotated := &rotatedGroupKey{key: &crypto.ConversationKey{
		ConversationID: conversationID,
		KeyVersion:     strconv.FormatInt(version, 10),
		Key:            newKey,
		CreatedAt:      time.Now(),
	}}
	for _, userID := range finalMembers {
		pk, ok := latest[userID]
		if !ok || pk.PublicKey.PublicKey == "" {
			return nil, fmt.Errorf("no X Chat public key for member %s", userID)
		}
		wrapped, err := crypto.WrapConversationKey(newKey, pk.PublicKey.PublicKey)
		if err != nil {
			return nil, fmt.Errorf("wrap key for %s: %w", userID, err)
		}
		rotated.participantKeys = append(rotated.participantKeys, &payload.ConversationParticipantKey{
			UserId:                   ptr.Ptr(userID),
			EncryptedConversationKey: ptr.Ptr(wrapped),
			PublicKeyVersion:         ptr.Ptr(pk.Version),
		})
	}
	if state.hasTitle {
		encrypted, err := crypto.SecretboxEncrypt(state.title, newKey)
		if err != nil {
			return nil, fmt.Errorf("encrypt group title: %w", err)
		}
		rotated.encTitle = ptr.Ptr(base64.StdEncoding.EncodeToString(encrypted))
	}
	if state.avatar != nil {
		uploaded, err := c.uploadXChatMediaWithKey(ctx, conversationID, fmt.Sprintf("avatar-%s-%s", conversationID, uuid.NewString()), state.avatar, rotated.key)
		if err != nil {
			return nil, fmt.Errorf("re-upload group avatar: %w", err)
		}
		rotated.avatarURL = ptr.Ptr(fmt.Sprintf("%s/%s/%s", endpoints.XCHAT_MEDIA_BASE_URL, conversationID, uploaded.MediaHashKey))
	}
	return rotated, nil
}

func (r *rotatedGroupKey) variableKeys() []payload.GroupParticipantKey {
	keys := make([]payload.GroupParticipantKey, len(r.participantKeys))
	for i, k := range r.participantKeys {
		keys[i] = payload.GroupParticipantKey{
			EncryptedConversationKey: *k.EncryptedConversationKey,
			PublicKeyVersion:         *k.PublicKeyVersion,
			UserID:                   *k.UserId,
		}
	}
	return keys
}

func (c *Client) signXChatKeyChange(ctx context.Context, conversationID string, rotated *rotatedGroupKey, forKeyRotation bool) (*payload.DeleteMessageActionSignature, error) {
	return c.signXChatEvent(ctx, &payload.MessageEventDetail{
		ConversationKeyChangeEvent: &payload.ConversationKeyChangeEvent{
			ConversationKeyVersion:      &rotated.key.KeyVersion,
			ConversationParticipantKeys: rotated.participantKeys,
			ForKeyRotation:              &forKeyRotation,
		},
	}, "ConversationKeyChangeEvent", conversationID,
		rotated.key.KeyVersion, base64.RawStdEncoding.EncodeToString(rotated.key.Key))
}

// AddXChatGroupMembers adds users to an X Chat group, rotating the group key.
func (c *Client) AddXChatGroupMembers(ctx context.Context, conversationID string, userIDs []string) error {
	state, err := c.fetchXChatGroupState(ctx, conversationID)
	if err != nil {
		return err
	}
	var newMembers []string
	for _, userID := range userIDs {
		if !slices.Contains(state.memberIDs, userID) && !slices.Contains(newMembers, userID) {
			newMembers = append(newMembers, userID)
		}
	}
	if len(newMembers) == 0 {
		return nil
	}
	rotated, err := c.rotateXChatGroupKey(ctx, conversationID, state, append(slices.Clone(newMembers), state.memberIDs...))
	if err != nil {
		return err
	}

	add := &payload.GroupMemberAddChange{
		MemberIds:                    newMembers,
		CurrentMemberIds:             state.memberIDs,
		CurrentAdminIds:              state.adminIDs,
		CurrentTitle:                 rotated.encTitle,
		CurrentAvatarUrl:             rotated.avatarURL,
		ConversationKeyVersion:       &rotated.key.KeyVersion,
		CurrentTtlMsec:               state.ttlMsec,
		ScreenCaptureBlockingEnabled: &state.screenBlocked,
		AdminSettings:                state.adminSettings,
	}
	if state.createdBy != "" {
		add.GroupCreatedByUserId = &state.createdBy
	}
	ttlString := "null"
	if state.ttlMsec != nil {
		ttlString = strconv.FormatInt(*state.ttlMsec, 10)
	}
	signatureFields := slices.Concat(newMembers, state.memberIDs, state.adminIDs, []string{
		stringOrNull(rotated.encTitle), stringOrNull(rotated.avatarURL), ttlString, rotated.key.KeyVersion, strconv.FormatBool(state.screenBlocked),
	})
	addSig, err := c.signXChatEvent(ctx, &payload.MessageEventDetail{
		GroupChangeEvent: &payload.GroupChangeEvent{GroupChange: &payload.GroupChange{GroupMemberAdd: add}},
	}, "GroupChangeEvent.GroupMemberAddChange", conversationID, signatureFields...)
	if err != nil {
		return err
	}
	keySig, err := c.signXChatKeyChange(ctx, conversationID, rotated, false)
	if err != nil {
		return err
	}
	vars := payload.AddGroupMembersMutationVariables{
		ConversationID:               conversationID,
		UserIDs:                      newMembers,
		ConversationKeyVersion:       rotated.key.KeyVersion,
		ConversationParticipantKeys:  rotated.variableKeys(),
		EncryptedAvatarURL:           rotated.avatarURL,
		EncryptedTitle:               rotated.encTitle,
		ActionSignatures:             []payload.DeleteMessageActionSignature{*addSig, *keySig},
		ScreenCaptureBlockingEnabled: state.screenBlocked,
	}
	if state.ttlMsec != nil {
		vars.TTLMsec = &ttlString
	}
	respBody, err := c.postXChatGroupMutation(ctx, endpoints.ADD_GROUP_MEMBERS_URL,
		payload.NewGroupMutationPayload("AddGroupMembersMutation", "ygL_9juTllEGdCTCZNnFqw", vars))
	if err != nil {
		return err
	} else if err = parseXChatGroupMutationResponse(respBody, "xchat_add_group_members"); err != nil {
		return err
	}
	return c.keyManager.PutConversationKey(ctx, rotated.key)
}

// RemoveXChatGroupMembers removes users from an X Chat group, rotating the group key.
func (c *Client) RemoveXChatGroupMembers(ctx context.Context, conversationID string, userIDs []string) error {
	state, err := c.fetchXChatGroupState(ctx, conversationID)
	if err != nil {
		return err
	}
	var removed, remaining []string
	for _, memberID := range state.memberIDs {
		if slices.Contains(userIDs, memberID) {
			removed = append(removed, memberID)
		} else {
			remaining = append(remaining, memberID)
		}
	}
	if len(removed) == 0 {
		return nil
	} else if len(remaining) == 0 {
		return errors.New("can't remove every member of a group")
	}
	rotated, err := c.rotateXChatGroupKey(ctx, conversationID, state, remaining)
	if err != nil {
		return err
	}

	sigs := make([]payload.DeleteMessageActionSignature, 0, 4)
	removeSig, err := c.signXChatEvent(ctx, &payload.MessageEventDetail{
		GroupChangeEvent: &payload.GroupChangeEvent{GroupChange: &payload.GroupChange{
			GroupMemberRemove: &payload.GroupMemberRemoveChange{MemberIds: removed, ConversationKeyVersion: &rotated.key.KeyVersion},
		}},
	}, "GroupChangeEvent.GroupMemberRemoveChange", conversationID, removed...)
	if err != nil {
		return err
	}
	sigs = append(sigs, *removeSig)
	if rotated.encTitle != nil {
		titleSig, err := c.signXChatEvent(ctx, &payload.MessageEventDetail{
			GroupChangeEvent: &payload.GroupChangeEvent{
				GroupChange: &payload.GroupChange{GroupTitleChange: &payload.GroupTitleChange{
					CustomTitle: rotated.encTitle, ConversationKeyVersion: &rotated.key.KeyVersion,
				}},
				ForKeyRotation: ptr.Ptr(true),
			},
		}, "GroupChangeEvent.GroupTitleChange", conversationID, rotated.key.KeyVersion, *rotated.encTitle)
		if err != nil {
			return err
		}
		sigs = append(sigs, *titleSig)
	}
	if rotated.avatarURL != nil {
		avatarSig, err := c.signXChatEvent(ctx, &payload.MessageEventDetail{
			GroupChangeEvent: &payload.GroupChangeEvent{
				GroupChange: &payload.GroupChange{GroupAvatarChange: &payload.GroupAvatarUrlChange{
					CustomAvatarUrl: rotated.avatarURL, ConversationKeyVersion: &rotated.key.KeyVersion,
				}},
				ForKeyRotation: ptr.Ptr(true),
			},
		}, "GroupChangeEvent.GroupAvatarUrlChange", conversationID, rotated.key.KeyVersion, *rotated.avatarURL)
		if err != nil {
			return err
		}
		sigs = append(sigs, *avatarSig)
	}
	keySig, err := c.signXChatKeyChange(ctx, conversationID, rotated, true)
	if err != nil {
		return err
	}
	sigs = append(sigs, *keySig)

	vars := payload.RemoveFromGroupMutationVariables{
		ConversationID:   conversationID,
		UserIDs:          removed,
		ActionSignatures: sigs,
		GroupKeyRotationData: payload.GroupKeyRotationData{
			ConversationKeyVersion:      rotated.key.KeyVersion,
			ConversationParticipantKeys: rotated.variableKeys(),
			EncryptedAvatar:             rotated.avatarURL,
			EncryptedTitle:              rotated.encTitle,
		},
	}
	respBody, err := c.postXChatGroupMutation(ctx, endpoints.REMOVE_GROUP_MEMBERS_URL,
		payload.NewGroupMutationPayload("RemoveFromGroupMutation", "Lnm2b8oz4dnV8V-RtP1SYA", vars))
	if err != nil {
		return err
	} else if err = parseXChatGroupMutationResponse(respBody, "xchat_remove_group_members"); err != nil {
		return err
	}
	return c.keyManager.PutConversationKey(ctx, rotated.key)
}

// stringOrNull formats optional signature fields the way the web client joins them.
func stringOrNull(s *string) string {
	if s == nil {
		return "null"
	}
	return *s
}

func compareNumericStrings(a, b string) int {
	if len(a) != len(b) {
		return len(a) - len(b)
	}
	if a < b {
		return -1
	} else if a > b {
		return 1
	}
	return 0
}
