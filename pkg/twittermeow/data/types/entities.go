package types

import "fmt"

type Attachment struct {
	Video       *AttachmentInfo  `json:"video,omitempty"`
	AnimatedGif *AttachmentInfo  `json:"animated_gif,omitempty"`
	Photo       *AttachmentInfo  `json:"photo,omitempty"`
	Card        *AttachmentCard  `json:"card,omitempty"`
	Tweet       *AttachmentTweet `json:"tweet,omitempty"`

	// XChat-specific fields for URL attachment images
	URLBannerMediaHashKey string `json:"url_banner_media_hash_key,omitempty"`
}
type URLs struct {
	URL         string `json:"url,omitempty"`
	ExpandedURL string `json:"expanded_url,omitempty"`
	DisplayURL  string `json:"display_url,omitempty"`
	Indices     []int  `json:"indices,omitempty"`
}

type UserMention struct {
	ID         int64  `json:"id,omitempty"`
	IDStr      string `json:"id_str,omitempty"`
	Name       string `json:"name,omitempty"`
	ScreenName string `json:"screen_name,omitempty"`
	Indices    []int  `json:"indices,omitempty"`
}

type Entities struct {
	Hashtags     []any            `json:"hashtags,omitempty"`
	Symbols      []any            `json:"symbols,omitempty"`
	UserMentions []UserMention    `json:"user_mentions,omitempty"`
	URLs         []URLs           `json:"urls,omitempty"`
	Media        []AttachmentInfo `json:"media,omitempty"`
}
type OriginalInfo struct {
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}
type Thumb struct {
	W      int    `json:"w,omitempty"`
	H      int    `json:"h,omitempty"`
	Resize string `json:"resize,omitempty"`
}
type Small struct {
	W      int    `json:"w,omitempty"`
	H      int    `json:"h,omitempty"`
	Resize string `json:"resize,omitempty"`
}
type Large struct {
	W      int    `json:"w,omitempty"`
	H      int    `json:"h,omitempty"`
	Resize string `json:"resize,omitempty"`
}
type Medium struct {
	W      int    `json:"w,omitempty"`
	H      int    `json:"h,omitempty"`
	Resize string `json:"resize,omitempty"`
}
type Sizes struct {
	Thumb  Thumb  `json:"thumb,omitempty"`
	Small  Small  `json:"small,omitempty"`
	Large  Large  `json:"large,omitempty"`
	Medium Medium `json:"medium,omitempty"`
}
type Variant struct {
	Bitrate     int    `json:"bitrate,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	URL         string `json:"url,omitempty"`
}
type VideoInfo struct {
	AspectRatio    []int     `json:"aspect_ratio,omitempty"`
	DurationMillis int       `json:"duration_millis,omitempty"`
	Variants       []Variant `json:"variants,omitempty"`
}

func (v *VideoInfo) GetHighestBitrateVariant() (Variant, error) {
	if len(v.Variants) == 0 {
		return Variant{}, fmt.Errorf("no variants available")
	}

	maxVariant := v.Variants[0]
	for _, variant := range v.Variants[1:] {
		if variant.Bitrate > maxVariant.Bitrate {
			maxVariant = variant
		}
	}

	return maxVariant, nil
}

type Features struct {
}
type Rgb struct {
	Red   int `json:"red,omitempty"`
	Green int `json:"green,omitempty"`
	Blue  int `json:"blue,omitempty"`
}
type Palette struct {
	Rgb        Rgb     `json:"rgb,omitempty"`
	Percentage float64 `json:"percentage,omitempty"`
}
type ExtMediaColor struct {
	Palette []Palette `json:"palette,omitempty"`
}
type MediaStats struct {
	R   string `json:"r,omitempty"`
	TTL int    `json:"ttl,omitempty"`
}
type Ok struct {
	Palette   []Palette `json:"palette,omitempty"`
	ViewCount string    `json:"view_count,omitempty"`
}
type R struct {
	Ok any `json:"ok,omitempty"`
}
type MediaColor struct {
	R   any `json:"r,omitempty"`
	TTL int `json:"ttl,omitempty"`
}
type AltTextR struct {
	Ok string `json:"ok,omitempty"`
}
type AltText struct {
	// this is weird, it can be both string or AltTextR struct object
	R   any `json:"r,omitempty"`
	TTL int `json:"ttl,omitempty"`
}

// different for video/image/gif
type Ext struct {
	MediaStats any        `json:"mediaStats,omitempty"`
	MediaColor MediaColor `json:"mediaColor,omitempty"`
	AltText    AltText    `json:"altText,omitempty"`
}
type AttachmentInfo struct {
	ID            int64         `json:"id,omitempty"`
	IDStr         string        `json:"id_str,omitempty"`
	Indices       []int         `json:"indices,omitempty"`
	MediaURL      string        `json:"media_url,omitempty"`
	MediaURLHTTPS string        `json:"media_url_https,omitempty"`
	URL           string        `json:"url,omitempty"`
	DisplayURL    string        `json:"display_url,omitempty"`
	ExpandedURL   string        `json:"expanded_url,omitempty"`
	Type          string        `json:"type,omitempty"`
	Filename      string        `json:"filename,omitempty"`
	FilesizeBytes int64         `json:"filesize_bytes,omitempty"`
	OriginalInfo  OriginalInfo  `json:"original_info,omitempty"`
	Sizes         Sizes         `json:"sizes,omitempty"`
	VideoInfo     VideoInfo     `json:"video_info,omitempty"`
	Features      Features      `json:"features,omitempty"`
	ExtMediaColor ExtMediaColor `json:"ext_media_color,omitempty"`
	ExtAltText    string        `json:"ext_alt_text,omitempty"`
	Ext           Ext           `json:"ext,omitempty"`
	AudioOnly     bool          `json:"audio_only,omitempty"`
	MediaHashKey  string        `json:"media_hash_key,omitempty"`
}

type AttachmentCard struct {
	BindingValues AttachmentCardBinding `json:"binding_values,omitempty"`
}

type AttachmentCardBinding struct {
	CardURL     AttachmentCardBindingValue `json:"card_url,omitempty"`
	Description AttachmentCardBindingValue `json:"description,omitempty"`
	Domain      AttachmentCardBindingValue `json:"domain,omitempty"`
	Title       AttachmentCardBindingValue `json:"title,omitempty"`
	VanityUrl   AttachmentCardBindingValue `json:"vanity_url,omitempty"`
}

type AttachmentCardBindingValue struct {
	StringValue string `json:"string_value,omitempty"`
}

type AttachmentTweet struct {
	DisplayURL  string                `json:"display_url,omitempty"`
	ExpandedURL string                `json:"expanded_url,omitempty"`
	Status      AttachmentTweetStatus `json:"status,omitempty"`
}

type AttachmentTweetStatus struct {
	FullText string   `json:"full_text,omitempty"`
	Entities Entities `json:"entities,omitempty"`
	User     User     `json:"user,omitempty"`
}
