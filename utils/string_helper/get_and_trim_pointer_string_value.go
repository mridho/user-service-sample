package string_helper

import "strings"

func GetAndTrimPointerStringValue(s *string) string {
	if s != nil {
		return strings.TrimSpace(*s)
	}

	return ""
}
