package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func proxyHandler(c *gin.Context, targetBaseURL string) {
	fullURL := targetBaseURL
	if c.Request.URL.RawQuery != "" {
		fullURL = targetBaseURL + "?" + c.Request.URL.RawQuery
	}

	if strings.HasPrefix(fullURL, "https:/") && !strings.HasPrefix(fullURL, "https://") {
		fullURL = "https://" + fullURL[7:]
	}

	outReq, err := http.NewRequest(c.Request.Method, fullURL, c.Request.Body)
	if err != nil {
		c.String(500, "server error: "+err.Error())
		return
	}

	for k, vs := range c.Request.Header {
		if k != "Host" {
			for _, v := range vs {
				outReq.Header.Add(k, v)
			}
		}
	}

	transport := &http.Transport{
		DisableCompression: true,
	}

	client := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(outReq)
	if err != nil {
		c.String(502, "proxy error: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if cfg.SizeLimit > 0 && resp.ContentLength > cfg.SizeLimit {
		c.Redirect(302, fullURL)
		return
	}

	for k, vs := range resp.Header {
		for _, v := range vs {
			c.Header(k, v)
		}
	}

	c.Header("Content-Security-Policy", "")
	c.Header("Content-Security-Policy-Report-Only", "")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Expose-Headers", "*")
	c.Header("Transfer-Encoding", "")

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		location := resp.Header.Get("Location")
		if location != "" {
			if expID, _ := checkURL(location); expID > 0 {
				c.Header("Location", "/"+location)
			} else {
				resp.Body.Close()
				c.Request.Body = nil
				proxyHandler(c, location)
				return
			}
		}
	}

	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
