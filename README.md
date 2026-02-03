# UniFi Protect → Pushover Webhook Bridge

A Go application that bridges UniFi Protect camera alerts to Pushover emergency notifications.

## Features

- Receives webhook POST requests from UniFi Protect Alarm Manager
- Sends emergency priority notifications to Pushover (bypasses DND)
- Configurable retry interval and expiration time
- Health check endpoint for Kubernetes probes
- Builds with `ko` for easy container deployment

## Configuration

### Required Environment Variables

| Variable | Description |
|----------|-------------|
| `PUSHOVER_APP_TOKEN` | Pushover application token |
| `PUSHOVER_USER_KEY` | Pushover user/group key |

### Optional Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `PUSHOVER_PRIORITY` | `2` | Priority level (2 = emergency) |
| `PUSHOVER_RETRY` | `60` | Retry interval in seconds |
| `PUSHOVER_EXPIRE` | `3600` | Expiration time in seconds |
| `LOG_LEVEL` | `info` | Logging level |

## Endpoints

- `POST /webhook` - Receives UniFi Protect webhook payloads
- `GET /health` - Health check endpoint

## Building

### With ko

```bash
ko build .
```

### With Go

```bash
go build -o unifi-pushover-bridge .
```

## Deployment

### Kubernetes

1. Update the secret in `k8s/secret.yaml` with your Pushover credentials
2. Apply the manifests:

```bash
kubectl apply -f k8s/
```

Or with ko:

```bash
ko apply -f k8s/
```

## Testing

```bash
# Run tests
go test ./...

# Manual test with curl
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

## UniFi Protect Setup

1. Go to UniFi Protect settings → Alarm Manager
2. Create a new alarm with desired triggers (person, vehicle, motion, etc.)
3. Add a webhook action pointing to `http://your-bridge-address:8080/webhook`
4. Method: POST

## License

MIT
