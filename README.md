# SpecGate

A CLI tool for enforcing OpenAPI specification readiness.

SpecGate evaluates OAS files against a set of quality rules and surfaces errors and warnings before they reach production. It exits with a non-zero status code when errors are detected, making it suitable for use as a CI quality gate.

## Installation

SpecGate currently requires building from source. Homebrew support is planned for a future release.

### Requirements

- Go 1.21 or later

### Build from source
```bash
git clone https://github.com/itsdeannat/specgate-cli.git
cd specgate-cli
make build
```

To run SpecGate from any directory, move the binary to your PATH:
```bash
mv specgate /usr/local/bin/
```

## Usage
```bash
# Check a spec for errors and warnings
specgate check oas.json

# Check a spec and treat warnings as errors
specgate check oas.json --strict

# Output results as JSON
specgate check oas.json --format json

# Get AI-powered suggestions for missing documentation
specgate advise oas.json

# View the rules SpecGate enforces
specgate rules
```

## Example
```
$ specgate check oas.json

ERRORS
------

Missing operation summaries for 1 operation(s):

- GET /menu/{itemId}

Missing error responses (4xx/5xx/default) for 4 operation(s):

- POST /orders
- GET /loyalty/{customerId}
- POST /loyalty/{customerId}
- GET /menu


WARNINGS
--------

Missing operation descriptions for 3 operation(s):

- GET /loyalty/{customerId}
- POST /loyalty/{customerId}
- GET /menu/{itemId}


Run with --strict to treat warnings as errors.
```

## Rules

SpecGate enforces the following rules. Run `specgate rules` to view them at any time.

### Errors

| Rule | Description |
|------|-------------|
| Missing operation summary | Operations must include a summary |
| Missing 2xx response | Operations must include at least one success response |
| Missing error response | Operations must include at least one error response (4xx, 5xx, or default) |
| Missing 2xx response description | Success responses must include a description |
| Missing error response description | Error responses must include a description |
| Missing servers object | A `servers` object must be included |
| Placeholder server URL | Server URLs must not contain `example.com` or `localhost` |

### Warnings

| Rule | Description |
|------|-------------|
| Missing operation description | Operations should include a description |

## Flags

| Flag | Description |
|------|-------------|
| `--strict` | Treat warnings as errors |
| `--format json` | Output results as JSON |

## License

MIT