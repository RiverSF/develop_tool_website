package bundle

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	pkgnet "develop_tools/pkg/net"
)

const (
	PlatformApple   = "apple"
	PlatformGoogle  = "google"
	PlatformUnknown = "unknown"
)

var (
	appleIDRegex    = regexp.MustCompile(`(?i)(?:^|/)id(\d{5,})(?:\b|$)`)
	digitsOnlyRegex = regexp.MustCompile(`^\d{5,}$`)
	googlePkgRegex  = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)+$`)
	playIDRegex     = regexp.MustCompile(`(?i)[?&]id=([a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)+)`)
	metaTagRegex    = regexp.MustCompile(`(?is)<meta\s+[^>]*?(?:property|name)\s*=\s*["']([^"']+)["'][^>]*?content\s*=\s*["']([^"']*)["'][^>]*?/?>`)
	metaTagRegex2   = regexp.MustCompile(`(?is)<meta\s+[^>]*?content\s*=\s*["']([^"']*)["'][^>]*?(?:property|name)\s*=\s*["']([^"']+)["'][^>]*?/?>`)
	ldJSONRegex     = regexp.MustCompile(`(?is)<script[^>]*type=["']application/ld\+json["'][^>]*>(.*?)</script>`)
	ds5DataRegex    = regexp.MustCompile(`AF_initDataCallback\(\{key:\s*'ds:5',\s*hash:\s*'[^']*',\s*data:(.*?),\s*sideChannel:`)
	ds10DataRegex   = regexp.MustCompile(`AF_initDataCallback\(\{key:\s*'ds:10',\s*hash:\s*'[^']*',\s*data:(.*?),\s*sideChannel:`)
	htmlTagRegex    = regexp.MustCompile(`(?i)<[^>]+>`)
)

// DetectedBundle is a normalized input after platform detection.
type DetectedBundle struct {
	Raw      string `json:"raw"`
	Bundle   string `json:"bundle"`
	Platform string `json:"platform"`
}

// AppResult is one bundle's crawl result.
type AppResult struct {
	Bundle string                 `json:"bundle"`
	Info   map[string]interface{} `json:"info,omitempty"`
	Error  string                 `json:"error,omitempty"`
}

// DetectPlatform classifies a single bundle / URL / id.
func DetectPlatform(raw string) DetectedBundle {
	raw = strings.TrimSpace(raw)
	out := DetectedBundle{Raw: raw, Bundle: raw, Platform: PlatformUnknown}
	if raw == "" {
		return out
	}

	lower := strings.ToLower(raw)

	if strings.Contains(lower, "play.google.com") {
		if m := playIDRegex.FindStringSubmatch(raw); len(m) > 1 {
			out.Bundle = m[1]
			out.Platform = PlatformGoogle
			return out
		}
	}

	if strings.Contains(lower, "apps.apple.com") || strings.Contains(lower, "itunes.apple.com") {
		if m := appleIDRegex.FindStringSubmatch(raw); len(m) > 1 {
			out.Bundle = m[1]
			out.Platform = PlatformApple
			return out
		}
		if u, err := url.Parse(raw); err == nil {
			q := u.Query().Get("id")
			if digitsOnlyRegex.MatchString(q) {
				out.Bundle = q
				out.Platform = PlatformApple
				return out
			}
		}
	}

	if strings.HasPrefix(lower, "id") && digitsOnlyRegex.MatchString(raw[2:]) {
		out.Bundle = raw[2:]
		out.Platform = PlatformApple
		return out
	}

	if digitsOnlyRegex.MatchString(raw) {
		out.Bundle = raw
		out.Platform = PlatformApple
		return out
	}

	if googlePkgRegex.MatchString(raw) {
		out.Bundle = raw
		out.Platform = PlatformGoogle
		return out
	}

	return out
}

// ParseBundles splits multiline / comma-separated input and detects each item.
// Empty lines are skipped; order and duplicates are preserved so card count matches input.
func ParseBundles(content string) []DetectedBundle {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	content = strings.ReplaceAll(content, ",", "\n")
	content = strings.ReplaceAll(content, "，", "\n")
	content = strings.ReplaceAll(content, "\t", "\n")
	content = strings.ReplaceAll(content, " ", "\n")

	var list []DetectedBundle
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		list = append(list, DetectPlatform(line))
	}
	return list
}

