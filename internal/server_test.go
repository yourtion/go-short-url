package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func buildCreateRequest(url string) *http.Request {
	body := strings.NewReader(`{"url": "` + url + `"}`)
	r, _ := http.NewRequest("POST", "/api/create", body)
	r.Header.Set("Content-Type", "application/json")
	return r
}

func getResponseUrl(resp []byte) string {
	return jsoniter.Get(resp, "data").Get("short").ToString()
}

func TestCreateShortURL(t *testing.T) {
	server := Server()

	r := buildCreateRequest("aaa")
	w := httptest.NewRecorder()

	// url not verify
	server.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Status not ok: %d", w.Code)
	}
	resp := string(w.Body.Bytes())
	if resp != `"url not verify"` {
		t.Fatalf("respones not ok: %s", resp)
	}
	t.Logf("url not verify pass!")

	// create url ok
	r = buildCreateRequest("https://baidu.com")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Status not ok: %d", w.Code)
	}
	url1 := getResponseUrl(w.Body.Bytes())
	if url1 == "" {
		t.Fatalf("URL Error: %s", string(w.Body.Bytes()))
	}
	t.Logf("create url pass! %s", url1)

	// create url ok
	r = buildCreateRequest("https://baidu.com")
	w = httptest.NewRecorder()
	server.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("Status not ok: %d", w.Code)
	}
	url2 := getResponseUrl(w.Body.Bytes())
	if url2 == "" {
		t.Fatalf("URL Error: %s", string(w.Body.Bytes()))
	}
	t.Logf("create url pass! %s", url1)

	if url1 != url2 {
		t.Fatalf("URL not match: %s != %s", url1, url2)
	}
}
