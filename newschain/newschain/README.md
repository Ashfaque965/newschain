# NewsChain

A decentralized news platform: articles are written and served through a
normal Go REST API, but every publish/edit anchors a `keccak256` hash of
the content on-chain via a Solidity registry contract. Readers (or your
frontend) can hit `/articles/:id/verify` at any time to recompute the
hash from what's stored and confirm it still matches the on-chain record —
giving tamper-evident proof of authorship and publish time without
putting full article text on-chain (that lives in SQLite + IPFS).

## Why this design (portfolio framing)

- **Security-first**: content integrity is provable, not just claimed.
  Anyone can independently verify a story hasn't been silently altered
  after publication.
- **Wallet-based auth (SIWE/EIP-4361)**: no passwords — journalists sign a
  nonce with their wallet, server verifies the ECDSA signature and issues
  a JWT session.
- **Hybrid on-chain/off-chain storage**: full text lives in IPFS/SQLite
  (cheap, fast); only a 32-byte hash + CID pointer is anchored on-chain
  (gas-efficient). This is the same pattern used by real content-provenance
  systems (e.g. C2PA-adjacent designs).
- **Amend / retract with audit trail**: edits emit `ArticleAmended`,
  retractions emit `ArticleRetracted` — a public, immutable edit history.

## Architecture

```
newschain/
├── contracts/
│   └── NewsRegistry.sol       # on-chain hash anchoring, verify(), moderation
├── cmd/server/main.go         # entrypoint, env config, wiring
├── internal/
│   ├── api/handlers.go        # REST endpoints
│   ├── auth/auth.go           # SIWE signature verification + JWT sessions
│   ├── blockchain/client.go   # go-ethereum client for the registry contract
│   ├── models/article.go      # data types
│   └── storage/
│       ├── store.go           # SQLite persistence
│       └── ipfs.go            # IPFS pinning (HTTP API + mock for dev)
└── go.mod
```

## Setup

### 1. Deploy the contract

Compile `contracts/NewsRegistry.sol` with Foundry or Hardhat and deploy to
a testnet first (Base Sepolia is a good default — cheap and fast):

```bash
forge create contracts/NewsRegistry.sol:NewsRegistry \
  --rpc-url https://sepolia.base.org \
  --private-key $DEPLOYER_KEY
```

Note the deployed address for `CONTRACT_ADDRESS` below.

### 2. Configure environment

```bash
export DB_PATH=newschain.db
export RPC_URL=https://sepolia.base.org
export CHAIN_ID=84532
export CONTRACT_ADDRESS=0xYourDeployedAddress
export SIGNER_PRIVATE_KEY=0xYourServerSignerKey   # pays gas for publishes
export IPFS_API_URL=http://127.0.0.1:5001         # omit to use MockIPFSPinner locally
export PORT=8080
```

### 3. Run

```bash
go mod tidy
go run ./cmd/server
```

## API

| Method | Path                     | Auth | Description                                  |
|--------|--------------------------|------|-----------------------------------------------|
| GET    | `/auth/nonce?address=0x` | –    | Get a SIWE nonce/message to sign               |
| POST   | `/auth/verify`           | –    | Submit signed message, receive JWT             |
| GET    | `/articles`              | –    | List published articles                       |
| GET    | `/articles/:id`          | –    | Get one article                                |
| GET    | `/articles/:id/verify`   | –    | Recompute hash and check against chain         |
| POST   | `/articles`              | JWT  | Publish new article (pins IPFS + anchors hash) |
| PUT    | `/articles/:id`          | JWT  | Amend article (author only)                    |
| DELETE | `/articles/:id`          | JWT  | Retract article (author only)                  |

### Publish flow example

```bash
# 1. get nonce
curl "http://localhost:8080/auth/nonce?address=0xAbc..."

# 2. sign the returned `message` client-side with the wallet, then:
curl -X POST http://localhost:8080/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"address":"0xAbc...","message":"...","signature":"0x..."}'
# -> {"token": "..."}

# 3. publish
curl -X POST http://localhost:8080/articles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"L2 Fees Drop 40% After Upgrade","body":"..."}'
```

## Next steps / roadmap ideas

- Swap the hand-written ABI in `blockchain/client.go` for `abigen`-generated
  bindings once the contract is finalized.
- Add a subgraph (The Graph) to index `ArticlePublished`/`ArticleAmended`
  events for fast historical queries instead of hitting the DB only.
- Rate-limit `/articles` POST per wallet to prevent spam publishing.
- Add moderator dispute UI backed by `disputeArticle`.
- Move `jwtSecret` in `auth.go` to env-based config before any deployment.
