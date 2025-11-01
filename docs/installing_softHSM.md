# Install SoftHSM

SoftHSM is a software implementation of a Hardware Security Module (HSM) that provides a secure environment for cryptographic operations. It is commonly used for testing and development purposes when a physical HSM is not available.

## Ubuntu 24.04

Installing softHSM

```bash
sudo apt update
sudo apt install softhsm2 libsofthsm2
softhsm2-util --show-slots
```

Installing pkcs11-tool

```bash
sudo apt update
sudo apt install opensc opensc-pkcs11
sudo apt install p11-kit
sudo apt install pkcs11-tool
sudo apt install libengine-pkcs11-openssl
```
