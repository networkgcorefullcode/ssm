# SoftHSM y PKCS#11 - Comandos de Administración

Guía completa para administrar SoftHSM y realizar operaciones PKCS#11 desde línea de comandos.

## 1. Instalación y Configuración Inicial

### Instalar SoftHSM

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install softhsm2 opensc-pkcs11
```

**CentOS/RHEL:**

```bash
sudo yum install softhsm opensc
# o con dnf
sudo dnf install softhsm opensc
```

**Windows:**

```powershell
# Descargar desde https://github.com/disig/SoftHSM2-for-Windows
# O usar chocolatey
choco install softhsm
```

### Configurar SoftHSM

**Ver configuración actual:**

```bash
softhsm2-util --show-slots
```

**Verificar ruta de configuración:**

```bash
# Linux
cat /etc/softhsm/softhsm2.conf
# O crear configuración personal
export SOFTHSM2_CONF="$HOME/.softhsm2.conf"
```

**Archivo de configuración ejemplo (`~/.softhsm2.conf`):**

```ini
# SoftHSM v2 configuration file

directories.tokendir = /home/user/.softhsm/tokens/
objectstore.backend = file
log.level = INFO
slots.removable = false
slots.mechanisms = ALL
```

## 2. Gestión de Tokens

### Inicializar Token

**Crear nuevo token:**

```bash
# Slot 0 con PIN de usuario y SO-PIN
softhsm2-util --init-token --slot 0 --label "MyToken" --pin 1234 --so-pin 5678

# Token específico para SSM
softhsm2-util --init-token --slot 0 --label "SSM-Token" --pin 1234 --so-pin 5678
```

**Inicializar múltiples tokens:**

```bash
# Token de desarrollo
softhsm2-util --init-token --slot 0 --label "DEV-Token" --pin 1234 --so-pin 5678

# Token de producción
softhsm2-util --init-token --slot 1 --label "PROD-Token" --pin 9876 --so-pin 5432
```

### Listar y Verificar Tokens

**Ver todos los slots:**

```bash
softhsm2-util --show-slots
```

**Ver slots con detalles:**

```bash
# Información detallada de slots
pkcs11-tool --list-slots --verbose

# Información de tokens
pkcs11-tool --list-token-slots
```

**Ver información específica de un token:**

```bash
# Token en slot específico
pkcs11-tool --slot 0 --list-token-slots
```

### Eliminar Token

**Eliminar token completo:**

```bash
softhsm2-util --delete-token --token "MyToken"

# O por slot
softhsm2-util --delete-token --slot 0
```

## 3. Gestión de Claves

### Generar Claves

**Generar clave AES:**

```bash
# AES-256
pkcs11-tool --slot 0 --login --pin 1234 \
  --keygen --key-type AES:32 \
  --label "MyAESKey" --id 01

# AES-128
pkcs11-tool --slot 0 --login --pin 1234 \
  --keygen --key-type AES:16 \
  --label "AES128Key" --id 02

# AES-192  
pkcs11-tool --slot 0 --login --pin 1234 \
  --keygen --key-type AES:24 \
  --label "AES192Key" --id 03
```

**Generar par de claves RSA:**

```bash
# RSA 2048 bits
pkcs11-tool --slot 0 --login --pin 1234 \
  --keypairgen --key-type RSA:2048 \
  --label "RSAKeyPair" --id 10

# RSA 4096 bits
pkcs11-tool --slot 0 --login --pin 1234 \
  --keypairgen --key-type RSA:4096 \
  --label "RSA4096Key" --id 11
```

**Generar par de claves ECDSA:**

```bash
# ECDSA P-256
pkcs11-tool --slot 0 --login --pin 1234 \
  --keypairgen --key-type EC:secp256r1 \
  --label "ECDSAKey" --id 20

# ECDSA P-384
pkcs11-tool --slot 0 --login --pin 1234 \
  --keypairgen --key-type EC:secp384r1 \
  --label "ECDSA384Key" --id 21
```

### Listar y Ver Claves

**Listar todas las claves:**

```bash
# Listar objetos públicos
pkcs11-tool --slot 0 --list-objects

# Listar objetos privados (requiere login)
pkcs11-tool --slot 0 --login --pin 1234 --list-objects
```

**Listar solo claves secretas:**

```bash
pkcs11-tool --slot 0 --login --pin 1234 \
  --list-objects --type secrkey
```

**Listar solo claves públicas:**

```bash
pkcs11-tool --slot 0 --list-objects --type pubkey
```

**Listar solo claves privadas:**

```bash
pkcs11-tool --slot 0 --login --pin 1234 \
  --list-objects --type privkey
```

**Ver detalles de una clave específica:**

```bash
# Por ID
pkcs11-tool --slot 0 --login --pin 1234 \
  --read-object --id 01 --type secrkey

# Por label
pkcs11-tool --slot 0 --login --pin 1234 \
  --read-object --label "MyAESKey" --type secrkey
```

### Eliminar Claves

**Eliminar clave por ID:**

```bash
pkcs11-tool --slot 0 --login --pin 1234 \
  --delete-object --id 01 --type secrkey
