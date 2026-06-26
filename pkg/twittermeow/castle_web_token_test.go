package twittermeow

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode/utf16"
)

func TestCurrentCastleWrappedTokenRoundTrip(t *testing.T) {
	const payload = `[null,"abc",123,true,{"k":"v"}]`
	const timestampMillis = int64(1760000000123)

	token, err := createCurrentCastleWrappedToken(payload, timestampMillis)
	if err != nil {
		t.Fatalf("createCurrentCastleWrappedToken() error = %v", err)
	}
	if !strings.HasPrefix(token, currentCastleTokenPrefix) {
		t.Fatalf("token missing current Castle prefix: %q", token[:min(len(token), len(currentCastleTokenPrefix))])
	}
	wire, err := decodeCurrentCastleWrappedTokenForTest(token)
	if err != nil {
		t.Fatalf("decodeCurrentCastleWrappedTokenForTest() error = %v", err)
	}
	if len(wire) < 2 {
		t.Fatalf("decoded wire too short: %d", len(wire))
	}
	encodedTimestampLen := int(wire[0])
	if len(wire) <= 1+encodedTimestampLen {
		t.Fatalf("decoded wire missing payload: wire_len=%d encoded_timestamp_len=%d", len(wire), encodedTimestampLen)
	}
	timestamp := strconv.FormatInt(timestampMillis, 10)
	expectedTimestamp := encodeCurrentCastleTimestamp(timestamp)
	if got := wire[1 : 1+encodedTimestampLen]; got != expectedTimestamp {
		t.Fatalf("encoded timestamp = %q, want %q", got, expectedTimestamp)
	}
	gotPayload := xorCurrentCastleString(wire[1+encodedTimestampLen:], timestamp)
	if wantPayload := insertCurrentCastleChecksum(payload); gotPayload != wantPayload {
		t.Fatalf("decoded payload = %q, want %q", gotPayload, wantPayload)
	}
	if !strings.Contains(gotPayload, currentCastleChecksumLabel) {
		t.Fatalf("decoded payload missing checksum label %q: %q", currentCastleChecksumLabel, gotPayload)
	}
}

