# SSM API Testing - Complete Setup

## 📁 Structure Created

```
ssm/
├── docs/
│   └── json_examples/
│       ├── requests/
│       │   ├── encrypt_request.json
│       │   ├── decrypt_request.json
│       │   ├── gen_aes_key_request.json
│       │   ├── gen_des_key_request.json
│       │   ├── gen_des3_key_request.json
│       │   ├── store_key_request.json
│       │   ├── get_data_keys_request.json
│       │   └── get_key_request.json
│       └── responses/
│           ├── encrypt_response.json
│           ├── decrypt_response.json
│           ├── gen_aes_key_response.json
│           ├── gen_des_key_response.json
│           ├── gen_des3_key_response.json
│           ├── store_key_response.json
│           ├── get_data_keys_response.json
│           ├── get_key_response.json
│           ├── get_all_keys_response.json
│           └── health_check_response.json
└── utils/
    ├── api_test.sh              # Bash script for Linux/macOS
    ├── api_test.ps1             # PowerShell script for Windows
    ├── README.md                # Complete documentation
    └── QUICK_EXAMPLES.md        # Quick reference examples
```

## 🚀 Quick Start

### Windows (PowerShell)

```powershell
# Navigate to utils folder
cd utils

# Test health check
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"

# Generate a key
.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json"

# Get all keys
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST"
```

### Linux/macOS (Bash)

```bash
# Navigate to utils folder
cd utils

# Make script executable
chmod +x api_test.sh

# Test health check
./api_test.sh /health-check GET

# Generate a key
./api_test.sh /generate-aes-key POST ../docs/json_examples/requests/gen_aes_key_request.json

# Get all keys
./api_test.sh /get-all-keys POST
```

## 📝 Features

### Both Scripts Support:

✅ All HTTP methods (GET, POST, PUT, DELETE, PATCH)
✅ JSON request body from file
✅ Custom base URL
✅ Pretty-printed JSON output
✅ HTTP status codes
✅ Request timing
✅ Color-coded success/error messages
✅ Help documentation

### PowerShell Script (`api_test.ps1`)

- Native Windows integration
- No external dependencies
- Parameter validation
- Detailed error messages

### Bash Script (`api_test.sh`)

- Cross-platform (Linux/macOS/WSL)
- Uses `curl` for requests
- Optional `jq` for JSON formatting
- Colorized terminal output

## 🎯 Common Use Cases

### 1. Test Complete Encryption Flow

```powershell
# Generate key
.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json"

# Encrypt data
.\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"

# Decrypt data
.\api_test.ps1 -Endpoint "/decrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\decrypt_request.json"
```

### 2. Key Management

```powershell
# Store custom key
.\api_test.ps1 -Endpoint "/store-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\store_key_request.json"

# Get all keys
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST"

# Get specific keys by label
.\api_test.ps1 -Endpoint "/get-data-keys" -Method "POST" -JsonFile "..\docs\json_examples\requests\get_data_keys_request.json"
```

### 3. Different Environments

```powershell
# Local development
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"

# Test environment
.\api_test.ps1 -Endpoint "/health-check" -Method "GET" -BaseUrl "http://test-server:8080"

# Production (with HTTPS)
.\api_test.ps1 -Endpoint "/health-check" -Method "GET" -BaseUrl "https://prod-server:8443"
```

## 📚 JSON Examples

All JSON examples are located in `../docs/json_examples/`:

### Request Examples

- **encrypt_request.json** - Encrypt plaintext data
- **decrypt_request.json** - Decrypt ciphertext
- **gen_aes_key_request.json** - Generate AES key (128/192/256 bits)
- **gen_des_key_request.json** - Generate DES key
- **gen_des3_key_request.json** - Generate DES3 key
- **store_key_request.json** - Store custom key in HSM
- **get_data_keys_request.json** - Query keys by label
- **get_key_request.json** - Get single key info

### Response Examples

Each request has a corresponding response example showing expected output format.

## 🔧 Customization

### Modify JSON Examples

Edit the JSON files in `docs/json_examples/requests/` to match your needs:

```json
{
  "key_label": "your-custom-label",
  "id": 123,
  "bits": 256
}
```

### Add New Endpoints

1. Create new JSON request/response files
2. Update the endpoint list in README.md
3. Test with the scripts

## 🐛 Troubleshooting

### PowerShell Execution Policy

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Bash Script Permissions

```bash
chmod +x api_test.sh
```

### Install curl (if missing)

```bash
# Ubuntu/Debian
sudo apt-get install curl

# macOS
brew install curl
```

### Install jq (optional, for pretty JSON)

```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq
```

## 📖 Full Documentation

For complete documentation, see:
- `README.md` - Full script documentation
- `QUICK_EXAMPLES.md` - Quick command reference

## 🎨 Output Examples

### Success Response

```
=== SSM API Request ===
URL:     http://localhost:8080/health-check
Method:  GET

=== Response ===
{
  "status": "OK",
  "message": "Service is healthy"
}

HTTP Status: 200
Time: 0.023s

✓ Request successful
```

### Error Response

```
=== SSM API Request ===
URL:     http://localhost:8080/invalid-endpoint
Method:  POST

=== Response ===
{
  "title": "Not Found",
  "detail": "The requested endpoint does not exist",
  "type": "ENDPOINT_NOT_FOUND",
  "status": 404
}

HTTP Status: 404
Time: 0.015s

⚠ Client error
```

## 🤝 Contributing

When adding new endpoints or modifying the API:

1. Update JSON examples
2. Test with both scripts
3. Update documentation
4. Commit all changes together

---

**Happy Testing! 🎉**
