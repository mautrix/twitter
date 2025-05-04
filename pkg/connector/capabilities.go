// mautrix-twitter - A Matrix-Twitter puppeting bridge.
// Copyright (C) 2025 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package connector

import (
	"context"
	"time"

	"go.mau.fi/util/ffmpeg"
	"go.mau.fi/util/jsontime"
	"go.mau.fi/util/ptr"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/event"
)

func (tc *TwitterConnector) GetCapabilities() *bridgev2.NetworkGeneralCapabilities {
	return &bridgev2.NetworkGeneralCapabilities{}
}

func (tc *TwitterConnector) GetBridgeInfoVersion() (info, caps int) {
	return 1, 3
}

const MaxTextLength = 10000

func supportedIfFFmpeg() event.CapabilitySupportLevel {
	if ffmpeg.Supported() {
		return event.CapLevelPartialSupport
	}
	return event.CapLevelRejected
}

func (tc *TwitterClient) GetCapabilities(_ context.Context, _ *bridgev2.Portal) *event.RoomFeatures {
	return &event.RoomFeatures{
		ID: "fi.mau.twitter.capabilities.2025_02_05",
		//Formatting: map[event.FormattingFeature]event.CapabilitySupportLevel{
		//	event.FmtUserLink: event.CapLevelFullySupported,
		//},
		File: event.FileFeatureMap{
			event.MsgImage: {
				MimeTypes: map[string]event.CapabilitySupportLevel{
					"image/jpeg": event.CapLevelFullySupported,
					"image/png":  event.CapLevelFullySupported,
					"image/gif":  event.CapLevelFullySupported,
					"image/webp": event.CapLevelFullySupported,
				},
				Caption:          event.CapLevelFullySupported,
				MaxCaptionLength: MaxTextLength,
				MaxSize:          5 * 1024 * 1024,
			},
			event.MsgVideo: {
				MimeTypes: map[string]event.CapabilitySupportLevel{
					"video/mp4":       event.CapLevelFullySupported,
					"video/quicktime": event.CapLevelFullySupported,
				},
				Caption:          event.CapLevelFullySupported,
				MaxCaptionLength: MaxTextLength,
				MaxSize:          15 * 1024 * 1024,
			},
			event.CapMsgVoice: {
				MimeTypes: map[string]event.CapabilitySupportLevel{
					"audio/aac": supportedIfFFmpeg(),
					"audio/ogg": supportedIfFFmpeg(),
				},
				Caption:          event.CapLevelFullySupported,
				MaxCaptionLength: MaxTextLength,
				MaxSize:          5 * 1024 * 1024,
			},
			event.CapMsgGIF: {
				MimeTypes: map[string]event.CapabilitySupportLevel{
					"image/gif": event.CapLevelFullySupported,
					"video/mp4": event.CapLevelFullySupported,
				},
				Caption:          event.CapLevelFullySupported,
				MaxCaptionLength: MaxTextLength,
				MaxSize:          5 * 1024 * 1024,
			},
		},

		MaxTextLength: MaxTextLength,

		Reply: event.CapLevelFullySupported,

		Edit:          event.CapLevelFullySupported,
		EditMaxCount:  10,
		EditMaxAge:    ptr.Ptr(jsontime.S(15 * time.Minute)),
		Reaction:      event.CapLevelFullySupported,
		ReactionCount: 1,
	}
}
