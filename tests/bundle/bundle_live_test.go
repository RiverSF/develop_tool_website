package bundle_test

import (
	"testing"

	"develop_tools/pkg/bundle"
)

func TestFetchAppleLive(t *testing.T) {
	if testing.Short() {
		t.Skip("skip live network test")
	}
	info, err := bundle.FetchApple("310633997")
	if err != nil {
		t.Fatalf("FetchApple: %v", err)
	}
	appInfo, ok := info["app"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing app: %+v", info)
	}
	developer, ok := info["developer"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing developer: %+v", info)
	}
	if name, _ := appInfo["trackName"].(string); name == "" {
		t.Fatalf("expected trackName, got %+v", appInfo)
	}
	if bid, _ := appInfo["bundleId"].(string); bid == "" {
		t.Fatalf("expected bundleId, got %+v", appInfo)
	}
	if developer["artistName"] == nil || developer["artistName"] == "" {
		t.Fatalf("expected developer.artistName, got %+v", developer)
	}
}

func TestFetchGoogleLive(t *testing.T) {
	if testing.Short() {
		t.Skip("skip live network test")
	}
	info, err := bundle.FetchGoogle("com.whatsapp")
	if err != nil {
		t.Fatalf("FetchGoogle: %v", err)
	}
	appInfo, ok := info["app"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing app: %+v", info)
	}
	developer, ok := info["developer"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing developer: %+v", info)
	}

	appRequired := []string{
		"id", "url", "locale", "country", "name", "description", "developerName",
		"icon", "screenshots", "score", "installsText", "cover", "category",
		"privacyPoliceUrl", "recentChange", "editorsChoice", "installs",
		"numberVoters", "histogramRating", "price", "currency", "offersIAP",
		"contentRating", "released", "releasedTimestamp", "updated", "updatedTimestamp",
		"numberReviews", "reviews",
	}
	for _, key := range appRequired {
		if _, ok := appInfo[key]; !ok {
			t.Fatalf("missing app.%s in %+v", key, keysOf(appInfo))
		}
	}
	if appInfo["id"] != "com.whatsapp" {
		t.Fatalf("unexpected id: %v", appInfo["id"])
	}
	if name, _ := appInfo["name"].(string); name == "" {
		t.Fatalf("empty name")
	}
	cat, _ := appInfo["category"].(map[string]interface{})
	if cat["id"] == nil && cat["name"] == nil {
		t.Fatalf("empty category: %+v", cat)
	}
	hist, _ := appInfo["histogramRating"].(map[string]interface{})
	for _, k := range []string{"one", "two", "three", "four", "five"} {
		if _, ok := hist[k]; !ok {
			t.Fatalf("histogramRating missing %s", k)
		}
	}
	reviews, _ := appInfo["reviews"].([]map[string]interface{})
	if len(reviews) == 0 {
		// json decode may produce []interface{}
		if arr, ok := appInfo["reviews"].([]interface{}); !ok || len(arr) == 0 {
			t.Fatalf("expected reviews, got %#v", appInfo["reviews"])
		}
	}

	devRequired := []string{"id", "url", "name", "website", "email", "address"}
	for _, key := range devRequired {
		if _, ok := developer[key]; !ok {
			t.Fatalf("missing developer.%s in %+v", key, keysOf(developer))
		}
	}
	if developer["email"] == nil || developer["email"] == "" {
		t.Fatalf("developer email empty: %+v", developer)
	}
}

func keysOf(m map[string]interface{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestExtractFilterLive(t *testing.T) {
	if testing.Short() {
		t.Skip("skip live network test")
	}
	results := bundle.Extract("310633997", "trackName,bundleId")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0]
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	appInfo, ok := r.Info["app"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected app object, got %+v", r.Info)
	}
	if len(appInfo) != 2 {
		t.Fatalf("expected 2 app fields, got %+v", appInfo)
	}
	if _, hasPlatform := r.Info["platform"]; hasPlatform {
		t.Fatal("platform should not appear in result payload")
	}
}
