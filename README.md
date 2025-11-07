# SSM Aether

This SSM is writing in Go, the goal of this project is to provide a point to comunicate to a SoftHSM.

## Features

- API HTTP or HTTPS, that can be used to comunicate with SoftHSM
- Support for PKCS#11 operations
- We can storage keys that can be used by other services for example desencrypt/encrypt data
- Support decrypt operation using keys stored in SoftHSM. Support symmetric keys like as AES128, AES192, AES256, DES, 3DES
- Support encrypt operation using keys stored in SoftHSM. Support symmetric keys like as AES128, AES192, AES256, DES, 3DES
- Generate models and http client using OpenAPI Generator
- Support secure comunication using mTLS 1.3
- Support AES256-GCM encryption and decryption with IV and AAD
- Support audit logging for security-related events
- Support rate limiting to prevent common attacks
- Support CORS for web applications
- Support authentication and authorization mechanisms

See more details in docs folder.
