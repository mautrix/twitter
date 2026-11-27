package twittermeow

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"strings"
	"sync/atomic"
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

	errRequestNotSent = errors.New("request failed before obtaining a connection")
)

func (c *Client) MakeRequest(ctx context.Context, url string, method string, headers http.Header, payload []byte, contentType types.ContentType) (*http.Response, []byte, error) {
	return c.MakeRequestFull(ctx, MakeRequestParams{
		URL:         url,
		Method:      method,
		Headers:     headers,
		Payload:     payload,
		ContentType: contentType,
	})
}

func (c *Client) makeIdempotentPost(ctx context.Context, url string, headers http.Header, payload []byte, contentType types.ContentType) (*http.Response, []byte, error) {
	return c.MakeRequestFull(ctx, MakeRequestParams{
		URL:         url,
		Method:      http.MethodPost,
		Headers:     headers,
		Payload:     payload,
		ContentType: contentType,
		Idempotent:  true,
	})
}

type MakeRequestParams struct {
	URL          string
	Method       string
	Headers      http.Header
	Payload      []byte
	ContentType  types.ContentType
	DontReadBody bool
	// Idempotent allows retrying a non-GET request after it may have reached the server.
	Idempotent bool
}

func (c *Client) MakeRequestFull(ctx context.Context, params MakeRequestParams) (*http.Response, []byte, error) {
	log := zerolog.Ctx(ctx).With().
		Str("url", params.URL).
		Str("method", params.Method).
		Str("function", "MakeRequest").
		Logger()
	var attempts int
	for {
		attempts++
		start := time.Now()
		resp, respDat, err := c.makeRequestDirect(ctx, params)
		dur := time.Since(start)
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
		} else if errors.Is(err, ErrRedirectAttempted) {
			c.Logger.Err(err).
				Str("location", resp.Header.Get("Location")).
				Dur("duration", dur).
				Msg("Redirect attempted")
			return resp, nil, err
		} else if attempts > MaxHTTPRetries {
			log.Err(err).
				Dur("duration", dur).
				Msg("Request failed, giving up")
			return nil, nil, fmt.Errorf("%w: %w", ErrMaxRetriesReached, err)
		} else if ctx.Err() != nil {
			return resp, nil, ctx.Err()
		} else if !canRetryRequest(params, err) {
			log.Err(err).
				Dur("duration", dur).
				Msg("Request failed after it may have reached the server, not retrying")
			return nil, respDat, err
		}
		log.Err(err).
			Dur("duration", dur).
			Str("response_body", string(respDat)).
			Msg("Request failed, retrying")
		time.Sleep(time.Duration(attempts) * 3 * time.Second)
	}
}

func canRetryRequest(params MakeRequestParams, err error) bool {
	switch params.Method {
	case "", http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return params.Idempotent || errors.Is(err, errRequestNotSent)
}

func (c *Client) makeRequestDirect(ctx context.Context, params MakeRequestParams) (*http.Response, []byte, error) {
	var getConn, gotConn atomic.Bool
	ctx = httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{
		GetConn: func(string) { getConn.Store(true) },
		GotConn: func(httptrace.GotConnInfo) { gotConn.Store(true) },
	})
	var body io.Reader
	if params.Payload != nil {
		body = bytes.NewReader(params.Payload)
	}
	newRequest, err := http.NewRequestWithContext(ctx, params.Method, params.URL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrRequestCreateFailed, err)
	}

	// An abandoned attempt's transport goroutine may still read the previous header map.
	newRequest.Header = params.Headers.Clone()
	if params.ContentType != types.ContentTypeNone {
		newRequest.Header.Set("content-type", string(params.ContentType))
	}

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
		err = fmt.Errorf("%w: %w", ErrRequestFailed, err)
		// Transports that don't report GetConn give no proof the request wasn't sent.
		if getConn.Load() && !gotConn.Load() {
			err = fmt.Errorf("%w: %w", errRequestNotSent, err)
		}
		return nil, nil, err
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
	if response.StatusCode >= 300 && response.StatusCode < 400 {
		return response, responseBody, fmt.Errorf("HTTP %d: %w", response.StatusCode, ErrRedirectAttempted)
	}
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
