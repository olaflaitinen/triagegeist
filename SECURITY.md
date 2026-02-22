# Security policy

This document describes how security vulnerabilities are handled for triagegeist and how to report them.

---

## Supported versions

Security fixes are applied to the **current major version** of the project. We do not maintain separate branches for older major versions unless explicitly stated in a release.

| Version | Supported          | Notes |
|---------|--------------------|-------|
| Latest major (e.g. v1.x) | Yes   | All security fixes and patches |
| Older majors       | No    | Upgrade to a supported version |

When a new major version is released, support for the previous major version will be announced (e.g. in release notes or this file) and may be time-limited.

---

## Reporting a vulnerability

**Do not report security vulnerabilities in public issues or pull requests.** Use one of the following private channels.

### Option 1: GitHub Security Advisories (preferred)

1. Go to the repository: [github.com/olaflaitinen/triagegeist](https://github.com/olaflaitinen/triagegeist)
2. Open the **Security** tab.
3. Click **Advisories**, then **Report a vulnerability** (or use [this link](https://github.com/olaflaitinen/triagegeist/security/advisories/new) if available).
4. Fill in:
   - **Title**: Short description of the issue
   - **Description**: Steps to reproduce, affected versions, impact
   - **Severity**: Use CVSS or "Low / Medium / High / Critical" if applicable

You will receive an initial response and can discuss the fix and disclosure timeline in the advisory thread.

### Option 2: Contact maintainers

If you cannot use GitHub Security Advisories, contact the maintainers listed in [GOVERNANCE.md](GOVERNANCE.md) or in the repository profile. Prefer encrypted communication if available. Include:

- Type of vulnerability (e.g. denial of service, information disclosure, input validation)
- Affected component and version
- Steps to reproduce
- Suggested fix (optional)
- Whether you want to be credited in the advisory (and how)

---

## What to expect

| Stage | Action |
|-------|--------|
| Acknowledgement | We will confirm receipt of the report and assign someone to triage it |
| Triage | We will assess severity and impact and decide on fix and release plan |
| Fix | A patch will be developed and tested; we may coordinate with you |
| Release | A new version (patch or minor, as appropriate) will be released |
| Disclosure | After the fix is available, we will publish a security advisory (e.g. GitHub Security Advisories) with details and credit, unless you prefer to remain anonymous |

We aim to acknowledge reports within a few business days and to provide a fix within a reasonable timeframe depending on severity. Critical issues will be prioritised.

---

## Scope

**In scope:**

- Bugs in triagegeist code that could lead to security impact (e.g. panic from malicious input, resource exhaustion, incorrect scoring due to logic errors that could be exploited).
- Vulnerabilities in dependencies used by the project (report so we can bump or replace them).
- Any issue that could lead to denial of service, information disclosure, or integrity violation when the library is used as intended.

**Out of scope:**

- General security hardening of systems that *use* triagegeist (e.g. network, OS, deployment). Those remain the responsibility of the integrator.
- Issues in third-party code that triagegeist does not depend on.

---

## Disclosure policy

- We do not disclose details of unfixed vulnerabilities publicly until a fix is available or a coordinated disclosure date is agreed.
- After a fix is released, we publish an advisory with sufficient detail for users to assess impact and upgrade. We will credit the reporter unless they prefer anonymity.
- We do not pursue legal action against researchers who report in good faith and follow this policy.

---

## Security-related configuration

triagegeist is a library; it does not open network ports or read secrets by default. When integrating it:

- Validate and sanitise inputs (vitals, resource counts) in the application layer.
- Do not expose internal triage parameters or patient data in logs or errors without proper controls.
- Follow your organisation's policies for handling health-related data (e.g. GDPR, HIPAA) and use triagegeist in a compliant pipeline.
