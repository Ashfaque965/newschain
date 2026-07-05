
# Changelog

All notable changes to **NewsChain** will be documented in this file.

The format is based on **Keep a Changelog** and this project follows **Semantic Versioning (SemVer)**.

## Versioning

NewsChain uses the following version format:

```text
MAJOR.MINOR.PATCH
```

* **MAJOR**: Incompatible API or architecture changes.
* **MINOR**: New features added in a backward-compatible manner.
* **PATCH**: Backward-compatible bug fixes and security updates.

---

## [Unreleased]

### Added

* Placeholder for upcoming features.
* Placeholder for documentation improvements.

### Changed

* Placeholder for upcoming enhancements.

### Fixed

* Placeholder for bug fixes.

### Security

* Placeholder for security improvements.

---

## [1.0.0] - 2026-07-05

### 🎉 Initial Release

First public release of **NewsChain**.

### Added

#### Backend

* Go REST API server
* Article management endpoints
* SQLite persistence layer
* Environment-based configuration

#### Authentication

* Sign-In with Ethereum (SIWE)
* JWT-based authentication
* Wallet ownership verification
* Nonce generation and signature verification

#### Blockchain

* Ethereum smart contract registry
* Keccak-256 content hashing
* On-chain article verification
* Immutable publish records
* Article amendment events
* Article retraction events

#### Storage

* SQLite for metadata and article storage
* IPFS integration for decentralized content storage
* Mock IPFS implementation for local development

#### API

* Publish articles
* Retrieve articles
* List published articles
* Amend articles
* Retract articles
* Verify article integrity

#### Security

* Cryptographic integrity verification
* Author ownership validation
* JWT authorization
* Blockchain-backed audit trail

#### Documentation

* Comprehensive README
* CONTRIBUTING guide
* Code of Conduct
* Security Policy
* MIT License

### Changed

* Initial project structure established.

### Fixed

* Initial implementation.

### Security

* Implemented wallet-based authentication.
* Added immutable blockchain verification.
* Added tamper-evident article integrity checks.

---

## Future Releases

The following sections are placeholders for future versions.

### [1.1.0]

#### Added

* Docker support
* Docker Compose
* Swagger/OpenAPI documentation
* PostgreSQL backend
* Pagination
* Search API

---

### [1.2.0]

#### Added

* Redis caching
* GitHub Actions CI/CD
* Metrics endpoint
* Structured logging
* Request tracing

---

### [2.0.0]

#### Added

* React frontend
* Admin dashboard
* Role-Based Access Control (RBAC)
* Multi-chain support
* The Graph integration
* Kubernetes deployment

---

## Release Notes

Each release should include:

* New features
* Improvements
* Bug fixes
* Breaking changes
* Security updates
* Documentation updates
* Dependency updates

---

## Contributing

When submitting a Pull Request, please update this changelog if your contribution includes:

* New features
* Breaking changes
* Bug fixes
* Security improvements
* Performance enhancements
* Documentation changes that affect users

Follow the existing format and add your changes under the **Unreleased** section.

---

## References

* Keep a Changelog: https://keepachangelog.com/
* Semantic Versioning: https://semver.org/
