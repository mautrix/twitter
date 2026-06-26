package twittermeow

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/bits"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

const (
	currentCastleTokenPrefix   = "IBYIll|"
	currentCastleChecksumLabel = "Nqkju"
	currentCastlePayloadSlots  = 494

	currentCastleSlot6ComponentCount  = 47
	currentCastleSlot6LengthMask      = 0x7fff
	currentCastleSlot6LengthContinue  = 0x8000
	currentCastleSlot6LengthShiftBits = 15

	currentCastleAcceptedAutomationBitfield = 0x1353003b
	currentCastleDefaultLanguage            = "en-US"
	currentCastleDefaultLanguages           = "en-US,en"
	currentCastleDefaultOrientation         = "landscape-primary"
	currentCastleDefaultPlatform            = "Win32"
	currentCastleDefaultHost                = "x.com"
	currentCastleDefaultWebGLVendor         = "Google Inc."
	currentCastleDefaultWebGLRenderer       = "ANGLE (NVIDIA, NVIDIA GeForce GTX 1080 Ti (0x00001B06) Direct3D11 vs_5_0 ps_5_0, D3D11)"
	currentCastleDefaultFeatureBits         = "1111111111111111111111111111111111111000000000"
	currentCastleDefaultAudioProbe          = "20030107"
	currentCastleDefaultPrecisionProbe      = "1000"
	currentCastleDefaultNumericProbe        = "43530826428"
	currentCastleDefaultWASMProbeHex        = "000002c3c600000789c60007"
	currentCastleSlot6TimestampSkewMillis   = 16
	currentCastleDefaultHash5               = "39a01d2b"
	currentCastleDefaultHash7               = "bd4be70d"
	currentCastleDefaultHash10              = "89dd3552"
	currentCastleDefaultHash11              = "8fc90f7d"
	currentCastleDefaultHash16              = "5d55117f"
	currentCastleDefaultHash18              = "42d68d60"
	currentCastleDefaultHash23              = "34b454a7"
	currentCastleDefaultHash39              = "db77b610"
)

var currentCastleAcceptedNumericSlots = []int{
	2, 20, 26, 27, 28, 31, 42, 44, 47, 49, 53, 56, 64, 67, 68, 71, 72, 73,
	87, 90, 98, 102, 103, 104, 106, 108, 109, 115, 116, 121, 123, 130, 133,
	147, 149, 150, 152, 154, 156, 159, 168, 174, 176, 177, 180, 182, 183, 184,
	188, 194, 196, 198, 203, 205, 206, 207, 209, 212, 220, 221, 223, 224, 228,
	235, 237, 240, 245, 246, 247, 256, 257, 258, 260, 261, 266, 267, 274, 277,
	279, 281, 288, 289, 292, 293, 297, 298, 304, 308, 311, 312, 316, 318, 321,
	326, 329, 330, 334, 339, 341, 344, 348, 349, 354, 356, 359, 366, 369, 374,
	375, 376, 378, 383, 386, 388, 390, 393, 395, 399, 406, 421, 422, 425, 428,
	429, 430,
}

var currentCastleAcceptedObjectSlots = []int{
	41, 46, 77, 99, 100, 128, 173, 214, 290, 342, 352, 416, 424,
}

const currentCastleAcceptedBooleanSlot = 273

func createCurrentCastleWrappedToken(payloadJSON string, timestampMillis int64) (string, error) {
	if payloadJSON == "" {
		return "", fmt.Errorf("current Castle payload is empty")
	}
	timestamp := strconv.FormatInt(timestampMillis, 10)
	payloadWithChecksum := insertCurrentCastleChecksum(payloadJSON)
	encodedTimestamp := encodeCurrentCastleTimestamp(timestamp)
	xoredPayload := xorCurrentCastleString(payloadWithChecksum, timestamp)
	wire := string([]byte{byte(len(encodedTimestamp))}) + encodedTimestamp + xoredPayload
	compressed, err := deflateCurrentCastleWireString(wire)
	if err != nil {
		return "", err
	}
	return currentCastleTokenPrefix + base64.StdEncoding.EncodeToString(compressed), nil
}

func createCurrentCastleWrappedPayloadToken(payload []any, timestampMillis int64) (string, error) {
	payloadJSON, err := marshalCurrentCastlePayload(payload)
	if err != nil {
		return "", err
	}
	return createCurrentCastleWrappedToken(payloadJSON, timestampMillis)
}

// newCurrentCastlePayloadScaffold returns a null-free payload skeleton. Unknown
// slots remain empty until the corresponding active X Castle slot is ported.
func newCurrentCastlePayloadScaffold() []any {
	payload := make([]any, currentCastlePayloadSlots)
	for i := range payload {
		payload[i] = ""
	}
	for _, slot := range currentCastleAcceptedNumericSlots {
		payload[slot] = 0
	}
	for _, slot := range currentCastleAcceptedObjectSlots {
		payload[slot] = map[string]any{}
	}
	payload[currentCastleAcceptedBooleanSlot] = false
	return payload
}

type currentCastlePayloadInput struct {
	AutomationBitfield uint32
	Slot6Components    []string
	NumericValues      map[int]uint32
	LowerFloatValues   map[int]float64
	ArrayValues        map[int][]float64
	ObjectValues       map[int]map[string]any
	PackedStringValues map[int][]string
	UnitPackedStrings  map[int][]string
	StringValues       map[int]string
	HighTimingValues   map[int]float64
}

type currentCastleSlot6Fingerprint struct {
	TimestampMillis int64
	Timezone        string
	WebGLRenderer   string
	WebGLVendor     string
	Orientation     string
	UserAgent       string
	Language        string
	Languages       string
	Platform        string
	Host            string
	PublicKey       string
	ClientUUID      string
	FeatureBits     string
	PrecisionProbe  string
	AudioProbe      string
	ProbeDateString string
	NumericProbe    string
	WASMProbeHex    string

	Hash5  string
	Hash7  string
	Hash10 string
	Hash11 string
	Hash16 string
	Hash18 string
	Hash23 string
	Hash39 string
}

func createCurrentCastleRequestToken(clientUUID string) (string, error) {
	timestampMillis := time.Now().UnixMilli()
	components, err := buildCurrentCastleSlot6Components(defaultCurrentCastleSlot6Fingerprint(timestampMillis, clientUUID))
	if err != nil {
		return "", err
	}
	payload, err := buildCurrentCastlePayload(currentCastlePayloadInput{
		AutomationBitfield: currentCastleAcceptedAutomationBitfield,
		Slot6Components:    components,
	})
	if err != nil {
		return "", err
	}
	return createCurrentCastleWrappedPayloadToken(payload, timestampMillis)
}

func defaultCurrentCastleSlot6Fingerprint(timestampMillis int64, clientUUID string) currentCastleSlot6Fingerprint {
	if timestampMillis <= 0 {
		timestampMillis = time.Now().UnixMilli()
	}
	timezone := jetfuelTimezone()
	renderer := currentCastleDefaultWebGLRenderer
	return currentCastleSlot6Fingerprint{
		TimestampMillis: timestampMillis,
		Timezone:        timezone,
		WebGLRenderer:   renderer,
		WebGLVendor:     currentCastleDefaultWebGLVendor,
		Orientation:     currentCastleDefaultOrientation,
		UserAgent:       UserAgent,
		Language:        currentCastleDefaultLanguage,
		Languages:       currentCastleDefaultLanguages,
		Platform:        currentCastleDefaultPlatform,
		Host:            currentCastleDefaultHost,
		PublicKey:       castlePublicKey,
		ClientUUID:      clientUUID,
		FeatureBits:     currentCastleDefaultFeatureBits,
		PrecisionProbe:  currentCastleDefaultPrecisionProbe,
		AudioProbe:      currentCastleDefaultAudioProbe,
		ProbeDateString: currentCastleSlot6ProbeDateString(timezone),
		NumericProbe:    currentCastleDefaultNumericProbe,
		WASMProbeHex:    currentCastleDefaultWASMProbeHex,
		Hash5:           currentCastleDefaultHash5,
		Hash7:           currentCastleDefaultHash7,
		Hash10:          currentCastleDefaultHash10,
		Hash11:          currentCastleDefaultHash11,
		Hash16:          currentCastleDefaultHash16,
		Hash18:          currentCastleDefaultHash18,
		Hash23:          currentCastleDefaultHash23,
		Hash39:          currentCastleDefaultHash39,
	}
}

func buildCurrentCastleSlot6Components(fp currentCastleSlot6Fingerprint) ([]string, error) {
	return encodeCurrentCastleSlot6ComponentValues(buildCurrentCastleSlot6RawValues(fp))
}

func buildCurrentCastleSlot6RawValues(fp currentCastleSlot6Fingerprint) []string {
	fp = normalizeCurrentCastleSlot6Fingerprint(fp)
	values := make([]string, currentCastleSlot6ComponentCount)
	values[0] = "toString"
	values[1] = "TypeError: Cyclic __proto__ value"
	values[2] = fp.Timezone
	values[3] = fp.WebGLRenderer
	values[5] = fp.Hash5
	values[6] = fp.Orientation
	values[7] = fp.Hash7
	values[8] = "RangeError"
	values[9] = "[]"
	values[10] = fp.Hash10
	values[11] = fp.Hash11
	values[12] = "probably"
	values[13] = "maybe"
	values[14] = fp.FeatureBits
	values[15] = "Cannot read properties of undefined (reading 'b')"
	values[16] = fp.Hash16
	values[18] = fp.Hash18
	values[19] = fp.ClientUUID
	values[20] = fp.PrecisionProbe
	values[21] = fp.WebGLVendor
	values[22] = fp.WebGLRenderer
	values[23] = fp.Hash23
	values[24] = fp.AudioProbe
	values[25] = "{}"
	values[27] = fp.PublicKey
	values[28] = "Illegal invocation"
	values[29] = fp.UserAgent
	values[30] = "r:1"
	values[31] = fp.Language
	values[33] = strconv.FormatInt(currentCastleSlot6ComponentTimestampMillis(fp.TimestampMillis), 10)
	values[34] = fp.Language
	values[35] = "Maximum call stack size exceeded"
	values[36] = fp.Platform
	values[37] = fp.ProbeDateString
	values[38] = "r:1"
	values[39] = fp.Hash39
	values[41] = fp.Host
	values[43] = fp.WASMProbeHex
	values[44] = fp.Languages
	values[45] = "probably"
	values[46] = fp.NumericProbe
	return values
}

func normalizeCurrentCastleSlot6Fingerprint(fp currentCastleSlot6Fingerprint) currentCastleSlot6Fingerprint {
	if fp.TimestampMillis <= 0 {
		fp.TimestampMillis = time.Now().UnixMilli()
	}
	if fp.Timezone == "" {
		fp.Timezone = jetfuelTimezone()
	}
	if fp.WebGLRenderer == "" {
		fp.WebGLRenderer = currentCastleDefaultWebGLRenderer
	}
	if fp.WebGLVendor == "" {
		fp.WebGLVendor = currentCastleDefaultWebGLVendor
	}
	if fp.Orientation == "" {
		fp.Orientation = currentCastleDefaultOrientation
	}
	if fp.UserAgent == "" {
		fp.UserAgent = UserAgent
	}
	if fp.Language == "" {
		fp.Language = currentCastleDefaultLanguage
	}
	if fp.Languages == "" {
		fp.Languages = currentCastleDefaultLanguages
	}
	if fp.Platform == "" {
		fp.Platform = currentCastleDefaultPlatform
	}
	if fp.Host == "" {
		fp.Host = currentCastleDefaultHost
	}
	if fp.PublicKey == "" {
		fp.PublicKey = castlePublicKey
	}
	if fp.FeatureBits == "" {
		fp.FeatureBits = currentCastleDefaultFeatureBits
	}
	if fp.PrecisionProbe == "" {
		fp.PrecisionProbe = currentCastleDefaultPrecisionProbe
	}
	if fp.AudioProbe == "" {
		fp.AudioProbe = currentCastleDefaultAudioProbe
	}
	if fp.ProbeDateString == "" {
		fp.ProbeDateString = currentCastleSlot6ProbeDateString(fp.Timezone)
	}
	if fp.NumericProbe == "" {
		fp.NumericProbe = currentCastleDefaultNumericProbe
	}
	if fp.WASMProbeHex == "" {
		fp.WASMProbeHex = currentCastleDefaultWASMProbeHex
	}
	if fp.Hash5 == "" {
		fp.Hash5 = currentCastleDefaultHash5
	}
	if fp.Hash7 == "" {
		fp.Hash7 = currentCastleDefaultHash7
	}
	if fp.Hash10 == "" {
		fp.Hash10 = currentCastleDefaultHash10
	}
	if fp.Hash11 == "" {
		fp.Hash11 = currentCastleDefaultHash11
	}
	if fp.Hash16 == "" {
		fp.Hash16 = currentCastleDefaultHash16
	}
	if fp.Hash18 == "" {
		fp.Hash18 = currentCastleDefaultHash18
	}
	if fp.Hash23 == "" {
		fp.Hash23 = currentCastleDefaultHash23
	}
	if fp.Hash39 == "" {
		fp.Hash39 = currentCastleDefaultHash39
	}
	return fp
}

func currentCastleSlot6ComponentTimestampMillis(timestampMillis int64) int64 {
	if timestampMillis <= currentCastleSlot6TimestampSkewMillis {
		return timestampMillis
	}
	return timestampMillis - currentCastleSlot6TimestampSkewMillis
}

func currentCastleSlot6ProbeDateString(timezone string) string {
	switch timezone {
	case "America/Chicago":
		return "3/3/1970, 6:00:00 PM"
	case "America/Los_Angeles":
		return "3/3/1970, 4:00:00 PM"
	case "America/New_York":
		return "3/3/1970, 7:00:00 PM"
	case "America/Denver":
		return "3/3/1970, 5:00:00 PM"
	case "America/Anchorage":
		return "3/3/1970, 3:00:00 PM"
	case "Pacific/Honolulu":
		return "3/3/1970, 2:00:00 PM"
	default:
		return "3/4/1970, 12:00:00 AM"
	}
}

