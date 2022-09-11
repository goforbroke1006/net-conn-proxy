package internal

import (
	"encoding/hex"
	"strings"
)

func GetHostFromAddr(addr string) string {
	parts := strings.Split(addr, ":")
	return parts[0]
}

func GetPrettyHexString(buffer []byte) string {
	result := ""
	for _, b := range buffer {
		result += hex.EncodeToString([]byte{b}) + " "
	}
	result = strings.TrimSpace(result)
	return result
}
