package twittermeow

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mau.fi/util/ffmpeg"
	"go.mau.fi/util/random"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/crypto"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/types"
)

func (c *Client) UploadMedia(ctx context.Context, params *payload.UploadMediaQuery, mediaBytes []byte) (*response.FinalizedUploadMediaResponse, error) {
	params.Command = "INIT"
	if mediaBytes != nil {
		params.TotalBytes = len(mediaBytes)
	}

	encodedQuery, err := params.Encode()
	if err != nil {
		return nil, err
	}

	var finalizedMediaResultBytes []byte

	url := fmt.Sprintf("%s?%s", endpoints.UPLOAD_MEDIA_URL, string(encodedQuery))
	headerOpts := HeaderOpts{
		WithNonAuthBearer: true,
		WithXCsrfToken:    true,
		WithCookies:       true,
		Origin:            endpoints.BASE_URL,
		Referer:           endpoints.BASE_URL + "/",
		Extra: map[string]string{
			"sec-fetch-dest": "empty",
			"sec-fetch-mode": "cors",
			"sec-fetch-site": "same-origin",
			"accept":         "*/*",
		},
	}
	headers := c.buildHeaders(headerOpts)

	_, respBody, err := c.MakeRequest(ctx, url, http.MethodPost, headers, nil, types.ContentTypeNone)
	if err != nil {
		return nil, err
	}

	initUploadResponse := &response.InitUploadMediaResponse{}
	err = json.Unmarshal(respBody, initUploadResponse)
	if err != nil {
		return nil, err
	}

	segmentIndex := 0
	if mediaBytes != nil {
		for chunk := range slices.Chunk(mediaBytes, 6*1024*1024) {
			appendMediaPayload, contentType, err := c.newMediaAppendPayload(chunk)
			if err != nil {
				return nil, err
			}
			headers.Add("content-type", contentType)

			url = fmt.Sprintf("%s?command=APPEND&media_id=%s&segment_index=%d", endpoints.UPLOAD_MEDIA_URL, initUploadResponse.MediaIDString, segmentIndex)
			resp, respBody, err := c.MakeRequest(ctx, url, http.MethodPost, headers, appendMediaPayload, types.ContentTypeNone)
			if err != nil {
				return nil, err
			}

			if resp.StatusCode > 204 {
				return nil, fmt.Errorf("failed to append media bytes for media with id %s (status_code=%d, response_body=%s)", initUploadResponse.MediaIDString, resp.StatusCode, string(respBody))
			}
			segmentIndex++
		}

		var originalMd5 string
		if params.MediaCategory == payload.MEDIA_CATEGORY_DM_IMAGE {
			md5Hash := md5.Sum(mediaBytes)
			originalMd5 = hex.EncodeToString(md5Hash[:])
		}

		finalizeMediaQuery := &payload.UploadMediaQuery{
			Command:     "FINALIZE",
			MediaID:     initUploadResponse.MediaIDString,
			OriginalMD5: originalMd5,
		}

		encodedQuery, err = finalizeMediaQuery.Encode()
		if err != nil {
			return nil, err
		}

		url = fmt.Sprintf("%s?%s", endpoints.UPLOAD_MEDIA_URL, string(encodedQuery))
		headers.Del("content-type")
		resp, respBody, err := c.MakeRequest(ctx, url, http.MethodPost, headers, nil, types.ContentTypeNone)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode > 204 {
			return nil, fmt.Errorf("failed to finalize media with id %s (status_code=%d, response_body=%s)", initUploadResponse.MediaIDString, resp.StatusCode, string(respBody))
		}

		finalizedMediaResultBytes = respBody
	} else {
		_, finalizedMediaResultBytes, err = c.GetMediaUploadStatus(ctx, initUploadResponse.MediaIDString, headers)
		if err != nil {
			return nil, err
		}
	}

	finalizedMediaResult := &response.FinalizedUploadMediaResponse{}
	err = json.Unmarshal(finalizedMediaResultBytes, finalizedMediaResult)
	if err != nil {
		return nil, err
	}

	if finalizedMediaResult.ProcessingInfo.State == response.ProcessingStatePending || finalizedMediaResult.ProcessingInfo.State == response.ProcessingStateInProgress {
		// might need to check for error processing states here, I have not encountered any though so I wouldn't know what they look like/are
		for finalizedMediaResult.ProcessingInfo.State != response.ProcessingStateSucceeded {
			finalizedMediaResult, _, err = c.GetMediaUploadStatus(ctx, finalizedMediaResult.MediaIDString, headers)
			if err != nil {
				return nil, err
			}
			c.Logger.Debug().
				Int("progress_percent", finalizedMediaResult.ProcessingInfo.ProgressPercent).
				Int("status_check_interval_seconds", finalizedMediaResult.ProcessingInfo.CheckAfterSecs).
				Str("media_id", finalizedMediaResult.MediaIDString).
				Str("state", string(finalizedMediaResult.ProcessingInfo.State)).
				Msg("Waiting for X to finish processing our media upload...")
			time.Sleep(time.Second * time.Duration(finalizedMediaResult.ProcessingInfo.CheckAfterSecs))
		}
	}

	return finalizedMediaResult, nil
}

