package cookie_service

import (
	"strings"

	cookie_flag "github.com/Rfluid/insta-tools/src/cookie/flag"
	log_service "github.com/Rfluid/insta-tools/src/log/service"
	"github.com/pterm/pterm"
)

// ParseCookies converts the cookie string into a map[string]string
func ParseCookies() map[string]string {
	cookieMap := make(map[string]string)

	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		"Parsing cookies to map...",
	)

	// If Cookies is empty, return an empty map
	if cookie_flag.Cookies == "" {
		return cookieMap
	}

	// Split the cookie string by `; ` to separate each key-value pair
	pairs := strings.Split(cookie_flag.Cookies, "; ")
	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		"Executing loop to parse cookies...",
	)
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			cookieMap[parts[0]] = parts[1]
		}
	}

	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		"Parsed cookies to map.",
	)

	return cookieMap
}