func buildCurrentCastlePayload(input currentCastlePayloadInput) ([]any, error) {
	payload := newCurrentCastlePayloadScaffold()
	if err := populateCurrentCastleAutomationSlot(payload, input.AutomationBitfield); err != nil {
		return nil, err
	}
	if err := populateCurrentCastleNumericSlots(payload, input.NumericValues); err != nil {
		return nil, err
	}
	if err := populateCurrentCastleLowerFloatSlots(payload, input.LowerFloatValues); err != nil {
		return nil, err
	}
	if err := populateCurrentCastleArraySlots(payload, input.ArrayValues); err != nil {
		return nil, err
	}
	if err := populateCurrentCastlePackedStringSlots(payload, input.PackedStringValues); err != nil {
		return nil, err
	}
	if err := populateCurrentCastleUnitPackedStringSlots(payload, input.UnitPackedStrings); err != nil {
		return nil, err
	}
	if err := populateCurrentCastleStringSlots(payload, input.StringValues); err != nil {
		return nil, err
	}
	if err := populateCurrentCastleObjectSlots(payload, input.ObjectValues); err != nil {
		return nil, err
	}
	slot6, err := encodeCurrentCastleSlot6(input.Slot6Components)
	if err != nil {
		return nil, err
	}
	payload[6] = slot6
	if err = populateCurrentCastleHighTimingSlots(payload, input.HighTimingValues); err != nil {
		return nil, err
	}
	populateCurrentCastleObservedStringSlots(payload, input)
	if _, err = marshalCurrentCastlePayload(payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func populateCurrentCastleAutomationSlot(payload []any, rawBitfield uint32) error {
	if len(payload) <= 2 {
		return fmt.Errorf("current Castle payload has %d slots, want at least 3", len(payload))
	}
	payload[2] = currentCastleTB(rawBitfield)
	return nil
}

var currentCastleNumericSlotEncoders = map[int]func(uint32) uint32{
	20:  encodeCurrentCastleNumericSlot20,
	26:  currentCastleNumericAffineEncoder(currentCastleIndexR7(0x8dadb484), currentCastleIndexR7(0xa02eff37), 0x6ba90087),
	27:  currentCastleNumericAffineEncoder(0xf4730f04, 0x5aa87861, 0xe1ed70a8),
	28:  currentCastleNumericAffineEncoder(0x0d877048, 0x49b0e02d, currentCastleIndexR7(0x058e30f8)),
	31:  currentCastleNumericAffineEncoder(0xc3b6e1f8, 0xf35393df, currentCastleIndexTG(0x0fb221a7)),
	42:  currentCastleNumericAffineEncoder(0xeaf85686, 0xf6253667, currentCastleIndexR7(0x2a606d93)),
	44:  encodeCurrentCastleNumericIdentity,
	47:  encodeCurrentCastleNumericIdentity,
	49:  currentCastleNumericAffineEncoder(0x14ba9b74, 0x6381aa33, 0xc8fde421),
	53:  currentCastleNumericAffineEncoder(0xf967d101, 0x6627399f, 0x4badc824),
	56:  currentCastleNumericAffineEncoder(0x0fdaf2c8, 0xd32ce769, 0x35531149),
	64:  currentCastleNumericAffineEncoder(0x40a5541c, currentCastleIndexR7(0x187ef987), currentCastleIndexTG(0xca19a6f1)),
	67:  currentCastleNumericAffineEncoder(0x5e57a336, 0x4cc0a093, 0x90ebbe3b),
	68:  currentCastleNumericAffineEncoder(0x93d2c136, currentCastleIndexTG(0xfeae7a21), currentCastleIndexR7(0xcc851a80)),
	71:  currentCastleNumericAffineEncoder(currentCastleIndexTG(0xe5c05ac6), currentCastleIndexTG(0xd9dea0bd), 0x2ba0704d),
	72:  encodeCurrentCastleNumericIdentity,
	73:  currentCastleNumericAffineEncoder(0xcac4cfca, 0x19c3286b, 0x1540a806),
	87:  currentCastleNumericAffineEncoder(0xb438dadb, 0xbd6a0d4f, currentCastleIndexTG(0x4e27a0eb)),
	90:  encodeCurrentCastleNumericIdentity,
	98:  currentCastleNumericAffineEncoder(0x5f9d021c, currentCastleIndexR7(0x2171a35d), currentCastleIndexR7(0x8016951b)),
	102: currentCastleNumericAffineEncoder(0xe4cc9400, currentCastleIndexR7(0x1d8d20c9), 0xfbf84330),
	103: encodeCurrentCastleNumericIdentity,
	104: encodeCurrentCastleNumericIdentity,
	106: currentCastleNumericAffineEncoder(currentCastleIndexR7(0x95b45d4c), 0xfcc3ad03, 0x67847013),
	108: encodeCurrentCastleNumericIdentity,
	109: encodeCurrentCastleNumericIdentity,
	115: currentCastleNumericAffineEncoder(0xe33fe88d, 0x5aca0647, 0x58d3320e),
	116: currentCastleNumericAffineEncoder(0xd726d1f7, currentCastleIndexR7(0x2416a9d1), currentCastleIndexTG(0x2fca6a65)),
	121: currentCastleNumericAffineEncoder(0x78e89baf, currentCastleIndexR7(0x6f76f97e), 0x1639fdd5),
	123: currentCastleNumericAffineEncoder(0x0472b08c, currentCastleIndexR7(0x64137da7), currentCastleIndexR7(0x5cf504e4)),
	130: encodeCurrentCastleNumericIdentity,
	133: currentCastleNumericAffineEncoder(0x4aa2215c, 0x0f3c0647, 0x9cb82d4f),
	147: currentCastleNumericAffineEncoder(0x92f56e95, 0x5feb3bd7, 0x0308f549),
	149: currentCastleNumericAffineEncoder(0x2a8c6716, 0xb94477db, 0x496e2ec4),
	150: currentCastleNumericAffineEncoder(0x7e49a49f, 0xc276ed9b, 0x1ad0a021),
	152: currentCastleNumericAffineEncoder(0xc60062fb, 0xf6dd85bd, 0xf1d97dda),
	154: encodeCurrentCastleNumericSlot154,
	156: encodeCurrentCastleNumericIdentity,
	159: encodeCurrentCastleNumericIdentity,
	168: currentCastleNumericAffineEncoder(0xc99ace12, currentCastleIndexR7(0x30387440), currentCastleIndexR7(0x8e8ca6a7)),
	174: encodeCurrentCastleNumericIdentity,
	176: encodeCurrentCastleNumericIdentity,
	177: encodeCurrentCastleNumericIdentity,
	180: encodeCurrentCastleNumericIdentity,
	182: currentCastleTB,
	183: encodeCurrentCastleNumericIdentity,
	184: encodeCurrentCastleNumericIdentity,
	188: encodeCurrentCastleNumericIdentity,
	194: currentCastleNumericAffineEncoder(0xa3e41d03, 0xd468aef7, 0x4fc32d5a),
	196: currentCastleNumericAffineEncoder(0x338ca7e4, 0x004c6c1b, 0xcf397719),
	198: currentCastleNumericAffineEncoder(currentCastleIndexTG(0x3a6b8b4a), currentCastleIndexR7(0x15997926), 0x109ad9c4),
	203: currentCastleNumericAffineEncoder(currentCastleIndexTG(0xbeff5831), currentCastleIndexR7(0x1bb9f08f), currentCastleIndexR7(0xf15b90bc)),
	205: currentCastleNumericAffineEncoder(0xafbd8da7, 0xdf5818a9, 0x532b20e8),
	206: currentCastleNumericAffineEncoder(0xe27d001f, 0xc7a3e3bf, 0x094f1534),
	207: currentCastleNumericAffineEncoder(currentCastleIndexTG(0x553a4119), currentCastleIndexR7(0xcf68a2a6), 0x7342df93),
	209: currentCastleNumericAffineEncoder(0xc970ba2f, currentCastleIndexR7(0x95562fdc), 0x396c0bfc),
	212: currentCastleNumericAffineEncoder(0xa84e4ee0, currentCastleIndexTG(0x8e122ec3), 0x19f6e853),
	220: encodeCurrentCastleNumericIdentity,
	221: currentCastleNumericAffineEncoder(0xaf9fb005, currentCastleIndexR7(0xd6326bbb), 0x0cd83feb),
	223: currentCastleNumericAffineEncoder(0x17dbb40b, currentCastleIndexTG(0xb13dfc57), currentCastleIndexR7(0x8ae57b7f)),
	224: currentCastleNumericAffineEncoder(currentCastleIndexTG(0xe9ad1a81), 0x4c0f0ccd, 0x7651692c),
	228: currentCastleNumericAffineEncoder(0x081ff5ca, currentCastleIndexR7(0x0b50b5e1), 0xddfdc282),
	235: encodeCurrentCastleNumericIdentity,
	237: encodeCurrentCastleNumericIdentity,
	240: currentCastleNumericAffineEncoder(0x15608129, currentCastleIndexR7(0x7f143c5a), 0x31aaf010),
	245: encodeCurrentCastleNumericIdentity,
	246: currentCastleNumericAffineEncoder(0x2da919dc, 0x24b26a51, currentCastleIndexTG(0x8db5b1b5)),
	247: currentCastleNumericAffineEncoder(0xf1398f6f, 0x031a5829, 0xa39c022e),
	256: encodeCurrentCastleNumericSlot256,
	257: currentCastleNumericAffineEncoder(0x63d3c42d, currentCastleIndexTG(0xa499760d), 0xb61b007f),
	258: currentCastleNumericAffineEncoder(currentCastleIndexR7(0xddb05bd9), currentCastleIndexTG(0xda4f6e1b), currentCastleIndexR7(0x44f87d7f)),
	260: currentCastleNumericAffineEncoder(0x5c2d1334, 0x82f4555d, 0xfc40a036),
	261: encodeCurrentCastleNumericIdentity,
	266: currentCastleNumericAffineEncoder(currentCastleIndexR7(0x6c2338c7), currentCastleIndexTG(0x5e089af3), 0x47add5ce),
	267: currentCastleNumericAffineEncoder(currentCastleIndexTG(0x82a9dbff), 0x7e034531, 0xc32531d9),
	274: encodeCurrentCastleNumericIdentity,
	277: currentCastleNumericAffineEncoder(0x87fda9a5, currentCastleIndexTG(0x061ad74d), 0x97429312),
	279: encodeCurrentCastleNumericIdentity,
	281: encodeCurrentCastleNumericIdentity,
	289: currentCastleTB,
	288: currentCastleNumericAffineEncoder(0x7539f159, 0xdd1ef321, currentCastleIndexR7(0x0de88265)),
	292: encodeCurrentCastleNumericIdentity,
	293: currentCastleNumericAffineEncoder(0xb2bb9131, 0x3ddd3435, 0x808c247f),
	297: currentCastleNumericAffineEncoder(0x6c4cf908, 0xabee7b25, 0x00f0dead),
	298: encodeCurrentCastleNumericIdentity,
	304: encodeCurrentCastleNumericIdentity,
	308: encodeCurrentCastleNumericIdentity,
	311: currentCastleNumericAffineEncoder(0xab4e6a90, 0xb2896847, 0xbddd98cc),
	312: currentCastleNumericAffineEncoder(currentCastleIndexTG(0x7d8c8409), 0xb2c0779f, 0xdf29574a),
	316: currentCastleIndexEUnderscore,
	318: currentCastleNumericAffineEncoder(0xb70bef86, 0xa653b485, currentCastleIndexTG(0x1921c4cc)),
	321: currentCastleNumericAffineEncoder(0x022eeaef, 0x42e2bfb5, 0x25d5ac75),
	326: encodeCurrentCastleNumericIdentity,
	329: currentCastleNumericAffineEncoder(0xf7a020e8, 0xfb500465, 0x8f5dd1f8),
	330: encodeCurrentCastleNumericIdentity,
	334: currentCastleNumericAffineEncoder(0x073bb038, currentCastleIndexR7(0x5e87765b), currentCastleIndexTG(0xc7d82a02)),
	339: currentCastleNumericAffineEncoder(0x31a8b928, 0x74c029f1, 0x39a0b04b),
	341: encodeCurrentCastleNumericSlot341,
	344: encodeCurrentCastleNumericIdentity,
	348: encodeCurrentCastleNumericSlot348,
	349: currentCastleNumericAffineEncoder(0x29789f9c, 0x81d1cfb9, 0x5fafa5dd),
	354: currentCastleNumericAffineEncoder(0xa5d12073, 0x5a141b6f, 0x5c04da70),
	356: currentCastleNumericAffineEncoder(currentCastleIndexTG(0xefdfa846), 0x7479d8cd, 0xb5e1823b),
	359: currentCastleNumericAffineEncoder(0x5bca9e62, 0x7115f42f, currentCastleIndexTG(0xd01dd90c)),
	366: currentCastleNumericAffineEncoder(currentCastleIndexTG(0xd819e362), 0x1865a9bb, 0x4a8e00a2),
	369: currentCastleNumericAffineEncoder(currentCastleIndexTG(0x6252cdfb), 0x77cf48f5, 0x16577410),
	374: currentCastleNumericAffineEncoder(0x5ad18bb4, currentCastleIndexR7(0x35f9af04), 0xa86efa9e),
	375: currentCastleNumericAffineEncoder(0x52aa8eda, currentCastleIndexTG(0x79d65887), 0xf14f9c27),
	376: currentCastleNumericAffineEncoder(0x9118a572, 0xc1804b39, currentCastleIndexR7(0xd0d41fd7)),
	378: encodeCurrentCastleNumericIdentity,
	383: encodeCurrentCastleNumericIdentity,
	386: currentCastleNumericAffineEncoder(currentCastleIndexTG(0x8994cbfe), currentCastleIndexR7(0xa78374d8), 0xaef88e50),
	388: currentCastleNumericAffineEncoder(0x5f1db9d0, 0x48572a5d, 0x58b3854a),
	390: currentCastleNumericAffineEncoder(0x3d4f1683, 0x97d2840b, 0x754b73a7),
	393: currentCastleNumericAffineEncoder(0x569dbc0a, 0xa03ef08d, 0x9efe8725),
	395: currentCastleNumericAffineEncoder(0xd1f0afd2, 0x456a2621, 0x8313e4bc),
	399: currentCastleNumericAffineEncoder(0x2d54860a, 0x68fac53f, currentCastleIndexR7(0x060755c2)),
	406: currentCastleNumericAffineEncoder(0x7837ed09, 0xdc5b6193, 0xf2dbed27),
	421: encodeCurrentCastleNumericIdentity,
	422: currentCastleNumericAffineEncoder(currentCastleIndexTG(0x7f83d97a), 0xa935b369, currentCastleIndexR7(0x45368801)),
	425: currentCastleNumericAffineEncoder(0x74c852e2, 0x23ff47e7, 0x29394da4),
	428: encodeCurrentCastleNumericSlot428,
	429: encodeCurrentCastleNumericSlot429,
	430: currentCastleNumericAffineEncoder(0x47d64b14, currentCastleIndexTG(0x98f93ed5), 0x464a9053),
}

func defaultCurrentCastleNumericSlotValues() map[int]uint32 {
	return map[int]uint32{
		// Browser-like defaults: the media query does not match, navigator.webdriver
		// is false, and the T4/T6 counter probes fall back to zero on errors.
		20: 0,
		26: 0,
		27: 0,
		28: 0,
		// navigator.permissions exists in normal Chrome and is wrapped through tB.
		31: currentCastleTB(1),
		// The input-type support map is initialized false and only read by the
		// active bundle. tel/text/password/email all wrap false through e_.
		42: currentCastleIndexEUnderscore(0),
		// Pending formula ports below use final numeric values observed from
		// the current visible Chrome begin_login token on this Windows machine.
		44: 0x0637b4dd,
		47: 0x5b4fa7aa,
		// EO.length from the current Chrome begin_login token on this machine.
		49: 253,
		53: 0,
		56: 0,
		64: currentCastleIndexEUnderscore(0),
		// The active Gi-gated slot-67 branch falls back to the e_ wrapper's false value.
		67: currentCastleIndexEUnderscore(0),
		68: 0,
		71: currentCastleIndexEUnderscore(0),
		72: 0xe5371923,
		73: 0,
		87: currentCastleIndexEUnderscore(0),
		90: 0x443a1261,
		98: 0,
		// EX[1] from the current Chrome begin_login token on this machine.
		102: 64,
		103: 0x975b0b86,
		104: 0xddc55158,
		// navigator.credentials exists in normal Chrome and is wrapped through tB.
		106: currentCastleTB(1),
		108: 0x810588a6,
		109: 0x4901e987,
		115: 0,
		// pX() returns window.outerWidth - window.innerWidth. Keep the default
		// neutral and allow live Chrome-captured deltas to override it.
		116: 0,
		// Ho(r) returns 0 in the no-event/default path.
		121: 0,
		// Hv() is false unless the window/visualViewport dimensions match the
		// full screen, then wraps through the e_ transform.
		123: currentCastleIndexEUnderscore(0),
		130: 0x605ff24d,
		// navigator.javaEnabled() is present but false in normal Chrome.
		133: currentCastleTB(0),
		147: 0,
		149: 0,
		150: 0,
		152: 0,
		154: 0,
		156: 0x96efd7bf,
		// HM() is a window flag probe; the browser-like missing flag branch is false.
		159: 21203,
		168: 0,
		174: 0x6ae4238a,
		176: 0x4fa64045,
		177: 0x834bc08c,
		180: 0x49eaa0e7,
		// Ta/QC/TH all default false, so only the inverted TH bit is set.
		182: 8,
		183: 0x8894a25b,
		184: 0x1851c3d2,
		188: 0xaad37e98,
		194: 0,
		196: 0,
		// The Gi-gated screen/window branch is disabled in this Go port, so
		// slot 198 uses the tB(false) wrapper default.
		198: currentCastleTB(0),
		203: 0,
		205: 0,
		206: 0,
		207: 0,
		209: 0,
		// Slot 212 is window.outerHeight - window.innerHeight.
		212: 0,
		220: 0x960e16fd,
		221: 0,
		223: 0,
		224: currentCastleIndexEUnderscore(0),
		// The Gi-gated storage branch is false in the current Go port, so slot
		// 228 takes f0(false), which is the e_ wrapper's false value.
		228: currentCastleIndexEUnderscore(0),
		235: 0x0350a000,
		237: 0x522c3f07,
		240: 0,
		245: 0x5abc844a,
		246: 0,
		247: currentCastleIndexEUnderscore(0),
		256: 0,
		// iR(false) returns tB(false) before slot 257's affine mix.
		257: currentCastleTB(0),
		258: 0,
		260: 0,
		261: 0xa9ca285e,
		266: 0,
		// No-event/default callbacks for the QT/Qu event probes.
		267: 0,
		274: 0x330f0003,
		277: 0,
		279: 0xe1895c5b,
		281: 0x4877c85d,
		288: 0,
		// fD(false) returns tB(false), then slot 289 runs that through tB again.
		289: currentCastleTB(0),
		292: 0x521f133b,
		293: 0,
		// No-event defaults for GP/G_/G$ event-array probes.
		297: 0,
		298: 0xa4d06c96,
		304: 0x78571199,
		308: 0xe48da40e,
		311: 0,
		312: 0,
		// window.external is present in Chrome but its native toString() does not
		// contain the Sequentum marker; QN wraps false through tB, then slot 316
		// applies the e_ index transform inline.
		316: currentCastleTB(0),
		// Qy/QU return the e_ wrapper's false value in the no-event path.
		318: currentCastleIndexEUnderscore(0),
		321: currentCastleIndexEUnderscore(0),
		326: 0x1350806a,
		329: 0,
		330: 0xdd2ae869,
		// FM.hidden is false in the static input-type support map.
		334: currentCastleTB(0),
		339: 0,
		// p9() is another missing window flag probe; the default branch is false.
		341: 0,
		344: 0x58c2cc46,
		348: 0,
		349: 0,
		354: 0,
		356: currentCastleTB(0),
		// Normal Chrome reports navigator.webdriver as false.
		359: currentCastleTB(0),
		// Hn() returns 0 unless the active aV() probe throws a named error.
		366: 0,
		// navigator.pdfViewerEnabled exists and is true in normal Chrome, so the
		// active predicate `pdfViewerEnabled in navigator && false === value` is false.
		369: currentCastleIndexEUnderscore(0),
		374: 0,
		// QD() reads navigator.webdriver and wraps the normal false value through tB.
		375: currentCastleTB(0),
		// CanvasRenderingContext2D.prototype.getImageData is native in Chrome, so
		// the active anti-tamper predicate is false and wraps through e_.
		376: currentCastleIndexEUnderscore(0),
		378: 0x3979e7ca,
		383: 0x82c994ec,
		386: 0,
		388: 0,
		390: 0,
		393: 0,
		395: 0,
		// Qf()/io() current Chrome default from the visible begin_login path.
		399: 2770,
		406: 0,
		421: 0x749316dd,
		// screen.availHeight from the current Chrome profile on this machine.
		422: 1152,
		425: 0,
		// The active Castle closure shadows tp with an array before the slot-428
		// assignment, so the bundle takes the typeof-not-function sentinel.
		428: 0xdae88a03,
		// Hu() falls back to e_(false) when the canvas copy probe is unavailable.
		429: currentCastleIndexEUnderscore(0),
		// Ha returns the tB(false) wrapper when no event/mn array is present.
		430: currentCastleTB(0),
	}
}

func populateCurrentCastleNumericSlots(payload []any, values map[int]uint32) error {
	if len(payload) != currentCastlePayloadSlots {
		return fmt.Errorf("current Castle payload has %d slots, want %d", len(payload), currentCastlePayloadSlots)
	}
	defaults := defaultCurrentCastleNumericSlotValues()
	finalDefaults := defaultCurrentCastleNumericSlotFinalValues()
	for slot, encoder := range currentCastleNumericSlotEncoders {
		if _, hasOverride := values[slot]; !hasOverride {
			if finalValue, ok := finalDefaults[slot]; ok {
				payload[slot] = finalValue
				continue
			}
		}
		value, ok := defaults[slot]
		if override, hasOverride := values[slot]; hasOverride {
			value = override
			ok = true
		}
		if !ok {
			return fmt.Errorf("current Castle numeric slot %d has no default", slot)
		}
		payload[slot] = encoder(value)
	}
	return nil
}

func defaultCurrentCastleNumericSlotFinalValues() map[int]uint32 {
	return map[int]uint32{
		// Current Chrome begin_login numeric payload defaults observed from the
		// native Windows profile. Slot 2 is populated by the automation encoder.
		20:  1607583818,
		26:  3643181490,
		27:  1696677852,
		28:  609078933,
		31:  1676013394,
		42:  3525692521,
		44:  104314077,
		47:  1560894515,
		49:  4152372946,
		53:  2655580655,
		56:  3433073095,
		64:  1726107147,
		67:  3800479293,
		68:  3971848222,
		71:  1354671007,
		72:  3845593379,
		73:  3755269648,
		87:  255423802,
		90:  1144656481,
		98:  1125315884,
		102: 3724615792,
		103: 3143232709,
		104: 3720696152,
		106: 1535610710,
		108: 2164623526,
		109: 2608958949,
		115: 4005085634,
		116: 697758972,
		121: 2010197273,
		123: 1687974328,
		130: 1616900685,
		133: 2388381063,
		147: 4144616434,
		149: 1675400182,
		150: 1577344424,
		152: 4164810793,
		154: 4047298945,
		156: 2532300735,
		159: 21203,
		168: 1297108559,
		174: 1793336202,
		176: 1336295493,
		177: 2202779788,
		180: 1240113383,
		182: 1686,
		183: 2291442267,
		184: 408011730,
		188: 2865987224,
		194: 3796548060,
		196: 2406903333,
		198: 2949066494,
		203: 2202744921,
		205: 2658569404,
		206: 261761365,
		207: 2341113330,
		209: 4275254847,
		212: 469580049,
		220: 2517505789,
		221: 3617915702,
		223: 4066925288,
		224: 1663298151,
		228: 2333063244,
		235: 55615488,
		237: 1378631431,
		240: 2095554485,
		245: 1522304074,
		246: 4049567695,
		247: 1912840949,
		256: 1319576667,
		257: 2968100594,
		258: 1031170005,
		260: 3845770389,
		261: 2848598110,
		266: 2747539663,
		267: 724571498,
		274: 856621059,
		277: 3017337424,
		279: 3783875675,
		281: 450570634,
		288: 2851106031,
		289: 119280,
		292: 1377768251,
		293: 682616773,
		297: 2166863061,
		298: 2765122710,
		304: 2018972057,
		308: 3834487822,
		311: 1812913062,
		312: 1106886207,
		316: 62586880,
		318: 269607942,
		321: 2331875440,
		326: 324042858,
		329: 1748866603,
		330: 3710576745,
		334: 1588140542,
		339: 1910411492,
		341: 63143,
		344: 1489161286,
		348: 1412005389,
		349: 2996412825,
		354: 1192889677,
		356: 1432694549,
		359: 533304594,
		366: 3605025046,
		369: 2077348073,
		374: 3870528754,
		375: 3354097035,
		376: 1590110334,
		378: 964290506,
		383: 2194248940,
		386: 693294050,
		388: 1338730888,
		390: 649623368,
		393: 331217063,
		395: 3022502606,
		399: 825292510,
		406: 2109928018,
		421: 1955796701,
		422: 1587800882,
		425: 2383857042,
		428: 3942703131,
		429: 1915333822,
		430: 3276067091,
	}
}

var currentCastleObjectSlotEncoders = map[int]func([]any) (map[string]any, error){
	41:  encodeCurrentCastleObjectSlot41,
	77:  encodeCurrentCastleObjectSlot77,
	99:  encodeCurrentCastleObjectSlot99,
	100: encodeCurrentCastleObjectSlot100,
	46:  encodeCurrentCastleObjectSlot46,
	128: encodeCurrentCastleObjectSlot128,
	173: encodeCurrentCastleObjectSlot173,
	214: encodeCurrentCastleObjectSlot214,
	290: encodeCurrentCastleObjectSlot290,
	342: encodeCurrentCastleObjectSlot342,
	352: encodeCurrentCastleObjectSlot352,
	416: encodeCurrentCastleObjectSlot416,
	424: encodeCurrentCastleObjectSlot424,
}

var currentCastleObjectSlotOrder = []int{
	41,
	77,
	99,
	100,
	46,
	128,
	173,
	214,
	290,
	342,
	352,
	416,
	424,
}

func populateCurrentCastleObjectSlots(payload []any, values map[int]map[string]any) error {
	if len(payload) != currentCastlePayloadSlots {
		return fmt.Errorf("current Castle payload has %d slots, want %d", len(payload), currentCastlePayloadSlots)
	}
	finalDefaults := defaultCurrentCastleObjectSlotFinalValues()
	for _, slot := range currentCastleObjectSlotOrder {
		encoder := currentCastleObjectSlotEncoders[slot]
		if override, ok := values[slot]; ok {
			payload[slot] = override
			continue
		}
		if finalValue, ok := finalDefaults[slot]; ok {
			payload[slot] = finalValue
			continue
		}
		encoded, err := encoder(payload)
		if err != nil {
			return fmt.Errorf("current Castle object slot %d: %w", slot, err)
		}
		payload[slot] = encoded
	}
	return nil
}

func defaultCurrentCastleObjectSlotFinalValues() map[int]map[string]any {
	return map[int]map[string]any{
		// Current Chrome begin_login object payload defaults observed from the
		// native Windows profile.
		41: {
			"V":  "Ls9BGR3QrImUhLGh8DdWMQ==",
			"wh": "",
			"y":  "",
			"G":  "",
			"Hrc": map[string]any{
				"j": "",
			},
		},
		46: {
			"mx": "",
			"h":  "",
			"r":  "",
			"AI": "",
		},
		77: {
			"AFS": "JcW/1uh5mugs0/qG3eXd5Q==",
			"tu":  "",
			"Xd":  "",
			"BgV": "",
		},
		99: {
			"n":   "",
			"cOw": "",
		},
		100: {
			"rj": "5e4bQb/W",
			"KwJ": map[string]any{
				"T": "CzejbOYagxFvwg==",
			},
			"L": map[string]any{
				"T": "",
			},
			"E": "",
		},
		128: {
			"ekC": "jM4d0BtBVjEbQeXusaE=",
			"fCH": "8Dc=",
		},
		173: {
			"t":   "gY/mGnrO",
			"JZX": "rImUhLGh8DdWMQ==",
			"ga": map[string]any{
				"sE":  "",
				"OkL": "",
				"hD":  "",
			},
			"wPg": "",
			"FqX": "",
		},
		214: {
			"h":   "LRTQyw==",
			"yhj": "S+hTpr/W8DdL6PA3",
			"XI":  "",
		},
		290: {
			"d": "LRRvwgs3+WmjbOYaCzcMuQ==",
		},
		342: {
			"kFx": "CzfQy+YaCzdac9DLLRSQHQ==",
		},
		352: {
			"XYw": "gY/mGnrO0MsLN5Ad9+ejbA==",
		},
		416: {
			"y":   "LRR6zlpzo2xac/lpkB335w==",
			"OiJ": "",
			"HRd": map[string]any{
				"j":   "",
				"hCP": "",
			},
		},
		424: {
			"kr": "jM4d0BtBVjE=",
			"vj": map[string]any{
				"GXV": "G0Hl7rGh8Dc=",
			},
			"RK": "",
			"wXd": map[string]any{
				"RtN": map[string]any{
					"o":   "",
					"ApQ": "",
					"G":   "",
				},
			},
		},
	}
}

func encodeCurrentCastleObjectSlot41(payload []any) (map[string]any, error) {
	units, err := currentCastleTQHashUnitsForSlot(payload, 40)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"V":  encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 0, 9)),
		"wh": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 9, 15)),
		"y":  encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 15, 16)),
		"G":  encodeCurrentCastleUnits(currentCastleSliceUnits(units, 16, 21)),
		"Hrc": map[string]any{
			"j": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 21, len(units))),
		},
	}, nil
}