func (c *Client) GetMediaUploadStatus(ctx context.Context, mediaID string, h http.Header) (*response.FinalizedUploadMediaResponse, []byte, error) {
	url := fmt.Sprintf("%s?command=STATUS&media_id=%s", endpoints.UPLOAD_MEDIA_URL, mediaID)
	resp, respBody, err := c.MakeRequest(ctx, url, http.MethodGet, h, nil, types.ContentTypeNone)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode > 204 {
		return nil, nil, fmt.Errorf("failed to get status of uploaded media with id %s (status_code=%d, response_body=%s)", mediaID, resp.StatusCode, string(respBody))
	}

	mediaStatusResult := &response.FinalizedUploadMediaResponse{}
	return mediaStatusResult, respBody, json.Unmarshal(respBody, mediaStatusResult)
}

func (c *Client) newMediaAppendPayload(mediaBytes []byte) ([]byte, string, error) {
	var appendMediaPayload bytes.Buffer
	writer := multipart.NewWriter(&appendMediaPayload)

	err := writer.SetBoundary("----WebKitFormBoundary" + random.String(16))
	if err != nil {
		return nil, "", fmt.Errorf("failed to set boundary (%s)", err.Error())
	}

	partHeader := textproto.MIMEHeader{
		"Content-Disposition": []string{`form-data; name="media"; filename="blob"`},
		"Content-Type":        []string{"application/octet-stream"},
	}

	mediaPart, err := writer.CreatePart(partHeader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create multipart writer (%s)", err.Error())
	}

	_, err = mediaPart.Write(mediaBytes)
	if err != nil {
		return nil, "", fmt.Errorf("failed to write data to multipart section (%s)", err.Error())
	}

	err = writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to close multipart writer (%s)", err.Error())
	}

	return appendMediaPayload.Bytes(), writer.FormDataContentType(), nil
}

func (c *Client) ConvertAudioPayload(ctx context.Context, mediaBytes []byte, mimeType string) ([]byte, error) {
	if !ffmpeg.Supported() {
		return nil, errors.New("ffmpeg is required to send voice message")
	}

	// A video part is required to send voice message.
	return ffmpeg.ConvertBytes(ctx, mediaBytes, ".mp4", []string{"-f", "lavfi", "-i", "color=black:s=854x480:r=30"}, []string{"-c:v", "h264", "-c:a", "aac", "-tune", "stillimage", "-shortest"}, mimeType)
}

// XChatMediaUploadResult contains the result of an XChat media upload.
type XChatMediaUploadResult struct {
	MediaHashKey string
	KeyVersion   string
}

