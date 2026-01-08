package connector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/mediaproxy"
)

var _ bridgev2.DirectMediableNetwork = (*TwitterConnector)(nil)

func (tc *TwitterConnector) SetUseDirectMedia() {
	tc.directMedia = true
}

func (tc *TwitterConnector) Download(ctx context.Context, mediaID networkid.MediaID, params map[string]string) (mediaproxy.GetMediaResponse, error) {
	parsed, err := ParseMediaID(mediaID)
	if err != nil {
		return nil, err
	}

	switch info := parsed.(type) {
	case *MediaInfo:
		return tc.downloadLegacyMedia(ctx, info)
	case *EncryptedMediaInfo:
		return tc.downloadEncryptedMedia(ctx, info)
	default:
		return nil, fmt.Errorf("unknown media info type: %T", parsed)
	}
}

// downloadLegacyMedia downloads non-encrypted media directly from Twitter's servers.
func (tc *TwitterConnector) downloadLegacyMedia(ctx context.Context, info *MediaInfo) (mediaproxy.GetMediaResponse, error) {
	zerolog.Ctx(ctx).Trace().Any("mediaInfo", info).Msg("download legacy direct media")
	ul := tc.br.GetCachedUserLoginByID(info.UserID)
	if ul == nil || !ul.Client.IsLoggedIn() {
		return nil, fmt.Errorf("no logged in user found")
	}
	client := ul.Client.(*TwitterClient)
	resp, err := downloadFile(ctx, client.client, info.URL)
	if err != nil {
		return nil, err
	}
	return &mediaproxy.GetMediaResponseData{
		Reader:        resp.Body,
		ContentType:   resp.Header.Get("content-type"),
		ContentLength: resp.ContentLength,
	}, nil
}

// downloadEncryptedMedia downloads and decrypts encrypted XChat media.
func (tc *TwitterConnector) downloadEncryptedMedia(ctx context.Context, info *EncryptedMediaInfo) (mediaproxy.GetMediaResponse, error) {
	log := zerolog.Ctx(ctx)
	log.Trace().
		Str("conversation_id", info.ConversationID).
		Str("media_hash_key", info.MediaHashKey).
		Str("key_version", info.KeyVersion).
		Msg("download encrypted direct media")

	ul := tc.br.GetCachedUserLoginByID(info.UserID)
	if ul == nil || !ul.Client.IsLoggedIn() {
		return nil, fmt.Errorf("no logged in user found")
	}
	client := ul.Client.(*TwitterClient)

	// Use existing DownloadXChatMedia which handles decryption
	decrypted, err := client.client.DownloadXChatMedia(ctx, info.ConversationID, info.MediaHashKey, info.KeyVersion)
	if err != nil {
		return nil, fmt.Errorf("download and decrypt XChat media: %w", err)
	}

	return &mediaproxy.GetMediaResponseData{
		Reader:        io.NopCloser(bytes.NewReader(decrypted)),
		ContentType:   http.DetectContentType(decrypted),
		ContentLength: int64(len(decrypted)),
	}, nil
}