func encodeCurrentCastleObjectSlot290(payload []any) (map[string]any, error) {
	return encodeCurrentCastleTQTTObjectSlot(payload, 289, "d")
}

func encodeCurrentCastleObjectSlot352(payload []any) (map[string]any, error) {
	return encodeCurrentCastleTQTTObjectSlot(payload, 351, "XYw")
}

func encodeCurrentCastleObjectSlot46(_ []any) (map[string]any, error) {
	return encodeCurrentCastleObjectSlot46FromMatches(nil), nil
}

func encodeCurrentCastleObjectSlot46FromMatches(matches []string) map[string]any {
	source := ""
	if len(matches) > 0 {
		source = currentCastleSlot6HashHex(strings.Join(matches, ","))
	}
	encoded := encodeCurrentCastleLR(source)
	units := jsUTF16Units(encoded)
	return map[string]any{
		"mx": encodeCurrentCastleUnits(currentCastleSliceUnits(units, 0, 8)),
		"h":  encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 8, 14)),
		"r":  encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 14, 16)),
		"AI": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 16, len(units))),
	}
}

func encodeCurrentCastleObjectSlot342(payload []any) (map[string]any, error) {
	return encodeCurrentCastleTITTObjectSlot(payload, 341, "kFx")
}

func encodeCurrentCastleObjectSlot77(_ []any) (map[string]any, error) {
	return encodeCurrentCastleObjectSlot77FromSource("[]"), nil
}

func encodeCurrentCastleObjectSlot77FromSource(source string) map[string]any {
	token := encodeCurrentCastleUnits(jsUTF16Units(source))
	units := jsUTF16Units(token)
	return map[string]any{
		"AFS": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 0, 10)),
		"tu":  encodeCurrentCastleUnits(currentCastleSliceUnits(units, 10, 16)),
		"Xd":  encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 16, 24)),
		"BgV": encodeCurrentCastleUnits(currentCastleSliceUnits(units, 24, len(units))),
	}
}

func encodeCurrentCastleObjectSlot99(_ []any) (map[string]any, error) {
	return encodeCurrentCastleObjectSlot99FromCanPlayType(""), nil
}

func encodeCurrentCastleObjectSlot99FromCanPlayType(canPlayType string) map[string]any {
	return encodeCurrentCastleObjectSlot99FromSource(encodeCurrentCastleRawPackedUnits(jsUTF16Units(canPlayType)))
}

func encodeCurrentCastleObjectSlot99FromSource(source string) map[string]any {
	units := jsUTF16Units(source)
	return map[string]any{
		"n":   encodeCurrentCastleUnits(currentCastleSliceUnits(units, 0, 9)),
		"cOw": encodeCurrentCastleUnits(currentCastleSliceUnits(units, 9, len(units))),
	}
}

func encodeCurrentCastleObjectSlot100(payload []any) (map[string]any, error) {
	if 99 >= len(payload) {
		return nil, fmt.Errorf("source slot 99 outside payload length %d", len(payload))
	}
	slot99, ok := payload[99].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("source slot 99 has type %T, want map[string]any", payload[99])
	}
	return encodeCurrentCastleObjectSlot100FromSlot99(slot99)
}

func encodeCurrentCastleObjectSlot100FromSlot99(slot99 map[string]any) (map[string]any, error) {
	slotJSON, err := currentCastleJSONStringifySlot99(slot99)
	if err != nil {
		return nil, err
	}
	units := jsUTF16Units(currentCastleTQHash(slotJSON))
	return map[string]any{
		"rj": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 0, 1)),
		"KwJ": map[string]any{
			"T": encodeCurrentCastleUnits(currentCastleSliceUnits(units, 1, 10)),
		},
		"L": map[string]any{
			"T": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 10, 14)),
		},
		"E": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 14, len(units))),
	}, nil
}

func encodeCurrentCastleObjectSlot128(payload []any) (map[string]any, error) {
	return encodeCurrentCastleTIToObjectSlot(payload, 127, map[string]currentCastleObjectFieldSpan{
		"ekC": {start: 0, end: 7},
		"fCH": {start: 7, end: 8},
	})
}

func encodeCurrentCastleObjectSlot214(payload []any) (map[string]any, error) {
	return encodeCurrentCastleTIObjectSlot214(payload)
}

func encodeCurrentCastleObjectSlot173(payload []any) (map[string]any, error) {
	units, err := currentCastleTQHashUnitsForSlot(payload, 172)
	if err != nil {
		return nil, err
	}
	nested := currentCastleSliceUnits(units, 12, 17)
	return map[string]any{
		"t":   encodeCurrentCastleUnits(currentCastleSliceUnits(units, 0, 3)),
		"JZX": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 3, 12)),
		"ga": map[string]any{
			"sE":  encodeCurrentCastleUnits(currentCastleSliceUnits(nested, 0, 6)),
			"OkL": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(nested, 6, 16)),
			"hD":  encodeCurrentCastleUnits(currentCastleSliceUnits(nested, 16, len(nested))),
		},
		"wPg": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 17, 25)),
		"FqX": encodeCurrentCastleUnits(currentCastleSliceUnits(units, 25, len(units))),
	}, nil
}

func encodeCurrentCastleObjectSlot416(payload []any) (map[string]any, error) {
	units, err := currentCastleTIHashUnitsForSlot(payload, 415)
	if err != nil {
		return nil, err
	}
	tail := currentCastleSliceUnits(units, 14, len(units))
	return map[string]any{
		"y":   encodeCurrentCastleUnits(currentCastleSliceUnits(units, 0, 9)),
		"OiJ": encodeCurrentCastleUnits(currentCastleSliceUnits(units, 9, 14)),
		"HRd": map[string]any{
			"j":   encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(tail, 0, 6)),
			"hCP": encodeCurrentCastleUnits(currentCastleSliceUnits(tail, 6, len(tail))),
		},
	}, nil
}

func encodeCurrentCastleObjectSlot424(payload []any) (map[string]any, error) {
	units, err := currentCastleTIHashUnitsForSlot(payload, 423)
	if err != nil {
		return nil, err
	}
	value := currentCastleSliceUnits(units, 4, 11)
	tail := currentCastleSliceUnits(units, 13, len(units))
	return map[string]any{
		"kr": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(units, 0, 4)),
		"vj": map[string]any{
			"GXV": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(value, 0, len(value))),
		},
		"RK": encodeCurrentCastleUnits(currentCastleSliceUnits(units, 11, 13)),
		"wXd": map[string]any{
			"RtN": encodeCurrentCastleRawPackedUnits(currentCastleSliceUnits(tail, 0, len(tail))),
		},
	}, nil
}

type currentCastleObjectFieldSpan struct {
	start int
	end   int
}

func currentCastleTIHashUnitsForSlot(payload []any, sourceSlot int) ([]uint16, error) {
	return currentCastleHashUnitsForSlot(payload, sourceSlot, currentCastleTIHash)
}

func currentCastleTQHashUnitsForSlot(payload []any, sourceSlot int) ([]uint16, error) {
	return currentCastleHashUnitsForSlot(payload, sourceSlot, currentCastleTQHash)
}

func currentCastleHashUnitsForSlot(payload []any, sourceSlot int, hash func(string) string) ([]uint16, error) {
	if sourceSlot < 0 || sourceSlot >= len(payload) {
		return nil, fmt.Errorf("source slot %d outside payload length %d", sourceSlot, len(payload))
	}
	slotJSON, err := currentCastleJSONStringify(payload[sourceSlot])
	if err != nil {
		return nil, err
	}
	return jsUTF16Units(hash(slotJSON)), nil
}

func currentCastleSliceUnits(units []uint16, start, end int) []uint16 {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}
	if start > len(units) {
		start = len(units)
	}
	if end > len(units) {
		end = len(units)
	}
	return units[start:end]
}

func encodeCurrentCastleTIToObjectSlot(payload []any, sourceSlot int, fields map[string]currentCastleObjectFieldSpan) (map[string]any, error) {
	units, err := currentCastleTIHashUnitsForSlot(payload, sourceSlot)
	if err != nil {
		return nil, err
	}
	out := make(map[string]any, len(fields))
	for field, span := range fields {
		if span.start < 0 || span.end < span.start || span.end > len(units) {
			return nil, fmt.Errorf("field %s span %d:%d outside hash length %d", field, span.start, span.end, len(units))
		}
		out[field] = encodeCurrentCastleRawPackedUnits(units[span.start:span.end])
	}
	return out, nil
}

func encodeCurrentCastleTIObjectSlot214(payload []any) (map[string]any, error) {
	if 213 >= len(payload) {
		return nil, fmt.Errorf("source slot 213 outside payload length %d", len(payload))
	}
	slotJSON, err := currentCastleJSONStringify(payload[213])
	if err != nil {
		return nil, err
	}
	hash := currentCastleTIHash(slotJSON)
	units := jsUTF16Units(hash)
	return map[string]any{
		"h":   encodeCurrentCastleUnits(units[:min(len(units), 2)]),
		"yhj": encodeCurrentCastleRawPackedUnits(units[min(len(units), 2):min(len(units), 10)]),
		"XI":  encodeCurrentCastleRawPackedUnits(units[min(len(units), 10):]),
	}, nil
}

func encodeCurrentCastleTQTTObjectSlot(payload []any, sourceSlot int, field string) (map[string]any, error) {
	if sourceSlot < 0 || sourceSlot >= len(payload) {
		return nil, fmt.Errorf("source slot %d outside payload length %d", sourceSlot, len(payload))
	}
	slotJSON, err := currentCastleJSONStringify(payload[sourceSlot])
	if err != nil {
		return nil, err
	}
	return map[string]any{
		field: encodeCurrentCastleUnits(jsUTF16Units(currentCastleTQHash(slotJSON))),
	}, nil
}

func encodeCurrentCastleTITTObjectSlot(payload []any, sourceSlot int, field string) (map[string]any, error) {
	if sourceSlot < 0 || sourceSlot >= len(payload) {
		return nil, fmt.Errorf("source slot %d outside payload length %d", sourceSlot, len(payload))
	}
	slotJSON, err := currentCastleJSONStringify(payload[sourceSlot])
	if err != nil {
		return nil, err
	}
	return map[string]any{
		field: encodeCurrentCastleUnits(jsUTF16Units(currentCastleTIHash(slotJSON))),
	}, nil
}

