# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest  | ✅ Yes    |
| < 1.0   | ❌ No     |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

If you discover a security vulnerability in LiteLog, report it responsibly by opening a [GitHub Security Advisory](https://github.com/yashnaiduu/Litelog/security/advisories/new).

Include the following in your report:

- A description of the vulnerability and its potential impact
- Steps to reproduce the issue
- Any suggested mitigations or patches

You can expect an acknowledgement within **72 hours** and a status update within **7 days**.

## Scope

The following are in scope for security reports:

- Arbitrary code execution via the ingestion API or CLI
- Path traversal or unauthorized file access
- SQL injection in the query CLI
- Denial-of-service vulnerabilities in the HTTP server

## Out of Scope

- Vulnerabilities in third-party dependencies not yet patched upstream
- Issues requiring physical access to the machine running LiteLog