// UploadXChatMedia uploads media for XChat messages.
// Media is encrypted using secretstream (XChaCha20-Poly1305) with the conversation key.
func (c *Client) UploadXChatMedia(ctx context.Context, conversationID, messageID string, mediaBytes []byte) (*XChatMediaUploadResult, error) {
	// Get the conversation key for encryption
	convKey, err := c.keyManager.GetLatestConversationKey(ctx, conversationID)
	if err != nil {
		return nil, fmt.Errorf("get conversation key: %w", err)
	}

	// Encrypt media using secretstream (XChaCha20-Poly1305)
	encryptedBytes, err := crypto.SecretstreamEncrypt(mediaBytes, convKey.Key)
	if err != nil {
		return nil, fmt.Errorf("encrypt media: %w", err)
	}

	c.Logger.Debug().
		Int("plaintext_size", len(mediaBytes)).
		Int("encrypted_size", len(encryptedBytes)).
		Str("key_version", convKey.KeyVersion).
		Msg("Encrypted media for XChat upload")

	// Initialize upload with encrypted size
	initResp, err := c.initializeXChatMediaUpload(ctx, conversationID, messageID, len(encryptedBytes))
	if err != nil {
		return nil, fmt.Errorf("initialize upload: %w", err)
	}

	// Upload encrypted bytes in chunks
	numParts, err := c.uploadMediaBytes(ctx, initResp.Data.XChatInitializeMediaUpload.ResumeUploadURL, encryptedBytes)
	if err != nil {
		return nil, fmt.Errorf("upload bytes: %w", err)
	}

	// Finalize upload with the actual number of parts
	_, err = c.finalizeXChatMediaUpload(
		ctx,
		conversationID,
		messageID,
		initResp.Data.XChatInitializeMediaUpload.MediaHashKey,
		initResp.Data.XChatInitializeMediaUpload.ResumeID,
		numParts,
	)
	if err != nil {
		return nil, fmt.Errorf("finalize upload: %w", err)
	}

	return &XChatMediaUploadResult{
		MediaHashKey: initResp.Data.XChatInitializeMediaUpload.MediaHashKey,
		KeyVersion:   convKey.KeyVersion,
	}, nil
}

// initializeXChatMediaUpload initializes an XChat media upload.
func (c *Client) initializeXChatMediaUpload(ctx context.Context, conversationID, messageID string, totalBytes int) (*response.InitializeXChatMediaUploadResponse, error) {
	pl := (&payload.InitializeXChatMediaUploadPayload{}).Default()
	pl.Variables = payload.InitializeXChatMediaUploadVariables{
		ConversationID: conversationID,
		MessageID:      messageID,
		TotalBytes:     strconv.Itoa(totalBytes),
	}

	// Extract sha256 hash from endpoint URL path
	u, err := url.Parse(endpoints.INITIALIZE_XCHAT_MEDIA_UPLOAD_URL)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint URL: %w", err)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("unexpected endpoint path: %s", u.Path)
	}
	pl.Extensions.PersistedQuery.Sha256Hash = parts[len(parts)-2]

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("InitializeXChatMediaUpload request")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.INITIALIZE_XCHAT_MEDIA_UPLOAD_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		RawJSON("response", respBody).
		Msg("InitializeXChatMediaUpload response")

	var resp response.InitializeXChatMediaUploadResponse
	return &resp, json.Unmarshal(respBody, &resp)
}

// uploadMediaBytes uploads media bytes to ton.x.com in 512KB chunks.
// Returns the number of parts uploaded.
func (c *Client) uploadMediaBytes(ctx context.Context, resumeUploadURL string, mediaBytes []byte) (int, error) {
	const chunkSize = 512 * 1024 // 512 KB chunks

	headerOpts := HeaderOpts{
		WithNonAuthBearer: true,
		WithXCsrfToken:    true,
		WithCookies:       true,
		Origin:            endpoints.BASE_URL,
		Referer:           endpoints.BASE_URL + "/",
		Extra: map[string]string{
			"accept": "*/*",
		},
	}
	headers := c.buildHeaders(headerOpts)
	headers.Set("Content-Type", "application/octet-stream")

	partNumber := 0
	for chunk := range slices.Chunk(mediaBytes, chunkSize) {
		uploadURL := endpoints.TON_UPLOAD_BASE_URL + "/i/ton/data" + resumeUploadURL + "&partNumber=" + strconv.Itoa(partNumber)

		resp, respBody, err := c.MakeRequest(ctx, uploadURL, http.MethodPost, headers, chunk, types.ContentTypeNone)
		if err != nil {
			return 0, fmt.Errorf("upload part %d: %w", partNumber, err)
		}

		if resp.StatusCode > 204 {
			return 0, fmt.Errorf("upload part %d failed (status_code=%d, response_body=%s)", partNumber, resp.StatusCode, string(respBody))
		}

		c.Logger.Debug().
			Int("part_number", partNumber).
			Int("chunk_size", len(chunk)).
			Msg("Uploaded media chunk")

		partNumber++
	}

	return partNumber, nil
}

