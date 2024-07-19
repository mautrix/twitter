package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	neturl "net/url"
	"strconv"
	"strings"
	"time"
)

type TransactionData struct {
	VerificationToken string
}

const (
	timestampConstant = 1682924400000
	// constant color/easing string (i dont think its relevant)
	hexStr     = "74c5e50fd70a3d70a3d701eb851eb851eb801eb851eb851eb80fd70a3d70a3d700"
	twitterStr = "bird"
)

// it may not be this simple.
// yes this does generate/sign the x-client-transaction header but there might be a lot more to it.
// this doesn't properly grab the svg data embedded in the page, it generates brand new svg data
// there's also some RTC connection involved under the hood which is hard to interpret but I have not investigated it
// the RTC connection being established does send/receive data related to the transaction being signed so it might have something to do with verification

func SignTransaction(verificationToken, url, method string) (string, error) {
	verificationTokenBytes, err := base64.StdEncoding.DecodeString(verificationToken)
	if err != nil {
		return "", fmt.Errorf("failed to decode verification token while signing client transaction id: %e", err)
	}

	parsedUrl, err := neturl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("failed to parse path for request (url=%s, method=%s): %e", url, method, err)
	}

	ts, tsBytes := signTimestamp()
	unsignedString := fmt.Sprintf("%s!%s!%s%s%s", method, parsedUrl.Path, ts, twitterStr, hexStr)

	hash := sha256.Sum256([]byte(unsignedString))
	hashSlice := hash[:16]

	resultBytes := bytes.NewBuffer([]byte{})
	resultBytes.WriteByte(byte(randNum()))
	resultBytes.Write(verificationTokenBytes)
	resultBytes.Write(tsBytes)
	resultBytes.Write(hashSlice)
	resultBytes.WriteByte(1)

	signedBytes := encodeXOR(resultBytes.Bytes())
	return base64.RawStdEncoding.EncodeToString(signedBytes), nil
}

func randNum() int {
	return rand.Intn(256)
}

func signTimestamp() (string, []byte) {
	elapsed := time.Now().UnixNano()/int64(time.Millisecond) - timestampConstant
	result := int(math.Floor(float64(elapsed) / 1000.0))
	resultInt32 := make([]byte, 4)

	binary.LittleEndian.PutUint32(resultInt32, uint32(result))

	return strconv.Itoa(result), resultInt32
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

//lint:ignore U1000 TODO fix unused method
func decodeSVGStr(svgStr string) [][]int {
	segmentStrings := strings.Split(strings.TrimSpace(svgStr), "C")[1:]
	byteArrays := [][]int{}

	for _, segment := range segmentStrings {
		segment = strings.ReplaceAll(strings.ReplaceAll(segment, "h", ""), "s", "")
		coords := strings.Fields(segment)
		bytesList := []int{}

		for _, coord := range coords {
			if strings.Contains(coord, ",") {
				for _, s := range strings.Split(coord, ",") {
					num, err := strconv.Atoi(s)
					if err == nil {
						bytesList = append(bytesList, num)
					}
				}
			} else {
				num, err := strconv.Atoi(coord)
				if err == nil {
					bytesList = append(bytesList, num)
				}
			}
		}
		byteArrays = append(byteArrays, bytesList)
	}

	return byteArrays
}

//lint:ignore U1000 TODO fix unused method
func buildColorStr(arr []int) []string {
	colors := []string{}
	s := ""
	for i := 0; i < 3; i++ {
		if arr[i] < 16 {
			s += "0"
		}
		s += strconv.FormatInt(int64(arr[i]), 16)
	}
	colors = append(colors, s)
	s = ""
	for i := 3; i < 6; i++ {
		if arr[i] < 16 {
			s += "0"
		}
		s += strconv.FormatInt(int64(arr[i]), 16)
	}
	colors = append(colors, s)
	return colors
}

//lint:ignore U1000 TODO fix unused method
func buildEasingStr(arr []int) []float64 {
	nums := []float64{}
	t := 1.0

	for i := 0; i < len(arr); i++ {
		b := float64(arr[i])
		o := 0.0
		if i%2 != 0 {
			o = -1.0
		}
		val := (((t - o) * b) / 255.0) + o
		roundedVal, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", val), 64)
		nums = append(nums, roundedVal)
	}

	return nums
}
