# Release v1.5.0 doc English

## Features

- HTTP or HTTPS API, which can be used to communicate with HSM
- Support for PKCS#11 operations
- We can store keys that can be used by other services, for example to decrypt/encrypt data
- Support for decryption operation using keys stored in HSM. Supports symmetric keys such as AES128, AES192, AES256, DES, 3DES
- Support for encryption operation using keys stored in HSM. Supports symmetric keys such as AES128, AES192, AES256, DES, 3DES
- Memory hardening for greater security (using mlockall)
- Model and HTTP client generation using OpenAPI Generator
- Support for secure communication using mTLS 1.3
- Support for AES256-GCM encryption and decryption with IV and AAD
- Support for audit logging of security-related events
- Support for rate limiting to prevent common attacks
- Support for CORS for web applications
- Support for authentication and authorization mechanisms

## What's missing, what's coming next for the new release
- Documentation improvements
- Bug fixes as they are found
- Improvements in unit and integration tests, need to add tests for better continuous development
- SIEM implemented according to security requirements
- Endpoint to reset access passwords, without having to do it manually
- Implementation of a technology to visualize audit logs in a centralized way
- Apply these changes to the other NFs (webconsole and udm)