func currentCastleJSONStringify(value any) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func currentCastleJSONStringifySlot99(slot99 map[string]any) (string, error) {
	n, ok := slot99["n"].(string)
	if !ok {
		return "", fmt.Errorf("slot 99 field n has type %T, want string", slot99["n"])
	}
	cOw, ok := slot99["cOw"].(string)
	if !ok {
		return "", fmt.Errorf("slot 99 field cOw has type %T, want string", slot99["cOw"])
	}
	nJSON, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	cOwJSON, err := json.Marshal(cOw)
	if err != nil {
		return "", err
	}
	return `{"n":` + string(nJSON) + `,"cOw":` + string(cOwJSON) + `}`, nil
}

var currentCastleLowerFloatSlotEncoders = map[int]func(float64) string{
	1:   encodeCurrentCastleLowerFloatSlot1,
	3:   encodeCurrentCastleLowerFloatSlot3,
	4:   encodeCurrentCastleLowerFloatSlot4,
	5:   encodeCurrentCastleLowerFloatSlot5,
	9:   encodeCurrentCastleLowerFloatSlot9,
	10:  encodeCurrentCastleLowerFloatSlot10,
	12:  encodeCurrentCastleLowerFloatSlot12,
	13:  encodeCurrentCastleLowerFloatSlot13,
	14:  encodeCurrentCastleLowerFloatSlot14,
	15:  encodeCurrentCastleLowerFloatSlot15,
	16:  encodeCurrentCastleLowerFloatSlot16,
	17:  encodeCurrentCastleLowerFloatSlot17,
	18:  encodeCurrentCastleLowerFloatSlot18,
	19:  encodeCurrentCastleLowerFloatSlot19,
	21:  encodeCurrentCastleLowerFloatSlot21,
	22:  encodeCurrentCastleLowerFloatSlot22,
	23:  encodeCurrentCastleLowerFloatSlot23,
	24:  encodeCurrentCastleLowerFloatSlot24,
	29:  encodeCurrentCastleLowerFloatSlot29,
	30:  encodeCurrentCastleLowerFloatSlot30,
	32:  encodeCurrentCastleLowerFloatSlot32,
	33:  encodeCurrentCastleLowerFloatSlot33,
	34:  encodeCurrentCastleLowerFloatSlot34,
	35:  encodeCurrentCastleLowerFloatSlot35,
	36:  encodeCurrentCastleLowerFloatSlot36,
	37:  encodeCurrentCastleLowerFloatSlot37,
	38:  encodeCurrentCastleLowerFloatSlot38,
	39:  encodeCurrentCastleLowerFloatSlot39,
	43:  encodeCurrentCastleLowerFloatSlot43,
	45:  encodeCurrentCastleLowerFloatSlot45,
	48:  encodeCurrentCastleLowerFloatSlot48,
	50:  encodeCurrentCastleLowerFloatSlot50,
	51:  encodeCurrentCastleLowerFloatSlot51,
	52:  encodeCurrentCastleLowerFloatSlot52,
	54:  encodeCurrentCastleLowerFloatSlot54,
	55:  encodeCurrentCastleLowerFloatSlot55,
	57:  encodeCurrentCastleLowerFloatSlot57,
	58:  encodeCurrentCastleLowerFloatSlot58,
	59:  encodeCurrentCastleLowerFloatSlot59,
	61:  encodeCurrentCastleLowerFloatSlot61,
	62:  encodeCurrentCastleLowerFloatSlot62,
	63:  encodeCurrentCastleLowerFloatSlot63,
	65:  encodeCurrentCastleLowerFloatSlot65,
	66:  encodeCurrentCastleLowerFloatSlot66,
	69:  encodeCurrentCastleLowerFloatSlot69,
	70:  encodeCurrentCastleLowerFloatSlot70,
	74:  encodeCurrentCastleLowerFloatSlot74,
	76:  encodeCurrentCastleLowerFloatSlot76,
	80:  encodeCurrentCastleLowerFloatSlot80,
	81:  encodeCurrentCastleLowerFloatSlot81,
	82:  encodeCurrentCastleLowerFloatSlot82,
	84:  encodeCurrentCastleLowerFloatSlot84,
	85:  encodeCurrentCastleLowerFloatSlot85,
	86:  encodeCurrentCastleLowerFloatSlot86,
	88:  encodeCurrentCastleLowerFloatSlot88,
	91:  encodeCurrentCastleLowerFloatSlot91,
	92:  encodeCurrentCastleLowerFloatSlot92,
	93:  encodeCurrentCastleLowerFloatSlot93,
	94:  encodeCurrentCastleLowerFloatSlot94,
	95:  encodeCurrentCastleLowerFloatSlot95,
	96:  encodeCurrentCastleLowerFloatSlot96,
	97:  encodeCurrentCastleLowerFloatSlot97,
	101: encodeCurrentCastleLowerFloatSlot101,
	105: encodeCurrentCastleLowerFloatSlot105,
	107: encodeCurrentCastleLowerFloatSlot107,
	110: encodeCurrentCastleLowerFloatSlot110,
	111: encodeCurrentCastleLowerFloatSlot111,
	112: encodeCurrentCastleLowerFloatSlot112,
	114: encodeCurrentCastleLowerFloatSlot114,
	117: encodeCurrentCastleLowerFloatSlot117,
	118: encodeCurrentCastleLowerFloatSlot118,
	122: encodeCurrentCastleLowerFloatSlot122,
	124: encodeCurrentCastleLowerFloatSlot124,
	125: encodeCurrentCastleLowerFloatSlot125,
	126: encodeCurrentCastleLowerFloatSlot126,
	129: encodeCurrentCastleLowerFloatSlot129,
	131: encodeCurrentCastleLowerFloatSlot131,
	132: encodeCurrentCastleLowerFloatSlot132,
	134: encodeCurrentCastleLowerFloatSlot134,
	135: encodeCurrentCastleLowerFloatSlot135,
	136: encodeCurrentCastleLowerFloatSlot136,
	138: encodeCurrentCastleLowerFloatSlot138,
	139: encodeCurrentCastleLowerFloatSlot139,
	140: encodeCurrentCastleLowerFloatSlot140,
	141: encodeCurrentCastleLowerFloatSlot141,
	142: encodeCurrentCastleLowerFloatSlot142,
	143: encodeCurrentCastleLowerFloatSlot143,
	144: encodeCurrentCastleLowerFloatSlot144,
	145: encodeCurrentCastleLowerFloatSlot145,
	146: encodeCurrentCastleLowerFloatSlot146,
	148: encodeCurrentCastleLowerFloatSlot148,
	151: encodeCurrentCastleLowerFloatSlot151,
	153: encodeCurrentCastleLowerFloatSlot153,
	155: encodeCurrentCastleLowerFloatSlot155,
	157: encodeCurrentCastleLowerFloatSlot157,
	158: encodeCurrentCastleLowerFloatSlot158,
	160: encodeCurrentCastleLowerFloatSlot160,
	161: encodeCurrentCastleLowerFloatSlot161,
	162: encodeCurrentCastleLowerFloatSlot162,
	163: encodeCurrentCastleLowerFloatSlot163,
	164: encodeCurrentCastleLowerFloatSlot164,
	165: encodeCurrentCastleLowerFloatSlot165,
	166: encodeCurrentCastleLowerFloatSlot166,
	167: encodeCurrentCastleLowerFloatSlot167,
	169: encodeCurrentCastleLowerFloatSlot169,
	170: encodeCurrentCastleLowerFloatSlot170,
	171: encodeCurrentCastleLowerFloatSlot171,
	178: encodeCurrentCastleLowerFloatSlot178,
	179: encodeCurrentCastleLowerFloatSlot179,
	181: encodeCurrentCastleLowerFloatSlot181,
	185: encodeCurrentCastleLowerFloatSlot185,
	186: encodeCurrentCastleLowerFloatSlot186,
	187: encodeCurrentCastleLowerFloatSlot187,
	189: encodeCurrentCastleLowerFloatSlot189,
	190: encodeCurrentCastleLowerFloatSlot190,
	191: encodeCurrentCastleLowerFloatSlot191,
	192: encodeCurrentCastleLowerFloatSlot192,
	193: encodeCurrentCastleLowerFloatSlot193,
	195: encodeCurrentCastleLowerFloatSlot195,
	197: encodeCurrentCastleLowerFloatSlot197,
	199: encodeCurrentCastleLowerFloatSlot199,
	200: encodeCurrentCastleLowerFloatSlot200,
	201: encodeCurrentCastleLowerFloatSlot201,
	202: encodeCurrentCastleLowerFloatSlot202,
	208: encodeCurrentCastleLowerFloatSlot208,
	210: encodeCurrentCastleLowerFloatSlot210,
	211: encodeCurrentCastleLowerFloatSlot211,
	215: encodeCurrentCastleLowerFloatSlot215,
	216: encodeCurrentCastleLowerFloatSlot216,
	218: encodeCurrentCastleLowerFloatSlot218,
	219: encodeCurrentCastleLowerFloatSlot219,
	222: encodeCurrentCastleLowerFloatSlot222,
	225: encodeCurrentCastleLowerFloatSlot225,
	226: encodeCurrentCastleLowerFloatSlot226,
	227: encodeCurrentCastleLowerFloatSlot227,
	229: encodeCurrentCastleLowerFloatSlot229,
	230: encodeCurrentCastleLowerFloatSlot230,
	232: encodeCurrentCastleLowerFloatSlot232,
	233: encodeCurrentCastleLowerFloatSlot233,
	234: encodeCurrentCastleLowerFloatSlot234,
	236: encodeCurrentCastleLowerFloatSlot236,
	238: encodeCurrentCastleLowerFloatSlot238,
	241: encodeCurrentCastleLowerFloatSlot241,
	242: encodeCurrentCastleLowerFloatSlot242,
	243: encodeCurrentCastleLowerFloatSlot243,
	244: encodeCurrentCastleLowerFloatSlot244,
	248: encodeCurrentCastleLowerFloatSlot248,
	249: encodeCurrentCastleLowerFloatSlot249,
	250: encodeCurrentCastleLowerFloatSlot250,
	251: encodeCurrentCastleLowerFloatSlot251,
	252: encodeCurrentCastleLowerFloatSlot252,
	253: encodeCurrentCastleLowerFloatSlot253,
	255: encodeCurrentCastleLowerFloatSlot255,
	259: encodeCurrentCastleLowerFloatSlot259,
	262: encodeCurrentCastleLowerFloatSlot262,
	263: encodeCurrentCastleLowerFloatSlot263,
	264: encodeCurrentCastleLowerFloatSlot264,
	265: encodeCurrentCastleLowerFloatSlot265,
	268: encodeCurrentCastleLowerFloatSlot268,
	269: encodeCurrentCastleLowerFloatSlot269,
	270: encodeCurrentCastleLowerFloatSlot270,
	271: encodeCurrentCastleLowerFloatSlot271,
	272: encodeCurrentCastleLowerFloatSlot272,
	275: encodeCurrentCastleLowerFloatSlot275,
	276: encodeCurrentCastleLowerFloatSlot276,
	278: encodeCurrentCastleLowerFloatSlot278,
	280: encodeCurrentCastleLowerFloatSlot280,
	282: encodeCurrentCastleLowerFloatSlot282,
	283: encodeCurrentCastleLowerFloatSlot283,
	284: encodeCurrentCastleLowerFloatSlot284,
	285: encodeCurrentCastleLowerFloatSlot285,
	286: encodeCurrentCastleLowerFloatSlot286,
	287: encodeCurrentCastleLowerFloatSlot287,
	291: encodeCurrentCastleLowerFloatSlot291,
	294: encodeCurrentCastleLowerFloatSlot294,
	295: encodeCurrentCastleLowerFloatSlot295,
	296: encodeCurrentCastleLowerFloatSlot296,
	299: encodeCurrentCastleLowerFloatSlot299,
	300: encodeCurrentCastleLowerFloatSlot300,
	301: encodeCurrentCastleLowerFloatSlot301,
	303: encodeCurrentCastleLowerFloatSlot303,
	305: encodeCurrentCastleLowerFloatSlot305,
	306: encodeCurrentCastleLowerFloatSlot306,
	307: encodeCurrentCastleLowerFloatSlot307,
	309: encodeCurrentCastleLowerFloatSlot309,
	310: encodeCurrentCastleLowerFloatSlot310,
	313: encodeCurrentCastleLowerFloatSlot313,
	314: encodeCurrentCastleLowerFloatSlot314,
	315: encodeCurrentCastleLowerFloatSlot315,
	317: encodeCurrentCastleLowerFloatSlot317,
	319: encodeCurrentCastleLowerFloatSlot319,
	320: encodeCurrentCastleLowerFloatSlot320,
	322: encodeCurrentCastleLowerFloatSlot322,
	323: encodeCurrentCastleLowerFloatSlot323,
	324: encodeCurrentCastleLowerFloatSlot324,
	325: encodeCurrentCastleLowerFloatSlot325,
	327: encodeCurrentCastleLowerFloatSlot327,
	331: encodeCurrentCastleLowerFloatSlot331,
	332: encodeCurrentCastleLowerFloatSlot332,
	333: encodeCurrentCastleLowerFloatSlot333,
	335: encodeCurrentCastleLowerFloatSlot335,
	336: encodeCurrentCastleLowerFloatSlot336,
	337: encodeCurrentCastleLowerFloatSlot337,
	338: encodeCurrentCastleLowerFloatSlot338,
	340: encodeCurrentCastleLowerFloatSlot340,
	343: encodeCurrentCastleLowerFloatSlot343,
	345: encodeCurrentCastleLowerFloatSlot345,
	346: encodeCurrentCastleLowerFloatSlot346,
	347: encodeCurrentCastleLowerFloatSlot347,
	353: encodeCurrentCastleLowerFloatSlot353,
	355: encodeCurrentCastleLowerFloatSlot355,
	357: encodeCurrentCastleLowerFloatSlot357,
	358: encodeCurrentCastleLowerFloatSlot358,
	360: encodeCurrentCastleLowerFloatSlot360,
	361: encodeCurrentCastleLowerFloatSlot361,
	362: encodeCurrentCastleLowerFloatSlot362,
	363: encodeCurrentCastleLowerFloatSlot363,
	364: encodeCurrentCastleLowerFloatSlot364,
	365: encodeCurrentCastleLowerFloatSlot365,
	367: encodeCurrentCastleLowerFloatSlot367,
	368: encodeCurrentCastleLowerFloatSlot368,
	370: encodeCurrentCastleLowerFloatSlot370,
	371: encodeCurrentCastleLowerFloatSlot371,
	372: encodeCurrentCastleLowerFloatSlot372,
	373: encodeCurrentCastleLowerFloatSlot373,
	377: encodeCurrentCastleLowerFloatSlot377,
	379: encodeCurrentCastleLowerFloatSlot379,
	380: encodeCurrentCastleLowerFloatSlot380,
	381: encodeCurrentCastleLowerFloatSlot381,
	382: encodeCurrentCastleLowerFloatSlot382,
	384: encodeCurrentCastleLowerFloatSlot384,
	385: encodeCurrentCastleLowerFloatSlot385,
	387: encodeCurrentCastleLowerFloatSlot387,
	389: encodeCurrentCastleLowerFloatSlot389,
	391: encodeCurrentCastleLowerFloatSlot391,
	392: encodeCurrentCastleLowerFloatSlot392,
	396: encodeCurrentCastleLowerFloatSlot396,
	397: encodeCurrentCastleLowerFloatSlot397,
	398: encodeCurrentCastleLowerFloatSlot398,
	401: encodeCurrentCastleLowerFloatSlot401,
	402: encodeCurrentCastleLowerFloatSlot402,
	403: encodeCurrentCastleLowerFloatSlot403,
	404: encodeCurrentCastleLowerFloatSlot404,
	405: encodeCurrentCastleLowerFloatSlot405,
	407: encodeCurrentCastleLowerFloatSlot407,
	408: encodeCurrentCastleLowerFloatSlot408,
	409: encodeCurrentCastleLowerFloatSlot409,
	410: encodeCurrentCastleLowerFloatSlot410,
	411: encodeCurrentCastleLowerFloatSlot411,
	412: encodeCurrentCastleLowerFloatSlot412,
	413: encodeCurrentCastleLowerFloatSlot413,
	414: encodeCurrentCastleLowerFloatSlot414,
	417: encodeCurrentCastleLowerFloatSlot417,
	418: encodeCurrentCastleLowerFloatSlot418,
	419: encodeCurrentCastleLowerFloatSlot419,
	420: encodeCurrentCastleLowerFloatSlot420,
	426: encodeCurrentCastleLowerFloatSlot426,
	427: encodeCurrentCastleLowerFloatSlot427,
}

func populateCurrentCastleLowerFloatSlots(payload []any, values map[int]float64) error {
	if len(payload) < currentCastlePayloadSlots {
		return fmt.Errorf("current Castle payload has %d slots, want at least %d", len(payload), currentCastlePayloadSlots)
	}
	for slot, encoder := range currentCastleLowerFloatSlotEncoders {
		payload[slot] = encoder(values[slot])
	}
	return nil
}

var currentCastleArraySlotEncoders = map[int]func([]float64) string{
	254: encodeCurrentCastleArraySlot254,
	302: encodeCurrentCastleArraySlot302,
}

func defaultCurrentCastleArraySlotValues() map[int][]float64 {
	return map[int][]float64{
		// The visible Chrome begin_login token encodes sparse one-hot samples.
		254: currentCastleFloat64Samples(256, 199),
		302: currentCastleFloat64Samples(16, 13),
	}
}

func currentCastleFloat64Samples(count int, oneIndexes ...int) []float64 {
	if count <= 0 {
		return nil
	}
	values := make([]float64, count)
	for _, index := range oneIndexes {
		if index >= 0 && index < len(values) {
			values[index] = 1
		}
	}
	return values
}

