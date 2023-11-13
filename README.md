# Installation

```shell
git checkout server
git pull
```

# Configuration

In order to configure the server, several steps are required.
```shell
#Generate certification authority
openssl genpkey -algorithm RSA -out ./cert/ca.key
openssl req -new -x509 -key ./cert/ca.key -out ./cert/ca.crt

# Generate server certificate
openssl genpkey -algorithm RSA -out ./cert/server.key
openssl req -new -key ./cert/server.key -out ./cert/server.csr
openssl x509 -req -in ./cert/server.csr -CA ./cert/ca.crt -CAkey ./cert/ca.key -CAcreateserial -out ./cert/server.crt

# Generate client certificate
openssl genpkey -algorithm RSA -out ./cert/client.key
openssl req -new -key ./cert/client.key -out ./cert/client.csr
openssl x509 -req -in ./cert/client.csr -CA ./cert/ca.crt -CAkey ./cert/ca.key -CAcreateserial -out ./cert/client.crt
openssl x509 -in ./cert/client.crt -outform PEM -out ./cert/client.pem

# Generate the encryption key
openssl rand -out ./cert/encrypt.key 32
```

# Run

```shell
./server
```