func TestCurrentCastleKnownTransforms(t *testing.T) {
	const payload = `[1,2,3]`
	const timestamp = `1760000000123`

	if got, want := currentCastleChecksum(payload), "9cdea041"; got != want {
		t.Fatalf("currentCastleChecksum() = %q, want %q", got, want)
	}
	tqCases := map[string]string{
		"":           "969740e5",
		"abc":        "4a84ecf2",
		"119280":     "f394e296",
		"3695375093": "a3e047f9",
	}
	for input, want := range tqCases {
		if got := currentCastleTQHash(input); got != want {
			t.Fatalf("currentCastleTQHash(%q) = %q, want %q", input, got, want)
		}
	}
	tiCases := map[string]string{
		"":           "83b623b0",
		"abc":        "2debeae9",
		`"GraL"`:     "f0b537b7",
		`"override"`: "37d6d39f",
	}
	for input, want := range tiCases {
		if got := currentCastleTIHash(input); got != want {
			t.Fatalf("currentCastleTIHash(%q) = %q, want %q", input, got, want)
		}
	}
	if got, want := encodeCurrentCastleTimestamp(timestamp), "WnP35wy50MvQy9DL0MvQy9DL0Mtac+Yab8I="; got != want {
		t.Fatalf("encodeCurrentCastleTimestamp() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleUnits([]uint16{0x002f, 0x8000, 0xffff}), "RyQxU6ar"; got != want {
		t.Fatalf("encodeCurrentCastleUnits() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleRawPackedUnits(jsUTF16Units("b537b7")), "S+hTpr/W8DdL6PA3"; got != want {
		t.Fatalf("encodeCurrentCastleRawPackedUnits() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleLR("abc"), "oGAg"; got != want {
		t.Fatalf("encodeCurrentCastleLR() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleFloat64(0), "xsbGxsbGxsY="; got != want {
		t.Fatalf("encodeCurrentCastleFloat64(0) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleFloat64(1), "OTbGxsbGxsY="; got != want {
		t.Fatalf("encodeCurrentCastleFloat64(1) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleFloat64(123.456), "hhjbKVyZ+PE="; got != want {
		t.Fatalf("encodeCurrentCastleFloat64(123.456) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6TU("abc"), "QXe5"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6TU() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6TF("abc"), "hxTl"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6TF() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6Component0("abc"), "cLDw"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6Component0() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6Component1("abc"), "Sjsq"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6Component1() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6Component1("X Castle"), "vB8NSisamgo="; got != want {
		t.Fatalf("encodeCurrentCastleSlot6Component1() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6Component2("abc"), "SXTz"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6Component2() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6Component5Fallback("abc"), "aGlm"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6Component5Fallback() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6UZ("abc"), "Y29r"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6UZ() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6UZ("X Castle"), "hybrYys3V3M="; got != want {
		t.Fatalf("encodeCurrentCastleSlot6UZ() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6IH("abc"), "odoT"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6IH() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6IH("X Castle"), "oCjzoaPcFIU="; got != want {
		t.Fatalf("encodeCurrentCastleSlot6IH() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6FJ("abc"), "Eh4a"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6FJ() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6FJ("X Castle"), "9heaElpGJgI="; got != want {
		t.Fatalf("encodeCurrentCastleSlot6FJ() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6IB("abc"), "BnWg"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6IB() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6IB("X Castle"), "cwsABpDPd8o="; got != want {
		t.Fatalf("encodeCurrentCastleSlot6IB() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6AH("abc"), "1LSU"; got != want {
		t.Fatalf("encodeCurrentCastleSlot6AH() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6AH("X Castle"), "tfyY1BLys1Q="; got != want {
		t.Fatalf("encodeCurrentCastleSlot6AH() = %q, want %q", got, want)
	}
	indexR7Cases := map[uint32]uint32{
		0x03718000: 264,
		0x0372e000: 275,
		0x03730000: 276,
		0x0373c000: 282,
		0x0376c000: 306,
		0x0377e000: 315,
		0x03788000: 320,
		0x037a2000: 333,
		0x03816000: 391,
		0x0383a000: 409,
		0x0385e000: 427,
		0x0387c000: 442,
		0x03880000: 444,
		0x03894000: 454,
		0x038a0000: 460,
		0x038b2000: 469,
		0x038d4000: 486,
		0x038d6000: 487,
		0x038d8000: 488,
		0x038da000: 489,
		0x038e2000: 493,
	}
	for input, want := range indexR7Cases {
		if got := currentCastleIndexR7(input); got != want {
			t.Fatalf("currentCastleIndexR7(%#x) = %d, want %d", input, got, want)
		}
	}
	indexTGCases := map[uint32]uint32{
		38521: 271,
		46027: 325,
		47278: 334,
		51170: 362,
		51865: 367,
		52560: 372,
		52699: 373,
		53811: 381,
		57842: 410,
		58120: 412,
		58259: 413,
		59232: 420,
		61456: 436,
		61873: 439,
		62151: 441,
		63402: 450,
		65348: 464,
		66460: 472,
		67155: 477,
		67989: 483,
		68128: 484,
	}
	for input, want := range indexTGCases {
		if got := currentCastleIndexTG(input); got != want {
			t.Fatalf("currentCastleIndexTG(%d) = %d, want %d", input, got, want)
		}
	}
	tbCases := map[uint32]uint32{
		0:          852,
		1:          991,
		0x12345678: 3798660732,
		0xffffffff: 713,
	}
	for input, want := range tbCases {
		got := currentCastleTB(input)
		if got != want {
			t.Fatalf("currentCastleTB(%#x) = %d, want %d", input, got, want)
		}
		if roundTrip := invertCurrentCastleTB(got); roundTrip != input {
			t.Fatalf("invertCurrentCastleTB(currentCastleTB(%#x)) = %#x, want %#x", input, roundTrip, input)
		}
	}
	if got, want := invertCurrentCastleTB(2115052381), uint32(0x1353003b); got != want {
		t.Fatalf("invertCurrentCastleTB(accepted slot 2) = %#x, want %#x", got, want)
	}
	highSlots := currentCastleHighTimingSlotIndexes()
	if len(highSlots) != 63 {
		t.Fatalf("currentCastleHighTimingSlotIndexes() len = %d, want 63", len(highSlots))
	}
	sortedHighSlots := append([]int(nil), highSlots...)
	sort.Ints(sortedHighSlots)
	for i, got := range sortedHighSlots {
		if want := 431 + i; got != want {
			t.Fatalf("currentCastleHighTimingSlotIndexes()[sorted %d] = %d, want %d", i, got, want)
		}
	}
	if got, want := insertCurrentCastleChecksum(payload), `[1,2,Nqkju9cdea0413]`; got != want {
		t.Fatalf("insertCurrentCastleChecksum() = %q, want %q", got, want)
	}
}

func TestCurrentCastleLowerFloatTransforms(t *testing.T) {
	cases := []struct {
		name string
		fn   func(float64) string
		zero string
		one  string
		some string
	}{
		{name: "slot1", fn: encodeCurrentCastleLowerFloatSlot1, zero: "4ODg4ODg4OA=", one: "KWfg4ODg4OA=", some: "4jI+qRAsNWs="},
		{name: "slot3", fn: encodeCurrentCastleLowerFloatSlot3, zero: "bm5ubm5ubm4=", one: "G39ubm5ubm4=", some: "agkBHM1lE+c="},
		{name: "slot4", fn: encodeCurrentCastleLowerFloatSlot4, zero: "7Ozs7Ozs7Ow=", one: "0xzs7Ozs7Ow=", some: "rLIxw/ZzUps="},
		{name: "slot5", fn: encodeCurrentCastleLowerFloatSlot5, zero: "YWFhYWFhYWE=", one: "II1hYWFhYWE=", some: "cfjYHOdIAC4="},
		{name: "slot9", fn: encodeCurrentCastleLowerFloatSlot9, zero: "FRUVFRUVFRU=", one: "YMUVFRUVFRU=", some: "VQuWEPdA6/g="},
		{name: "slot10", fn: encodeCurrentCastleLowerFloatSlot10, zero: "m5ubm5ubm5s=", one: "6Yebm5ubm5s=", some: "qxVp5RTxPds="},
		{name: "slot12", fn: encodeCurrentCastleLowerFloatSlot12, zero: "AAAAAAAAAAA=", one: "LWoAAAAAAAA=", some: "2AmpN7GRzQo="},
		{name: "slot13", fn: encodeCurrentCastleLowerFloatSlot13, zero: "7Ozs7Ozs7Ow=", one: "Jwvs7Ozs7Ow=", some: "iAGpDqUtz0M="},
		{name: "slot14", fn: encodeCurrentCastleLowerFloatSlot14, zero: "/f39/f39/f0=", one: "uo39/f39/f0=", some: "vZMISrdaM0I="},
		{name: "slot15", fn: encodeCurrentCastleLowerFloatSlot15, zero: "Hh4eHh4eHh4=", one: "YP8eHh4eHh4=", some: "nqKlQCohY/A="},
		{name: "slot16", fn: encodeCurrentCastleLowerFloatSlot16, zero: "PDw8PDw8PDw=", one: "llQ8PDw8PDw=", some: "HFutzi2ma8I="},
		{name: "slot17", fn: encodeCurrentCastleLowerFloatSlot17, zero: "6+vr6+vr6+s=", one: "/Pvr6+vr6+s=", some: "qy3+DLGczUQ="},
		{name: "slot18", fn: encodeCurrentCastleLowerFloatSlot18, zero: "AAAAAAAAAAA=", one: "fuEAAAAAAAA=", some: "gLy7XjQ/fe4="},
		{name: "slot19", fn: encodeCurrentCastleLowerFloatSlot19, zero: "AAAAAAAAAAA=", one: "4ZQAAAAAAAA=", some: "MOq7dR/p8ls="},
		{name: "slot21", fn: encodeCurrentCastleLowerFloatSlot21, zero: "1dXV1dXV1dU=", one: "hNDV1dXV1dU=", some: "2XqigzaOfAk="},
		{name: "slot22", fn: encodeCurrentCastleLowerFloatSlot22, zero: "GBgYGBgYGBg=", one: "/wYYGBgYGBg=", some: "ENOj/Vvrz/Y="},
		{name: "slot23", fn: encodeCurrentCastleLowerFloatSlot23, zero: "6Ojo6Ojo6Og=", one: "Z8ro6Ojo6Og=", some: "aaWkRx0gZtc="},
		{name: "slot24", fn: encodeCurrentCastleLowerFloatSlot24, zero: "8/Pz8/Pz8/M=", one: "L3bz8/Pz8/M=", some: "9WqgskkyZfI="},
		{name: "slot29", fn: encodeCurrentCastleLowerFloatSlot29, zero: "MzMzMzMzMzM=", one: "/A8zMzMzMzM=", some: "I6RE+LXUnO4="},
		{name: "slot30", fn: encodeCurrentCastleLowerFloatSlot30, zero: "rKysrKysrKw=", one: "S7KsrKysrKw=", some: "pGcXSe9fe0I="},
		{name: "slot32", fn: encodeCurrentCastleLowerFloatSlot32, zero: "2dnZ2dnZ2dk=", one: "xFLZ2dnZ2dk=", some: "28vPRGndyJ4="},
		{name: "slot33", fn: encodeCurrentCastleLowerFloatSlot33, zero: "MDAwMDAwMDA=", one: "/2wwMDAwMDA=", some: "QMen+7YX3w0="},
		{name: "slot34", fn: encodeCurrentCastleLowerFloatSlot34, zero: "Tk5OTk5OTk4=", one: "vXJOTk5OTk4=", some: "XtEtodCF+a8="},
		{name: "slot35", fn: encodeCurrentCastleLowerFloatSlot35, zero: "AwMDAwMDAwM=", one: "5jMDAwMDAwM=", some: "QylMFnXGCT4="},
		{name: "slot36", fn: encodeCurrentCastleLowerFloatSlot36, zero: "WFhYWFhYWFg=", one: "OfdYWFhYWFg=", some: "qhp2uai0lUs="},
		{name: "slot37", fn: encodeCurrentCastleLowerFloatSlot37, zero: "Li4uLi4uLi4=", one: "OV4uLi4uLi4=", some: "bhArqfTZMDE="},
		{name: "slot38", fn: encodeCurrentCastleLowerFloatSlot38, zero: "MzMzMzMzMzM=", one: "LLozMzMzMzM=", some: "NSUhrAMvKO4="},
		{name: "slot39", fn: encodeCurrentCastleLowerFloatSlot39, zero: "AAAAAAAAAAA=", one: "mqEAAAAAAAA=", some: "gXSOOl3atOo="},
		{name: "slot43", fn: encodeCurrentCastleLowerFloatSlot43, zero: "q6urq6urq6s=", one: "RJCrq6urq6s=", some: "/2KKcZ7G7Gg="},
		{name: "slot45", fn: encodeCurrentCastleLowerFloatSlot45, zero: "q6urq6urq6s=", one: "2turq6urq6s=", some: "a21cCml6DRI="},
		{name: "slot48", fn: encodeCurrentCastleLowerFloatSlot48, zero: "AAAAAAAAAAA=", one: "LYIAAAAAAAA=", some: "Bld6r/UuUGw="},
		{name: "slot50", fn: encodeCurrentCastleLowerFloatSlot50, zero: "WVlZWVlZWVk=", one: "IslZWVlZWVk=", some: "mQtckjfCK6o="},
		{name: "slot51", fn: encodeCurrentCastleLowerFloatSlot51, zero: "GRkZGRkZGRk=", one: "9J4ZGRkZGRk=", some: "G8uhcinzyrU="},
		{name: "slot52", fn: encodeCurrentCastleLowerFloatSlot52, zero: "AAAAAAAAAAA=", one: "RVAAAAAAAAA=", some: "wCovlX5lSi0="},
		{name: "slot54", fn: encodeCurrentCastleLowerFloatSlot54, zero: "xcXFxcXFxcU=", one: "mMbFxcXFxcU=", some: "wal4mWGSoxw="},
		{name: "slot55", fn: encodeCurrentCastleLowerFloatSlot55, zero: "AAAAAAAAAAA=", one: "3qQAAAAAAAA=", some: "hZuQfR4X06M="},
		{name: "slot57", fn: encodeCurrentCastleLowerFloatSlot57, zero: "sLCwsLCwsLA=", one: "jVCwsLCwsLA=", some: "0wceiszNpN0="},
		{name: "slot58", fn: encodeCurrentCastleLowerFloatSlot58, zero: "XFxcXFxcXFw=", one: "UNxcXFxcXFw=", some: "Wk1H0G9VShY="},
		{name: "slot59", fn: encodeCurrentCastleLowerFloatSlot59, zero: "6urq6urq6uo=", one: "BGrq6urq6uo=", some: "7PwRgxsH/0Y="},
		{name: "slot61", fn: encodeCurrentCastleLowerFloatSlot61, zero: "UFBQUFBQUFA=", one: "6VlQUFBQUFA=", some: "XHUL4r/jf2I="},
		{name: "slot62", fn: encodeCurrentCastleLowerFloatSlot62, zero: "3d3d3d3d3d0=", one: "yI3d3d3d3d0=", some: "HVdKuNNod4A="},
		{name: "slot63", fn: encodeCurrentCastleLowerFloatSlot63, zero: "a2tra2tra2s=", one: "cetra2tra2s=", some: "aX6C8Zh0fbc="},
		{name: "slot65", fn: encodeCurrentCastleLowerFloatSlot65, zero: "g4ODg4ODg4M=", one: "gLODg4ODg4M=", some: "w+UqsLFgxdg="},
		{name: "slot66", fn: encodeCurrentCastleLowerFloatSlot66, zero: "np6enp6enp4=", one: "da6enp6enp4=", some: "3jBfxVxVEN0="},
		{name: "slot69", fn: encodeCurrentCastleLowerFloatSlot69, zero: "gICAgICAgIA=", one: "xqSAgICAgIA=", some: "sBVL6j7uPeg="},
		{name: "slot70", fn: encodeCurrentCastleLowerFloatSlot70, zero: "ycnJycnJyck=", one: "dlHJycnJyck=", some: "6eYnfsQGthI="},
		{name: "slot74", fn: encodeCurrentCastleLowerFloatSlot74, zero: "n5+fn5+fn58=", one: "Du+fn5+fn58=", some: "X53MHuHufVY="},
		{name: "slot76", fn: encodeCurrentCastleLowerFloatSlot76, zero: "jY2NjY2NjY0=", one: "IEWNjY2NjY0=", some: "bXCvGJLQYAQ="},
		{name: "slot80", fn: encodeCurrentCastleLowerFloatSlot80, zero: "8PDw8PDw8PA=", one: "IzDw8PDw8PA=", some: "39Y8Y/eBRTI="},
		{name: "slot81", fn: encodeCurrentCastleLowerFloatSlot81, zero: "OTk5OTk5OTk=", one: "6CU5OTk5OTk=", some: "KaJE1KvQitY="},
		{name: "slot82", fn: encodeCurrentCastleLowerFloatSlot82, zero: "c3Nzc3Nzc3M=", one: "Tplzc3Nzc3M=", some: "q7pqtAJy/qk="},
		{name: "slot84", fn: encodeCurrentCastleLowerFloatSlot84, zero: "+/v7+/v7+/s=", one: "t6L7+/v7+/s=", some: "8D4UQjEHjYw="},
		{name: "slot85", fn: encodeCurrentCastleLowerFloatSlot85, zero: "a2tra2tra2s=", one: "BJ5ra2tra2s=", some: "SdmlhBu30KI="},
		{name: "slot88", fn: encodeCurrentCastleLowerFloatSlot88, zero: "wsLCwsLCwsI=", one: "Q6rCwsLCwsI=", some: "4tMUS7Xzg2c="},
		{name: "slot91", fn: encodeCurrentCastleLowerFloatSlot91, zero: "FRUVFRUVFRU=", one: "ZBoVFRUVFRU=", some: "Eb0Jaf5mv+o="},
		{name: "slot93", fn: encodeCurrentCastleLowerFloatSlot93, zero: "XV1dXV1dXV0=", one: "I7xdXV1dXV0=", some: "3eHmA2liILM="},
		{name: "slot94", fn: encodeCurrentCastleLowerFloatSlot94, zero: "V1dXV1dXV1c=", one: "mEdXV1dXV1c=", some: "l3nyiDW4GdA="},
		{name: "slot95", fn: encodeCurrentCastleLowerFloatSlot95, zero: "GxsbGxsbGxs=", one: "9usbGxsbGxs=", some: "202ohrFWLf4="},
		{name: "slot96", fn: encodeCurrentCastleLowerFloatSlot96, zero: "3Nzc3Nzc3Nw=", one: "3iPc3Nzc3Nw=", some: "XJi5/hC9/04="},
		{name: "slot97", fn: encodeCurrentCastleLowerFloatSlot97, zero: "goKCgoKCgoI=", one: "gfKCgoKCgoI=", some: "wmg3MdQhiPk="},
		{name: "slot101", fn: encodeCurrentCastleLowerFloatSlot101, zero: "uLi4uLi4uLg=", one: "Xhm4uLi4uLg=", some: "OcTqvtwehA4="},
		{name: "slot105", fn: encodeCurrentCastleLowerFloatSlot105, zero: "AAAAAAAAAAA=", one: "1iwAAAAAAAA=", some: "EIVswpDOvdw="},
		{name: "slot107", fn: encodeCurrentCastleLowerFloatSlot107, zero: "LCwsLCwsLCw=", one: "21wsLCwsLCw=", some: "bOIpS0Z7AsM="},
		{name: "slot110", fn: encodeCurrentCastleLowerFloatSlot110, zero: "2tra2tra2to=", one: "1Yra2tra2to=", some: "GjBrhdy1EO0="},
		{name: "slot111", fn: encodeCurrentCastleLowerFloatSlot111, zero: "g4ODg4ODg4M=", one: "9sKDg4ODg4M=", some: "guhfNrt3aZU="},
		{name: "slot112", fn: encodeCurrentCastleLowerFloatSlot112, zero: "AAAAAAAAAAA=", one: "l3AAAAAAAAA=", some: "QA7lByr3bo8="},
		{name: "slot114", fn: encodeCurrentCastleLowerFloatSlot114, zero: "aWlpaWlpaWk=", one: "hWVpaWlpaWk=", some: "edgHkf2NwIs="},
		{name: "slot117", fn: encodeCurrentCastleLowerFloatSlot117, zero: "AwMDAwMDAwM=", one: "vAgDAwMDAwM=", some: "B2YauSO+YDI="},
		{name: "slot118", fn: encodeCurrentCastleLowerFloatSlot118, zero: "fX19fX19fX0=", one: "Ak19fX19fX0=", some: "PSOgEmfig8o="},
		{name: "slot122", fn: encodeCurrentCastleLowerFloatSlot122, zero: "Q0NDQ0NDQ0M=", one: "K0FDQ0NDQ0M=", some: "Sw/+KYY3GzI="},
		{name: "slot124", fn: encodeCurrentCastleLowerFloatSlot124, zero: "CAgICAgICAg=", one: "lckICAgICAg=", some: "tAObZoe/LdE="},
		{name: "slot125", fn: encodeCurrentCastleLowerFloatSlot125, zero: "KioqKioqKio=", one: "qAsqKioqKio=", some: "qu7tiGZprzg="},
		{name: "slot126", fn: encodeCurrentCastleLowerFloatSlot126, zero: "ra2tra2tra0=", one: "9B2tra2tra0=", some: "bbuS5NdUG7w="},
		{name: "slot129", fn: encodeCurrentCastleLowerFloatSlot129, zero: "FxcXFxcXFxc=", one: "dgcXFxcXFxc=", some: "V7VUZvHWFa4="},
		{name: "slot131", fn: encodeCurrentCastleLowerFloatSlot131, zero: "PDw8PDw8PDw=", one: "8wA8PDw8PDw=", some: "LKtL97rbk+E="},
		{name: "slot132", fn: encodeCurrentCastleLowerFloatSlot132, zero: "W1tbW1tbW1s=", one: "QstbW1tbW1s=", some: "mw108qHirfo="},
		{name: "slot134", fn: encodeCurrentCastleLowerFloatSlot134, zero: "xcXFxcXFxcU=", one: "5NXFxcXFxcU=", some: "BaeG9OMExyw="},
		{name: "slot136", fn: encodeCurrentCastleLowerFloatSlot136, zero: "hYWFhYWFhYU=", one: "TYOFhYWFhYU=", some: "bTEha8hZPVQ="},
		{name: "slot138", fn: encodeCurrentCastleLowerFloatSlot138, zero: "Y2NjY2NjY2M=", one: "9BNjY2NjY2M=", some: "IzkKpDUU2dw="},
		{name: "slot139", fn: encodeCurrentCastleLowerFloatSlot139, zero: "j4+Pj4+Pj48=", one: "KP+Pj4+Pj48=", some: "TyH6mEUIASA="},
		{name: "slot140", fn: encodeCurrentCastleLowerFloatSlot140, zero: "eHh4eHh4eHg=", one: "rch4eHh4eHg=", some: "uIKHfZZNIiU="},
		{name: "slot141", fn: encodeCurrentCastleLowerFloatSlot141, zero: "UFBQUFBQUFA=", one: "TThQUFBQUFA=", some: "MJ0aVbf9DSE="},
		{name: "slot142", fn: encodeCurrentCastleLowerFloatSlot142, zero: "rq6urq6urq4=", one: "LY6urq6urq4=", some: "L2toDeLtKp0="},
		{name: "slot143", fn: encodeCurrentCastleLowerFloatSlot143, zero: "SUlJSUlJSUk=", one: "/uNJSUlJSUk=", some: "ISp6+NLiDlM="},
		{name: "slot144", fn: encodeCurrentCastleLowerFloatSlot144, zero: "U1NTU1NTU1M=", one: "6HtTU1NTU1M=", some: "sxgNkCL4KBw="},
		{name: "slot145", fn: encodeCurrentCastleLowerFloatSlot145, zero: "p6enp6enp6c=", one: "Rp+np6enp6c=", some: "x9YVPrR2BmI="},
		{name: "slot146", fn: encodeCurrentCastleLowerFloatSlot146, zero: "Dw8PDw8PDw8=", one: "vh8PDw8PDw8=", some: "T418ToneLXY="},
		{name: "slot148", fn: encodeCurrentCastleLowerFloatSlot148, zero: "2NjY2NjY2Ng=", one: "pt/Y2NjY2Ng=", some: "WLSpBvyNm/Y="},
		{name: "slot151", fn: encodeCurrentCastleLowerFloatSlot151, zero: "6urq6urq6uo=", one: "V/Dq6urq6uo=", some: "8qUOXSdLuUg="},
		{name: "slot153", fn: encodeCurrentCastleLowerFloatSlot153, zero: "eHh4eHh4eHg=", one: "C1l4eHh4eHg=", some: "+V1u6hXKHPo="},
		{name: "slot155", fn: encodeCurrentCastleLowerFloatSlot155, zero: "UlJSUlJSUlI=", one: "mNZSUlJSUlI=", some: "VMQNHeGdwV0="},
		{name: "slot157", fn: encodeCurrentCastleLowerFloatSlot157, zero: "GRkZGRkZGRk=", one: "tAsZGRkZGRk=", some: "IeBwtlioxLs="},
		{name: "slot158", fn: encodeCurrentCastleLowerFloatSlot158, zero: "1dXV1dXV1dU=", one: "/iXV1dXV1dU=", some: "FavEzjceS6Y="},
		{name: "slot160", fn: encodeCurrentCastleLowerFloatSlot160, zero: "MDAwMDAwMDA=", one: "LfMwMDAwMDA=", some: "Mamn7JiuKw4="},
		{name: "slot162", fn: encodeCurrentCastleLowerFloatSlot162, zero: "Nzc3Nzc3Nzc=", one: "ysM3Nzc3Nzc=", some: "B+JCBgXyqnA="},
		{name: "slot163", fn: encodeCurrentCastleLowerFloatSlot163, zero: "ampqampqamo=", one: "tFtqampqamo=", some: "6hZB1B6F5yQ="},
		{name: "slot164", fn: encodeCurrentCastleLowerFloatSlot164, zero: "DQ0NDQ0NDQ0=", one: "fFYNDQ0NDQ0=", some: "Cemne6vEMvs="},
		{name: "slot165", fn: encodeCurrentCastleLowerFloatSlot165, zero: "ycnJycnJyck=", one: "lvnJycnJyck=", some: "CUdUhkM2504="},
		{name: "slot166", fn: encodeCurrentCastleLowerFloatSlot166, zero: "tLS0tLS0tLQ=", one: "SHe0tLS0tLQ=", some: "tc3DCNzKTmk="},
		{name: "slot167", fn: encodeCurrentCastleLowerFloatSlot167, zero: "Pz8/Pz8/Pz8=", one: "5q8/Pz8/Pz8=", some: "f21UVslGzV4="},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.fn(0); got != tc.zero {
				t.Fatalf("%s(0) = %q, want %q", tc.name, got, tc.zero)
			}
			if got := tc.fn(1); got != tc.one {
				t.Fatalf("%s(1) = %q, want %q", tc.name, got, tc.one)
			}
			if got := tc.fn(123.456); got != tc.some {
				t.Fatalf("%s(123.456) = %q, want %q", tc.name, got, tc.some)
			}
		})
	}
}

func TestCurrentCastleSlot6Hash(t *testing.T) {
	cases := []struct {
		input string
		hash  uint32
		hex   string
	}{
		{input: "", hash: 0, hex: "0"},
		{input: "abc", hash: 1712826425, hex: "6617a839"},
		{input: "X Castle", hash: 2975336394, hex: "b15807ca"},
		{input: "function test() { return 1; }", hash: 1625178382, hex: "60de410e"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			if got := currentCastleSlot6Hash(tc.input); got != tc.hash {
				t.Fatalf("currentCastleSlot6Hash() = %d/%#x, want %d/%#x", got, got, tc.hash, tc.hash)
			}
			if got := currentCastleSlot6HashHex(tc.input); got != tc.hex {
				t.Fatalf("currentCastleSlot6HashHex() = %q, want %q", got, tc.hex)
			}
		})
	}
	if got, want := encodeCurrentCastleSlot6Component5Fallback(currentCastleSlot6HashHex("abc")), "vb24umi7trA="; got != want {
		t.Fatalf("encoded component-5 hash = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6AH(currentCastleSlot6HashHex("abc")), "urpamtS5Gpk="; got != want {
		t.Fatalf("encoded aH hash = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6EJ(currentCastleSlot6HashHex("abc")), "RERKRupYTlo="; got != want {
		t.Fatalf("encoded ej hash = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot6A8(currentCastleSlot6HashHex("abc")), "Dw//H/Jv338="; got != want {
		t.Fatalf("encoded a8 hash = %q, want %q", got, want)
	}
}

func TestCurrentCastleSlot6AdditionalByteTransforms(t *testing.T) {
	cases := []struct {
		name string
		got  string
		want string
	}{
		{name: "iT abc", got: encodeCurrentCastleSlot6IT("abc"), want: "fjss"},
		{name: "iT x castle", got: encodeCurrentCastleSlot6IT("X Castle"), want: "dihefhEbtHM="},
		{name: "aG abc", got: encodeCurrentCastleSlot6AG("abc"), want: "UO0O"},
		{name: "aG x castle", got: encodeCurrentCastleSlot6AG("X Castle"), want: "J+/uUB67s8w="},
		{name: "aD abc", got: encodeCurrentCastleSlot6AD("abc"), want: "4jKC"},
		{name: "aD x castle", got: encodeCurrentCastleSlot6AD("X Castle"), want: "p+as4u290qI="},
		{name: "t7 abc", got: encodeCurrentCastleSlot6T7("abc"), want: "LcW9"},
		{name: "t7 x castle", got: encodeCurrentCastleSlot6T7("X Castle"), want: "a0W1LblQVgw="},
		{name: "id abc", got: encodeCurrentCastleSlot6ID("abc"), want: "1vYX"},
		{name: "id x castle", got: encodeCurrentCastleSlot6ID("X Castle"), want: "ta4T1hk5OFc="},
		{name: "fu abc", got: encodeCurrentCastleSlot6FU("abc"), want: "JLSE"},
		{name: "fu x castle", got: encodeCurrentCastleSlot6FU("X Castle"), want: "V2WMJIBQUiM="},
		{name: "ej abc", got: encodeCurrentCastleSlot6EJ("abc"), want: "6uzu"},
		{name: "ej x castle", got: encodeCurrentCastleSlot6EJ("X Castle"), want: "mGiu6s7A8OI="},
		{name: "l9 abc", got: encodeCurrentCastleSlot6L9("abc"), want: "bFxM"},
		{name: "l9 x castle", got: encodeCurrentCastleSlot6L9("X Castle"), want: "/wBObE29PKw="},
		{name: "tx abc", got: encodeCurrentCastleSlot6TX("abc"), want: "uLm+"},
		{name: "tx x castle", got: encodeCurrentCastleSlot6TX("X Castle"), want: "o/veuI6Pt7w="},
		{name: "li abc", got: encodeCurrentCastleSlot6LI("abc"), want: "h4mL"},
		{name: "li x castle", got: encodeCurrentCastleSlot6LI("X Castle"), want: "dQVLh6utnY8="},
		{name: "fL abc", got: encodeCurrentCastleSlot6FL("abc"), want: "B4YG"},
		{name: "fL x castle", got: encodeCurrentCastleSlot6FL("X Castle"), want: "g2cWBw6xjQk="},
		{name: "iO abc", got: encodeCurrentCastleSlot6IO("abc"), want: "LTQ/"},
		{name: "iO x castle", got: encodeCurrentCastleSlot6IO("X Castle"), want: "cuofLa+2/gk="},
		{name: "f3 abc", got: encodeCurrentCastleSlot6F3("abc"), want: "HEqw"},
		{name: "f3 x castle", got: encodeCurrentCastleSlot6F3("X Castle"), want: "pCK6HLPpaNQ="},
		{name: "aj abc", got: encodeCurrentCastleSlot6AJ("abc"), want: "u24h"},
		{name: "aj x castle", got: encodeCurrentCastleSlot6AJ("X Castle"), want: "ZwSiu2AXdYo="},
		{name: "component17 abc", got: encodeCurrentCastleSlot6Component17("abc"), want: "np+g"},
		{name: "component17 x castle", got: encodeCurrentCastleSlot6Component17("X Castle"), want: "pd3AnpCJkZo="},
		{name: "component19 abc", got: encodeCurrentCastleSlot6Component19("abc"), want: "Wt9k"},
		{name: "component19 x castle", got: encodeCurrentCastleSlot6Component19("X Castle"), want: "rZXEWrQ5EW4="},
		{name: "component27 abc", got: encodeCurrentCastleSlot6Component27("abc"), want: "hYmN"},
		{name: "component27 x castle", got: encodeCurrentCastleSlot6Component27("X Castle"), want: "YYANhc3RsZU="},
		{name: "tm abc", got: encodeCurrentCastleSlot6TM("abc"), want: "kYmB"},
		{name: "tm x castle", got: encodeCurrentCastleSlot6TM("X Castle"), want: "WJuAkQE5+bE="},
		{name: "tF abc", got: encodeCurrentCastleSlot6TUpperF("abc"), want: "veDf"},
		{name: "tF x castle", got: encodeCurrentCastleSlot6TUpperF("X Castle"), want: "tn6/vc/S6uE="},
		{name: "fs abc", got: encodeCurrentCastleSlot6FS("abc"), want: "ZZj/"},
		{name: "fs x castle", got: encodeCurrentCastleSlot6FS("X Castle"), want: "KB7XZcsvCXY="},
		{name: "component31 abc", got: encodeCurrentCastleSlot6Component31("abc"), want: "s7W3"},
		{name: "component31 x castle", got: encodeCurrentCastleSlot6Component31("X Castle"), want: "oTF3s9fZybs="},
		{name: "fF abc", got: encodeCurrentCastleSlot6FF("abc"), want: "+J1B"},
		{name: "fF x castle", got: encodeCurrentCastleSlot6FF("X Castle"), want: "sLQx+Amuigo="},
		{name: "aU abc", got: encodeCurrentCastleSlot6AU("abc"), want: "C7N6"},
		{name: "aU x castle", got: encodeCurrentCastleSlot6AU("X Castle"), want: "mPwKC2Kq5mk="},
		{name: "aQ abc", got: encodeCurrentCastleSlot6AQ("abc"), want: "eOdT"},
		{name: "aQ x castle", got: encodeCurrentCastleSlot6AQ("X Castle"), want: "bA3ReJMff0o="},
		{name: "lW abc", got: encodeCurrentCastleSlot6LW("abc"), want: "45xV"},
		{name: "lW x castle", got: encodeCurrentCastleSlot6LW("X Castle"), want: "Yuo14+We1sc="},
		{name: "tM abc", got: encodeCurrentCastleSlot6TUpperM("abc"), want: "UGN0"},
		{name: "tM x castle", got: encodeCurrentCastleSlot6TUpperM("X Castle"), want: "rr40UJWmFpg="},
		{name: "th abc", got: encodeCurrentCastleSlot6TH("abc"), want: "lKS0"},
		{name: "th x castle", got: encodeCurrentCastleSlot6TH("X Castle"), want: "BICylLXFRdQ="},
		{name: "is abc", got: encodeCurrentCastleSlot6IS("abc"), want: "xtYm"},
		{name: "is x castle", got: encodeCurrentCastleSlot6IS("X Castle"), want: "cfIgxic3twY="},
		{name: "a8 abc", got: encodeCurrentCastleSlot6A8("abc"), want: "8sLS"},
		{name: "a8 x castle", got: encodeCurrentCastleSlot6A8("X Castle"), want: "Ye7Q8tMjozI="},
		{name: "lM abc", got: encodeCurrentCastleSlot6LUpperM("abc"), want: "w0zV"},
		{name: "lM x castle", got: encodeCurrentCastleSlot6LUpperM("X Castle"), want: "8vq1w2Xupuc="},
		{name: "eA abc", got: encodeCurrentCastleSlot6EA("abc"), want: "IRsd"},
		{name: "eA x castle", got: encodeCurrentCastleSlot6EA("X Castle"), want: "L55dIfz2Bxk="},
		{name: "e3 abc", got: encodeCurrentCastleSlot6E3("abc"), want: "qwBV"},
		{name: "e3 x castle", got: encodeCurrentCastleSlot6E3("X Castle"), want: "Bh61qyV6Avc="},
		{name: "t2 abc", got: encodeCurrentCastleSlot6T2("abc"), want: "CNeS"},
		{name: "t2 x castle", got: encodeCurrentCastleSlot6T2("X Castle"), want: "1Q3yCOKxORw="},
		{name: "tX abc", got: encodeCurrentCastleSlot6TUpperX("abc"), want: "Bh4W"},
		{name: "tX x castle", got: encodeCurrentCastleSlot6TUpperX("X Castle"), want: "zwwXBpaubiY="},
		{name: "l5 abc", got: encodeCurrentCastleSlot6L5("abc"), want: "WKr8"},
		{name: "l5 x castle", got: encodeCurrentCastleSlot6L5("X Castle"), want: "doa8WBxu3qA="},
		{name: "component46 abc", got: encodeCurrentCastleSlot6Component46("abc"), want: "fX5/"},
		{name: "component46 x castle", got: encodeCurrentCastleSlot6Component46("X Castle"), want: "ZLxffY+QeIE="},
		{name: "tz abc", got: encodeCurrentCastleSlot6TZ("abc"), want: "cFQ4"},
		{name: "tz x castle", got: encodeCurrentCastleSlot6TZ("X Castle"), want: "ioM0cDoeHeA="},
		{name: "tP abc", got: encodeCurrentCastleSlot6TP("abc"), want: "+WTP"},
		{name: "tP x castle", got: encodeCurrentCastleSlot6TP("X Castle"), want: "Ns5v+X/qkqU="},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("got %q, want %q", tc.got, tc.want)
			}
		})
	}
}

func TestPopulateCurrentCastleHighTimingSlots(t *testing.T) {
	payload := make([]any, currentCastlePayloadSlots)
	if err := populateCurrentCastleHighTimingSlots(payload, map[int]float64{
		431: 0,
		480: 123.456,
		493: 1,
	}); err != nil {
		t.Fatalf("populateCurrentCastleHighTimingSlots() error = %v", err)
	}
	for slot := 431; slot <= 493; slot++ {
		if payload[slot] == nil {
			t.Fatalf("payload[%d] is nil", slot)
		}
	}
	if got, want := payload[431], "xsbGxsbGxsY="; got != want {
		t.Fatalf("payload[431] = %q, want %q", got, want)
	}
	if got, want := payload[442], "hoZgICDGxsY="; got != want {
		t.Fatalf("payload[442] = %q, want %q", got, want)
	}
	if got, want := payload[480], "hhjbKVyZ+PE="; got != want {
		t.Fatalf("payload[480] = %q, want %q", got, want)
	}
	if got, want := payload[493], "OTbGxsbGxsY="; got != want {
		t.Fatalf("payload[493] = %q, want %q", got, want)
	}
	if err := populateCurrentCastleHighTimingSlots(make([]any, currentCastlePayloadSlots-1), nil); err == nil {
		t.Fatalf("populateCurrentCastleHighTimingSlots() with short payload succeeded")
	}
}

func TestPopulateCurrentCastleLowerFloatSlots(t *testing.T) {
	payload := newCurrentCastlePayloadScaffold()
	if err := populateCurrentCastleLowerFloatSlots(payload, map[int]float64{
		1:  1,
		3:  123.456,
		19: 1,
	}); err != nil {
		t.Fatalf("populateCurrentCastleLowerFloatSlots() error = %v", err)
	}
	if got, want := payload[1], "KWfg4ODg4OA="; got != want {
		t.Fatalf("payload[1] = %q, want %q", got, want)
	}
	if got, want := payload[3], "agkBHM1lE+c="; got != want {
		t.Fatalf("payload[3] = %q, want %q", got, want)
	}
	if got, want := payload[4], "7Ozs7Ozs7Ow="; got != want {
		t.Fatalf("payload[4] = %q, want %q", got, want)
	}
	if got, want := payload[19], "4ZQAAAAAAAA="; got != want {
		t.Fatalf("payload[19] = %q, want %q", got, want)
	}
	if got, want := payload[21], "1dXV1dXV1dU="; got != want {
		t.Fatalf("payload[21] = %q, want %q", got, want)
	}
	if got, want := payload[39], "AAAAAAAAAAA="; got != want {
		t.Fatalf("payload[39] = %q, want %q", got, want)
	}
	if got, want := payload[43], "q6urq6urq6s="; got != want {
		t.Fatalf("payload[43] = %q, want %q", got, want)
	}
	if got, want := payload[63], "a2tra2tra2s="; got != want {
		t.Fatalf("payload[63] = %q, want %q", got, want)
	}
	if got, want := payload[65], "g4ODg4ODg4M="; got != want {
		t.Fatalf("payload[65] = %q, want %q", got, want)
	}
	if got, want := payload[107], "LCwsLCwsLCw="; got != want {
		t.Fatalf("payload[107] = %q, want %q", got, want)
	}
	if got, want := payload[167], "Pz8/Pz8/Pz8="; got != want {
		t.Fatalf("payload[167] = %q, want %q", got, want)
	}
	if got, want := payload[169], "ioqKioqKioo="; got != want {
		t.Fatalf("payload[169] = %q, want %q", got, want)
	}
	if got, want := payload[226], "o6Ojo6Ojo6M="; got != want {
		t.Fatalf("payload[226] = %q, want %q", got, want)
	}
	if got, want := payload[255], "MjIyMjIyMjI="; got != want {
		t.Fatalf("payload[255] = %q, want %q", got, want)
	}
	if got, want := payload[294], "NDQ0NDQ0NDQ="; got != want {
		t.Fatalf("payload[294] = %q, want %q", got, want)
	}
	if got, want := payload[291], "hoaGhoaGhoY="; got != want {
		t.Fatalf("payload[291] = %q, want %q", got, want)
	}
	if got, want := payload[324], "WlpaWlpaWlo="; got != want {
		t.Fatalf("payload[324] = %q, want %q", got, want)
	}
	if got, want := payload[325], "GhoaGhoaGho="; got != want {
		t.Fatalf("payload[325] = %q, want %q", got, want)
	}
	if got, want := payload[373], "pKSkpKSkpKQ="; got != want {
		t.Fatalf("payload[373] = %q, want %q", got, want)
	}
	if got, want := payload[407], "2tra2tra2to="; got != want {
		t.Fatalf("payload[407] = %q, want %q", got, want)
	}
	if got, want := payload[427], "h4eHh4eHh4c="; got != want {
		t.Fatalf("payload[427] = %q, want %q", got, want)
	}
	for slot := range currentCastleLowerFloatSlotEncoders {
		if _, ok := payload[slot].(string); !ok || payload[slot] == "" {
			t.Fatalf("payload[%d] was not populated with a non-empty string: %#v", slot, payload[slot])
		}
	}
	if err := populateCurrentCastleLowerFloatSlots(make([]any, currentCastlePayloadSlots-1), nil); err == nil {
		t.Fatalf("populateCurrentCastleLowerFloatSlots() with short payload succeeded")
	}
}

func TestCurrentCastleNumericSlotEncoders(t *testing.T) {
	if got, want := currentCastleIndexEUnderscore(0), uint32(0x03508000); got != want {
		t.Fatalf("currentCastleIndexEUnderscore(0) = %#x, want %#x", got, want)
	}
	if got, want := currentCastleIndexR7(currentCastleIndexEUnderscore(0)), uint32(0); got != want {
		t.Fatalf("currentCastleIndexR7(currentCastleIndexEUnderscore(0)) = %#x, want %#x", got, want)
	}

	simpleCases := []struct {
		name string
		fn   func(uint32) uint32
		zero uint32
		one  uint32
	}{
		{name: "slot20", fn: encodeCurrentCastleNumericSlot20, zero: 0xd2c1284a, one: 0x5fd1c84a},
		{name: "slot154", fn: encodeCurrentCastleNumericSlot154, zero: 0x8eca32a4, one: 0xf13ce181},
		{name: "slot159", fn: encodeCurrentCastleNumericIdentity, zero: 0, one: 1},
		{name: "slot182", fn: currentCastleTB, zero: 0x00000354, one: 0x000003df},
		{name: "slot316", fn: currentCastleIndexEUnderscore, zero: 0x03508000, one: 0x0350a000},
		{name: "slot256", fn: encodeCurrentCastleNumericSlot256, zero: 0xe39b0b56, one: 0x4ea7245b},
		{name: "slot348", fn: encodeCurrentCastleNumericSlot348, zero: 0x54297e0d, one: 0x4ff13582},
	}
	for _, tc := range simpleCases {
		if got := tc.fn(0); got != tc.zero {
			t.Fatalf("%s(0) = %#x, want %#x", tc.name, got, tc.zero)
		}
		if got := tc.fn(1); got != tc.one {
			t.Fatalf("%s(1) = %#x, want %#x", tc.name, got, tc.one)
		}
	}
	affineCases := []struct {
		slot int
		zero uint32
		one  uint32
	}{
		{slot: 26, zero: 0x60bd09b2, one: 0x670022bf},
		{slot: 27, zero: 0xb330012c, one: 0x0dd8798d},
		{slot: 28, zero: 0x244dce95, one: 0x6dfeaec2},
		{slot: 31, zero: 0xb78dc961, one: 0xaae15d40},
		{slot: 42, zero: 0x9e224c69, one: 0x944782d0},
		{slot: 49, zero: 0x06d9e43d, one: 0x6a5b8e70},
		{slot: 53, zero: 0x8dd8d0c3, one: 0x27b19724},
		{slot: 56, zero: 0xd9931d51, one: 0xacc004ba},
		{slot: 64, zero: 0xc538ce0b, one: 0x9171777e},
		{slot: 67, zero: 0xf3c0363d, one: 0x4080d6d0},
		{slot: 68, zero: 0xecbd981e, one: 0x76dc5125},
		{slot: 71, zero: 0x677c239f, one: 0xce58e7fa},
		{slot: 73, zero: 0x96d91174, one: 0xb09c39df},
		{slot: 87, zero: 0xaa10f53a, one: 0xeca6e7eb},
		{slot: 98, zero: 0x95bd552c, one: 0xb0a64635},
		{slot: 102, zero: 0x4bcca730, one: 0x52157915},
		{slot: 106, zero: 0xb5f76f6d, one: 0xb2bb1c70},
		{slot: 115, zero: 0xf041ff29, one: 0x9577f8e2},
		{slot: 116, zero: 0x8c84e89a, one: 0x3dfbe269},
		{slot: 121, zero: 0xdbab50b2, one: 0x0fb7ef7f},
		{slot: 123, zero: 0x89d7f1b8, one: 0x7712f7cf},
		{slot: 133, zero: 0x520795d3, one: 0x61439c1a},
		{slot: 147, zero: 0xcd6c2b6c, one: 0x6d80ef95},
		{slot: 149, zero: 0x3b4d9896, one: 0xf4921071},
		{slot: 150, zero: 0x1e187f66, one: 0x5ba191cb},
		{slot: 152, zero: 0xf83df829, one: 0x0160726c},
		{slot: 168, zero: 0x4d504e4f, one: 0xef51b58e},
		{slot: 194, zero: 0xc1cb353f, one: 0xed628648},
		{slot: 196, zero: 0x6d935c25, one: 0x6ddfc840},
		{slot: 198, zero: 0x1a205ab2, one: 0xe350ecf9},
		{slot: 203, zero: 0x45c65c75, one: 0xc14d992a},
		{slot: 205, zero: 0x46fd4c27, one: 0x67a5337e},
		{slot: 206, zero: 0x086aa955, one: 0x40c6c596},
		{slot: 207, zero: 0x2496e7c2, one: 0x0f608701},
		{slot: 209, zero: 0xfed3363f, one: 0x7feea612},
		{slot: 212, zero: 0x2c3ea5b3, one: 0x665c3ee0},
		{slot: 221, zero: 0xfd5c8336, one: 0x9f7dec27},
		{slot: 223, zero: 0xdeb6d52a, one: 0xcb0594c1},
		{slot: 224, zero: 0xa9cd6a67, one: 0x5dbe5d9a},
		{slot: 228, zero: 0xa7e0384c, one: 0x56e8784d},
		{slot: 240, zero: 0xf4661fb5, one: 0x11924198},
		{slot: 246, zero: 0x3c8624df, one: 0x61388f30},
		{slot: 247, zero: 0xed6822f5, one: 0xea4dcacc},
		{slot: 257, zero: 0xac01bfae, one: 0x0e46a163},
		{slot: 258, zero: 0x3d7667d5, one: 0xceb0120a},
		{slot: 260, zero: 0x2590de1a, one: 0xa8853377},
		{slot: 266, zero: 0xa3c420cf, one: 0xa4ee6d12},
		{slot: 267, zero: 0x2b30156a, one: 0xad2cd039},
		{slot: 277, zero: 0xb939ef29, one: 0x4aadc51e},
		{slot: 288, zero: 0x4d25ec39, one: 0x7006f918},
		{slot: 293, zero: 0x9f2c27a4, one: 0x614ef36f},
		{slot: 297, zero: 0x8127b4d5, one: 0x2d162ffa},
		{slot: 311, zero: 0x42f7a6bc, one: 0xf5810f03},
		{slot: 312, zero: 0xca4c0beb, one: 0x178b944c},
		{slot: 318, zero: 0x0fc46406, one: 0xb618188b},
		{slot: 321, zero: 0xb7481870, one: 0x746558bb},
		{slot: 329, zero: 0x828e6d80, one: 0x7dde71e5},
		{slot: 334, zero: 0x79a7bbd2, one: 0x2c829589},
		{slot: 339, zero: 0xfd1e66f3, one: 0x71de90e4},
		{slot: 349, zero: 0xcb462199, one: 0x4d17f152},
		{slot: 354, zero: 0x471a0d4d, one: 0xed05f1de},
		{slot: 356, zero: 0xbdfe52d9, one: 0x32782ba6},
		{slot: 359, zero: 0xa583bc26, one: 0x1699b055},
		{slot: 366, zero: 0xa9ac1690, one: 0xc211c04b},
		{slot: 369, zero: 0xc33c50e9, one: 0x4b6d07f4},
		{slot: 374, zero: 0xe6b394f2, one: 0x5ed52a3b},
		{slot: 375, zero: 0xcf565231, one: 0x4d737d2a},
		{slot: 376, zero: 0x30b3a87e, one: 0xf233f3b7},
		{slot: 386, zero: 0xd1ee53e2, one: 0x78b37579},
		{slot: 388, zero: 0x6cac25da, one: 0xb5035037},
		{slot: 390, zero: 0x9bc2f748, one: 0x03f0733d},
		{slot: 393, zero: 0xfe9478a7, one: 0x9ed36934},
		{slot: 395, zero: 0xb427bace, one: 0xf991e0ef},
		{slot: 399, zero: 0xd6c6c42c, one: 0x3fc1896b},
		{slot: 406, zero: 0x29fc7252, one: 0x4da110bf},
		{slot: 422, zero: 0xecfc92b2, one: 0x9632461b},
		{slot: 425, zero: 0x8e16c592, one: 0xb2160d79},
		{slot: 430, zero: 0x2b90120f, one: 0x40ec1bb2},
	}
	for _, tc := range affineCases {
		fn, ok := currentCastleNumericSlotEncoders[tc.slot]
		if !ok {
			t.Fatalf("missing numeric slot encoder for slot %d", tc.slot)
		}
		if got := fn(0); got != tc.zero {
			t.Fatalf("numeric slot %d raw 0 = %#x, want %#x", tc.slot, got, tc.zero)
		}
		if got := fn(1); got != tc.one {
			t.Fatalf("numeric slot %d raw 1 = %#x, want %#x", tc.slot, got, tc.one)
		}
	}
	if got, want := encodeCurrentCastleNumericSlot20(0xdae88a03), uint32(0x15bbadc1); got != want {
		t.Fatalf("encodeCurrentCastleNumericSlot20(sentinel) = %#x, want %#x", got, want)
	}
	if got, want := encodeCurrentCastleNumericSlot154(0xdae88a03), uint32(0xbd4808ff); got != want {
		t.Fatalf("encodeCurrentCastleNumericSlot154(sentinel) = %#x, want %#x", got, want)
	}
	if got, want := currentCastleTB(8), uint32(0x000007ac); got != want {
		t.Fatalf("currentCastleTB(8) = %#x, want %#x", got, want)
	}
	if got, want := currentCastleNumericSlotEncoders[257](currentCastleTB(0)), uint32(0xb0e99ef2); got != want {
		t.Fatalf("numeric slot 257 false default = %#x, want %#x", got, want)
	}
	if got, want := currentCastleNumericSlotEncoders[257](currentCastleTB(1)), uint32(0x3e5af065); got != want {
		t.Fatalf("numeric slot 257 true-like default = %#x, want %#x", got, want)
	}
	if got, want := currentCastleNumericSlotEncoders[289](currentCastleTB(0)), uint32(0x0001d1f0); got != want {
		t.Fatalf("numeric slot 289 false default = %#x, want %#x", got, want)
	}
	if got, want := currentCastleNumericSlotEncoders[289](currentCastleTB(1)), uint32(0x00021d69); got != want {
		t.Fatalf("numeric slot 289 true-like default = %#x, want %#x", got, want)
	}
	if got, want := currentCastleNumericSlotEncoders[152](0xdae88a03), uint32(0xb8ee54f2); got != want {
		t.Fatalf("numeric slot 152 sentinel = %#x, want %#x", got, want)
	}
	if got, want := currentCastleNumericSlotEncoders[258](0xdae88a03), uint32(0x5c577ba0); got != want {
		t.Fatalf("numeric slot 258 sentinel = %#x, want %#x", got, want)
	}
	sentinelCases := []struct {
		slot int
		want uint32
	}{
		{slot: 26, want: 0xb8d522a5},
		{slot: 28, want: 0xd800b11c},
		{slot: 31, want: 0x20b7bafe},
		{slot: 42, want: 0xff1f0002},
		{slot: 49, want: 0x34e764d6},
		{slot: 64, want: 0xe3b2c864},
		{slot: 67, want: 0xe83187aa},
		{slot: 68, want: 0x572da517},
		{slot: 71, want: 0x08140144},
		{slot: 87, want: 0x938e374d},
		{slot: 98, want: 0x813ade47},
		{slot: 102, want: 0x919f8edf},
		{slot: 106, want: 0x5338546a},
		{slot: 116, want: 0xde374007},
		{slot: 123, want: 0xa65269fd},
		{slot: 133, want: 0x1cb5eea8},
		{slot: 168, want: 0xfdcef110},
		{slot: 198, want: 0xa79b0e6b},
		{slot: 203, want: 0x884f8094},
		{slot: 207, want: 0x1032cf7f},
		{slot: 209, want: 0xd57543b8},
		{slot: 212, want: 0xa8aa2f3a},
		{slot: 221, want: 0xc1423045},
		{slot: 223, want: 0x510aadef},
		{slot: 224, want: 0x945b9200},
		{slot: 228, want: 0xf62f824b},
		{slot: 240, want: 0xf5609fd2},
		{slot: 246, want: 0x6a5ffdd2},
		{slot: 247, want: 0x52eb007a},
		{slot: 257, want: 0x3f3e4bf9},
		{slot: 266, want: 0xb4e6c28c},
		{slot: 277, want: 0xe1cc5734},
		{slot: 289, want: 0xdc42f2f5},
		{slot: 312, want: 0x4f85ef0e},
		{slot: 316, want: 0x1490fb5d},
		{slot: 318, want: 0x2858fd81},
		{slot: 321, want: 0xd4c04751},
		{slot: 356, want: 0x4f2df80c},
		{slot: 366, want: 0xdfa64ed5},
		{slot: 374, want: 0x6746facd},
		{slot: 375, want: 0xa8e6ed38},
		{slot: 376, want: 0x2e9f1745},
		{slot: 386, want: 0x6185984b},
		{slot: 399, want: 0xc363f8ed},
		{slot: 422, want: 0xaf6c7949},
		{slot: 430, want: 0xe14d50f8},
	}
	for _, tc := range sentinelCases {
		if got := currentCastleNumericSlotEncoders[tc.slot](0xdae88a03); got != tc.want {
			t.Fatalf("numeric slot %d sentinel = %#x, want %#x", tc.slot, got, tc.want)
		}
	}

	slot341Cases := []struct {
		name string
		raw  uint32
		want uint32
	}{
		{name: "falseFlag", raw: 0, want: 0x0000f6a7},
		{name: "trueFlag", raw: 1, want: 0x0000793c},
	}
	for _, tc := range slot341Cases {
		if got := encodeCurrentCastleNumericSlot341(tc.raw); got != tc.want {
			t.Fatalf("encodeCurrentCastleNumericSlot341(%s) = %#x, want %#x", tc.name, got, tc.want)
		}
	}

	slot428Cases := []struct {
		name string
		raw  uint32
		want uint32
	}{
		{name: "false", raw: 0, want: 0xeb00e01b},
		{name: "true", raw: 1, want: 0xc0f29076},
		{name: "notFunction", raw: 0xdae88a03, want: 0xc64546ec},
		{name: "undefinedTProtect", raw: 0x7b862786, want: 0xfc5c96bd},
	}
	for _, tc := range slot428Cases {
		if got := encodeCurrentCastleNumericSlot428(tc.raw); got != tc.want {
			t.Fatalf("encodeCurrentCastleNumericSlot428(%s) = %#x, want %#x", tc.name, got, tc.want)
		}
	}

	slot429Cases := []struct {
		name string
		raw  uint32
		want uint32
	}{
		{name: "huFalseFallback", raw: currentCastleIndexEUnderscore(0), want: 0x14808cbe},
		{name: "huTrueFallback", raw: currentCastleIndexEUnderscore(1), want: 0x7229acbe},
		{name: "zero", raw: 0, want: 0x7bac0cbe},
		{name: "one", raw: 1, want: 0x57fefa07},
	}
	for _, tc := range slot429Cases {
		if got := encodeCurrentCastleNumericSlot429(tc.raw); got != tc.want {
			t.Fatalf("encodeCurrentCastleNumericSlot429(%s) = %#x, want %#x", tc.name, got, tc.want)
		}
	}
}

func TestPopulateCurrentCastleNumericSlots(t *testing.T) {
	payload := newCurrentCastlePayloadScaffold()
	if err := populateCurrentCastleNumericSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleNumericSlots() error = %v", err)
	}
	for slot, want := range defaultCurrentCastleNumericSlotFinalValues() {
		if got := payload[slot]; got != want {
			t.Fatalf("payload[%d] = %#v, want %#v", slot, got, want)
		}
	}

	if err := populateCurrentCastleNumericSlots(payload, map[int]uint32{
		20:  1,
		26:  1,
		27:  1,
		28:  1,
		31:  currentCastleTB(0),
		42:  currentCastleIndexEUnderscore(1),
		49:  1,
		64:  currentCastleIndexEUnderscore(1),
		67:  currentCastleIndexEUnderscore(1),
		68:  1,
		71:  currentCastleIndexEUnderscore(1),
		87:  currentCastleIndexEUnderscore(1),
		98:  1,
		102: 1,
		106: currentCastleTB(0),
		116: 1,
		121: 1,
		123: currentCastleIndexEUnderscore(1),
		133: currentCastleTB(1),
		152: 1,
		154: 1,
		159: 1,
		168: 1,
		182: 1,
		198: currentCastleTB(1),
		203: 1,
		207: 1,
		209: 1,
		212: 1,
		221: 1,
		223: 1,
		224: currentCastleIndexEUnderscore(1),
		228: currentCastleIndexEUnderscore(1),
		240: 1,
		246: 1,
		247: currentCastleIndexEUnderscore(1),
		256: 1,
		257: currentCastleTB(1),
		258: 1,
		266: 1,
		267: 1,
		277: 1,
		288: 1,
		289: currentCastleTB(1),
		297: 1,
		312: 1,
		316: currentCastleTB(1),
		318: currentCastleIndexEUnderscore(1),
		321: currentCastleIndexEUnderscore(1),
		334: currentCastleTB(1),
		341: 1,
		348: 1,
		356: currentCastleTB(1),
		359: currentCastleTB(1),
		366: 1,
		369: currentCastleIndexEUnderscore(1),
		374: 1,
		375: currentCastleTB(1),
		376: currentCastleIndexEUnderscore(1),
		386: 1,
		399: 1,
		406: 1,
		422: 1,
		425: 1,
		428: 1,
		429: 1,
		430: currentCastleTB(1),
	}); err != nil {
		t.Fatalf("populateCurrentCastleNumericSlots() with override error = %v", err)
	}
	if got, want := payload[428], uint32(0xc0f29076); got != want {
		t.Fatalf("override payload[428] = %#v, want %#v", got, want)
	}
	if got, want := payload[429], uint32(0x57fefa07); got != want {
		t.Fatalf("override payload[429] = %#v, want %#v", got, want)
	}
	if got, want := payload[20], uint32(0x5fd1c84a); got != want {
		t.Fatalf("override payload[20] = %#v, want %#v", got, want)
	}
	if got, want := payload[26], uint32(0x670022bf); got != want {
		t.Fatalf("override payload[26] = %#v, want %#v", got, want)
	}
	if got, want := payload[28], uint32(0x6dfeaec2); got != want {
		t.Fatalf("override payload[28] = %#v, want %#v", got, want)
	}
	if got, want := payload[31], uint32(0xce51c22d); got != want {
		t.Fatalf("override payload[31] = %#v, want %#v", got, want)
	}
	if got, want := payload[42], uint32(0x78f2ac69); got != want {
		t.Fatalf("override payload[42] = %#v, want %#v", got, want)
	}
	if got, want := payload[49], uint32(0x6a5b8e70); got != want {
		t.Fatalf("override payload[49] = %#v, want %#v", got, want)
	}
	if got, want := payload[64], uint32(0x7c10ae0b); got != want {
		t.Fatalf("override payload[64] = %#v, want %#v", got, want)
	}
	if got, want := payload[67], uint32(0xce74563d); got != want {
		t.Fatalf("override payload[67] = %#v, want %#v", got, want)
	}
	if got, want := payload[68], uint32(0x76dc5125); got != want {
		t.Fatalf("override payload[68] = %#v, want %#v", got, want)
	}
	if got, want := payload[71], uint32(0xe94a039f); got != want {
		t.Fatalf("override payload[71] = %#v, want %#v", got, want)
	}
	if got, want := payload[87], uint32(0x50e3553a); got != want {
		t.Fatalf("override payload[87] = %#v, want %#v", got, want)
	}
	if got, want := payload[98], uint32(0xb0a64635); got != want {
		t.Fatalf("override payload[98] = %#v, want %#v", got, want)
	}
	if got, want := payload[102], uint32(0x52157915); got != want {
		t.Fatalf("override payload[102] = %#v, want %#v", got, want)
	}
	if got, want := payload[106], uint32(0xdc9222f1); got != want {
		t.Fatalf("override payload[106] = %#v, want %#v", got, want)
	}
	if got, want := payload[116], uint32(0x3dfbe269); got != want {
		t.Fatalf("override payload[116] = %#v, want %#v", got, want)
	}
	if got, want := payload[123], uint32(0x03d991b8); got != want {
		t.Fatalf("override payload[123] = %#v, want %#v", got, want)
	}
	if got, want := payload[133], uint32(0xe032d1a4); got != want {
		t.Fatalf("override payload[133] = %#v, want %#v", got, want)
	}
	if got, want := payload[154], uint32(0xf13ce181); got != want {
		t.Fatalf("override payload[154] = %#v, want %#v", got, want)
	}
	if got, want := payload[212], uint32(0x665c3ee0); got != want {
		t.Fatalf("override payload[212] = %#v, want %#v", got, want)
	}
	if got, want := payload[224], uint32(0x44bd8a67); got != want {
		t.Fatalf("override payload[224] = %#v, want %#v", got, want)
	}
	if got, want := payload[228], uint32(0x830f984c); got != want {
		t.Fatalf("override payload[228] = %#v, want %#v", got, want)
	}
	if got, want := payload[121], uint32(0x0fb7ef7f); got != want {
		t.Fatalf("override payload[121] = %#v, want %#v", got, want)
	}
	if got, want := payload[256], uint32(0x4ea7245b); got != want {
		t.Fatalf("override payload[256] = %#v, want %#v", got, want)
	}
	if got, want := payload[257], uint32(0x3e5af065); got != want {
		t.Fatalf("override payload[257] = %#v, want %#v", got, want)
	}
	if got, want := payload[247], uint32(0xbd08c2f5); got != want {
		t.Fatalf("override payload[247] = %#v, want %#v", got, want)
	}
	if got, want := payload[267], uint32(0xad2cd039); got != want {
		t.Fatalf("override payload[267] = %#v, want %#v", got, want)
	}
	if got, want := payload[288], uint32(0x7006f918); got != want {
		t.Fatalf("override payload[288] = %#v, want %#v", got, want)
	}
	if got, want := payload[289], uint32(0x00021d69); got != want {
		t.Fatalf("override payload[289] = %#v, want %#v", got, want)
	}
	if got, want := payload[348], uint32(0x4ff13582); got != want {
		t.Fatalf("override payload[348] = %#v, want %#v", got, want)
	}
	if got, want := payload[27], uint32(0x0dd8798d); got != want {
		t.Fatalf("override payload[27] = %#v, want %#v", got, want)
	}
	if got, want := payload[152], uint32(0x0160726c); got != want {
		t.Fatalf("override payload[152] = %#v, want %#v", got, want)
	}
	if got, want := payload[159], uint32(1); got != want {
		t.Fatalf("override payload[159] = %#v, want %#v", got, want)
	}
	if got, want := payload[182], uint32(0x000003df); got != want {
		t.Fatalf("override payload[182] = %#v, want %#v", got, want)
	}
	if got, want := payload[168], uint32(0xef51b58e); got != want {
		t.Fatalf("override payload[168] = %#v, want %#v", got, want)
	}
	if got, want := payload[198], uint32(0x97d20b6f); got != want {
		t.Fatalf("override payload[198] = %#v, want %#v", got, want)
	}
	if got, want := payload[203], uint32(0xc14d992a); got != want {
		t.Fatalf("override payload[203] = %#v, want %#v", got, want)
	}
	if got, want := payload[258], uint32(0xceb0120a); got != want {
		t.Fatalf("override payload[258] = %#v, want %#v", got, want)
	}
	if got, want := payload[341], uint32(0x0000793c); got != want {
		t.Fatalf("override payload[341] = %#v, want %#v", got, want)
	}
	if got, want := payload[318], uint32(0x99814406); got != want {
		t.Fatalf("override payload[318] = %#v, want %#v", got, want)
	}
	if got, want := payload[321], uint32(0x3306f870); got != want {
		t.Fatalf("override payload[321] = %#v, want %#v", got, want)
	}
	if got, want := payload[334], uint32(0x4dcfb4eb); got != want {
		t.Fatalf("override payload[334] = %#v, want %#v", got, want)
	}
	if got, want := payload[356], uint32(0xc1a68330); got != want {
		t.Fatalf("override payload[356] = %#v, want %#v", got, want)
	}
	if got, want := payload[359], uint32(0xc25d59db); got != want {
		t.Fatalf("override payload[359] = %#v, want %#v", got, want)
	}
	if got, want := payload[369], uint32(0x92b330e9); got != want {
		t.Fatalf("override payload[369] = %#v, want %#v", got, want)
	}
	if got, want := payload[374], uint32(0x5ed52a3b); got != want {
		t.Fatalf("override payload[374] = %#v, want %#v", got, want)
	}
	if got, want := payload[386], uint32(0x78b37579); got != want {
		t.Fatalf("override payload[386] = %#v, want %#v", got, want)
	}
	if got, want := payload[406], uint32(0x4da110bf); got != want {
		t.Fatalf("override payload[406] = %#v, want %#v", got, want)
	}
	if got, want := payload[399], uint32(0x3fc1896b); got != want {
		t.Fatalf("override payload[399] = %#v, want %#v", got, want)
	}
	if got, want := payload[422], uint32(0x9632461b); got != want {
		t.Fatalf("override payload[422] = %#v, want %#v", got, want)
	}
	if got, want := payload[297], uint32(0x2d162ffa); got != want {
		t.Fatalf("override payload[297] = %#v, want %#v", got, want)
	}
	if got, want := payload[316], uint32(0x03cc6000); got != want {
		t.Fatalf("override payload[316] = %#v, want %#v", got, want)
	}
	if got, want := payload[425], uint32(0xb2160d79); got != want {
		t.Fatalf("override payload[425] = %#v, want %#v", got, want)
	}
	if got, want := payload[366], uint32(0xc211c04b); got != want {
		t.Fatalf("override payload[366] = %#v, want %#v", got, want)
	}
	if got, want := payload[375], uint32(0xc3e69804); got != want {
		t.Fatalf("override payload[375] = %#v, want %#v", got, want)
	}
	if got, want := payload[376], uint32(0x5560087e); got != want {
		t.Fatalf("override payload[376] = %#v, want %#v", got, want)
	}
	if got, want := payload[430], uint32(0x5c3e0c94); got != want {
		t.Fatalf("override payload[430] = %#v, want %#v", got, want)
	}
	if err := populateCurrentCastleNumericSlots(make([]any, currentCastlePayloadSlots-1), nil); err == nil {
		t.Fatalf("populateCurrentCastleNumericSlots() with short payload succeeded")
	}
}

func TestPopulateCurrentCastleObjectSlots(t *testing.T) {
	payload := newCurrentCastlePayloadScaffold()
	if err := populateCurrentCastleNumericSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleNumericSlots() error = %v", err)
	}
	if err := populateCurrentCastleStringSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() error = %v", err)
	}
	slot41, ok := payload[41].(map[string]any)
	if !ok {
		t.Fatalf("payload[41] type = %T, want map[string]any", payload[41])
	}
	if got, want := slot41["V"], "Ls9BGR3QrImUhLGh8DdWMQ=="; got != want {
		t.Fatalf("payload[41].V = %q, want %q", got, want)
	}
	if got, want := slot41["wh"], ""; got != want {
		t.Fatalf("payload[41].wh = %q, want %q", got, want)
	}
	if got, want := slot41["y"], ""; got != want {
		t.Fatalf("payload[41].y = %q, want %q", got, want)
	}
	if got, want := slot41["G"], ""; got != want {
		t.Fatalf("payload[41].G = %q, want %q", got, want)
	}
	slot41Hrc, ok := slot41["Hrc"].(map[string]any)
	if !ok {
		t.Fatalf("payload[41].Hrc type = %T, want map[string]any", slot41["Hrc"])
	}
	if got, want := slot41Hrc["j"], ""; got != want {
		t.Fatalf("payload[41].Hrc.j = %q, want %q", got, want)
	}
	slot77, ok := payload[77].(map[string]any)
	if !ok {
		t.Fatalf("payload[77] type = %T, want map[string]any", payload[77])
	}
	if got, want := slot77["AFS"], "JcW/1uh5mugs0/qG3eXd5Q=="; got != want {
		t.Fatalf("payload[77].AFS = %q, want %q", got, want)
	}
	if got, want := slot77["tu"], ""; got != want {
		t.Fatalf("payload[77].tu = %q, want %q", got, want)
	}
	if got, want := slot77["Xd"], ""; got != want {
		t.Fatalf("payload[77].Xd = %q, want %q", got, want)
	}
	if got, want := slot77["BgV"], ""; got != want {
		t.Fatalf("payload[77].BgV = %q, want %q", got, want)
	}
	slot77OverrideSource := encodeCurrentCastleObjectSlot77FromSource("override")
	if got, want := slot77OverrideSource["AFS"], "LNP2Gx3Q2X7QdIp/Bza9OsmqjM4="; got != want {
		t.Fatalf("override-source payload[77].AFS = %q, want %q", got, want)
	}
	if got, want := slot77OverrideSource["tu"], "WnPIiEckrrMC9K6z"; got != want {
		t.Fatalf("override-source payload[77].tu = %q, want %q", got, want)
	}
	if got, want := slot77OverrideSource["Xd"], "84ixof8NxS9L6Bba3eXd5Q=="; got != want {
		t.Fatalf("override-source payload[77].Xd = %q, want %q", got, want)
	}
	if got, want := slot77OverrideSource["BgV"], ""; got != want {
		t.Fatalf("override-source payload[77].BgV = %q, want %q", got, want)
	}
	slot99, ok := payload[99].(map[string]any)
	if !ok {
		t.Fatalf("payload[99] type = %T, want map[string]any", payload[99])
	}
	if got, want := slot99["n"], ""; got != want {
		t.Fatalf("payload[99].n = %q, want %q", got, want)
	}
	if got, want := slot99["cOw"], ""; got != want {
		t.Fatalf("payload[99].cOw = %q, want %q", got, want)
	}
	slot99Probably := encodeCurrentCastleObjectSlot99FromCanPlayType("probably")
	if got, want := slot99Probably["n"], "xwZQrlpzYPkEdsCAOFtie66z"; got != want {
		t.Fatalf("probably payload[99].n = %q, want %q", got, want)
	}
	if got, want := slot99Probably["cOw"], "vx4WQxnF0Mtie4ybxwZ6zgR2vx56zgy5O90z1TPV"; got != want {
		t.Fatalf("probably payload[99].cOw = %q, want %q", got, want)
	}
	slot99OverrideSource := encodeCurrentCastleObjectSlot99FromSource("override")
	if got, want := slot99OverrideSource["n"], "xwaMm6NsZf1l/SuSGcWjbA=="; got != want {
		t.Fatalf("override-source payload[99].n = %q, want %q", got, want)
	}
	if got, want := slot99OverrideSource["cOw"], ""; got != want {
		t.Fatalf("override-source payload[99].cOw = %q, want %q", got, want)
	}
	slot100, ok := payload[100].(map[string]any)
	if !ok {
		t.Fatalf("payload[100] type = %T, want map[string]any", payload[100])
	}
	if got, want := slot100["rj"], "5e4bQb/W"; got != want {
		t.Fatalf("payload[100].rj = %q, want %q", got, want)
	}
	slot100KwJ, ok := slot100["KwJ"].(map[string]any)
	if !ok {
		t.Fatalf("payload[100].KwJ type = %T, want map[string]any", slot100["KwJ"])
	}
	if got, want := slot100KwJ["T"], "CzejbOYagxFvwg=="; got != want {
		t.Fatalf("payload[100].KwJ.T = %q, want %q", got, want)
	}
	slot100L, ok := slot100["L"].(map[string]any)
	if !ok {
		t.Fatalf("payload[100].L type = %T, want map[string]any", slot100["L"])
	}
	if got, want := slot100L["T"], ""; got != want {
		t.Fatalf("payload[100].L.T = %q, want %q", got, want)
	}
	if got, want := slot100["E"], ""; got != want {
		t.Fatalf("payload[100].E = %q, want %q", got, want)
	}
	slot100Probably, err := encodeCurrentCastleObjectSlot100FromSlot99(slot99Probably)
	if err != nil {
		t.Fatalf("encodeCurrentCastleObjectSlot100FromSlot99(probably) error = %v", err)
	}
	if got, want := slot100Probably["rj"], "Ls8="; got != want {
		t.Fatalf("probably payload[100].rj = %q, want %q", got, want)
	}
	slot100ProbablyKwJ := slot100Probably["KwJ"].(map[string]any)
	if got, want := slot100ProbablyKwJ["T"], "LRR6zgy5es7352/CgxE="; got != want {
		t.Fatalf("probably payload[100].KwJ.T = %q, want %q", got, want)
	}
	slot100OverrideSource, err := encodeCurrentCastleObjectSlot100FromSlot99(slot99OverrideSource)
	if err != nil {
		t.Fatalf("encodeCurrentCastleObjectSlot100FromSlot99(override source) error = %v", err)
	}
	if got, want := slot100OverrideSource["rj"], "U6Y="; got != want {
		t.Fatalf("override-source payload[100].rj = %q, want %q", got, want)
	}
	slot100OverrideSourceKwJ := slot100OverrideSource["KwJ"].(map[string]any)
	if got, want := slot100OverrideSourceKwJ["T"], "5hrmGi0UkB0LN/fnkB0="; got != want {
		t.Fatalf("override-source payload[100].KwJ.T = %q, want %q", got, want)
	}
	slot46, ok := payload[46].(map[string]any)
	if !ok {
		t.Fatalf("payload[46] type = %T, want map[string]any", payload[46])
	}
	if got, want := slot46["mx"], ""; got != want {
		t.Fatalf("payload[46].mx = %q, want %q", got, want)
	}
	if got, want := slot46["h"], ""; got != want {
		t.Fatalf("payload[46].h = %q, want %q", got, want)
	}
	if got, want := slot46["r"], ""; got != want {
		t.Fatalf("payload[46].r = %q, want %q", got, want)
	}
	if got, want := slot46["AI"], ""; got != want {
		t.Fatalf("payload[46].AI = %q, want %q", got, want)
	}
	slot46Matches := encodeCurrentCastleObjectSlot46FromMatches([]string{"1", "2"})
	if got, want := slot46Matches["mx"], "iZfmGsWEiBWjbBTBAvSMmw=="; got != want {
		t.Fatalf("matches payload[46].mx = %q, want %q", got, want)
	}
	if got, want := slot46Matches["h"], "k5JfO+Xu3eU="; got != want {
		t.Fatalf("matches payload[46].h = %q, want %q", got, want)
	}
	if got, want := slot46Matches["r"], ""; got != want {
		t.Fatalf("matches payload[46].r = %q, want %q", got, want)
	}
	if got, want := slot46Matches["AI"], ""; got != want {
		t.Fatalf("matches payload[46].AI = %q, want %q", got, want)
	}
	slot128, ok := payload[128].(map[string]any)
	if !ok {
		t.Fatalf("payload[128] type = %T, want map[string]any", payload[128])
	}
	if got, want := slot128["ekC"], "jM4d0BtBVjEbQeXusaE="; got != want {
		t.Fatalf("payload[128].ekC = %q, want %q", got, want)
	}
	if got, want := slot128["fCH"], "8Dc="; got != want {
		t.Fatalf("payload[128].fCH = %q, want %q", got, want)
	}
	slot173, ok := payload[173].(map[string]any)
	if !ok {
		t.Fatalf("payload[173] type = %T, want map[string]any", payload[173])
	}
	if got, want := slot173["t"], "gY/mGnrO"; got != want {
		t.Fatalf("payload[173].t = %q, want %q", got, want)
	}
	if got, want := slot173["JZX"], "rImUhLGh8DdWMQ=="; got != want {
		t.Fatalf("payload[173].JZX = %q, want %q", got, want)
	}
	slot173Ga, ok := slot173["ga"].(map[string]any)
	if !ok {
		t.Fatalf("payload[173].ga type = %T, want map[string]any", slot173["ga"])
	}
	if got, want := slot173Ga["sE"], ""; got != want {
		t.Fatalf("payload[173].ga.sE = %q, want %q", got, want)
	}
	if got, want := slot173Ga["OkL"], ""; got != want {
		t.Fatalf("payload[173].ga.OkL = %q, want %q", got, want)
	}
	if got, want := slot173Ga["hD"], ""; got != want {
		t.Fatalf("payload[173].ga.hD = %q, want %q", got, want)
	}
	if got, want := slot173["wPg"], ""; got != want {
		t.Fatalf("payload[173].wPg = %q, want %q", got, want)
	}
	if got, want := slot173["FqX"], ""; got != want {
		t.Fatalf("payload[173].FqX = %q, want %q", got, want)
	}
	slot214, ok := payload[214].(map[string]any)
	if !ok {
		t.Fatalf("payload[214] type = %T, want map[string]any", payload[214])
	}
	if got, want := slot214["h"], "LRTQyw=="; got != want {
		t.Fatalf("payload[214].h = %q, want %q", got, want)
	}
	if got, want := slot214["yhj"], "S+hTpr/W8DdL6PA3"; got != want {
		t.Fatalf("payload[214].yhj = %q, want %q", got, want)
	}
	if got, want := slot214["XI"], ""; got != want {
		t.Fatalf("payload[214].XI = %q, want %q", got, want)
	}
	slot290, ok := payload[290].(map[string]any)
	if !ok {
		t.Fatalf("payload[290] type = %T, want map[string]any", payload[290])
	}
	if got, want := slot290["d"], "LRRvwgs3+WmjbOYaCzcMuQ=="; got != want {
		t.Fatalf("payload[290].d = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleObjectSlot290ValueForTest(payload), "LRRvwgs3+WmjbOYaCzcMuQ=="; got != want {
		t.Fatalf("encodeCurrentCastleObjectSlot290ValueForTest() = %q, want %q", got, want)
	}
	slot342, ok := payload[342].(map[string]any)
	if !ok {
		t.Fatalf("payload[342] type = %T, want map[string]any", payload[342])
	}
	if got, want := slot342["kFx"], "CzfQy+YaCzdac9DLLRSQHQ=="; got != want {
		t.Fatalf("payload[342].kFx = %q, want %q", got, want)
	}
	slot352, ok := payload[352].(map[string]any)
	if !ok {
		t.Fatalf("payload[352] type = %T, want map[string]any", payload[352])
	}
	if got, want := slot352["XYw"], "gY/mGnrO0MsLN5Ad9+ejbA=="; got != want {
		t.Fatalf("payload[352].XYw = %q, want %q", got, want)
	}
	slot416, ok := payload[416].(map[string]any)
	if !ok {
		t.Fatalf("payload[416] type = %T, want map[string]any", payload[416])
	}
	if got, want := slot416["y"], "LRR6zlpzo2xac/lpkB335w=="; got != want {
		t.Fatalf("payload[416].y = %q, want %q", got, want)
	}
	if got, want := slot416["OiJ"], ""; got != want {
		t.Fatalf("payload[416].OiJ = %q, want %q", got, want)
	}
	slot416HRd, ok := slot416["HRd"].(map[string]any)
	if !ok {
		t.Fatalf("payload[416].HRd type = %T, want map[string]any", slot416["HRd"])
	}
	if got, want := slot416HRd["j"], ""; got != want {
		t.Fatalf("payload[416].HRd.j = %q, want %q", got, want)
	}
	if got, want := slot416HRd["hCP"], ""; got != want {
		t.Fatalf("payload[416].HRd.hCP = %q, want %q", got, want)
	}
	slot424, ok := payload[424].(map[string]any)
	if !ok {
		t.Fatalf("payload[424] type = %T, want map[string]any", payload[424])
	}
	if got, want := slot424["kr"], "jM4d0BtBVjE="; got != want {
		t.Fatalf("payload[424].kr = %q, want %q", got, want)
	}
	slot424Vj, ok := slot424["vj"].(map[string]any)
	if !ok {
		t.Fatalf("payload[424].vj type = %T, want map[string]any", slot424["vj"])
	}
	if got, want := slot424Vj["GXV"], "G0Hl7rGh8Dc="; got != want {
		t.Fatalf("payload[424].vj.GXV = %q, want %q", got, want)
	}
	if got, want := slot424["RK"], ""; got != want {
		t.Fatalf("payload[424].RK = %q, want %q", got, want)
	}
	slot424Wxd, ok := slot424["wXd"].(map[string]any)
	if !ok {
		t.Fatalf("payload[424].wXd type = %T, want map[string]any", slot424["wXd"])
	}
	slot424RtN, ok := slot424Wxd["RtN"].(map[string]any)
	if !ok {
		t.Fatalf("payload[424].wXd.RtN type = %T, want map[string]any", slot424Wxd["RtN"])
	}
	if got, want := slot424RtN["o"], ""; got != want {
		t.Fatalf("payload[424].wXd.RtN.o = %q, want %q", got, want)
	}
	if got, want := slot424RtN["ApQ"], ""; got != want {
		t.Fatalf("payload[424].wXd.RtN.ApQ = %q, want %q", got, want)
	}
	if got, want := slot424RtN["G"], ""; got != want {
		t.Fatalf("payload[424].wXd.RtN.G = %q, want %q", got, want)
	}
	if err := populateCurrentCastleNumericSlots(payload, map[int]uint32{289: 0xdae88a03}); err != nil {
		t.Fatalf("populateCurrentCastleNumericSlots() sentinel override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after override error = %v", err)
	}
	slot290 = payload[290].(map[string]any)
	if got, want := slot290["d"], "LRRvwgs3+WmjbOYaCzcMuQ=="; got != want {
		t.Fatalf("sentinel payload[290].d = %q, want %q", got, want)
	}
	if err := populateCurrentCastleNumericSlots(payload, map[int]uint32{341: 1}); err != nil {
		t.Fatalf("populateCurrentCastleNumericSlots() slot 341 override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after slot 341 override error = %v", err)
	}
	slot342 = payload[342].(map[string]any)
	if got, want := slot342["kFx"], "CzfQy+YaCzdac9DLLRSQHQ=="; got != want {
		t.Fatalf("override payload[342].kFx = %q, want %q", got, want)
	}
	if err := populateCurrentCastleStringSlots(payload, map[int]string{351: "override"}); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() slot 351 override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after string override error = %v", err)
	}
	slot352 = payload[352].(map[string]any)
	if got, want := slot352["XYw"], "gY/mGnrO0MsLN5Ad9+ejbA=="; got != want {
		t.Fatalf("override payload[352].XYw = %q, want %q", got, want)
	}
	if err := populateCurrentCastleStringSlots(payload, map[int]string{213: "override"}); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() slot 213 override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after slot 213 override error = %v", err)
	}
	slot214 = payload[214].(map[string]any)
	if got, want := slot214["h"], "LRTQyw=="; got != want {
		t.Fatalf("override payload[214].h = %q, want %q", got, want)
	}
	if got, want := slot214["yhj"], "S+hTpr/W8DdL6PA3"; got != want {
		t.Fatalf("override payload[214].yhj = %q, want %q", got, want)
	}
	if got, want := slot214["XI"], ""; got != want {
		t.Fatalf("override payload[214].XI = %q, want %q", got, want)
	}
	if err := populateCurrentCastleStringSlots(payload, map[int]string{127: "override"}); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() slot 127 override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after slot 127 override error = %v", err)
	}
	slot128 = payload[128].(map[string]any)
	if got, want := slot128["ekC"], "jM4d0BtBVjEbQeXusaE="; got != want {
		t.Fatalf("override payload[128].ekC = %q, want %q", got, want)
	}
	if got, want := slot128["fCH"], "8Dc="; got != want {
		t.Fatalf("override payload[128].fCH = %q, want %q", got, want)
	}
	if err := populateCurrentCastleStringSlots(payload, map[int]string{172: "override"}); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() slot 172 override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after slot 172 override error = %v", err)
	}
	slot173 = payload[173].(map[string]any)
	if got, want := slot173["t"], "gY/mGnrO"; got != want {
		t.Fatalf("override payload[173].t = %q, want %q", got, want)
	}
	if got, want := slot173["JZX"], "rImUhLGh8DdWMQ=="; got != want {
		t.Fatalf("override payload[173].JZX = %q, want %q", got, want)
	}
	slot173Ga = slot173["ga"].(map[string]any)
	if got, want := slot173Ga["sE"], ""; got != want {
		t.Fatalf("override payload[173].ga.sE = %q, want %q", got, want)
	}
	if got, want := slot173Ga["OkL"], ""; got != want {
		t.Fatalf("override payload[173].ga.OkL = %q, want %q", got, want)
	}
	if got, want := slot173Ga["hD"], ""; got != want {
		t.Fatalf("override payload[173].ga.hD = %q, want %q", got, want)
	}
	if got, want := slot173["wPg"], ""; got != want {
		t.Fatalf("override payload[173].wPg = %q, want %q", got, want)
	}
	if got, want := slot173["FqX"], ""; got != want {
		t.Fatalf("override payload[173].FqX = %q, want %q", got, want)
	}
	if err := populateCurrentCastleStringSlots(payload, map[int]string{415: "override"}); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() slot 415 override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after slot 415 override error = %v", err)
	}
	slot416 = payload[416].(map[string]any)
	if got, want := slot416["y"], "LRR6zlpzo2xac/lpkB335w=="; got != want {
		t.Fatalf("override payload[416].y = %q, want %q", got, want)
	}
	if got, want := slot416["OiJ"], ""; got != want {
		t.Fatalf("override payload[416].OiJ = %q, want %q", got, want)
	}
	slot416HRd = slot416["HRd"].(map[string]any)
	if got, want := slot416HRd["j"], ""; got != want {
		t.Fatalf("override payload[416].HRd.j = %q, want %q", got, want)
	}
	if got, want := slot416HRd["hCP"], ""; got != want {
		t.Fatalf("override payload[416].HRd.hCP = %q, want %q", got, want)
	}
	if err := populateCurrentCastleStringSlots(payload, map[int]string{423: "override"}); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() slot 423 override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after slot 423 override error = %v", err)
	}
	slot424 = payload[424].(map[string]any)
	if got, want := slot424["kr"], "jM4d0BtBVjE="; got != want {
		t.Fatalf("override payload[424].kr = %q, want %q", got, want)
	}
	slot424Vj = slot424["vj"].(map[string]any)
	if got, want := slot424Vj["GXV"], "G0Hl7rGh8Dc="; got != want {
		t.Fatalf("override payload[424].vj.GXV = %q, want %q", got, want)
	}
	if got, want := slot424["RK"], ""; got != want {
		t.Fatalf("override payload[424].RK = %q, want %q", got, want)
	}
	slot424Wxd = slot424["wXd"].(map[string]any)
	slot424RtN = slot424Wxd["RtN"].(map[string]any)
	if got, want := slot424RtN["o"], ""; got != want {
		t.Fatalf("override payload[424].wXd.RtN.o = %q, want %q", got, want)
	}
	if got, want := slot424RtN["ApQ"], ""; got != want {
		t.Fatalf("override payload[424].wXd.RtN.ApQ = %q, want %q", got, want)
	}
	if got, want := slot424RtN["G"], ""; got != want {
		t.Fatalf("override payload[424].wXd.RtN.G = %q, want %q", got, want)
	}
	if err := populateCurrentCastleStringSlots(payload, map[int]string{40: "override"}); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() slot 40 override error = %v", err)
	}
	if err := populateCurrentCastleObjectSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() after slot 40 override error = %v", err)
	}
	slot41 = payload[41].(map[string]any)
	if got, want := slot41["V"], "Ls9BGR3QrImUhLGh8DdWMQ=="; got != want {
		t.Fatalf("override payload[41].V = %q, want %q", got, want)
	}
	if got, want := slot41["wh"], ""; got != want {
		t.Fatalf("override payload[41].wh = %q, want %q", got, want)
	}
	if got, want := slot41["y"], ""; got != want {
		t.Fatalf("override payload[41].y = %q, want %q", got, want)
	}
	if got, want := slot41["G"], ""; got != want {
		t.Fatalf("override payload[41].G = %q, want %q", got, want)
	}
	slot41Hrc = slot41["Hrc"].(map[string]any)
	if got, want := slot41Hrc["j"], ""; got != want {
		t.Fatalf("override payload[41].Hrc.j = %q, want %q", got, want)
	}
	if err := populateCurrentCastleObjectSlots(payload, map[int]map[string]any{
		41:  {"V": "override41"},
		77:  {"AFS": "override77"},
		99:  {"n": "override99"},
		100: {"rj": "override100"},
		46:  {"mx": "override46"},
		128: {"ekC": "override128"},
		173: {"t": "override173"},
		214: {"h": "override214"},
		290: {"d": "override"},
		342: {"kFx": "override342"},
		352: {"XYw": "override352"},
		416: {"y": "override416"},
		424: {"kr": "override424"},
	}); err != nil {
		t.Fatalf("populateCurrentCastleObjectSlots() override error = %v", err)
	}
	slot41 = payload[41].(map[string]any)
	if got, want := slot41["V"], "override41"; got != want {
		t.Fatalf("override payload[41].V = %q, want %q", got, want)
	}
	slot77 = payload[77].(map[string]any)
	if got, want := slot77["AFS"], "override77"; got != want {
		t.Fatalf("override payload[77].AFS = %q, want %q", got, want)
	}
	slot99 = payload[99].(map[string]any)
	if got, want := slot99["n"], "override99"; got != want {
		t.Fatalf("override payload[99].n = %q, want %q", got, want)
	}
	slot100 = payload[100].(map[string]any)
	if got, want := slot100["rj"], "override100"; got != want {
		t.Fatalf("override payload[100].rj = %q, want %q", got, want)
	}
	slot46 = payload[46].(map[string]any)
	if got, want := slot46["mx"], "override46"; got != want {
		t.Fatalf("override payload[46].mx = %q, want %q", got, want)
	}
	slot128 = payload[128].(map[string]any)
	if got, want := slot128["ekC"], "override128"; got != want {
		t.Fatalf("override payload[128].ekC = %q, want %q", got, want)
	}
	slot173 = payload[173].(map[string]any)
	if got, want := slot173["t"], "override173"; got != want {
		t.Fatalf("override payload[173].t = %q, want %q", got, want)
	}
	slot214 = payload[214].(map[string]any)
	if got, want := slot214["h"], "override214"; got != want {
		t.Fatalf("override payload[214].h = %q, want %q", got, want)
	}
	slot290 = payload[290].(map[string]any)
	if got, want := slot290["d"], "override"; got != want {
		t.Fatalf("override payload[290].d = %q, want %q", got, want)
	}
	slot342 = payload[342].(map[string]any)
	if got, want := slot342["kFx"], "override342"; got != want {
		t.Fatalf("override payload[342].kFx = %q, want %q", got, want)
	}
	slot352 = payload[352].(map[string]any)
	if got, want := slot352["XYw"], "override352"; got != want {
		t.Fatalf("override payload[352].XYw = %q, want %q", got, want)
	}
	slot416 = payload[416].(map[string]any)
	if got, want := slot416["y"], "override416"; got != want {
		t.Fatalf("override payload[416].y = %q, want %q", got, want)
	}
	slot424 = payload[424].(map[string]any)
	if got, want := slot424["kr"], "override424"; got != want {
		t.Fatalf("override payload[424].kr = %q, want %q", got, want)
	}
	if err := populateCurrentCastleObjectSlots(make([]any, currentCastlePayloadSlots-1), nil); err == nil {
		t.Fatalf("populateCurrentCastleObjectSlots() with short payload succeeded")
	}
}

