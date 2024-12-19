package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	neturl "net/url"
	"strconv"
	"time"

	"go.mau.fi/util/random"
)

func SignTransaction(animationToken, verificationToken, url, method string) (string, error) {
	verificationTokenBytes, err := base64.StdEncoding.DecodeString(verificationToken)
	if err != nil {
		return "", fmt.Errorf("failed to decode verification token: %w", err)
	}
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL %q: %w", url, err)
	}

	ts, tsBytes := makeTSBytes()
	salt := base64.RawStdEncoding.EncodeToString([]byte{161, 183, 226, 163, 7, 171, 122, 24, 171, 138, 120})
	hashInput := fmt.Sprintf("%s!%s!%s%s%s", method, parsedURL.Path, ts, salt, animationToken)
	rawHash := sha256.Sum256([]byte(hashInput))
	sdp := generateSDP()
	sdpHash := []byte{sdp[verificationTokenBytes[5]%8], sdp[verificationTokenBytes[8]%8]}
	hash := append(rawHash[:], sdpHash...)

	resultBytes := bytes.NewBuffer(make([]byte, 0, 1+len(verificationTokenBytes)+len(tsBytes)+16+1))
	resultBytes.WriteByte(byte(rand.IntN(256)))
	resultBytes.Write(verificationTokenBytes)
	resultBytes.Write(tsBytes)
	resultBytes.Write(hash[:16])
	resultBytes.WriteByte(3)

	return base64.RawStdEncoding.EncodeToString(encodeXOR(resultBytes.Bytes())), nil
}

const sdpTemplate = "v=0\r\no=- %d 2 IN IP4 127.0.0.1\r\ns=-\r\nt=0 0\r\na=group:BUNDLE 0\r\na=extmap-allow-mixed\r\n" +
	"a=msid-semantic: WMS\r\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\r\nc=IN IP4 0.0.0.0\r\na=ice-ufrag:%s\r\n" +
	"a=ice-pwd:%s\r\na=ice-options:trickle\r\na=fingerprint:sha-256 00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:" +
	"00:00:00:00:00:00:00:00:00:00:00:00:00\r\na=setup:actpass\r\na=mid:0\r\na=sctp-port:5000\r\na=max-message-size:262144\r\n"

func generateSDP() string {
	return fmt.Sprintf(sdpTemplate, rand.Int64(), random.String(4), random.String(24))
}

func makeTSBytes() (string, []byte) {
	ts := time.Now().Unix() - 1682924400
	resultInt32 := make([]byte, 4)
	binary.LittleEndian.PutUint32(resultInt32, uint32(ts))
	return strconv.FormatInt(ts, 10), resultInt32
}

func encodeXOR(plainArr []byte) []byte {
	encodedArr := make([]byte, len(plainArr))
	for i := 0; i < len(plainArr); i++ {
		if i == 0 {
			encodedArr[i] = plainArr[i]
		} else {
			encodedArr[i] = plainArr[i] ^ plainArr[0]
		}
	}
	return encodedArr
}
