# Generación de Certificados Autofirmados para HTTPS en Ubuntu 24.04

Este documento describe los pasos necesarios para generar certificados autofirmados que permitan la comunicación HTTPS en Ubuntu 24.04.

## Requisitos Previos

Asegúrate de tener OpenSSL instalado en tu sistema Ubuntu 24.04:

```bash
sudo apt update
sudo apt install openssl
```

## Método 1: Certificado Básico Autofirmado

### Generar clave privada y certificado en un solo paso

```bash
# Generar clave privada (2048 bits) y certificado autofirmado válido por 365 días
openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes
```

Durante el proceso se te pedirán los siguientes datos:

- **Country Name (2 letter code)**: CO (para Colombia)
- **State or Province Name**: Nombre de tu departamento/estado
- **City or Locality Name**: Nombre de tu ciudad
- **Organization Name**: Nombre de tu organización
- **Organizational Unit Name**: Nombre de tu departamento/unidad
- **Common Name**: **IMPORTANTE** - Debe ser el dominio o IP donde correrá el servidor (ej: localhost, 192.168.1.100, midominio.com)
- **Email Address**: Tu dirección de email

## Método 2: Proceso en Dos Pasos (Mayor Control)

### Paso 1: Generar la clave privada

```bash
# Generar clave privada RSA de 2048 bits
openssl genrsa -out server.key 2048

# O generar clave privada RSA de 4096 bits (más segura)
openssl genrsa -out server.key 4096
```

### Paso 2: Generar el certificado usando la clave privada

```bash
# Crear certificado autofirmado válido por 365 días
openssl req -new -x509 -key server.key -out server.crt -days 365
```

## Método 3: Usando Archivo de Configuración

### Crear archivo de configuración

Crea un archivo `server.conf` con el siguiente contenido:

```ini
[req]
default_bits = 2048
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = v3_req

[dn]
C = CO
ST = Cundinamarca
L = Bogota
O = Mi Organizacion
OU = Departamento IT
CN = localhost

[v3_req]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = *.localhost
IP.1 = 127.0.0.1
IP.2 = ::1
```

### Generar certificado usando el archivo de configuración

```bash
# Generar clave privada y certificado usando configuración
openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes -config server.conf
```

## Verificación del Certificado

### Ver detalles del certificado generado

```bash
# Mostrar información del certificado
openssl x509 -in server.crt -text -noout

# Verificar fechas de validez
openssl x509 -in server.crt -dates -noout

# Verificar el Common Name
openssl x509 -in server.crt -subject -noout
```

### Verificar la clave privada

```bash
# Mostrar información de la clave privada
openssl rsa -in server.key -text -noout

# Verificar que la clave privada coincida con el certificado
openssl x509 -noout -modulus -in server.crt | openssl md5
openssl rsa -noout -modulus -in server.key | openssl md5
# Los dos comandos anteriores deben devolver el mismo hash
```

## Configuración para Diferentes Escenarios

### Para uso en localhost

```bash
# Certificado específico para localhost
openssl req -x509 -newkey rsa:2048 -keyout localhost.key -out localhost.crt -days 365 -nodes \
  -subj "/C=CO/ST=Cundinamarca/L=Bogota/O=Local Development/CN=localhost"
```

### Para uso en red local (IP específica)

```bash
# Reemplaza 192.168.1.100 con tu IP real
openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -days 365 -nodes \
  -subj "/C=CO/ST=Cundinamarca/L=Bogota/O=Mi Organizacion/CN=192.168.1.100"
```

### Para múltiples dominios/IPs

Usa el método con archivo de configuración y modifica la sección `[alt_names]` para incluir todos los dominios e IPs necesarios.

## Configuración de Permisos

```bash
# Establecer permisos seguros para la clave privada
chmod 600 server.key

# Establecer permisos para el certificado
chmod 644 server.crt
```

## Integración con tu Aplicación Go

Para usar estos certificados en tu aplicación SSM, puedes configurar el servidor HTTPS de la siguiente manera:

```go
package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    // Configurar tus rutas aquí
    
    log.Println("Servidor iniciando en https://localhost:8443")
    err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", mux)
    if err != nil {
        log.Fatal("Error al iniciar servidor HTTPS: ", err)
    }
}
```

## Notas Importantes

1. **Advertencias del navegador**: Los certificados autofirmados generarán advertencias de seguridad en los navegadores web. Esto es normal.

2. **Validez limitada**: Configura la duración del certificado según tus necesidades (parámetro `-days`).

3. **Common Name**: Es crucial que el Common Name coincida exactamente con el dominio o IP desde donde accederás al servidor.

4. **Seguridad**: Los certificados autofirmados solo deben usarse para desarrollo o entornos internos controlados.

5. **Backup**: Guarda copias de seguridad de tus certificados y claves privadas en un lugar seguro.

## Comandos de Limpieza

```bash
# Remover certificados y claves generados
rm -f server.key server.crt localhost.key localhost.crt server.conf
```

## Troubleshooting

### Error: "unable to load Private Key"

- Verifica que el archivo de la clave privada existe y tiene los permisos correctos
- Asegúrate de que la clave no esté corrupta

### Error: "certificate verify failed"

- Verifica que el Common Name del certificado coincida con el hostname/IP usado
- Considera agregar el certificado a las autoridades certificadoras confiables del sistema

### Navegador muestra "conexión no es privada"

- Es normal para certificados autofirmados
- Puedes proceder haciendo clic en "Avanzado" → "Continuar al sitio"