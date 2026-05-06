# Security Policy

## Supported Versions
Security fixes are applied to the default branch first. Release backports are best-effort.

## Reporting a Vulnerability
Do not open public issues for undisclosed vulnerabilities.

Use one of these methods:
- repository security advisories (preferred)
- private maintainer contact channel listed in repository settings

Include:
- affected component/path
- reproduction steps or proof of concept
- impact assessment
- mitigation ideas (if available)

## Secret Handling Rules
- Never commit live credentials, API keys, tokens, private keys, or production endpoints with privileged access.
- Use `.env.example` templates in git and keep real `.env` values outside source control.
- Rotate exposed credentials immediately if leakage is suspected.
