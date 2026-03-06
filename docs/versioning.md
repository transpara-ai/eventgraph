# Versioning

## Semantic Versioning

EventGraph follows [Semantic Versioning 2.0.0](https://semver.org/).

```
MAJOR.MINOR.PATCH
```

- **MAJOR** — Incompatible API changes
- **MINOR** — New functionality, backwards compatible
- **PATCH** — Bug fixes, backwards compatible

## Pre-1.0 Policy

While the version is `0.x.y`:

- The project is in active development
- Interfaces may change between minor versions
- Breaking changes will be documented in CHANGELOG.md
- We will provide migration guidance for all breaking changes
- Patch versions are always backwards compatible

## Post-1.0 Policy

Once `1.0.0` is released:

- Public interfaces are stable
- Breaking changes only in major versions
- Interface changes require the RFC process (see CONTRIBUTING.md)
- Deprecation before removal (minimum one minor version)

## What Counts as a Public Interface

- The `Store` interface
- The `Primitive` interface (Process, lifecycle, mutations)
- The `IIntelligence` interface
- The `IDecisionMaker` interface
- The `Bus` interface
- Event struct fields
- EGIP message formats and envelope structure
- The `eg` CLI command interface

## What Doesn't Count

- Internal package APIs (anything not importable)
- Default primitive behaviour (can be improved without breaking)
- Decision tree structure (evolves by design)
- Documentation structure

## Language Package Versions

Each language package follows the same version number. When the Go reference implementation is at `0.3.0`, the Rust package should also target `0.3.0` compatibility.

If a language package falls behind, it declares which version of the specification it implements.

## Version File

The current version is in the `VERSION` file at the repository root. This is the single source of truth for the version number. Build scripts and packages should read from this file.
