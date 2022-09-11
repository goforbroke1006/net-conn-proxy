package internal

import "strings"

func GetHostFromAddr(addr string) string {
	parts := strings.Split(addr, ":")
	return parts[0]
}
