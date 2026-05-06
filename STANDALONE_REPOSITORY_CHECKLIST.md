# Standalone Repository Checklist

Use this checklist before pushing this project as its own repository.

## 1. Repository Metadata
- [ ] `README.md` is current and includes setup/run/validation commands
- [ ] `LICENSE` exists
- [ ] `CONTRIBUTING.md`, `SECURITY.md`, and `CODE_OF_CONDUCT.md` exist

## 2. Configuration Hygiene
- [ ] `.env.example` templates exist for required runtime/config surfaces
- [ ] real `.env` and secret files are ignored by `.gitignore`
- [ ] no private keys or live credentials in tracked files

Secret scan example:

```bash
rg -n --hidden --glob '!**/.git/**' --glob '!**/node_modules/**' --glob '!**/.next/**' '(AKIA[0-9A-Z]{16}|AIza[0-9A-Za-z\-_]{35}|ghp_[0-9A-Za-z]{36}|xox[baprs]-[0-9A-Za-z-]+|-----BEGIN [A-Z ]+PRIVATE KEY-----)' .
```

## 3. Runtime Validation
- [ ] project setup instructions execute as documented
- [ ] core tests and lint/build checks pass for changed modules
- [ ] Docker profiles/services start with documented commands

## 4. Push Preparation
- [ ] review `git status` for unintended files
- [ ] split unrelated changes into separate commits
- [ ] tag release baseline after first stable publish
