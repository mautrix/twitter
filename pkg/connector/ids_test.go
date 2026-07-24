package connector

import (
	"testing"

	"maunium.net/go/mautrix/bridgev2/networkid"
)

func TestXChatGroupPortalAliasKey(t *testing.T) {
	loginID := networkid.UserLoginID("1155463061127467008")
	tests := []struct {
		name string
		key  networkid.PortalKey
		want networkid.PortalKey
		ok   bool
	}{
		{
			name: "REST group",
			key:  networkid.PortalKey{ID: "1580281760675504141", Receiver: loginID},
			want: networkid.PortalKey{ID: "g1580281760675504141", Receiver: loginID},
			ok:   true,
		},
		{
			name: "direct message",
			key:  networkid.PortalKey{ID: "1155463061127467008-1247940250015588353", Receiver: loginID},
		},
		{
			name: "XChat group",
			key:  networkid.PortalKey{ID: "g1580281760675504141", Receiver: loginID},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := xchatGroupPortalAliasKey(test.key)
			if ok != test.ok || got != test.want {
				t.Fatalf("xchatGroupPortalAliasKey() = (%#v, %t), want (%#v, %t)", got, ok, test.want, test.ok)
			}
		})
	}
}

func TestRESTGroupPortalAliasKey(t *testing.T) {
	loginID := networkid.UserLoginID("1155463061127467008")
	tests := []struct {
		name string
		key  networkid.PortalKey
		want networkid.PortalKey
		ok   bool
	}{
		{
			name: "xchat group",
			key:  networkid.PortalKey{ID: "g1580281760675504141", Receiver: loginID},
			want: networkid.PortalKey{ID: "1580281760675504141", Receiver: loginID},
			ok:   true,
		},
		{
			name: "rest group",
			key:  networkid.PortalKey{ID: "1580281760675504141", Receiver: loginID},
		},
		{
			name: "xchat direct message",
			key:  networkid.PortalKey{ID: "1155463061127467008-1247940250015588353", Receiver: loginID},
		},
		{
			name: "non-numeric g prefix",
			key:  networkid.PortalKey{ID: "group", Receiver: loginID},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := restGroupPortalAliasKey(test.key)
			if ok != test.ok || got != test.want {
				t.Fatalf("restGroupPortalAliasKey() = (%#v, %t), want (%#v, %t)", got, ok, test.want, test.ok)
			}
		})
	}
}
