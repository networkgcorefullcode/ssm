# SoftHSM - Comandos Administrativos Avanzados

Comandos avanzados para administración, monitoreo y troubleshooting de SoftHSM.

## Información del Sistema

### Estado General del HSM

```bash
# Ver información completa del módulo
pkcs11-tool --module /usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so --show-info

# Ver versión de SoftHSM
softhsm2-util --version

# Ver configuración activa
softhsm2-util --show-slots --verbose
```

### Capacidades y Límites

```bash
# Ver todos los mecanismos disponibles
pkcs11-tool --slot 0 --list-mechanisms | grep -E "(AES|RSA|ECDSA|SHA)"

# Ver límites de memoria del token
pkcs11-tool --slot 0 --token-info | grep -E "(memory|space)"

# Contar objetos por tipo
echo "=== Resumen de Objetos ==="
echo "Claves secretas: $(pkcs11-tool --slot 0 --login --pin 1234 --list-objects --type secrkey | grep "Secret Key" | wc -l)"
echo "Claves públicas: $(pkcs11-tool --slot 0 --list-objects --type pubkey | grep "Public Key" | wc -l)"
echo "Claves privadas: $(pkcs11-tool --slot 0 --login --pin 1234 --list-objects --type privkey | grep "Private Key" | wc -l)"
echo "Certificados: $(pkcs11-tool --slot 0 --list-objects --type cert | grep "Certificate" | wc -l)"
```

## Operaciones Batch

### Generar Múltiples Claves

```bash
#!/bin/bash
# Generar conjunto de claves para testing

SLOT=0
PIN=1234

echo "Generando claves de prueba..."

# Claves AES para diferentes propósitos
for i in {1..5}; do
    pkcs11-tool --slot $SLOT --login --pin $PIN \
      --keygen --key-type AES:32 \
      --label "TestAES$i" --id $(printf "%02d" $i)
    echo "Generada clave AES TestAES$i"
done

# Pares de claves RSA
for size in 2048 3072 4096; do
    id=$((10 + size/1024))
    pkcs11-tool --slot $SLOT --login --pin $PIN \
      --keypairgen --key-type RSA:$size \
      --label "RSA${size}Key" --id $id
    echo "Generado par RSA de $size bits"
done
```

### Backup Selectivo

```bash
#!/bin/bash
# Backup de claves específicas

BACKUP_DIR="hsm_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p $BACKUP_DIR

echo "Creando backup en $BACKUP_DIR..."

# Exportar información de objetos
pkcs11-tool --slot 0 --login --pin 1234 --list-objects > $BACKUP_DIR/objects_list.txt

# Exportar claves públicas
pkcs11-tool --slot 0 --list-objects --type pubkey | grep "ID:" | while read line; do
    id=$(echo $line | grep -o 'ID: [0-9a-f]*' | cut -d' ' -f2)
    pkcs11-tool --slot 0 --read-object --type pubkey --id $id \
      --output-file $BACKUP_DIR/pubkey_$id.der 2>/dev/null || true
done

# Exportar certificados
pkcs11-tool --slot 0 --list-objects --type cert | grep "ID:" | while read line; do
    id=$(echo $line | grep -o 'ID: [0-9a-f]*' | cut -d' ' -f2)
    pkcs11-tool --slot 0 --read-object --type cert --id $id \
      --output-file $BACKUP_DIR/cert_$id.der 2>/dev/null || true
done

echo "Backup completado en $BACKUP_DIR"
```

## Monitoreo y Performance

### Scripts de Monitoreo

