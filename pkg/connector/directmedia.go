package connector

import (
	"context"
	"io"

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
	mediaInfo, err := ParseMediaID(mediaID)
	if err != nil {
		return nil, err
	}
	zerolog.Ctx(ctx).Trace().Any("mediaInfo", mediaInfo).Any("err", err).Msg("download direct media")
	ul := tc.br.GetCachedUserLoginByID(mediaInfo.UserID)
	client := ul.Client.(*TwitterClient)
	return &mediaproxy.GetMediaResponseCallback{
		Callback: func(w io.Writer) (int64, error) {
			resp, err := client.downloadFile(ctx, mediaInfo.URL)
			if err != nil {
				return 0, err
			}
			return io.Copy(w, resp.Body)
		},
	}, nil
}