```

**Eliminar clave por label:**

```bash
pkcs11-tool --slot 0 --login --pin 1234 \
  --delete-object --label "MyAESKey" --type secrkey
```

**Eliminar todas las claves de un tipo:**

```bash
# ⚠️ CUIDADO: Elimina TODAS las claves secretas
pkcs11-tool --slot 0 --login --pin 1234 \
  --delete-object --type secrkey
```

## 4. Operaciones Criptográficas

### Cifrado/Descifrado

**Cifrar archivo con AES:**

```bash
# Crear archivo de prueba
echo "Mensaje secreto" > test.txt

# Cifrar
pkcs11-tool --slot 0 --login --pin 1234 \
  --encrypt --mechanism AES-CBC --input-file test.txt \
  --output-file test.enc --id 01

# Descifrar
pkcs11-tool --slot 0 --login --pin 1234 \
  --decrypt --mechanism AES-CBC --input-file test.enc \
  --output-file test.dec --id 01
```

**Cifrar con RSA:**

```bash
# Cifrar con clave pública RSA
pkcs11-tool --slot 0 --encrypt \
  --mechanism RSA-PKCS --input-file test.txt \
  --output-file test-rsa.enc --id 10

# Descifrar con clave privada RSA
pkcs11-tool --slot 0 --login --pin 1234 \
  --decrypt --mechanism RSA-PKCS --input-file test-rsa.enc \
  --output-file test-rsa.dec --id 10
```

### Firma Digital

**Firmar con RSA:**

```bash
# Crear hash del archivo
openssl dgst -sha256 -binary test.txt > test.hash

# Firmar
pkcs11-tool --slot 0 --login --pin 1234 \
  --sign --mechanism SHA256-RSA-PKCS --input-file test.hash \
  --output-file test.sig --id 10

# Verificar firma
pkcs11-tool --slot 0 --verify \
  --mechanism SHA256-RSA-PKCS --input-file test.hash \
  --signature-file test.sig --id 10
```

**Firmar con ECDSA:**

```bash
# Firmar con ECDSA
pkcs11-tool --slot 0 --login --pin 1234 \
  --sign --mechanism ECDSA-SHA256 --input-file test.hash \
  --output-file test-ec.sig --id 20

# Verificar firma ECDSA
pkcs11-tool --slot 0 --verify \
  --mechanism ECDSA-SHA256 --input-file test.hash \
  --signature-file test-ec.sig --id 20
```

## 5. Información del Sistema

### Capacidades del HSM

**Ver mecanismos soportados:**

```bash
# Todos los mecanismos
pkcs11-tool --slot 0 --list-mechanisms

# Mecanismos detallados
pkcs11-tool --slot 0 --list-mechanisms --verbose
```

**Información del módulo:**

```bash
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so --show-info
```

**Información detallada del token:**

```bash
pkcs11-tool --slot 0 --token-info
```

### Estadísticas y Estado

**Ver espacio usado/libre:**

```bash
pkcs11-tool --slot 0 --login --pin 1234 --list-objects | wc -l
```

**Ver sesiones activas:**

```bash
# Información de sesión
pkcs11-tool --slot 0 --login --pin 1234 --test
```

## 6. Backup y Restauración

### Exportar/Importar Objetos

**Exportar clave pública:**

```bash
pkcs11-tool --slot 0 --read-object \
  --type pubkey --id 10 --output-file rsa-pub.der
```

**Importar certificado:**

```bash
# Crear certificado auto-firmado primero
openssl req -new -x509 -key rsa-key.pem -out cert.pem -days 365

# Importar al HSM
pkcs11-tool --slot 0 --login --pin 1234 \
  --write-object cert.pem --type cert \
  --id 10 --label "MyCertificate"
```

### Backup del Token Completo

**Backup de directorio de tokens (SoftHSM):**

```bash
# Ubicación típica en Linux
sudo tar -czf softhsm-backup.tar.gz /var/lib/softhsm/tokens/

# Ubicación personal
tar -czf softhsm-backup.tar.gz ~/.softhsm/tokens/
```

**Restaurar backup:**

```bash
# Restaurar tokens
sudo tar -xzf softhsm-backup.tar.gz -C /

