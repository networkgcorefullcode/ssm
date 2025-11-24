# mTLS1.3 certificates  generation

REMPLAZA LOS VALORES DE LOS CAMPOS SEGUN TU ORGANIZACION Y NECESIDADES
GENERA UN TOKEN EN SOFTHSM2 ANTES DE EJECUTAR ESTOS COMANDOS

```bash
sudo apt update
sudo apt install pkcs11-provider

sudo mkdir -p /var/certs
sudo tee /var/certs/openssl.cnf > /dev/null <<EOF
[ ca ]
default_ca = myca

[ myca ]
dir               = /var/certs
database          = \$dir/index.txt
new_certs_dir     = /var/certs/newcerts
certificate       = /var/certs/ca.crt
serial            = /var/certs/serial
default_md        = sha256
policy            = mypolicy
private_key       = pkcs11:object=RootCAKey;type=private;pin-value=1234
default_days      = 365
default_crl_days  = 30

[ mypolicy ]
commonName        = supplied
EOF

sudo mkdir -p /var/certs/newcerts
sudo touch /var/certs/index.txt
echo 1000 | sudo tee /var/certs/serial
```

```bash
# Init token and create a keypair for the Root CA
sudo pkcs11-tool --module /usr/lib/softhsm/libsofthsm2.so --slot 1121449042 --login --pin 1234 \
  --keypairgen --key-type rsa:4096 --id 01 --label "RootCAKey"
```

```bash
# Generate a self-signed Root CA certificate
sudo openssl req -engine pkcs11 -new -x509 -days 365 -sha256 \
  -keyform engine -key "pkcs11:object=RootCAKey;type=private" \
  -subj "/C=CU/ST=Havana/L=Marianao/O=MyOrg/OU=Security/CN=MyRootCA/emailAddress=admin@myorg.cu" \
  -addext "basicConstraints=critical,CA:TRUE,pathlen:1" \
  -addext "keyUsage=critical,keyCertSign,cRLSign" \
  -addext "subjectKeyIdentifier=hash" \
  -out /var/certs/ca.crt
```

```bash
# Create a token and keypair for the server certificate
sudo openssl req -new -newkey rsa:2048 -nodes -keyout /var/certs/server.key -out /var/certs/server.csr \
  -subj "/C=CU/ST=Havana/O=MyOrg/CN=myserver.local"
```

```bash
# Sign the server certificate with the Root CA
# Firmar el CSR del servidor con la Root CA almacenada en SoftHSM
sudo openssl ca -engine pkcs11   -config /var/certs/openssl.cnf   -keyform engine   -key "pkcs11:object=RootCAKey;type=private"   -in /var/certs/server.csr   -out /var/certs/server.crt -batch
```

```bash
sudo openssl req -new -newkey rsa:2048 -nodes -keyout /var/certs/client-udm.key -out /var/certs/client-udm.csr \
  -subj "/C=CU/ST=Havana/O=MyOrg/CN=Client1"
sudo openssl ca -engine pkcs11   -config /var/certs/openssl.cnf   -keyform engine   -key "pkcs11:object=RootCAKey;type=private"   -in /var/certs/client-udm.csr   -out /var/certs/client-udm.crt -batch
```

```bash
sudo openssl req -new -newkey rsa:2048 -nodes -keyout /var/certs/client-webconsole.key -out /var/certs/client-webconsole.csr \
  -subj "/C=CU/ST=Havana/O=MyOrg/CN=WebConsole"
sudo openssl ca -engine pkcs11   -config /var/certs/openssl.cnf   -keyform engine   -key "pkcs11:object=RootCAKey;type=private"   -in /var/certs/client-webconsole.csr   -out /var/certs/client-webconsole.crt -batch
```
