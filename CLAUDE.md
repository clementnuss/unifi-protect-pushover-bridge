# UniFi Protect → Pushover Webhook Bridge

## Context

Building a Go application to bridge UniFi Protect camera alerts to Pushover emergency notifications. This will wake up the user when motion/person is detected.

## Project Requirements

### Functionality
- Receive webhook POST requests from UniFi Protect Alarm Manager
- Parse UniFi Protect webhook JSON payload
- Send emergency priority notifications to Pushover (Priority 2)
- Support for emergency alert features:
  - Bypasses Do Not Disturb mode
  - Repeats alert every N seconds until acknowledged
  - Expires after configured time

### Deployment
- Container image built with `ko` (à-la https://github.com/clementnuss/truckflow-user-importer)
- Deploy to Kubernetes
- Configuration via environment variables

## UniFi Protect Webhook Payload Structure

UniFi Protect sends POST requests with JSON payload like:

```json
{
  "alarm": {
    "name": "Front Door Person Detected",
    "sources": [{
      "device": "F4E2C60E6104",
      "type": "include"
    }],
    "conditions": [
      {
        "condition": {
          "type": "is",
          "source": "person"
        }
      }
    ],
    "triggers": [
      {
        "key": "person",
        "device": "F4E2C60E6104"
      }
    ]
  },
  "timestamp": 1725883107267
}
```

Key fields to extract:
- `alarm.name` - Human readable alarm name
- `alarm.triggers[].key` - Type of detection (person, vehicle, animal, motion, etc.)
- `timestamp` - Unix timestamp in milliseconds

## Pushover API Integration

### Go Library
Use [`github.com/gregdel/pushover`](https://github.com/gregdel/pushover) - well-maintained, supports emergency priority

### Emergency Priority Configuration
```go
message := pushover.NewMessage("Alert message")
message.Priority = pushover.PriorityEmergency
message.Retry = 60 * time.Second     // Retry every 60 seconds
message.Expire = 3600 * time.Second  // Expire after 1 hour
```

### Required Pushover Credentials
- `PUSHOVER_APP_TOKEN` - Application API token (register at https://pushover.net/apps/build)
- `PUSHOVER_USER_KEY` - User/group key to send notifications to

## Technical Stack

### Dependencies
```go
require (
    github.com/gregdel/pushover v1.3.0
    // Standard library for HTTP server
)
```

### Application Structure
```
unifi-pushover-bridge/
├── main.go              # Main application entry point
├── handlers/
│   └── webhook.go       # HTTP handler for webhook endpoint
├── notifier/
│   └── pushover.go      # Pushover notification logic
├── types/
│   └── unifi.go         # UniFi webhook payload structs
├── config/
│   └── config.go        # Environment variable configuration
├── Dockerfile           # Multi-stage Dockerfile (if not using ko)
├── .ko.yaml             # ko configuration
├── go.mod
├── go.sum
└── README.md
```

## Configuration via Environment Variables

Required:
- `PUSHOVER_APP_TOKEN` - Pushover application token
- `PUSHOVER_USER_KEY` - Pushover user/group key
- `PORT` - HTTP server port (default: 8080)

Optional:
- `PUSHOVER_PRIORITY` - Priority level (default: 2 for emergency)
- `PUSHOVER_RETRY` - Retry interval in seconds (default: 60)
- `PUSHOVER_EXPIRE` - Expiration time in seconds (default: 3600)
- `LOG_LEVEL` - Logging level (default: info)

## Implementation Tasks

### 1. Core HTTP Server
- Create HTTP server listening on `/webhook` endpoint
- Accept POST requests only
- Parse JSON body into UniFi webhook struct
- Add basic request validation (check for required fields)
- Health check endpoint on `/health`

### 2. UniFi Webhook Parser
- Define Go structs matching UniFi webhook payload
- Parse JSON payload
- Extract relevant information:
  - Alarm name
  - Trigger type
  - Timestamp
  - Device ID (optional)

### 3. Pushover Notifier
- Initialize Pushover client with credentials
- Construct notification message from webhook data
- Format: "🚨 [Alarm Name] - [Trigger Type] detected at [Time]"
- Set emergency priority with configurable retry/expire
- Handle API errors gracefully
- Log successful/failed notifications

### 4. Error Handling & Logging
- Structured logging (use `log/slog` from stdlib)
- Log incoming webhook requests
- Log Pushover API responses
- Handle and log errors appropriately
- Return appropriate HTTP status codes

### 5. Ko Configuration
Create `.ko.yaml`:
```yaml
defaultBaseImage: gcr.io/distroless/static:nonroot
builds:
  - id: unifi-pushover-bridge
    main: .
    env:
      - CGO_ENABLED=0
```

### 6. Kubernetes Manifests
Create basic deployment and service manifests:
- Deployment with environment variables from ConfigMap/Secret
- Service exposing the webhook endpoint
- Optional: Ingress for external access

## Example Usage Flow

1. User configures UniFi Protect Alarm Manager:
   - Trigger: Person detected on Front Door camera
   - Action: Webhook → `http://unifi-pushover-bridge:8080/webhook`
   - Method: POST

2. Person detected → UniFi sends webhook POST

3. Bridge receives request:
   - Parses JSON payload
   - Extracts: "Front Door Person Detected" alarm
   - Constructs message: "🚨 Front Door Person Detected - person detected at 02:15 AM"

4. Bridge sends to Pushover:
   - Priority 2 (Emergency)
   - Retry every 60s
   - Expire after 1 hour

5. User's phone:
   - Alarm sounds even in DND mode
   - Shows notification requiring acknowledgment
   - Repeats every 60s until acknowledged

## Security Considerations

1. Store Pushover tokens as Kubernetes Secrets
2. Consider adding authentication to webhook endpoint (optional):
   - Bearer token
   - IP whitelist
3. HTTPS termination at ingress level
4. Rate limiting to prevent abuse

## Testing Strategy

1. Unit tests for webhook parsing
2. Mock Pushover client for testing notification logic
3. Integration test with test webhook from UniFi
4. Manual E2E test with actual camera detection

## References

- [Pushover API Documentation](https://pushover.net/api)
- [Pushover Go Library](https://github.com/gregdel/pushover)
- [UniFi Protect Webhooks Documentation](https://help.ui.com/hc/en-us/articles/25478744592023)
- [Example UniFi webhook integration](https://github.com/patfelst/Unifi-webhook-integration-with-Home-Assistant)
- [Ko documentation](https://ko.build/)
- [Reference project structure](https://github.com/clementnuss/truckflow-user-importer)

## Next Steps in Claude Code

1. Initialize Go module: `go mod init github.com/yourusername/unifi-pushover-bridge`
2. Create basic project structure
3. Implement HTTP server with `/webhook` and `/health` endpoints
4. Add UniFi webhook payload parsing
5. Integrate Pushover notification sending
6. Add configuration from environment variables
7. Test with mock webhook payloads
8. Create ko configuration
9. Create Kubernetes manifests
10. Document deployment instructions

## Sample Test Webhook

For testing, you can send:
```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "alarm": {
      "name": "Front Door Person",
      "triggers": [{"key": "person"}]
    },
    "timestamp": 1725883107267
  }'
```

## Desired Output

A working Go application that:
- ✅ Receives UniFi Protect webhooks
- ✅ Sends emergency Pushover notifications
- ✅ Deploys easily to Kubernetes with ko
- ✅ Configurable via environment variables
- ✅ Production-ready with proper error handling and logging
