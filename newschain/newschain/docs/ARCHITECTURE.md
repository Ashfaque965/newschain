

# NewsChain Architecture

## Overview

NewsChain is a hybrid Web2/Web3 news publishing platform designed to provide **tamper-evident, verifiable, and decentralized article publishing**.

Instead of storing complete articles on-chain, NewsChain stores only cryptographic proofs on the blockchain while keeping article content off-chain for efficiency and scalability.

This architecture provides:

* Fast API responses
* Low blockchain transaction costs
* Cryptographic integrity verification
* Decentralized content storage
* Transparent audit trails

---

# High-Level Architecture

```text
                           +----------------------+
                           |   Web / Mobile App   |
                           +----------+-----------+
                                      |
                                      |
                               HTTPS REST API
                                      |
                                      ▼
                         +-------------------------+
                         |      Go API Server      |
                         +-----------+-------------+
                                     |
      +------------------------------+------------------------------+
      |                              |                              |
      ▼                              ▼                              ▼
+-------------+              +---------------+             +----------------+
| Authentication|            | SQLite DB     |             | IPFS Node      |
| SIWE + JWT    |            | Article Store |             | Content Storage|
+-------------+              +---------------+             +----------------+
                                     |
                                     ▼
                         +--------------------------+
                         | Ethereum Smart Contract  |
                         |     News Registry        |
                         +--------------------------+
```

---

# Core Components

## Client Applications

Client applications interact with NewsChain through the REST API.

Possible clients include:

* Web frontend
* Mobile applications
* CLI tools
* Third-party integrations
* Automated publishing systems

Responsibilities:

* Wallet authentication
* Article publishing
* Article retrieval
* Verification requests

---

## Go REST API

The Go server is the central application layer.

Responsibilities:

* Authentication
* Authorization
* Validation
* Business logic
* Database operations
* IPFS interaction
* Blockchain interaction
* Verification
* Error handling

---

## Authentication Layer

Authentication uses **Sign-In with Ethereum (SIWE)**.

Workflow:

1. Client requests a nonce.
2. Server generates a nonce.
3. User signs the SIWE message using their wallet.
4. Server verifies the signature.
5. JWT token is issued.
6. Protected API endpoints accept the JWT.

Benefits:

* Passwordless authentication
* Wallet ownership verification
* No password storage
* Cryptographically secure identity

---

## SQLite Database

SQLite stores application metadata.

Stored information includes:

* Article metadata
* Titles
* Authors
* Publication dates
* IPFS CIDs
* Blockchain transaction hashes
* Article status
* Edit history

The database does **not** serve as the source of trust. It acts as the operational data store.

---

## IPFS

IPFS stores complete article content.

Benefits:

* Decentralized storage
* Reduced blockchain costs
* Content addressing
* High availability
* Efficient distribution

Each article receives a unique Content Identifier (CID).

---

## Ethereum Smart Contract

The smart contract stores immutable proofs.

Stored information includes:

* Article hash
* IPFS CID
* Author wallet
* Timestamp
* Version information
* Status

The contract enables anyone to independently verify article authenticity.

---

# Repository Structure

```text
newschain/

├── cmd/
│   └── server/
│       └── main.go

├── contracts/
│   └── NewsRegistry.sol

├── internal/
│   ├── api/
│   ├── auth/
│   ├── blockchain/
│   ├── models/
│   └── storage/

├── docs/
│   └── ARCHITECTURE.md

├── README.md
└── go.mod
```

---

# Publish Workflow

```text
Author

↓

Connect Wallet

↓

Request Nonce

↓

Sign SIWE Message

↓

Receive JWT

↓

Submit Article

↓

Generate Keccak-256 Hash

↓

Upload Content to IPFS

↓

Receive CID

↓

Save Metadata in SQLite

↓

Anchor Hash + CID on Ethereum

↓

Return Success
```

---

# Verification Workflow

```text
Reader

↓

Request Verification

↓

Load Article Metadata

↓

Retrieve Stored Article

↓

Recompute Keccak-256 Hash

↓

Read Blockchain Record

↓

Compare Hashes

↓

Verified ✓
```

If the hashes match, the article has not been altered since it was anchored on-chain.

---

# Update Workflow

```text
Authenticated Author

↓

Update Article

↓

Generate New Hash

↓

Upload Updated Content to IPFS

↓

Receive New CID

↓

Update SQLite

↓

Submit Amendment Transaction

↓

Emit ArticleAmended Event
```

Each update creates a new immutable blockchain record while preserving the article's history.

---

# Data Flow

## Authentication

```text
Wallet
   │
   ▼
SIWE
   │
   ▼
JWT
   │
   ▼
Protected API
```

---

## Publishing

```text
Article

↓

Hash

↓

IPFS

↓

SQLite

↓

Ethereum
```

---

## Verification

```text
SQLite

+

Ethereum

↓

Hash Comparison

↓

Verified
```

---

# Security Architecture

NewsChain uses multiple layers of security.

## Identity

* Wallet ownership verification
* SIWE authentication
* JWT authorization

---

## Integrity

* Keccak-256 hashing
* Immutable blockchain records
* Tamper detection

---

## Storage

* Hybrid architecture
* Off-chain content
* On-chain verification

---

## Authorization

Only authenticated article owners may:

* Edit articles
* Retract articles

Ownership is validated using the wallet address associated with the article.

---

# Design Decisions

## Why Hybrid Storage?

Storing full articles on-chain is expensive and inefficient.

Instead:

* Full content is stored in IPFS.
* Metadata is stored in SQLite.
* Integrity proof is stored on Ethereum.

This approach balances decentralization, cost, and performance.

---

## Why SIWE?

Wallet-based authentication:

* Eliminates passwords.
* Verifies identity cryptographically.
* Integrates naturally with Ethereum.

---

## Why SQLite?

SQLite provides:

* Simplicity
* Fast local development
* Minimal operational overhead

Future releases may add PostgreSQL support for production deployments.

---

# Scalability

Planned improvements include:

* PostgreSQL
* Redis caching
* Horizontal API scaling
* Kubernetes deployment
* Load balancing
* Event indexing with The Graph
* Multi-chain support
* Background job processing

---

# Reliability

Future production improvements:

* Automatic backups
* Health checks
* Retry mechanisms
* Circuit breakers
* Structured logging
* Metrics collection
* Distributed tracing

---

# Future Architecture

```text
                    Load Balancer
                           │
        ┌──────────────────┴──────────────────┐
        ▼                                     ▼
+----------------+                   +----------------+
| API Instance 1 |                   | API Instance 2 |
+-------+--------+                   +-------+--------+
        │                                    │
        └──────────────────┬─────────────────┘
                           ▼
                    PostgreSQL Cluster
                           │
                    Redis Cache Layer
                           │
                    Background Workers
                           │
          +----------------+----------------+
          ▼                                 ▼
      Ethereum                        IPFS Cluster
```

---

# Architecture Principles

NewsChain is built around the following principles:

* Security by default
* Cryptographic verification
* Separation of concerns
* Modular architecture
* Developer-friendly APIs
* Scalability
* Open-source collaboration
* Extensibility
* Maintainability
* Simplicity where practical

---

# Conclusion

NewsChain combines modern backend engineering with blockchain technology to create a decentralized publishing platform that emphasizes authenticity, transparency, and verifiable integrity.

The architecture is intentionally modular, allowing contributors to improve or replace individual components without fundamentally changing the overall system design.
