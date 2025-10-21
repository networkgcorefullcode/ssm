# SSM API Test - Practical Example
# This script demonstrates a complete workflow using the SSM API

Write-Host "`n=== SSM API Complete Workflow Test ===" -ForegroundColor Cyan
Write-Host "This will test the complete encryption/decryption workflow`n" -ForegroundColor Yellow

$BaseUrl = "http://localhost:8080"

# Step 1: Health Check
Write-Host "`n[Step 1] Checking service health..." -ForegroundColor Green
.\api_test.ps1 -Endpoint "/health-check" -Method "GET" -BaseUrl $BaseUrl

Read-Host "`nPress Enter to continue to Step 2"

# Step 2: Generate AES Key
Write-Host "`n[Step 2] Generating AES-256 key..." -ForegroundColor Green
.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json" -BaseUrl $BaseUrl

Read-Host "`nPress Enter to continue to Step 3"

# Step 3: Encrypt Data
Write-Host "`n[Step 3] Encrypting data..." -ForegroundColor Green
.\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json" -BaseUrl $BaseUrl

Read-Host "`nPress Enter to continue to Step 4"

# Step 4: Decrypt Data
Write-Host "`n[Step 4] Decrypting data..." -ForegroundColor Green
Write-Host "Note: Update the decrypt_request.json with the cipher and IV from Step 3" -ForegroundColor Yellow
.\api_test.ps1 -Endpoint "/decrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\decrypt_request.json" -BaseUrl $BaseUrl

Read-Host "`nPress Enter to continue to Step 5"

# Step 5: Get All Keys
Write-Host "`n[Step 5] Retrieving all keys..." -ForegroundColor Green
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST" -BaseUrl $BaseUrl

Write-Host "`n=== Workflow Complete ===" -ForegroundColor Cyan
Write-Host "All steps executed successfully!`n" -ForegroundColor Green