// finalizeXChatMediaUpload finalizes an XChat media upload.
func (c *Client) finalizeXChatMediaUpload(ctx context.Context, conversationID, messageID, mediaHashKey, resumeID string, numParts int) (*response.FinalizeXChatMediaUploadResponse, error) {
	pl := (&payload.FinalizeXChatMediaUploadPayload{}).Default()
	pl.Variables = payload.FinalizeXChatMediaUploadVariables{
		ConversationID: conversationID,
		MessageID:      messageID,
		MediaHashKey:   mediaHashKey,
		ResumeID:       resumeID,
		NumParts:       strconv.Itoa(numParts),
		TTLMsec:        nil,
	}

	// Extract sha256 hash from endpoint URL path
	u, err := url.Parse(endpoints.FINALIZE_XCHAT_MEDIA_UPLOAD_URL)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint URL: %w", err)
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("unexpected endpoint path: %s", u.Path)
	}
	pl.Extensions.PersistedQuery.Sha256Hash = parts[len(parts)-2]

	jsonBody, err := json.Marshal(pl)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	c.Logger.Debug().
		RawJSON("payload", jsonBody).
		Msg("FinalizeXChatMediaUpload request")

	_, respBody, err := c.makeAPIRequest(ctx, apiRequestOpts{
		URL:            endpoints.FINALIZE_XCHAT_MEDIA_UPLOAD_URL,
		Method:         http.MethodPost,
		WithClientUUID: true,
		Origin:         endpoints.BASE_URL,
		ContentType:    types.ContentTypeJSON,
		Body:           jsonBody,
	})
	if err != nil {
		return nil, err
	}

	c.Logger.Debug().
		RawJSON("response", respBody).
		Msg("FinalizeXChatMediaUpload response")

	var resp response.FinalizeXChatMediaUploadResponse
	return &resp, json.Unmarshal(respBody, &resp)
}

// DownloadXChatMedia downloads and decrypts encrypted media from XChat.
// The media is fetched from ton.x.com and decrypted using secretstream with the conversation key.
// If keyVersion is empty, the latest conversation key will be used.
func (c *Client) DownloadXChatMedia(ctx context.Context, conversationID, mediaHashKey, keyVersion string) ([]byte, error) {
	// Get the conversation key for decryption
	var convKey *crypto.ConversationKey
	var err error
	if keyVersion != "" {
		convKey, err = c.keyManager.GetConversationKey(ctx, conversationID, keyVersion)
	} else {
		convKey, err = c.keyManager.GetLatestConversationKey(ctx, conversationID)
	}
	if err != nil {
		return nil, fmt.Errorf("get conversation key: %w", err)
	}

	// Construct download URL
	downloadURL := endpoints.TON_UPLOAD_BASE_URL + "/i/ton/data/xchat_media/" + conversationID + "/" + mediaHashKey

	c.Logger.Info().
		Str("download_url", downloadURL).
		Str("conversation_id", conversationID).
		Str("media_hash_key", mediaHashKey).
		Str("key_version", convKey.KeyVersion).
		Msg("Downloading XChat encrypted media")

	headerOpts := HeaderOpts{
		WithNonAuthBearer: true,
		WithXCsrfToken:    true,
		WithCookies:       true,
		Origin:            endpoints.BASE_URL,
		Referer:           endpoints.BASE_URL + "/",
		Extra: map[string]string{
			"accept": "*/*",
		},
	}
	headers := c.buildHeaders(headerOpts)

	resp, respBody, err := c.MakeRequest(ctx, downloadURL, http.MethodGet, headers, nil, types.ContentTypeNone)
	if err != nil {
		return nil, fmt.Errorf("download request: %w", err)
	}

	if resp.StatusCode > 204 {
		return nil, fmt.Errorf("download failed (status_code=%d, response_body=%s)", resp.StatusCode, string(respBody))
	}

	// Decrypt using secretstream (XChaCha20-Poly1305)
	plaintext, err := crypto.SecretstreamDecrypt(respBody, convKey.Key)
	if err != nil {
		return nil, fmt.Errorf("decrypt media: %w", err)
	}

	return plaintext, nil
}
