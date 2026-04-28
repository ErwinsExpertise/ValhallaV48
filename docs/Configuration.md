# Configuration Guide

This guide focuses on the settings most people actually need for local **dev mode**.

## Most Important Command

```powershell
.\Valhalla.exe -type dev -config config_dev.toml
```

## Most Important File

For almost everyone, the file you want is:

```text
config_dev.toml
```

## Command Line Flags

| Flag | What it does |
|---|---|
| `-type` | Which server mode to start (`dev`, `login`, `world`, `channel`, `cashshop`) |
| `-config` | Which TOML file to use |
| `-nx` | Path to `Data.nx` or a folder full of `.nx` files |
| `-metrics-port` | Metrics port |
| `-channels` | Number of channels in dev mode |

## Recommended Dev Mode Command

```powershell
.\Valhalla.exe -type dev -config config_dev.toml
```

## Recommended Dev Mode Command With Custom NX Path

```powershell
.\Valhalla.exe -type dev -config config_dev.toml -nx "C:\path\to\your\nx"
```

## Database Settings

These are the settings most users need to edit first:

```toml
[database]
address = "127.0.0.1"
port = "3306"
user = "root"
password = "your_mysql_password"
database = "maplestory"
```

## NX Data Settings

Valhalla supports either:

- a single `Data.nx` file, or
- a folder containing many `.nx` files

### Config file option

```toml
[nx]
path = ""
```

Set `path` only if your NX files are not in the normal easy locations.

### Search order when `path` is empty

Valhalla will try these locations automatically:

1. `Data.nx`
2. `nx`
3. `Data.nx` next to the executable
4. `nx` next to the executable
5. the legacy `../v48/wz/nx` path

## Dev Config Highlights

The sample `config_dev.toml` already uses easy local defaults:

- localhost database/network settings
- `autoRegister = true`
- boosted rates for testing

## Separate Config Files

If you are not using dev mode, these files are available:

- `config_login.toml`
- `config_world.toml`
- `config_channel_1.toml`
- `config_channel_2.toml`
- `config_channel_3.toml`
- `config_cashshop.toml`

They now also support:

```toml
[nx]
path = ""
```

## Environment Variable Option

You can also set the NX path with:

```text
VALHALLA_NX_PATH
```

Example:

```powershell
$env:VALHALLA_NX_PATH = 'C:\path\to\your\nx'
.\Valhalla.exe -type dev -config config_dev.toml
```

## Common Fixes

### Wrong MySQL password

Update `config_dev.toml`.

### NX files are in the wrong place

Use either:

- `[nx].path`, or
- `-nx`

### You converted only one WZ file

For v48, convert the **whole MapleStory folder**, not just `Data.wz`.

## See Also

- [Installation Guide](Installation.md)
- [Local / Dev Mode Guide](Local.md)
