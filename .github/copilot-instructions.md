# Terraform Provider for Equinix Platform - Development Guide

## Repository Overview

This is the official Terraform provider for Equinix Platform, enabling lifecycle management of Equinix resources through Terraform. The codebase is approximately 64,000 lines of Go code across 789 files, implementing a Terraform plugin that supports both Equinix Metal and Fabric services.

**Key Technologies:**
- Language: Go 1.23.0+
- Framework: Terraform Plugin Framework v1.15.0 and Plugin SDK v2.37.0 (using mux for compatibility)
- Architecture: Mixed SDKv2 and Framework provider using terraform-plugin-mux
- Build Tool: GNU Make + Go toolchain
- Testing: Terraform Plugin Testing framework with acceptance tests

## Project Structure

### Core Directories
- `main.go` - Provider entry point, muxes SDKv2 and Framework providers
- `equinix/` - Legacy SDKv2 provider code (avoid adding new files here per validation workflow)
- `internal/` - Main implementation code
  - `internal/provider/` - Framework-based provider setup
  - `internal/resources/` - Resource implementations (metal, fabric subdirectories)
  - `internal/config/` - Provider configuration and client setup
  - `internal/fabric/`, `internal/network/` - Service-specific helpers
  - `internal/acceptance/` - Acceptance test helpers
  - `internal/sweep/` - Resource cleanup for tests
- `docs/` - Generated documentation (do not edit manually)
- `templates/` - Documentation templates for generation
- `examples/` - Example Terraform configurations for docs
- `scripts/` - Build and validation scripts

### Configuration Files
- `GNUmakefile` - Build automation (primary interface)
- `.golangci.yaml` - Linter configuration (golangci-lint v2.3.0)
- `go.mod` - Go dependencies
- `.github/workflows/` - CI/CD workflows

## Build, Test, and Validation

### Prerequisites
- Go 1.23.0+ (required version from go.mod)
- Terraform 1.0.11+ (for testing)
- GNU Make
- Git

### Essential Commands (Run in Order)

#### 1. Initial Setup
```bash
# Download dependencies (always run first after clone or when go.mod changes)
go mod download
```

#### 2. Build
```bash
# Clean and build the provider binary (takes ~30 seconds)
make clean build

# The binary is created as: ./terraform-provider-equinix
```

#### 3. Unit Tests
```bash
# Run unit tests (takes ~2-3 minutes, skips acceptance tests)
make test

# Run specific test:
# go test -v ./path/to/package -run TestName
```

#### 4. Code Formatting
```bash
# Check formatting (always run before committing)
make fmtcheck

# Auto-format code if needed
make fmt
```

#### 5. Linting
```bash
# Run linter (takes ~3-5 minutes)
# Only checks files changed from origin/main
make lint

# The linter uses golangci-lint v2.3.0 specified in GNUmakefile
```

#### 6. Documentation Generation
```bash
# Generate docs from templates (requires network access)
make docs

# Verify docs are up to date
make docs-check
```

**NOTE:** `make docs-check` may fail in offline/restricted network environments due to Terraform CLI download. This is expected in sandboxed environments.

### Complete Pre-Commit Checklist
```bash
# Run these commands in order before committing changes:
make fmt          # Auto-format code
make fmtcheck     # Verify formatting
make build        # Ensure code compiles
make test         # Run unit tests
make lint         # Check code quality
make docs         # Update documentation (if resources/datasources changed)
```

### Acceptance Tests
Acceptance tests create real infrastructure and require credentials:
```bash
export EQUINIX_API_ENDPOINT=https://api.equinix.com
export EQUINIX_API_CLIENTID=<your-client-id>
export EQUINIX_API_CLIENTSECRET=<your-client-secret>
make testacc  # Takes 180+ minutes, runs all acceptance tests

# Run specific acceptance test:
TF_ACC=1 go test -v -timeout=20m ./... -run=TestAccMetalDevice_Basic
```

**WARNING:** Acceptance tests are expensive and create real resources. Only run locally if necessary.

## GitHub Actions CI/CD

### Pull Request Validation Workflows

All PRs trigger these checks automatically:

1. **validate_pr.yml** - PR Title & File Location Checks
   - Validates PR title follows Conventional Commits (feat:, fix:, docs:, etc.)
   - Blocks new files in `equinix/` package (use `internal/` instead)
   
2. **test.yml** - Unit Tests & Docs Validation
   - Runs `go build` and `go test` with coverage
   - Runs `make docs-check` to ensure docs are generated
   
3. **golangci-lint.yml** - Code Quality
   - Runs golangci-lint with `--whole-files` flag
   - Only checks changed files (`--new-from-rev=origin/main`)
   
