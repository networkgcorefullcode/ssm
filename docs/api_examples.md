# SSM API Examples

Este documento contiene ejemplos de curl para testear la API del SSM (Secure Storage Manager).

## Configuración Previa

Asegúrate de que el servidor SSM esté ejecutándose:

```bash
./ssm --cfg factory/ssmConfig.yml
```

El servidor escucha en un socket Unix según la configuración en `ssmConfig.yml`.

## 1. Generar Clave AES

### Generar clave AES de 256 bits

```bash
curl -X POST http://dummy/generate-aes-key \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "label": "MySecretKey",
    "id": "key001",
    "bits": 256
  }'
```

### Generar clave AES de 128 bits

```bash
curl -X POST http://dummy/generate-aes-key \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "label": "TestKey128",
    "id": "test128",
    "bits": 128
  }'
```

### Generar clave AES de 192 bits

```bash
curl -X POST http://dummy/generate-aes-key \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "label": "ProductionKey",
    "id": "prod001",
    "bits": 192
  }'
```

**Respuesta esperada (éxito):**

```json
{
  "handle": 123456789,
  "label": "MySecretKey",
  "id": "key001",
  "bits": 256
}
```

**Respuesta de error (validación):**

```json
{
  "title": "Validation Error",
  "detail": "Label is required",
  "status": 400
}
```

## 2. Cifrar Datos

### Cifrar texto plano

```bash
curl -X POST http://dummy/encrypt \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "MySecretKey",
    "plain_b64": "SGVsbG8gV29ybGQh"
  }'
```

### Cifrar datos sensibles
```bash
curl -X POST http://dummy/encrypt \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "keyLabel": "ProductionKey",
    "plaintext": "eyJ1c2VyIjoiYWRtaW4iLCJwYXNzd29yZCI6InNlY3JldDEyMyJ9",
    "iv": "cmFuZG9taXZlY3RvcjEyMw=="
  }'
```

## 3. Descifrar Datos

### Descifrar texto
```bash
curl -X POST http://dummy/decrypt \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "keyLabel": "MySecretKey",
    "ciphertext": "YWJjZGVmZ2hpams...",
    "iv": "MTIzNDU2Nzg5MEFCQ0RFRg=="
  }'
```

## 4. Casos de Error

### Error: Método HTTP incorrecto
```bash
curl -X GET http://dummy/generate-aes-key \
  --unix-socket /var/run/socket.so
```

**Respuesta:**
```json
{
  "title": "Method Not Allowed",
  "detail": "Only POST method is allowed",
  "status": 405
}
```

### Error: JSON inválido
```bash
curl -X POST http://dummy/generate-aes-key \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{"label": "test", "invalid_json"}'
```

### Error: Bits inválidos
```bash
curl -X POST http://dummy/generate-aes-key \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "label": "TestKey",
    "id": "test001",
    "bits": 512
  }'
```

**Respuesta:**
```json
{
  "title": "Validation Error",
  "detail": "Bits must be 128, 192, or 256",
  "status": 400
}
```

### Error: Campo requerido faltante
```bash
curl -X POST http://dummy/generate-aes-key \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "id": "test001",
    "bits": 256
  }'
```

## 5. Testing con Scripts

### Script de prueba completo
```bash
#!/bin/bash

SOCKET_PATH="/var/run/socket.so"
BASE_URL="http://dummy"

echo "=== Testing SSM API ==="

# Test 1: Generate AES Key
echo "1. Generating AES Key..."
KEY_RESPONSE=$(curl -s -X POST \
  --unix-socket $SOCKET_PATH \
  $BASE_URL/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{
    "label": "TestAPIKey",
    "id": "api_test_001",
    "bits": 256
  }')

echo "Response: $KEY_RESPONSE"

# Test 2: Encrypt Data
echo "2. Encrypting data..."
ENCRYPT_RESPONSE=$(curl -s -X POST \
  --unix-socket $SOCKET_PATH \
  $BASE_URL/encrypt \
  -H "Content-Type: application/json" \
  -d '{
    "keyLabel": "TestAPIKey",
    "plaintext": "VGVzdCBtZXNzYWdl",
    "iv": "MTIzNDU2Nzg5MEFCQ0RFRg=="
  }')

echo "Response: $ENCRYPT_RESPONSE"

# Test 3: Decrypt Data
echo "3. Decrypting data..."
# Usar el ciphertext de la respuesta anterior
curl -s -X POST \
  --unix-socket $SOCKET_PATH \
  $BASE_URL/decrypt \
  -H "Content-Type: application/json" \
  -d '{
    "keyLabel": "TestAPIKey", 
    "ciphertext": "...",
    "iv": "MTIzNDU2Nzg5MEFCQ0RFRg=="
  }'
```

## 6. Troubleshooting

### Verificar que el socket existe
```bash
ls -la /var/run/socket.so
```

### Verificar que el servidor está ejecutándose
```bash
ps aux | grep ssm
```

### Probar conectividad básica
```bash
curl --unix-socket /var/run/socket.so http://dummy/
```

## 7. Datos de Prueba en Base64

Para las pruebas, puedes usar estos datos pre-codificados:

### Texto plano
- `"Hello World!"` → `SGVsbG8gV29ybGQh`
- `"Secret Message"` → `U2VjcmV0IE1lc3NhZ2U=`
- `"Test data 123"` → `VGVzdCBkYXRhIDEyMw==`

### IVs de ejemplo (16 bytes en Base64)
- `1234567890ABCDEF` → `MTIzNDU2Nzg5MEFCQ0RFRg==`
- `randomvector123` → `cmFuZG9tdmVjdG9yMTIz`
- `0000000000000000` → `MDAwMDAwMDAwMDAwMDAwMA==`

## Notas Importantes

1. **Socket Path**: Ajusta el path del socket según tu configuración en `ssmConfig.yml`
2. **Datos en Base64**: Todos los datos binarios (plaintext, ciphertext, IV) deben estar en Base64
3. **IV Consistency**: Usa el mismo IV para cifrar y descifrar
4. **Key Management**: Las claves se almacenan permanentemente en el HSM hasta que se eliminen
5. **Error Handling**: Siempre verifica los códigos de estado HTTP y los mensajes de error