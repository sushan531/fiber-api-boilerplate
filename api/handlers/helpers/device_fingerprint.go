package helpers

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/ua-parser/uap-go/uaparser"
)

// DeviceFingerprint represents a device fingerprint based on User-Agent
type DeviceFingerprint struct {
	Hash     string `json:"hash"`
	Platform string `json:"platform"`
	Browser  string `json:"browser"`
	Version  string `json:"version"`
}

// GenerateDeviceFingerprint creates a device fingerprint from User-Agent string
func GenerateDeviceFingerprint(userAgent string) *DeviceFingerprint {
	if userAgent == "" {
		// Return a default fingerprint for empty user agents
		return &DeviceFingerprint{
			Hash:     generateHash("unknown"),
			Platform: "unknown",
			Browser:  "unknown",
			Version:  "unknown",
		}
	}

	parser := uaparser.NewFromSaved()
	client := parser.Parse(userAgent)

	// Extract key components for fingerprinting
	platform := normalizeString(client.Os.Family)
	browser := normalizeString(client.UserAgent.Family)
	version := normalizeString(client.UserAgent.Major)

	// Create a composite string for hashing
	// Include OS family, browser family, and major version
	composite := fmt.Sprintf("%s|%s|%s", platform, browser, version)

	// Generate SHA256 hash of the composite string
	hash := generateHash(composite)

	return &DeviceFingerprint{
		Hash:     hash,
		Platform: platform,
		Browser:  browser,
		Version:  version,
	}
}

// ValidateDeviceFingerprint compares current User-Agent against stored fingerprint
func ValidateDeviceFingerprint(userAgent string, storedHash string) bool {
	if storedHash == "" {
		return false
	}

	currentFingerprint := GenerateDeviceFingerprint(userAgent)
	return currentFingerprint.Hash == storedHash
}

// generateHash creates a SHA256 hash of the input string
func generateHash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}

// normalizeString normalizes strings for consistent fingerprinting
func normalizeString(s string) string {
	if s == "" {
		return "unknown"
	}
	return strings.ToLower(strings.TrimSpace(s))
}

// GetFingerprintFromClaims extracts device fingerprint from JWT claims
func GetFingerprintFromClaims(claims map[string]interface{}) (string, bool) {
	fingerprint, exists := claims["device_fingerprint"]
	if !exists {
		return "", false
	}

	fingerprintStr, ok := fingerprint.(string)
	return fingerprintStr, ok
}