4. **fabric_acctest.yml** & **metal_acctest.yml** - Acceptance Tests
   - Run only when specific paths change (fabric/** or metal/**)
   - Require approval for external contributors

### Common CI Failures & Fixes

**Problem:** "PR title doesn't follow Conventional Commits"
- **Fix:** Start PR title with: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`, `chore:`, etc.

**Problem:** "New files added to equinix package"
- **Fix:** Move new .go files to appropriate `internal/` subdirectory (e.g., `internal/resources/metal/` or `internal/resources/fabric/`)

**Problem:** "Uncommitted doc changes"
- **Fix:** Run `make docs` and commit the generated files in `docs/`

**Problem:** Linter failures
- **Fix:** Run `make lint` locally to see issues, then `make fmt` to auto-fix formatting

## Development Workflow

### Making Code Changes

1. **For New Resources/Data Sources:**
   - Add implementation to `internal/resources/{metal|fabric}/resource_name/`
   - Add examples to `examples/resources/equinix_resource_name/`
   - Add template (if needed) to `templates/resources/`
   - Run `make docs` to generate documentation
   - Add tests following existing patterns in `*_test.go` files

2. **For Bug Fixes:**
   - Locate code in `internal/` or `equinix/` (prefer `internal/`)
   - Make minimal changes to fix issue
   - Add/update tests to cover the fix
   - Run full validation checklist

3. **Coding Standards:**
   - Use existing code style (gofmt enforced)
   - Add comments for exported functions/types
   - Handle all errors (errcheck enforced)
   - Follow Terraform plugin development best practices
   - Use descriptive variable names (revive linter enforced)

### Testing Strategy

- **Unit Tests:** Test business logic, converters, validators
- **Acceptance Tests:** Test actual provider behavior with Terraform
  - Named: `TestAcc<Resource>_<scenario>`
  - Use `resource.Test()` from terraform-plugin-testing
  - Located in `*_test.go` files alongside resource code

### Documentation Requirements

All resource/data source changes MUST include documentation:
- Add/update examples in `examples/`
- Documentation is auto-generated via `make docs`
- Templates in `templates/` control doc structure
- Default templates: `templates/resources.md.tmpl`, `templates/data-sources.md.tmpl`
- CI enforces docs are up-to-date via `make docs-check`

## Common Pitfalls & Workarounds

### Known Issues

1. **Docs generation requires network access** - May fail in restricted environments when downloading Terraform CLI binary. Skip if sandboxed.

2. **Equinix package restrictions** - Cannot add new files to `equinix/*.go`. Use `internal/` packages instead. This is enforced by PR validation.

3. **Multiple provider versions** - The codebase uses both SDKv2 (legacy) and Framework (new). New resources should use Framework in `internal/provider/`.

4. **Test timeouts** - Acceptance tests can be very slow. Default timeout is 180 minutes (`ACCTEST_TIMEOUT=180m`).

5. **API rate limiting** - Acceptance tests may hit API rate limits. Use `ACCTEST_PARALLELISM=8` (default) or lower.

### Build Environment Notes

- **Go Module Mode:** This project uses Go modules; safe to work outside GOPATH
- **Binary Output:** Build creates `./terraform-provider-equinix` (~60MB binary)
- **Clean Builds:** Run `make clean` to remove binary and build artifacts
- **Vendor Directory:** Not used; dependencies from go.mod only

### Debugging Tips

1. **Enable Terraform Debug Logs:**
   ```bash
   TF_LOG=DEBUG terraform apply
   ```

2. **Test Specific Provider Code:**
   ```bash
   TF_ACC=1 TF_LOG=DEBUG go test -v ./internal/resources/metal/device -run=TestAccMetalDevice_Basic
   ```

3. **Check Provider Registration:**
   - Both providers registered in `main.go` via muxServer
   - Framework provider in `internal/provider/`
   - SDKv2 provider in `equinix/`

## Quick Reference

### Makefile Targets
- `make build` - Build provider binary
- `make test` - Run unit tests (10 min timeout)
- `make testacc` - Run acceptance tests (180 min timeout, needs credentials)
- `make lint` - Run golangci-lint
- `make fmt` - Auto-format Go code
- `make fmtcheck` - Check code formatting
- `make docs` - Generate documentation
- `make docs-check` - Verify docs are current
- `make clean` - Remove build artifacts
- `make sweep` - Clean up leaked test resources (needs credentials)

### Environment Variables for Testing
- `TF_ACC=1` - Enable acceptance tests
- `TF_LOG=DEBUG` - Enable debug logging
- `EQUINIX_API_ENDPOINT` - API endpoint (default: https://api.equinix.com)
- `EQUINIX_API_CLIENTID` - API client ID
- `EQUINIX_API_CLIENTSECRET` - API client secret
- Various `TF_ACC_*` variables for test parametrization (see DEVELOPMENT.md)

### File Naming Conventions
- Resources: `resource.go`, `resource_test.go`, `resource_schema.go`
- Data Sources: `datasource.go`, `datasource_test.go`, `datasource_schema.go`
- Models: `models.go` - Terraform schema models
- Sweepers: `sweeper.go` - Resource cleanup

## Important: Trust These Instructions

These instructions have been validated by testing actual commands in the repository. Follow them precisely to avoid wasted time exploring or debugging failed commands. If you encounter issues not covered here, check:
1. DEVELOPMENT.md for detailed developer docs
2. CONTRIBUTING.md for contribution process
3. GitHub workflow files in `.github/workflows/` for CI requirements
4. GNUmakefile for available build targets

Only search or explore the codebase if these instructions are incomplete or if you've confirmed they're incorrect.
