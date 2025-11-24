# Despliegue del SSM con Aether Core

![alt text](../img/image4.png)
Figura 1: Arquitectura del SSM integrado con UDM y Webconsole

A continuación una serie de pasos a seguir para desplegar el SSM y ser utilizado por los componentes de Aether.

## 1. Requisitos Previos Sugeridos

1- Server Ubuntu 24.04 LTS (4 núcleos, 16 GB RAM, 60 GB disco)

Se recomiendan estos recursos para un entorno más estables, pero se han hecho simulaciones con menos recursos (2 núcleos y 8 GB RAM) y funciona correctamente para entornos de pruebas y simulacione pequeñas.

2- MongoDB 7.0 Community Edition desplegado como ReplicaSet

(Ver [Instalación de MongoDB en Ubuntu 24.04](../../commandsAndManagament/installing_mongodb_ubuntu24.04.md))

3- SoftHSM2 instalado y configurado

(Ver [Instalación de SoftHSM2](../../commandsAndManagament/installing_softHSM.md))

(Ver [Configuración de SoftHSM2](../../commandsAndManagament/softhsm_pkcs11_commands.md))

## 2. Iniciar el SSM

1- Descargar el binario del SSM desde la página de releases de GitLab o compilarlo desde el código fuente.

Si descarga el código fuente simplemente dirijase al directorio raíz del proyecto y ejecute:

```bash
go build -o ssm ./ssm.go
```

2- Crear el archivo de configuración config.yml con los parámetros necesarios para la conexión a SoftHSM2 y MongoDB, así como los parámetros de seguridad.

Desde el código fuente puede utilizar la plantilla de configuración que se encuentra en el directorio `factory` y adaptarla a sus necesidades.

```bash
vim factory/ssmConfig.yml
```

## 3. Desplegar Aether Core

1- Generar los certificados TLS para el SSM y  los demás componentes (si se desea utilizar HTTPS).

Ver [Generación de certificados TLS](../../commandsAndManagament/gen_certificates_mTLS1.3.md)

Utilice los certificados server.* para el ssm.

2- Desplegar el SSM con la configuración creada.

```bash
./ssm --cfg factory/ssmConfig.yml
```

3- Desplegar el Core de Aether

Este paso depende de la forma en que haya desplegado Aether Core. Le sugerimos utilizar Aether OnRamp siguiendo la guía oficial de Aether y la documentación generada por nosotros.

Antes de hacer el deploy de Aether Core, debe generar nuevos secrets para el Webconsole y UDM.

Para ello utilice la siguiente plantilla de los comandos a ejecutar. 

```bash
kubectl create secret generic udm-tls \
  --from-file=tls.crt=/var/certs/client-udm.crt \
  --from-file=tls.key=/var/certs/client-udm.key \
  --from-file=ca.crt=/var/certs/ca.crt \
  --namespace aether-5gc

kubectl create secret generic udm-credentials \
  --from-literal=service_id='<your_service_id>' \
  --from-literal=password='<your_password>' \
  --namespace aether-5gc

kubectl create secret generic webui-tls \
  --from-file=tls.crt=/var/certs/client-webconsole.crt \
  --from-file=tls.key=/var/certs/client-webconsole.key \
  --from-file=ca.crt=/var/certs/ca.crt \
  --namespace aether-5gc

kubectl create secret generic webui-credentials \
  --from-literal=service_id='<your_service_id>' \
  --from-literal=password='<your_password>' \
  --namespace aether-5gc
```
Para este depliegue es necesario utilizar un helm charts diferentes a los que vienen por defecto en Aether OnRamp, ya que es necesario configurar los parámetros de seguridad para que el Webconsole y UDM puedan comunicarse con el SSM.

Estos helm charts personalizados aún no están disponibles públicamente, por lo que debe crearlos usted mismo siguiendo la documentación oficial de Aether OnRamp y adaptando los valores necesarios para la configuración del SSM. Una implementación crítica pendiente es exponer estos charts en un repositorio público para facilitar su uso.

Asegurese de configurar correctamente los parámetros de conexión al SSM en los valores del helm chart, incluyendo la URL de conexión a MongoDB, los certificados TLS y las credenciales necesarias.

A continuación se muestra un ejemplo de cómo se verían los valores del helm chart para el UDM y el Webconsole:

