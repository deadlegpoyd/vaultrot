# vaultrot

> A CLI tool for rotating secrets across multiple secret backends (Vault, AWS SSM, and Doppler) with dry-run support.

---

## Installation

```bash
go install github.com/youruser/vaultrot@latest
```

Or download a pre-built binary from the [releases page](https://github.com/youruser/vaultrot/releases).

---

## Usage

```bash
# Rotate all secrets in a Vault backend
vaultrot rotate --backend vault --path secret/myapp

# Rotate secrets in AWS SSM with a dry run first
vaultrot rotate --backend ssm --path /myapp/prod --dry-run

# Rotate secrets in Doppler for a specific project
vaultrot rotate --backend doppler --project myapp --config prd
```

### Flags

| Flag | Description |
|------|-------------|
| `--backend` | Secret backend to target (`vault`, `ssm`, `doppler`) |
| `--path` | Path or prefix of secrets to rotate |
| `--dry-run` | Preview changes without applying them |
| `--config` | Path to a vaultrot config file (default: `./vaultrot.yaml`) |

### Config File

```yaml
backend: vault
vault:
  address: https://vault.example.com
  token_env: VAULT_TOKEN
  path: secret/myapp
dry_run: false
```

---

## Supported Backends

- **HashiCorp Vault**
- **AWS SSM Parameter Store**
- **Doppler**

---

## License

[MIT](LICENSE)