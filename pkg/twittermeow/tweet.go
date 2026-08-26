package twittermeow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

const (
	tweetPreviewTimeout        = 10 * time.Second
	maxTweetOEmbedResponseSize = 1 << 20
)

type tweetOEmbedResponse struct {
	URL        string `json:"url"`
	AuthorName string `json:"author_name"`
	HTML       string `json:"html"`
}

// GetTweetPreview fetches public post text and author metadata from X's oEmbed
// endpoint. It intentionally sends no account cookies or authorization headers.
func (c *Client) GetTweetPreview(ctx context.Context, postURL string) (*types.AttachmentTweet, error) {
	ctx, cancel := context.WithTimeout(ctx, tweetPreviewTimeout)
	defer cancel()
	if c == nil || c.HTTP == nil {
		return nil, fmt.Errorf("twitter HTTP client is not configured")
	}
	requestURL := endpoints.TWEET_OEMBED_URL + "?dnt=true&omit_script=true&url=" + url.QueryEscape(postURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare tweet oEmbed request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.GetUserAgent())

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tweet oEmbed response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("tweet oEmbed returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxTweetOEmbedResponseSize+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read tweet oEmbed response: %w", err)
	}
	if len(body) > maxTweetOEmbedResponseSize {
		return nil, fmt.Errorf("tweet oEmbed response exceeds %d bytes", maxTweetOEmbedResponseSize)
	}
	var oembed tweetOEmbedResponse
	if err = json.Unmarshal(body, &oembed); err != nil {
		return nil, fmt.Errorf("failed to parse tweet oEmbed response: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(oembed.HTML))
	if err != nil {
		return nil, fmt.Errorf("failed to parse tweet oEmbed HTML: %w", err)
	}
	post := doc.Find("blockquote.twitter-tweet p").First()
	post.Find("br").Each(func(_ int, br *goquery.Selection) {
		br.ReplaceWithHtml("\n")
	})
	postText := strings.TrimSpace(post.Text())
	if postText == "" {
		return nil, fmt.Errorf("tweet oEmbed post text is empty")
	}
	canonicalURL := strings.TrimSpace(oembed.URL)
	if canonicalURL == "" {
		canonicalURL = postURL
	}

	return &types.AttachmentTweet{
		ExpandedURL: canonicalURL,
		Status: types.AttachmentTweetStatus{
			FullText: postText,
			User: types.User{
				Name: strings.TrimSpace(oembed.AuthorName),
			},
		},
	}, nil
}
