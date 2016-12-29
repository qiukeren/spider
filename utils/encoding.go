package utils

import (
	"github.com/endeveit/enca"
)

func EncodingTest(content *[]byte) (string, error) {
	analyzer, err := enca.New("zh")
	if err == nil {
		encoding, err := analyzer.FromBytes(*content, enca.NAME_STYLE_ICONV)
		defer analyzer.Free()

		if err == nil {
			return encoding, err
		}
	}
	return "", err
}