# Verificar restauración
softhsm2-util --show-slots
```

## 7. Troubleshooting y Diagnósticos

### Verificar Instalación

**Comprobar librerías:**

```bash
# Verificar SoftHSM
ls -la /usr/lib/*/softhsm/libsofthsm2.so
ldd /usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so

# Verificar OpenSC
which pkcs11-tool
pkcs11-tool --version
```

**Test básico de conectividad:**

```bash
# Test simple
pkcs11-tool --list-slots

# Test con módulo específico
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so --list-slots
```

### Logs y Debugging

**Habilitar logging detallado:**

```bash
# Modificar softhsm2.conf
echo "log.level = DEBUG" >> ~/.softhsm2.conf

# Ver logs en tiempo real (systemd)
journalctl -u softhsm -f

# Logs manuales
export SOFTHSM2_CONF="$HOME/.softhsm2.conf"
softhsm2-util --show-slots 2>&1 | tee debug.log
```

**Debugging de aplicaciones:**

```bash
# Variables de debug para PKCS#11
export PKCS11_DEBUG=1
export PKCS11SPY=/usr/lib/x86_64-linux-gnu/pkcs11-spy.so

# Ejecutar aplicación con tracing
strace -e trace=file ./ssm --cfg config.yml
```

### Problemas Comunes

**Error: "CKR_TOKEN_NOT_PRESENT"**

```bash
# Verificar que el token esté inicializado
softhsm2-util --show-slots

# Reinicializar si es necesario
softhsm2-util --init-token --slot 0 --label "NewToken" --pin 1234 --so-pin 5678
```

**Error: "CKR_USER_PIN_NOT_INITIALIZED"**

```bash
# Inicializar PIN de usuario
pkcs11-tool --slot 0 --login --so-pin 5678 --init-pin --pin 1234
```

**Error: "Module not found"**

```bash
# Verificar rutas de módulos
find /usr -name "*softhsm*" -type f 2>/dev/null
find /usr -name "*pkcs11*" -type f 2>/dev/null

# Instalar si falta
sudo apt install softhsm2 opensc-pkcs11
```

## 8. Scripts de Automatización

### Script de Setup Inicial

```bash
#!/bin/bash
# setup-hsm.sh

set -e

SLOT=0
TOKEN_LABEL="SSM-Production"
USER_PIN="1234"
SO_PIN="5678"

echo "🔧 Configurando SoftHSM para SSM..."

# 1. Verificar instalación
if ! command -v softhsm2-util &> /dev/null; then
    echo "❌ SoftHSM no encontrado. Instalando..."
    sudo apt update && sudo apt install -y softhsm2 opensc-pkcs11
fi

# 2. Crear directorio de tokens
mkdir -p ~/.softhsm/tokens

# 3. Configurar SoftHSM
cat > ~/.softhsm2.conf << EOF
directories.tokendir = $HOME/.softhsm/tokens/
objectstore.backend = file
log.level = INFO
slots.removable = false
EOF

export SOFTHSM2_CONF="$HOME/.softhsm2.conf"

# 4. Inicializar token
echo "📝 Inicializando token..."
softhsm2-util --init-token --slot $SLOT --label "$TOKEN_LABEL" --pin $USER_PIN --so-pin $SO_PIN

# 5. Verificar
echo "✅ Verificando configuración..."
softhsm2-util --show-slots

# 6. Generar clave de prueba
echo "🔑 Generando clave de prueba..."
pkcs11-tool --slot $SLOT --login --pin $USER_PIN \
  --keygen --key-type AES:32 --label "TestKey" --id 99

echo "🎉 Setup completado!"
echo "Token: $TOKEN_LABEL"
echo "Slot: $SLOT"
echo "PIN: $USER_PIN"
```

### Script de Limpieza

```bash
#!/bin/bash
# cleanup-hsm.sh

set -e

echo "🧹 Limpiando SoftHSM..."

# 1. Eliminar todas las claves
pkcs11-tool --slot 0 --login --pin 1234 --delete-object --type secrkey || true
pkcs11-tool --slot 0 --login --pin 1234 --delete-object --type privkey || true
pkcs11-tool --slot 0 --login --pin 1234 --delete-object --type pubkey || true

# 2. Eliminar token
softhsm2-util --delete-token --slot 0 || true

# 3. Limpiar directorio
rm -rf ~/.softhsm/tokens/*

echo "✅ Limpieza completada"
```

## 9. Integración con SSM

### Variables de Entorno para SSM

```bash
# Configuración para el proyecto SSM
export SOFTHSM2_CONF="$PWD/factory/.softhsm2.conf"
export PKCS11_MODULE="/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so"
export HSM_SLOT=0
export HSM_PIN=1234
```

### Configuración Recomendada para Desarrollo

```yaml
# factory/ssmConfig.yml
configuration:
  ssmName: "SSM-Dev"
  ssmId: "ssm-dev-001"  
  socketPath: "/tmp/ssm.sock"
  pkcsPath: "/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so"
  lotsNumber: 0
  pin: "1234"
```

### Test de Integración

```bash
#!/bin/bash
# test-ssm-integration.sh

# 1. Setup HSM
./setup-hsm.sh

# 2. Iniciar SSM
./ssm --cfg factory/ssmConfig.yml &
SSM_PID=$!

# 3. Wait for startup
sleep 2

# 4. Test API
curl -X POST --unix-socket /tmp/ssm.sock http://localhost/generate-aes-key \
  -H "Content-Type: application/json" \
  -d '{"label":"integration-test","id":"int001","bits":256}'

# 5. Cleanup
kill $SSM_PID
```

## Notas Importantes

- ⚠️ **Seguridad**: Nunca uses PINs débiles en producción
- 📁 **Backup**: Siempre haz backup de los tokens antes de cambios importantes
- 🔄 **Testing**: Usa tokens separados para desarrollo y producción
- 📝 **Logging**: Habilita logs detallados durante troubleshooting
- 🔧 **Permisos**: Asegúrate de que el usuario tenga permisos en los directorios de SoftHSM
