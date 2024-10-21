#  “Word of Wisdom” tcp server.
• TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
• The choice of the POW algorithm should be explained.
• After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
• Docker file should be provided both for the server and for the client that solves the POW challenge


## How to run

Run the client and server.
```
docker-compose -f docker/docker-compose.yml up -d
```

## Chosen algorithm

SHA-256 algorithm was chosen for the Proof of Work algorithm. It is a widely used cryptographic hash function that produces a 256-bit (32-byte) hash value. It is used in various security applications and protocols. 
It is also the basis for Bitcoin mining.

Advantages of SHA-256:
* Resistance to Spam and DoS Attacks Hashcash was designed to mitigate spam by requiring a sender to do computational work (solving a cryptographic puzzle) before sending a message or request. This makes bulk spamming or DoS attacks expensive and inefficient because attackers would need to solve multiple hash puzzles to overwhelm the system
* The difficulty of the SHA-256 puzzle can be adjusted dynamically. For example, if more computational power becomes available (due to advances in hardware or distributed botnets), the difficulty of finding a valid hash can be increased.
* SHA-256 is considered secure for the foreseeable future, with no known vulnerabilities that would make Hashcash ineffective
* Hashcash using SHA-256 is easy to implement with existing cryptographic libraries. Most modern programming languages provide native support for the SHA-256 hash function, making it straightforward to integrate Hashcash mechanisms into various systems.
