# KNOWLEDGE.md — terraform-provider-ome

<!-- yaml-metadata-start -->
scope_paths: ["./"]
capture_git_sha: "aa83a72e8e01871bd84c97a8d37caf9f322727b4"
status: "current"
auto_update: false
preview_before_apply: true
scaffold_version: "1.0"
# session_state: { is_complete: true }
<!-- yaml-metadata-end -->

<!-- quick-reference-start -->
## Agent Quick Reference

| Section | Heading | Summary | never_again_count |
|---------|---------|---------|-------------------|
| Component Overview | `## Component Overview` | Dell OpenManage Enterprise fleet management provider | — |
| Architectural Rationale | `## Architectural Rationale` | Internal SDK strategy; Plugin Framework architecture | — |
| Failure Modes & Gotchas | `## Failure Modes & Gotchas` | Endpoint format, SDK versioning, state secrets | 0 |
| Implicit Contracts | `## Implicit Contracts` | Env var precedence, auth validation, TLS defaults | — |
<!-- quick-reference-end -->

## Five Questions Quick Reference

### What does it do?
Terraform provider for Dell OpenManage Enterprise fleet management. Exposes 14 resources and 11 data sources covering templates, baselines, firmware repositories, configuration compliance, devices, discovery, user management, VLAN profiles, and deployment operations
through HashiCorp's Terraform Plugin Framework. Communicates with
the hardware REST API via Internal client (no external SDK).

### How do you modify it?
Create `resource_<name>.go` (or `*_resource.go`) implementing
`resource.Resource`, add model structs, register in `provider.go`,
add unit tests with mockey mocks, add acceptance tests, create
example HCL, and run `make generate` for docs.

### What breaks?
**Endpoint is the OpenManage Enterprise console IP or FQDN.** Requires OME to be installed and running. Acceptance tests against live hardware create real
resources — failed test runs may leave orphaned resources. State files
contain secrets — use encrypted remote backends.

### What depends on it?
Terraform Core (gRPC go-plugin), Internal client (no external SDK),
`hashicorp/terraform-plugin-framework` v1.19.0.

### What's undocumented?
The `clients/` package implements the OME REST API client directly — no external SDK. Session management and authentication are handled internally.

---

## Component Overview

Terraform provider for Dell OpenManage Enterprise fleet management.
14 resources and 11 data sources covering templates, baselines, firmware repositories, configuration compliance, devices, discovery, user management, VLAN profiles, and deployment operations. Resources use `resource_*.go` naming under `ome/`. Data sources use `datasource_*.go` naming.

---

## Architectural Rationale

The provider follows the standard Terraform Plugin Framework architecture
— a standalone Go binary communicating with Terraform Core over gRPC.

**SDK strategy (Internal):** No external SDK dependency. REST calls implemented in provider code via `clients/` package. The internal client handles OME REST API session management, authentication, and API operations.

All providers in the Dell Terraform family share this architecture:
Terraform Plugin Framework interfaces, `resource.Resource` for CRUD
resources, `datasource.DataSource` for read-only queries, models with
`tfsdk` struct tags, and mockey-based unit testing.

### Evolution

TBD — requires SME input on how the architecture changed over time.

---

## Failure Modes & Gotchas

### 1. Endpoint URL format

**Endpoint is the OpenManage Enterprise console IP or FQDN.** Requires OME to be installed and running.

### 2. Sensitive attributes must be marked

All credential fields must have `Sensitive: true` in the schema.
Without this, passwords appear in `terraform plan` output and state
files. This is enforced by code convention, not by the framework.

### 3. State file contains secrets

Terraform state files contain full resource representations including
credentials. Always use encrypted remote backends (S3+KMS, Terraform
Cloud) in production.


### Never Again

No incident-derived constraints recorded. If you know of past
incidents affecting this component, please record them during the
next Knowledge Extraction session.

### Evolution

TBD — requires SME input.

---

## Performance Characteristics

TBD — requires SME input for bottlenecks, scaling limits, tuning
parameters, benchmarks, and known performance cliffs.

### Evolution

TBD — requires SME input.

---

## Implicit Contracts

**Environment variable precedence:** env vars (`OME_*`)
override HCL provider block values when both are set. This is
implemented in `Configure()` and is not documented as an explicit
contract.

**Authentication validation:** `Configure()` makes a dummy API call
to validate credentials before any resource operations proceed. If
this call fails, all resource operations are blocked.

**TLS verification default:** `insecure` defaults to `false` —
TLS verification is on by default. Setting `insecure = true` is
a lab-only setting and must never be used in production.

**Acceptance test gating:** tests guarded by `TF_ACC=1` — never
run without live hardware credentials. Tests create real resources
that must be cleaned up manually if the test run fails.

### Evolution

TBD — requires SME input.

---

## Threading & Synchronization

Terraform Plugin Framework handles concurrency at the provider level.
Individual resource operations are not concurrent by default.

### Evolution

TBD — requires SME input.

---

## Build System & Configuration

Standard Makefile targets shared across all Dell Terraform providers:

| Target | Purpose | Hardware Required |
|--------|---------|-------------------|
| `make build` | Compile provider binary | No |
| `make install` | Install to `~/.terraform.d/plugins/` | No |
| `make test` | Run unit tests | No |
| `make testacc` | Run acceptance tests | **Yes** |
| `make check` | Format, lint, vet | No |
| `make gosec` | Security scan | No |
| `make cover` | Generate coverage report | No |
| `make generate` | Generate documentation | No |

GoReleaser configuration: CGO_ENABLED=0, platforms (freebsd, windows,
linux, darwin), architectures (amd64, 386, arm, arm64).

### Evolution

TBD — requires SME input.

---

## Operational Knowledge

**Unit tests:** `bytedance/mockey` for runtime function patching.
No hardware required. Run with `make test`.

**Acceptance tests:** `terraform-plugin-testing` against live hardware.
Creates real resources. Run with `TF_ACC=1 make testacc`. Clean up
manually if tests fail mid-run.

### Evolution

TBD — requires SME input.

---

## General Context

### Open Issues

TBD — requires code scanning for TODO/FIXME/HACK markers.

### Glossary

| Term | Definition |
|------|------------|
| Plugin Framework | HashiCorp's Terraform Plugin Framework (`terraform-plugin-framework`) |
| mockey | `bytedance/mockey` — runtime function patching for unit tests |
| OME | Environment variable prefix for this provider |

---

## References

- [Terraform Plugin Framework Docs](https://developer.hashicorp.com/terraform/plugin/framework)
- [Dell Terraform Registry](https://registry.terraform.io/namespaces/dell)

---

## Governance Spec Discrepancies

No discrepancies detected between code/SME knowledge and loaded
governance specs.