```bash
#!/bin/bash
# hsm-monitor.sh - Monitoreo continuo del HSM

SLOT=0
INTERVAL=5

while true; do
    clear
    echo "=== SoftHSM Monitor $(date) ==="
    echo
    
    # Estado del token
    echo "📊 Estado del Token:"
    pkcs11-tool --slot $SLOT --token-info | grep -E "(Label|Available|Total)"
    echo
    
    # Conteo de objetos
    echo "🔑 Objetos almacenados:"
    secret_keys=$(pkcs11-tool --slot $SLOT --login --pin 1234 --list-objects --type secrkey 2>/dev/null | grep "Secret Key" | wc -l)
    pub_keys=$(pkcs11-tool --slot $SLOT --list-objects --type pubkey 2>/dev/null | grep "Public Key" | wc -l)
    priv_keys=$(pkcs11-tool --slot $SLOT --login --pin 1234 --list-objects --type privkey 2>/dev/null | grep "Private Key" | wc -l)
    
    echo "  Claves secretas: $secret_keys"
    echo "  Claves públicas: $pub_keys"
    echo "  Claves privadas: $priv_keys"
    echo
    
    # Uso de memoria (SoftHSM específico)
    echo "💾 Uso de espacio:"
    du -sh ~/.softhsm/tokens/* 2>/dev/null | head -5
    echo
    
    sleep $INTERVAL
done
```

### Test de Performance

```bash
#!/bin/bash
# performance-test.sh

SLOT=0
PIN=1234
ITERATIONS=100

echo "🚀 Test de Performance SoftHSM"
echo "Iteraciones: $ITERATIONS"
echo

# Test generación de claves AES
echo "📝 Test: Generación de claves AES..."
start_time=$(date +%s.%N)
for i in $(seq 1 $ITERATIONS); do
    pkcs11-tool --slot $SLOT --login --pin $PIN \
      --keygen --key-type AES:32 \
      --label "PerfTest$i" --id $(printf "%04d" $i) >/dev/null 2>&1
done
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
rate=$(echo "scale=2; $ITERATIONS / $duration" | bc)
echo "  Tiempo total: ${duration}s"
echo "  Rate: ${rate} claves/segundo"

# Test cifrado
echo
echo "🔒 Test: Operaciones de cifrado..."
echo "Mensaje de prueba" > test_message.txt
start_time=$(date +%s.%N)
for i in $(seq 1 50); do
    pkcs11-tool --slot $SLOT --login --pin $PIN \
      --encrypt --mechanism AES-CBC \
      --input-file test_message.txt --output-file /tmp/test_enc_$i \
      --id 0001 >/dev/null 2>&1
done
end_time=$(date +%s.%N)
duration=$(echo "$end_time - $start_time" | bc)
rate=$(echo "scale=2; 50 / $duration" | bc)
echo "  Tiempo total: ${duration}s"
echo "  Rate: ${rate} operaciones/segundo"

# Cleanup
rm -f test_message.txt /tmp/test_enc_*
for i in $(seq 1 $ITERATIONS); do
    pkcs11-tool --slot $SLOT --login --pin $PIN \
      --delete-object --type secrkey --id $(printf "%04d" $i) >/dev/null 2>&1
done

echo
echo "✅ Test completado"
```

## Troubleshooting Avanzado

### Diagnóstico de Problemas