func encodeCurrentCastleObjectSlot290ValueForTest(payload []any) string {
	encoded, err := encodeCurrentCastleObjectSlot290(payload)
	if err != nil {
		return "ERROR: " + err.Error()
	}
	value, _ := encoded["d"].(string)
	return value
}

func TestCurrentCastleLowerFloatEncoderSixthBatch(t *testing.T) {
	cases := []struct {
		slot int
		want string
	}{
		{slot: 169, want: "kNqJtno023Y="},
		{slot: 170, want: "wjo4fSk/u54="},
		{slot: 171, want: "Bw3Aapka7RI="},
		{slot: 178, want: "lugV26QrSCM="},
		{slot: 179, want: "nvAp5xT3kN8="},
		{slot: 181, want: "I+6oFcmOm/g="},
		{slot: 185, want: "yfI99uK0cdc="},
		{slot: 186, want: "Bra+Pl6or+g="},
		{slot: 187, want: "vM2SNbNt/WE="},
		{slot: 189, want: "51GaJEYYWQo="},
		{slot: 191, want: "jjVVwgTmLfw="},
		{slot: 193, want: "Hjvh3DH0K9A="},
		{slot: 195, want: "+pxbCViZ/HE="},
		{slot: 197, want: "7FYto4qzNis="},
		{slot: 199, want: "HVla/tDdGI8="},
		{slot: 200, want: "OsCZA8zz4Cs="},
		{slot: 201, want: "NExCiV1Lz+g="},
		{slot: 202, want: "uTE3fMA+s10="},
		{slot: 208, want: "+e4VYAvk7SY="},
		{slot: 210, want: "QI6lR6o37s8="},
		{slot: 211, want: "HN+v+Vfny+I="},
		{slot: 216, want: "Aj7xzHYVp7w="},
		{slot: 219, want: "c1HKtEVksdw="},
		{slot: 222, want: "IznkXo0O2TY="},
		{slot: 225, want: "mAlFU+jTDJM="},
		{slot: 226, want: "ohn2/Iu8m9w="},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.slot), func(t *testing.T) {
			encoder := currentCastleLowerFloatSlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing encoder for slot %d", tc.slot)
			}
			if got := encoder(123.456); got != tc.want {
				t.Fatalf("slot %d sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}
}

func TestCurrentCastleLowerFloatEncoderSeventhBatch(t *testing.T) {
	cases := []struct {
		slot int
		want string
	}{
		{slot: 227, want: "gN+3yhvT4U0="},
		{slot: 229, want: "oCOWa+lDkz8="},
		{slot: 230, want: "XH6fDUKdHsU="},
		{slot: 232, want: "vJLiyWbWNNA="},
		{slot: 233, want: "EJdX64bHr/0="},
		{slot: 234, want: "ulzfjZAdPHU="},
		{slot: 236, want: "OtTxp4j39H8="},
		{slot: 238, want: "CIfnTJZwr34="},
		{slot: 241, want: "D34+l5zfTvM="},
		{slot: 242, want: "5xKC3Mo6/t0="},
		{slot: 243, want: "pLrtRzY3Wj8="},
		{slot: 244, want: "NKJoeoP6pzo="},
		{slot: 248, want: "ti+S+364r1s="},
		{slot: 250, want: "BOXd8qH563c="},
		{slot: 251, want: "wITFYgxBg9I="},
		{slot: 252, want: "aUvOdA+E67w="},
		{slot: 253, want: "GOIXHfbtAjU="},
		{slot: 255, want: "cpQVY1jT9Ls="},
		{slot: 259, want: "TReYPlPut/Y="},
		{slot: 262, want: "uJ0tM2U1MSI="},
		{slot: 263, want: "eJmhjt2Flws="},
		{slot: 265, want: "lDYL8aIhFtk="},
		{slot: 268, want: "/ErazXMOdDo="},
		{slot: 269, want: "sEhKD1tNyew="},
		{slot: 270, want: "7ZWPXiLkmEE="},
		{slot: 272, want: "2s3BUivXzhQ="},
		{slot: 278, want: "0DEJJnUtP6M="},
		{slot: 280, want: "QjIuuRA8Nfs="},
		{slot: 284, want: "qIPjHok0HyM="},
		{slot: 285, want: "oFbzLSL9drU="},
		{slot: 286, want: "m4UG9MFEZaw="},
		{slot: 287, want: "vLo3aXb52pE="},
		{slot: 294, want: "dKJFf/ZPQoc="},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.slot), func(t *testing.T) {
			encoder := currentCastleLowerFloatSlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing encoder for slot %d", tc.slot)
			}
			if got := encoder(123.456); got != tc.want {
				t.Fatalf("slot %d sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}
}

func TestCurrentCastleLowerFloatEncoderEighthBatch(t *testing.T) {
	cases := []struct {
		slot int
		want string
	}{
		{slot: 295, want: "P0sxhaQmCxY="},
		{slot: 296, want: "VW86452i7oI="},
		{slot: 299, want: "I32YVhlG3W4="},
		{slot: 301, want: "GHm+V8txfVo="},
		{slot: 303, want: "ycEfdODWEmU="},
		{slot: 305, want: "NBaXhVLV9k0="},
		{slot: 306, want: "VpPjtRurh6w="},
		{slot: 307, want: "ygl5J4ExFSw="},
		{slot: 309, want: "vT9OXLOsnwQ="},
		{slot: 310, want: "7o1St8S8jz0="},
		{slot: 313, want: "s23OvClMDWQ="},
		{slot: 315, want: "qyMhZhIopIc="},
		{slot: 317, want: "hg3tURx9NUM="},
		{slot: 319, want: "SdTyth3WFUY="},
		{slot: 320, want: "BJrHWZ7pOoE="},
		{slot: 322, want: "udt86p96e6I="},
		{slot: 323, want: "9BEpqt3Vf9s="},
		{slot: 324, want: "W2NlnnJc4P8="},
		{slot: 325, want: "Wnz7SUC53JE="},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.slot), func(t *testing.T) {
			encoder := currentCastleLowerFloatSlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing encoder for slot %d", tc.slot)
			}
			if got := encoder(123.456); got != tc.want {
				t.Fatalf("slot %d sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}
}

func TestCurrentCastleLowerFloatEncoderSkippedMidBatch(t *testing.T) {
	cases := []struct {
		slot int
		want string
	}{
		{slot: 264, want: "TmgJ/7QPiDc="},
		{slot: 271, want: "XUPAcgeC46o="},
		{slot: 275, want: "CuY/I5bnZsY="},
		{slot: 276, want: "Tzc58iYwtJM="},
		{slot: 282, want: "gLy7XjQ/fe4="},
		{slot: 283, want: "NRRsClgDEoY="},
		{slot: 291, want: "rg8/TelXQ14="},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.slot), func(t *testing.T) {
			encoder := currentCastleLowerFloatSlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing encoder for slot %d", tc.slot)
			}
			if got := encoder(123.456); got != tc.want {
				t.Fatalf("slot %d sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}
}

func TestCurrentCastleLowerFloatEncoderNinthBatch(t *testing.T) {
	cases := []struct {
		slot int
		want string
	}{
		{slot: 327, want: "ahpGDRAsubc="},
		{slot: 331, want: "DsRPybD5pLE="},
		{slot: 332, want: "tT3ccSxNFUM="},
		{slot: 333, want: "ANr/FY7F+r0="},
		{slot: 335, want: "525Ool2+hrQ="},
		{slot: 336, want: "7FKBU/bjcks="},
		{slot: 337, want: "LX8nUxjReBQ="},
		{slot: 338, want: "YVuK46XKa+8="},
		{slot: 340, want: "3ewsVcoNnHk="},
		{slot: 343, want: "iqgpm+xrCEM="},
		{slot: 345, want: "za9IfmvODyY="},
		{slot: 346, want: "70YeEaqSIKQ="},
		{slot: 347, want: "rdMWoJ9Qs8g="},
		{slot: 353, want: "pPKJKw4bUrM="},
		{slot: 355, want: "Fkp6MN0iRik="},
		{slot: 357, want: "IW9YGosqDyI="},
		{slot: 358, want: "98XUPhmuJUY="},
		{slot: 360, want: "m21mpJG0DZw="},
		{slot: 361, want: "8JfuP6W3x4M="},
		{slot: 362, want: "STiYbHNIUNo="},
		{slot: 363, want: "DVz+oc67WKI="},
		{slot: 364, want: "AVtHrQpu2k8="},
		{slot: 365, want: "Cp3VqJFpU2c="},
		{slot: 367, want: "lDzTSCt0JHo="},
		{slot: 368, want: "YdkXPGi+ur0="},
		{slot: 370, want: "o5SYKvKul+0="},
		{slot: 371, want: "7jBde2yr0JM="},
		{slot: 372, want: "txswobMsJlE="},
		{slot: 373, want: "5MZDWZIJplE="},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.slot), func(t *testing.T) {
			encoder := currentCastleLowerFloatSlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing encoder for slot %d", tc.slot)
			}
			if got := encoder(123.456); got != tc.want {
				t.Fatalf("slot %d sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}
}

func TestCurrentCastleLowerFloatEncoderTenthBatch(t *testing.T) {
	cases := []struct {
		slot int
		want string
	}{
		{slot: 377, want: "GB/eh/2/T7M="},
		{slot: 379, want: "D4WY8hnCpfo="},
		{slot: 380, want: "/BKN38ZP8rc="},
		{slot: 381, want: "YGn19Dv8OfA="},
		{slot: 382, want: "UbiYdIsI0MI="},
		{slot: 384, want: "pZsaTF/cu3Q="},
		{slot: 385, want: "QI4lx6q37k8="},
		{slot: 387, want: "EO5N3ypPTic="},
		{slot: 389, want: "SrqmMZi0vfM="},
		{slot: 391, want: "gLy7XjQ/fe4="},
		{slot: 396, want: "m7k+aH34GaA="},
		{slot: 397, want: "Nx2aLEHcveQ="},
		{slot: 398, want: "O6mkvv3uCaY="},
		{slot: 401, want: "cMsmLFlsSQw="},
		{slot: 402, want: "/+KligDCcpY="},
		{slot: 403, want: "RjiTkZyBmKk="},
		{slot: 404, want: "R+EDnPaQ6YI="},
		{slot: 405, want: "SquTvO+3pTk="},
		{slot: 407, want: "20D8lRHUwDU="},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.slot), func(t *testing.T) {
			encoder := currentCastleLowerFloatSlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing encoder for slot %d", tc.slot)
			}
			if got := encoder(123.456); got != tc.want {
				t.Fatalf("slot %d sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}
}

func TestCurrentCastleLowerFloatEncoderEleventhBatch(t *testing.T) {
	cases := []struct {
		slot int
		want string
	}{
		{slot: 408, want: "pLfhbuAN954="},
		{slot: 409, want: "qElxXg1VR9s="},
		{slot: 410, want: "VLWNpvmhvyM="},
		{slot: 411, want: "DTOw4vdyk1o="},
		{slot: 412, want: "8HeXK2YHTz0="},
		{slot: 413, want: "15OSNQsWVMU="},
		{slot: 414, want: "8AaD1cpFZh0="},
		{slot: 417, want: "NIpDDRZdqmU="},
		{slot: 418, want: "VJCPMggTUcI="},
		{slot: 419, want: "4h4dwJah31A="},
		{slot: 420, want: "CNY0a0Rx0mg="},
		{slot: 426, want: "ZFy6j/vxLQA="},
		{slot: 427, want: "x71ooinSnYo="},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.slot), func(t *testing.T) {
			encoder := currentCastleLowerFloatSlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing encoder for slot %d", tc.slot)
			}
			if got := encoder(123.456); got != tc.want {
				t.Fatalf("slot %d sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}
}

func TestCurrentCastleLowerFloatEncoderScatteredBatch(t *testing.T) {
	cases := []struct {
		slot int
		want string
	}{
		{slot: 86, want: "kCmz3FqbqXw="},
		{slot: 92, want: "8ftCGK+Im8A="},
		{slot: 135, want: "572GzjnMqc8="},
		{slot: 161, want: "dPubMOoM0wI="},
		{slot: 190, want: "iOpruaYpCsE="},
		{slot: 192, want: "BVX/KPGqVOg="},
		{slot: 215, want: "fVB4U4z0vvo="},
		{slot: 218, want: "DCqbtX7liu0="},
		{slot: 249, want: "vC/Pa1K/V60="},
		{slot: 300, want: "bqiVq5z7yAM="},
		{slot: 314, want: "/nmZJWgJQTM="},
		{slot: 392, want: "UXMUgjcSEzo="},
	}
	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.slot), func(t *testing.T) {
			encoder := currentCastleLowerFloatSlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing encoder for slot %d", tc.slot)
			}
			if got := encoder(123.456); got != tc.want {
				t.Fatalf("slot %d sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}
}

func TestCurrentCastleArraySlotEncoders(t *testing.T) {
	sample := []float64{0, 1, 123.456}
	cases := []struct {
		name   string
		slot   int
		values []float64
		want   string
	}{
		{name: "slot254 sample", slot: 254, values: sample, want: "Y2NjY2NjY2Oe02NjY2NjY6O5dE5lvlm2"},
		{name: "slot302 sample", slot: 302, values: sample, want: "wMDAwMDAwMBRycDAwMDAwMzlclqvW+/a"},
		{name: "slot254 empty", slot: 254, values: nil, want: ""},
		{name: "slot302 empty", slot: 302, values: nil, want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			encoder := currentCastleArraySlotEncoders[tc.slot]
			if encoder == nil {
				t.Fatalf("missing array encoder for slot %d", tc.slot)
			}
			if got := encoder(tc.values); got != tc.want {
				t.Fatalf("slot %d array sample = %q, want %q", tc.slot, got, tc.want)
			}
		})
	}

	payload := newCurrentCastlePayloadScaffold()
	if err := populateCurrentCastleArraySlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleArraySlots() error = %v", err)
	}
	if got, want := payload[254], encodeCurrentCastleArraySlot254(currentCastleFloat64Samples(256, 199)); got != want {
		t.Fatalf("payload[254] = %q, want %q", got, want)
	}
	if got, want := payload[302], encodeCurrentCastleArraySlot302(currentCastleFloat64Samples(16, 13)); got != want {
		t.Fatalf("payload[302] = %q, want %q", got, want)
	}
	if err := populateCurrentCastleArraySlots(make([]any, currentCastlePayloadSlots-1), nil); err == nil {
		t.Fatalf("populateCurrentCastleArraySlots() with short payload succeeded")
	}
}

func TestCurrentCastleStringSlotEncoders(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  string
	}{
		{name: "empty", value: "", want: ""},
		{name: "ascii", value: "abc", want: "qeqtnO+D"},
		{name: "number", value: "123.456", want: "9FL4BDprU243uXGgdFI="},
		{name: "timezone", value: "America/Chicago", want: "j1otnCU4vWSQTu+DqeqtVc1zVmeQTu+Dqepr0VO1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := encodeCurrentCastleUTF16Scramble(tc.value); got != tc.want {
				t.Fatalf("encodeCurrentCastleUTF16Scramble(%q) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
	if got, want := encodeCurrentCastleDoubleUTF16Scramble("123.456"), "/DYPWrElN7lzDMlBf0u9ZBRw+AQ3uTprNoBGn0lBa9Hr0Q9adz54BA=="; got != want {
		t.Fatalf("encodeCurrentCastleDoubleUTF16Scramble() = %q, want %q", got, want)
	}

	payload := newCurrentCastlePayloadScaffold()
	if err := populateCurrentCastleStringSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() error = %v", err)
	}
	expectedDefaults := map[int]string{
		11:  "isz",
		40:  "NoAtnP9LeAQ=",
		89:  "NoAtnP9LeAQ=",
		127: "NoAtnP9LeAQ=",
		172: "NoAtnP9LeAQ=",
		175: "NoAtnP9LeAQ=",
		213: "GraL",
		217: "NoAtnP9LeAQ=",
		231: "NoAtnP9LeAQ=",
		351: "NoAtnP9LeAQ=",
		400: "NoAtnP9LeAQ=",
		415: "NoAtnP9LeAQ=",
		423: "NoAtnP9LeAQ=",
	}
	for slot, want := range expectedDefaults {
		if got := payload[slot]; got != want {
			t.Fatalf("payload[%d] = %q, want %q", slot, got, want)
		}
	}
	if err := populateCurrentCastleStringSlots(payload, map[int]string{89: "123.456"}); err != nil {
		t.Fatalf("populateCurrentCastleStringSlots() with override error = %v", err)
	}
	if got, want := payload[89], "9FL4BDprU243uXGgdFI="; got != want {
		t.Fatalf("payload[89] override = %q, want %q", got, want)
	}
	if err := populateCurrentCastleStringSlots(make([]any, currentCastlePayloadSlots-1), nil); err == nil {
		t.Fatalf("populateCurrentCastleStringSlots() with short payload succeeded")
	}
}

func TestCurrentCastlePackedStringSlotEncoders(t *testing.T) {
	const slot0MissingProbe = "\u216c\u2166\u216c\u2129\u2097\u598c"
	slot0Cases := []struct {
		name string
		fn   func(string) string
		want string
	}{
		{name: "primary", fn: encodeCurrentCastleSlot0Primary, want: "Wrf7WrdxWrf7WvW2WnIXn3HD"},
		{name: "secondary", fn: encodeCurrentCastleSlot0Secondary, want: "bvpIbvorbvpIbsC5bmxS/CtC"},
		{name: "tertiary", fn: encodeCurrentCastleSlot0Tertiary, want: "HnddHneyHnddHmD4HjK3g7JZ"},
		{name: "quaternary", fn: encodeCurrentCastleSlot0Quaternary, want: "ROQFROTMROQFRAB5RFiu8MwB"},
		{name: "quinary", fn: encodeCurrentCastleSlot0Quinary, want: "A5UGA5VHA5UGA6U2A8UY00fI"},
		{name: "senary", fn: encodeCurrentCastleSlot0Senary, want: "8JO68JO08JO68JK38JCl87Sa"},
	}
	for _, tc := range slot0Cases {
		t.Run("slot0 "+tc.name, func(t *testing.T) {
			if got := tc.fn(slot0MissingProbe); got != tc.want {
				t.Fatalf("slot 0 %s encoder = %q, want %q", tc.name, got, tc.want)
			}
		})
	}
	if got, want := encodeCurrentCastleSlot0Primary("abc"), "uHo8"; got != want {
		t.Fatalf("encodeCurrentCastleSlot0Primary(sample) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot0Secondary("abc"), "PGaQ"; got != want {
		t.Fatalf("encodeCurrentCastleSlot0Secondary(sample) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot0Quaternary("abc"), "cFQ4"; got != want {
		t.Fatalf("encodeCurrentCastleSlot0Quaternary(sample) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot0Quinary("abc"), "Gwv7"; got != want {
		t.Fatalf("encodeCurrentCastleSlot0Quinary(sample) = %q, want %q", got, want)
	}

	slot7Cases := []struct {
		name string
		fn   func(string) string
		want string
	}{
		{name: "canvas primary", fn: encodeCurrentCastleSlot7CanvasPrimary, want: "K6zp"},
		{name: "worker timing", fn: encodeCurrentCastleSlot7WorkerTiming, want: "JNn2"},
		{name: "media ratio", fn: encodeCurrentCastleSlot7MediaRatio, want: "p1vL"},
		{name: "plugin state", fn: encodeCurrentCastleSlot7PluginState, want: "1qXk"},
		{name: "frame ratio", fn: encodeCurrentCastleSlot7FrameRatio, want: "32it"},
		{name: "worker canvas", fn: encodeCurrentCastleSlot7WorkerCanvas, want: "pM6W"},
		{name: "worker feature", fn: encodeCurrentCastleSlot7WorkerFeature, want: "Acgz"},
		{name: "navigator probe", fn: encodeCurrentCastleSlot7NavigatorProbe, want: "2qQO"},
		{name: "viewport ratio", fn: encodeCurrentCastleSlot7ViewportRatio, want: "jOua"},
		{name: "worker signal", fn: encodeCurrentCastleSlot7WorkerSignal, want: "0zoY"},
		{name: "pointer ratio", fn: encodeCurrentCastleSlot7PointerRatio, want: "KV8X"},
	}
	for _, tc := range slot7Cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.fn("w:1"); got != tc.want {
				t.Fatalf("slot 7 %s encoder = %q, want %q", tc.name, got, tc.want)
			}
		})
	}
	if got, want := encodeCurrentCastleSlot7CanvasPrimary("abc"), "+WTP"; got != want {
		t.Fatalf("encodeCurrentCastleSlot7CanvasPrimary(sample) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot7WorkerTiming("abc"), "RtEw"; got != want {
		t.Fatalf("encodeCurrentCastleSlot7WorkerTiming(sample) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot7WorkerSignal("abc"), "FjFR"; got != want {
		t.Fatalf("encodeCurrentCastleSlot7WorkerSignal(sample) = %q, want %q", got, want)
	}

	if got, want := encodeCurrentCastleWebGPUVendor("g:1"), "3eT7"; got != want {
		t.Fatalf("encodeCurrentCastleWebGPUVendor() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleWebGPULimits("g:1"), "ElQu"; got != want {
		t.Fatalf("encodeCurrentCastleWebGPULimits() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleWebGPUArchitecture("g:1"), "Q4Su"; got != want {
		t.Fatalf("encodeCurrentCastleWebGPUArchitecture() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastlePackedStringComponents(nil), "qj4="; got != want {
		t.Fatalf("encodeCurrentCastlePackedStringComponents(nil) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastlePackedStringComponents([]string{"a", "bb"}), "To4Q9k6OHdBL6Evo"; got != want {
		t.Fatalf("encodeCurrentCastlePackedStringComponents(sample) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastlePackedStringComponents([]string{"", "", ""}), "tUeqPqo+qj4="; got != want {
		t.Fatalf("encodeCurrentCastlePackedStringComponents(empty3) = %q, want %q", got, want)
	}

	payload := newCurrentCastlePayloadScaffold()
	if err := populateCurrentCastlePackedStringSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastlePackedStringSlots() error = %v", err)
	}
	if got, want := payload[0].(string), 192; len(got) != want {
		t.Fatalf("payload[0] len = %d, want %d", len(got), want)
	}
	if got, want := payload[7].(string), 1560; len(got) != want {
		t.Fatalf("payload[7] len = %d, want %d", len(got), want)
	}
	if got, want := payload[119].(string), 96; len(got) != want {
		t.Fatalf("payload[119] len = %d, want %d", len(got), want)
	}
	if err := populateCurrentCastlePackedStringSlots(payload, map[int][]string{0: {"a", "bb"}, 7: {"a", "bb"}, 119: {"a", "bb"}}); err != nil {
		t.Fatalf("populateCurrentCastlePackedStringSlots() with override error = %v", err)
	}
	if got, want := payload[0], "To4Q9k6OHdBL6Evo"; got != want {
		t.Fatalf("payload[0] override = %q, want %q", got, want)
	}
	if got, want := payload[7], "To4Q9k6OHdBL6Evo"; got != want {
		t.Fatalf("payload[7] override = %q, want %q", got, want)
	}
	if got, want := payload[119], "To4Q9k6OHdBL6Evo"; got != want {
		t.Fatalf("payload[119] override = %q, want %q", got, want)
	}
	if err := populateCurrentCastlePackedStringSlots(make([]any, currentCastlePayloadSlots-1), nil); err == nil {
		t.Fatalf("populateCurrentCastlePackedStringSlots() with short payload succeeded")
	}
}

func TestCurrentCastleUnitPackedStringSlotEncoders(t *testing.T) {
	if got, want := encodeCurrentCastleSlot8Primary("s:1"), "ASz5"; got != want {
		t.Fatalf("encodeCurrentCastleSlot8Primary() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot8Secondary("s:1"), "SKnA"; got != want {
		t.Fatalf("encodeCurrentCastleSlot8Secondary() = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot8Primary("abc"), "6wur"; got != want {
		t.Fatalf("encodeCurrentCastleSlot8Primary(sample) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleSlot8Secondary("abc"), "dpa2"; got != want {
		t.Fatalf("encodeCurrentCastleSlot8Secondary(sample) = %q, want %q", got, want)
	}
	if got, want := encodeCurrentCastleUnitPackedStringComponents([]string{"a", "bb"}), "RaK7+kWies4EdgR2"; got != want {
		t.Fatalf("encodeCurrentCastleUnitPackedStringComponents(sample) = %q, want %q", got, want)
	}

	payload := newCurrentCastlePayloadScaffold()
	if err := populateCurrentCastleUnitPackedStringSlots(payload, nil); err != nil {
		t.Fatalf("populateCurrentCastleUnitPackedStringSlots() error = %v", err)
	}
	if got, want := payload[8].(string), 52; len(got) != want {
		t.Fatalf("payload[8] len = %d, want %d", len(got), want)
	}
	if err := populateCurrentCastleUnitPackedStringSlots(payload, map[int][]string{8: {"a", "bb"}}); err != nil {
		t.Fatalf("populateCurrentCastleUnitPackedStringSlots() with override error = %v", err)
	}
	if got, want := payload[8], "RaK7+kWies4EdgR2"; got != want {
		t.Fatalf("payload[8] override = %q, want %q", got, want)
	}
	if err := populateCurrentCastleUnitPackedStringSlots(make([]any, currentCastlePayloadSlots-1), nil); err == nil {
		t.Fatalf("populateCurrentCastleUnitPackedStringSlots() with short payload succeeded")
	}
}

func TestCurrentCastlePayloadScaffold(t *testing.T) {
	payload := newCurrentCastlePayloadScaffold()
	if len(payload) != currentCastlePayloadSlots {
		t.Fatalf("newCurrentCastlePayloadScaffold() len = %d, want %d", len(payload), currentCastlePayloadSlots)
	}
	for slot, value := range payload {
		if value == nil {
			t.Fatalf("payload[%d] is nil", slot)
		}
	}
	if err := populateCurrentCastleAutomationSlot(payload, 0x1353003b); err != nil {
		t.Fatalf("populateCurrentCastleAutomationSlot() error = %v", err)
	}
	if err := populateCurrentCastleHighTimingSlots(payload, map[int]float64{
		431: 0,
		480: 123.456,
		493: 1,
	}); err != nil {
		t.Fatalf("populateCurrentCastleHighTimingSlots() error = %v", err)
	}
	payloadJSON, err := marshalCurrentCastlePayload(payload)
	if err != nil {
		t.Fatalf("marshalCurrentCastlePayload() error = %v", err)
	}
	if strings.Contains(payloadJSON, "null") {
		t.Fatalf("payload JSON contains null: %s", payloadJSON)
	}
	var decoded []any
	if err = json.Unmarshal([]byte(payloadJSON), &decoded); err != nil {
		t.Fatalf("payload JSON did not unmarshal: %v", err)
	}
	if len(decoded) != currentCastlePayloadSlots {
		t.Fatalf("decoded payload len = %d, want %d", len(decoded), currentCastlePayloadSlots)
	}
	if got, want := uint32(decoded[2].(float64)), uint32(2115052381); got != want {
		t.Fatalf("decoded payload[2] = %d, want %d", got, want)
	}
	if got, want := decoded[431], "xsbGxsbGxsY="; got != want {
		t.Fatalf("decoded payload[431] = %q, want %q", got, want)
	}
	if got, want := decoded[480], "hhjbKVyZ+PE="; got != want {
		t.Fatalf("decoded payload[480] = %q, want %q", got, want)
	}
	token, err := createCurrentCastleWrappedPayloadToken(payload, 1760000000123)
	if err != nil {
		t.Fatalf("createCurrentCastleWrappedPayloadToken() error = %v", err)
	}
	wire, err := decodeCurrentCastleWrappedTokenForTest(token)
	if err != nil {
		t.Fatalf("decodeCurrentCastleWrappedTokenForTest() error = %v", err)
	}
	encodedTimestampLen := int(wire[0])
	decodedPayload := xorCurrentCastleString(wire[1+encodedTimestampLen:], "1760000000123")
	if strings.Contains(decodedPayload, "null") {
		t.Fatalf("wrapped payload contains null: %s", decodedPayload)
	}
	if _, err = marshalCurrentCastlePayload(payload[:currentCastlePayloadSlots-1]); err == nil {
		t.Fatalf("marshalCurrentCastlePayload() with short payload succeeded")
	}
	payload[7] = nil
	if _, err = marshalCurrentCastlePayload(payload); err == nil {
		t.Fatalf("marshalCurrentCastlePayload() with nil slot succeeded")
	}
	if err = populateCurrentCastleAutomationSlot(make([]any, 2), 0); err == nil {
		t.Fatalf("populateCurrentCastleAutomationSlot() with short payload succeeded")
	}
}

func TestCurrentCastlePayloadScaffoldAcceptedTypes(t *testing.T) {
	payload := newCurrentCastlePayloadScaffold()
	typeCounts := map[string]int{}
	numericIndexes := make([]int, 0)
	objectIndexes := make([]int, 0)
	booleanIndexes := make([]int, 0)
	for index, value := range payload {
		switch value.(type) {
		case int, int64, uint32, float64:
			typeCounts["number"]++
			numericIndexes = append(numericIndexes, index)
		case map[string]any:
			typeCounts["object"]++
			objectIndexes = append(objectIndexes, index)
		case bool:
			typeCounts["boolean"]++
			booleanIndexes = append(booleanIndexes, index)
		case string:
			typeCounts["string"]++
		default:
			t.Fatalf("payload[%d] has unexpected type %T", index, value)
		}
	}
	if got, want := typeCounts["string"], 355; got != want {
		t.Fatalf("string slot count = %d, want %d", got, want)
	}
	if got, want := typeCounts["number"], 125; got != want {
		t.Fatalf("number slot count = %d, want %d", got, want)
	}
	if got, want := typeCounts["object"], 13; got != want {
		t.Fatalf("object slot count = %d, want %d", got, want)
	}
	if got, want := typeCounts["boolean"], 1; got != want {
		t.Fatalf("boolean slot count = %d, want %d", got, want)
	}
	if !sameInts(numericIndexes, currentCastleAcceptedNumericSlots) {
		t.Fatalf("numeric slots = %#v, want %#v", numericIndexes, currentCastleAcceptedNumericSlots)
	}
	if !sameInts(objectIndexes, currentCastleAcceptedObjectSlots) {
		t.Fatalf("object slots = %#v, want %#v", objectIndexes, currentCastleAcceptedObjectSlots)
	}
	if !sameInts(booleanIndexes, []int{currentCastleAcceptedBooleanSlot}) {
		t.Fatalf("boolean slots = %#v, want [%d]", booleanIndexes, currentCastleAcceptedBooleanSlot)
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal(payload) error = %v", err)
	}
	if strings.Contains(string(data), "null") {
		t.Fatalf("payload JSON contains null: %s", string(data))
	}
	var decoded []any
	if err = json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal(payload) error = %v", err)
	}
	decodedCounts := map[string]int{}
	for _, value := range decoded {
		switch value.(type) {
		case float64:
			decodedCounts["number"]++
		case map[string]any:
			decodedCounts["object"]++
		case bool:
			decodedCounts["boolean"]++
		case string:
			decodedCounts["string"]++
		default:
			t.Fatalf("decoded payload has unexpected type %T", value)
		}
	}
	for key, want := range typeCounts {
		if got := decodedCounts[key]; got != want {
			t.Fatalf("decoded %s count = %d, want %d", key, got, want)
		}
	}
}

func TestBuildCurrentCastlePayload(t *testing.T) {
	components := make([]string, currentCastleSlot6ComponentCount)
	components[0] = encodeCurrentCastleSlot6Component0("ua")
	components[1] = encodeCurrentCastleSlot6Component1("probe")
	components[2] = encodeCurrentCastleSlot6Component2("en-US")
	components[6] = encodeCurrentCastleSlot6TU("window")
	components[10] = encodeCurrentCastleSlot6AH(currentCastleSlot6HashHex("function test() { return 1; }"))
	components[18] = encodeCurrentCastleSlot6EJ(currentCastleSlot6HashHex("canvas"))
	components[39] = encodeCurrentCastleSlot6A8(currentCastleSlot6HashHex("document"))
	expectedSlot6, err := encodeCurrentCastleSlot6(components)
	if err != nil {
		t.Fatalf("encodeCurrentCastleSlot6() error = %v", err)
	}
	payload, err := buildCurrentCastlePayload(currentCastlePayloadInput{
		AutomationBitfield: 0x1353003b,
		Slot6Components:    components,
		NumericValues: map[int]uint32{
			20:  1,
			26:  1,
			27:  1,
			68:  1,
			98:  1,
			152: 1,
			154: 1,
			159: 1,
			182: 1,
			203: 1,
			256: 1,
			257: currentCastleTB(1),
			258: 1,
			289: currentCastleTB(1),
			341: 1,
			348: 1,
			374: 1,
			386: 1,
			406: 1,
			428: 1,
			429: 1,
		},
		LowerFloatValues: map[int]float64{
			1: 1,
			3: 123.456,
		},
		HighTimingValues: map[int]float64{
			431: 0,
			480: 123.456,
			493: 1,
		},
	})
	if err != nil {
		t.Fatalf("buildCurrentCastlePayload() error = %v", err)
	}
	if len(payload) != currentCastlePayloadSlots {
		t.Fatalf("buildCurrentCastlePayload() len = %d, want %d", len(payload), currentCastlePayloadSlots)
	}
	if got, want := payload[2], uint32(2115052381); got != want {
		t.Fatalf("payload[2] = %v, want %v", got, want)
	}
	if got, want := payload[20], uint32(0x5fd1c84a); got != want {
		t.Fatalf("payload[20] = %v, want %v", got, want)
	}
	if got, want := payload[26], uint32(0x670022bf); got != want {
		t.Fatalf("payload[26] = %v, want %v", got, want)
	}
	if got, want := payload[68], uint32(0x76dc5125); got != want {
		t.Fatalf("payload[68] = %v, want %v", got, want)
	}
	if got, want := payload[98], uint32(0xb0a64635); got != want {
		t.Fatalf("payload[98] = %v, want %v", got, want)
	}
	if got, want := payload[154], uint32(0xf13ce181); got != want {
		t.Fatalf("payload[154] = %v, want %v", got, want)
	}
	if got, want := payload[77].(map[string]any)["AFS"], "JcW/1uh5mugs0/qG3eXd5Q=="; got != want {
		t.Fatalf("payload[77].AFS = %v, want %v", got, want)
	}
	if got, want := payload[99].(map[string]any)["n"], ""; got != want {
		t.Fatalf("payload[99].n = %v, want %v", got, want)
	}
	if got, want := payload[99].(map[string]any)["cOw"], ""; got != want {
		t.Fatalf("payload[99].cOw = %v, want %v", got, want)
	}
	if got, want := payload[100].(map[string]any)["rj"], "5e4bQb/W"; got != want {
		t.Fatalf("payload[100].rj = %v, want %v", got, want)
	}
	if got, want := payload[100].(map[string]any)["KwJ"].(map[string]any)["T"], "CzejbOYagxFvwg=="; got != want {
		t.Fatalf("payload[100].KwJ.T = %v, want %v", got, want)
	}
	if got, want := payload[46].(map[string]any)["mx"], ""; got != want {
		t.Fatalf("payload[46].mx = %v, want %v", got, want)
	}
	if got, want := payload[128].(map[string]any)["ekC"], "jM4d0BtBVjEbQeXusaE="; got != want {
		t.Fatalf("payload[128].ekC = %v, want %v", got, want)
	}
	if got, want := payload[128].(map[string]any)["fCH"], "8Dc="; got != want {
		t.Fatalf("payload[128].fCH = %v, want %v", got, want)
	}
	if got, want := payload[173].(map[string]any)["t"], "gY/mGnrO"; got != want {
		t.Fatalf("payload[173].t = %v, want %v", got, want)
	}
	if got, want := payload[173].(map[string]any)["JZX"], "rImUhLGh8DdWMQ=="; got != want {
		t.Fatalf("payload[173].JZX = %v, want %v", got, want)
	}
	if got, want := payload[256], uint32(0x4ea7245b); got != want {
		t.Fatalf("payload[256] = %v, want %v", got, want)
	}
	if got, want := payload[257], uint32(0x3e5af065); got != want {
		t.Fatalf("payload[257] = %v, want %v", got, want)
	}
	if got, want := payload[289], uint32(0x00021d69); got != want {
		t.Fatalf("payload[289] = %v, want %v", got, want)
	}
	if got, want := payload[214].(map[string]any)["h"], "LRTQyw=="; got != want {
		t.Fatalf("payload[214].h = %v, want %v", got, want)
	}
	if got, want := payload[214].(map[string]any)["yhj"], "S+hTpr/W8DdL6PA3"; got != want {
		t.Fatalf("payload[214].yhj = %v, want %v", got, want)
	}
	if got, want := payload[214].(map[string]any)["XI"], ""; got != want {
		t.Fatalf("payload[214].XI = %v, want %v", got, want)
	}
	if got, want := payload[290].(map[string]any)["d"], "LRRvwgs3+WmjbOYaCzcMuQ=="; got != want {
		t.Fatalf("payload[290].d = %v, want %v", got, want)
	}
	if got, want := payload[348], uint32(0x4ff13582); got != want {
		t.Fatalf("payload[348] = %v, want %v", got, want)
	}
	if got, want := payload[352].(map[string]any)["XYw"], "gY/mGnrO0MsLN5Ad9+ejbA=="; got != want {
		t.Fatalf("payload[352].XYw = %v, want %v", got, want)
	}
	if got, want := payload[27], uint32(0x0dd8798d); got != want {
		t.Fatalf("payload[27] = %v, want %v", got, want)
	}
	if got, want := payload[152], uint32(0x0160726c); got != want {
		t.Fatalf("payload[152] = %v, want %v", got, want)
	}
	if got, want := payload[159], uint32(1); got != want {
		t.Fatalf("payload[159] = %v, want %v", got, want)
	}
	if got, want := payload[182], uint32(0x000003df); got != want {
		t.Fatalf("payload[182] = %v, want %v", got, want)
	}
	if got, want := payload[203], uint32(0xc14d992a); got != want {
		t.Fatalf("payload[203] = %v, want %v", got, want)
	}
	if got, want := payload[258], uint32(0xceb0120a); got != want {
		t.Fatalf("payload[258] = %v, want %v", got, want)
	}
	if got, want := payload[341], uint32(0x0000793c); got != want {
		t.Fatalf("payload[341] = %v, want %v", got, want)
	}
	if got, want := payload[342].(map[string]any)["kFx"], "CzfQy+YaCzdac9DLLRSQHQ=="; got != want {
		t.Fatalf("payload[342].kFx = %v, want %v", got, want)
	}
	if got, want := payload[374], uint32(0x5ed52a3b); got != want {
		t.Fatalf("payload[374] = %v, want %v", got, want)
	}
	if got, want := payload[386], uint32(0x78b37579); got != want {
		t.Fatalf("payload[386] = %v, want %v", got, want)
	}
	if got, want := payload[406], uint32(0x4da110bf); got != want {
		t.Fatalf("payload[406] = %v, want %v", got, want)
	}
	if got, want := payload[416].(map[string]any)["y"], "LRR6zlpzo2xac/lpkB335w=="; got != want {
		t.Fatalf("payload[416].y = %v, want %v", got, want)
	}
	if got, want := payload[424].(map[string]any)["kr"], "jM4d0BtBVjE="; got != want {
		t.Fatalf("payload[424].kr = %v, want %v", got, want)
	}
	if got, want := payload[428], uint32(0xc0f29076); got != want {
		t.Fatalf("payload[428] = %v, want %v", got, want)
	}
	if got, want := payload[429], uint32(0x57fefa07); got != want {
		t.Fatalf("payload[429] = %v, want %v", got, want)
	}
	if got := payload[6]; got != expectedSlot6 {
		t.Fatalf("payload[6] = %q, want %q", got, expectedSlot6)
	}
	if got, want := payload[1], "KWfg4ODg4OA="; got != want {
		t.Fatalf("payload[1] = %q, want %q", got, want)
	}
	if got, want := payload[3], "agkBHM1lE+c="; got != want {
		t.Fatalf("payload[3] = %q, want %q", got, want)
	}
	if got, want := payload[4], "7Ozs7Ozs7Ow="; got != want {
		t.Fatalf("payload[4] = %q, want %q", got, want)
	}
	if got, want := payload[0].(string), 192; len(got) != want {
		t.Fatalf("payload[0] len = %d, want %d", len(got), want)
	}
	if got, want := payload[7].(string), 1560; len(got) != want {
		t.Fatalf("payload[7] len = %d, want %d", len(got), want)
	}
	if got, want := payload[8].(string), 52; len(got) != want {
		t.Fatalf("payload[8] len = %d, want %d", len(got), want)
	}
	if got, want := payload[254], encodeCurrentCastleArraySlot254(currentCastleFloat64Samples(256, 199)); got != want {
		t.Fatalf("payload[254] = %q, want %q", got, want)
	}
	if got, want := payload[302], encodeCurrentCastleArraySlot302(currentCastleFloat64Samples(16, 13)); got != want {
		t.Fatalf("payload[302] = %q, want %q", got, want)
	}
	if got, want := payload[119].(string), 96; len(got) != want {
		t.Fatalf("payload[119] len = %d, want %d", len(got), want)
	}
	if got, want := payload[431], "xsbGxsbGxsY="; got != want {
		t.Fatalf("payload[431] = %q, want %q", got, want)
	}
	if got, want := payload[480], "hhjbKVyZ+PE="; got != want {
		t.Fatalf("payload[480] = %q, want %q", got, want)
	}
	if got, want := payload[493], "OTbGxsbGxsY="; got != want {
		t.Fatalf("payload[493] = %q, want %q", got, want)
	}
	payloadJSON, err := marshalCurrentCastlePayload(payload)
	if err != nil {
		t.Fatalf("marshalCurrentCastlePayload() error = %v", err)
	}
	if strings.Contains(payloadJSON, "null") {
		t.Fatalf("payload JSON contains null: %s", payloadJSON)
	}
	var decoded []any
	if err = json.Unmarshal([]byte(payloadJSON), &decoded); err != nil {
		t.Fatalf("payload JSON did not unmarshal: %v", err)
	}
	if got := decoded[6]; got != expectedSlot6 {
		t.Fatalf("decoded payload[6] = %q, want %q", got, expectedSlot6)
	}
	token, err := createCurrentCastleWrappedPayloadToken(payload, 1760000000123)
	if err != nil {
		t.Fatalf("createCurrentCastleWrappedPayloadToken() error = %v", err)
	}
	wire, err := decodeCurrentCastleWrappedTokenForTest(token)
	if err != nil {
		t.Fatalf("decodeCurrentCastleWrappedTokenForTest() error = %v", err)
	}
	encodedTimestampLen := int(wire[0])
	decodedPayload := xorCurrentCastleString(wire[1+encodedTimestampLen:], "1760000000123")
	if !strings.Contains(decodedPayload, expectedSlot6) {
		t.Fatalf("wrapped payload missing slot 6 value")
	}
	if strings.Contains(decodedPayload, "null") {
		t.Fatalf("wrapped payload contains null: %s", decodedPayload)
	}
	if _, err = buildCurrentCastlePayload(currentCastlePayloadInput{
		AutomationBitfield: 0x1353003b,
		Slot6Components:    components[:len(components)-1],
	}); err == nil {
		t.Fatalf("buildCurrentCastlePayload() with short slot 6 components succeeded")
	}
}

func TestCurrentCastleLengthPrefixUnits(t *testing.T) {
	cases := []struct {
		name   string
		input  uint32
		output []uint16
	}{
		{name: "zero", input: 0, output: []uint16{0}},
		{name: "component count", input: 47, output: []uint16{47}},
		{name: "single max", input: 0x7fff, output: []uint16{0x7fff}},
		{name: "first continuation", input: 0x8000, output: []uint16{0x8000, 0x0001}},
		{name: "continuation with payload bits", input: 0x8001, output: []uint16{0x8001, 0x0001}},
		{name: "three units", input: 0x40000000, output: []uint16{0x8000, 0x8000, 0x0001}},
		{name: "uint32 max", input: 0xffffffff, output: []uint16{0xffff, 0xffff, 0x0003}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := currentCastleLengthPrefixUnits(tc.input); !sameUint16s(got, tc.output) {
				t.Fatalf("currentCastleLengthPrefixUnits(%#x) = %#v, want %#v", tc.input, got, tc.output)
			}
		})
	}
}

func TestCurrentCastleSlot6LengthPrefixUnits(t *testing.T) {
	lengths := make([]int, currentCastleSlot6ComponentCount)
	for i := range lengths {
		lengths[i] = i
	}
	lengths[1] = 0x8000
	units, err := currentCastleSlot6LengthPrefixUnits(lengths)
	if err != nil {
		t.Fatalf("currentCastleSlot6LengthPrefixUnits() error = %v", err)
	}
	wantPrefix := []uint16{47, 0, 0x8000, 1, 2, 3, 4}
	if !sameUint16s(units[:len(wantPrefix)], wantPrefix) {
		t.Fatalf("slot 6 prefix start = %#v, want %#v", units[:len(wantPrefix)], wantPrefix)
	}
	if _, err = currentCastleSlot6LengthPrefixUnits(lengths[:len(lengths)-1]); err == nil {
		t.Fatalf("currentCastleSlot6LengthPrefixUnits() with short component list succeeded")
	}
	lengths[7] = -1
	if _, err = currentCastleSlot6LengthPrefixUnits(lengths); err == nil {
		t.Fatalf("currentCastleSlot6LengthPrefixUnits() with negative component length succeeded")
	}
}

func TestCurrentCastleSlot6RawUnits(t *testing.T) {
	components := make([]string, currentCastleSlot6ComponentCount)
	components[0] = "A"
	components[1] = "BC"
	units, err := currentCastleSlot6RawUnits(components)
	if err != nil {
		t.Fatalf("currentCastleSlot6RawUnits() error = %v", err)
	}
	wantPrefix := []uint16{47, 1, 2, 0, 0}
	if !sameUint16s(units[:len(wantPrefix)], wantPrefix) {
		t.Fatalf("slot 6 raw prefix = %#v, want %#v", units[:len(wantPrefix)], wantPrefix)
	}
	wantSuffix := []uint16{'A', 'B', 'C'}
	if !sameUint16s(units[len(units)-len(wantSuffix):], wantSuffix) {
		t.Fatalf("slot 6 raw suffix = %#v, want %#v", units[len(units)-len(wantSuffix):], wantSuffix)
	}
	if len(units) != currentCastleSlot6ComponentCount+1+len(wantSuffix) {
		t.Fatalf("slot 6 raw unit len = %d, want %d", len(units), currentCastleSlot6ComponentCount+1+len(wantSuffix))
	}
	encoded, err := encodeCurrentCastleSlot6(components)
	if err != nil {
		t.Fatalf("encodeCurrentCastleSlot6() error = %v", err)
	}
	const wantEncoded = "RyS7+kWiMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTMFMwUzBTO93FhE8s"
	if encoded != wantEncoded {
		t.Fatalf("encodeCurrentCastleSlot6() = %q, want %q", encoded, wantEncoded)
	}
	if _, err = currentCastleSlot6RawUnits(components[:len(components)-1]); err == nil {
		t.Fatalf("currentCastleSlot6RawUnits() with short components succeeded")
	}
	if _, err = encodeCurrentCastleSlot6(components[:len(components)-1]); err == nil {
		t.Fatalf("encodeCurrentCastleSlot6() with short components succeeded")
	}
}

func TestEncodeCurrentCastleSlot6ComponentValues(t *testing.T) {
	values := make([]string, currentCastleSlot6ComponentCount)
	values[0] = "toString"
	values[1] = "TypeError: Cyclic __proto__ value"
	values[2] = "America/Chicago"
	values[3] = "ANGLE (NVIDIA, NVIDIA GeForce GTX 1080 Ti (0x00001B06) Direct3D11 vs_5_0 ps_5_0, D3D11)"
	values[5] = currentCastleSlot6HashHex("function test() { return 1; }")
	values[6] = "landscape-primary"
	values[8] = "RangeError"
	values[12] = "probably"
	values[13] = "maybe"
	values[14] = "1111111111111111111111111111111111111000000000"
	values[15] = "Cannot read properties of undefined (reading 'b')"
	values[19] = "091a1f32-3826-4bad-9250-aa14e3c0a2b2"
	values[20] = "1000"
	values[21] = "Google Inc."
	values[22] = values[3]
	values[25] = "{}"
	values[29] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/149.0.0.0 Safari/537.36"
	values[31] = "en-US"
	values[34] = "en-US"
	values[36] = "Win32"
	values[41] = "x.com"
	values[44] = "en-US,en"
	values[45] = "probably"
	components, err := encodeCurrentCastleSlot6ComponentValues(values)
	if err != nil {
		t.Fatalf("encodeCurrentCastleSlot6ComponentValues() error = %v", err)
	}
	if len(components) != currentCastleSlot6ComponentCount {
		t.Fatalf("components len = %d, want %d", len(components), currentCastleSlot6ComponentCount)
	}
	cases := map[int]string{
		0:  encodeCurrentCastleSlot6Component0(values[0]),
		1:  encodeCurrentCastleSlot6Component1(values[1]),
		2:  encodeCurrentCastleSlot6Component2(values[2]),
		3:  encodeCurrentCastleSlot6UZ(values[3]),
		5:  encodeCurrentCastleSlot6Component5Fallback(values[5]),
		6:  encodeCurrentCastleSlot6TU(values[6]),
		8:  encodeCurrentCastleSlot6FJ(values[8]),
		12: encodeCurrentCastleSlot6AG(values[12]),
		13: encodeCurrentCastleSlot6AD(values[13]),
		14: encodeCurrentCastleSlot6T7(values[14]),
		15: encodeCurrentCastleSlot6ID(values[15]),
		19: encodeCurrentCastleSlot6Component19(values[19]),
		20: encodeCurrentCastleSlot6L9(values[20]),
		21: encodeCurrentCastleSlot6TX(values[21]),
		22: encodeCurrentCastleSlot6LI(values[22]),
		25: encodeCurrentCastleSlot6F3(values[25]),
		29: encodeCurrentCastleSlot6TUpperF(values[29]),
		31: encodeCurrentCastleSlot6Component31(values[31]),
		34: encodeCurrentCastleSlot6AQ(values[34]),
		36: encodeCurrentCastleSlot6TUpperM(values[36]),
		41: encodeCurrentCastleSlot6EA(values[41]),
		44: encodeCurrentCastleSlot6TUpperX(values[44]),
		45: encodeCurrentCastleSlot6L5(values[45]),
	}
	for index, want := range cases {
		if got := components[index]; got != want {
			t.Fatalf("components[%d] = %q, want %q", index, got, want)
		}
	}
	slot6, err := encodeCurrentCastleSlot6(components)
	if err != nil {
		t.Fatalf("encodeCurrentCastleSlot6() error = %v", err)
	}
	units, err := currentCastleSlot6RawUnits(components)
	if err != nil {
		t.Fatalf("currentCastleSlot6RawUnits() error = %v", err)
	}
	if slot6 != encodeCurrentCastleUnits(units) {
		t.Fatalf("slot 6 encoding mismatch")
	}
	if _, err = encodeCurrentCastleSlot6ComponentValues(values[:len(values)-1]); err == nil {
		t.Fatalf("encodeCurrentCastleSlot6ComponentValues() with short values succeeded")
	}
}

func TestBuildCurrentCastleSlot6RawValues(t *testing.T) {
	fp := currentCastleSlot6Fingerprint{
		TimestampMillis: 1782425262310,
		Timezone:        "America/Chicago",
		WebGLRenderer:   "ANGLE (NVIDIA, NVIDIA GeForce GTX 1080 Ti (0x00001B06) Direct3D11 vs_5_0 ps_5_0, D3D11)",
		WebGLVendor:     "Google Inc.",
		Orientation:     "landscape-primary",
		UserAgent:       UserAgent,
		Language:        "en-US",
		Languages:       "en-US,en",
		Platform:        "Win32",
		Host:            "x.com",
		PublicKey:       castlePublicKey,
		ClientUUID:      "091a1f32-3826-4bad-9250-aa14e3c0a2b2",
		FeatureBits:     currentCastleDefaultFeatureBits,
		PrecisionProbe:  "1000",
		AudioProbe:      "20030107",
		ProbeDateString: "3/3/1970, 6:00:00 PM",
		NumericProbe:    "51833740944",
		WASMProbeHex:    "c0ffee",
		Hash5:           "39a01d2b",
		Hash7:           "bd4be70d",
		Hash10:          "45eb4280",
		Hash11:          "8fc90f7d",
		Hash16:          "5d55117f",
		Hash18:          "42d68d60",
		Hash23:          "34b454a7",
		Hash39:          "db77b610",
	}
	values := buildCurrentCastleSlot6RawValues(fp)
	if len(values) != currentCastleSlot6ComponentCount {
		t.Fatalf("raw values len = %d, want %d", len(values), currentCastleSlot6ComponentCount)
	}
	cases := map[int]string{
		0:  "toString",
		1:  "TypeError: Cyclic __proto__ value",
		2:  fp.Timezone,
		3:  fp.WebGLRenderer,
		5:  fp.Hash5,
		6:  fp.Orientation,
		7:  fp.Hash7,
		8:  "RangeError",
		9:  "[]",
		10: fp.Hash10,
		11: fp.Hash11,
		12: "probably",
		13: "maybe",
		14: fp.FeatureBits,
		15: "Cannot read properties of undefined (reading 'b')",
		16: fp.Hash16,
		18: fp.Hash18,
		19: fp.ClientUUID,
		20: fp.PrecisionProbe,
		21: fp.WebGLVendor,
		22: fp.WebGLRenderer,
		23: fp.Hash23,
		24: fp.AudioProbe,
		25: "{}",
		27: fp.PublicKey,
		28: "Illegal invocation",
		29: fp.UserAgent,
		30: "r:1",
		31: fp.Language,
		33: strconv.FormatInt(currentCastleSlot6ComponentTimestampMillis(fp.TimestampMillis), 10),
		34: fp.Language,
		35: "Maximum call stack size exceeded",
		36: fp.Platform,
		37: fp.ProbeDateString,
		38: "r:1",
		39: fp.Hash39,
		41: fp.Host,
		43: fp.WASMProbeHex,
		44: fp.Languages,
		45: "probably",
		46: fp.NumericProbe,
	}
	for index, want := range cases {
		if got := values[index]; got != want {
			t.Fatalf("values[%d] = %q, want %q", index, got, want)
		}
	}
	for _, index := range []int{4, 17, 26, 32, 40, 42} {
		if values[index] != "" {
			t.Fatalf("values[%d] = %q, want empty", index, values[index])
		}
	}

	components, err := buildCurrentCastleSlot6Components(fp)
	if err != nil {
		t.Fatalf("buildCurrentCastleSlot6Components() error = %v", err)
	}
	if got, want := components[2], encodeCurrentCastleSlot6Component2(fp.Timezone); got != want {
		t.Fatalf("components[2] = %q, want %q", got, want)
	}
	if got, want := components[19], encodeCurrentCastleSlot6Component19(fp.ClientUUID); got != want {
		t.Fatalf("components[19] = %q, want %q", got, want)
	}
	if got, want := components[27], encodeCurrentCastleSlot6Component27(fp.PublicKey); got != want {
		t.Fatalf("components[27] = %q, want %q", got, want)
	}
}

func TestDefaultCurrentCastleSlot6Fingerprint(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_TIMEZONE", "America/Chicago")
	const timestampMillis = int64(1782425262310)
	const clientUUID = "091a1f32-3826-4bad-9250-aa14e3c0a2b2"

	fp := defaultCurrentCastleSlot6Fingerprint(timestampMillis, clientUUID)
	values := buildCurrentCastleSlot6RawValues(fp)
	if got, want := values[2], "America/Chicago"; got != want {
		t.Fatalf("values[2] = %q, want %q", got, want)
	}
	if got, want := values[19], clientUUID; got != want {
		t.Fatalf("values[19] = %q, want %q", got, want)
	}
	if got, want := values[21], currentCastleDefaultWebGLVendor; got != want {
		t.Fatalf("values[21] = %q, want %q", got, want)
	}
	if got, want := values[22], currentCastleDefaultWebGLRenderer; got != want {
		t.Fatalf("values[22] = %q, want %q", got, want)
	}
	if got, want := values[27], castlePublicKey; got != want {
		t.Fatalf("values[27] = %q, want %q", got, want)
	}
	if got, want := values[29], UserAgent; got != want {
		t.Fatalf("values[29] = %q, want %q", got, want)
	}
	if got, want := values[33], strconv.FormatInt(currentCastleSlot6ComponentTimestampMillis(timestampMillis), 10); got != want {
		t.Fatalf("values[33] = %q, want %q", got, want)
	}
	if got, want := values[37], "3/3/1970, 6:00:00 PM"; got != want {
		t.Fatalf("values[37] = %q, want %q", got, want)
	}
	exactDefaults := map[int]string{
		5:  currentCastleDefaultHash5,
		7:  currentCastleDefaultHash7,
		10: currentCastleDefaultHash10,
		11: currentCastleDefaultHash11,
		16: currentCastleDefaultHash16,
		18: currentCastleDefaultHash18,
		23: currentCastleDefaultHash23,
		39: currentCastleDefaultHash39,
		43: currentCastleDefaultWASMProbeHex,
		46: currentCastleDefaultNumericProbe,
	}
	for index, want := range exactDefaults {
		if got := values[index]; got != want {
			t.Fatalf("values[%d] = %q, want %q", index, got, want)
		}
	}
	components, err := buildCurrentCastleSlot6Components(fp)
	if err != nil {
		t.Fatalf("buildCurrentCastleSlot6Components() error = %v", err)
	}
	payload, err := buildCurrentCastlePayload(currentCastlePayloadInput{
		AutomationBitfield: currentCastleAcceptedAutomationBitfield,
		Slot6Components:    components,
	})
	if err != nil {
		t.Fatalf("buildCurrentCastlePayload() error = %v", err)
	}
	if got, want := payload[2], uint32(2115052381); got != want {
		t.Fatalf("payload[2] = %v, want %v", got, want)
	}
	if payload[6] == "" {
		t.Fatalf("payload[6] is empty")
	}
	token, err := createCurrentCastleWrappedPayloadToken(payload, timestampMillis)
	if err != nil {
		t.Fatalf("createCurrentCastleWrappedPayloadToken() error = %v", err)
	}
	if !strings.HasPrefix(token, currentCastleTokenPrefix) {
		t.Fatalf("token missing current Castle prefix")
	}
}

func TestCreateCurrentCastleRequestToken(t *testing.T) {
	t.Setenv("TWITTER_JETFUEL_TIMEZONE", "America/Chicago")
	token, err := createCurrentCastleRequestToken("091a1f32-3826-4bad-9250-aa14e3c0a2b2")
	if err != nil {
		t.Fatalf("createCurrentCastleRequestToken() error = %v", err)
	}
	if !strings.HasPrefix(token, currentCastleTokenPrefix) {
		t.Fatalf("token missing current Castle prefix")
	}
	t.Logf("current Castle token length: %d", len(token))
	if len(token) < 1000 {
		t.Fatalf("token length = %d, want a current-format request token", len(token))
	}
	if _, err = decodeCurrentCastleWrappedTokenForTest(token); err != nil {
		t.Fatalf("decodeCurrentCastleWrappedTokenForTest() error = %v", err)
	}
}

func TestCompareChromeCurrentCastleTokenProbe(t *testing.T) {
	token := strings.TrimSpace(os.Getenv("TWITTER_CASTLE_TOKEN_DIFF_PROBE"))
	if token == "" {
		t.Skip("TWITTER_CASTLE_TOKEN_DIFF_PROBE is required")
	}
	chromePayload, timestampMillis := decodeCurrentCastlePayloadFromTokenForTest(t, token)
	components, err := buildCurrentCastleSlot6Components(defaultCurrentCastleSlot6Fingerprint(timestampMillis, "091a1f32-3826-4bad-9250-aa14e3c0a2b2"))
	if err != nil {
		t.Fatalf("buildCurrentCastleSlot6Components() error = %v", err)
	}
	goPayload, err := buildCurrentCastlePayload(currentCastlePayloadInput{
		AutomationBitfield: currentCastleAcceptedAutomationBitfield,
		Slot6Components:    components,
	})
	if err != nil {
		t.Fatalf("buildCurrentCastlePayload() error = %v", err)
	}
	if len(chromePayload) != len(goPayload) {
		t.Fatalf("payload lengths: chrome=%d go=%d", len(chromePayload), len(goPayload))
	}
	typeCounts := map[string]int{}
	mismatches := make([]string, 0)
	for i := range chromePayload {
		typeCounts[fmt.Sprintf("chrome_%s", currentCastleJSONTypeForTest(chromePayload[i]))]++
		typeCounts[fmt.Sprintf("go_%s", currentCastleJSONTypeForTest(goPayload[i]))]++
		chromeJSON := currentCastleJSONSummaryForTest(t, chromePayload[i])
		goJSON := currentCastleJSONSummaryForTest(t, goPayload[i])
		if chromeJSON != goJSON {
			mismatches = append(mismatches, fmt.Sprintf(
				"%d:%s/%d/%s -> %s/%d/%s",
				i,
				currentCastleJSONTypeForTest(chromePayload[i]),
				len(chromeJSON),
				currentCastleShortHashForTest(chromeJSON),
				currentCastleJSONTypeForTest(goPayload[i]),
				len(goJSON),
				currentCastleShortHashForTest(goJSON),
			))
		}
	}
	sort.Strings(mismatches)
	t.Logf("decoded timestamp millis: %d", timestampMillis)
	t.Logf("type counts: %v", typeCounts)
	t.Logf("mismatch count: %d", len(mismatches))
	limit := min(len(mismatches), 120)
	for _, mismatch := range mismatches[:limit] {
		t.Logf("mismatch %s", mismatch)
	}
	if len(mismatches) > limit {
		t.Logf("mismatch output truncated: %d more", len(mismatches)-limit)
	}
	if slots := strings.TrimSpace(os.Getenv("TWITTER_CASTLE_TOKEN_DIFF_VALUE_SLOTS")); slots != "" {
		for _, part := range strings.Split(slots, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			slot, err := strconv.Atoi(part)
			if err != nil {
				t.Fatalf("invalid slot %q in TWITTER_CASTLE_TOKEN_DIFF_VALUE_SLOTS: %v", part, err)
			}
			if slot < 0 || slot >= len(chromePayload) {
				t.Fatalf("slot %d outside payload length %d", slot, len(chromePayload))
			}
			chromeString, chromeOK := chromePayload[slot].(string)
			goString, goOK := goPayload[slot].(string)
			if !chromeOK || !goOK {
				t.Logf("value %d chrome_type=%s go_type=%s", slot, currentCastleJSONTypeForTest(chromePayload[slot]), currentCastleJSONTypeForTest(goPayload[slot]))
				continue
			}
			t.Logf("value %d chrome=%q go=%q", slot, chromeString, goString)
		}
	}
	if slots := strings.TrimSpace(os.Getenv("TWITTER_CASTLE_TOKEN_DIFF_ARRAY_SLOTS")); slots != "" {
		for _, part := range strings.Split(slots, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			slot, err := strconv.Atoi(part)
			if err != nil {
				t.Fatalf("invalid slot %q in TWITTER_CASTLE_TOKEN_DIFF_ARRAY_SLOTS: %v", part, err)
			}
			chromeString, ok := chromePayload[slot].(string)
			if !ok {
				t.Logf("array %d chrome_type=%s", slot, currentCastleJSONTypeForTest(chromePayload[slot]))
				continue
			}
			values := decodeCurrentCastleFloat64ArraySlotForTest(t, slot, chromeString)
			nonZero := 0
			firstNonZero := make([]string, 0, 12)
			for i, value := range values {
				if value == 0 {
					continue
				}
				nonZero++
				if len(firstNonZero) < 12 {
					firstNonZero = append(firstNonZero, fmt.Sprintf("%d:%g", i, value))
				}
			}
			preview := make([]float64, min(16, len(values)))
			copy(preview, values[:len(preview)])
			t.Logf("array %d len=%d nonzero=%d first=%v first_nonzero=%v", slot, len(values), nonZero, preview, firstNonZero)
		}
	}
	if os.Getenv("TWITTER_CASTLE_TOKEN_DIFF_SLOT6") != "" {
		chromeSlot6, chromeOK := chromePayload[6].(string)
		goSlot6, goOK := goPayload[6].(string)
		if !chromeOK || !goOK {
			t.Logf("slot6 chrome_type=%s go_type=%s", currentCastleJSONTypeForTest(chromePayload[6]), currentCastleJSONTypeForTest(goPayload[6]))
		} else {
			chromeComponents := decodeCurrentCastleSlot6ComponentsForTest(t, chromeSlot6)
			goComponents := decodeCurrentCastleSlot6ComponentsForTest(t, goSlot6)
			componentMismatches := make([]string, 0)
			for i := range chromeComponents {
				chromeJSON := currentCastleJSONSummaryForTest(t, chromeComponents[i])
				goJSON := currentCastleJSONSummaryForTest(t, goComponents[i])
				if chromeJSON != goJSON {
					componentMismatches = append(componentMismatches, fmt.Sprintf(
						"%d:%d/%s -> %d/%s",
						i,
						len(chromeComponents[i]),
						currentCastleShortHashForTest(chromeJSON),
						len(goComponents[i]),
						currentCastleShortHashForTest(goJSON),
					))
				}
			}
			t.Logf("slot6 component mismatch count: %d", len(componentMismatches))
			for _, mismatch := range componentMismatches {
				t.Logf("slot6 component %s", mismatch)
			}
			if os.Getenv("TWITTER_CASTLE_TOKEN_DIFF_SLOT6_RAW") != "" {
				for i := range chromeComponents {
					chromeRaw := decodeCurrentCastleSlot6EncodedComponentForTest(t, i, chromeComponents[i])
					goRaw := decodeCurrentCastleSlot6EncodedComponentForTest(t, i, goComponents[i])
					if chromeRaw != goRaw {
						t.Logf("slot6 raw component %d: chrome=%q go=%q", i, chromeRaw, goRaw)
					}
				}
			}
		}
	}
}

func decodeCurrentCastlePayloadFromTokenForTest(t *testing.T, token string) ([]any, int64) {
	t.Helper()
	wire, err := decodeCurrentCastleWrappedTokenForTest(token)
	if err != nil {
		t.Fatalf("decodeCurrentCastleWrappedTokenForTest() error = %v", err)
	}
	if len(wire) < 2 {
		t.Fatalf("decoded wire too short: %d", len(wire))
	}
	encodedTimestampLen := int(wire[0])
	if len(wire) <= 1+encodedTimestampLen {
		t.Fatalf("decoded wire missing payload: wire_len=%d encoded_timestamp_len=%d", len(wire), encodedTimestampLen)
	}
	timestamp := decodeCurrentCastleTimestampForTest(t, wire[1:1+encodedTimestampLen])
	timestampMillis, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		t.Fatalf("decoded timestamp %q parse error: %v", timestamp, err)
	}
	payloadWithChecksum := xorCurrentCastleString(wire[1+encodedTimestampLen:], timestamp)
	payloadJSON := stripCurrentCastleChecksumForTest(t, payloadWithChecksum)
	var payload []any
	if err = json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		t.Fatalf("json.Unmarshal(decoded payload) error = %v", err)
	}
	return payload, timestampMillis
}

func decodeCurrentCastleTimestampForTest(t *testing.T, encoded string) string {
	t.Helper()
	return string(utf16.Decode(decodeCurrentCastleUnitsForTest(t, encoded)))
}

func decodeCurrentCastleUnitsForTest(t *testing.T, encoded string) []uint16 {
	t.Helper()
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("base64 Castle units decode error = %v", err)
	}
	if len(raw)%2 != 0 {
		t.Fatalf("encoded Castle units decoded to odd byte length %d", len(raw))
	}
	units := make([]uint16, 0, len(raw)/2)
	for i := 0; i < len(raw); i += 2 {
		n := uint32(raw[i])<<8 | uint32(raw[i+1])
		n = currentCastleRotR16ForTest(n, 5)
		n = (n - 60834) & 0xffff
		n = currentCastleRotR16ForTest(n, 4)
		n = ((n - 385) * 13069) & 0xffff
		n = (46488 ^ n) & 0xffff
		n = (n + 54655) & 0xffff
		units = append(units, uint16(n))
	}
	return units
}

func currentCastleRotR16ForTest(value uint32, shift uint) uint32 {
	return ((value >> shift) | (value << (16 - shift))) & 0xffff
}

func stripCurrentCastleChecksumForTest(t *testing.T, payload string) string {
	t.Helper()
	re := regexp.MustCompile(regexp.QuoteMeta(currentCastleChecksumLabel) + `[0-9a-f]{8}`)
	loc := re.FindStringIndex(payload)
	if loc == nil {
		t.Fatalf("decoded payload missing checksum marker")
	}
	return payload[:loc[0]] + payload[loc[1]:]
}

func currentCastleJSONSummaryForTest(t *testing.T, value any) string {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal(summary) error = %v", err)
	}
	return string(data)
}

func currentCastleJSONTypeForTest(value any) string {
	switch value.(type) {
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "bool"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("%T", value)
	}
}

func currentCastleShortHashForTest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum[:4])
}

func decodeCurrentCastleFloat64ArraySlotForTest(t *testing.T, slot int, encoded string) []float64 {
	t.Helper()
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("base64 slot %d decode error = %v", slot, err)
	}
	if len(raw)%8 != 0 {
		t.Fatalf("slot %d decoded to byte length %d, not divisible by 8", slot, len(raw))
	}
	for i, b := range raw {
		switch slot {
		case 254:
			raw[i] = decodeCurrentCastleArraySlot254ByteForTest(b)
		case 302:
			raw[i] = decodeCurrentCastleArraySlot302ByteForTest(b)
		default:
			t.Fatalf("unsupported array slot %d", slot)
		}
	}
	values := make([]float64, len(raw)/8)
	for i := range values {
		values[i] = math.Float64frombits(binary.BigEndian.Uint64(raw[i*8:]))
	}
	return values
}

func decodeCurrentCastleArraySlot254ByteForTest(b byte) byte {
	n := uint32(b)
	n = (n - 195) & 0xff
	n = (n * 13) & 0xff // modular inverse of 185*109 mod 256
	n ^= 32
	return byte(n)
}

func decodeCurrentCastleArraySlot302ByteForTest(b byte) byte {
	n := uint32(b) ^ 192
	n = currentCastleRotR8ForTest(n, 4)
	n = (n * 23) & 0xff // modular inverse of 167 mod 256
	return byte(n)
}

func currentCastleRotR8ForTest(value uint32, shift uint) uint32 {
	return ((value >> shift) | (value << (8 - shift))) & 0xff
}

func decodeCurrentCastleSlot6ComponentsForTest(t *testing.T, encoded string) []string {
	t.Helper()
	units := decodeCurrentCastleUnitsForTest(t, encoded)
	index := 0
	count := readCurrentCastleLengthPrefixForTest(t, units, &index)
	if count != currentCastleSlot6ComponentCount {
		t.Fatalf("slot 6 component count = %d, want %d", count, currentCastleSlot6ComponentCount)
	}
	lengths := make([]int, count)
	for i := range lengths {
		lengths[i] = readCurrentCastleLengthPrefixForTest(t, units, &index)
	}
	components := make([]string, count)
	for i, length := range lengths {
		if index+length > len(units) {
			t.Fatalf("slot 6 component %d length %d overflows decoded unit len %d at %d", i, length, len(units), index)
		}
		components[i] = string(utf16.Decode(units[index : index+length]))
		index += length
	}
	if index != len(units) {
		t.Fatalf("slot 6 decoded %d units but consumed %d", len(units), index)
	}
	return components
}

func readCurrentCastleLengthPrefixForTest(t *testing.T, units []uint16, index *int) int {
	t.Helper()
	length := 0
	shift := 0
	for {
		if *index >= len(units) {
			t.Fatalf("slot 6 length prefix overflow")
		}
		unit := units[*index]
		*index = *index + 1
		length |= int(unit&currentCastleSlot6LengthMask) << shift
		if unit&currentCastleSlot6LengthContinue == 0 {
			return length
		}
		shift += currentCastleSlot6LengthShiftBits
	}
}

func decodeCurrentCastleSlot6EncodedComponentForTest(t *testing.T, index int, encoded string) string {
	t.Helper()
	if index < 0 || index >= currentCastleSlot6ComponentCount {
		t.Fatalf("slot 6 component index %d outside range", index)
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
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("base64 slot 6 component %d decode error: %v", index, err)
	}
	inverse := map[byte]byte{}
	for b := 0; b < 256; b++ {
		probe, err := base64.StdEncoding.DecodeString(encoders[index](string([]byte{byte(b)})))
		if err != nil {
			t.Fatalf("base64 slot 6 component %d inverse probe decode error: %v", index, err)
		}
		if len(probe) != 1 {
			t.Fatalf("slot 6 component %d inverse probe len = %d, want 1", index, len(probe))
		}
		inverse[probe[0]] = byte(b)
	}
	out := make([]byte, len(raw))
	for i, b := range raw {
		value, ok := inverse[b]
		if !ok {
			t.Fatalf("slot 6 component %d byte %d has no inverse for %#x", index, i, b)
		}
		out[i] = value
	}
	return string(out)
}

func decodeCurrentCastleWrappedTokenForTest(token string) (string, error) {
	if !strings.HasPrefix(token, currentCastleTokenPrefix) {
		return "", errCurrentCastleTestMissingPrefix
	}
	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(token, currentCastleTokenPrefix))
	if err != nil {
		return "", err
	}
	reader := flate.NewReader(bytes.NewReader(raw))
	defer reader.Close()
	out, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func sameUint16s(a, b []uint16) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sameInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

var errCurrentCastleTestMissingPrefix = errors.New("current Castle token missing prefix")
