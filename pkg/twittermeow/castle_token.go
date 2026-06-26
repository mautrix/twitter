package twittermeow

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"net/url"
	"time"
)

const castlePublicKey = "pk_AvRa79bHyJSYSQHnRpcVtzyxetSvFerx"

const (
	castleValueB2H                 = 3
	castleValueSerializedByteArray = 4
	castleValueB2HWithChecks       = 5
	castleValueB2HRounded          = 6
	castleValueJustAppend          = 7
	castleValueEmpty               = -1
)

func addCastleTokenToForm(form url.Values) error {
	if form.Get("$castle_token") != "" {
		return nil
	}
	token, err := createCastleRequestToken()
	if err != nil {
		return err
	}
	form.Set("$castle_token", token)
	return nil
}

func createCastleRequestToken() (string, error) {
	initDelta, err := cryptoRandInt(2*60*1000, 30*60*1000+1)
	if err != nil {
		return "", err
	}
	initTime := time.Now().UnixMilli() - int64(initDelta)
	tokenUUID, err := randomBytes(16)
	if err != nil {
		return "", err
	}
	tokenUUIDHex := hex.EncodeToString(tokenUUID)

	fpOne, err := castleFPOne(initTime)
	if err != nil {
		return "", err
	}
	fpTwo, err := castleFPTwo(initTime)
	if err != nil {
		return "", err
	}
	fpThree, err := castleFPThree(initTime)
	if err != nil {
		return "", err
	}
	eventLog, err := castleEventLog()
	if err != nil {
		return "", err
	}
	eventValues, err := castleFPEventValues()
	if err != nil {
		return "", err
	}

	fpData := append(append(append(append(fpOne, fpTwo...), fpThree...), eventLog...), eventValues...)
	fpData = append(fpData, 0xff)
	fpDataKey := castleEncodeTimestampEncrypted(time.Now().UnixMilli())
	encryptedFPData := deriveKeyAndXORBytes(hex.EncodeToString(fpDataKey), 4, hex.EncodeToString(fpDataKey)[3], fpData)
	encryptedFPDataTwo := deriveKeyAndXORBytes(tokenUUIDHex, 8, tokenUUIDHex[9], append(fpDataKey, encryptedFPData...))

	header, err := castleHeader(tokenUUIDHex, initTime)
	if err != nil {
		return "", err
	}
	fingerprint := append(header, encryptedFPDataTwo...)
	encryptedFP := xxteaEncrypt(fingerprint, []uint32{1164413191, 3891440048, 185273099, 2746598870})
	encryptedFP = append([]byte{0x0b, byte(len(encryptedFP) - len(fingerprint))}, encryptedFP...)
	encryptedFP = append(encryptedFP, byte((len(encryptedFP)*2)&0xff))

	randomByte, err := cryptoRandInt(0, 256)
	if err != nil {
		return "", err
	}
	final := append([]byte{byte(randomByte)}, xorBytes(encryptedFP, []byte{byte(randomByte)})...)
	return base64.RawURLEncoding.EncodeToString(final), nil
}

func castleHeader(uuidHex string, initTime int64) ([]byte, error) {
	uuidBytes, err := hex.DecodeString(uuidHex)
	if err != nil {
		return nil, err
	}
	out := castleEncodeTimestampEncrypted(initTime)
	version := uint16((3 << 13) | (1 << 11) | (6 << 6))
	out = binary.BigEndian.AppendUint16(out, version)
	out = append(out, []byte(castlePublicKey)...)
	out = append(out, uuidBytes...)
	return out, nil
}

func castleEncodeTimestampEncrypted(timestampMillis int64) []byte {
	seconds := int64(math.Floor(float64(timestampMillis)/1e3 - 1535e6))
	if seconds < 0 {
		seconds = 0
	} else if seconds > 268435455 {
		seconds = 268435455
	}
	timeBytes := []byte{byte(seconds >> 24), byte(seconds >> 16), byte(seconds >> 8), byte(seconds)}
	msBytes := []byte{byte((timestampMillis % 1000) >> 8), byte(timestampMillis % 1000)}
	key, _ := cryptoRandInt(0, 16)
	return append(xorAndAppendNibbleKey(timeBytes, byte(key)), xorAndAppendNibbleKey(msBytes, byte(key))...)
}

