# Governance

This document describes how the triagegeist project is maintained, who has decision-making authority, and how disputes or major changes are handled. The project follows lightweight, transparent governance consistent with Google-style open source projects.

---

## Roles

### Maintainers

Maintainers have write access to the repository and are responsible for:

- Reviewing and merging pull requests
- Releasing versions and maintaining the changelog
- Enforcing code of conduct and security policy
- Setting technical and roadmap direction within the project scope

| Role        | Responsibility |
|-------------|----------------|
| **Merge**   | Approve and merge PRs; ensure CI and quality standards are met |
| **Release** | Tag releases, update CHANGELOG, publish release notes |
| **Steward** | Code of conduct and conflict resolution; final say on conduct and scope |

Current maintainers are listed in the repository (e.g. in README or GitHub "People" / "CODEOWNERS"). The initial authors are: **Gustav Olaf Yunus Laitinen-Fredriksson Lundström-Imanov**.

### Contributors

Anyone who submits a patch, issue, or documentation improvement is a contributor. Contributors do not have write access unless they are also maintainers. All contributions are subject to the [LICENSE](LICENSE) (EUPL-1.2) and [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

---

## Decision-making

| Type of decision | Process |
|------------------|---------|
| **Code and API** | Propose via issue or PR; maintainers review and merge. Breaking changes require explicit discussion and, where appropriate, a major version or deprecation path. |
| **Dependencies** | New dependencies must be justified (e.g. necessity, licence compatibility, maintenance). Prefer standard library or minimal, well-maintained packages. |
| **Roadmap and scope** | Maintainers set scope and priorities. Large or controversial changes should be discussed in an issue before a large PR. |
| **Conduct and membership** | Handled per [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md). The steward (or designated maintainer) decides on warnings and bans. |

There is no formal voting process; the project relies on maintainer consensus and open discussion. Disagreements are resolved through discussion; if consensus cannot be reached, the steward or repository owner may make a final decision and document the rationale.

---

## Adding or removing maintainers

- **Adding**: Existing maintainers may invite a new maintainer based on sustained, high-quality contributions and alignment with project values. The invitation is made in private; acceptance is documented (e.g. in GOVERNANCE or CODEOWNERS).
- **Removing**: A maintainer may step down at any time. Removal for cause (e.g. violation of code of conduct, prolonged inactivity) is decided by the remaining maintainers and the repository owner.

---

## Repository and licence

- The canonical repository is **github.com/olaflaitinen/triagegeist** (or as updated in the README).
- The project is licensed under the **European Union Public Licence v. 1.2 (EUPL-1.2)**. All contributions must be compatible with this licence; there is no CLA beyond the licence grant implied by submitting a PR.

---

## Transparency

- **Discussions**: Technical and scope discussions take place in the open (issues, PRs) unless they involve private security or conduct matters.
- **Releases**: Version tags, CHANGELOG, and release notes are public. Security fixes are disclosed in line with [SECURITY.md](SECURITY.md).
- **Governance changes**: Changes to this document are proposed via pull request and merged by maintainers. Significant changes (e.g. new roles, new licence) should be announced in a release or issue.