// ParseFields splits optional display field list. Empty means show all.
func ParseFields(fields string) []string {
	fields = strings.ReplaceAll(fields, "\r\n", "\n")
	fields = strings.ReplaceAll(fields, "\r", "\n")
	fields = strings.ReplaceAll(fields, ",", "\n")
	fields = strings.ReplaceAll(fields, "，", "\n")
	fields = strings.ReplaceAll(fields, "\t", "\n")
	fields = strings.ReplaceAll(fields, " ", "\n")

	var out []string
	seen := make(map[string]struct{})
	for _, f := range strings.Split(fields, "\n") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		if _, ok := seen[f]; ok {
			continue
		}
		seen[f] = struct{}{}
		out = append(out, f)
	}
	return out
}

// FilterInfo keeps only requested keys. Empty fields returns all.
// For results shaped as {app, developer}, fields apply to app;
// include developer only when "developer" is requested.
func FilterInfo(info map[string]interface{}, fields []string) map[string]interface{} {
	if len(fields) == 0 || info == nil {
		return info
	}
	if appInfo, ok := info["app"].(map[string]interface{}); ok {
		out := map[string]interface{}{
			"app": filterTopKeys(appInfo, fields),
		}
		if fieldRequested(fields, "developer") {
			if d, ok := info["developer"]; ok {
				out["developer"] = d
			}
		}
		return out
	}
	return filterTopKeys(info, fields)
}

func filterTopKeys(info map[string]interface{}, fields []string) map[string]interface{} {
	out := make(map[string]interface{}, len(fields))
	for _, f := range fields {
		if v, ok := info[f]; ok {
			out[f] = v
		}
	}
	return out
}

func fieldRequested(fields []string, name string) bool {
	for _, f := range fields {
		if f == name {
			return true
		}
	}
	return false
}

// Extract fetches store info for each bundle and optionally filters fields.
func Extract(content string, fieldsRaw string) []AppResult {
	bundles := ParseBundles(content)
	fields := ParseFields(fieldsRaw)
	if len(bundles) == 0 {
		return nil
	}

	results := make([]AppResult, len(bundles))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for i, b := range bundles {
		wg.Add(1)
		go func(i int, b DetectedBundle) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			res := AppResult{Bundle: b.Bundle}
			switch b.Platform {
			case PlatformApple:
				info, err := FetchApple(b.Bundle)
				if err != nil {
					res.Error = err.Error()
				} else {
					res.Info = FilterInfo(info, fields)
				}
			case PlatformGoogle:
				info, err := FetchGoogle(b.Bundle)
				if err != nil {
					res.Error = err.Error()
				} else {
					res.Info = FilterInfo(info, fields)
				}
			default:
				res.Error = "无法识别 bundle 格式（纯数字为 Apple，包名如 com.xxx 为 Google）"
			}
			results[i] = res
		}(i, b)
	}
	wg.Wait()
	return results
}

