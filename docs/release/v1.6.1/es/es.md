# Release v1.6.1 doc Español

## Características

- API HTTP o HTTPS, que puede ser utilizada para comunicarse con HSM
- Soporte para operaciones PKCS#11
- Podemos almacenar claves que pueden ser utilizadas por otros servicios, por ejemplo para descifrar/cifrar datos
- Soporte para operación de descifrado usando claves almacenadas en HSM. Soporta claves simétricas como AES128, AES192, AES256, DES, 3DES
- Soporte para operación de cifrado usando claves almacenadas en HSM. Soporta claves simétricas como AES128, AES192, AES256, DES, 3DES
- Endurecimiento de memoria para mayor seguridad (usando mlockall)
- Generación de modelos y cliente http usando OpenAPI Generator
- Soporte para comunicación segura usando mTLS 1.3
- Soporte para cifrado y descifrado AES256-GCM con IV y AAD
- Soporte para registro de auditoría de eventos relacionados con la seguridad
- Soporte para limitación de velocidad para prevenir ataques comunes ( se quito esta funcionalidades en este release, ya que presento fallos en la implementación y se decidió no incluirla en este release)
- Soporte para CORS para aplicaciones web
- Soporte para mecanismos de autenticación y autorización
- Este release es el utilizado en el core de Aether y ha sido probado con simulaciones para comprobar el correcto funcionamiento