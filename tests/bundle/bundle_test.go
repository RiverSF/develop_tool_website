package bundle_test

import (
	"testing"

	"develop_tools/pkg/bundle"
)

func TestDetectPlatform(t *testing.T) {
	cases := []struct {
		in       string
		platform string
		bundle   string
	}{
		{"310633997", bundle.PlatformApple, "310633997"},
		{"id310633997", bundle.PlatformApple, "310633997"},
		{"https://apps.apple.com/us/app/whatsapp-messenger/id310633997", bundle.PlatformApple, "310633997"},
		{"https://itunes.apple.com/app/id310633997", bundle.PlatformApple, "310633997"},
		{"com.whatsapp", bundle.PlatformGoogle, "com.whatsapp"},
		{"https://play.google.com/store/apps/details?id=com.whatsapp", bundle.PlatformGoogle, "com.whatsapp"},
		{"not-a-bundle", bundle.PlatformUnknown, "not-a-bundle"},
		{"", bundle.PlatformUnknown, ""},
	}
	for _, c := range cases {
		got := bundle.DetectPlatform(c.in)
		if got.Platform != c.platform || got.Bundle != c.bundle {
			t.Fatalf("DetectPlatform(%q)=(%s,%s), want (%s,%s)", c.in, got.Platform, got.Bundle, c.platform, c.bundle)
		}
	}
}

func TestParseBundles(t *testing.T) {
	list := bundle.ParseBundles("com.whatsapp\n310633997,com.whatsapp\n")
	if len(list) != 3 {
		t.Fatalf("expected 3 bundles aligned with input, got %d: %+v", len(list), list)
	}
	if list[0].Bundle != "com.whatsapp" || list[1].Bundle != "310633997" || list[2].Bundle != "com.whatsapp" {
		t.Fatalf("unexpected order/content: %+v", list)
	}
}

func TestParseFields(t *testing.T) {
	fields := bundle.ParseFields("trackName, bundleId\nartistName")
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d: %v", len(fields), fields)
	}
	if len(bundle.ParseFields("  \n ")) != 0 {
		t.Fatal("empty fields should be empty slice")
	}
}

func TestFilterInfo(t *testing.T) {
	info := map[string]interface{}{
		"trackName": "WhatsApp",
		"bundleId":  "net.whatsapp.WhatsApp",
		"version":   "1.0",
	}
	all := bundle.FilterInfo(info, nil)
	if len(all) != 3 {
		t.Fatalf("nil fields should keep all, got %d", len(all))
	}
	filtered := bundle.FilterInfo(info, []string{"trackName", "missing"})
	if len(filtered) != 1 || filtered["trackName"] != "WhatsApp" {
		t.Fatalf("unexpected filtered: %+v", filtered)
	}

	google := map[string]interface{}{
		"app": map[string]interface{}{
			"name":  "WhatsApp Messenger",
			"id":    "com.whatsapp",
			"score": 4.6,
		},
		"developer": map[string]interface{}{
			"name":  "WhatsApp LLC",
			"email": "android@support.whatsapp.com",
		},
	}
	gFiltered := bundle.FilterInfo(google, []string{"name", "developer"})
	appInfo, _ := gFiltered["app"].(map[string]interface{})
	if len(appInfo) != 1 || appInfo["name"] != "WhatsApp Messenger" {
		t.Fatalf("unexpected google app filter: %+v", gFiltered)
	}
	if _, ok := gFiltered["developer"]; !ok {
		t.Fatalf("developer should be kept when requested: %+v", gFiltered)
	}

	apple := map[string]interface{}{
		"app": map[string]interface{}{
			"trackName": "WhatsApp",
			"bundleId":  "net.whatsapp.WhatsApp",
			"version":   "1.0",
		},
		"developer": map[string]interface{}{
			"artistName": "WhatsApp Inc.",
			"artistId":   12345,
		},
	}
	aFiltered := bundle.FilterInfo(apple, []string{"trackName", "developer"})
	appleApp, _ := aFiltered["app"].(map[string]interface{})
	if len(appleApp) != 1 || appleApp["trackName"] != "WhatsApp" {
		t.Fatalf("unexpected apple app filter: %+v", aFiltered)
	}
	dev, _ := aFiltered["developer"].(map[string]interface{})
	if dev["artistName"] != "WhatsApp Inc." {
		t.Fatalf("apple developer should keep original field names: %+v", aFiltered)
	}
}
