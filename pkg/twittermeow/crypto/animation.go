package crypto

import (
	"encoding/base64"
	"math"
	"strconv"
	"strings"
)

func GenerateAnimationState(variableIndexes *[4]int, loadingAnims *[4][16][11]int, verificationToken string) string {
	verificationTokenBytes, err := base64.StdEncoding.DecodeString(verificationToken)
	if err != nil {
		return ""
	}
	svgData := loadingAnims[verificationTokenBytes[5]%4][verificationTokenBytes[variableIndexes[0]]%16]
	animationTime := int(verificationTokenBytes[variableIndexes[1]]%16) * int(verificationTokenBytes[variableIndexes[2]]%16) * int(verificationTokenBytes[variableIndexes[3]]%16)
	randomPart := generateAnimationStateWithParams(svgData[:], animationTime)
	return randomPart
}

const totalTime = 4096.0

func generateAnimationStateWithParams(row []int, animTime int) string {
	if animTime >= totalTime-1 {
		return "000"
	}
	fromColor := []float64{float64(row[0]), float64(row[1]), float64(row[2]), 1.0}
	toColor := []float64{float64(row[3]), float64(row[4]), float64(row[5]), 1.0}
	fromRotation := []float64{0.0}
	toRotation := []float64{math.Floor(mapValueToRange(float64(row[6]), 60.0, 360.0))}
	row = row[7:]
	curves := [4]float64{}
	for i := 0; i < len(row); i++ {
		curves[i] = toFixed(mapValueToRange(float64(row[i]), isEven(i), 1.0), 2)
	}
	c := &cubic{Curves: curves}
	val := c.getValue(math.Round(float64(animTime)/10.0) * 10.0 / totalTime)
	color := interpolate(fromColor, toColor, val)
	rotation := interpolate(fromRotation, toRotation, val)
	matrix := convertRotationToMatrix(rotation[0])
	strArr := []string{}
	for i := 0; i < len(color)-1; i++ {
		if color[i] < 0 {
			color[i] = 0
		}
		roundedColor := math.Round(color[i])
		if roundedColor < 0 {
			roundedColor = 0
		}
		hexColor := strconv.FormatInt(int64(roundedColor), 16)
		strArr = append(strArr, hexColor)
	}
	for i := 0; i < len(matrix)-2; i++ {
		rounded := toFixed(matrix[i], 2)
		if rounded < 0 {
			rounded = -rounded
		}
		strArr = append(strArr, floatToHex(rounded))
	}
	strArr = append(strArr, "0", "0")
	return strings.Join(strArr, "")
}

func interpolate(from, to []float64, f float64) []float64 {
	out := []float64{}
	for i := 0; i < len(from); i++ {
		out = append(out, interpolateNum(from[i], to[i], f))
	}
	return out
}

func interpolateNum(from, to, f float64) float64 {
	return from*(1.0-f) + to*f
}

func convertRotationToMatrix(degrees float64) []float64 {
	radians := degrees * math.Pi / 180
	c := math.Cos(radians)
	s := math.Sin(radians)
	return []float64{c, s, -s, c, 0, 0}
}

type cubic struct {
	Curves [4]float64
}

func (c *cubic) getValue(time float64) float64 {
	startGradient := 0.0
	endGradient := 0.0
	if time <= 0.0 {
		if c.Curves[0] > 0.0 {
			startGradient = c.Curves[1] / c.Curves[0]
		} else if c.Curves[1] == 0.0 && c.Curves[2] > 0.0 {
			startGradient = c.Curves[3] / c.Curves[2]
		}
		return startGradient * time
	}

	if time >= 1.0 {
		if c.Curves[2] < 1.0 {
			endGradient = (c.Curves[3] - 1.0) / (c.Curves[2] - 1.0)
		} else if c.Curves[2] == 1.0 && c.Curves[0] < 1.0 {
			endGradient = (c.Curves[1] - 1.0) / (c.Curves[0] - 1.0)
		}
		return 1.0 + endGradient*(time-1.0)
	}

	start := 0.0
	end := 1.0
	mid := 0.0
	for start < end {
		mid = (start + end) / 2
		xEst := bezierCurve(c.Curves[0], c.Curves[2], mid)
		if math.Abs(time-xEst) < 0.00001 {
			return bezierCurve(c.Curves[1], c.Curves[3], mid)
		}
		if xEst < time {
			start = mid
		} else {
			end = mid
		}
	}
	return bezierCurve(c.Curves[1], c.Curves[3], mid)
}

func bezierCurve(a, b, m float64) float64 {
	return 3.0*a*(1-m)*(1-m)*m + 3.0*b*(1-m)*m*m + m*m*m
}

func roundPositive(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(roundPositive(num*output)) / output
}

func floatToHex(x float64) string {
	quotient := int(x)
	fraction := x - float64(quotient)
	result := strconv.FormatInt(int64(quotient), 16)
	if fraction == 0 {
		return result
	}

	for fraction > 0 {
		fraction *= 16
		integer := int64(fraction)
		result += strconv.FormatInt(integer, 16)
		fraction -= float64(integer)
	}

	return result
}

func mapValueToRange(val, min, max float64) float64 {
	return val*(max-min)/255.0 + min
}

func isEven(val int) float64 {
	if val%2 == 1 {
		return -1.0
	}
	return 0.0
}
