package twittermeow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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
	log := zerolog.Ctx(ctx).With().
		Str("url", url).
		Str("method", method).
		Str("request_body", string(payload)).
		Str("function", "MakeRequest").
		Logger()
	var attempts int
	for {
		attempts++
		start := time.Now()
		resp, respDat, err := c.makeRequestDirect(ctx, url, method, headers, payload, contentType)
		dur := time.Since(start)
		if err == nil {
			logLevel := zerolog.DebugLevel
			if strings.Contains(url, "dm/user_updates.json") {
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
				Str("request_body", string(payload)).
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

func (c *Client) makeRequestDirect(ctx context.Context, url string, method string, headers http.Header, payload []byte, contentType types.ContentType) (*http.Response, []byte, error) {
	newRequest, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrRequestCreateFailed, err)
	}

	if contentType != types.ContentTypeNone {
		headers.Set("content-type", string(contentType))
	}

	newRequest.Header = headers

	response, err := c.HTTP.Do(newRequest)
	defer func() {
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
	}()
	if err != nil {
		if errors.Is(err, ErrRedirectAttempted) {
			return response, nil, err
		}
		return nil, nil, fmt.Errorf("%w: %w", ErrRequestFailed, err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrResponseReadFailed, err)
	}
	c.Logger.Trace().
		Int("status_code", response.StatusCode).
		Str("response_body", string(responseBody)).
		Msg("Raw HTTP response")
	if response.StatusCode >= 400 {
		var respErr TwitterErrors
		if json.Unmarshal(responseBody, &respErr) == nil {
			return response, responseBody, fmt.Errorf("HTTP %d: %w", response.StatusCode, &respErr)
		} else if len(responseBody) == 0 {
			return response, responseBody, fmt.Errorf("HTTP %d (no response body)", response.StatusCode)
		} else if len(responseBody) < 512 {
			return response, responseBody, fmt.Errorf("HTTP %d: %s", response.StatusCode, responseBody)
		}
		return response, responseBody, fmt.Errorf("HTTP %d (%d bytes of data)", response.StatusCode, len(responseBody))
	}

	return response, responseBody, nil
}
