# Mezamero / 目覚めろ

> [!NOTE]
> This project is almost 100% ai slop by Composer 2 Fast (because I ran out of Opus 4.6 since first week), including below readme

Small **Wake-on-LAN** web UI and API: one Go binary serves a static SPA (`index.html` and assets via `go:embed`) and sends magic packets for configured devices.

## Requirements

- Go 1.26+
- Optional: `golangci-lint` for `make lint`

## Configuration

Provide a YAML file (path via `-config`, default `config.yaml`). Each device entry includes:

| Field              | Required | Description                                                  |
| ------------------ | -------- | ------------------------------------------------------------ |
| `mac`              | yes      | MAC address (`:` or `-` separators)                          |
| `ip`               | no       | Optional host IP (shown next to the MAC on the card)        |
| `name`             | yes      | Display name                                                 |
| `name_original`    | no       | Secondary label (smaller gray text in UI)                    |
| `description`      | no       | Longer human-readable blurb (shown as body text on the card) |
| `spec`             | no       | Short technical line (shown under a “Spec” label, monospace) |
| `image`            | yes      | Public image URL (e.g. S3 HTTPS URL)                         |
| `background_color` | yes      | Hex color for button gradient                                |
| `foreground_color` | yes      | Hex color for text on the card                               |

See `config.example.yaml`.

## Run locally

```bash
make install   # download modules (like npm install)
cp config.example.yaml config.yaml
# edit config.yaml
make dev
```

Open [http://127.0.0.1:8080](http://127.0.0.1:8080).

### Flags

- `-config` — path to YAML (default `config.yaml`)
- `-addr` — listen address (default `:8080`)

## API

- `GET /api/devices` — JSON array of configured devices (same fields as YAML).
- `POST /api/wake` — body `{"mac":"aa:bb:cc:dd:ee:ff"}`. MAC must match a configured device. Responds `{"status":"ok"}` on success.

## Logging

On successful startup, one line is written to stdout with the listen address and a browser URL; the `host:port` part (e.g. `localhost:8080`) is **bold** in ANSI-capable terminals.

Routine logs also appear when a WoL packet is sent successfully (device name, MAC, UTC timestamp RFC3339Nano). Configuration and fatal server errors go to stderr and exit non-zero.

## Docker

Build:

```bash
docker build -t mezamero .
```

Run with your config mounted (image expects config at `/config/config.yaml`):

```bash
docker run --rm -p 8080:8080 \
  -v /absolute/path/to/config.yaml:/config/config.yaml:ro \
  mezamero
```

## Makefile

| Target         | Purpose                                             |
| -------------- | --------------------------------------------------- |
| `make install` | `go mod download` + `go mod verify` (fetch modules) |
| `make dev`     | `go run` with `config.example.yaml` and `:8080`     |
| `make build`   | Produce `./mezamero` binary                         |
| `make preview` | Build and run the binary                            |
| `make format`  | `gofmt` (and `goimports` if installed)              |
| `make lint`    | `golangci-lint run`                                 |
| `make test`    | `go test ./...`                                     |

Override config or address: `make dev CONFIG=config.yaml ADDR=:3000`.

## Notes

- WoL uses UDP broadcast to `255.255.255.255:9`. The host running Mezamero must be able to reach the target subnet (same LAN or routed broadcast path, depending on your network).
- Images are loaded by the browser from the URLs you configure; ensure CORS and HTTPS as needed for your bucket.

## License

Use and modify as you like for your deployment.