func populateCurrentCastleArraySlots(payload []any, values map[int][]float64) error {
	if len(payload) < currentCastlePayloadSlots {
		return fmt.Errorf("current Castle payload has %d slots, want at least %d", len(payload), currentCastlePayloadSlots)
	}
	defaults := defaultCurrentCastleArraySlotValues()
	for slot, encoder := range currentCastleArraySlotEncoders {
		slotValues, ok := values[slot]
		if !ok {
			slotValues = defaults[slot]
		}
		payload[slot] = encoder(slotValues)
	}
	return nil
}

var currentCastlePackedStringSlotEncoders = map[int]func([]string) string{
	0:   encodeCurrentCastlePackedStringComponents,
	7:   encodeCurrentCastlePackedStringComponents,
	119: encodeCurrentCastlePackedStringComponents,
}

func defaultCurrentCastlePackedStringSlotValues() map[int][]string {
	return map[int][]string{
		0:   defaultCurrentCastleSlot0PackedComponents(),
		7:   defaultCurrentCastleSlot7PackedComponents(),
		119: defaultCurrentCastleWebGPUPackedComponents(),
	}
}

func defaultCurrentCastleSlot0PackedComponents() []string {
	return []string{
		// Current Chrome begin_login defaults from the native Windows profile.
		"pLJf4l/i",
		"hoom",
		"vbq6wVWT2EHZILqMkw==",
		"",
		"nLpq61raGg==",
		"P0JHPD48RUZARTw/R0Q=",
	}
}

func defaultCurrentCastleSlot7PackedComponents() []string {
	return []string{
		// Current Chrome begin_login defaults from the native Windows profile.
		"CJk=",
		"JEbZeg==",
		"m3s=",
		"sJ5lbMPD9N7IWWe3L7ZsGYue1nK3+Vu35GdZZyq3tmwZUTsqt19RO6y3VCcnw1i2WH2abPveyLIWWbJRty+az1uwIwO3w2w6WLdGWAI6nqy3Yu8tnlBY3uQ7HFlnWWdZZ7fS9OH0LWzeyLIWWbJR",
		"i+vAAJAdJIbAaMBvZ2dPnsmaY5KQZ+HJwF7AnMmQ09VnJMBowK3WJeRq5J/iUJ/krSVcwHFe68AAkB0khsBowGOSkGfh1VnhwF7AnMmQ09VnJMBowK3WJeRq5J/iUJ/krSVcwHFe68AAkB0khsBowERnFpU9qyCQHSSGwF7AnMmQ09VnJMBowFDW5GrkauRqwHER",
		"ltZunm6e",
		"YYsEQpkB5Q==",
		"2gykDtI=",
		"kHDTk+PKP+HT69NcxMQ8qa59eMXjxHau0+nTx67jqJLEP9Pr05qxItPG6XDTk+PKP+HT69N4xePEdpI+dtPp08eu46iSxD/T69OasSLTxulw05Pjyj/h0+vTX8Rx0uqws+PKP+HT6dPHruOoksQ/0+vTI7HTxuY=",
		"8/u7",
		"Fy9X/g/+J08fJ/4XVz8=",
	}
}

func defaultCurrentCastleWebGPUPackedComponents() []string {
	return []string{
		// Current Chrome begin_login defaults from the native Windows profile.
		"ACiTIpOL",
		"Lga4aAbykN6kkA==",
		"mY47m47x",
	}
}

func populateCurrentCastlePackedStringSlots(payload []any, values map[int][]string) error {
	if len(payload) < currentCastlePayloadSlots {
		return fmt.Errorf("current Castle payload has %d slots, want at least %d", len(payload), currentCastlePayloadSlots)
	}
	defaults := defaultCurrentCastlePackedStringSlotValues()
	for slot, encoder := range currentCastlePackedStringSlotEncoders {
		slotValues, ok := values[slot]
		if !ok {
			slotValues = defaults[slot]
		}
		payload[slot] = encoder(slotValues)
	}
	return nil
}

var currentCastleUnitPackedStringSlotEncoders = map[int]func([]string) string{
	8: encodeCurrentCastleUnitPackedStringComponents,
}

func defaultCurrentCastleUnitPackedStringSlotValues() map[int][]string {
	return map[int][]string{
		8: defaultCurrentCastleSlot8PackedComponents(),
	}
}

func defaultCurrentCastleSlot8PackedComponents() []string {
	return []string{
		// Current Chrome begin_login defaults from the native Windows profile.
		"a14=",
		"pD/f1v/ISA==",
	}
}

func populateCurrentCastleUnitPackedStringSlots(payload []any, values map[int][]string) error {
	if len(payload) < currentCastlePayloadSlots {
		return fmt.Errorf("current Castle payload has %d slots, want at least %d", len(payload), currentCastlePayloadSlots)
	}
	defaults := defaultCurrentCastleUnitPackedStringSlotValues()
	for slot, encoder := range currentCastleUnitPackedStringSlotEncoders {
		slotValues, ok := values[slot]
		if !ok {
			slotValues = defaults[slot]
		}
		payload[slot] = encoder(slotValues)
	}
	return nil
}

var currentCastleStringSlotEncoders = map[int]func(string) string{
	11:  encodeCurrentCastleRawString,
	40:  encodeCurrentCastleDoubleUTF16Scramble,
	89:  encodeCurrentCastleUTF16Scramble,
	127: encodeCurrentCastleUTF16Scramble,
	172: encodeCurrentCastleUTF16Scramble,
	175: encodeCurrentCastleUTF16Scramble,
	213: encodeCurrentCastleRawString,
	217: encodeCurrentCastleDoubleUTF16Scramble,
	231: encodeCurrentCastleUTF16Scramble,
	351: encodeCurrentCastleUTF16Scramble,
	400: encodeCurrentCastleDoubleUTF16Scramble,
	415: encodeCurrentCastleDoubleUTF16Scramble,
	423: encodeCurrentCastleDoubleUTF16Scramble,
}

func defaultCurrentCastleStringSlotValues() map[int]string {
	return map[int]string{
		11:  "isz",
		40:  "0",
		89:  "0.5",
		127: "0",
		172: "0.95",
		175: "0.207",
		213: "GraL",
		217: "0.217",
		231: "0.486",
		351: "0.279",
		400: "0.4",
		415: "0.415",
		423: "0.423",
	}
}

func populateCurrentCastleStringSlots(payload []any, values map[int]string) error {
	if len(payload) < currentCastlePayloadSlots {
		return fmt.Errorf("current Castle payload has %d slots, want at least %d", len(payload), currentCastlePayloadSlots)
	}
	defaults := defaultCurrentCastleStringSlotValues()
	finalDefaults := defaultCurrentCastleStringSlotEncodedValues()
	for slot, encoder := range currentCastleStringSlotEncoders {
		if _, hasOverride := values[slot]; !hasOverride {
			if encoded, ok := finalDefaults[slot]; ok {
				payload[slot] = encoded
				continue
			}
		}
		slotValue, ok := values[slot]
		if !ok {
			slotValue = defaults[slot]
		}
		payload[slot] = encoder(slotValue)
	}
	return nil
}

func defaultCurrentCastleStringSlotEncodedValues() map[int]string {
	const encodedZero = "NoAtnP9LeAQ="
	return map[int]string{
		// Current Chrome begin_login defaults from the native Windows profile.
		11:  "isz",
		40:  encodedZero,
		89:  encodedZero,
		127: encodedZero,
		172: encodedZero,
		175: encodedZero,
		213: "GraL",
		217: encodedZero,
		231: encodedZero,
		351: encodedZero,
		400: encodedZero,
		415: encodedZero,
		423: encodedZero,
	}
}

func currentCastleLengthPrefixUnits(length uint32) []uint16 {
	var out []uint16
	for length >= currentCastleSlot6LengthContinue {
		out = append(out, uint16(length&currentCastleSlot6LengthMask|currentCastleSlot6LengthContinue))
		length >>= currentCastleSlot6LengthShiftBits
	}
	return append(out, uint16(length&currentCastleSlot6LengthMask))
}

func currentCastleSlot6LengthPrefixUnits(componentLengths []int) ([]uint16, error) {
	if len(componentLengths) != currentCastleSlot6ComponentCount {
		return nil, fmt.Errorf("current Castle slot 6 has %d component lengths, want %d", len(componentLengths), currentCastleSlot6ComponentCount)
	}
	out := currentCastleLengthPrefixUnits(currentCastleSlot6ComponentCount)
	for i, length := range componentLengths {
		if length < 0 {
			return nil, fmt.Errorf("current Castle slot 6 component %d has negative length %d", i, length)
		}
		out = append(out, currentCastleLengthPrefixUnits(uint32(length))...)
	}
	return out, nil
}

func marshalCurrentCastlePayload(payload []any) (string, error) {
	if len(payload) != currentCastlePayloadSlots {
		return "", fmt.Errorf("current Castle payload has %d slots, want %d", len(payload), currentCastlePayloadSlots)
	}
	for i, value := range payload {
		if value == nil {
			return "", fmt.Errorf("current Castle payload slot %d is nil", i)
		}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func insertCurrentCastleChecksum(payloadJSON string) string {
	units := jsUTF16Units(payloadJSON)
	checksum := currentCastleChecksum(payloadJSON)
	insertAt := int(math.Floor(0.76 * float64(len(units))))
	label := currentCastleChecksumLabel + checksum
	withChecksum := make([]uint16, 0, len(units)+len(label))
	withChecksum = append(withChecksum, units[:insertAt]...)
	withChecksum = append(withChecksum, jsUTF16Units(label)...)
	withChecksum = append(withChecksum, units[insertAt:]...)
	return string(utf16.Decode(withChecksum))
}

func currentCastleChecksum(value string) string {
	state := uint32(len(jsUTF16Units(value)))
	sum := uint32(0)
	for index, ch := range jsUTF16Units(value) {
		state = bitswap32(state)
		state = bitswap32(state)
		state = state + currentCastleRotate16Into32(ch, 7, 25)
		state = bitswap32(state)
		sum += uint32(ch) * 0x9e3779b1
		if index&3 == 3 {
			state = state + currentCastleRotate16Into32(uint16(state), 4, 28)
		}
	}
	state ^= sum
	state = (state << 21) | (state >> 11)
	state ^= state >> 16
	state *= 0x85ebca6b
	state ^= state >> 13
	state *= 0xc2b2ae35
	state ^= state >> 16
	var out [4]byte
	out[0] = byte(state >> 24)
	out[1] = byte(state >> 16)
	out[2] = byte(state >> 8)
	out[3] = byte(state)
	return hex.EncodeToString(out[:])
}

func currentCastleTQHash(value string) string {
	units := jsUTF16Units(value)
	state := uint32(0xc1960396) ^ uint32(len(units))
	sum := uint32(0)
	state ^= 0 * 0x38da94b1
	state += ((state & 0xffff) << 5) | ((state & 0xffff) >> 27)
	state ^= ((state & 0xffff) << 24) | ((state & 0xffff) >> 8)
	state = (state^(state>>11))*0x88a4e807 + 0x89f085ba
	for index, ch := range units {
		raw := uint32(ch)
		state ^= raw
		state = (state ^ raw) * 0x69d69459
		state ^= 0xa87dabc4
		state ^= raw
		state ^= ((raw & 0xffff) << 14) | ((raw & 0xffff) >> 18)
		state ^= 0x83274064
		state += ((raw & 0xffff) << 14) | ((raw & 0xffff) >> 18)
		state ^= (sum << 13) | (sum >> 19)
		state += raw * 0x858ff2cb
		state += ((raw & 0xffff) << 12) | ((raw & 0xffff) >> 20)
		state = (state ^ raw) * 0x9dfcf351
		state = bits.RotateLeft32(state, 17)
		state = bitswap32(state)
		sum += raw * 0x9e3779b1
		if index&3 == 3 {
			state += (state ^ (state >> 11)) * 0xd1a79fef
			state += ((state & 0xffff) << 8) | ((state & 0xffff) >> 24)
			state = (state ^ (state >> 11)) * 0xa3cfce21
			state ^= state >> 5
		}
	}
	state ^= sum
	state ^= sum
	state = bitswap32(state)
	state += ((sum & 0xffff) << 15) | ((sum & 0xffff) >> 17)
	state = bits.RotateLeft32(state, 15)
	state ^= state >> 16
	state *= 0x85ebca6b
	state ^= state >> 13
	state *= 0xc2b2ae35
	state ^= state >> 16
	return fmt.Sprintf("%08x", state)
}

func currentCastleTIHash(value string) string {
	units := jsUTF16Units(value)
	state := uint32(0xdd3b1081) ^ uint32(len(units))
	sum := uint32(0)
	state ^= 0 * 0x1ab5d0cb
	state = (state ^ (state >> 11)) * 0xfcc4435f
	state += (state ^ (state >> 11)) * 0x7d8bc77f
	for index, ch := range units {
		raw := uint32(ch)
		state ^= state >> 8
		state ^= state >> 16
		state += ((raw & 0xffff) << 1) | ((raw & 0xffff) >> 31)
		state ^= ((raw & 0xffff) << 5) | ((raw & 0xffff) >> 27)
		state ^= raw
		state += ((raw & 0xffff) << 5) | ((raw & 0xffff) >> 27)
		state += ((raw & 0xffff) << 14) | ((raw & 0xffff) >> 18)
		state = bits.RotateLeft32(state, 20)
		state = (state ^ raw) * 0x98f0350b
		state += raw * 0xd9853933
		sum += raw * 0x9e3779b1
		if index&3 == 3 {
			state ^= uint32(len(units)) * 0x6963131f
			state = (state^(state>>11))*0xe0acd687 + 0x75c53c2d
			state = bits.RotateLeft32(state, 12)
		}
	}
	state ^= sum
	state = (state ^ sum) * 0xb7fb31df
	state ^= ((sum & 0xffff) << 12) | ((sum & 0xffff) >> 20)
	state ^= (sum << 19) | (sum >> 13)
	state ^= state >> 16
	state *= 0x85ebca6b
	state ^= state >> 13
	state *= 0xc2b2ae35
	state ^= state >> 16
	return fmt.Sprintf("%08x", state)
}

func encodeCurrentCastleTimestamp(timestamp string) string {
	return encodeCurrentCastleUnits(jsUTF16Units(timestamp))
}

func encodeCurrentCastleUnits(units []uint16) string {
	out := make([]byte, 0, len(units)*2)
	for _, ch := range units {
		n := (uint32(ch) - 54655) & 0xffff
		n = (46488 ^ n) & 0xffff
		n = (uint32(uint16(n))*54213 + 385) & 0xffff
		n = ((n >> 12) | (n << 4)) & 0xffff
		n = (60834 + n) & 0xffff
		n = ((n >> 11) | (n << 5)) & 0xffff
		out = append(out, byte(n>>8), byte(n))
	}
	return base64.StdEncoding.EncodeToString(out)
}

func encodeCurrentCastleRawPackedUnits(units []uint16) string {
	out := make([]byte, 0, len(units)*2)
	for _, ch := range units {
		n := (uint32(ch)*25467 + 59068) & 0xffff
		n = (n ^ 17750) & 0xffff
		n = currentCastleRotL16(n, 14)
		n = currentCastleRotL16(n, 14)
		out = append(out, byte(n>>8), byte(n))
	}
	return base64.StdEncoding.EncodeToString(out)
}

func encodeCurrentCastleUTF16Scramble(value string) string {
	out := make([]byte, 0, len(jsUTF16Units(value))*2)
	for _, ch := range jsUTF16Units(value) {
		n := (uint32(ch) ^ 27722) & 0xffff
		n = (n - 40955) & 0xffff
		n = ((n << 15) | (n >> 1)) & 0xffff
		n = (uint32(uint16(n))*49687 + 51498) & 0xffff
		n = (n ^ 36024) & 0xffff
		out = append(out, byte(n>>8), byte(n))
	}
	return base64.StdEncoding.EncodeToString(out)
}

func encodeCurrentCastleDoubleUTF16Scramble(value string) string {
	return encodeCurrentCastleUTF16Scramble(encodeCurrentCastleUTF16Scramble(value))
}

func encodeCurrentCastleRawString(value string) string {
	return value
}

func encodeCurrentCastleLR(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 8
		b = currentCastleRotL8(b, 6)
		b = byte(uint16(b) * 207)
		b += 218
		return b
	})
}

func currentCastleLengthPrefixedStringComponentUnits(values []string) []uint16 {
	componentUnits := make([][]uint16, len(values))
	componentLengths := make([]int, len(values))
	totalComponentUnits := 0
	for i, value := range values {
		units := jsUTF16Units(value)
		componentUnits[i] = units
		componentLengths[i] = len(units)
		totalComponentUnits += len(units)
	}
	rawUnits := make([]uint16, 0, len(values)+1+totalComponentUnits)
	rawUnits = append(rawUnits, currentCastleLengthPrefixUnits(uint32(len(values)))...)
	for _, length := range componentLengths {
		rawUnits = append(rawUnits, currentCastleLengthPrefixUnits(uint32(length))...)
	}
	for _, units := range componentUnits {
		rawUnits = append(rawUnits, units...)
	}
	return rawUnits
}

func encodeCurrentCastlePackedStringComponents(values []string) string {
	rawUnits := currentCastleLengthPrefixedStringComponentUnits(values)
	return encodeCurrentCastleRawPackedUnits(rawUnits)
}

func encodeCurrentCastleUnitPackedStringComponents(values []string) string {
	return encodeCurrentCastleUnits(currentCastleLengthPrefixedStringComponentUnits(values))
}

func encodeCurrentCastleSlot0Primary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 1
		b = byte(uint16(b) * 7)
		b += 52
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleSlot0Secondary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 163)
		b = currentCastleRotL8(b, 2)
		b = currentCastleRotL8(b, 3)
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleSlot0Tertiary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 185)
		b = currentCastleRotL8(b, 5)
		b += 212
		return b
	})
}

func encodeCurrentCastleSlot0Quaternary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 223)
		b = currentCastleRotL8(b, 5)
		b ^= 159
		return b
	})
}

