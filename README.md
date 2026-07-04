
# 📰 NewsChain

> **A decentralized, tamper-evident news publishing platform powered by Go, Ethereum, IPFS, SQLite, and Sign-In with Ethereum (SIWE).**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge\&logo=go)](https://go.dev)
[![Solidity](https://img.shields.io/badge/Solidity-0.8+-363636?style=for-the-badge\&logo=solidity)](https://soliditylang.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Open Source](https://img.shields.io/badge/Open%20Source-Welcome-blue?style=for-the-badge)](CONTRIBUTING.md)
[![PRs Welcome](https://img.shields.io/badge/PRs-Welcome-brightgreen?style=for-the-badge)](CONTRIBUTING.md)

---

## Overview

NewsChain is an open-source decentralized publishing platform that combines traditional web infrastructure with blockchain technology to provide **cryptographically verifiable news authenticity**.

Instead of storing entire articles on-chain (which is expensive), NewsChain stores:

* Full article content in **SQLite**
* Content backup on **IPFS**
* Immutable **Keccak-256 content hash** on Ethereum

Anyone can independently verify that an article has **never been modified** after publication.

---

# Why NewsChain?

Traditional news platforms rely on trust.

NewsChain enables **cryptographic trust**.

Every published article receives an immutable blockchain proof containing:

* Publish timestamp
* Author wallet
* Content hash
* IPFS CID
* Edit history

Readers can verify authenticity without trusting NewsChain itself.

---

# Features

### Authentication

* Sign-In with Ethereum (SIWE)
* Wallet-based authentication
* JWT session management
* No passwords

---

### Article Publishing

* Publish articles
* Edit articles
* Retract articles
* Ownership validation
* Immutable audit trail

---

### Blockchain Verification

* Keccak256 hashing
* Ethereum smart contract registry
* Tamper detection
* On-chain proof of publication

---

### Storage

Hybrid architecture:

* SQLite
* IPFS
* Ethereum

Fast and inexpensive while maintaining blockchain integrity.

---

### REST API

Production-ready REST API built with Go.

Endpoints include:

* Authentication
* Publishing
* Editing
* Verification
* Listing
* Moderation

---

### Security

* SIWE authentication
* JWT authorization
* Hash verification
* Wallet ownership validation
* Immutable blockchain records

---

# Architecture

```
                    Wallet (MetaMask)
                           │
                           ▼
                 SIWE Authentication
                           │
                           ▼
                 Go REST API Server
                           │
      ┌────────────────────┼─────────────────────┐
      ▼                    ▼                     ▼
 SQLite Database      Ethereum Network        IPFS
   Articles          Hash Registry         Article Backup
```

---

# Repository Structure

```
newschain/

├── cmd/
│   └── server/
│       └── main.go

├── contracts/
│   └── NewsRegistry.sol

├── internal/
│
├── api/
│   ├── handlers.go
│
├── auth/
│   └── auth.go
│
├── blockchain/
│   └── client.go
│
├── models/
│   └── article.go
│
├── storage/
│   ├── store.go
│   └── ipfs.go
│
├── README.md
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── SECURITY.md
├── LICENSE
└── go.mod
```

---

# Technology Stack

Backend

* Go
* SQLite

Blockchain

* Solidity
* Ethereum
* Base Sepolia
* go-ethereum

Storage

* SQLite
* IPFS

Authentication

* SIWE (EIP-4361)
* JWT

Cryptography

* Keccak256
* secp256k1
* ECDSA

---

# Quick Start

## Clone

```bash
git clone https://github.com/YOUR_USERNAME/newschain.git

cd newschain
```

---

## Install Dependencies

```bash
go mod tidy
```

---

## Deploy Smart Contract

Example using Foundry:

```bash
forge create contracts/NewsRegistry.sol:NewsRegistry \
--rpc-url https://sepolia.base.org \
--private-key $DEPLOYER_KEY
```

Save the deployed contract address.

---

## Configure Environment

```bash
export DB_PATH=newschain.db

export RPC_URL=https://sepolia.base.org

export CHAIN_ID=84532

export CONTRACT_ADDRESS=0x...

export SIGNER_PRIVATE_KEY=0x...

export IPFS_API_URL=http://127.0.0.1:5001

export PORT=8080
```

---

## Run

```bash
go run ./cmd/server
```

---

# API

## Authentication

### Get Nonce

```
GET /auth/nonce
```

### Verify Wallet

```
POST /auth/verify
```

Returns

```
JWT Token
```

---

## Articles

### List

```
GET /articles
```

### Get

```
GET /articles/{id}
```

### Verify

```
GET /articles/{id}/verify
```

### Publish

```
POST /articles
```

### Update

```
PUT /articles/{id}
```

### Retract

```
DELETE /articles/{id}
```

---

# Verification Flow

```
Author

↓

Creates Article

↓

Hash Generated

↓

Pinned to IPFS

↓

Hash Stored on Ethereum

↓

Article Saved in SQLite

↓

Reader Requests Verification

↓

Server Recomputes Hash

↓

Compares Against Blockchain

↓

Verified ✓
```

---

# Security Model

NewsChain protects against:

* Silent article modification
* Database tampering
* Unauthorized editing
* Fake publication dates
* Identity spoofing

Every article is cryptographically verifiable.

---

# Roadmap

## Phase 1

* REST API
* Ethereum integration
* SIWE
* IPFS
* SQLite

---

## Phase 2

* PostgreSQL support
* Redis caching
* Docker
* Docker Compose
* OpenAPI
* Swagger

---

## Phase 3

* Kubernetes deployment
* Prometheus metrics
* Grafana dashboards
* Structured logging
* Rate limiting

---

## Phase 4

* The Graph indexing
* ENS support
* Multi-chain support
* Arbitrum
* Polygon
* Optimism
* Base

---

## Phase 5

* Frontend (React)
* Mobile App
* WebSocket notifications
* AI-assisted moderation
* DAO governance

---

# Open Source

We welcome contributors from around the world.

Whether you're fixing a typo, improving documentation, adding tests, or implementing major features, every contribution matters.

Please read:

* CONTRIBUTING.md
* CODE_OF_CONDUCT.md
* SECURITY.md

before submitting a Pull Request.

---

# Good First Issues

Ideas suitable for first-time contributors:

* Improve API documentation
* Add unit tests
* Improve error handling
* Add Docker support
* Add GitHub Actions CI
* Add Swagger docs
* Add PostgreSQL backend
* Improve logging
* Add pagination
* Add search endpoint
* Add article categories
* Improve SIWE validation

---

# Development Workflow

```text
Fork Repository

↓

Create Feature Branch

↓

Write Code

↓

Add Tests

↓

Run Linter

↓

Commit

↓

Open Pull Request

↓

Code Review

↓

Merge
```

---

# Community

We believe trustworthy journalism should be transparent, verifiable, and decentralized.

Join us in building the future of authenticated digital publishing.

Contributions are welcome from:

* Go developers
* Solidity developers
* Security researchers
* Blockchain engineers
* DevOps engineers
* UI/UX designers
* Technical writers
* Students
* Open-source enthusiasts

---

# License

Licensed under the MIT License.

See the LICENSE file for details.

---

# Acknowledgements

Inspired by:

* Ethereum
* IPFS
* EIP-4361 (Sign-In with Ethereum)
* go-ethereum
* Foundry
* OpenZeppelin
* The Graph

---

# Star History

If you find this project useful:

⭐ Star the repository

🍴 Fork it

🐛 Report bugs

💡 Suggest new features

🤝 Submit pull requests







👋 About Me
Hi, I'm Ashfaque Quraishi, a passionate Full Stack & Blockchain Developer focused on building scalable web applications and decentralized systems.
I enjoy turning ideas into real-world products using modern web technologies, and I’m constantly exploring new tools in MERN stack, Web3, and smart contract development.

🚀 What I Do
🌐 Full Stack Web Development (MERN)
⛓️ Blockchain & Smart Contracts Development
🔐 Web3 Applications & DApps
📱 API Development & Integration
🧠 Learning System Design & Scalable Architectures

🛠️ Tech Stack
Frontend:
React.js
HTML, CSS, JavaScript
Tailwind CSS
Backend:
Node.js
Express.js
Database:
MongoDB
Blockchain:
Solidity
Ethereum
Web3.js / Ethers.js
Tools & Platforms:
Git & GitHub
Postman
VS Code

📊 GitHub Stats
🔭 Repositories: 74+
⭐ Stars: 87+
👥 Followers: 11
🔁 Following: 330

🌍 Connect With Me
GitHub: https://github.com/Ashfaque965
Instagram: https://www.instagram.com/ashfaquequraishi111/
Telegram: https://t.me/King_world1111
LinkedIn: https://in/ashfaque-quraishi-51b6a0222
X (Twitter): @king_quraishi1

💡 Goals
Build impactful Web3 & SaaS products
Contribute to open-source projects
Grow as a top-tier Full Stack + Blockchain engineer
Launch scalable tech solutions


Together we can build a more transparent and trustworthy news ecosystem.
