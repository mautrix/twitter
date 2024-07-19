package types

type Urls struct {
	URL         string `json:"url,omitempty"`
	ExpandedURL string `json:"expanded_url,omitempty"`
	DisplayURL  string `json:"display_url,omitempty"`
	Indices     []int  `json:"indices,omitempty"`
}
type Entities struct {
	Hashtags     []any  `json:"hashtags,omitempty"`
	Symbols      []any  `json:"symbols,omitempty"`
	UserMentions []any  `json:"user_mentions,omitempty"`
	Urls         []Urls `json:"urls,omitempty"`
}
type OriginalInfo struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
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
type Variants struct {
	Bitrate     int    `json:"bitrate,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	URL         string `json:"url,omitempty"`
}
type VideoInfo struct {
	AspectRatio []int      `json:"aspect_ratio,omitempty"`
	Variants    []Variants `json:"variants,omitempty"`
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
	Palette []Palette `json:"palette,omitempty"`
}
type R struct {
	Ok Ok `json:"ok,omitempty"`
}
type MediaColor struct {
	R   R   `json:"r,omitempty"`
	TTL int `json:"ttl,omitempty"`
}
type AltTextR struct {
	Ok string `json:"ok,omitempty"`
}
type AltText struct {
	// this is weird, it can be both string or AltTextR struct object
	R   interface{} `json:"r,omitempty"`
	TTL int         `json:"ttl,omitempty"`
}
type Ext struct {
	MediaStats MediaStats `json:"mediaStats,omitempty"`
	MediaColor MediaColor `json:"mediaColor,omitempty"`
	AltText    AltText    `json:"altText,omitempty"`
}
type AnimatedGif struct {
	ID            int64         `json:"id,omitempty"`
	IDStr         string        `json:"id_str,omitempty"`
	Indices       []int         `json:"indices,omitempty"`
	MediaURL      string        `json:"media_url,omitempty"`
	MediaURLHTTPS string        `json:"media_url_https,omitempty"`
	URL           string        `json:"url,omitempty"`
	DisplayURL    string        `json:"display_url,omitempty"`
	ExpandedURL   string        `json:"expanded_url,omitempty"`
	Type          string        `json:"type,omitempty"`
	OriginalInfo  OriginalInfo  `json:"original_info,omitempty"`
	Sizes         Sizes         `json:"sizes,omitempty"`
	VideoInfo     VideoInfo     `json:"video_info,omitempty"`
	Features      Features      `json:"features,omitempty"`
	ExtMediaColor ExtMediaColor `json:"ext_media_color,omitempty"`
	ExtAltText    string        `json:"ext_alt_text,omitempty"`
	Ext           Ext           `json:"ext,omitempty"`
	AudioOnly     bool          `json:"audio_only,omitempty"`
}
type Attachment struct {
	AnimatedGif AnimatedGif `json:"animated_gif,omitempty"`
}
