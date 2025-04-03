package types

type ContentType string

const (
	ContentTypeNone ContentType = ""
	ContentTypeJSON ContentType = "application/json"
	ContentTypeForm ContentType = "application/x-www-form-urlencoded"
)
