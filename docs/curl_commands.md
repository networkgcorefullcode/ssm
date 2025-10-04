
``` bash
sudo curl -X POST http://dummy/store-key \
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K43DES",
    "id": "1",
    "key_value": "1234567890abcdef1234567890abcdef1234567890abcdef",
    "key_type": "DES3"
  }'
```

``` bash
sudo curl -X POST http://dummy/decrypt\
  --unix-socket /var/run/socket.so \
  -H "Content-Type: application/json" \
  -d '{
    "key_label": "K43DES",
    "cipher_b64": "DC1D1221FA595EBE23E93399D48CBEBF",
    "encryption_algoritme": 4
  }'
```