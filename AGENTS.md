# AGENTS.md - Dell Terraform Provider for OpenManage Enterprise

## Project Overview

This is the Terraform provider for Dell OpenManage Enterprise (OME) fleet management. It implements resources and data sources using HashiCorp's Terraform Plugin Framework, enabling infrastructure-as-code management of Dell server fleets through OME.

- **Language:** Go 1.25
- **Module path:** `terraform-provider-ome`
- **Terraform Plugin Framework:** v1.19.0
- **SDK:** Internal client (no external SDK dependency)
- **Registry address:** `registry.terraform.io/dell/ome`
- **License:** Mozilla Public License 2.0

## Architecture

The provider follows the standard Terraform Plugin Framework architecture. It runs as a gRPC server that Terraform Core communicates with to manage OME resources.

### Provider Configuration

The provider authenticates to an OME server using endpoint, username, and password. Configuration can be supplied via HCL provider block or environment variables (`OME_ENDPOINT`, `OME_USERNAME`, `OME_PASSWORD`, `OME_INSECURE`, `OME_TIMEOUT`).

### SDK Strategy

Uses an **internal client** — REST calls are implemented directly in the provider code under `clients/`. No external SDK dependency. This gives full control but requires more maintenance.

### Resources and Data Sources

The provider exposes approximately 14 resources and 10 data sources covering OME entities such as templates, firmware baselines, configuration baselines, device groups, deployments, firmware catalogs, and discovery jobs.

## Directory Structure

```
main.go                           Entry point (providerserver.Serve)
ome/
  provider.go                     Provider configuration, resource/datasource registration
  *_resource.go                   Resource implementations
  *_datasource.go                 Data source implementations
  *_test.go                       Unit and acceptance tests
clients/                          OME REST API client implementations
  json_data/                      JSON fixture data for tests
models/                           Terraform state model structs
helper/                           Shared helper functions
utils/                            Shared utility functions
examples/                         Example HCL configurations
docs/                             Generated documentation
templates/                        Documentation templates
testdata/                         Test fixture data
tools/                            Build and generation tools
about/                            Provider metadata
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Compile the provider binary |
| `make install` | Build and install to `~/.terraform.d/plugins/` |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests (`TF_ACC=1`, requires live hardware) |
| `make check` | Run `gofmt`, `golangci-lint`, `go vet` |
| `make gosec` | Run security scan with `gosec` |
| `make cover` | Generate HTML coverage report |
| `make generate` | Run `go generate` (docs generation) |

## Testing

### Unit Tests (mockey)

- Test files follow `*_test.go` convention in `ome/`.
- Frameworks: `github.com/stretchr/testify` (assertions), `github.com/bytedance/mockey` (function-level mocking).
- Run with `make test`.
- No hardware required.

### Acceptance Tests (terraform-plugin-testing)

- **Requires live OME hardware** with credentials set via environment variables.
- Creates real resources — clean up after failures.
- Run with `make testacc`.

### Running Tests

```bash
# Unit tests (no hardware)
make test

# Acceptance tests (requires live hardware)
export OME_ENDPOINT="https://ome-ip"
export OME_USERNAME="admin"
export OME_PASSWORD="secret"
export OME_INSECURE="true"
make testacc
```

## Code Style and Conventions

### Code Organization Patterns

- **Resource pattern:** Each resource is a Go file in `ome/` implementing `resource.Resource`.
- **Models:** Terraform state structs in `models/` using `tfsdk` struct tags.
- **Clients:** `clients/` contains OME REST API implementations with methods for each OME endpoint.
- **Helpers and utils:** Shared logic split between `helper/` and `utils/`.

### File Header

All source files must include the Dell copyright and MPL 2.0 license header.

## Common Development Tasks

### Adding a New Resource

1. Create `ome/<name>_resource.go` implementing `resource.Resource`.
2. Create the model struct in `models/`.
3. Add API client methods in `clients/`.
4. Register in `ome/provider.go`.
5. Add unit and acceptance tests.
6. Create example HCL in `examples/resources/ome_<name>/`.
7. Run `make generate` to produce documentation.

## CI/CD

GitHub Actions workflows in `.github/workflows/`. GoReleaser configuration in `.goreleaser.yml` builds cross-platform binaries.

## Code Ownership

All files are owned by the maintainers defined in `.github/CODEOWNERS`.
