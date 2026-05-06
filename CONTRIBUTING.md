# Contributing

## Development Flow
1. Fork/branch from the default branch.
2. Keep changes scoped and atomic.
3. Update docs and tests with code changes.
4. Run local validation before opening a PR.

## Validation Baseline
Backend-oriented projects:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
go generate ./...
go test ./...
```

Frontend-oriented projects:

```bash
npm install
npm run lint
npm run build
npm run test -- --runInBand
```

## Commit and PR Guidelines
- Use clear commit messages describing intent and impact.
- Reference changed modules/files in PR descriptions.
- Mention any migration, infra, or env changes explicitly.

## Secrets and Privacy
- Never commit real secrets.
- Keep `.env` files local and commit only `.env.example` templates.
- Redact sensitive data from logs, screenshots, and issue reports.
