package matrixfmt

import (
	"fmt"
)

type BodyRangeValue interface {
	String() string
	Format(message string) string
}

type Mention struct {
	ID string
}

func (m Mention) String() string {
	return fmt.Sprintf("Mention{ID: (%s)}", m.ID)
}

func (m Mention) Format(message string) string {
	return message
}

type Style int

const (
	StyleNone Style = iota
	StyleBold
	StyleItalic
	StyleStrikethrough
	StyleSourceCode
	StyleMonospace // 5
	StyleHidden
	StyleMonospaceBlock
	StyleUnderline
	StyleFontColor
)

func (s Style) String() string {
	return fmt.Sprintf("Style(%d)", s)
}

func (s Style) Format(message string) string {
	return message
}