func encodeCurrentCastleSlot0Quinary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b += 220
		b = currentCastleRotL8(b, 2)
		b = byte(uint16(b) * 31)
		return b
	})
}

func encodeCurrentCastleSlot0Senary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 249
		b += 85
		b += 192
		return b
	})
}

func encodeCurrentCastleSlot7CanvasPrimary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 133)
		b += 242
		b = byte(uint16(b) * 215)
		b = byte(uint16(b) * 105)
		return b
	})
}

func encodeCurrentCastleSlot7WorkerTiming(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 83
		b ^= 38
		b = byte(uint16(b) * 105)
		b = byte(uint16(b) * 211)
		return b
	})
}

func encodeCurrentCastleSlot7MediaRatio(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b ^= 207
		b += 239
		return b
	})
}

func encodeCurrentCastleSlot7PluginState(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 135)
		b ^= 37
		b += 242
		return b
	})
}

func encodeCurrentCastleSlot7FrameRatio(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 238
		b ^= 80
		b = byte(uint16(b) * 67)
		return b
	})
}

func encodeCurrentCastleSlot7WorkerCanvas(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = currentCastleRotL8(b, 5)
		b ^= 31
		return b
	})
}

func encodeCurrentCastleSlot7WorkerFeature(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 105
		b ^= 187
		b = byte(uint16(b) * 211)
		return b
	})
}

func encodeCurrentCastleSlot7NavigatorProbe(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 119)
		b = byte(uint16(b) * 227)
		b = currentCastleRotL8(b, 5)
		b += 96
		return b
	})
}

func encodeCurrentCastleSlot7ViewportRatio(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 55)
		b = byte(uint16(b) * 129)
		b ^= 157
		return b
	})
}

func encodeCurrentCastleSlot7WorkerSignal(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 70
		b = currentCastleRotL8(b, 5)
		b += 151
		b ^= 157
		return b
	})
}

func encodeCurrentCastleSlot7PointerRatio(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b += 202
		b += 162
		b = currentCastleRotL8(b, 5)
		return b
	})
}

func encodeCurrentCastleSlot8Primary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = byte(uint16(b) * 45)
		b ^= 87
		return b
	})
}

func encodeCurrentCastleSlot8Secondary(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b += 154
		b = byte(uint16(b) * 201)
		return b
	})
}

func encodeCurrentCastleWebGPUVendor(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 79)
		b = byte(uint16(b) * 211)
		b ^= 118
		return b
	})
}

func encodeCurrentCastleWebGPULimits(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = currentCastleRotL8(b, 5)
		b += 72
		b = byte(uint16(b) * 187)
		return b
	})
}

func encodeCurrentCastleWebGPUArchitecture(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 59)
		b = currentCastleRotL8(b, 1)
		b ^= 56
		return b
	})
}

func currentCastleSlot6RawUnits(components []string) ([]uint16, error) {
	if len(components) != currentCastleSlot6ComponentCount {
		return nil, fmt.Errorf("current Castle slot 6 has %d components, want %d", len(components), currentCastleSlot6ComponentCount)
	}
	componentUnits := make([][]uint16, len(components))
	componentLengths := make([]int, len(components))
	totalComponentUnits := 0
	for i, component := range components {
		units := jsUTF16Units(component)
		componentUnits[i] = units
		componentLengths[i] = len(units)
		totalComponentUnits += len(units)
	}
	prefixUnits, err := currentCastleSlot6LengthPrefixUnits(componentLengths)
	if err != nil {
		return nil, err
	}
	out := make([]uint16, 0, len(prefixUnits)+totalComponentUnits)
	out = append(out, prefixUnits...)
	for _, units := range componentUnits {
		out = append(out, units...)
	}
	return out, nil
}

func encodeCurrentCastleSlot6ComponentValues(values []string) ([]string, error) {
	if len(values) != currentCastleSlot6ComponentCount {
		return nil, fmt.Errorf("current Castle slot 6 has %d raw component values, want %d", len(values), currentCastleSlot6ComponentCount)
	}
	encoders := [...]func(string) string{
		encodeCurrentCastleSlot6Component0,
		encodeCurrentCastleSlot6Component1,
		encodeCurrentCastleSlot6Component2,
		encodeCurrentCastleSlot6UZ,
		encodeCurrentCastleSlot6TF,
		encodeCurrentCastleSlot6Component5Fallback,
		encodeCurrentCastleSlot6TU,
		encodeCurrentCastleSlot6IH,
		encodeCurrentCastleSlot6FJ,
		encodeCurrentCastleSlot6IB,
		encodeCurrentCastleSlot6AH,
		encodeCurrentCastleSlot6IT,
		encodeCurrentCastleSlot6AG,
		encodeCurrentCastleSlot6AD,
		encodeCurrentCastleSlot6T7,
		encodeCurrentCastleSlot6ID,
		encodeCurrentCastleSlot6FU,
		encodeCurrentCastleSlot6Component17,
		encodeCurrentCastleSlot6EJ,
		encodeCurrentCastleSlot6Component19,
		encodeCurrentCastleSlot6L9,
		encodeCurrentCastleSlot6TX,
		encodeCurrentCastleSlot6LI,
		encodeCurrentCastleSlot6FL,
		encodeCurrentCastleSlot6IO,
		encodeCurrentCastleSlot6F3,
		encodeCurrentCastleSlot6AJ,
		encodeCurrentCastleSlot6Component27,
		encodeCurrentCastleSlot6TM,
		encodeCurrentCastleSlot6TUpperF,
		encodeCurrentCastleSlot6FS,
		encodeCurrentCastleSlot6Component31,
		encodeCurrentCastleSlot6FF,
		encodeCurrentCastleSlot6AU,
		encodeCurrentCastleSlot6AQ,
		encodeCurrentCastleSlot6LW,
		encodeCurrentCastleSlot6TUpperM,
		encodeCurrentCastleSlot6TH,
		encodeCurrentCastleSlot6IS,
		encodeCurrentCastleSlot6A8,
		encodeCurrentCastleSlot6LUpperM,
		encodeCurrentCastleSlot6EA,
		encodeCurrentCastleSlot6E3,
		encodeCurrentCastleSlot6T2,
		encodeCurrentCastleSlot6TUpperX,
		encodeCurrentCastleSlot6L5,
		encodeCurrentCastleSlot6Component46,
	}
	components := make([]string, currentCastleSlot6ComponentCount)
	for i, value := range values {
		components[i] = encoders[i](value)
	}
	return components, nil
}

func encodeCurrentCastleSlot6(components []string) (string, error) {
	units, err := currentCastleSlot6RawUnits(components)
	if err != nil {
		return "", err
	}
	return encodeCurrentCastleUnits(units), nil
}

func encodeCurrentCastleTransformedUTF8(value string, transform func(byte) byte) string {
	input := []byte(value)
	output := make([]byte, len(input))
	for i, b := range input {
		output[i] = transform(b)
	}
	return base64.StdEncoding.EncodeToString(output)
}

func currentCastleRotL8(value byte, shift uint) byte {
	return value<<shift | value>>(8-shift)
}

func currentCastleRotL16(value uint32, shift uint) uint32 {
	return ((value << shift) | (value >> (16 - shift))) & 0xffff
}

func encodeCurrentCastleSlot6TU(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 113)
		b ^= 246
		b += 41
		return currentCastleRotL8(b, 2)
	})
}

func encodeCurrentCastleSlot6TF(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 187)
		b = byte(uint16(b) * 195)
		b ^= 22
		b ^= 64
		return b
	})
}

func encodeCurrentCastleSlot6Component0(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 57
		b += 39
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleSlot6Component1(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b = byte(uint16(b) * 99)
		b = byte(uint16(b) * 229)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleSlot6Component2(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 231
		b ^= 237
		b = byte(uint16(b) * 213)
		return b
	})
}

func encodeCurrentCastleSlot6Component5Fallback(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 239
		b ^= 202
		b += 206
		return b
	})
}

func encodeCurrentCastleSlot6UZ(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 114
		b = currentCastleRotL8(b, 2)
		b = byte(uint16(b) * 193)
		b ^= 47
		return b
	})
}

func encodeCurrentCastleSlot6IH(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 205)
		b = byte(uint16(b) * 29)
		b += 8
		return b
	})
}

func encodeCurrentCastleSlot6FJ(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b ^= 188
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleSlot6IB(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 51)
		b += 205
		b ^= 155
		b ^= 189
		return b
	})
}

func encodeCurrentCastleSlot6AH(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = currentCastleRotL8(b, 7)
		return byte(uint16(b) * 63)
	})
}

func encodeCurrentCastleSlot6IT(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 43)
		b ^= 227
		b = currentCastleRotL8(b, 4)
		return byte(uint16(b) * 179)
	})
}

func encodeCurrentCastleSlot6AG(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 230
		b ^= 224
		b += 233
		return byte(uint16(b) * 33)
	})
}

func encodeCurrentCastleSlot6AD(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 3
		b = currentCastleRotL8(b, 4)
		return byte(uint16(b) * 107)
	})
}

func encodeCurrentCastleSlot6T7(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 85
		b = byte(uint16(b) * 33)
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleSlot6ID(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 85
		b = currentCastleRotL8(b, 7)
		b = currentCastleRotL8(b, 4)
		return currentCastleRotL8(b, 2)
	})
}

func encodeCurrentCastleSlot6FU(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 63)
		b = currentCastleRotL8(b, 4)
		b = currentCastleRotL8(b, 2)
		b += 45
		return b
	})
}

func encodeCurrentCastleSlot6EJ(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 203
		b = currentCastleRotL8(b, 3)
		b = currentCastleRotL8(b, 6)
		b ^= 191
		return b
	})
}

func encodeCurrentCastleSlot6L9(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b += 90
		b += 98
		b ^= 190
		return b
	})
}

func encodeCurrentCastleSlot6TX(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 244
		b += 253
		b ^= 234
		return b
	})
}

func encodeCurrentCastleSlot6LI(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b = currentCastleRotL8(b, 2)
		b += 38
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleSlot6FL(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = currentCastleRotL8(b, 2)
		b += 230
		b ^= 145
		return b
	})
}

func encodeCurrentCastleSlot6IO(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 249)
		b += 215
		b ^= 29
		return b
	})
}

func encodeCurrentCastleSlot6F3(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 99)
		b = currentCastleRotL8(b, 3)
		b ^= 18
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleSlot6AJ(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 173)
		b += 38
		b += 59
		return currentCastleRotL8(b, 2)
	})
}

func encodeCurrentCastleSlot6Component17(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 252
		b += 214
		b += 43
		return b
	})
}

func encodeCurrentCastleSlot6Component19(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 133)
		b += 122
		b += 123
		return b
	})
}

func encodeCurrentCastleSlot6Component27(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = currentCastleRotL8(b, 6)
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleSlot6TM(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b = currentCastleRotL8(b, 4)
		b ^= 154
		return b
	})
}

func encodeCurrentCastleSlot6TUpperF(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 30
		b ^= 145
		b += 207
		return b
	})
}

func encodeCurrentCastleSlot6FS(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 109)
		b ^= 216
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleSlot6Component31(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = currentCastleRotL8(b, 3)
		b += 241
		return b
	})
}

func encodeCurrentCastleSlot6FF(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 73)
		b ^= 129
		b = currentCastleRotL8(b, 7)
		b += 228
		return b
	})
}

func encodeCurrentCastleSlot6AU(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 111)
		b = currentCastleRotL8(b, 3)
		b ^= 200
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleSlot6AQ(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 29)
		b = currentCastleRotL8(b, 2)
		b ^= 143
		return b
	})
}

func encodeCurrentCastleSlot6LW(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 27
		b += 255
		return byte(uint16(b) * 185)
	})
}

func encodeCurrentCastleSlot6TUpperM(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 137)
		b += 63
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleSlot6TH(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 69
		b += 178
		b += 241
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleSlot6IS(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 229
		b ^= 42
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleSlot6A8(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 196
		b ^= 10
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleSlot6LUpperM(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 107)
		b = byte(uint16(b) * 23)
		b = byte(uint16(b) * 221)
		b += 218
		return b
	})
}

func encodeCurrentCastleSlot6EA(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 190
		b += 177
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleSlot6E3(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 179)
		b ^= 115
		b += 11
		return b
	})
}

func encodeCurrentCastleSlot6T2(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 9
		b = byte(uint16(b) * 139)
		return byte(uint16(b) * 111)
	})
}

func encodeCurrentCastleSlot6TUpperX(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = currentCastleRotL8(b, 7)
		b ^= 13
		return b
	})
}

func encodeCurrentCastleSlot6L5(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b += 235
		b = currentCastleRotL8(b, 1)
		return byte(uint16(b) * 41)
	})
}

func encodeCurrentCastleSlot6Component46(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b ^= 151
		b ^= 95
		b += 189
		b += 23
		return b
	})
}

func encodeCurrentCastleSlot6TZ(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 223)
		b = currentCastleRotL8(b, 5)
		b ^= 159
		return b
	})
}

func encodeCurrentCastleSlot6TP(value string) string {
	return encodeCurrentCastleTransformedUTF8(value, func(b byte) byte {
		b = byte(uint16(b) * 133)
		b += 242
		b = byte(uint16(b) * 215)
		return byte(uint16(b) * 105)
	})
}

func currentCastleSlot6Hash(value string) uint32 {
	units := jsUTF16Units(value)
	hash := uint32(0)
	index := 0
	remaining := len(units) & 3
	blockEnd := len(units) - remaining
	blockMul := currentCastleIndexTG(0xea02e370)
	for index < blockEnd {
		k := uint32(units[index]&0xff) |
			uint32(units[index+1]&0xff)<<8 |
			uint32(units[index+2]&0xff)<<16 |
			uint32(units[index+3]&0xff)<<24
		index += 4
		k *= 0xcc9e2d51
		k = bits.RotateLeft32(k, 15)
		k *= blockMul
		hash ^= k
		hash = bits.RotateLeft32(hash, 13)
		hashTimesFive := hash * 5
		hash = (hashTimesFive & 0xffff) +
			currentCastleIndexTG(0x10bd0000) +
			((((hashTimesFive >> 16) + currentCastleIndexTG(8196848)) & 0xffff) << 16)
	}
	k := uint32(0)
	switch remaining {
	case 3:
		k ^= uint32(units[index+2]&0xff) << 16
		fallthrough
	case 2:
		k ^= uint32(units[index+1]&0xff) << 8
		fallthrough
	case 1:
		k ^= uint32(units[index] & 0xff)
		k *= 0xcc9e2d51
		k = bits.RotateLeft32(k, 15)
		k *= blockMul
		hash ^= k
	}
	hash ^= uint32(len(units))
	hash ^= hash >> 16
	hash *= 0x85ebca6b
	hash ^= hash >> 13
	hash *= 0xc2b2ae35
	hash ^= hash >> 16
	return hash
}

func currentCastleSlot6HashHex(value string) string {
	return strconv.FormatUint(uint64(currentCastleSlot6Hash(value)), 16)
}

func encodeCurrentCastleFloat64WithTransform(value float64, transform func(byte) byte) string {
	var raw [8]byte
	binary.BigEndian.PutUint64(raw[:], math.Float64bits(value))
	for i, b := range raw {
		raw[i] = transform(b)
	}
	return base64.StdEncoding.EncodeToString(raw[:])
}

func encodeCurrentCastleFloat64(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		return byte((uint16(b)*65)&0xff) ^ 0xc6
	})
}

func encodeCurrentCastleFloat64ArrayWithTransform(values []float64, transform func(byte) byte) string {
	raw := make([]byte, len(values)*8)
	for i, value := range values {
		binary.BigEndian.PutUint64(raw[i*8:], math.Float64bits(value))
	}
	for i, b := range raw {
		raw[i] = transform(b)
	}
	return base64.StdEncoding.EncodeToString(raw)
}

func encodeCurrentCastleArraySlot254(values []float64) string {
	return encodeCurrentCastleFloat64ArrayWithTransform(values, func(b byte) byte {
		b ^= 32
		b = byte(uint16(b) * 185)
		b = byte(uint16(b) * 109)
		b += 195
		return b
	})
}

func encodeCurrentCastleArraySlot302(values []float64) string {
	return encodeCurrentCastleFloat64ArrayWithTransform(values, func(b byte) byte {
		b = byte(uint16(b) * 167)
		b = currentCastleRotL8(b, 4)
		b ^= 192
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot1(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 155
		b += 129
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot3(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b += 233
		b += 224
		b ^= 167
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot4(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 45
		b = currentCastleRotL8(b, 2)
		b ^= 7
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot5(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 32
		b = currentCastleRotL8(b, 6)
		b += 89
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot9(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 123)
		b = byte(uint16(b) * 91)
		b += 89
		return byte(uint16(b) * 29)
	})
}

func encodeCurrentCastleLowerFloatSlot10(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 228
		b = byte(uint16(b) * 79)
		b ^= 50
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot12(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 129)
		b = currentCastleRotL8(b, 5)
		return byte(uint16(b) * 251)
	})
}

func encodeCurrentCastleLowerFloatSlot13(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b ^= 153
		b = currentCastleRotL8(b, 3)
		return byte(uint16(b) * 25)
	})
}

func encodeCurrentCastleLowerFloatSlot14(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 43)
		b = byte(uint16(b) * 107)
		b ^= 253
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot15(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b = currentCastleRotL8(b, 7)
		b ^= 30
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot16(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 166
		b = byte(uint16(b) * 235)
		b ^= 26
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot17(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 182
		b += 207
		return byte(uint16(b) * 175)
	})
}

func encodeCurrentCastleLowerFloatSlot18(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b = currentCastleRotL8(b, 7)
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot19(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 187)
		b = currentCastleRotL8(b, 6)
		return byte(uint16(b) * 161)
	})
}

