package healthlake

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Call(ctx context.Context, httpClient HTTPClient, req *http.Request, response interface{}) error {
	var client HTTPClient
	if httpClient != nil {
		client = httpClient
	} else {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, response); err != nil {
		return err
	}

	fmt.Println(response)

	return nil
}
