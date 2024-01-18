# TCP-server and client with protection from DDOS attacks

## 1. Getting started
### 1.1 Requirements
+ [Go 1.20+](https://go.dev/dl/) to run tests and server, client without Docker
+ [Docker](https://docs.docker.com/engine/install/) to run with Docker

### 1.2 Install dependencies (run without Docker):
```
make deps
```

### 1.3 Start server without Docker:
```
make run-server
```

### 1.4 Start client without Docker:
```
make run-client
```

### 1.5 Start server and client with docker-compose:
```
make start
```

### 1.6 Run tests:
```
make test
```

## 2. Task description
Design and implement “Word of Wisdom” tcp server.
+ TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
+ The choice of the POW algorithm should be explained.
+ After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
+ Docker file should be provided both for the server and for the client that solves the POW challenge


## 3. Choice of Hashcash as the POW algorithm.
### Why Hashcash?
Hashcash is a widely recognized PoW algorithm that was originally designed to combat email spam by requiring senders to prove that they have expended a certain amount of computational effort before sending an email. It has been adapted for various applications, including cryptocurrency mining (e.g., Bitcoin).

Advantages of Using Hashcash in this Project:

+ **Simplicity:** Hashcash is relatively simple to understand and implement. This simplicity makes it suitable for rapid development.
+ **Resource-Efficient:** Compared to more complex PoW algorithms used in cryptocurrencies, Hashcash is less resource-intensive. It doesn't require powerful mining hardware or significant computational resources. This makes it suitable for lightweight applications like protecting a TCP server from basic DDoS attacks.
+ **Customization:** Hashcash can be easily customized to adjust the level of difficulty for solving the challenge. In this project, the number of leading zeros in the proof can be adjusted to control the difficulty level.

## 4. Hash Function:
In the project, I used SHA-256 as the hash function. SHA-256 is widely used and considered secure.