func xorAndAppendNibbleKey(buf []byte, key byte) []byte {
	nibbles := make([]byte, 0, len(buf)*2)
	for _, b := range buf {
		nibbles = append(nibbles, b>>4, b&0x0f)
	}
	outNibbles := make([]byte, 0, len(nibbles))
	for _, nibble := range nibbles[1:] {
		outNibbles = append(outNibbles, nibble^key)
	}
	outNibbles = append(outNibbles, key)
	out := make([]byte, len(outNibbles)/2)
	for i := range out {
		out[i] = outNibbles[i*2]<<4 | outNibbles[i*2+1]
	}
	return out
}

func processCastleFPValue(index, valueType int, val []byte, intVal int, initTime int64) ([]byte, error) {
	out := []byte{byte(((31 & index) << 3) | (7 & valueType))}
	switch valueType {
	case castleValueB2HRounded, castleValueB2H:
		out = append(out, byte(intVal))
	case castleValueB2HWithChecks:
		if intVal <= 127 {
			out = append(out, byte(intVal))
		} else {
			out = binary.BigEndian.AppendUint16(out, uint16((1<<15)|(32767&intVal)))
		}
	case castleValueSerializedByteArray:
		encrypted := xxteaEncrypt(val, []uint32{uint32(index), uint32(initTime), 16373134, 643144773, 1762804430, 1186572681, 1164413191})
		out = append(out, byte(len(encrypted)))
		out = append(out, encrypted...)
	case castleValueJustAppend:
		out = append(out, val...)
	case castleValueEmpty:
	case 1, 2:
	default:
		return nil, fmt.Errorf("unsupported castle fp value type %d", valueType)
	}
	return out, nil
}

func castleValuesToBytes(first, second int) []byte {
	r := uint16(32767 & first)
	e := uint16(65535 & second)
	if r == e {
		return []byte{byte((32768 | int(r)) >> 8), byte(32768 | int(r))}
	}
	out := make([]byte, 4)
	binary.BigEndian.PutUint16(out[:2], r)
	binary.BigEndian.PutUint16(out[2:], e)
	return out
}

func castleEncodeBits(bits []int, bitSize int) []byte {
	numBytes := bitSize / 8
	out := make([]byte, numBytes)
	for _, bit := range bits {
		byteIndex := (numBytes - 1) - (bit / 8)
		bitPosition := bit % 8
		if byteIndex >= 0 && byteIndex < numBytes {
			out[byteIndex] |= 1 << bitPosition
		}
	}
	return out
}

func castleBoolArrayToBinary(values []bool, size int) int {
	e := values
	if size > 0 && len(values) > size {
		e = values[:size]
	}
	out := 0
	for i := len(e) - 1; i >= 0; i-- {
		if e[i] {
			out |= 1 << (len(e) - i - 1)
		}
	}
	if size > 0 && len(e) < size {
		out <<= size - len(e)
	}
	return out
}

