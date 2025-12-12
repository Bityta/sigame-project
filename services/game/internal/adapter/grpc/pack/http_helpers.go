package pack

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func buildURL(baseURL, path string, args ...interface{}) string {
	formattedPath := fmt.Sprintf(path, args...)
	return fmt.Sprintf("%s%s", baseURL, formattedPath)
}

func doGetRequest(ctx context.Context, client *http.Client, url string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, ErrCreateRequest(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, ErrExecuteRequest(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, ErrReadResponse(err)
	}

	return body, resp.StatusCode, nil
}

