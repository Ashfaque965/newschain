# Security Policy

Thank you for helping keep **NewsChain** and its users secure.

The security of our project is a top priority. We greatly appreciate responsible disclosure of security vulnerabilities and will make every effort to acknowledge, investigate, and address valid reports in a timely manner.

---

# Supported Versions

The following table indicates which versions of NewsChain currently receive security updates.

| Version               | Supported |
| --------------------- | --------- |
| Latest `main` branch  | ✅ Yes     |
| Latest stable release | ✅ Yes     |
| Older releases        | ❌ No      |

If you are using an unsupported version, please upgrade to the latest release before reporting an issue.

---

# Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub Issues or Discussions.**

Instead, report vulnerabilities privately using one of the following methods:

* **GitHub Security Advisories** (preferred): Use GitHub's **Report a Vulnerability** feature available under the repository's **Security** tab.
* **Email** (optional): Contact the project maintainers at **[security@example.com](mailto:security@example.com)** (replace with your actual security contact address).

Include as much information as possible:

* Description of the vulnerability
* Steps to reproduce
* Proof of concept (if available)
* Affected components
* Potential impact
* Suggested mitigation (optional)
* Environment details (OS, Go version, browser, wallet, etc.)

Providing complete information helps us investigate and resolve issues more quickly.

---

# What to Expect

After submitting a report, you can generally expect the following process:

1. We will acknowledge receipt of your report, typically within **72 hours**.
2. We will investigate the issue and assess its severity.
3. We may contact you for additional information if needed.
4. We will develop, test, and validate a fix.
5. We will publish a security update when appropriate.
6. We will publicly disclose the issue after a fix is available, unless there is a compelling reason to delay disclosure.

Response times may vary depending on the complexity and impact of the vulnerability.

---

# Responsible Disclosure

We ask that you:

* Give us a reasonable amount of time to investigate and fix the issue before publicly disclosing it.
* Avoid accessing, modifying, or deleting data that does not belong to you.
* Avoid actions that could disrupt the availability of the project or services.
* Act in good faith and comply with applicable laws.

We will not pursue action against researchers who follow this policy and conduct responsible security research.

---

# Scope

This policy applies to the NewsChain project, including:

* Go backend
* REST API
* Authentication (SIWE / JWT)
* Smart contracts
* Blockchain integration
* IPFS integration
* SQLite storage
* Configuration and deployment
* Official documentation and project infrastructure

Third-party services and dependencies should be reported to their respective maintainers when appropriate.

---

# Out of Scope

The following are generally outside the scope of this policy:

* Issues in unsupported versions
* Missing security headers in local development environments
* Social engineering attacks
* Denial-of-service testing without prior approval
* Vulnerabilities in third-party services or dependencies that are not specific to NewsChain
* Reports requiring unrealistic user interaction with no meaningful security impact
* Best-practice suggestions that do not demonstrate a security vulnerability

---

# Security Best Practices for Contributors

When contributing code, please:

* Never commit secrets, API keys, or private keys.
* Use environment variables for sensitive configuration.
* Validate all user input.
* Handle authentication and authorization carefully.
* Keep dependencies up to date.
* Write tests for security-sensitive functionality.
* Follow the principle of least privilege.
* Review changes for potential security implications before submitting a pull request.

---

# Security Features

NewsChain includes several security-focused design elements, including:

* Wallet-based authentication using **Sign-In with Ethereum (SIWE)**.
* JWT-based session management.
* Cryptographic integrity verification using **Keccak-256**.
* Immutable blockchain anchoring of article hashes.
* Hybrid off-chain/on-chain storage with IPFS and Ethereum.
* Author ownership validation for article updates.
* Immutable audit trail for article amendments and retractions.

---

# Security Roadmap

Future security enhancements may include:

* Multi-signature administrative controls
* Hardware Security Module (HSM) support
* Rate limiting and abuse prevention
* Role-based access control (RBAC)
* Audit logging
* OpenAPI request validation
* Dependency vulnerability scanning
* Static Application Security Testing (SAST)
* Dynamic Application Security Testing (DAST)
* Smart contract security audits
* Continuous security monitoring
* Supply chain security (SBOM generation)
* Container image scanning

---

# Hall of Fame

We appreciate members of the security community who help improve NewsChain through responsible disclosure.

With the reporter's permission, we may acknowledge valid vulnerability reports in this section after the issue has been resolved.

---

# License

By submitting a vulnerability report, you agree that the information you provide may be used to investigate, remediate, and document the issue.

---

Thank you for helping make **NewsChain** more secure for everyone.
