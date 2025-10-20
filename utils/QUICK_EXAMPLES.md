# Quick Test Examples for SSM API

## PowerShell Commands (Windows)

## 1. Health Check

.\api_test.ps1 -Endpoint "/health-check" -Method "GET"

## 2. Generate AES Key (256 bits)

.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json"

## 3. Generate DES Key

.\api_test.ps1 -Endpoint "/generate-des-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_des_key_request.json"

## 4. Generate DES3 Key

.\api_test.ps1 -Endpoint "/generate-des3-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_des3_key_request.json"

## 5. Store Key

.\api_test.ps1 -Endpoint "/store-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\store_key_request.json"

## 6. Encrypt Data

.\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"

## 7. Decrypt Data

.\api_test.ps1 -Endpoint "/decrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\decrypt_request.json"

## 8. Get Keys by Label

.\api_test.ps1 -Endpoint "/get-data-keys" -Method "POST" -JsonFile "..\docs\json_examples\requests\get_data_keys_request.json"

## 9. Get Single Key

.\api_test.ps1 -Endpoint "/get-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\get_key_request.json"

## 10. Get All Keys
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST"


## Bash Commands (Linux/macOS)

## 1. Health Check
./api_test.sh /health-check GET

## 2. Generate AES Key (256 bits)
./api_test.sh /generate-aes-key POST ../docs/json_examples/requests/gen_aes_key_request.json

## 3. Generate DES Key
./api_test.sh /generate-des-key POST ../docs/json_examples/requests/gen_des_key_request.json

## 4. Generate DES3 Key
./api_test.sh /generate-des3-key POST ../docs/json_examples/requests/gen_des3_key_request.json

## 5. Store Key
./api_test.sh /store-key POST ../docs/json_examples/requests/store_key_request.json

## 6. Encrypt Data
./api_test.sh /encrypt POST ../docs/json_examples/requests/encrypt_request.json

## 7. Decrypt Data
./api_test.sh /decrypt POST ../docs/json_examples/requests/decrypt_request.json

## 8. Get Keys by Label
./api_test.sh /get-data-keys POST ../docs/json_examples/requests/get_data_keys_request.json

## 9. Get Single Key
./api_test.sh /get-key POST ../docs/json_examples/requests/get_key_request.json

## 10. Get All Keys
./api_test.sh /get-all-keys POST
