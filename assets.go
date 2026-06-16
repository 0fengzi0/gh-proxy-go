package main

import (
	"io"
	"net/http"
	"time"
)

var (
	indexHTML string
	favicon   []byte
)

func fetchAssets() {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(cfg.AssetURL)
	if err != nil {
		indexHTML = fallbackHTML
	} else {
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil || resp.StatusCode != 200 {
			indexHTML = fallbackHTML
		} else {
			indexHTML = string(data)
		}
	}

	resp, err = client.Get(cfg.AssetURL + "/favicon.ico")
	if err != nil {
		favicon = nil
	} else {
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil || resp.StatusCode != 200 {
			favicon = nil
		} else {
			favicon = data
		}
	}
}

var fallbackHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>gh-proxy-go</title>
</head>
<body>
    <h2>gh-proxy-go</h2>
    <p>Use <code>/https://github.com/user/repo</code> (full URL) or <code>/user/repo</code> (shorthand) to access resources.</p>
</body>
</html>`
