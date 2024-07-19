package response

import "go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"

type SearchResponse struct {
	NumResults      int          `json:"num_results,omitempty"`
	Users           []types.User `json:"users,omitempty"`
	Topics          []any        `json:"topics,omitempty"`
	Events          []any        `json:"events,omitempty"`
	Lists           []any        `json:"lists,omitempty"`
	OrderedSections []any        `json:"ordered_sections,omitempty"`
	Oneclick        []any        `json:"oneclick,omitempty"`
	Hashtags        []any        `json:"hashtags,omitempty"`
	CompletedIn     float32      `json:"completed_in,omitempty"`
	Query           string       `json:"query,omitempty"`
}
