## go::abac::cerbos

### introduction

this repository hosts the code for the authenication microservices for kale::capital.
it exposes a `rest api` that allows clients to authenticate themselves.
it will also expose a `grpc` server that internal microservices can call to verify paseto tokens or refresh expired tokens using a refresh token mechanism

### tech stack

- go
- paseto
- postgres
- docker(+compose)
- make
- air
- fiber
