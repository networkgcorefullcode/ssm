# List for Things To Do

- [x] Task 1 Implement feature for simetrict description that uses IV vector
- [x] Task 2 Implement rotation for simetric description (see in the webconsole)
- [x] Task 4 Update documentation for new simetric description features
- [x] Task 5 Update the open api generator to support new features
- [x] Task 5.1 Update the open api generator to support new features
- [x] Task 6 Review and refactor code for performance improvements
- [x] Task 7 Implement syncronization to aether core using the webconsole component
- [x] Task 7.1 Implement new api for syncronization functions
- [x] Task 7.2 Migrate to Gin framework for the HTTP API
- [ ] Task 8 Add security modules for HTTP API using middlewares for security
  - [x] Add authentication and authorization mechanisms
  - [x] Implement CORS
  - [x] Implement rate limiting to prevent common attacks
  - [x] Add audit logging for security-related events
  - [x] Implement TLS for secure communication (mTLS1.3)
  - [x] Update OpenAPI documentation to reflect security features
  - [ ] Write tests to verify the effectiveness of security measures
  - [ ] Test all security features and fix any vulnerabilities found
- [x] Task 9 Implement memory hardening techniques
- [ ] Task 10 Implement SIEM integration for audit logs
- [ ] Task 11 Performance testing and optimization, add benchmarks
- [x] Task 12 Add function to save data in a secure way using AES256-GCM (encrypt and decrypt)
- [ ] Task 12 Add a option to reset the user data for the login
- [ ] Task 13 Add a new function to save the ssm data in other softHSMv2 instance securely
- [ ] Task 13.1 Implement a frontend technology to see status information and logs
- [ ] Task 14 Final review and documentation update

## Requirements

### Criterios Técnicos

•Migración completa de AES128-ECB a AES-256-GCM
Se pueden trabajar con los siguientes algoritmos: DES, DES3, AES128 y AES256 tanto ECB como CBC. Por defecto para nuestros usuarios o nuevos subscriptores de nuestra red se utilizará el algoritmo AES256_CBC_PAD. Queda por implementar AES256_GCM y probar su implementación en el softHSMv2. Se le da soporte a otros algoritmos para usuarios legacy.

•Implementación de AAD robusto y IV único por operación
Implementado y probado, solo se pueden agregar llaves con keyLabels únicos y establecidos por nuestras constantes, tampoco se pueden repetir Id de un mismo keyLabel y se da soporte para que el sistema retorne un id válido en caso de que no se especifique este dato. Sobre los IV estos deben ser generados durante el proceso de encriptación, en nuestro caso se generan de forma aleatoria, donde la probabilidad de que se repitan es mínima, haciendo a estos únicos por operación.

•API gRPC UDM↔SSM con mTLS 1.3 operativa
Api http con openapi incluido para generar clientes para los demas servicios. GRPC no implementado solo sería necesario en casos que se necesite un mayor rendimiento en el numero de operaciones, pero el api http en Gin ofrece un buen rendimiento. Queda por implementar mTLS 1.3.

•Endurecimiento de memoria (mlock, zeroización) implementado
Paso final en todo el desarrollo, para evitar la exposición en memoria de datos sensibles. No implementado aún.

•Sistema de auditoría con logs firmados y SIEM
El sistema de auditoría está implementado de forma centralizada en un middleware del API, ya que a través del API se realizan todas las operaciones de la aplicación. También se firman estos logs utilizando el propio softHSMv2 con una llave asimétrica RSA256. Falta implementar el SIEM, donde a través de ciertas operaciones se generen alarmas, que pueden ser notificadas por diversas vías.

•Rendimiento: ≤ 5ms p95 para operaciones de desencriptación
Debe ser chequeado al final.

### Criterios de Seguridad

•Cumplimiento con checklist GSMA FS.31
Ver estos requisitos.

•Protección contra reuso de IV implementada
Ya implementado.

•Control de acceso RBAC operativo
Implementado y pendiente a ser probado en profundidad y a integrarse nuevos cambios de ser necesario. Para este sistema se implementan tokens JWT firmados por el propio softHSMv2 que incluyen los claims necesarios para el UDM y el Webconsole. El flujo que sigue este proceso es el siguiente -> En cuanto inicializa el SSM este debe crear un service_id con su password para los udm y webconsole. Estos datos son sensibles y se guardan en un archivo txt que debe ser guardado en un lugar seguro, luego serán utilizados por el UDM y el webconsole para hacer login en el API, si el loguin es correcto entonces se crea un token jwt firmado por el propio softHSMv2 utilizando firmas asimétricas RSA256 que tendrá un tiempo de expiración de 24 horas. Este token será utilizado entonces para autenticarse y poder hacer las demás operaciones. El role de UDM solo permitirá hacer operaciones en el API para desencriptar datos, ya que es lo unico que necesita para las operaciones de generar datos de autenticación. El role de Webconsole es el administrativo por lo que con este token se podrán hacer todas las operaciones.

•Rate limiting y protección DoS implementados
Implementado, pendiente a ser sometido a pruebas de estrés

•Plan de rotación de K4 automatizado
Implementado un plan de rotación de K4. Este plan es destructivo y tiene varias desventajas, queda pendiente analizar otros planes más modernos y su implementación. 