func castleFPOne(initTime int64) ([]byte, error) {
	timezone := jetfuelTimezone()
	values := [][]byte{}
	addInt := func(index, valueType, value int) error {
		v, err := processCastleFPValue(index, valueType, nil, value, initTime)
		values = append(values, v)
		return err
	}
	addBytes := func(index, valueType int, value []byte) error {
		v, err := processCastleFPValue(index, valueType, value, 0, initTime)
		values = append(values, v)
		return err
	}
	if err := addInt(0, castleValueB2H, 1); err != nil {
		return nil, err
	}
	if err := addInt(1, castleValueB2H, 0); err != nil {
		return nil, err
	}
	if err := addBytes(2, castleValueSerializedByteArray, []byte("en-US")); err != nil {
		return nil, err
	}
	if err := addInt(3, castleValueB2HRounded, 80); err != nil {
		return nil, err
	}
	screenDims := append(castleValuesToBytes(1920, 1920), castleValuesToBytes(1080, 1032)...)
	if err := addBytes(4, castleValueJustAppend, screenDims); err != nil {
		return nil, err
	}
	if err := addInt(5, castleValueB2HWithChecks, 24); err != nil {
		return nil, err
	}
	if err := addInt(6, castleValueB2HWithChecks, 24); err != nil {
		return nil, err
	}
	if err := addInt(7, castleValueB2HRounded, 10); err != nil {
		return nil, err
	}
	if err := addBytes(8, castleValueJustAppend, castleTimezoneDiff(timezone)); err != nil {
		return nil, err
	}
	if err := addBytes(9, castleValueJustAppend, []byte{0x02, 0x7d, 0x5f, 0xc9, 0xa7}); err != nil {
		return nil, err
	}
	if err := addBytes(10, castleValueJustAppend, []byte{0x05, 0x72, 0x93, 0x02, 0x08}); err != nil {
		return nil, err
	}
	if err := addBytes(11, castleValueJustAppend, append([]byte{12}, castleEncodeBits([]int{0, 1, 2, 3, 4, 5, 6}, 16)...)); err != nil {
		return nil, err
	}
	ua := xxteaEncrypt([]byte(UserAgent), []uint32{12, uint32(initTime), 16373134, 643144773, 1762804430, 1186572681, 1164413191})
	if err := addBytes(12, castleValueJustAppend, append([]byte{1, byte(len(ua))}, ua...)); err != nil {
		return nil, err
	}
	if err := addBytes(13, castleValueSerializedByteArray, []byte("54b4b5cf")); err != nil {
		return nil, err
	}
	if err := addBytes(14, castleValueJustAppend, append([]byte{3}, castleEncodeBits([]int{0, 1, 2}, 8)...)); err != nil {
		return nil, err
	}
	if err := addInt(17, castleValueB2H, 0); err != nil {
		return nil, err
	}
	if err := addBytes(18, castleValueSerializedByteArray, []byte("c6749e76")); err != nil {
		return nil, err
	}
	gpu := "ANGLE (NVIDIA, NVIDIA GeForce RTX 3060 (0x00002504) Direct3D11 vs_5_0 ps_5_0, D3D11)"
	if err := addBytes(19, castleValueSerializedByteArray, []byte(gpu)); err != nil {
		return nil, err
	}
	if err := addBytes(20, castleValueSerializedByteArray, []byte(castleEpochLocaleString(timezone))); err != nil {
		return nil, err
	}
	if err := addBytes(21, castleValueJustAppend, append([]byte{8}, castleEncodeBits(nil, 8)...)); err != nil {
		return nil, err
	}
	if err := addInt(22, castleValueB2HWithChecks, 33); err != nil {
		return nil, err
	}
	if err := addInt(24, castleValueB2HWithChecks, 12549); err != nil {
		return nil, err
	}
	if err := addInt(25, castleValueB2H, 0); err != nil {
		return nil, err
	}
	if err := addInt(26, castleValueB2H, 1); err != nil {
		return nil, err
	}
	if err := addInt(27, castleValueB2HWithChecks, 4644); err != nil {
		return nil, err
	}
	if err := addBytes(28, castleValueJustAppend, []byte{0}); err != nil {
		return nil, err
	}
	if err := addInt(29, castleValueB2H, 3); err != nil {
		return nil, err
	}
	if err := addBytes(30, castleValueJustAppend, []byte{0x5d, 0xc5, 0xab, 0xb5, 0x88}); err != nil {
		return nil, err
	}
	if err := addBytes(31, castleValueJustAppend, []byte{0xa2, 0x6a}); err != nil {
		return nil, err
	}
	out := []byte{byte(31 & len(values))}
	for _, value := range values {
		out = append(out, value...)
	}
	return out, nil
}

