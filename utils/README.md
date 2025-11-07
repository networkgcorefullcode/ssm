# SSM API Testing Utilities

This directory contains utility scripts for testing the SSM API endpoints.

## Scripts

### 1. `api_test.sh` (Bash/Linux/macOS)
Shell script for testing API endpoints using curl.

### 2. `api_test.ps1` (PowerShell/Windows)
PowerShell script for testing API endpoints using Invoke-RestMethod.

## Prerequisites

### For Bash script (`api_test.sh`)
- **curl**: Required for making HTTP requests
- **jq**: Optional but recommended for pretty JSON output

### For PowerShell script (`api_test.ps1`)
- PowerShell 5.1 or higher (pre-installed on Windows 10/11)
- No additional dependencies required

## Usage

### Bash Script

```bash
# Make the script executable (first time only)
chmod +x api_test.sh

# Basic usage
./api_test.sh <endpoint> <method> [json_file] [base_url]

# Examples
./api_test.sh /health-check GET
./api_test.sh /encrypt POST ../docs/json_examples/requests/encrypt_request.json
./api_test.sh /decrypt POST ../docs/json_examples/requests/decrypt_request.json http://localhost:9000

# Show help
./api_test.sh --help
```

### PowerShell Script

```powershell
# Basic usage
.\api_test.ps1 -Endpoint <endpoint> -Method <method> [-JsonFile <file>] [-BaseUrl <url>]

# Examples
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"
.\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"
.\api_test.ps1 -Endpoint "/decrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\decrypt_request.json" -BaseUrl "http://localhost:9000"

# Show help
.\api_test.ps1 -Help
```

## Available Endpoints

| Endpoint | Method | Description | Example JSON |
|----------|--------|-------------|--------------|
| `/health-check` | GET/POST | Check service health | N/A |
| `/encrypt` | POST | Encrypt data | `encrypt_request.json` |
| `/decrypt` | POST | Decrypt data | `decrypt_request.json` |
| `/generate-aes-key` | POST | Generate AES key | `gen_aes_key_request.json` |
| `/generate-des-key` | POST | Generate DES key | `gen_des_key_request.json` |
| `/generate-des3-key` | POST | Generate DES3 key | `gen_des3_key_request.json` |
| `/store-key` | POST | Store a key | `store_key_request.json` |
| `/get-data-keys` | POST | Get keys by label | `get_data_keys_request.json` |
| `/get-key` | POST | Get single key | `get_key_request.json` |
| `/get-all-keys` | POST | Get all keys | N/A (no body needed) |

## JSON Examples

Example JSON files are located in `../docs/json_examples/`:
- `requests/` - Request body examples
- `responses/` - Expected response examples

### Request Examples

```bash
# Encrypt data
.\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"

# Generate AES key
.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json"

# Get all keys
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST"
```

## Output Format

Both scripts provide:
- Request details (URL, method, body)
- Response body (pretty-printed JSON)
- HTTP status code
- Request duration
- Color-coded success/error indicators

### Success Response (200-299)
```
✓ Request successful
```

### Client Error (400-499)
```
⚠ Client error
```

### Server Error (500+)
```
✗ Server error
```

## Customization

### Change Default Base URL

**Bash:**
```bash
export SSM_BASE_URL="http://localhost:9000"
./api_test.sh /health-check GET
```

**PowerShell:**
```powershell
$env:SSM_BASE_URL = "http://localhost:9000"
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"
```

Or pass it directly:
```bash
./api_test.sh /health-check GET "" http://localhost:9000
```

```powershell
.\api_test.ps1 -Endpoint "/health-check" -Method "GET" -BaseUrl "http://localhost:9000"
```

## Testing Workflow

### 1. Start the SSM server
```bash
go run ssm.go
```

### 2. Test health check
```powershell
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"
```

### 3. Generate a key
```powershell
.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json"
```

### 4. Encrypt data
```powershell
.\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"
```

### 5. Decrypt data
```powershell
.\api_test.ps1 -Endpoint "/decrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\decrypt_request.json"
```

### 6. Get all keys
```powershell
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST"
```

## Troubleshooting

### Bash Script Issues

**curl not found:**
```bash
# Ubuntu/Debian
sudo apt-get install curl

# macOS
brew install curl
```

**jq not found (optional):**
```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq
```

### PowerShell Script Issues

**Execution policy error:**
```powershell
# Run as Administrator
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

**SSL/TLS errors (for HTTPS):**
```powershell
# Bypass SSL validation (development only!)
[System.Net.ServicePointManager]::ServerCertificateValidationCallback = {$true}
```

## Contributing

When adding new endpoints:
1. Create corresponding JSON examples in `../docs/json_examples/requests/`
2. Create expected response examples in `../docs/json_examples/responses/`
3. Update this README with the new endpoint information

## License

Same as the main SSM project.
