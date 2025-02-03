// Package parsers - functions to parse misc parameters
package parsers

import (
	"fmt"
	"regexp"
	"strconv"
)

func ParseInterval(s string) (int64, error) {
	re := regexp.MustCompile(`^[1-9]{1}[0-9]{0,2}[sm]$`)
	match := re.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("wrong interval")
	}
	intervalStr := s[:len(s)-1]
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		return 0, fmt.Errorf("wrong interval")
	}
	if string(s[len(s)-1]) == "m" {
		interval = interval * 60
	}
	return int64(interval), nil
}
