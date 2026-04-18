package handlers

import "strings"

func (apicfg *ApiConfig) publicAssetURL(key string) string {
	if key == "" {
		return ""
	}
	if strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return key
	}
	if apicfg.S3BaseURL == "" {
		return key
	}
	return strings.TrimRight(apicfg.S3BaseURL, "/") + "/" + strings.TrimLeft(key, "/")
}
