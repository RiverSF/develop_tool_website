package common

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func CreateUuid() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func String(ad interface{}) string {
	b, err := json.Marshal(ad)
	if err != nil {
		return fmt.Sprintf("%v", ad)
	}
	return string(b)
}

func HmacSha256(ad string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(ad))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5(str string) string {
	if len(str) == 0 {
		return ""
	}
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func StringToFloat64(str string) float64 {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return res
}

func Float64ToString(value float64) string {
	return fmt.Sprintf("%v", value)
}

func StringToInt(value string) int {
	res, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return res
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}
