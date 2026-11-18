package twittermeow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

var MaxHTTPRetries = 5

var (
	ErrRedirectAttempted   = errors.New("redirect attempted")
	ErrRequestCreateFailed = errors.New("failed to create request")
	ErrRequestFailed       = errors.New("failed to send request")
	ErrResponseReadFailed  = errors.New("failed to read response body")
	ErrMaxRetriesReached   = errors.New("maximum retries reached")
)

func (c *Client) MakeRequest(ctx context.Context, url string, method string, headers http.Header, payload []byte, contentType types.ContentType) (*http.Response, []byte, error) {
	var pl io.Reader
	if payload != nil {
		pl = bytes.NewReader(payload)
	}
	return c.MakeRequestFull(ctx, MakeRequestParams{
		URL:         url,
		Method:      method,
		Headers:     headers,
		Payload:     pl,
		ContentType: contentType,
	})
}

type MakeRequestParams struct {
	URL          string
	Method       string
	Headers      http.Header
	Payload      io.Reader
	ContentType  types.ContentType
	DontReadBody bool
}

func (c *Client) MakeRequestFull(ctx context.Context, params MakeRequestParams) (*http.Response, []byte, error) {
	log := zerolog.Ctx(ctx).With().
		Str("url", params.URL).
		Str("method", params.Method).
		Str("function", "MakeRequest").
		Logger()
	rateLimitKey := getRateLimitKey(params.URL)
	var attempts int
	for {
		attempts++
		if err := c.checkRateLimit(rateLimitKey); err != nil {
			log.Warn().Err(err).Msg("Not sending request to rate limited endpoint")
			return nil, nil, err
		}
		start := time.Now()
		resp, respDat, err := c.makeRequestDirect(ctx, params)
		dur := time.Since(start)
		c.updateRateLimit(rateLimitKey, resp)
		if err == nil {
			logLevel := zerolog.DebugLevel
			if strings.Contains(params.URL, "dm/user_updates.json") {
				// Don't spam log full of user_updates requests
				logLevel = zerolog.TraceLevel
			}
			log.WithLevel(logLevel).
				Dur("duration", dur).
				Msg("Request successful")
			return resp, respDat, nil
		} else if IsClientHTTPError(err) {
			// The request never left the user's device, so retrying it server-side is
			// pointless (and would ask the client to run the same request again).
			log.Err(err).
				Dur("duration", dur).
				Msg("Request failed on the login client")
			return nil, nil, err
		} else if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 {
			log.Error().
				Err(err).
				Dur("duration", dur).
				Str("response_body", string(respDat)).
				Msg("Request failed")
			return nil, respDat, err
		} else if attempts > MaxHTTPRetries {
			log.Err(err).
				Dur("duration", dur).
				Msg("Request failed, giving up")
			return nil, nil, fmt.Errorf("%w: %w", ErrMaxRetriesReached, err)
		} else if errors.Is(err, ErrRedirectAttempted) {
			c.Logger.Err(err).
				Str("location", resp.Header.Get("Location")).
				Dur("duration", dur).
				Msg("Redirect attempted")
			return resp, nil, err
		} else if ctx.Err() != nil {
			return resp, nil, ctx.Err()
		}
		log.Err(err).
			Dur("duration", dur).
			Str("response_body", string(respDat)).
			Msg("Request failed, retrying")
		time.Sleep(time.Duration(attempts) * 3 * time.Second)
	}
}

func getRateLimitKey(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	return parsed.Host + parsed.Path
}

func (c *Client) checkRateLimit(key string) error {
	c.rateLimitLock.Lock()
	defer c.rateLimitLock.Unlock()
	reset, ok := c.rateLimitResets[key]
	if !ok {
		return nil
	} else if time.Now().After(reset) {
		delete(c.rateLimitResets, key)
		return nil
	}
	return RateLimitError{Reset: reset}
}

func (c *Client) updateRateLimit(key string, resp *http.Response) {
	if resp == nil {
		return
	}
	var reset time.Time
	if resp.Header.Get("x-rate-limit-remaining") == "0" {
		if resetUnix, err := strconv.ParseInt(resp.Header.Get("x-rate-limit-reset"), 10, 64); err == nil {
			reset = time.Unix(resetUnix, 0)
		}
	}
	if retryAfter := resp.Header.Get("Retry-After"); resp.StatusCode == http.StatusTooManyRequests && retryAfter != "" {
		var retryAt time.Time
		if seconds, err := strconv.Atoi(retryAfter); err == nil {
			retryAt = time.Now().Add(time.Duration(seconds) * time.Second)
		} else if date, err := http.ParseTime(retryAfter); err == nil {
			retryAt = date
		}
		if retryAt.After(reset) {
			reset = retryAt
		}
	}
	if time.Now().After(reset) {
		return
	}
	c.rateLimitLock.Lock()
	if reset.After(c.rateLimitResets[key]) {
		c.rateLimitResets[key] = reset
	}
	c.rateLimitLock.Unlock()
}

func (c *Client) makeRequestDirect(ctx context.Context, params MakeRequestParams) (*http.Response, []byte, error) {
	newRequest, err := http.NewRequestWithContext(ctx, params.Method, params.URL, params.Payload)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrRequestCreateFailed, err)
	}

	if params.ContentType != types.ContentTypeNone {
		params.Headers.Set("content-type", string(params.ContentType))
	}

	newRequest.Header = params.Headers

	response, err := c.HTTP.Do(newRequest)
	dontCloseBody := false
	defer func() {
		if response != nil && response.Body != nil && !dontCloseBody {
			_ = response.Body.Close()
		}
	}()
	if err != nil {
		if errors.Is(err, ErrRedirectAttempted) {
			return response, nil, err
		}
		return nil, nil, fmt.Errorf("%w: %w", ErrRequestFailed, err)
	}

	var responseBody []byte
	if !params.DontReadBody || response.StatusCode >= 300 || response.StatusCode < 200 {
		responseBody, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %w", ErrResponseReadFailed, err)
		}
	} else {
		dontCloseBody = true
	}
	c.Logger.Trace().
		Int("status_code", response.StatusCode).
		Str("response_body", string(responseBody)).
		Msg("Raw HTTP response")
	if response.StatusCode >= 400 {
		var respErr TwitterErrors
		if json.Unmarshal(responseBody, &respErr) == nil {
			respErr.StatusCode = response.StatusCode
			return response, responseBody, fmt.Errorf("HTTP %d: %w", response.StatusCode, &respErr)
		} else if response.StatusCode == http.StatusTooManyRequests && bytes.Contains(responseBody, []byte("1015")) && bytes.Contains(responseBody, []byte("Cloudflare")) {
			return response, responseBody, fmt.Errorf("HTTP %d: %w", response.StatusCode, ErrCloudflareRateLimited)
		} else if response.StatusCode == http.StatusTooManyRequests {
			return response, responseBody, fmt.Errorf("HTTP %d: %w", response.StatusCode, ErrRatelimitExceeded)
		} else if len(responseBody) == 0 {
			return response, responseBody, fmt.Errorf("HTTP %d (no response body)", response.StatusCode)
		} else if len(responseBody) < 512 {
			return response, responseBody, fmt.Errorf("HTTP %d: %s", response.StatusCode, responseBody)
		}
		return response, responseBody, fmt.Errorf("HTTP %d (%d bytes of data)", response.StatusCode, len(responseBody))
	}

	return response, responseBody, nil
}