func castleTimezoneDiff(timezone string) []byte {
	switch timezone {
	case "America/Chicago":
		return []byte{20, 4}
	case "America/Los_Angeles":
		return []byte{28, 4}
	case "America/New_York":
		return []byte{16, 4}
	case "America/Denver":
		return []byte{24, 4}
	case "America/Anchorage":
		return []byte{32, 4}
	case "Pacific/Honolulu":
		return []byte{40, 0}
	default:
		_, offset := time.Now().Zone()
		return []byte{byte((-offset / 60) / 15), 0}
	}
}

func castleEpochLocaleString(timezone string) string {
	switch timezone {
	case "America/Chicago":
		return "12/31/1969, 6:00:00 PM"
	case "America/Los_Angeles":
		return "12/31/1969, 4:00:00 PM"
	case "America/New_York":
		return "12/31/1969, 7:00:00 PM"
	case "America/Denver":
		return "12/31/1969, 5:00:00 PM"
	case "America/Anchorage":
		return "12/31/1969, 3:00:00 PM"
	case "Pacific/Honolulu":
		return "12/31/1969, 2:00:00 PM"
	default:
		return "1/1/1970, 12:00:00 AM"
	}
}

func castleFPTwo(initTime int64) ([]byte, error) {
	timezone := jetfuelTimezone()
	values := [][]byte{}
	addInt := func(index, valueType, value int) error {
		v, err := processCastleFPValue(index, valueType, nil, value, initTime)
		values = append(values, v)
		return err
	}
	addBytes := func(index, valueType int, value []byte) error {
		v, err := processCastleFPValue(index, valueType, value, 0, initTime)
		values = append(values, v)
		return err
	}
	if err := addInt(0, castleValueB2H, 0); err != nil {
		return nil, err
	}
	if enum, ok := map[string]int{
		"America/New_York":    0,
		"America/Sao_Paulo":   1,
		"America/Chicago":     2,
		"America/Los_Angeles": 3,
		"America/Mexico_City": 4,
		"Asia/Shanghai":       5,
	}[timezone]; ok {
		if err := addInt(1, castleValueB2H, enum); err != nil {
			return nil, err
		}
	} else if err := addBytes(1, castleValueSerializedByteArray, []byte(timezone)); err != nil {
		return nil, err
	}
	if err := addBytes(2, castleValueSerializedByteArray, []byte("en-US,en")); err != nil {
		return nil, err
	}
	if err := addInt(6, castleValueB2HWithChecks, 0); err != nil {
		return nil, err
	}
	if err := addBytes(10, castleValueJustAppend, append([]byte{4}, castleEncodeBits([]int{1, 2, 3}, 8)...)); err != nil {
		return nil, err
	}
	if err := addInt(12, castleValueB2HWithChecks, 80); err != nil {
		return nil, err
	}
	if err := addBytes(13, castleValueJustAppend, []byte{9, 0, 0}); err != nil {
		return nil, err
	}
	if err := addBytes(17, castleValueJustAppend, append([]byte{0x0d}, castleEncodeBits([]int{1, 5, 8, 9, 10}, 16)...)); err != nil {
		return nil, err
	}
	if err := addInt(18, 1, 0); err != nil {
		return nil, err
	}
	if err := addBytes(21, castleValueJustAppend, []byte{0, 0, 0, 0}); err != nil {
		return nil, err
	}
	if err := addBytes(22, castleValueSerializedByteArray, []byte("en-US")); err != nil {
		return nil, err
	}
	if err := addBytes(23, castleValueJustAppend, append([]byte{2}, castleEncodeBits([]int{0}, 8)...)); err != nil {
		return nil, err
	}
	heightDiff, err := cryptoRandInt(10, 31)
	if err != nil {
		return nil, err
	}
	if err := addBytes(24, castleValueJustAppend, []byte{0, 0, 0, byte(heightDiff)}); err != nil {
		return nil, err
	}
	out := []byte{byte((7&4)<<5 | (31 & len(values)))}
	for _, value := range values {
		out = append(out, value...)
	}
	return out, nil
}

