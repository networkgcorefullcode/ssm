# 🎯 SSM API Testing - Instrucciones Rápidas

## ✅ Lo que se ha creado:

### 1. **Ejemplos JSON** (18 archivos)
   - 📁 `docs/json_examples/requests/` - 8 archivos de solicitud
   - 📁 `docs/json_examples/responses/` - 10 archivos de respuesta

### 2. **Scripts de Testing** (6 archivos en `utils/`)
   - `api_test.ps1` - Script PowerShell para Windows
   - `api_test.sh` - Script Bash para Linux/macOS
   - `test_workflow.ps1` - Workflow completo (PowerShell)
   - `test_workflow.sh` - Workflow completo (Bash)
   - `README.md` - Documentación completa
   - `QUICK_EXAMPLES.md` - Ejemplos rápidos
   - `SETUP_SUMMARY.md` - Resumen del setup

## 🚀 Uso Rápido (Windows/PowerShell)

### Navega a la carpeta utils:
```powershell
cd d:\projects\aether-forks-gitlab\ssm-gitlab\ssm\utils
```

### Prueba básica:
```powershell
# Health check
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"
```

### Con archivo JSON:
```powershell
# Generar clave AES
.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json"
```

### Workflow completo:
```powershell
.\test_workflow.ps1
```

## 📝 Sintaxis del Script

### PowerShell:
```powershell
.\api_test.ps1 -Endpoint "<endpoint>" -Method "<method>" [-JsonFile "<path>"] [-BaseUrl "<url>"]
```

### Ejemplos comunes:

```powershell
# Sin body JSON (GET o POST sin datos)
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST"

# Con body JSON
.\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"

# Servidor diferente
.\api_test.ps1 -Endpoint "/health-check" -Method "GET" -BaseUrl "http://localhost:9000"
```

## 📋 Endpoints Disponibles

| Endpoint | Método | Necesita JSON | Archivo Ejemplo |
|----------|--------|---------------|-----------------|
| `/health-check` | GET/POST | No | - |
| `/encrypt` | POST | Sí | `encrypt_request.json` |
| `/decrypt` | POST | Sí | `decrypt_request.json` |
| `/generate-aes-key` | POST | Sí | `gen_aes_key_request.json` |
| `/generate-des-key` | POST | Sí | `gen_des_key_request.json` |
| `/generate-des3-key` | POST | Sí | `gen_des3_key_request.json` |
| `/store-key` | POST | Sí | `store_key_request.json` |
| `/get-data-keys` | POST | Sí | `get_data_keys_request.json` |
| `/get-key` | POST | Sí | `get_key_request.json` |
| `/get-all-keys` | POST | No | - |

## 🔥 Casos de Uso Comunes

### 1. Verificar que el servidor está funcionando:
```powershell
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"
```

### 2. Generar una clave AES de 256 bits:
```powershell
.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json"
```

### 3. Ver todas las claves del HSM:
```powershell
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST"
```

### 4. Cifrar datos:
```powershell
.\api_test.ps1 -Endpoint "/encrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\encrypt_request.json"
```

### 5. Descifrar datos:
```powershell
# Primero actualiza decrypt_request.json con el cipher e IV que obtuviste del paso anterior
.\api_test.ps1 -Endpoint "/decrypt" -Method "POST" -JsonFile "..\docs\json_examples\requests\decrypt_request.json"
```

## 🛠️ Modificar los Ejemplos JSON

Los archivos JSON están en:
- **Requests:** `docs\json_examples\requests\`
- **Responses:** `docs\json_examples\responses\` (solo como referencia)

### Ejemplo: Modificar gen_aes_key_request.json

```json
{
  "id": 123,        // Cambia el ID
  "bits": 256       // Cambia el tamaño (128, 192, 256)
}
```

### Ejemplo: Modificar encrypt_request.json

```json
{
  "key_label": "mi-clave-personalizada",  // Tu label
  "plain": "48656c6c6f",                   // Datos en HEX
  "encryption_algorithm": 1                // 1=AES, 3=DES, 4=DES3
}
```

## 🎬 Workflow Completo

Para probar todo el flujo de cifrado/descifrado:

```powershell
.\test_workflow.ps1
```

Este script:
1. ✅ Verifica que el servidor está funcionando
2. 🔑 Genera una clave AES
3. 🔒 Cifra datos
4. 🔓 Descifra datos
5. 📋 Lista todas las claves

## 💡 Tips

### Ver ayuda:
```powershell
.\api_test.ps1 -Help
```

### Cambiar servidor por defecto:
```powershell
# Opción 1: En cada comando
.\api_test.ps1 -Endpoint "/health-check" -Method "GET" -BaseUrl "http://192.168.1.100:8080"

# Opción 2: Variable de entorno (edita el script)
$BaseUrl = "http://192.168.1.100:8080"
```

### Ejecutar múltiples tests:
```powershell
# Crear un archivo test_suite.ps1
.\api_test.ps1 -Endpoint "/health-check" -Method "GET"
.\api_test.ps1 -Endpoint "/generate-aes-key" -Method "POST" -JsonFile "..\docs\json_examples\requests\gen_aes_key_request.json"
.\api_test.ps1 -Endpoint "/get-all-keys" -Method "POST"
```

## ❓ Solución de Problemas

### Error: "no se puede cargar... no está firmado digitalmente"
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Error: "No se puede encontrar la ruta"
```powershell
# Asegúrate de estar en la carpeta utils
cd d:\projects\aether-forks-gitlab\ssm-gitlab\ssm\utils
```

### El servidor no responde:
1. Verifica que el servidor SSM está corriendo
2. Verifica el puerto (por defecto 8080)
3. Prueba con: `.\api_test.ps1 -Endpoint "/health-check" -Method "GET"`

## 📚 Más Información

- `README.md` - Documentación completa de los scripts
- `QUICK_EXAMPLES.md` - Lista de todos los comandos disponibles
- `SETUP_SUMMARY.md` - Resumen completo del setup

---

**¡Todo listo para empezar a probar la API! 🎉**

Ejecuta: `.\api_test.ps1 -Endpoint "/health-check" -Method "GET"`
