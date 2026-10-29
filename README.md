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

   The bundled `httpbin` profile needs no secrets. For authenticated APIs, reference the OS keyring as `keyring:service/user` (see `examples/xurl.yaml.example` comments). Plaintext secrets in YAML require `allow_plaintext: true` on that profile.

2. **Run a named request**

   ```bash
   xurl run examples/example.http health -p httpbin -f json
   ```

   **Library** (managed `.http` files under `~/.config/xurl/library` by default):

   ```bash
   xurl lib path
   xurl lib add examples/example.http
   xurl run example.http health -p httpbin -f json   # resolves in library if not a local path
   ```

   **History** (stored in `~/.local/share/xurl/history.json`):

   ```bash
   xurl history list -f json
   xurl history run 1          # re-run by index (1 = most recent)
   xurl run --last             # shortcut for most recent entry
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
| `xurl run <target> [request-name]` | Run one request (filesystem path or library name; `--last` re-runs history) |
| `xurl lib path` / `lib list` / `lib add <file>` | Managed library directory (`library_dir` in config) |
| `xurl history list` / `history run <index>` | Run history (newest = index 1) |
| `xurl profile list` | List profiles |
| `xurl profile show <name>` | Show profile (secrets masked) |

Global flags:

- `-p, --profile` — profile name (or `XURL_PROFILE`)
- `-f, --format json\|table` — output format (or `XURL_FORMAT`)
- `--config` — path to `xurl.yaml` (default: `./xurl.yaml` then `~/.config/xurl/xurl.yaml`)
- `--attr key=value` — extra `{{var}}` bindings (repeatable)

Priority: **CLI > ENV > config file**.

Config keys:

- `library_dir` — default `~/.config/xurl/library` (created on first use)
- History file path is fixed at `~/.local/share/xurl/history.json` (or `$XDG_DATA_HOME/xurl/history.json`), max 100 entries

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