func encodeCurrentCastleLowerFloatSlot21(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 42
		b += 51
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleLowerFloatSlot22(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 192
		b = currentCastleRotL8(b, 6)
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot23(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b += 29
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot24(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 185)
		b = currentCastleRotL8(b, 3)
		b += 243
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot29(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b = currentCastleRotL8(b, 3)
		b ^= 51
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot30(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b = currentCastleRotL8(b, 4)
		b ^= 172
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot32(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b += 213
		b ^= 45
		b += 225
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot33(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = currentCastleRotL8(b, 7)
		b = currentCastleRotL8(b, 1)
		b += 48
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot34(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 113)
		b ^= 57
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot35(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 243)
		b += 173
		return byte(uint16(b) * 111)
	})
}

func encodeCurrentCastleLowerFloatSlot36(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b += 152
		return byte(uint16(b) * 41)
	})
}

func encodeCurrentCastleLowerFloatSlot37(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 87
		b ^= 113
		b = byte(uint16(b) * 137)
		b += 216
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot38(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = currentCastleRotL8(b, 5)
		b += 51
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot39(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 115)
		b = currentCastleRotL8(b, 6)
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot43(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = byte(uint16(b) * 21)
		b ^= 171
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot45(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 203
		b ^= 29
		b += 175
		return byte(uint16(b) * 111)
	})
}

func encodeCurrentCastleLowerFloatSlot48(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 251)
		b = byte(uint16(b) * 97)
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot50(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 250
		b ^= 171
		return byte(uint16(b) * 137)
	})
}

func encodeCurrentCastleLowerFloatSlot51(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 84
		b = byte(uint16(b) * 59)
		b = currentCastleRotL8(b, 3)
		b ^= 251
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot52(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 123)
		b = currentCastleRotL8(b, 3)
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleLowerFloatSlot54(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 106
		b ^= 142
		b = byte(uint16(b) * 207)
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleLowerFloatSlot55(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = byte(uint16(b) * 11)
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleLowerFloatSlot57(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 72
		b = currentCastleRotL8(b, 1)
		return byte(uint16(b) * 163)
	})
}

func encodeCurrentCastleLowerFloatSlot58(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 63)
		b = currentCastleRotL8(b, 3)
		b ^= 92
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot59(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 14
		b += 79
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot61(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 243
		b = byte(uint16(b) * 39)
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleLowerFloatSlot62(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 221
		b += 214
		return byte(uint16(b) * 47)
	})
}

func encodeCurrentCastleLowerFloatSlot63(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 178
		b ^= 223
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot65(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 109)
		b = byte(uint16(b) * 17)
		b ^= 210
		b ^= 81
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot66(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 229
		b += 213
		return byte(uint16(b) * 27)
	})
}

func encodeCurrentCastleLowerFloatSlot69(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 11)
		b += 250
		b = byte(uint16(b) * 85)
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot70(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b += 144
		b ^= 112
		b ^= 41
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot74(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 177)
		b += 195
		b ^= 92
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot76(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b ^= 230
		b += 167
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot80(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 48
		b = currentCastleRotL8(b, 2)
		b += 80
		return byte(uint16(b) * 239)
	})
}

func encodeCurrentCastleLowerFloatSlot81(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 249)
		b = currentCastleRotL8(b, 6)
		b ^= 178
		b ^= 139
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot82(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = byte(uint16(b) * 59)
		b ^= 115
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot84(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = byte(uint16(b) * 75)
		b ^= 239
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot85(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 14
		b ^= 182
		b = currentCastleRotL8(b, 3)
		return byte(uint16(b) * 111)
	})
}

func encodeCurrentCastleLowerFloatSlot86(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 29)
		b = currentCastleRotL8(b, 2)
		b += 143
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot88(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 215
		b ^= 73
		b = currentCastleRotL8(b, 7)
		b += 115
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot91(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 203)
		b += 81
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleLowerFloatSlot92(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 69)
		b ^= 38
		b += 139
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot93(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 174
		b = currentCastleRotL8(b, 6)
		b = currentCastleRotL8(b, 6)
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleLowerFloatSlot94(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 222
		b ^= 114
		b ^= 81
		b += 90
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot95(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 94
		b += 229
		return byte(uint16(b) * 73)
	})
}

func encodeCurrentCastleLowerFloatSlot96(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b = byte(uint16(b) * 15)
		b ^= 28
		return byte(uint16(b) * 81)
	})
}

func encodeCurrentCastleLowerFloatSlot97(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 123
		b += 78
		b ^= 219
		return byte(uint16(b) * 121)
	})
}

func encodeCurrentCastleLowerFloatSlot101(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 76
		b = byte(uint16(b) * 109)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot105(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 229)
		b = currentCastleRotL8(b, 5)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot107(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 233)
		b = byte(uint16(b) * 97)
		b ^= 58
		b ^= 22
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot110(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 33
		b += 241
		return byte(uint16(b) * 69)
	})
}

func encodeCurrentCastleLowerFloatSlot111(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 227)
		b += 224
		return currentCastleRotL8(b, 2)
	})
}

func encodeCurrentCastleLowerFloatSlot112(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 63)
		b = byte(uint16(b) * 85)
		b = byte(uint16(b) * 3)
		return byte(uint16(b) * 233)
	})
}

func encodeCurrentCastleLowerFloatSlot114(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 141)
		b = currentCastleRotL8(b, 6)
		b ^= 105
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot117(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 69)
		b = currentCastleRotL8(b, 4)
		b ^= 3
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot118(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 209
		b += 207
		b ^= 221
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot122(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 4
		b += 22
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleLowerFloatSlot124(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = byte(uint16(b) * 47)
		b ^= 32
		b ^= 40
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot125(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b += 141
		b += 49
		b ^= 148
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot126(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 149)
		b ^= 55
		b = byte(uint16(b) * 25)
		return byte(uint16(b) * 115)
	})
}

func encodeCurrentCastleLowerFloatSlot129(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 175
		b += 72
		return byte(uint16(b) * 225)
	})
}

func encodeCurrentCastleLowerFloatSlot131(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b ^= 116
		b = currentCastleRotL8(b, 3)
		b ^= 159
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot132(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 89)
		b += 253
		b ^= 166
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot134(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 111)
		b = byte(uint16(b) * 209)
		b ^= 197
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot135(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 247
		b = byte(uint16(b) * 57)
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleLowerFloatSlot136(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 205
		b = currentCastleRotL8(b, 5)
		b += 74
		b += 130
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot138(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 151
		b ^= 160
		b = byte(uint16(b) * 83)
		b ^= 182
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot139(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 165)
		b = byte(uint16(b) * 203)
		b ^= 56
		b += 87
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot140(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 99)
		b = byte(uint16(b) * 239)
		b ^= 72
		return byte(uint16(b) * 87)
	})
}

func encodeCurrentCastleLowerFloatSlot141(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b = byte(uint16(b) * 195)
		b ^= 149
		b ^= 197
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot142(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 87
		b = currentCastleRotL8(b, 4)
		b = currentCastleRotL8(b, 4)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot143(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = byte(uint16(b) * 125)
		b ^= 31
		b += 42
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot144(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b = byte(uint16(b) * 75)
		b += 83
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot145(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 42
		b = currentCastleRotL8(b, 2)
		b = currentCastleRotL8(b, 5)
		b += 146
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot146(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 207)
		b ^= 207
		return byte(uint16(b) * 193)
	})
}

func encodeCurrentCastleLowerFloatSlot148(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b += 78
		b ^= 38
		return byte(uint16(b) * 167)
	})
}

func encodeCurrentCastleLowerFloatSlot151(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 211)
		b = currentCastleRotL8(b, 5)
		b ^= 234
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot153(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 7)
		b ^= 79
		b = currentCastleRotL8(b, 1)
		b ^= 230
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot155(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 103)
		b = currentCastleRotL8(b, 3)
		b ^= 82
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot157(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b ^= 144
		b ^= 38
		b += 99
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot158(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 57
		b = byte(uint16(b) * 205)
		b += 48
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot160(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 12
		b = currentCastleRotL8(b, 4)
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot161(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 145
		b = currentCastleRotL8(b, 1)
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleLowerFloatSlot162(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = byte(uint16(b) * 115)
		b ^= 55
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot163(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = currentCastleRotL8(b, 5)
		b = byte(uint16(b) * 81)
		b ^= 106
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot164(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b += 153
		b = byte(uint16(b) * 73)
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot165(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 113)
		b += 232
		b ^= 33
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot166(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b ^= 210
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot167(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 121)
		b = byte(uint16(b) * 33)
		b += 13
		b += 50
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot169(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 139)
		b = currentCastleRotL8(b, 3)
		b += 138
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot170(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b += 249
		b += 200
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot171(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 78
		b += 149
		return byte(uint16(b) * 205)
	})
}

func encodeCurrentCastleLowerFloatSlot178(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 187
		b = byte(uint16(b) * 83)
		b ^= 87
		return byte(uint16(b) * 145)
	})
}

func encodeCurrentCastleLowerFloatSlot179(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 71)
		b += 123
		b += 99
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot181(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b = currentCastleRotL8(b, 5)
		b = byte(uint16(b) * 153)
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot185(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 41
		b = byte(uint16(b) * 207)
		b = currentCastleRotL8(b, 2)
		b ^= 86
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot186(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 129)
		b = currentCastleRotL8(b, 6)
		b = byte(uint16(b) * 147)
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleLowerFloatSlot187(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b += 53
		b ^= 133
		b += 236
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot189(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 91)
		b = currentCastleRotL8(b, 6)
		b += 183
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot190(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 66
		b += 41
		b ^= 35
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot191(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 41
		b = currentCastleRotL8(b, 6)
		b ^= 212
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot192(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 177)
		b = byte(uint16(b) * 27)
		b = currentCastleRotL8(b, 3)
		b ^= 3
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot193(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 203)
		b = currentCastleRotL8(b, 7)
		b ^= 126
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot195(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 43
		b ^= 237
		return byte(uint16(b) * 63)
	})
}

func encodeCurrentCastleLowerFloatSlot197(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 150
		b += 6
		return byte(uint16(b) * 29)
	})
}

func encodeCurrentCastleLowerFloatSlot199(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 101
		b = currentCastleRotL8(b, 7)
		b += 117
		return currentCastleRotL8(b, 2)
	})
}

func encodeCurrentCastleLowerFloatSlot200(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 30
		b = currentCastleRotL8(b, 2)
		b = currentCastleRotL8(b, 6)
		return byte(uint16(b) * 115)
	})
}

func encodeCurrentCastleLowerFloatSlot201(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b = currentCastleRotL8(b, 1)
		b ^= 212
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot202(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 3
		b = currentCastleRotL8(b, 2)
		b ^= 180
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot208(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 162
		b ^= 221
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot210(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 133)
		b = byte(uint16(b) * 111)
		return byte(uint16(b) * 59)
	})
}

func encodeCurrentCastleLowerFloatSlot211(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 160
		b = currentCastleRotL8(b, 6)
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot215(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = byte(uint16(b) * 5)
		b = byte(uint16(b) * 193)
		b ^= 105
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot216(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 89
		b = currentCastleRotL8(b, 1)
		b = byte(uint16(b) * 87)
		return byte(uint16(b) * 255)
	})
}

func encodeCurrentCastleLowerFloatSlot218(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 173
		b ^= 224
		b += 143
		return byte(uint16(b) * 5)
	})
}

func encodeCurrentCastleLowerFloatSlot219(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 232
		b ^= 173
		b = byte(uint16(b) * 173)
		b += 18
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot222(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 159)
		b = byte(uint16(b) * 19)
		b += 207
		b ^= 44
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot225(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 249)
		b = currentCastleRotL8(b, 3)
		b += 150
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot226(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 217)
		b = currentCastleRotL8(b, 6)
		b ^= 58
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleLowerFloatSlot227(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b ^= 244
		b = currentCastleRotL8(b, 2)
		b += 169
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot229(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = currentCastleRotL8(b, 3)
		b = byte(uint16(b) * 199)
		return byte(uint16(b) * 75)
	})
}

func encodeCurrentCastleLowerFloatSlot230(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 25
		b ^= 255
		b += 54
		return byte(uint16(b) * 161)
	})
}

func encodeCurrentCastleLowerFloatSlot232(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 127)
		b = currentCastleRotL8(b, 4)
		return byte(uint16(b) * 229)
	})
}

func encodeCurrentCastleLowerFloatSlot233(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 129)
		b = currentCastleRotL8(b, 5)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot234(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 25
		b += 147
		b ^= 86
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot236(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 132
		b += 58
		b = byte(uint16(b) * 33)
		return byte(uint16(b) * 3)
	})
}

func encodeCurrentCastleLowerFloatSlot238(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 17
		b = currentCastleRotL8(b, 6)
		b ^= 92
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot241(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 219
		b = currentCastleRotL8(b, 7)
		b += 166
		b ^= 60
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot242(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = byte(uint16(b) * 199)
		b ^= 46
		b ^= 241
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot243(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 1
		b = byte(uint16(b) * 179)
		b ^= 215
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot244(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 71)
		b = currentCastleRotL8(b, 3)
		b ^= 50
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot248(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 93)
		b = currentCastleRotL8(b, 2)
		b ^= 183
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot249(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = byte(uint16(b) * 181)
		b ^= 146
		b ^= 126
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot250(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = currentCastleRotL8(b, 5)
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot251(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b = byte(uint16(b) * 191)
		b += 10
		b += 54
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot252(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 158
		b ^= 101
		b += 127
		b ^= 83
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot253(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 101)
		b += 40
		return byte(uint16(b) * 239)
	})
}

func encodeCurrentCastleLowerFloatSlot255(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 91
		b ^= 143
		b += 94
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot259(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 145
		b ^= 133
		b ^= 110
		b += 19
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot262(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = byte(uint16(b) * 207)
		return byte(uint16(b) * 25)
	})
}

func encodeCurrentCastleLowerFloatSlot263(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 199
		b = currentCastleRotL8(b, 5)
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot264(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 21
		b = byte(uint16(b) * 161)
		b ^= 59
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot265(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 252
		b = byte(uint16(b) * 21)
		return byte(uint16(b) * 191)
	})
}

func encodeCurrentCastleLowerFloatSlot268(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 157)
		b = byte(uint16(b) * 247)
		b = currentCastleRotL8(b, 4)
		return byte(uint16(b) * 149)
	})
}

func encodeCurrentCastleLowerFloatSlot269(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b += 65
		b ^= 148
		b ^= 102
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot270(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b = byte(uint16(b) * 125)
		b += 245
		b ^= 159
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot271(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 117
		b += 71
		b += 25
		b ^= 200
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot272(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 51
		b = currentCastleRotL8(b, 3)
		b ^= 165
		b ^= 228
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot275(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b = byte(uint16(b) * 181)
		b += 161
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot276(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = currentCastleRotL8(b, 6)
		b ^= 78
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot278(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b ^= 154
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot280(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b += 133
		b += 187
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot282(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b = currentCastleRotL8(b, 5)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot283(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 193
		b = currentCastleRotL8(b, 4)
		b += 80
		b ^= 85
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot284(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 207)
		b = byte(uint16(b) * 149)
		b = currentCastleRotL8(b, 5)
		return byte(uint16(b) * 39)
	})
}

func encodeCurrentCastleLowerFloatSlot285(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 219)
		b += 238
		b ^= 14
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot286(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b = currentCastleRotL8(b, 6)
		b ^= 219
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot287(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 49
		b += 28
		b += 47
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot291(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = byte(uint16(b) * 229)
		b = currentCastleRotL8(b, 7)
		b ^= 134
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot294(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 109)
		b = byte(uint16(b) * 25)
		b ^= 52
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot295(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 179)
		b = currentCastleRotL8(b, 1)
		return byte(uint16(b) * 191)
	})
}

func encodeCurrentCastleLowerFloatSlot296(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 141)
		b = currentCastleRotL8(b, 2)
		b += 84
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot299(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 245
		b = byte(uint16(b) * 243)
		return byte(uint16(b) * 109)
	})
}

func encodeCurrentCastleLowerFloatSlot300(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 37)
		b = byte(uint16(b) * 87)
		b += 174
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot301(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b = byte(uint16(b) * 43)
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot303(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b = byte(uint16(b) * 107)
		b = byte(uint16(b) * 179)
		b += 248
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot305(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 160
		b += 202
		b += 29
		b ^= 243
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot306(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b += 131
		b ^= 221
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot307(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 3)
		b = currentCastleRotL8(b, 2)
		b ^= 194
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot309(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 251
		b = byte(uint16(b) * 81)
		b ^= 22
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot310(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 51)
		b = currentCastleRotL8(b, 4)
		b += 226
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot313(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 151)
		b = byte(uint16(b) * 71)
		b ^= 243
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot314(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = currentCastleRotL8(b, 1)
		b ^= 77
		b ^= 163
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot315(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b = currentCastleRotL8(b, 3)
		b += 240
		b += 186
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot317(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 96
		b = currentCastleRotL8(b, 7)
		b = currentCastleRotL8(b, 7)
		b += 126
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot319(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 173)
		b ^= 52
		b = byte(uint16(b) * 189)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot320(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 138
		b = byte(uint16(b) * 143)
		b += 46
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot322(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 237)
		b = byte(uint16(b) * 203)
		b = byte(uint16(b) * 17)
		b += 249
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot323(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b = byte(uint16(b) * 31)
		return byte(uint16(b) * 163)
	})
}

func encodeCurrentCastleLowerFloatSlot324(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 20
		b = currentCastleRotL8(b, 2)
		b ^= 110
		b += 28
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot325(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 15
		b ^= 6
		b += 17
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot327(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b = currentCastleRotL8(b, 1)
		return byte(uint16(b) * 53)
	})
}

func encodeCurrentCastleLowerFloatSlot331(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 23)
		b += 235
		b = byte(uint16(b) * 211)
		b += 29
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot332(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 231
		b += 175
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot333(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 136
		b = byte(uint16(b) * 91)
		b ^= 213
		b ^= 205
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot335(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b = currentCastleRotL8(b, 1)
		b += 136
		b += 79
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot336(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 129)
		b = byte(uint16(b) * 133)
		b ^= 60
		return byte(uint16(b) * 229)
	})
}

func encodeCurrentCastleLowerFloatSlot337(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 79
		b = byte(uint16(b) * 139)
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot338(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 97)
		b = currentCastleRotL8(b, 6)
		b = byte(uint16(b) * 139)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot340(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 186
		b += 193
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot343(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 255)
		b += 96
		b ^= 170
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot345(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 158
		b ^= 96
		b = byte(uint16(b) * 161)
		b ^= 179
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot346(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 80
		b = currentCastleRotL8(b, 4)
		return byte(uint16(b) * 55)
	})
}

func encodeCurrentCastleLowerFloatSlot347(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 11)
		b += 91
		return byte(uint16(b) * 215)
	})
}

func encodeCurrentCastleLowerFloatSlot353(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 159)
		b += 188
		return byte(uint16(b) * 119)
	})
}

func encodeCurrentCastleLowerFloatSlot355(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 163
		b = currentCastleRotL8(b, 5)
		b ^= 215
		b ^= 189
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot357(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 31)
		b = byte(uint16(b) * 153)
		b ^= 78
		b += 147
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot358(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 26
		b ^= 187
		return byte(uint16(b) * 215)
	})
}

func encodeCurrentCastleLowerFloatSlot360(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 137)
		b += 245
		return byte(uint16(b) * 143)
	})
}

func encodeCurrentCastleLowerFloatSlot361(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 160
		b = currentCastleRotL8(b, 7)
		b = byte(uint16(b) * 131)
		return byte(uint16(b) * 35)
	})
}

func encodeCurrentCastleLowerFloatSlot362(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 169
		b = currentCastleRotL8(b, 6)
		b = byte(uint16(b) * 27)
		b ^= 151
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot363(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 107)
		b = currentCastleRotL8(b, 5)
		b ^= 21
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot364(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 23)
		b = byte(uint16(b) * 195)
		return currentCastleRotL8(b, 2)
	})
}

func encodeCurrentCastleLowerFloatSlot365(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b ^= 74
		return byte(uint16(b) * 115)
	})
}