func castleFPThree(initTime int64) ([]byte, error) {
	minute := time.UnixMilli(initTime).UTC().Minute()
	first, err := processCastleFPValue(3, castleValueB2HWithChecks, nil, 1, initTime)
	if err != nil {
		return nil, err
	}
	second, err := processCastleFPValue(4, castleValueB2HWithChecks, nil, minute, initTime)
	if err != nil {
		return nil, err
	}
	return append([]byte{byte((7 << 5) | 2)}, append(first, second...)...), nil
}

func castleEventLog() ([]byte, error) {
	simpleEvents := []int{21, 18, 25, 26, 27}
	targetEvents := []int{0, 6, 5}
	allEvents := append(simpleEvents, targetEvents...)
	count, err := cryptoRandInt(30, 71)
	if err != nil {
		return nil, err
	}
	events := []byte{}
	for range count {
		idx, err := cryptoRandInt(0, len(allEvents))
		if err != nil {
			return nil, err
		}
		eventID := allEvents[idx]
		if eventID == 0 || eventID == 6 || eventID == 5 {
			events = append(events, byte(eventID|128), 63)
		} else {
			events = append(events, byte(eventID))
		}
	}
	payload := append([]byte{0}, byte(count>>8), byte(count), 0)
	payload = payload[:3]
	payload = append(payload, events...)
	return append([]byte{byte(len(payload) >> 8), byte(len(payload))}, payload...), nil
}

func castleFPEventValues() ([]byte, error) {
	bits := make([]bool, 15)
	for _, bit := range []int{2, 3, 5, 6, 9, 11, 12} {
		bits[bit] = true
	}
	binaryNum := castleBoolArrayToBinary(bits, 16)
	encodedNum := (6 << 20) | (2 << 16) | (65535 & binaryNum)
	out := []byte{byte(encodedNum >> 16), byte(encodedNum >> 8), byte(encodedNum)}
	floatValues := []float64{
		randFloat(40, 50), -1, randFloat(70, 80), -1, randFloat(60, 70), -1,
		0, 0, randFloat(60, 80), randFloat(5, 10), randFloat(30, 40), randFloat(2, 5),
		-1, -1, -1, -1, -1, -1, -1, -1,
		randFloat(150, 180), randFloat(3, 6), randFloat(150, 180), randFloat(3, 6),
		randFloat(0, 2), randFloat(0, 2), 0, 0, -1, -1, -1, -1,
		0, 0, 0, 0, 0, 0, 1, 0, 1, 0,
		randFloat(0, 4), randFloat(0, 3), randFloat(25, 50), randFloat(25, 50),
		randFloat(25, 50), randFloat(25, 30), randFloat(0, 2), randFloat(0, 1),
		randFloat(0, 1), 1, 0,
	}
	for _, value := range floatValues {
		if value == -1 {
			out = append(out, 0)
		} else {
			out = append(out, byte(castleEncodeEventFloat(value)))
		}
	}
	mouseMove, err := cryptoRandInt(100, 200)
	if err != nil {
		return nil, err
	}
	keyUp, err := cryptoRandInt(1, 5)
	if err != nil {
		return nil, err
	}
	click, err := cryptoRandInt(1, 5)
	if err != nil {
		return nil, err
	}
	keyDown, err := cryptoRandInt(0, 5)
	if err != nil {
		return nil, err
	}
	wheel, err := cryptoRandInt(0, 5)
	if err != nil {
		return nil, err
	}
	unk1, err := cryptoRandInt(0, 11)
	if err != nil {
		return nil, err
	}
	out = append(out, byte(mouseMove), byte(keyUp), byte(click), 0, byte(keyDown), 0, 0, 0, byte(wheel), byte(unk1), 0, 11)
	return out, nil
}

