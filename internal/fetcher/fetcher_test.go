package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"lp_fetcher_golang/internal/models"
)

func withoutPagePause(t *testing.T) {
	t.Helper()
	orig := pagePause
	pagePause = func(time.Duration) {}
	t.Cleanup(func() { pagePause = orig })
}

func makeItemBriefs(n int) []models.ItemBrief {
	items := make([]models.ItemBrief, n)
	for i := range items {
		items[i] = models.ItemBrief{
			ItemID:   fmt.Sprintf("id-%d", i),
			ItemName: fmt.Sprintf("Item %d", i),
			Price:    "99",
		}
	}
	return items
}

func searchResponseJSON(code int, message string, itemNames ...string) string {
	itemList := make([]models.Item, len(itemNames))
	for i, name := range itemNames {
		itemList[i] = models.Item{
			ItemID:   fmt.Sprintf("id-%d", i),
			ItemName: name,
			Price:    "99",
		}
	}
	resp := models.SearchResponse{
		Status: models.Status{Code: code, Message: message},
		Result: models.SearchResult{
			AllSuccess: code == 0,
			HasData:    len(itemList) > 0,
			ItemList:   itemList,
		},
	}
	data, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func TestParseSearchResponse(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		body := searchResponseJSON(0, "", "Album LP")
		got, err := parseSearchResponse([]byte(body))
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].ItemName != "Album LP" {
			t.Fatalf("got = %+v", got)
		}
	})

	t.Run("api error with message", func(t *testing.T) {
		body := searchResponseJSON(1001, "token invalid")
		_, err := parseSearchResponse([]byte(body))
		if err == nil || err.Error() != "token invalid" {
			t.Fatalf("err = %v, want token invalid", err)
		}
	})

	t.Run("api error without message", func(t *testing.T) {
		body := searchResponseJSON(1001, "")
		_, err := parseSearchResponse([]byte(body))
		if err == nil || err.Error() != "未知错误" {
			t.Fatalf("err = %v, want 未知错误", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		_, err := parseSearchResponse([]byte("not json"))
		if err == nil {
			t.Fatal("expected error for invalid json")
		}
	})
}

func TestBuildFetchResult(t *testing.T) {
	items := makeItemBriefs(2)

	t.Run("success", func(t *testing.T) {
		got := buildFetchResult(items, nil)
		if got.Status != "success" || got.Count != 2 || len(got.Items) != 2 {
			t.Fatalf("got = %+v", got)
		}
	})

	t.Run("too many results", func(t *testing.T) {
		err := &TooManyResultsError{Message: "查询结果过多、请精确查询关键词"}
		got := buildFetchResult(nil, err)
		if got.Status != "failed" || got.Msg != err.Error() || got.Items != nil {
			t.Fatalf("got = %+v", got)
		}
	})

	t.Run("generic error", func(t *testing.T) {
		got := buildFetchResult(nil, errors.New("HTTP 503"))
		if got.Status != "failed" || got.Msg != "HTTP 503" {
			t.Fatalf("got = %+v", got)
		}
	})
}

func TestFetchAllPages(t *testing.T) {
	withoutPagePause(t)

	t.Run("single partial page", func(t *testing.T) {
		got, err := fetchAllPages(func(offset int) ([]models.ItemBrief, error) {
			if offset == 0 {
				return makeItemBriefs(5), nil
			}
			return nil, nil
		}, searchMaxPages)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 5 {
			t.Fatalf("len(got) = %d, want 5", len(got))
		}
	})

	t.Run("empty first page", func(t *testing.T) {
		got, err := fetchAllPages(func(offset int) ([]models.ItemBrief, error) {
			return nil, nil
		}, searchMaxPages)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Fatalf("len(got) = %d, want 0", len(got))
		}
	})

	t.Run("unlimited pages until empty", func(t *testing.T) {
		got, err := fetchAllPages(func(offset int) ([]models.ItemBrief, error) {
			switch offset {
			case 0:
				return makeItemBriefs(pageSize), nil
			case pageSize:
				return makeItemBriefs(3), nil
			default:
				return nil, nil
			}
		}, 0)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != pageSize+3 {
			t.Fatalf("len(got) = %d, want %d", len(got), pageSize+3)
		}
	})

	t.Run("exceeds max pages", func(t *testing.T) {
		_, err := fetchAllPages(func(offset int) ([]models.ItemBrief, error) {
			return makeItemBriefs(pageSize), nil
		}, searchMaxPages)
		if err == nil {
			t.Fatal("expected TooManyResultsError")
		}
		var tooMany *TooManyResultsError
		if !errors.As(err, &tooMany) {
			t.Fatalf("err = %T(%v), want *TooManyResultsError", err, err)
		}
	})

	t.Run("propagates fetch error", func(t *testing.T) {
		want := errors.New("network down")
		_, err := fetchAllPages(func(offset int) ([]models.ItemBrief, error) {
			return nil, want
		}, searchMaxPages)
		if !errors.Is(err, want) {
			t.Fatalf("err = %v, want %v", err, want)
		}
	})
}

func TestBuildThorURL(t *testing.T) {
	got, err := buildThorURL(
		"decorate/search.itemList/1.0",
		map[string]interface{}{"shopId": shopID, "key": "test"},
		"token123",
		map[string]string{"_": "1234567890"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got, "https://thor.weidian.com/decorate/search.itemList/1.0?") {
		t.Fatalf("unexpected url prefix: %s", got)
	}
	if !strings.Contains(got, "wdtoken=token123") {
		t.Fatalf("missing wdtoken in url: %s", got)
	}
	if !strings.Contains(got, "_=1234567890") {
		t.Fatalf("missing dash in url: %s", got)
	}
}

func TestDoGET(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if got := r.Header.Get("referer"); got != "https://example.com/" {
				t.Errorf("referer = %q", got)
			}
			if got := r.Header.Get("User-Agent"); got == "" {
				t.Error("missing User-Agent")
			}
			w.Write([]byte(`{"ok":true}`))
		}))
		defer srv.Close()

		orig := httpClient
		httpClient = srv.Client()
		t.Cleanup(func() { httpClient = orig })

		body, err := doGET(srv.URL, "https://example.com/")
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != `{"ok":true}` {
			t.Fatalf("body = %q", body)
		}
	})

	t.Run("non-200 status", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
		}))
		defer srv.Close()

		orig := httpClient
		httpClient = srv.Client()
		t.Cleanup(func() { httpClient = orig })

		_, err := doGET(srv.URL, "https://example.com/")
		if err == nil || err.Error() != "HTTP 502" {
			t.Fatalf("err = %v, want HTTP 502", err)
		}
	})
}

func TestShopReferer(t *testing.T) {
	want := fmt.Sprintf("https://shop%s.v.weidian.com/", shopID)
	if got := shopReferer(); got != want {
		t.Fatalf("shopReferer() = %q, want %q", got, want)
	}
}
