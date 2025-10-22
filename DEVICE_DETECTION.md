# Device Fingerprinting Implementation

This document outlines the device fingerprinting mechanism implemented in the Fiber API to prevent cross-device token usage and enhance security.

## Overview

The device fingerprinting system creates unique fingerprints based on User-Agent headers and embeds them in JWT tokens. This prevents tokens generated on one device from being used on another device, even if the token is compromised.

## Key Features

- **Device-Specific Tokens**: Each token is bound to the device that requested it
- **Cross-Device Protection**: Tokens cannot be used from different devices
- **Automatic Validation**: Every API request validates the device fingerprint
- **Refresh Token Security**: Refresh tokens also validate device fingerprints

## Implementation

### Device Fingerprint Generation

Located in `api/handlers/helpers/device_fingerprint.go`:

```go
type DeviceFingerprint struct {
    Hash     string `json:"hash"`      // SHA256 hash of device characteristics
    Platform string `json:"platform"`  // Operating system family
    Browser  string `json:"browser"`   // Browser family
    Version  string `json:"version"`   // Browser major version
}
```

The fingerprint is generated from:
- **OS Family**: Windows, macOS, Linux, Android, iOS, etc.
- **Browser Family**: Chrome, Firefox, Safari, Edge, etc.
- **Browser Version**: Major version number only

### JWT Claims Enhancement

JWT tokens now include a `device_fingerprint` claim:

```go
type JWTClaims struct {
    UserID            string `json:"user_id"`
    UserEmail         string `json:"user_email"`
    Role              string `json:"role"`
    DeviceFingerprint string `json:"device_fingerprint"` // New field
}
```

### Authentication Flow

#### 1. Login Process
1. Extract User-Agent from request headers
2. Generate device fingerprint hash
3. Include fingerprint in JWT claims
4. Return tokens bound to the device

#### 2. Token Validation
1. Extract fingerprint from JWT claims
2. Generate fingerprint from current User-Agent
3. Compare fingerprints - reject if different
4. Allow request if fingerprints match

#### 3. Token Refresh
1. Validate refresh token's device fingerprint
2. Ensure current device matches original device
3. Generate new tokens with same fingerprint
4. Reject refresh if device doesn't match

## Security Benefits

### Cross-Device Protection
- Stolen tokens cannot be used from different devices
- Even if an attacker has the token, they need the exact same User-Agent

### Session Integrity
- Each device maintains its own token session
- No token sharing between devices, even for the same user

### Audit Trail
- Device mismatches are logged for security monitoring
- Failed attempts include device information

## Error Responses

### Invalid Device Fingerprint
```json
{
  "error": "Token cannot be used from this device"
}
```

### Missing User-Agent
```json
{
  "error": "Missing User-Agent header"
}
```

### Missing Fingerprint in Token
```json
{
  "error": "Invalid token: missing device fingerprint"
}
```

## Testing Device Fingerprinting

### Test Different Devices

```bash
# Chrome on Windows
curl -H "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36" \
     -H "Authorization: Bearer <token>" \
     http://localhost:3000/api/protected

# Firefox on macOS  
curl -H "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0" \
     -H "Authorization: Bearer <token>" \
     http://localhost:3000/api/protected

# Mobile Safari on iOS
curl -H "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1" \
     -H "Authorization: Bearer <token>" \
     http://localhost:3000/api/protected
```

### Expected Behavior
- Token generated on Chrome/Windows will only work with Chrome/Windows User-Agent
- Same token will be rejected when used with Firefox/macOS User-Agent
- Each device needs its own login to get device-specific tokens

## Configuration

No additional configuration required. The system:
- Automatically generates fingerprints for all login requests
- Validates fingerprints on every protected route
- Handles missing User-Agents gracefully (assigns "unknown" fingerprint)

## Files Modified

### Core Implementation
- `api/models/auth.go` - Added `DeviceFingerprint` field to JWT claims
- `api/handlers/helpers/device_fingerprint.go` - New fingerprinting utilities
- `api/handlers/helpers/helpers.go` - Updated `CreateJWTClaims` function
- `api/handlers/auth.go` - Updated login and refresh handlers
- `api/middleware/jwk_middleware.go` - Added fingerprint validation

## Limitations

### User-Agent Spoofing
- Attackers can potentially spoof User-Agent headers
- However, this requires knowledge of the exact User-Agent string
- Most automated attacks won't have this information

### Browser Updates
- Major browser updates may change User-Agent strings
- Users may need to re-login after significant browser updates
- Only major version changes affect fingerprints (minor updates are ignored)

### Multiple Browser Sessions
- Each browser/device combination needs separate login
- Users cannot share tokens between Chrome and Firefox on same device
- This is intentional for enhanced security

## Migration Notes

### Existing Tokens
- Tokens generated before this implementation will be rejected
- All users need to re-login after deployment
- Consider implementing a grace period if needed

### Backward Compatibility
- No database schema changes required
- Existing user accounts work without modification
- Only token generation and validation logic changed