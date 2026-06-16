package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	exp1 = regexp.MustCompile(`^(?:https?://)?github\.com/([^/]+)/([^/]+)/(?:releases|archive)/.*$`)
	exp2 = regexp.MustCompile(`^(?:https?://)?github\.com/([^/]+)/([^/]+)/(?:blob|raw)/.*$`)
	exp3 = regexp.MustCompile(`^(?:https?://)?github\.com/([^/]+)/([^/]+)/(?:info|git-).*$`)
	exp4 = regexp.MustCompile(`^(?:https?://)?raw\.(?:githubusercontent|github)\.com/([^/]+)/([^/]+)/.+?/.+$`)
	exp5 = regexp.MustCompile(`^(?:https?://)?gist\.(?:githubusercontent|github)\.com/([^/]+)/.+?/.+$`)

	patterns = []*regexp.Regexp{exp1, exp2, exp3, exp4, exp5}
)

func checkURL(u string) (int, []string) {
	for i, re := range patterns {
		m := re.FindStringSubmatch(u)
		if m != nil {
			return i + 1, m
		}
	}
	return 0, nil
}

func cleanURL(u string) string {
	if !strings.HasPrefix(u, "http") {
		u = "https://" + u
	}
	if strings.HasPrefix(u, "https:/") && !strings.HasPrefix(u, "https://") {
		u = "https://" + u[7:]
	}
	return u
}

func jsDelivrRedirect(rawURL string, expID int) (string, bool) {
	if expID == 2 {
		newURL := strings.Replace(rawURL, "/blob/", "@", 1)
		newURL = strings.Replace(newURL, "github.com", "cdn.jsdelivr.net/gh", 1)
		return newURL, true
	}
	if expID == 4 {
		u, err := url.Parse(rawURL)
		if err != nil {
			return "", false
		}
		parts := strings.SplitN(strings.TrimPrefix(u.Path, "/"), "/", 4)
		if len(parts) < 4 {
			return "", false
		}
		newPath := fmt.Sprintf("/gh/%s/%s@%s/%s", parts[0], parts[1], parts[2], parts[3])
		newURL := fmt.Sprintf("https://cdn.jsdelivr.net%s", newPath)
		if u.RawQuery != "" {
			newURL += "?" + u.RawQuery
		}
		return newURL, true
	}
	return "", false
}

func blobToRaw(u string) string {
	return strings.Replace(u, "/blob/", "/raw/", 1)
}

func matchEntry(list []Entry, user, repo string) bool {
	for _, e := range list {
		if e.User == "*" || e.User == user {
			if e.Repo == "*" || e.Repo == repo {
				return true
			}
		}
	}
	return false
}

func Handler(c *gin.Context) {
	u := c.Param("path")
	u = strings.TrimPrefix(u, "/")

	if u == "favicon.ico" {
		if favicon != nil {
			c.Data(200, "image/vnd.microsoft.icon", favicon)
		} else {
			c.Status(204)
		}
		return
	}

	if u == "" {
		if c.Request.URL.Query().Get("q") != "" {
			c.Redirect(302, "/"+c.Query("q"))
			return
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(200, indexHTML)
		return
	}

	if !strings.Contains(u, "://") {
		u = "https://github.com/" + u
	} else {
		u = cleanURL(u)
	}

	expID, matched := checkURL(u)
	if matched == nil {
		c.String(403, "Invalid input.")
		return
	}

	user := matched[1]
	var repo string
	if len(matched) >= 3 {
		repo = matched[2]
	}

	if len(cfg.WhiteList) > 0 && !matchEntry(cfg.WhiteList, user, repo) {
		c.String(403, "Forbidden by white list.")
		return
	}
	if len(cfg.BlackList) > 0 && matchEntry(cfg.BlackList, user, repo) {
		c.String(403, "Forbidden by black list.")
		return
	}

	passBy := len(cfg.PassList) > 0 && matchEntry(cfg.PassList, user, repo)

	if cfg.JsDelivr {
		if expID == 2 {
			newURL, ok := jsDelivrRedirect(u, 2)
			if ok {
				c.Redirect(302, newURL)
				return
			}
		}
		if expID == 4 {
			newURL, ok := jsDelivrRedirect(u, 4)
			if ok {
				c.Redirect(302, newURL)
				return
			}
		}
	}

	if passBy {
		if expID == 2 {
			newURL, ok := jsDelivrRedirect(u, 2)
			if ok {
				c.Redirect(302, newURL)
				return
			}
		}
		if expID == 4 {
			newURL, ok := jsDelivrRedirect(u, 4)
			if ok {
				c.Redirect(302, newURL)
				return
			}
		}
	}

	if passBy {
		if expID == 2 {
			u = blobToRaw(u)
		}
		targetURL := u
		if c.Request.URL.RawQuery != "" {
			targetURL = u + "?" + c.Request.URL.RawQuery
		}
		c.Redirect(302, targetURL)
		return
	}

	if expID == 2 {
		u = blobToRaw(u)
	}
	proxyHandler(c, u)
}
