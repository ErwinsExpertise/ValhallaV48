# AGENTS.md

## Verification
- Primary CI path is `go build .` and `go test ./...` from the repo root. GitHub Actions only runs those two commands.
- Focused checks use normal Go package targeting, e.g. `go test ./channel`, `go test ./login`, `go test ./nx`, or `go test ./channel -run TestName`.

## Read First
- Highest-value docs for local work are `docs/Installation.md`, `docs/Local.md`, and `docs/Configuration.md`, in that order.
- `docs/Docker.md` and `docs/Kubernetes.md` are deployment guides, not the default dev workflow.
- `docs/Admin-Commands.md` is the in-game GM command reference; commands are chat commands prefixed with `/`.

## Entrypoints
- This is a single Go module (`go.mod` at root), not a monorepo.
- `main.go` is the only binary entrypoint. `-type` selects `login`, `world`, `channel`, `cashshop`, or `dev`.
- The root `server_*.go` files are process wrappers. Most gameplay logic lives under `channel/`; auth/account flow is under `login/`; inter-server coordination is under `world/`; NX parsing/loading is under `nx/`.

## Runtime Assets
- The server expects these repo-root runtime assets to exist alongside the binary or working directory: `drops.json`, `reactors.json`, `reactor_drops.json`, and `scripts/`.
- NPC, event, portal, and reactor scripts live under `scripts/` and are loaded directly by the channel server. They are plain `.js` files compiled with `goja`.
- Script files are watched with `fsnotify`; write/create/remove events hot-reload them without a rebuild.
- For v48 local setup, the supported path to optimize for is a repo-root `nx/` directory containing the converted `.nx` files. The helper script `setup/convert-wz-to-nx.ps1` writes there by default.

## Config And Env
- Config loading is `TOML first, then env overrides` via `viper` in `server_config.go`.
- Env vars use the `VALHALLA_` prefix with section/key names joined by underscores, e.g. `VALHALLA_DATABASE_PASSWORD`, `VALHALLA_NX_PATH`, `VALHALLA_CHANNEL_LISTENPORT`.
- `-nx` overrides `[nx].path` from config.
- Runtime still probes `Data.nx` paths in `server_nx.go`, but repo docs and setup scripts for v48 assume converting the full WZ set into `nx/`. Prefer documenting and testing with `nx/` unless you are specifically working on container packaging.

## Dev Mode
- The recommended local run path is `./Valhalla(.exe) -type dev -config config_dev.toml`.
- Dev mode starts login, world, cash shop, and multiple channel servers in one process.
- `-channels` only affects dev mode. Extra dev channels do not read `config_channel_*.toml`; `server_dev.go` clones the shared channel config and assigns ports starting at `8685`.
- The sample dev config is intentionally friendly for local testing: `autoRegister = true` and boosted world rates.

## Tests And Data Quirks
- `go test ./...` passes without local NX data because the NX-dependent tests skip when the legacy `../v48/wz/nx` directory is missing.
- If you need those smoke tests to execute instead of skip, place the converted v48 NX directory at `../v48/wz/nx` or adjust the tests deliberately.

## Deployment Notes
- Docker builds the root binary and copies the JSON data files plus `scripts/`.
- Some deployment docs/examples still refer to mounting `Data.nx`, but that is a packaging path, not the main local v48 workflow. Do not "fix" local docs or tests toward `Data.nx` without checking whether you are touching container-specific behavior.
- Releases are produced with GoReleaser on pushes to `main`; release archives intentionally include `config_*.toml`, the JSON data files, and `scripts/**`.

## MCP Servers
- `nx-mcp` is the fastest way to inspect converted game data. Call `nx_load` on the directory containing `.nx` files first, then use `nx_list_node`, `nx_get_node`, and `nx_search`.
- `ida-pro-mcp` is for static reverse-engineering work. Start with `survey_binary`; use `open_file` if no database is open, `select_instance` if multiple IDBs are open, then `analyze_function`, `decompile`, `xrefs_to`, `callgraph`, or rename/type tools as needed.
- `x64dbg` is for live debugging. Check `IsDebugActive`/`IsDebugging` first, then prefer high-level queries like `GetRegisterDump`, `GetModuleList`, `DisasmGetInstructionRange`, `XrefGet`, `GetCallStack`, and breakpoint helpers.
- `x64dbg` command execution is picky: `ExecCommand` arguments must be comma-separated after the command name, e.g. `findallmem 0x140001000,CC,20480`.
