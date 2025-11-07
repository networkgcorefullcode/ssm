# Release v1.5.0 doc Español

## Características

- API HTTP o HTTPS, que puede ser utilizada para comunicarse con HSM
- Soporte para operaciones PKCS#11
- Podemos almacenar claves que pueden ser utilizadas por otros servicios, por ejemplo para descifrar/cifrar datos
- Soporte para operación de descifrado usando claves almacenadas en HSM. Soporta claves simétricas como AES128, AES192, AES256, DES, 3DES
- Soporte para operación de cifrado usando claves almacenadas en HSM. Soporta claves simétricas como AES128, AES192, AES256, DES, 3DES
- Generación de modelos y cliente http usando OpenAPI Generator
- Soporte para comunicación segura usando mTLS 1.3
- Soporte para cifrado y descifrado AES256-GCM con IV y AAD
- Soporte para registro de auditoría de eventos relacionados con la seguridad
- Soporte para limitación de velocidad para prevenir ataques comunes
- Soporte para CORS para aplicaciones web
- Soporte para mecanismos de autenticación y autorización

## Que falta, qué viene después para el nuevo release 
- Mejoras en la documentación
- Corrección de bugs según se vayan encontrando
- Mejoras en pruebas unitarias y de integración, faltan agregar test para un mejor desarrollo continuo
- SIEM implementado según requerimientos de seguridad
- Endpoint para reiniciar los password de acceso, sin necesidad de hacerlo manualmente
- Implementación de una tecnología para visualizar los logs de auditoría de forma centralizada
- Aplicar estos cambios en las demás NFs (webconsole y udm)