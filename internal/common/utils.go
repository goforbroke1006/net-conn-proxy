package common

import (
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

func GetHostFromAddr(addr string) string {
	parts := strings.Split(addr, ":")
	return parts[0]
}

func GetPortFromAddr(addr string) uint16 {
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		panic(errors.New("invalid addr: " + addr))
	}

	port, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		panic(err)
	}

	return uint16(port)
}

func GetPrettyHexString(buffer []byte) string {
	result := ""
	for _, b := range buffer {
		result += hex.EncodeToString([]byte{b}) + " "
	}
	result = strings.TrimSpace(result)
	return result
}