func encodeCurrentCastleLowerFloatSlot367(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 83
		b = byte(uint16(b) * 31)
		b = currentCastleRotL8(b, 6)
		b ^= 231
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot368(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = currentCastleRotL8(b, 4)
		return byte(uint16(b) * 97)
	})
}

func encodeCurrentCastleLowerFloatSlot370(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 173
		b ^= 250
		b = currentCastleRotL8(b, 3)
		b ^= 27
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot371(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 37)
		b += 133
		b ^= 43
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot372(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b = byte(uint16(b) * 37)
		b ^= 55
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot373(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 244
		b += 152
		return byte(uint16(b) * 3)
	})
}

func encodeCurrentCastleLowerFloatSlot377(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b ^= 4
		b += 114
		b += 130
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot379(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 107)
		b ^= 139
		b += 196
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot380(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 181
		b += 115
		b ^= 57
		b += 171
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot381(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 95)
		b = byte(uint16(b) * 249)
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot382(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 6)
		b = byte(uint16(b) * 33)
		b += 71
		b ^= 6
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot384(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 52
		b += 254
		b += 51
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot385(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 233)
		b = currentCastleRotL8(b, 3)
		b = currentCastleRotL8(b, 4)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot387(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 225)
		b = byte(uint16(b) * 65)
		b += 208
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot389(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 4)
		b ^= 144
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot391(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b = currentCastleRotL8(b, 1)
		b = currentCastleRotL8(b, 7)
		return currentCastleRotL8(b, 2)
	})
}

func encodeCurrentCastleLowerFloatSlot392(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 35)
		b = byte(uint16(b) * 205)
		b = byte(uint16(b) * 105)
		b += 145
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot396(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 162
		b += 249
		b += 171
		b ^= 29
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot397(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 112
		b ^= 131
		b += 38
		b ^= 110
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot398(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 68
		b = byte(uint16(b) * 133)
		b += 39
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot401(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 89)
		b ^= 92
		b = currentCastleRotL8(b, 6)
		return currentCastleRotL8(b, 4)
	})
}

func encodeCurrentCastleLowerFloatSlot402(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 7)
		b ^= 218
		b += 151
		b ^= 110
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot403(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 113
		b = byte(uint16(b) * 231)
		b ^= 241
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot404(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 251)
		b += 93
		return currentCastleRotL8(b, 6)
	})
}

func encodeCurrentCastleLowerFloatSlot405(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 228
		b = currentCastleRotL8(b, 3)
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot407(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 157)
		b ^= 152
		b ^= 46
		return currentCastleRotL8(b, 2)
	})
}

func encodeCurrentCastleLowerFloatSlot408(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 107)
		b = currentCastleRotL8(b, 1)
		b += 35
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot409(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 202
		b = currentCastleRotL8(b, 1)
		return currentCastleRotL8(b, 3)
	})
}

func encodeCurrentCastleLowerFloatSlot410(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 96
		b = currentCastleRotL8(b, 5)
		b = currentCastleRotL8(b, 7)
		b ^= 94
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot411(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 199
		b += 9
		b ^= 29
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot412(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 5)
		b ^= 217
		b = currentCastleRotL8(b, 1)
		b ^= 83
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot413(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b += 238
		b += 192
		return currentCastleRotL8(b, 7)
	})
}

func encodeCurrentCastleLowerFloatSlot414(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 148
		b ^= 145
		b += 118
		b += 53
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot417(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b ^= 134
		b += 40
		b = byte(uint16(b) * 19)
		b += 138
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot418(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b += 47
		b += 187
		return currentCastleRotL8(b, 1)
	})
}

func encodeCurrentCastleLowerFloatSlot419(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 1)
		b += 92
		b += 6
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot420(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 107)
		b = byte(uint16(b) * 255)
		return currentCastleRotL8(b, 5)
	})
}

func encodeCurrentCastleLowerFloatSlot426(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = currentCastleRotL8(b, 2)
		b += 188
		b = byte(uint16(b) * 81)
		b += 151
		return b
	})
}

func encodeCurrentCastleLowerFloatSlot427(value float64) string {
	return encodeCurrentCastleFloat64WithTransform(value, func(b byte) byte {
		b = byte(uint16(b) * 135)
		b = byte(uint16(b) * 195)
		b += 135
		return b
	})
}

func currentCastleIndexR7(value uint32) uint32 {
	value = bits.RotateLeft32(value, -16)
	value -= 762
	value = bits.RotateLeft32(value, 3)
	value -= 836
	value += 144
	return value
}

func currentCastleIndexEUnderscore(value uint32) uint32 {
	value -= 144
	value += 836
	value = bits.RotateLeft32(value, -3)
	value += 762
	value = bits.RotateLeft32(value, 16)
	return value
}

func currentCastleIndexTG(value uint32) uint32 {
	value -= 956
	value -= 24
	value += 42
	value += 41
	value += 478
	value -= 433
	return value * 0x70586723
}

func currentCastleTB(value uint32) uint32 {
	value = value*139 + 433
	value -= 478
	value -= 41
	value -= 42
	value += 24
	value += 956
	return value
}

func invertCurrentCastleTB(value uint32) uint32 {
	value -= 956
	value -= 24
	value += 42
	value += 41
	value += 478
	value -= 433
	return value * 1884841763
}

func encodeCurrentCastleNumericSlot20(value uint32) uint32 {
	value = currentCastleIndexEUnderscore(value)
	return (0x9f5593c9^value)*0x5d546885 + 0x312038dd
}

func encodeCurrentCastleNumericSlot154(value uint32) uint32 {
	value = currentCastleTB(value)
	return (0x3cf29b3d^value)*0x5bc9f085 + 0x786d9417
}

func encodeCurrentCastleNumericSlot256(value uint32) uint32 {
	return (0x40ad2308^value)*0x6b0c1905 + 0xbb6d942e
}

func encodeCurrentCastleNumericSlot348(value uint32) uint32 {
	return (0x00076853^value)*0x0438488b + 0x1ca480fc
}

func currentCastleNumericAffineEncoder(xorKey, multiplier, addend uint32) func(uint32) uint32 {
	return func(value uint32) uint32 {
		return (xorKey^value)*multiplier + addend
	}
}

func encodeCurrentCastleNumericIdentity(value uint32) uint32 {
	return value
}

func encodeCurrentCastleNumericSlot341(value uint32) uint32 {
	if value != 0 {
		return currentCastleIndexR7(0x12780000)
	}
	return currentCastleIndexTG(8777729)
}

func encodeCurrentCastleNumericSlot428(value uint32) uint32 {
	value = currentCastleTB(value)
	return (0xd57437b0^value)*0x62d22c71 + 0x0e965777
}

func encodeCurrentCastleNumericSlot429(value uint32) uint32 {
	return (0xb280aa5f^value)*0x23ad12b7 + currentCastleIndexR7(0x95eb36d0)
}

func currentCastleHighTimingSlotIndexes() []int {
	return []int{
		431, 432, 433, 434, 435,
		int(currentCastleIndexTG(61456)),
		437, 438,
		int(currentCastleIndexTG(61873)),
		440,
		int(currentCastleIndexTG(62151)),
		int(currentCastleIndexR7(0x0387c000)),
		443,
		int(currentCastleIndexR7(0x03880000)),
		445, 446, 447, 448, 449,
		int(currentCastleIndexTG(63402)),
		451, 452, 453,
		int(currentCastleIndexR7(0x03894000)),
		455, 456, 457, 458, 459,
		int(currentCastleIndexR7(0x038a0000)),
		461, 462, 463,
		int(currentCastleIndexTG(65348)),
		465, 466, 467, 468,
		int(currentCastleIndexR7(0x038b2000)),
		470, 471,
		int(currentCastleIndexTG(66460)),
		473, 474, 475, 476,
		int(currentCastleIndexTG(67155)),
		478, 479, 480, 481, 482,
		int(currentCastleIndexTG(67989)),
		int(currentCastleIndexTG(68128)),
		485,
		int(currentCastleIndexR7(0x038d4000)),
		int(currentCastleIndexR7(0x038d6000)),
		int(currentCastleIndexR7(0x038d8000)),
		int(currentCastleIndexR7(0x038da000)),
		490, 491, 492,
		int(currentCastleIndexR7(0x038e2000)),
	}
}

func populateCurrentCastleHighTimingSlots(payload []any, values map[int]float64) error {
	if len(payload) < currentCastlePayloadSlots {
		return fmt.Errorf("current Castle payload has %d slots, want at least %d", len(payload), currentCastlePayloadSlots)
	}
	defaults := defaultCurrentCastleHighTimingSlotEncodedValues()
	for _, slot := range currentCastleHighTimingSlotIndexes() {
		if value, ok := values[slot]; ok {
			payload[slot] = encodeCurrentCastleFloat64(value)
		} else if encoded, ok := defaults[slot]; ok {
			payload[slot] = encoded
		} else {
			payload[slot] = encodeCurrentCastleFloat64(0)
		}
	}
	return nil
}

func populateCurrentCastleObservedStringSlots(payload []any, input currentCastlePayloadInput) {
	for slot, encoded := range defaultCurrentCastleObservedStringSlotEncodedValues() {
		if currentCastleInputHasStringSlotOverride(input, slot) {
			continue
		}
		payload[slot] = encoded
	}
}

func currentCastleInputHasStringSlotOverride(input currentCastlePayloadInput, slot int) bool {
	if _, ok := input.LowerFloatValues[slot]; ok {
		return true
	}
	if _, ok := input.ArrayValues[slot]; ok {
		return true
	}
	if _, ok := input.PackedStringValues[slot]; ok {
		return true
	}
	if _, ok := input.UnitPackedStrings[slot]; ok {
		return true
	}
	if _, ok := input.StringValues[slot]; ok {
		return true
	}
	if _, ok := input.HighTimingValues[slot]; ok {
		return true
	}
	return false
}

func defaultCurrentCastleObservedStringSlotEncodedValues() map[int]string {
	return map[int]string{
		// Current Chrome begin_login string payload defaults observed from the
		// native Windows profile. These cover lower browser-probe string slots,
		// one compact array slot, and high-timing slots that are not stable zeroes.
		25:  "NoAtnP9LeAQ=",
		34:  "vXZOTk5OTk4=",
		48:  "BukFAAAAAAA=",
		60:  "NoAtnP9LeAQ=",
		75:  "NoAtnP9LeAQ=",
		78:  "3Nzc3Nzc3Nw=",
		79:  "PT09PT09PT0=",
		83:  "5bKenp6enp4=",
		86:  "kE2Rj4+Pj48=",
		91:  "ZBoVFRUVFRU=",
		97:  "gfKCgoKCgoI=",
		110: "Gj062tra2to=",
		113: "NoAtnP9LeAQ=",
		120: "z8/Pz8/Pz88=",
		122: "K02P7J1heMI=",
		135: "D+3r//////8=",
		137: "NoAtnP9LeAQ=",
		148: "Jk/Y2NjY2Ng=",
		185: "VlZWVlZWVlY=",
		187: "PRScnJycnJw=",
		190: "iPkkFBQUFBs=",
		191: "1NTU1NTU1NQ=",
		204: "z8/Pz8/Pz88=",
		208: "4XnOkU1ytBo=",
		215: "DbSN7GNpaWk=",
		218: "DDxMTExMTEw=",
		222: "bh7W/YRTetY=",
		239: "NoAtnP9LeAQ=",
		242: "ThnLgCnt7EY=",
		243: "F4LS0tLS0tI=",
		249: "vDuC8vLy8rI=",
		262: "uN6u6urq6uo=",
		283: "NWYCT+EYstQ=",
		303: "tCv4+Pj4+Pg=",
		314: "/s0XIiIiIiI=",
		315: "pusz3/VxUJ0=",
		327: "jerv8g6b8dM=",
		328: "NoAtnP9LeAQ=",
		332: "daGlpaWlpaU=",
		345: "bv3T2KP8NIk=",
		350: "XFyWlpaWlpY=",
		353: "pCkEZGRkZGQ=",
		370: "qyChoaGhoaE=",
		387: "78DQ0NDQ0NA=",
		392: "cpmRkZGRkZE=",
		394: "zESL",
		396: "mIu1jgUq8l8=",
		397: "PGd3d3d3d3c=",
		407: "VJqyPvHtTOg=",
		411: "Ej0vAL/cZNE=",
		436: "OT8fXsbGxsY=",
		442: "hinGxsbGxsY=",
		443: "OT8f3MbGxsY=",
		445: "xsbGxsbGxsY=",
		447: "OT8f3MbGxsY=",
		449: "OT8f3MbGxsY=",
		456: "OaAgIIbGxsY=",
		459: "OT8f3MbGxsY=",
		460: "xsbGxsbGxsY=",
		462: "xsbGxsbGxsY=",
		466: "xsbGxsbGxsY=",
		468: "OT8fXsbGxsY=",
		476: "hpMfH17GxsY=",
		478: "OT8fXsbGxsY=",
		479: "Oc8f3MbGxsY=",
		480: "huGKCgoGxsY=",
		481: "hmt1NTWGxsY=",
		482: "hqqwICBGxsY=",
		483: "huVGxsbGxsY=",
		484: "huUPHx9GxsY=",
		485: "hqLFNTWGxsY=",
		486: "hqJPHx8GxsY=",
		487: "hlOfHx9GxsY=",
		488: "hpL/Hx9GxsY=",
		489: "OT8f3MbGxsY=",
		490: "hlM/Hx9GxsY=",
		491: "OVU1NcbGxsY=",
		492: "hhhKCsvGxsY=",
		493: "huWaCsvGxsY=",
	}
}

func defaultCurrentCastleHighTimingSlotEncodedValues() map[int]string {
	return map[int]string{
		// Current Chrome begin_login high-timing defaults observed on this
		// Windows profile. Other high-timing slots encode zero by default.
		442: "hoZgICDGxsY=",
		443: "OVU1NcbGxsY=",
		445: "OT8f3MbGxsY=",
		449: "OT8fXsbGxsY=",
		456: "OWU1NYbGxsY=",
		460: "OT8f3MbGxsY=",
		462: "OT8f3MbGxsY=",
		466: "OT8f3MbGxsY=",
		476: "hmQKCgrGxsY=",
		478: "OT8f3MbGxsY=",
		479: "OT8f3MbGxsY=",
		480: "hnMPHx9mxsY=",
		481: "hn8XHx9GxsY=",
		482: "hn/tNTWGxsY=",
		483: "hrZXHx9GxsY=",
		484: "hrZfHx9GxsY=",
		485: "hrbXHx9mxsY=",
		486: "hndAICBGxsY=",
		487: "hu1gICBGxsY=",
		488: "hu0VNTWGxsY=",
		489: "Oc8f3MbGxsY=",
		490: "hu1qCgoGxsY=",
		491: "Oc8f3MbGxsY=",
		492: "hiz6CgoGxsY=",
		493: "hrYuxsbGxsY=",
	}
}

func xorCurrentCastleString(value, key string) string {
	valueUnits := jsUTF16Units(value)
	keyUnits := jsUTF16Units(key)
	if len(valueUnits) == 0 || len(keyUnits) == 0 {
		return value
	}
	out := make([]uint16, len(valueUnits))
	for i, ch := range valueUnits {
		out[i] = ch ^ keyUnits[i%len(keyUnits)]
	}
	return string(utf16.Decode(out))
}

func deflateCurrentCastleWireString(value string) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.BestSpeed)
	if err != nil {
		return nil, err
	}
	if _, err = writer.Write([]byte(value)); err != nil {
		_ = writer.Close()
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func currentCastleRotate16Into32(value uint16, leftShift, rightShift uint) uint32 {
	return uint32(value)<<leftShift | uint32(value)>>rightShift
}

func bitswap32(value uint32) uint32 {
	return value>>24 | value>>8&0xff00 | value<<8&0xff0000 | value<<24
}

func jsUTF16Units(value string) []uint16 {
	return utf16.Encode([]rune(value))
}
