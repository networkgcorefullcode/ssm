# Quick Test Commands for SSM API

Comandos rápidos para probar la API SSM. Asegúrate de que el servidor esté ejecutándose.

## Variables de Entorno

```bash
export SOCKET_PATH="/tmp/ssm.sock"
export BASE_URL="http://localhost"
```

## 1. Test Básico - Generar Clave

```bash
curl -X POST \
  --unix-socket $SOCKET_PATH \
  $BASE_URL/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{"label":"test-key","id":"001","bits":256}'
```

## 2. Test Básico - Cifrar

```bash
curl -X POST \
  --unix-socket $SOCKET_PATH \
  $BASE_URL/encrypt \
  -H "Content-Type: application/json" \
  -d '{"keyLabel":"test-key","plaintext":"SGVsbG8gV29ybGQ=","iv":"MTIzNDU2Nzg5MEFCQ0RFRg=="}'
```

## 3. Test Básico - Descifrar

```bash
curl -X POST \
  --unix-socket $SOCKET_PATH \
  $BASE_URL/decrypt \
  -H "Content-Type: application/json" \
  -d '{"keyLabel":"test-key","ciphertext":"RESULTADO_DEL_CIFRADO","iv":"MTIzNDU2Nzg5MEFCQ0RFRg=="}'
```

## 4. Script de Test Automático

```bash
#!/bin/bash

# Configuración
SOCKET="/tmp/ssm.sock"
URL="http://localhost"

echo "🔧 Testing SSM API..."

# 1. Generar clave
echo "📝 1. Generating AES key..."
RESPONSE=$(curl -s -X POST --unix-socket $SOCKET $URL/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{"label":"auto-test","id":"auto001","bits":256}')
echo "Response: $RESPONSE"

# 2. Cifrar datos
echo "🔒 2. Encrypting data..."
ENCRYPT_RESP=$(curl -s -X POST --unix-socket $SOCKET $URL/encrypt \
  -H "Content-Type: application/json" \
  -d '{"keyLabel":"auto-test","plaintext":"VGVzdE1lc3NhZ2U=","iv":"dGVzdGl2ZWN0b3IxMjM0NQ=="}')
echo "Response: $ENCRYPT_RESP"

# 3. Extraer ciphertext y descifrar
CIPHERTEXT=$(echo $ENCRYPT_RESP | jq -r '.ciphertext')
if [ "$CIPHERTEXT" != "null" ] && [ "$CIPHERTEXT" != "" ]; then
    echo "🔓 3. Decrypting data..."
    DECRYPT_RESP=$(curl -s -X POST --unix-socket $SOCKET $URL/decrypt \
      -H "Content-Type: application/json" \
      -d "{\"keyLabel\":\"auto-test\",\"ciphertext\":\"$CIPHERTEXT\",\"iv\":\"dGVzdGl2ZWN0b3IxMjM0NQ==\"}")
    echo "Response: $DECRYPT_RESP"
else
    echo "❌ Error: No se pudo obtener ciphertext"
fi

echo "✅ Test completed"
```

## 5. Validación de Errores

### Error 405 - Método no permitido
```bash
curl -X GET --unix-socket /tmp/ssm.sock http://localhost/generate-aes-key
```

### Error 400 - Campo requerido faltante
```bash
curl -X POST --unix-socket /tmp/ssm.sock http://localhost/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{"id":"test","bits":256}'
```

### Error 400 - Bits inválidos
```bash
curl -X POST --unix-socket /tmp/ssm.sock http://localhost/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{"label":"test","id":"test","bits":512}'
```

## 6. Datos de Prueba

### Base64 Encoded Data
- `"Hello"` = `SGVsbG8=`
- `"World"` = `V29ybGQ=`
- `"TestMessage"` = `VGVzdE1lc3NhZ2U=`

### IVs de Prueba (16 bytes)
- `"testvector12345"` = `dGVzdHZlY3RvcjEyMzQ1`
- `"1234567890ABCDEF"` = `MTIzNDU2Nzg5MEFCQ0RFRg==`

## 7. Verificación de Estado

### Verificar que el socket existe
```bash
ls -la /tmp/ssm.sock
```

### Verificar proceso SSM
```bash
ps aux | grep ssm
```

### Test de conectividad
```bash
curl --unix-socket /tmp/ssm.sock http://localhost/ || echo "Socket no disponible"
```