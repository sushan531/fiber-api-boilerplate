package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ua-parser/uap-go/uaparser"
)

// DeviceType represents the type of device making the request
type DeviceType string

const (
	DeviceTypeWeb     DeviceType = "web"
	DeviceTypeAndroid DeviceType = "android"
	DeviceTypeIOS     DeviceType = "ios"
)

// DeviceDetectionMiddleware parses User-Agent to determine device type
// Web includes browsers on any platform (desktop, mobile web browsers)
// Android/iOS only for native apps, not mobile browsers
func DeviceDetectionMiddleware() fiber.Handler {
	parser := uaparser.NewFromSaved()

	return func(c *fiber.Ctx) error {
		userAgent := c.Get("User-Agent")
		if userAgent == "" {
			// Default to web if no User-Agent
			c.Locals("device_type", DeviceTypeWeb)
			return c.Next()
		}

		client := parser.Parse(userAgent)
		deviceType := determineDeviceType(client, userAgent)

		// Store device information in context
		c.Locals("device_type", deviceType)
		c.Locals("user_agent", userAgent)
		return c.Next()
	}
}

// determineDeviceType analyzes parsed user agent to determine device type
func determineDeviceType(client *uaparser.Client, userAgent string) DeviceType {
	osFamily := strings.ToLower(client.Os.Family)
	browserFamily := strings.ToLower(client.UserAgent.Family)
	userAgentLower := strings.ToLower(userAgent)

	// Check for native Android apps
	if osFamily == "android" {
		// If it's a browser on Android, treat as web
		if isBrowser(browserFamily) {
			return DeviceTypeWeb
		}
		// If it contains typical browser indicators, treat as web
		if containsBrowserIndicators(userAgentLower) {
			return DeviceTypeWeb
		}
		// Otherwise, it's likely a native Android app
		return DeviceTypeAndroid
	}

	// Check for native iOS apps
	if osFamily == "ios" {
		// If it's a browser on iOS, treat as web
		if isBrowser(browserFamily) {
			return DeviceTypeWeb
		}
		// If it contains typical browser indicators, treat as web
		if containsBrowserIndicators(userAgentLower) {
			return DeviceTypeWeb
		}
		// Otherwise, it's likely a native iOS app
		return DeviceTypeIOS
	}

	// Everything else is considered web (desktop browsers, other mobile browsers, etc.)
	return DeviceTypeWeb
}

// isBrowser checks if the user agent family indicates a web browser
func isBrowser(browserFamily string) bool {
	browserFamilyLower := strings.ToLower(browserFamily)
	browsers := []string{
		"chrome", "firefox", "safari", "edge", "opera", "internet explorer",
		"chrome mobile", "safari mobile", "firefox mobile", "opera mobile",
		"samsung internet", "uc browser", "mobile safari", "chrome mobile ios",
		"firefox ios", "opera mini", "android browser", "webview",
	}

	for _, browser := range browsers {
		if strings.Contains(browserFamilyLower, browser) {
			return true
		}
	}
	return false
}

// containsBrowserIndicators checks for common browser-related keywords in user agent
func containsBrowserIndicators(userAgent string) bool {
	indicators := []string{
		"mozilla", "webkit", "gecko", "chrome", "safari", "firefox",
		"browser", "webview", "mobile safari", "version/",
	}

	for _, indicator := range indicators {
		if strings.Contains(userAgent, indicator) {
			return true
		}
	}
	return false
}

// GetDeviceType is a helper function to extract device type from Fiber context
func GetDeviceType(c *fiber.Ctx) DeviceType {
	if deviceType, ok := c.Locals("device_type").(DeviceType); ok {
		return deviceType
	}
	return DeviceTypeWeb // Default fallback
}
