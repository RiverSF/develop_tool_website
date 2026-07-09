package common

import (
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
)

// urlEncode
func UrlEscape(str string) string {
	return url.QueryEscape(str)
}

// urlUncode
func UrlUnescape(str string) string {
	unescape, err := url.QueryUnescape(str)
	if err != nil {
		return ""
	}
	return unescape
}

// base64解码
func Base64URLDecode(data string) ([]byte, error) {
	var missing = (4 - len(data)%4) % 4
	data += strings.Repeat("=", missing)
	res, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 安全Base64编码
func Base64UrlSafeEncode(source []byte) string {
	// Base64 Url Safe is the same as Base64 but does not contain '/' and '+' (replaced by '_' and '-') and trailing '=' are removed.
	bytearr := base64.StdEncoding.EncodeToString(source)
	safeurl := strings.ReplaceAll(string(bytearr), "/", "_")
	safeurl = strings.ReplaceAll(safeurl, "+", "-")
	safeurl = strings.ReplaceAll(safeurl, "=", "")
	return safeurl
}

// ValidateRedirectURL checks redirect targets against an allowlist of hostnames.
func ValidateRedirectURL(raw string, allowedHosts ...string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("invalid redirect scheme")
	}
	host := strings.ToLower(u.Hostname())
	if host == "" {
		return "", errors.New("invalid redirect host")
	}
	for _, allowed := range allowedHosts {
		if host == strings.ToLower(allowed) {
			return u.String(), nil
		}
	}
	return "", errors.New("redirect host not allowed")
}

// HostnameFromURL extracts hostname from a URL string.
func HostnameFromURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return u.Hostname()
}
