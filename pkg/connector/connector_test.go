package connector

import (
	"testing"

	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/database"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func TestTwitterConnectorUsesCompletionAwarePortalHandling(t *testing.T) {
	previousBuffer := bridgev2.PortalEventBuffer
	t.Cleanup(func() {
		bridgev2.PortalEventBuffer = previousBuffer
	})
	bridgev2.PortalEventBuffer = 64

	connector := &TwitterConnector{}
	connector.Init(&bridgev2.Bridge{})
	if bridgev2.PortalEventBuffer != 0 {
		t.Fatalf("PortalEventBuffer = %d, want 0", bridgev2.PortalEventBuffer)
	}
	if xchatRemoteEventHandled(bridgev2.EventHandlingResultQueued) {
		t.Fatal("queued event was treated as completed")
	}
	if !xchatRemoteEventHandled(bridgev2.EventHandlingResultSuccess) {
		t.Fatal("successful completed event was rejected")
	}
	if !xchatRemoteEventHandled(bridgev2.EventHandlingResultIgnored) {
		t.Fatal("completed duplicate/ignored event was rejected")
	}
}

func TestXChatUserResultFallsBackToOuterID(t *testing.T) {
	userID, user := xchatUserFromResult(response.XChatUserResult{
		RestID: "123",
		Result: &response.XChatUser{
			Core: &response.XChatUserCore{Name: "Display Name", ScreenName: "username"},
		},
	})
	if userID != "123" || user == nil || user.IDStr != "123" {
		t.Fatalf("converted user ID/user = %q/%#v, want outer ID 123", userID, user)
	}

	client := &TwitterClient{userCache: make(map[string]*types.User)}
	missing := client.collectAndCacheUserResults([]response.XChatUserResult{{RestID: "missing"}})
	if len(missing) != 1 || missing[0] != "missing" {
		t.Fatalf("missing user IDs = %v, want [missing]", missing)
	}

	client.userCache["complete"] = &types.User{IDStr: "complete", ScreenName: "known", Name: "Known User"}
	missing = client.collectAndCacheUserResults([]response.XChatUserResult{{
		RestID: "complete",
		Result: &response.XChatUser{RestID: "complete"},
	}})
	if len(missing) != 0 || client.userCache["complete"].ScreenName != "known" {
		t.Fatalf("partial inline profile replaced complete cache: missing=%v cache=%#v", missing, client.userCache["complete"])
	}

	partial := response.XChatUserResult{RestID: "duplicate", Result: &response.XChatUser{RestID: "duplicate"}}
	complete := response.XChatUserResult{
		RestID: "duplicate",
		Result: &response.XChatUser{
			RestID: "duplicate",
			Core:   &response.XChatUserCore{ScreenName: "resolved"},
		},
	}
	preferred := preferXChatUserResult(partial, complete)
	_, preferredUser := xchatUserFromResult(preferred)
	if preferredUser == nil || preferredUser.ScreenName != "resolved" {
		t.Fatalf("preferred duplicate profile = %#v, want resolved profile", preferredUser)
	}
	users := map[string]*types.User{"duplicate": {IDStr: "duplicate", ScreenName: "resolved"}}
	mergeXChatUser(users, "duplicate", &types.User{IDStr: "duplicate"})
	if users["duplicate"].ScreenName != "resolved" {
		t.Fatalf("partial user replaced complete map entry: %#v", users["duplicate"])
	}
	partials := map[string]*types.User{"split": {IDStr: "split", ScreenName: "handle"}}
	merged := mergeXChatUser(partials, "split", &types.User{IDStr: "split", Name: "Display Name"})
	if merged.ScreenName != "handle" || merged.Name != "Display Name" {
		t.Fatalf("split profile metadata was not merged: %#v", merged)
	}
}

func TestXChatGroupWithoutMemberFieldIsNotAuthoritative(t *testing.T) {
	twitterClient := &twittermeow.Client{}
	twitterClient.SetCurrentUserID("self")
	client := &TwitterClient{
		connector: &TwitterConnector{},
		client:    twitterClient,
		userCache: make(map[string]*types.User),
	}
	item := &response.XChatInboxItem{
		ConversationDetail: response.XChatConversationDetail{
			ConversationID: "g123",
			GroupMetadata:  &response.XChatGroupMetadata{},
		},
	}
	info := client.xchatItemToChatInfo(t.Context(), item, nil, &types.Conversation{
		ConversationID: "g123",
		Type:           ConversationTypeGroupDM,
	})
	if info.Members == nil || info.Members.IsFull {
		t.Fatalf("group members = %#v, want non-authoritative partial list", info.Members)
	}
	item.ConversationDetail.GroupMembersResults = []response.XChatUserResult{{RestID: "123"}}
	info = client.xchatItemToChatInfo(t.Context(), item, nil, &types.Conversation{
		ConversationID: "g123",
		Type:           ConversationTypeGroupDM,
	})
	if info.Members == nil || !info.Members.IsFull {
		t.Fatalf("group members = %#v, want authoritative populated list", info.Members)
	}
}

func TestXChatGroupMemberProfilesTriggerInfoUpdate(t *testing.T) {
	info := &bridgev2.ChatInfo{
		Type: ptr.Ptr(database.RoomTypeDefault),
		Members: &bridgev2.ChatMemberList{
			MemberMap: bridgev2.ChatMemberMap{
				MakeUserID("123"): {EventSender: bridgev2.EventSender{Sender: MakeUserID("123")}},
			},
		},
	}
	if !shouldEmitChatInfoUpdate(info, database.RoomTypeDefault) {
		t.Fatal("group member profile update was skipped")
	}
}