```bash
#!/bin/bash
# hsm-diagnostic.sh

echo "🔍 Diagnóstico SoftHSM"
echo "===================="

# 1. Verificar instalación
echo "📦 1. Verificación de instalación:"
if command -v softhsm2-util &> /dev/null; then
    echo "  ✅ softhsm2-util encontrado: $(which softhsm2-util)"
    echo "  📄 Versión: $(softhsm2-util --version 2>&1 | head -1)"
else
    echo "  ❌ softhsm2-util no encontrado"
fi

if command -v pkcs11-tool &> /dev/null; then
    echo "  ✅ pkcs11-tool encontrado: $(which pkcs11-tool)"
else
    echo "  ❌ pkcs11-tool no encontrado"
fi

# 2. Verificar librerías
echo
echo "📚 2. Verificación de librerías:"
SOFTHSM_LIB="/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so"
if [ -f "$SOFTHSM_LIB" ]; then
    echo "  ✅ Librería SoftHSM: $SOFTHSM_LIB"
    echo "  🔗 Dependencias:"
    ldd "$SOFTHSM_LIB" | grep -E "(not found|missing)" && echo "  ❌ Dependencias faltantes" || echo "  ✅ Dependencias OK"
else
    echo "  ❌ Librería SoftHSM no encontrada"
    echo "  🔍 Buscando alternativas:"
    find /usr -name "*softhsm*" -type f 2>/dev/null | head -5
fi

# 3. Verificar configuración
echo
echo "⚙️ 3. Verificación de configuración:"
if [ -n "$SOFTHSM2_CONF" ]; then
    echo "  📄 Config file: $SOFTHSM2_CONF"
    if [ -f "$SOFTHSM2_CONF" ]; then
        echo "  ✅ Archivo existe"
        echo "  📋 Contenido:"
        cat "$SOFTHSM2_CONF" | sed 's/^/      /'
    else
        echo "  ❌ Archivo no existe"
    fi
else
    echo "  ⚠️  Variable SOFTHSM2_CONF no definida"
    echo "  🔍 Buscando archivos de configuración:"
    find /etc /home -name "softhsm*.conf" 2>/dev/null | head -3
fi

# 4. Test de conectividad básica
echo
echo "🔌 4. Test de conectividad:"
if pkcs11-tool --list-slots >/dev/null 2>&1; then
    echo "  ✅ Conexión PKCS#11 OK"
    slot_count=$(pkcs11-tool --list-slots 2>/dev/null | grep "Slot" | wc -l)
    echo "  📊 Slots disponibles: $slot_count"
else
    echo "  ❌ Error de conexión PKCS#11"
    echo "  🔍 Error details:"
    pkcs11-tool --list-slots 2>&1 | sed 's/^/      /'
fi

# 5. Verificar permisos
echo
echo "🔐 5. Verificación de permisos:"
TOKEN_DIR="$HOME/.softhsm/tokens"
if [ -d "$TOKEN_DIR" ]; then
    echo "  ✅ Directorio tokens: $TOKEN_DIR"
    echo "  📁 Permisos: $(ls -ld "$TOKEN_DIR" | awk '{print $1, $3, $4}')"
    token_count=$(ls -1 "$TOKEN_DIR" 2>/dev/null | wc -l)
    echo "  📊 Tokens encontrados: $token_count"
else
    echo "  ⚠️  Directorio tokens no existe: $TOKEN_DIR"
fi

# 6. Test de operación básica
echo
echo "🧪 6. Test de operación básica:"
if softhsm2-util --show-slots >/dev/null 2>&1; then
    echo "  ✅ softhsm2-util functional"
    echo "  📋 Slots info:"
    softhsm2-util --show-slots 2>/dev/null | sed 's/^/      /'
else
    echo "  ❌ Error en softhsm2-util"
    echo "  🔍 Error details:"
    softhsm2-util --show-slots 2>&1 | sed 's/^/      /'
fi

echo
echo "🏁 Diagnóstico completado"
```

### Reparación Automática

```bash
#!/bin/bash
# hsm-repair.sh

echo "🔧 Reparación automática SoftHSM"
echo "==============================="

# 1. Recrear directorios
echo "📁 Creando estructura de directorios..."
mkdir -p ~/.softhsm/tokens
chmod 755 ~/.softhsm/tokens

# 2. Regenerar configuración
echo "⚙️ Regenerando configuración..."
cat > ~/.softhsm2.conf << EOF
directories.tokendir = $HOME/.softhsm/tokens/
objectstore.backend = file
log.level = INFO
slots.removable = false
slots.mechanisms = ALL
EOF

export SOFTHSM2_CONF="$HOME/.softhsm2.conf"

# 3. Limpiar tokens corruptos
echo "🧹 Limpiando tokens corruptos..."
find ~/.softhsm/tokens -name "*.lock" -delete 2>/dev/null || true
find ~/.softhsm/tokens -size 0 -delete 2>/dev/null || true

# 4. Test básico
echo "🧪 Realizando test básico..."
if softhsm2-util --show-slots; then
    echo "✅ Reparación exitosa"
else
    echo "❌ Reparación fallida"
    exit 1
fi

# 5. Crear token de prueba si no existe
if ! softhsm2-util --show-slots | grep -q "Initialized"; then
    echo "🔑 Creando token de prueba..."
    softhsm2-util --init-token --slot 0 --label "TestToken" --pin 1234 --so-pin 5678
fi

echo "🎉 Reparación completada"
```

