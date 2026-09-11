.PHONY: check-version bump-patch bump-minor bump-major bump-beta test build-all

check-version:
	node scripts/bump-version.mjs --check

bump-patch:
	node scripts/bump-version.mjs patch

bump-minor:
	node scripts/bump-version.mjs minor

bump-major:
	node scripts/bump-version.mjs major

bump-beta:
	node scripts/bump-version.mjs beta

build-backend:
	cd registro-backend && go build -v -ldflags "-X registro-backend/pkg/version.GitCommit=$$(git rev-parse --short HEAD 2>/dev/null || echo dev) -X registro-backend/pkg/version.BuildDate=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o bin/api-server cmd/api-server/main.go

build-frontend:
	cd registro-frontend && npm run build
