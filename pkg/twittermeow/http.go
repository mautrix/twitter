package twittermeow

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"
)

const MaxHTTPRetries = 5

var (
	ErrRedirectAttempted   = errors.New("redirect attempted")
	ErrRequestCreateFailed = errors.New("failed to create request")
	ErrRequestFailed       = errors.New("failed to send request")
	ErrResponseReadFailed  = errors.New("failed to read response body")
	ErrMaxRetriesReached   = errors.New("maximum retries reached")
)

func (c *Client) MakeRequest(url string, method string, headers http.Header, payload []byte, contentType types.ContentType) (*http.Response, []byte, error) {
	var attempts int
	for {
		attempts++
		start := time.Now()
		resp, respDat, err := c.makeRequestDirect(url, method, headers, payload, contentType)
		dur := time.Since(start)
		if err == nil {
			var logEvt *zerolog.Event
			if strings.HasPrefix(url, endpoints.DM_USER_UPDATES_URL) {
				logEvt = c.Logger.Trace()
			} else {
				logEvt = c.Logger.Debug()
			}
			logEvt.
				Str("url", url).
				Str("method", method).
				Dur("duration", dur).
				Msg("Request successful")
			return resp, respDat, nil
		} else if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 502 {
			c.Logger.Err(err).
				Str("url", url).
				Str("method", method).
				Dur("duration", dur).
				Msg("Request failed")
			return nil, nil, err
		} else if attempts > MaxHTTPRetries {
			c.Logger.Err(err).
				Str("url", url).
				Str("method", method).
				Dur("duration", dur).
				Msg("Request failed, giving up")
			return nil, nil, fmt.Errorf("%w: %w", ErrMaxRetriesReached, err)
		} else if errors.Is(err, ErrRedirectAttempted) {
			location := resp.Header.Get("Location")
			c.Logger.Err(err).
				Str("url", url).
				Str("location", location).
				Str("method", method).
				Dur("duration", dur).
				Msg("Redirect attempted")
			return resp, nil, err
		}
		c.Logger.Err(err).
			Str("url", url).
			Str("method", method).
			Dur("duration", dur).
			Msg("Request failed, retrying")
		time.Sleep(time.Duration(attempts) * 3 * time.Second)
	}
}

func (c *Client) makeRequestDirect(url string, method string, headers http.Header, payload []byte, contentType types.ContentType) (*http.Response, []byte, error) {
	newRequest, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrRequestCreateFailed, err)
	}

	if contentType != types.NONE {
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
	if response.StatusCode >= 400 {
		if len(responseBody) < 512 {
			return response, responseBody, fmt.Errorf("HTTP %d: %s", response.StatusCode, responseBody)
		}
		return response, responseBody, fmt.Errorf("HTTP %d (%d bytes of data)", response.StatusCode, len(responseBody))
	}

	return response, responseBody, nil
}