## Integración con Aplicaciones

### Configuración para Aplicaciones Go

```go
// hsm_config.go
package main

import (
    "os"
    "path/filepath"
)

type HSMConfig struct {
    ModulePath string
    TokenLabel string
    UserPIN    string
    SlotID     uint
}

func GetHSMConfig() *HSMConfig {
    // Detectar SO y configurar ruta del módulo
    var modulePath string
    if runtime.GOOS == "windows" {
        modulePath = "C:\\SoftHSM2\\lib\\softhsm2-x64.dll"
    } else {
        modulePath = "/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so"
    }
    
    // Override con variable de entorno si existe
    if envPath := os.Getenv("PKCS11_MODULE"); envPath != "" {
        modulePath = envPath
    }
    
    return &HSMConfig{
        ModulePath: modulePath,
        TokenLabel: os.Getenv("HSM_TOKEN_LABEL"),
        UserPIN:    os.Getenv("HSM_USER_PIN"),
        SlotID:     0, // Default slot
    }
}
```

### Variables de Entorno para Producción

```bash
# production-env.sh
export SOFTHSM2_CONF="/opt/ssm/config/softhsm2.conf"
export PKCS11_MODULE="/usr/lib/x86_64-linux-gnu/softhsm/libsofthsm2.so"
export HSM_TOKEN_LABEL="SSM-Production"
export HSM_USER_PIN="$(cat /opt/ssm/secrets/hsm.pin)"
export HSM_SLOT=0

# Logging para troubleshooting (solo desarrollo)
# export PKCS11SPY=/usr/lib/x86_64-linux-gnu/pkcs11-spy.so
```

## Mejores Prácticas

### Seguridad

```bash
# 1. Usar PINs seguros generados aleatoriamente
openssl rand -base64 12

# 2. Cambiar PIN por defecto
pkcs11-tool --slot 0 --login --pin 1234 --change-pin

# 3. Backup seguro de tokens
tar -czf tokens_backup_$(date +%Y%m%d).tar.gz ~/.softhsm/tokens/
gpg -c tokens_backup_$(date +%Y%m%d).tar.gz
rm tokens_backup_$(date +%Y%m%d).tar.gz

# 4. Verificar integridad
find ~/.softhsm/tokens -type f -exec sha256sum {} \; > tokens_checksums.txt
```

### Performance

```bash
# 1. Configurar para mejor rendimiento
cat > ~/.softhsm2.conf << EOF
directories.tokendir = /dev/shm/softhsm/tokens/  # Usar RAM disk
objectstore.backend = file
log.level = ERROR  # Reducir logging
slots.removable = false
EOF

# 2. Pre-calentar el HSM
for i in {1..10}; do
    pkcs11-tool --slot 0 --login --pin 1234 --list-objects >/dev/null 2>&1
done
```

### Monitoreo en Producción

```bash
#!/bin/bash
# hsm-healthcheck.sh - Para usar en monitoring

SLOT=0
PIN=1234
CRITICAL_OBJECTS_MIN=5  # Mínimo de objetos esperados

# Test básico de conectividad
if ! pkcs11-tool --slot $SLOT --list-objects >/dev/null 2>&1; then
    echo "CRITICAL: No se puede conectar al HSM"
    exit 2
fi

# Verificar que el token no esté lleno
token_info=$(pkcs11-tool --slot $SLOT --token-info 2>/dev/null)
if echo "$token_info" | grep -q "Total.*space.*0"; then
    echo "CRITICAL: Token sin espacio disponible"
    exit 2
fi

# Contar objetos críticos
object_count=$(pkcs11-tool --slot $SLOT --login --pin $PIN --list-objects --type secrkey 2>/dev/null | grep "Secret Key" | wc -l)
if [ "$object_count" -lt "$CRITICAL_OBJECTS_MIN" ]; then
    echo "WARNING: Pocos objetos en el token ($object_count < $CRITICAL_OBJECTS_MIN)"
    exit 1
fi

echo "OK: HSM funcionando correctamente ($object_count objetos)"
exit 0
```