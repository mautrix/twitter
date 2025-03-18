package twittermeow

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"go.mau.fi/util/random"

	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/endpoints"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/payload"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/data/response"
	"go.mau.fi/mautrix-twitter/pkg/twittermeow/types"
)

func (c *Client) UploadMedia(params *payload.UploadMediaQuery, mediaBytes []byte) (*response.FinalizedUploadMediaResponse, error) {
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

	_, respBody, err := c.MakeRequest(url, http.MethodPost, headers, nil, types.NONE)
	if err != nil {
		return nil, err
	}

	initUploadResponse := &response.InitUploadMediaResponse{}
	err = json.Unmarshal(respBody, initUploadResponse)
	if err != nil {
		return nil, err
	}

	if mediaBytes != nil {
		appendMediaPayload, contentType, err := c.newMediaAppendPayload(mediaBytes)
		if err != nil {
			return nil, err
		}
		headers.Add("content-type", contentType)

		url = fmt.Sprintf("%s?command=APPEND&media_id=%s&segment_index=0", endpoints.UPLOAD_MEDIA_URL, initUploadResponse.MediaIDString)
		resp, respBody, err := c.MakeRequest(url, http.MethodPost, headers, appendMediaPayload, types.NONE)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode > 204 {
			return nil, fmt.Errorf("failed to append media bytes for media with id %s (status_code=%d, response_body=%s)", initUploadResponse.MediaIDString, resp.StatusCode, string(respBody))
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
		resp, respBody, err = c.MakeRequest(url, http.MethodPost, headers, nil, types.NONE)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode > 204 {
			return nil, fmt.Errorf("failed to finalize media with id %s (status_code=%d, response_body=%s)", initUploadResponse.MediaIDString, resp.StatusCode, string(respBody))
		}

		finalizedMediaResultBytes = respBody
	} else {
		_, finalizedMediaResultBytes, err = c.GetMediaUploadStatus(initUploadResponse.MediaIDString, headers)
		if err != nil {
			return nil, err
		}
	}

	finalizedMediaResult := &response.FinalizedUploadMediaResponse{}
	err = json.Unmarshal(finalizedMediaResultBytes, finalizedMediaResult)
	if err != nil {
		return nil, err
	}

	if finalizedMediaResult.ProcessingInfo.State == response.PROCESSING_STATE_PENDING || finalizedMediaResult.ProcessingInfo.State == response.PROCESSING_STATE_IN_PROGRESS {
		// might need to check for error processing states here, I have not encountered any though so I wouldn't know what they look like/are
		for finalizedMediaResult.ProcessingInfo.State != response.PROCESSING_STATE_SUCCEEDED {
			finalizedMediaResult, _, err = c.GetMediaUploadStatus(finalizedMediaResult.MediaIDString, headers)
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

func (c *Client) GetMediaUploadStatus(mediaID string, h http.Header) (*response.FinalizedUploadMediaResponse, []byte, error) {
	url := fmt.Sprintf("%s?command=STATUS&media_id=%s", endpoints.UPLOAD_MEDIA_URL, mediaID)
	resp, respBody, err := c.MakeRequest(url, http.MethodGet, h, nil, types.NONE)
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
