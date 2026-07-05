

# Contributing to NewsChain

First off, thank you for your interest in contributing to **NewsChain**! 🎉

We appreciate every contribution, whether it's fixing a typo, improving documentation, reporting bugs, writing tests, or implementing new features.

Our goal is to build a secure, decentralized, and transparent news publishing platform together with the open-source community.

---

# Table of Contents

* Welcome
* Code of Conduct
* Ways to Contribute
* Reporting Bugs
* Suggesting Features
* Development Setup
* Branch Naming
* Commit Message Guidelines
* Pull Request Process
* Coding Standards
* Testing
* Documentation
* Issue Labels
* Good First Issues
* Community

---

# Welcome

NewsChain is an open-source project focused on:

* Go backend development
* Ethereum smart contracts
* IPFS integration
* Decentralized authentication
* Cryptographic verification
* REST APIs
* Security-first architecture

Whether you're a beginner or an experienced developer, your contributions are welcome.

---

# Code of Conduct

Please read our **CODE_OF_CONDUCT.md** before participating.

We are committed to providing a welcoming, respectful, and inclusive environment for everyone.

---

# Ways to Contribute

You can contribute in many ways, including:

## 🐛 Bug Reports

Help us identify bugs by opening GitHub Issues with:

* Clear description
* Expected behavior
* Actual behavior
* Steps to reproduce
* Environment details
* Screenshots (if applicable)

---

## ✨ Feature Requests

We welcome new ideas.

When suggesting a feature, include:

* Problem statement
* Proposed solution
* Possible implementation
* Alternatives considered

---

## 📚 Documentation

Documentation improvements are always appreciated.

Examples include:

* README improvements
* API documentation
* Architecture diagrams
* Tutorials
* Examples
* Code comments

---

## 🧪 Testing

Help improve reliability by adding:

* Unit tests
* Integration tests
* Edge case tests
* Security tests

---

## 🔐 Security Improvements

Security contributions are highly valued.

Examples include:

* Authentication improvements
* Input validation
* Cryptography reviews
* Dependency updates
* Vulnerability fixes

Please report sensitive vulnerabilities privately as described in **SECURITY.md** instead of opening a public issue.

---

# Development Setup

## 1. Fork the Repository

Fork NewsChain to your own GitHub account.

---

## 2. Clone Your Fork

```bash
git clone https://github.com/<your-username>/newschain.git

cd newschain
```

---

## 3. Add the Upstream Remote

```bash
git remote add upstream https://github.com/<original-owner>/newschain.git
```

---

## 4. Install Dependencies

```bash
go mod tidy
```

---

## 5. Run the Project

```bash
go run ./cmd/server
```

---

# Branch Naming

Please create a new branch for every change.

Examples:

```
feature/add-search-api

feature/ipfs-cache

bugfix/auth-token-expiry

bugfix/sqlite-migration

docs/update-readme

docs/api-reference

refactor/storage-layer

test/article-verification

ci/github-actions
```

Avoid committing directly to the `main` branch.

---

# Commit Message Guidelines

Use clear and descriptive commit messages.

Examples:

```
feat: add article search endpoint

feat: integrate IPFS pinning

fix: prevent duplicate article hashes

fix: validate JWT expiration

docs: improve API documentation

test: add blockchain verification tests

refactor: simplify authentication middleware

ci: add GitHub Actions workflow
```

Prefer the Conventional Commits format when possible.

---

# Pull Request Process

Before opening a Pull Request:

* Ensure your branch is up to date with `main`.
* Run tests locally.
* Run formatting tools.
* Update documentation if your change affects behavior.
* Keep Pull Requests focused on a single topic whenever possible.

Your Pull Request description should include:

* Summary of the change
* Related issue (if any)
* Testing performed
* Screenshots (if UI changes)
* Breaking changes (if any)

Maintainers may request changes before merging.

---

# Coding Standards

## Go

* Follow standard Go formatting (`gofmt`).
* Keep functions small and focused.
* Prefer clear names over abbreviations.
* Handle errors explicitly.
* Avoid unnecessary global state.

---

## Solidity

* Follow Solidity style guidelines.
* Document public and external functions.
* Validate inputs.
* Emit events for important state changes.
* Minimize gas usage where practical.

---

## General

* Write readable code.
* Avoid duplicated logic.
* Add comments only when they improve understanding.
* Keep dependencies minimal.

---

# Testing

Please include tests whenever possible.

Run all tests before submitting:

```bash
go test ./...
```

If you add smart contract tests, include instructions for running them.

---

# Documentation

If your contribution changes behavior or adds a feature, update the relevant documentation.

This may include:

* README.md
* API documentation
* Examples
* Comments
* Architecture documentation

Good documentation helps everyone.

---

# Issue Labels

Common labels include:

* `good first issue`
* `help wanted`
* `bug`
* `enhancement`
* `documentation`
* `security`
* `performance`
* `question`
* `discussion`

If you're new to open source, look for issues labeled **good first issue**.

---

# Good First Issues

Some beginner-friendly contribution ideas:

* Improve README examples
* Fix typos
* Add unit tests
* Improve logging
* Add Docker support
* Add pagination
* Improve API validation
* Add request tracing
* Improve error messages
* Add OpenAPI/Swagger documentation
* Improve CI workflows

---

# Community

We strive to build an open, collaborative community.

Please:

* Be respectful.
* Assume good intent.
* Provide constructive feedback.
* Welcome newcomers.
* Help others learn.

Healthy discussions lead to better software.

---

# Recognition

Every merged contribution is appreciated.

Thank you for helping make NewsChain more secure, reliable, and useful for the community.

Happy coding! 🚀
