package response

type UserByScreenNameResponse struct {
	Data UserByScreenNameData `json:"data,omitempty"`
}

type UserByScreenNameData struct {
	User UserByScreenNameUser `json:"user,omitempty"`
}

type UserByScreenNameUser struct {
	Result UserResult `json:"result,omitempty"`
}

type UserResult struct {
	ID     string     `json:"id,omitempty"`
	RestID string     `json:"rest_id,omitempty"`
	Avatar UserAvatar `json:"avatar,omitempty"`
	Core   UserCore   `json:"core,omitempty"`
}

type UserAvatar struct {
	ImageURL string `json:"image_url,omitempty"`
}

type UserCore struct {
	Name       string `json:"name,omitempty"`
	ScreenName string `json:"screen_name,omitempty"`
}
