package connector

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
)

// twitterURLRegex matches X.com and twitter.com status URLs
// Captures: (1) username, (2) tweet ID
var twitterURLRegex = regexp.MustCompile(
	`https?://(?:twitter\.com|x\.com)/([^/\s]+)/status/(\d+)(?:\?[^\s]*)?`,
)

// extractTwitterURLs finds all tweet URLs in text
// Returns map[tweetID]canonicalURL to deduplicate same tweet with different params
func extractTwitterURLs(text string) map[string]string {
	matches := twitterURLRegex.FindAllStringSubmatch(text, -1)
	urls := make(map[string]string, len(matches))

	for _, match := range matches {
		if len(match) >= 3 {
			tweetID := match[2]
			// Normalize to x.com
			canonicalURL := fmt.Sprintf("https://x.com/%s/status/%s", match[1], tweetID)
			urls[tweetID] = canonicalURL
		}
	}

	return urls
}

// truncateText truncates text to maxLen characters, adding ellipsis if truncated
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

// isMinimalPreview checks if a preview has only title/URL with no description or image
func isMinimalPreview(preview *event.BeeperLinkPreview) bool {
	if preview == nil {
		return false
	}
	return preview.LinkPreview.Description == "" && preview.LinkPreview.ImageURL == ""
}

// getExistingPreviewURL extracts the canonical URL from a preview
func getExistingPreviewURL(preview *event.BeeperLinkPreview) string {
	if preview == nil {
		return ""
	}
	return preview.LinkPreview.CanonicalURL
}

// fetchAndGenerateTweetPreviews generates link previews for X.com URLs in text
// It returns previews and a map of replacement indices (which existing preview to replace, if any)
// The map key is the index in the returned previews, value is the index to replace in existing (or -1 if append)
func (tc *TwitterClient) fetchAndGenerateTweetPreviews(
	ctx context.Context,
	portal *bridgev2.Portal,
	intent bridgev2.MatrixAPI,
	text string,
	existingPreviews []*event.BeeperLinkPreview,
) ([]*event.BeeperLinkPreview, map[int]int) {
	log := zerolog.Ctx(ctx)

	// DEBUG: Log input parameters
	log.Info().
		Str("text", text).
		Int("existing_previews_count", len(existingPreviews)).
		Msg("[LinkPreview] fetchAndGenerateTweetPreviews called")

	// DEBUG: Log each existing preview
	for idx, preview := range existingPreviews {
		if preview != nil {
			log.Info().
				Int("index", idx).
				Str("canonical_url", preview.LinkPreview.CanonicalURL).
				Str("title", preview.LinkPreview.Title).
				Str("description", preview.LinkPreview.Description).
				Str("image_url", string(preview.LinkPreview.ImageURL)).
				Bool("is_minimal", isMinimalPreview(preview)).
				Msg("[LinkPreview] Existing preview")
		}
	}

	// Extract URLs from message text
	tweetURLs := extractTwitterURLs(text)
	log.Info().
		Int("urls_from_text", len(tweetURLs)).
		Msg("[LinkPreview] URLs extracted from message text")

	// Build map of existing preview URLs to their indices
	existingURLs := make(map[string]int)
	for idx, preview := range existingPreviews {
		url := getExistingPreviewURL(preview)
		if url != "" {
			existingURLs[url] = idx

			// ALSO extract tweet IDs from existing minimal preview URLs
			// This catches URLs from msg.Attachment.Tweet.ExpandedURL that aren't in msg.Text
			if isMinimalPreview(preview) {
				log.Info().
					Str("url", url).
					Msg("[LinkPreview] Checking minimal preview URL for X.com pattern")

				matches := twitterURLRegex.FindStringSubmatch(url)
				log.Info().
					Int("matches_count", len(matches)).
					Msg("[LinkPreview] Regex match result")

				if len(matches) >= 3 {
					tweetID := matches[2]
					if _, exists := tweetURLs[tweetID]; !exists {
						// Normalize to x.com format
						canonicalURL := fmt.Sprintf("https://x.com/%s/status/%s", matches[1], tweetID)
						tweetURLs[tweetID] = canonicalURL
						log.Info().
							Str("tweet_id", tweetID).
							Str("original_url", url).
							Str("canonical_url", canonicalURL).
							Msg("[LinkPreview] Found X.com URL in minimal preview (not in message text)")
					}
				}
			}
		}
	}

	log.Info().
		Int("total_tweet_urls", len(tweetURLs)).
		Msg("[LinkPreview] Total tweet URLs to process")

	if len(tweetURLs) == 0 {
		log.Info().Msg("[LinkPreview] No tweet URLs found, returning nil")
		return nil, nil
	}

	// Limit to prevent spam (max 3 previews per message)
	const maxPreviews = 3
	count := 0

	previews := make([]*event.BeeperLinkPreview, 0, len(tweetURLs))
	replacementMap := make(map[int]int) // map[new preview index] = existing preview index to replace

	for tweetID, canonicalURL := range tweetURLs {
		if count >= maxPreviews {
			log.Debug().
				Int("total_urls", len(tweetURLs)).
				Int("previews_created", count).
				Msg("Skipping additional tweet previews (max limit reached)")
			break
		}

		// Check if we already have a preview for this URL
		existingIdx, hasExisting := existingURLs[canonicalURL]
		if hasExisting {
			existingPreview := existingPreviews[existingIdx]
			// If it's not minimal, skip (already have good preview)
			if !isMinimalPreview(existingPreview) {
				log.Debug().
					Str("tweet_id", tweetID).
					Msg("Skipping X.com URL: already have good preview from attachment")
				continue
			}
			// It's minimal - we'll replace it
		}

		// Fetch tweet data
		tweetData, err := tc.client.FetchTweet(ctx, tweetID)
		if err != nil {
			log.Warn().
				Err(err).
				Str("tweet_id", tweetID).
				Msg("Failed to fetch tweet for link preview")
			continue
		}

		// Generate preview
		preview, err := tc.generateTweetPreview(ctx, portal, intent, canonicalURL, tweetData)
		if err != nil {
			log.Warn().
				Err(err).
				Str("tweet_id", tweetID).
				Msg("Failed to generate tweet preview")
			continue
		}

		previewIdx := len(previews)
		previews = append(previews, preview)

		// Track if this replaces an existing minimal preview
		if hasExisting {
			replacementMap[previewIdx] = existingIdx
			log.Debug().
				Str("tweet_id", tweetID).
				Int("replacement_index", existingIdx).
				Msg("Replacing minimal X.com preview with fetched data")
		}

		count++
	}

	return previews, replacementMap
}

