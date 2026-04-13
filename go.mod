module go.mau.fi/mautrix-twitter

go 1.25.0

toolchain go1.26.1

tool go.mau.fi/util/cmd/maubuild

require (
	github.com/PuerkitoBio/goquery v1.12.0
	github.com/apache/thrift v0.22.0
	github.com/bwesterb/go-ristretto v1.2.3
	github.com/coder/websocket v1.8.14
	github.com/fxamacker/cbor/v2 v2.9.0
	github.com/google/go-querystring v1.2.0
	github.com/google/uuid v1.6.0
	github.com/imroc/req/v3 v3.56.0
	github.com/openziti/secretstream v0.1.48
	github.com/rs/zerolog v1.35.0
	github.com/stretchr/testify v1.11.1
	github.com/tidwall/gjson v1.18.0
	go.mau.fi/util v0.9.8-0.20260406161447-0300c476893a
	golang.org/x/crypto v0.49.0
	golang.org/x/sync v0.20.0
	gopkg.in/yaml.v3 v3.0.1
	maunium.net/go/mautrix v0.26.5-0.20260412204845-73589d69756f
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/coreos/go-systemd/v22 v22.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/icholy/digest v1.1.0 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/kr/text v0.1.0 // indirect
	github.com/lib/pq v1.12.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.37 // indirect
	github.com/petermattis/goid v0.0.0-20260226131333-17d1149c6ac6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.56.0 // indirect
	github.com/refraction-networking/utls v1.8.2 // indirect
	github.com/rs/xid v1.6.0 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/yuin/goldmark v1.8.2 // indirect
	go.mau.fi/zeroconfig v0.2.0 // indirect
	golang.org/x/exp v0.0.0-20260312153236-7ab1446f8b90 // indirect
	golang.org/x/mod v0.34.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	maunium.net/go/mauflag v1.0.0 // indirect
)

replace github.com/imroc/req/v3 => github.com/beeper/req/v3 v3.0.0-20251116110214-7c681754bd16
