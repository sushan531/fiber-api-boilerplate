# Device Detection Middleware

This middleware uses the `uap-go` library to parse User-Agent headers and distinguish between web browsers, native Android apps, and native iOS apps.

## Device Types

- **Web**: Any browser on any platform (desktop, mobile web browsers)
- **Android**: Native Android applications only (not mobile browsers)
- **iOS**: Native iOS applications only (not mobile browsers)

## Usage

### Apply Middleware

```go
import "fiber-api/api/middleware"

// Apply to all routes in a group
app.Use("/api", middleware.DeviceDetectionMiddleware())

// Or apply to specific routes
app.Get("/endpoint", middleware.DeviceDetectionMiddleware(), handler)
```

### Access Device Information in Handlers

```go
func MyHandler() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Get device type
        deviceType := middleware.GetDeviceType(c)
        
        // Get additional device info
        userAgent := c.Locals("user_agent").(string)
        deviceFamily := c.Locals("device_family").(string)
        osFamily := c.Locals("os_family").(string)
        browserFamily := c.Locals("browser_family").(string)
        
        // Device-specific logic
        switch deviceType {
        case middleware.DeviceTypeAndroid:
            // Android app logic
        case middleware.DeviceTypeIOS:
            // iOS app logic  
        case middleware.DeviceTypeWeb:
            // Web browser logic
        }
        
        return c.JSON(fiber.Map{"device": deviceType})
    }
}
```

## Examples

### User Agent Classification

| User Agent | Device Type | Reason |
|------------|-------------|---------|
| `Mozilla/5.0 (Linux; Android 10) Chrome/91.0 Mobile Safari/537.36` | Web | Chrome browser on Android |
| `Mozilla/5.0 (iPhone; CPU iPhone OS 14_6) Safari/604.1` | Web | Safari browser on iOS |
| `MyApp/1.0 (Linux; Android 10; SM-G973F)` | Android | Native Android app |
| `MyApp/1.0 CFNetwork/1240.0.4 Darwin/20.6.0` | iOS | Native iOS app |
| `Mozilla/5.0 (Windows NT 10.0) Chrome/91.0 Safari/537.36` | Web | Desktop browser |

### Integration with Authentication

The middleware is integrated with the login handler to create device-specific session keys:

```go
// In LoginHandler
deviceType := middleware.GetDeviceType(c)
keyID, err := jwkManager.CreateSessionKey(userID, string(deviceType))
```

## Testing

Run the middleware tests:

```bash
go test ./api/middleware -v
```

## Context Variables

The middleware sets these variables in the Fiber context:

- `device_type`: DeviceType enum (web/android/ios)
- `user_agent`: Original User-Agent header
- `device_family`: Parsed device family
- `os_family`: Parsed OS family  
- `browser_family`: Parsed browser family