// generateTweetPreview creates a BeeperLinkPreview from tweet GraphQL data
func (tc *TwitterClient) generateTweetPreview(
	ctx context.Context,
	portal *bridgev2.Portal,
	intent bridgev2.MatrixAPI,
	canonicalURL string,
	tweetData *response.FetchPostQueryResponse,
) (*event.BeeperLinkPreview, error) {
	log := zerolog.Ctx(ctx)

	result := tweetData.Data.PostResult.Result

	// Extract author name
	authorName := "Unknown"
	if result.Core.UserResults.Result.Legacy.Name != "" {
		authorName = result.Core.UserResults.Result.Legacy.Name
	}

	// Build preview
	preview := &event.BeeperLinkPreview{
		LinkPreview: event.LinkPreview{
			CanonicalURL: canonicalURL,
			Title:        fmt.Sprintf("%s on X", authorName),
			Description:  truncateText(result.Legacy.FullText, 500),
		},
	}

	// Download and upload first image if available
	if result.Legacy.ExtendedEntities != nil && len(result.Legacy.ExtendedEntities.Media) > 0 {
		media := result.Legacy.ExtendedEntities.Media[0]

		if media.Type == "photo" {
			resp, err := downloadFile(ctx, tc.client, media.MediaURLHTTPS)
			if err != nil {
				log.Warn().
					Err(err).
					Str("tweet_id", result.RestID).
					Msg("Failed to download tweet image for preview")
			} else {
				preview.LinkPreview.ImageType = "image/jpeg"
				preview.LinkPreview.ImageWidth = event.IntOrString(media.OriginalInfo.Width)
				preview.LinkPreview.ImageHeight = event.IntOrString(media.OriginalInfo.Height)
				preview.LinkPreview.ImageSize = event.IntOrString(resp.ContentLength)

				imageURL, _, err := intent.UploadMediaStream(ctx, portal.MXID, resp.ContentLength, false, func(file io.Writer) (*bridgev2.FileStreamResult, error) {
					_, err := io.Copy(file, resp.Body)
					if err != nil {
						return nil, err
					}
					return &bridgev2.FileStreamResult{
						MimeType: "image/jpeg",
						FileName: "tweet-image.jpeg",
					}, nil
				})
				if err != nil {
					log.Warn().
						Err(err).
						Str("tweet_id", result.RestID).
						Msg("Failed to upload tweet image to Matrix")
				} else {
					preview.LinkPreview.ImageURL = imageURL
				}
			}
		}
	}

	return preview, nil
}