// FetchApple uses iTunes Lookup API and returns {app, developer} like Google.
func FetchApple(idOrBundle string) (map[string]interface{}, error) {
	apiURL := "https://itunes.apple.com/lookup?id=" + url.QueryEscape(idOrBundle)
	if !digitsOnlyRegex.MatchString(idOrBundle) {
		apiURL = "https://itunes.apple.com/lookup?bundleId=" + url.QueryEscape(idOrBundle)
	}

	_, body, err := pkgnet.HttpGetRequest(apiURL, pkgnet.HttpClient5000)
	if err != nil {
		return nil, fmt.Errorf("请求 App Store 失败: %w", err)
	}

	var resp struct {
		ResultCount int                      `json:"resultCount"`
		Results     []map[string]interface{} `json:"results"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("解析 App Store 响应失败: %w", err)
	}
	if resp.ResultCount == 0 || len(resp.Results) == 0 {
		return nil, fmt.Errorf("未找到应用: %s", idOrBundle)
	}
	return shapeAppleInfo(resp.Results[0]), nil
}

// shapeAppleInfo splits iTunes Lookup result into {app, developer}.
// Field names inside both objects keep Apple's original naming.
func shapeAppleInfo(raw map[string]interface{}) map[string]interface{} {
	if raw == nil {
		return map[string]interface{}{
			"app":       map[string]interface{}{},
			"developer": map[string]interface{}{},
		}
	}

	developerKeys := []string{
		"artistId",
		"artistName",
		"artistViewUrl",
		"sellerName",
		"sellerUrl",
	}
	developerKeySet := make(map[string]struct{}, len(developerKeys))
	for _, k := range developerKeys {
		developerKeySet[k] = struct{}{}
	}

	app := make(map[string]interface{}, len(raw))
	developer := make(map[string]interface{}, len(developerKeys))
	for k, v := range raw {
		if _, isDev := developerKeySet[k]; isDev {
			developer[k] = v
			continue
		}
		app[k] = v
	}

	return map[string]interface{}{
		"app":       app,
		"developer": developer,
	}
}

// FetchGoogle scrapes Google Play and returns {app, developer}.
func FetchGoogle(packageName string) (map[string]interface{}, error) {
	apiURL := "https://play.google.com/store/apps/details?id=" + url.QueryEscape(packageName) + "&hl=en&gl=us"
	resp, body, err := httpGetWithUA(apiURL, pkgnet.HttpClient120000)
	if err != nil {
		return nil, fmt.Errorf("请求 Google Play 失败: %w", err)
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("未找到应用: %s", packageName)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Google Play 返回状态码 %d", resp.StatusCode)
	}

	htmlBody := string(body)
	appInfo, developer, ok := parsePlayAppInfo(htmlBody, packageName)
	if !ok {
		return nil, fmt.Errorf("未能解析 Google Play 页面信息: %s", packageName)
	}

	// Fill gaps from ld+json / meta when ds:5 missed them.
	fallback := map[string]interface{}{}
	for k, v := range extractMetaTags(htmlBody) {
		fallback[k] = v
	}
	mergeLDJSON(htmlBody, fallback)
	applyGoogleFallback(appInfo, developer, fallback)

	out := map[string]interface{}{
		"app":       appInfo,
		"developer": developer,
	}
	return out, nil
}

func parsePlayAppInfo(page, packageName string) (map[string]interface{}, map[string]interface{}, bool) {
	m := ds5DataRegex.FindStringSubmatch(page)
	if len(m) < 2 {
		return nil, nil, false
	}
	var root interface{}
	if err := json.Unmarshal([]byte(m[1]), &root); err != nil {
		return nil, nil, false
	}
	app := pathGet(root, 1, 2)
	if app == nil {
		return nil, nil, false
	}

	name := pathString(app, 0, 0)
	pkg := pathString(app, 77, 0)
	if pkg == "" {
		pkg = packageName
	}
	if name == "" && pkg == "" {
		return nil, nil, false
	}

	descHTML := pathString(app, 72, 0, 1)
	recentHTML := pathString(app, 144, 1, 1)
	iapCost := pathString(app, 19, 0)
	devName := firstNonEmpty(pathString(app, 68, 0), pathString(app, 37, 0))
	devPath := pathString(app, 68, 1, 4, 2)
	devID := ""
	devURL := ""
	if devPath != "" {
		if strings.HasPrefix(devPath, "/") {
			devURL = "https://play.google.com" + devPath
		} else {
			devURL = devPath
		}
		if u, err := url.Parse(devPath); err == nil {
			devID = u.Query().Get("id")
		}
	}
	if devID == "" {
		devID = devName
	}

	releasedTS := pathNumber(app, 10, 1, 0)
	updatedTS := pathNumber(app, 145, 0, 1, 0)

	score := pathNumber(app, 51, 0, 1)
	numberVoters := pathNumber(app, 51, 2, 1)
	numberReviews := pathNumber(app, 51, 3, 1)
	installs := pathNumber(app, 13, 2)
	if installs == nil {
		installs = pathNumber(app, 13, 1)
	}

	price := 0.0
	currency := "USD"
	priceText := interface{}(nil)
	if p := pathNumber(app, 57, 0, 0, 0, 0, 1, 0, 0); p != nil {
		if f, ok := p.(float64); ok && f > 0 {
			price = f
			priceText = fmt.Sprintf("%v", f)
		}
	}
	if c := pathString(app, 57, 0, 0, 0, 0, 1, 0, 1); c != "" {
		currency = c
	}

	appInfo := map[string]interface{}{
		"id":                pkg,
		"url":               "https://play.google.com/store/apps/details?id=" + pkg,
		"locale":            "en_US",
		"country":           "us",
		"name":              name,
		"description":       stripHTML(descHTML),
		"developerName":     devName,
		"icon":              nilOrString(pathString(app, 95, 0, 3, 2)),
		"screenshots":       extractScreenshots(pathGet(app, 78, 0)),
		"score":             score,
		"priceText":         priceText,
		"installsText":      nilOrString(firstNonEmpty(pathString(app, 13, 0), pathString(app, 13, 3))),
		"cover":             nilOrString(pathString(app, 96, 0, 3, 2)),
		"category":          buildCategory(app),
		"categoryFamily":    nil,
		"video":             nil,
		"privacyPoliceUrl":  nilOrString(pathString(app, 99, 0, 5, 2)),
		"recentChange":      nilOrString(stripHTML(recentHTML)),
		"editorsChoice":     false,
		"installs":          installs,
		"numberVoters":      numberVoters,
		"histogramRating":   extractHistogramRating(pathGet(app, 51, 1)),
		"price":             price,
		"currency":          currency,
		"offersIAP":         iapCost != "",
		"offersIAPCost":     nilOrString(iapCost),
		"containsAds":       false,
		"appVersion":        nil,
		"androidVersion":    nil,
		"minAndroidVersion": nil,
		"contentRating":     pathString(app, 9, 0),
		"released":          formatPlayTime(pathString(app, 10, 0), releasedTS),
		"releasedTimestamp": releasedTS,
		"updated":           formatPlayTime(pathString(app, 145, 0, 0), updatedTS),
		"updatedTimestamp":  updatedTS,
		"numberReviews":     numberReviews,
		"reviews":           parsePlayReviews(page),
	}

	developer := map[string]interface{}{
		"id":          nilOrString(devID),
		"url":         nilOrString(devURL),
		"name":        nilOrString(devName),
		"description": nil,
		"website":     nilOrString(pathString(app, 69, 0, 5, 2)),
		"icon":        nil,
		"cover":       nil,
		"email":       nilOrString(firstNonEmpty(pathString(app, 69, 1, 0), pathString(app, 69, 4, 1, 0))),
		"address":     nilOrString(pathString(app, 69, 4, 2, 0)),
	}

	return appInfo, developer, true
}

func applyGoogleFallback(appInfo, developer, fallback map[string]interface{}) {
	if isEmpty(appInfo["name"]) {
		if v, ok := fallback["title"].(string); ok {
			appInfo["name"] = v
		}
	}
	if isEmpty(appInfo["description"]) {
		if v, ok := fallback["description"].(string); ok {
			appInfo["description"] = v
		}
	}
	if isEmpty(appInfo["icon"]) {
		if v, ok := fallback["icon"].(string); ok {
			appInfo["icon"] = v
		}
	}
	if isEmpty(appInfo["developerName"]) {
		if v, ok := fallback["developer"].(string); ok {
			appInfo["developerName"] = v
			if isEmpty(developer["name"]) {
				developer["name"] = v
			}
			if isEmpty(developer["id"]) {
				developer["id"] = v
			}
		}
	}
	if isEmpty(appInfo["score"]) {
		if v, ok := fallback["rating"]; ok {
			appInfo["score"] = v
		}
	}
	if isEmpty(appInfo["numberVoters"]) {
		if v, ok := fallback["ratingCount"]; ok {
			appInfo["numberVoters"] = v
		}
	}
	if cat, ok := appInfo["category"].(map[string]interface{}); !ok || isEmpty(cat["name"]) {
		if v, ok := fallback["category"].(string); ok && v != "" {
			appInfo["category"] = map[string]interface{}{"id": v, "name": v}
		}
	}
	if price, ok := fallback["price"]; ok {
		switch t := price.(type) {
		case float64:
			appInfo["price"] = t
		case string:
			if f, err := strconv.ParseFloat(t, 64); err == nil {
				appInfo["price"] = f
			}
		}
	}
	if c, ok := fallback["currency"].(string); ok && c != "" {
		appInfo["currency"] = c
	}
	if isEmpty(developer["website"]) {
		if v, ok := fallback["developerUrl"].(string); ok {
			developer["website"] = v
		}
	}
}

func buildCategory(app interface{}) interface{} {
	name := pathString(app, 79, 0, 0, 0)
	id := pathString(app, 79, 0, 0, 2)
	if name == "" && id == "" {
		return nil
	}
	return map[string]interface{}{
		"id":   nilOrString(id),
		"name": nilOrString(name),
	}
}

func extractHistogramRating(v interface{}) map[string]interface{} {
	arr, ok := v.([]interface{})
	if !ok || len(arr) < 6 {
		return map[string]interface{}{
			"five": nil, "four": nil, "three": nil, "two": nil, "one": nil,
		}
	}
	return map[string]interface{}{
		"one":   pathNumber(arr, 1, 1),
		"two":   pathNumber(arr, 2, 1),
		"three": pathNumber(arr, 3, 1),
		"four":  pathNumber(arr, 4, 1),
		"five":  pathNumber(arr, 5, 1),
	}
}

func parsePlayReviews(page string) []map[string]interface{} {
	m := ds10DataRegex.FindStringSubmatch(page)
	if len(m) < 2 {
		return []map[string]interface{}{}
	}
	var root interface{}
	if err := json.Unmarshal([]byte(m[1]), &root); err != nil {
		return []map[string]interface{}{}
	}
	list, ok := pathGet(root, 0).([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}
	out := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		id := pathString(item, 0)
		if id == "" {
			continue
		}
		ts := pathNumber(item, 5, 0)
		out = append(out, map[string]interface{}{
			"id":         id,
			"userName":   pathString(item, 1, 0),
			"text":       pathString(item, 4),
			"avatar":     nilOrString(pathString(item, 1, 1, 3, 2)),
			"appVersion": nilOrString(pathString(item, 10)),
			"date":       formatPlayTime("", ts),
			"timestamp":  ts,
			"score":      pathNumber(item, 2),
			"countLikes": pathNumber(item, 6),
			"reply":      nil,
		})
	}
	return out
}

func formatPlayTime(text string, ts interface{}) interface{} {
	if f, ok := ts.(float64); ok && f > 0 {
		return time.Unix(int64(f), 0).UTC().Format(time.RFC3339)
	}
	if text != "" {
		return text
	}
	return nil
}

func nilOrString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func isEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case string:
		return t == ""
	case map[string]interface{}:
		return len(t) == 0
	case []interface{}:
		return len(t) == 0
	}
	return false
}

func extractScreenshots(v interface{}) []string {
	arr, ok := v.([]interface{})
	if !ok {
		return []string{}
	}
	var out []string
	for _, item := range arr {
		if urlStr := pathString(item, 3, 2); urlStr != "" {
			out = append(out, urlStr)
		}
	}
	return out
}

func pathGet(v interface{}, idxs ...int) interface{} {
	cur := v
	for _, i := range idxs {
		arr, ok := cur.([]interface{})
		if !ok || i < 0 || i >= len(arr) {
			return nil
		}
		cur = arr[i]
	}
	return cur
}

func pathString(v interface{}, idxs ...int) string {
	cur := pathGet(v, idxs...)
	if cur == nil && len(idxs) == 0 {
		cur = v
	}
	switch t := cur.(type) {
	case string:
		return strings.TrimSpace(t)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	default:
		return ""
	}
}

func pathNumber(v interface{}, idxs ...int) interface{} {
	cur := pathGet(v, idxs...)
	switch t := cur.(type) {
	case float64:
		return t
	case json.Number:
		if f, err := t.Float64(); err == nil {
			return f
		}
	case string:
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f
		}
	}
	return nil
}

func stripHTML(s string) string {
	s = html.UnescapeString(s)
	s = strings.ReplaceAll(s, "<br>", "\n")
	s = strings.ReplaceAll(s, "<br/>", "\n")
	s = strings.ReplaceAll(s, "<br />", "\n")
	s = htmlTagRegex.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\u00a0", " ")
	return strings.TrimSpace(s)
}

func httpGetWithUA(apiURL string, client *http.Client) (*http.Response, []byte, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}
	return resp, body, nil
}

func extractMetaTags(page string) map[string]interface{} {
	out := make(map[string]interface{})
	put := func(name, content string) {
		name = strings.ToLower(strings.TrimSpace(name))
		content = html.UnescapeString(strings.TrimSpace(content))
		if content == "" {
			return
		}
		switch name {
		case "og:title", "twitter:title":
			setIfAbsent(out, "title", content)
		case "og:description", "twitter:description", "description":
			setIfAbsent(out, "description", content)
		case "og:image", "twitter:image":
			setIfAbsent(out, "icon", content)
		case "og:url":
			setIfAbsent(out, "storeUrl", content)
		case "appstore:developer_url":
			setIfAbsent(out, "developerUrl", content)
		}
	}

	for _, m := range metaTagRegex.FindAllStringSubmatch(page, -1) {
		if len(m) >= 3 {
			put(m[1], m[2])
		}
	}
	for _, m := range metaTagRegex2.FindAllStringSubmatch(page, -1) {
		if len(m) >= 3 {
			put(m[2], m[1])
		}
	}
	return out
}

func mergeLDJSON(page string, info map[string]interface{}) {
	matches := ldJSONRegex.FindAllStringSubmatch(page, -1)
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		raw := strings.TrimSpace(m[1])
		var any interface{}
		if err := json.Unmarshal([]byte(raw), &any); err != nil {
			continue
		}
		walkLDJSON(any, info)
	}
}

func walkLDJSON(v interface{}, info map[string]interface{}) {
	switch t := v.(type) {
	case []interface{}:
		for _, item := range t {
			walkLDJSON(item, info)
		}
	case map[string]interface{}:
		typeVal, _ := t["@type"].(string)
		if strings.Contains(strings.ToLower(typeVal), "softwareapplication") || typeVal == "MobileApplication" {
			if name, ok := t["name"].(string); ok {
				setIfAbsent(info, "title", name)
			}
			if desc, ok := t["description"].(string); ok {
				setIfAbsent(info, "description", desc)
			}
			if img, ok := t["image"].(string); ok {
				setIfAbsent(info, "icon", img)
			}
			if imgs, ok := t["image"].([]interface{}); ok && len(imgs) > 0 {
				if s, ok := imgs[0].(string); ok {
					setIfAbsent(info, "icon", s)
				}
			}
			if author, ok := t["author"].(map[string]interface{}); ok {
				if name, ok := author["name"].(string); ok {
					setIfAbsent(info, "developer", name)
				}
			}
			if offers, ok := t["offers"].(map[string]interface{}); ok {
				if price, ok := offers["price"]; ok {
					setIfAbsent(info, "price", price)
				}
				if currency, ok := offers["priceCurrency"].(string); ok {
					setIfAbsent(info, "currency", currency)
				}
			}
			if rating, ok := t["aggregateRating"].(map[string]interface{}); ok {
				if v, ok := rating["ratingValue"]; ok {
					setIfAbsent(info, "rating", toFloat(v))
				}
				if v, ok := rating["ratingCount"]; ok {
					setIfAbsent(info, "ratingCount", toFloat(v))
				}
			}
			if cat, ok := t["applicationCategory"].(string); ok {
				setIfAbsent(info, "category", cat)
			}
			if os, ok := t["operatingSystem"].(string); ok {
				setIfAbsent(info, "operatingSystem", os)
			}
		}
		for _, child := range t {
			walkLDJSON(child, info)
		}
	}
}

func setIfAbsent(m map[string]interface{}, key string, val interface{}) {
	if _, ok := m[key]; ok {
		return
	}
	m[key] = val
}

func toFloat(v interface{}) interface{} {
	switch t := v.(type) {
	case float64:
		return t
	case string:
		if f, err := strconv.ParseFloat(t, 64); err == nil {
			return f
		}
		return t
	default:
		return v
	}
}
