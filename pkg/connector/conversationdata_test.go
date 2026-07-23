package connector

import "testing"

func TestConversationDataResultID(t *testing.T) {
	tests := []struct {
		name      string
		requested string
		returned  string
		want      string
	}{
		{
			name:      "missing response ID",
			requested: "g1580281760675504141",
			want:      "g1580281760675504141",
		},
		{
			name:      "bare REST alias",
			requested: "g1580281760675504141",
			returned:  "1580281760675504141",
			want:      "g1580281760675504141",
		},
		{
			name:      "matching XChat ID",
			requested: "g1580281760675504141",
			returned:  "g1580281760675504141",
			want:      "g1580281760675504141",
		},
		{
			name:      "different response ID",
			requested: "g1580281760675504141",
			returned:  "g999",
			want:      "g999",
		},
		{
			name:      "missing direct message ID",
			requested: "1155463061127467008:1247940250015588353",
			want:      "1155463061127467008:1247940250015588353",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := conversationDataResultID(test.requested, test.returned); got != test.want {
				t.Fatalf("conversationDataResultID() = %q, want %q", got, test.want)
			}
		})
	}
}
