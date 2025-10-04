
``` bash
sudo curl -X POST http://dummy/store-key \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K43DES",
    "id": "1",
    "key_value": "1234567890abcdef1234567890abcdef1234567890abcdef"
    "key_type": "DES3",
  }'
```

``` bash
sudo curl -X POST http://dummy/decrypt\
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K43DES",
    "cipher_b64": "DC1D1221FA595EBE23E93399D48CBEBF",
    "key_type": "DES3",
    encryption_algoritme: 4
  }'
```

type DecryptRequest struct {
    // Etiqueta de la clave para descifrar
    KeyLabel string `json:"key_label"`
    // Datos cifrados en Base64
    CipherB64 string `json:"cipher_b64"`
    // Vector de inicialización en Base64 (mismo usado para cifrar)
    IvB64 string `json:"iv_b64,omitempty"`
    // ID opcional para tracking
    Id *int32 `json:"id,omitempty"`
    // Details for the encryption algoritme
    EncryptionAlgoritme int `json:"encryption_algoritme"`
}