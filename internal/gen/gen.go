// downloader for syukujitsu.csv
// https://www8.cao.go.jp/chosei/shukujitsu/gaiyou.html

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := _main(); err != nil {
		log.Fatal(err)
	}
}

func _main() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if _, err := download(ctx); err != nil {
		return err
	}
	return nil
}

func download(ctx context.Context) ([]byte, error) {
	const csvURL = "https://www8.cao.go.jp/chosei/shukujitsu/syukujitsu.csv"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, csvURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "https://github.com/shogo82148/holidays-jp")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// save the raw data
	if err := os.WriteFile("syukujitsu.csv", buf, 0644); err != nil {
		return nil, err
	}

	return buf, nil
}
