# Funcionamiento actual del SSM

![alt text](../img/image.png)
Figura 1: Diagrama de alto nivel del funcionamiento del SSM

La figura 1 muestra un diagrama de alto nivel del funcionamiento del SSM. 

1- API HTTP/HTPS: Esta API esta escrita en gin y con openapi generator se genera el codigo de golang para que las aplicaciones cliente puedan interactuar con el SSM.

2- PCKS11 Manager: El PCKS11 Manager es responsable de gestionar los módulos PKCS#11 que se utilizan para interactuar con los dispositivos de seguridad. Proporciona una interfaz para cargar y administrar estos módulos. Aquí se realizan operaciones como la inicialización del módulo, la gestión de sesiones y la ejecución de comandos PKCS#11 para encriptar, desencriptar, firmar y verificar datos. Aquí también se manejan los pools de sesiones pkcs11 para optimizar el rendimiento y las operaciones naturalmente concurrentes. 

3- HSM: Parte encargada de manejar las operaciones criptográficas de manera segura. Puede ser un HSM físico o un HSM virtual como en nuestro caso actual que se ha utilizado SoftHSM para casos de prueba y desarrollo. 

![alt text](../img/image1.png)
Figura 2: Ilustra los componentes de seguridad implementados en el SSM

1- La sección de init SSM se ejecuta cada vez que se inicia el SSM. Esta se encarga de dos cosas fundamentales. Primero se crean los usuarios udm y webconsole sino estan creados, estos se guardan en la base de datos (los passwords se cifran con una llave AES256 generada para este proceso en específico) y se exponen estos datos sensibles en un txt en el path /tmp/user-secret-ssm/, el cual debe ser manejado con cuidado. Estos datos serán los necesarios para hacer el login en el Webconsole y en el UDM. Por último se crear llaves de cifrado asimétrico para los procesos de audit y logging, si ya estan creadas estas llaves no se vuelven a crear.

2- El siguiente bloque es el principal, él cual brinda seguridad al API HTTP que se expone.

Todas las request al API deben proporcionar los certificados TLS correspondientes, en caso de no ser certificados válidos o no proporcionar ninguno se rechazará la conexión. 

Los clientes antes de realizar las operaciones deberán hacer login con sus credenciales, en caso de ser validas se les proporcionará un token JWT que deberán usar en las siguientes requests.

El proceso de audit siempre se ejecutará en cada petición al API y guardará la información en la base de datos, con una firma digital generada con las llaves creadas en el init SSM.

En el auth se verificará el token JWT proporcionado por el cliente, en caso de ser válido se permitirá el acceso a los recursos solicitados. Aquí se verifican los roles, en caso de ser webconsole se permitirá el acceso a todos los recursos, en caso de ser udm se permitirá el acceso solo a los recursos permitidos para este rol, hasta ahora solo tiene permitido utilizar el endpoint para desencriptar secretos.

Todos los procesos que implequen cifrado o firma digital utilizarán el PCKS11 Manager para interactuar con el HSM y realizar las operaciones criptográficas de manera segura.

![alt text](../img/image2.png)
Figura 3: Diagrama de secuencia de una petición al API HTTP del SSM por parte del UDM

![alt text](../img/image3.png)
Figura 4: Diagrama de secuencia de una petición al API HTTP del SSM por parte del Webconsole

![alt text](../img/image4.png)
Figura 5: Arquiterctura de la integración del SSM con el UDM y Webconsole