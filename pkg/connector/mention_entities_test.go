package connector

import "testing"

func TestBuildMentionEntities(t *testing.T) {
	tests := []struct {
		name string
		text string
		want [][2]int32
	}{
		{
			name: "start",
			text: "@zhoro_x hi",
			want: [][2]int32{{0, 8}},
		},
		{
			name: "preceded by text",
			text: "hey @zhoro_x!",
			want: [][2]int32{{4, 12}},
		},
		{
			name: "unicode before mention",
			text: "hi \u263a @zhoro_x",
			want: [][2]int32{{5, 13}},
		},
		{
			name: "email ignored",
			text: "test@example.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			entities := buildMentionEntities(test.text)
			if len(entities) != len(test.want) {
				t.Fatalf("expected %d entities, got %d", len(test.want), len(entities))
			}
			for i, want := range test.want {
				if entities[i].StartIndex == nil || entities[i].EndIndex == nil {
					t.Fatalf("entity %d missing indices", i)
				}
				if *entities[i].StartIndex != want[0] || *entities[i].EndIndex != want[1] {
					t.Fatalf("entity %d indices = %d..%d, want %d..%d", i, *entities[i].StartIndex, *entities[i].EndIndex, want[0], want[1])
				}
				if entities[i].Content == nil || entities[i].Content.Mention == nil {
					t.Fatalf("entity %d missing mention content", i)
				}
			}
		})
	}
}
