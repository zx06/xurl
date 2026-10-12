# xurl

**xurl** is an AI-first HTTP client CLI (Go), sibling to [xsql](https://github.com/zx06/xsql). It runs [VS Code / JetBrains style `.http` files](https://httpyac.github.io/) with YAML profiles, structured JSON output, and OS keyring-backed secrets.

## Install

```bash
git clone https://github.com/zx06/xurl.git
cd xurl
go install ./cmd/xurl
# or: go build -o xurl ./cmd/xurl
```

Requires **Go 1.22+**.

## Quick start (~5 minutes)

1. **Copy example config**

   ```bash
   mkdir -p ~/.config/xurl
   cp examples/xurl.yaml.example ~/.config/xurl/xurl.yaml
   ```

   Edit `profiles.httpbin.vars.host` if needed. For real secrets, store a token in your OS keyring and reference it:

   ```bash
   # Linux (secret-service); macOS Keychain; Windows Credential Manager
   # Service/account must match keyring:service/user in YAML
   ```

   Example profile snippet:

   ```yaml
   profiles:
     httpbin:
       base_url: https://httpbin.org
       allow_plaintext: false
       vars:
         host: https://httpbin.org
       auth:
         bearer: keyring:xurl/demo-token
   ```

   For local dev only, set `allow_plaintext: true` and use plaintext tokens in YAML.

2. **Run a named request**

   ```bash
   xurl run examples/example.http health -p httpbin -f json
   ```

   Output shape:

   ```json
   {
     "ok": true,
     "request": { "name": "health", "method": "GET", "url": "..." },
     "status": 200,
     "headers": { "...": "..." },
     "body": "..."
   }
   ```

3. **List profiles**

   ```bash
   xurl profile list -f json
   xurl profile show httpbin -f json
   ```

## CLI

| Command | Description |
|--------|-------------|
| `xurl run <file.http> [request-name]` | Run one request (`### name` blocks in multi-request files) |
| `xurl profile list` | List profiles |
| `xurl profile show <name>` | Show profile (secrets masked) |

Global flags:

- `-p, --profile` — profile name (or `XURL_PROFILE`)
- `-f, --format json\|table` — output format (or `XURL_FORMAT`)
- `--config` — path to `xurl.yaml` (default: `./xurl.yaml` then `~/.config/xurl/xurl.yaml`)
- `--attr key=value` — extra `{{var}}` bindings (repeatable)

Priority: **CLI > ENV > config file**.

## `.http` file format (subset)

```http
@host=https://api.example.com

### list-items
GET {{host}}/v1/items
Accept: application/json

### create-item
POST {{host}}/v1/items
Content-Type: application/json

{"name":"demo"}
```

- Request line: `METHOD URL` (optional `HTTP/1.1`)
- Headers until blank line; body after blank line
- Multiple requests separated by `### optional-name`
- Variables: `{{name}}` from profile `vars`, `base_url` / `host`, `@name=value` in file, `--attr`, and environment variables

## Secrets

Secret fields (`auth.bearer`, `auth.basic.password`, `auth.header.value`) support:

- `keyring:service/user` — read from OS keyring ([zalando/go-keyring](https://github.com/zalando/go-keyring))
- Plaintext only when `allow_plaintext: true` on the profile (otherwise xurl exits with a clear error)

## License

MIT — see [LICENSE](LICENSE).
