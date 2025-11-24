# Release v1.6.1 doc English

## Features

- HTTP or HTTPS API, which can be used to communicate with HSM
- Support for PKCS#11 operations
- We can store keys that can be used by other services, for example to decrypt/encrypt data
- Support for decryption operation using keys stored in HSM. Supports symmetric keys such as AES128, AES192, AES256, DES, 3DES
- Support for encryption operation using keys stored in HSM. Supports symmetric keys such as AES128, AES192, AES256, DES, 3DES
- Memory hardening for enhanced security (using mlockall)
- Model and HTTP client generation using OpenAPI Generator
- Support for secure communication using mTLS 1.3
- Support for AES256-GCM encryption and decryption with IV and AAD
- Support for audit logging of security-related events
- Support for rate limiting to prevent common attacks (this functionality was removed in this release, as it presented implementation failures and it was decided not to include it in this release)
- Support for CORS for web applications
- Support for authentication and authorization mechanisms
- This release is used in the Aether core and has been tested with simulations to verify correct operation