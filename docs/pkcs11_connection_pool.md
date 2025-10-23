# PKCS11 Connection Pool

## Overview

This implementation adds a connection pool for PKCS11 sessions to improve performance and resource management in the SSM (Simple Secret Manager) application.

## Benefits

### Before (Per-Request Connections)
```go
// Each API call creates a new session
mgr, err := pkcs11mgr.New(...)
err = mgr.OpenSession()
defer mgr.CloseSession()
defer mgr.Finalize()
```

**Problems:**
- High latency due to session creation overhead
- Resource waste creating/destroying sessions
- Potential bottleneck under high load

### After (Connection Pool)
```go
// Reuse existing sessions from pool
err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
    // Use the pooled connection
    return nil
})
```

**Benefits:**
- ⚡ **Faster response times** - No session creation overhead
- 💾 **Better resource usage** - Reuse existing connections  
- 🔄 **Concurrent support** - Multiple requests can share pool
- 🛡️ **Automatic cleanup** - Expired connections are recycled
- 📊 **Monitoring** - Pool statistics available

## Configuration

The pool is configured in `server/server.go`:

```go
poolConfig := pkcs11mgr.DefaultPoolConfig()
poolConfig.PkcsPath = factory.SsmConfig.Configuration.PkcsPath
poolConfig.SlotNumber = uint(factory.SsmConfig.Configuration.LotsNumber)
poolConfig.Pin = factory.SsmConfig.Configuration.Pin
poolConfig.MaxSize = 10  // Maximum connections
poolConfig.MinSize = 2   // Minimum connections to maintain
```

### Configuration Options

| Parameter | Description | Default |
|-----------|-------------|---------|
| `MaxSize` | Maximum connections in pool | 10 |
| `MinSize` | Minimum connections to maintain | 2 |
| `IdleTimeout` | How long idle connections live | 5 minutes |
| `MaxLifetime` | Maximum connection lifetime | 30 minutes |

## API Endpoints

### Connection Pool Statistics
```http
GET /pool/stats
```

**Response:**
```json
{
    "max_size": 10,
    "current_size": 5,
    "available": 3,
    "in_use": 2
}
```

### Pool-Optimized Encrypt Endpoint
```http
POST /encrypt-pool
```

Uses the connection pool instead of creating new sessions.

## Usage Examples

### Using the Pool in Handlers

```go
// Get connection from pool, use it, return it automatically
err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
    // Your PKCS11 operations here
    keyHandle, err := mgr.FindKeyLabelReturnRandom(keyLabel)
    if err != nil {
        return err
    }
    
    ciphertext, err := mgr.EncryptKey(keyHandle, iv, plaintext, mechanism)
    return err
})
```

### Manual Pool Management

```go
// Get connection from pool
mgr, err := pool.Get()
if err != nil {
    return err
}
defer pool.Put(mgr) // Important: return to pool

// Use connection
result, err := mgr.SomeOperation()
```

## Monitoring

Monitor pool health using the stats endpoint:

```bash
curl http://localhost:8080/pool/stats
```

Key metrics to watch:
- **High in_use**: May need to increase pool size
- **Low available**: Pool is at capacity
- **current_size < min_size**: Connections are failing

## Implementation Notes

### Thread Safety
- Pool is fully thread-safe using channels and mutexes
- Multiple goroutines can safely get/put connections

### Connection Lifecycle
1. **Creation**: Connections created on-demand up to MaxSize
2. **Validation**: Connections checked for expiration before use
3. **Cleanup**: Expired connections automatically closed and removed

### Error Handling
- Failed connections are automatically removed from pool
- Pool degrades gracefully when HSM is unavailable
- Timeouts prevent indefinite blocking

## Migration Guide

### Updating Existing Handlers

**Old Way:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    mgr, err := pkcs11mgr.New(...)
    if err != nil { return }
    
    err = mgr.OpenSession()
    if err != nil { return }
    
    defer mgr.CloseSession()
    defer mgr.Finalize()
    
    // Use mgr...
}
```

**New Way:**
```go
func handler(w http.ResponseWriter, r *http.Request) {
    err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
        // Use mgr...
        return nil
    })
    if err != nil { return }
}
```

## Performance Impact

Expected improvements:
- **Latency**: 50-80% reduction in response time
- **Throughput**: 2-3x increase in concurrent requests
- **Resource Usage**: Significant reduction in HSM load

## Troubleshooting

### Pool Exhaustion
```
Error: timeout waiting for available connection
```
**Solution**: Increase `MaxSize` or check for connection leaks

### Connection Failures
```
Error: failed to create connection
```
**Solution**: Check HSM connectivity and credentials

### High Memory Usage
**Solution**: Reduce `MaxSize` or `MaxLifetime`

## Future Enhancements

1. **Dynamic Pool Sizing**: Auto-adjust based on load
2. **Health Checks**: Periodic connection validation
3. **Metrics Export**: Prometheus/metrics integration
4. **Circuit Breaker**: Fail-fast when HSM is down