func castleEncodeEventFloat(value float64) int {
	n := math.Max(value, 0)
	if n <= 15 {
		return 64 | castleCustomFloatEncode(2, 4, n+1)
	}
	return 128 | castleCustomFloatEncode(4, 3, n-14)
}

func castleCustomFloatEncode(expBits, manBits int, n float64) int {
	if n == 0 {
		return 0
	}
	exponent := 0
	if n < 0 {
		n = -n
	}
	for 2 <= n {
		n /= 2
		exponent++
	}
	for n < 1 {
		n *= 2
		exponent--
	}
	maxExponent := (1 << expBits) - 1
	if exponent > maxExponent {
		exponent = maxExponent
	}
	fractional := n - math.Floor(n)
	mantissaBits := 0
	position := 1
	for fractional > 0 && position <= manBits {
		fractional *= 2
		bit := int(math.Floor(fractional))
		mantissaBits |= bit << (manBits - position)
		fractional -= float64(bit)
		position++
	}
	return exponent<<manBits | mantissaBits
}

func deriveKeyAndXORBytes(key string, sliceLen int, rotationKey byte, data []byte) []byte {
	substring := key[:sliceLen]
	rotationNumber := int(hexNibble(rotationKey))
	rotationIndex := rotationNumber % len(substring)
	rotated := substring[rotationIndex:] + substring[:rotationIndex]
	keyBytes, _ := hex.DecodeString(rotated)
	return xorBytes(data, keyBytes)
}

func hexNibble(ch byte) byte {
	switch {
	case ch >= '0' && ch <= '9':
		return ch - '0'
	case ch >= 'a' && ch <= 'f':
		return ch - 'a' + 10
	case ch >= 'A' && ch <= 'F':
		return ch - 'A' + 10
	default:
		return 0
	}
}

func xorBytes(data, key []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = b ^ key[i%len(key)]
	}
	return out
}

func xxteaEncrypt(data []byte, key []uint32) []byte {
	blocks := bytesToXXTEABlocks(data)
	u := len(blocks) - 1
	if u < 1 {
		return data
	}
	sum := uint32(0)
	o := blocks[u]
	rounds := 6 + 52/(u+1)
	for rounds > 0 {
		rounds--
		sum += 0x9e3779b9
		e := (sum >> 2) & 3
		for c := 0; c < u; c++ {
			a := blocks[c+1]
			blocks[c] += (((o >> 5) ^ (a << 2)) + ((a >> 3) ^ (o << 4))) ^ ((sum ^ a) + (key[(c&3)^int(e)] ^ o))
			o = blocks[c]
		}
		a := blocks[0]
		blocks[u] += (((o >> 5) ^ (a << 2)) + ((a >> 3) ^ (o << 4))) ^ ((sum ^ a) + (key[(u&3)^int(e)] ^ o))
		o = blocks[u]
	}
	out := make([]byte, len(blocks)*4)
	for i, block := range blocks {
		binary.LittleEndian.PutUint32(out[i*4:], block)
	}
	return out
}

func bytesToXXTEABlocks(data []byte) []uint32 {
	count := (len(data) + 3) / 4
	padded := make([]byte, count*4)
	copy(padded, data)
	blocks := make([]uint32, count)
	for i := range blocks {
		blocks[i] = binary.LittleEndian.Uint32(padded[i*4:])
	}
	return blocks
}

func randFloat(min, max float64) float64 {
	n, err := randomUint64()
	if err != nil {
		return min
	}
	return min + (float64(n>>11) / (1 << 53) * (max - min))
}

func cryptoRandInt(minInclusive, maxExclusive int) (int, error) {
	if maxExclusive <= minInclusive {
		return minInclusive, nil
	}
	n, err := randomUint64()
	if err != nil {
		return 0, err
	}
	return minInclusive + int(n%uint64(maxExclusive-minInclusive)), nil
}

func randomUint64() (uint64, error) {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(buf[:]), nil
}

func randomBytes(length int) ([]byte, error) {
	out := make([]byte, length)
	_, err := rand.Read(out)
	return out, err
}
