package main

import (
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	User string
	Repo string
}

type Config struct {
	Host      string
	Port      string
	AssetURL  string
	JsDelivr  bool
	SizeLimit int64
	WhiteList []Entry
	BlackList []Entry
	PassList  []Entry
	Debug     bool
}

var cfg *Config

func loadConfig() *Config {
	c := &Config{
		Host:     getEnv("HOST", "0.0.0.0"),
		Port:     getEnv("PORT", "8080"),
		AssetURL: getEnv("ASSET_URL", "https://hunshcn.github.io/gh-proxy"),
	}

	if v := os.Getenv("JSDELIVR"); v == "1" || v == "true" || v == "yes" {
		c.JsDelivr = true
	}

	if v := os.Getenv("SIZE_LIMIT"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			c.SizeLimit = n
		}
	}

	if v := os.Getenv("DEBUG"); v == "1" || v == "true" || v == "yes" {
		c.Debug = true
	}

	c.WhiteList = parseEntryList(os.Getenv("WHITE_LIST"))
	c.BlackList = parseEntryList(os.Getenv("BLACK_LIST"))
	c.PassList = parseEntryList(os.Getenv("PASS_LIST"))

	return c
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func parseEntryList(s string) []Entry {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	var entries []Entry
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.Contains(line, "/") {
			parts := strings.SplitN(line, "/", 2)
			entries = append(entries, Entry{User: parts[0], Repo: parts[1]})
		} else {
			entries = append(entries, Entry{User: line, Repo: "*"})
		}
	}
	return entries
}
