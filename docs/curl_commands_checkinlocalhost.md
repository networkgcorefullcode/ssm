# SSM API - cURL Commands for Testing

This document contains cURL commands to test all SSM API endpoints on localhost.

## Base URL

- **HTTPS**: `https://localhost:9000`

---

## 1. Health Check

### GET Request

```bash
curl -k -X GET https://localhost:9000/health-check \
  -H "Accept: application/json"
```

---

## 2. Generate Keys

### Generate AES Key (128 bits)

```bash
curl -k -X POST https://localhost:9000/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "bits": 128
  }'
```

### Generate AES Key (256 bits)

```bash
curl -k -X POST https://localhost:9000/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{
    "id": 2,
    "bits": 256
  }'
```

### Generate DES Key

```bash
curl -k -X POST https://localhost:9000/generate-des-key \
  -H "Content-Type: application/json" \
  -d '{
    "id": 3
  }'
```

### Generate DES3 Key

```bash
curl -k -X POST https://localhost:9000/generate-des3-key \
  -H "Content-Type: application/json" \
  -d '{
    "id": 4
  }'
```

---

## 3. Store Key

### Store AES Key

```bash
curl -k -X POST https://localhost:9000/store-key \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K4_AES",
    "id": 10,
    "key_value": "0123456789abcdef0123456789abcdef",
    "key_type": "AES"
  }'
```

### Store DES3 Key

```bash
curl -k -X POST https://localhost:9000/store-key \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K4_DES",
    "id": 11,
    "key_value": "1234567890abcdef1234567890abcdef1234567890abcdef",
    "key_type": "DES3"
  }'
```

### Store DES Key

```bash
curl -k -X POST https://localhost:9000/store-key \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K4_DES",
    "id": 12,
    "key_value": "0123456789abcdef",
    "key_type": "DES"
  }'
```

---

## 4. Encryption

### Encrypt with AES

```bash
curl -k -X POST https://localhost:9000/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "KEY_ENCRIPTION_AES",
    "plain": "48656c6c6f20576f726c6421",
    "encryption_algorithm": 1
  }'
```

### Encrypt with DES

```bash
curl -k -X POST https://localhost:9000/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "KEY_ENCRIPTION_DES",
    "plain": "48656c6c6f20576f726c6421",
    "encryption_algorithm": 3
  }'
```

### Encrypt with DES3

```bash
curl -k -X POST https://localhost:9000/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "KEY_ENCRIPTION_DES3",
    "plain": "48656c6c6f20576f726c6421",
    "encryption_algorithm": 4
  }'
```

---

## 5. Decryption

### Decrypt with AES

```bash
curl -k -X POST https://localhost:9000/decrypt \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "fill",
    "cipher": "fill",
    "iv": "fill",
    "id": 1,
    "encryption_algorithm": 1
  }'
```

### Decrypt with DES3

```bash
curl -k -X POST https://localhost:9000/decrypt \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K4_DES",
    "cipher": "f8e7d6c5b4a39281",
    "iv": "1234567890abcdef",
    "id": 4,
    "encryption_algorithm": 4
  }'
```

---

## 6. Key Retrieval

### Get Keys by Label

```bash
curl -k -X POST https://localhost:9000/get-data-keys \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "my-aes-key"
  }'
```

### Get Single Key by Label

```bash
curl -k -X POST https://localhost:9000/get-key \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "my-aes-key",
    "id": 1
  }'
```

### Get All Keys

```bash
curl -k -X POST https://localhost:9000/get-all-keys \
  -H "Content-Type: application/json"
```

---

## 7. Complete Workflow Example

### Step 1: Health Check

```bash
curl -k -X GET https://localhost:9000/health-check
```

### Step 2: Generate AES-256 Key

```bash
curl -k -X POST https://localhost:9000/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{
    "id": 100,
    "bits": 256
  }'
```

### Step 3: Encrypt Data

```bash
# Plain text "Hello World!" in hexadecimal: 48656c6c6f20576f726c6421
curl -k -X POST https://localhost:9000/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "my-aes-key",
    "plain": "48656c6c6f20576f726c6421",
    "encryption_algorithm": 1
  }'
```

### Step 4: Decrypt Data (use cipher and IV from Step 3)

```bash
curl -k -X POST https://localhost:9000/decrypt \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "my-aes-key",
    "cipher": "<CIPHER_FROM_STEP_3>",
    "iv": "<IV_FROM_STEP_3>",
    "id": 100,
    "encryption_algorithm": 1
  }'
```

### Step 5: Get All Keys

```bash
curl -k -X POST https://localhost:9000/get-all-keys \
  -H "Content-Type: application/json"
```

---

## 8. Pretty Print JSON (with jq)

If you have `jq` installed, you can format the JSON output:

```bash
curl -k -X GET https://localhost:9000/health-check | jq '.'
```

```bash
curl -k -X POST https://localhost:9000/get-all-keys \
  -H "Content-Type: application/json" | jq '.'
```

---

## 9. HTTPS Examples (with self-signed certificates)

### Health Check (HTTPS)

```bash
curl -k -k -X GET https://localhost:8443/health-check \
  -H "Accept: application/json"
```

### Generate Key (HTTPS)

```bash
curl -k -k -X POST https://localhost:8443/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "bits": 256
  }'
```

### Store Key (HTTPS)

```bash
curl -k -k -X POST https://localhost:8443/store-key \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K4_DES",
    "id": 1,
    "key_value": "1234567890abcdef1234567890abcdef1234567890abcdef",
    "key_type": "DES3"
  }'
```

---

## 10. Verbose Mode (Debug)

To see detailed request/response information:

```bash
curl -k -v -X POST https://localhost:9000/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "bits": 256
  }'
```

---

## Notes

### Encryption Algorithms

- `1` = AES
- `2` = AES (alternative)
- `3` = DES
- `4` = DES3

### Data Formats

- **plain**: Hexadecimal string
- **cipher**: Base64 or Hexadecimal string
- **iv**: Base64 or Hexadecimal string
- **key_value**: Hexadecimal string

### Key Sizes

- **AES**: 128, 192, or 256 bits
- **DES**: 56 bits (8 bytes)
- **DES3**: 168 bits (24 bytes)

### Common HTTP Status Codes

- `200` - Success
- `400` - Bad Request (invalid JSON or parameters)
- `404` - Not Found (endpoint doesn't exist)
- `500` - Internal Server Error

---

## Tips

1. **Save responses to file:**

   ```bash
   curl -k -X GET https://localhost:9000/health-check -o response.json
   ```

2. **Include response headers:**

   ```bash
   curl -k -i -X GET https://localhost:9000/health-check
   ```

3. **Set timeout:**

   ```bash
   curl -k --max-time 10 -X GET https://localhost:9000/health-check
   ```

4. **Use environment variables:**

   ```bash
   export SSM_URL="https://localhost:9000"
   curl -k -X GET $SSM_URL/health-check
   ```

5. **Test with different keys:**

   ```bash
   # Create a loop to test multiple keys
   for i in {1..5}; do
     curl -k -X POST https://localhost:9000/generate-aes-key \
       -H "Content-Type: application/json" \
       -d "{\"id\": $i, \"bits\": 256}"
     echo ""
   done
   ```
