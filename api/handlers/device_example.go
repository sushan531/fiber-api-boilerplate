package handlers

import (
	"fiber-api/api/middleware"

	"github.com/gofiber/fiber/v2"
)

// DeviceInfoHandler demonstrates how to use device detection in handlers
func DeviceInfoHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		deviceType := middleware.GetDeviceType(c)

		// Get additional device information from context
		userAgent := c.Locals("user_agent")
		deviceFamily := c.Locals("device_family")
		osFamily := c.Locals("os_family")
		browserFamily := c.Locals("browser_family")

		response := fiber.Map{
			"device_type":    deviceType,
			"user_agent":     userAgent,
			"device_family":  deviceFamily,
			"os_family":      osFamily,
			"browser_family": browserFamily,
		}

		// Device-specific logic
		switch deviceType {
		case middleware.DeviceTypeAndroid:
			response["message"] = "Welcome Android app user!"
			response["features"] = []string{"push_notifications", "biometric_auth", "offline_mode"}
		case middleware.DeviceTypeIOS:
			response["message"] = "Welcome iOS app user!"
			response["features"] = []string{"push_notifications", "face_id", "apple_pay"}
		case middleware.DeviceTypeWeb:
			response["message"] = "Welcome web user!"
			response["features"] = []string{"desktop_notifications", "file_upload", "advanced_ui"}
		}

		return c.JSON(response)
	}
}

// LoginResponseHandler shows device-specific login responses
func LoginResponseHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		deviceType := middleware.GetDeviceType(c)

		// Simulate different response formats based on device
		baseResponse := fiber.Map{
			"success": true,
			"message": "Login successful",
		}

		switch deviceType {
		case middleware.DeviceTypeAndroid:
			baseResponse["android_config"] = fiber.Map{
				"push_endpoint": "/api/push/android",
				"update_check":  "/api/version/android",
			}
		case middleware.DeviceTypeIOS:
			baseResponse["ios_config"] = fiber.Map{
				"push_endpoint": "/api/push/ios",
				"update_check":  "/api/version/ios",
			}
		case middleware.DeviceTypeWeb:
			baseResponse["web_config"] = fiber.Map{
				"websocket_url": "wss://api.example.com/ws",
				"cdn_base":      "https://cdn.example.com",
			}
		}

		return c.JSON(baseResponse)
	}
}
