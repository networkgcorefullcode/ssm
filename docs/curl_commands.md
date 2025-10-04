
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