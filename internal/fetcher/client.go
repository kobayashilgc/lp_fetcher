package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	shopID         = "1404154952"
	httpTimeout    = 30 * time.Second
)

var httpClient = &http.Client{Timeout: httpTimeout}

func shopReferer() string {
	return fmt.Sprintf("https://shop%s.v.weidian.com/", shopID)
}

func buildThorURL(apiPath string, param map[string]interface{}, wdtoken string, extraQuery map[string]string) (string, error) {
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("param", string(paramBytes))
	values.Set("wdtoken", wdtoken)
	for key, val := range extraQuery {
		values.Set(key, val)
	}

	return fmt.Sprintf("https://thor.weidian.com/%s?%s", apiPath, values.Encode()), nil
}

func doGET(reqURL, referer string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("referer", referer)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	return body, nil
}
