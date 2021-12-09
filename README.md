# Chat cifrado

## ¿Cómo hicieron el programa?

Para el programa decidimos usar, en lugar de Java, el lenguaje de programación Go. Esto por tres principales razones: Go tiene una mucho menor verbosidad, tiene un manejo muy sencillo para la concurrencia (era necesario, pues necesitabamos un hilo que recibiera mensajes y uno que los escribiera), y, siendo un lenguaje con el que tenemos menos experiencia, nos permitía aprender más cosas nuevas.

Respecto al flujo del programa, toda la comunicación se hace por TCP directamente entre los dos computadores que harán parte del chat, es decir, no hay un servidor intermedio. Uno de los dos computadores debe elegir ser el host del chat, y el otro se conecta por medio de su dirección ip y puerto, que se pasan por variables de entorno. Una vez establecida la comunicación, se procede al negociado de claves por medio de Diffie Hellman. Como Go no tiene soporte por defecto para Diffie Hellman, usamos la siguiente implementación, que ofrecía una API sencilla de utilizar: https://github.com/monnand/dhkx . El "intercambio" de claves se hace usando el grupo por defecto de la implementación, es decir, el grupo 14.
La clave negociada tiene 256 bytes, pero como para AES-128 solo necesitamos 16 bytes (128 bits = 16 bytes), tomamos este prefijo de la clave negociada.

Una vez finalizada la negociación de la clave, cliente y servidor inician el hilo encargado de recibir mensajes, que buffer de 8192 bytes para intentar leer lo que haya en el canal, y posteriormente desencripta la cadena de bytes con la clave negociada. El hilo de UI es el encargado de enviar mensajes: espera a que el usuario escriba un mensaje por consola con un fin de línea, lo convierte en cadena de bytes, lo encripta, y lo envía por la conexión. El encriptado y desencriptado usa AES-128 en modo CTR, usando por simplicidad como IV un prefijo del tamaño del bloque de la clave intercambiada.

## Dificultades del proyecto.

- El uso de un lenguaje de programación nuevo.
- Go no ofrece soporte nativo para Diffie-Hellman, por lo que tuvimos que buscar en internet implementaciones.
- Manejo del buffer de entrada para leer los bytes recibidos.

## Conclusiones

Gracias al desarrollo del proyecto tuvimos un acercamiento más real al uso de algoritmos para comunicación segura, como pueden ser Diffie-Hellman y AES. Desarrollarlo en un lenguaje nuevo nos ayudó a desarrollar las capacidades de aprender cosas nuevas. En síntesis podemos ver que gracias al cifrado, tenemos una capa de protección en caso de que un atacante malintencionado quiera interceptar nuestros mensajes.
