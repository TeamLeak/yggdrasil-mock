package utils

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateUUID() string {
	return strings.ReplaceAll(fmt.Sprintf("%s", base64.StdEncoding.EncodeToString([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))), "=", "")
}