```yaml
webui:
      tlsSecretName: "webui-tls"
      tlsSecretKeys:
        cert: tls.crt
        key: tls.key
        ca: ca.crt
      credentialsSecretName: "webui-credentials"
      ssm:
        enable: true
      credentials:
        serviceIdKey: service_id
        passwordKey: password
      serviceType: NodePort
      deploy: true
      urlport:
        port: 5000
        
        nodePort: 30001
      rest:
        port: 5001
        
        nodePort: 30002
      ingress:
        enabled: false
        hostname: aether.local
        path: /
        pathType: Prefix
        # extraHosts:
          # - host: aether.local
          #   path: /
      cfgFiles:
        webuicfg.yaml:
          info:
            version: 1.0.0
            description: WebUI initial local configuration
          configuration:
            # 5G mode
            spec-compliant-sdf: false
            enableAuthentication: false
            send-pebble-notifications: false
            # MongoDB configuration
            mongodb:
              name: aether
              url: "mongodb://mongodb-arbiter-headless"
              authKeysDbName: authentication
              authUrl: "mongodb://mongodb-arbiter-headless"
              webuiDbName: aether
              webuiDbUrl: "mongodb://mongodb-arbiter-headless"
              checkReplica: false
            enableAuthentication: false
            spec-compliant-sdf: false
            webui-tls:
              pem: /var/run/certs/tls.crt
              key: /var/run/certs/tls.key
            nfconfig-tls:
              pem: /var/run/certs/tls.crt
              key: /var/run/certs/tls.key
            managedByConfigPod:
              syncUrl: http://roc-service:8080/sync
              enabled: false
            ssm:
              ssm-uri: https://192.168.12.16:9000
              allow-ssm: true
              tls-insecure: true
              ssm-synchronize:
                enable: true
                interval-minute: 30
                max-keys-create: 5
                delete-missing: true
                max-sync-keys: 5
                max-sync-users: 5
                max-sync-rotations: 5
              m-tls:
                crt: /etc/webui/tls/tls.crt
                key: /etc/webui/tls/tls.key
                ca: /etc/webui/tls/ca.crt
              # read in cascade if gcm is false then cbc is used if gcm is true cbc is ignored, both is false no encryption is used (disabled SSM jjjj)
              is-encrypt-aes-gcm: true
              is-encrypt-aes-cbc: false
            metricsPort: ":8080"

    udm:
      deploy: true
      serviceType: ClusterIP
      credentialsSecretName: "udm-credentials"
      tlsSecretName: "udm-tls"
      tlsSecretKeys:
        cert: tls.crt
        key: tls.key
        ca: ca.crt
      ssm:
        enable: true
      credentials:
        serviceIdKey: service_id
        passwordKey: password
      sbi:
        port: 29503
      cfgFiles:
        udmcfg.yaml:
          info:
            version: 1.0.0
            description: UDM initial local configuration
          configuration:
            nrfUri: https://nrf:29510
            webuiUri: https://webui.aether-5gc:5001
            enableNrfCaching: true
            nrfCacheEvictionInterval: 900
            serviceList:
              - nudm-sdm
              - nudm-uecm
              - nudm-ueau
              - nudm-ee
              - nudm-pp
            sbi:
              scheme: https
              bindingIPv4: "0.0.0.0"
              registerIPv4: udm
              tls:
                pem: /var/run/certs/tls.crt
                key: /var/run/certs/tls.key
            keys:
              udmProfileAHNPublicKey: 5a8d38864820197c3394b92613b20b91633cbd897119273bf8e4a6f4eec0a650
              udmProfileAHNPrivateKey: c53c22208b61860b06c62e5406a7b330c2b577aa5558981510d128247d38bd1d
              udmProfileBHNPublicKey: 0472DA71976234CE833A6907425867B82E074D44EF907DFB4B3E21C1C2256EBCD15A7DED52FCBB097A4ED250E036C7B9C8C7004C4EEDC4F068CD7BF8D3F900E3B4
              udmProfileBHNPrivateKey: F1AB1074477EBCC7F554EA1C5FC368B1616730155E0041AC447D6301975FECDA
            ssm:
              enable: true
              host: https://192.168.12.16:9000
              tls_insecure: true
              m-tls:
                crt: /etc/udm/tls/tls.crt
                key: /etc/udm/tls/tls.key
                ca: /etc/udm/tls/ca.crt
            metricsPort: ":8080" 
```

Con todo este implementado ya debería tener el SSM funcionando y listo para ser utilizado por el UDM y el Webconsole de Aether Core. Solo tiene que ejecutar el comando de despliegue del Aether OnRamp con los nuevos charts y valores configurados.

## Sugerencias para el despliegue

1- Automatizar pasos utilizando Ansible
2- Implementar CI/CD en el SSM
3- Exponer los helm charts personalizados en un repositorio público y documentar cada versión de los helm charts