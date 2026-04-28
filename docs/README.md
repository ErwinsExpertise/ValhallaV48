# Valhalla Documentation

These docs are written with **Windows users** and **dev mode** in mind first.

If you are new, read these in order:

1. **[Installation Guide](Installation.md)**
2. **[Local / Dev Mode Guide](Local.md)**
3. **[Configuration Guide](Configuration.md)**

## Fastest Path

If your goal is just to get Valhalla running on your PC:

1. Get MapleStory **v48** from https://msdl.xyz/
2. Download the localhost-ready client from the main README
3. Install MySQL and import `sql/maplestory.sql`
4. Run `setup\convert-wz-to-nx.bat`
5. Edit `config_dev.toml`
6. Run:

```powershell
.\Valhalla.exe -type dev -config config_dev.toml
```

## Important Note About WZ Files

v48 uses **multiple WZ files**.

So do **not** convert only `Data.wz`.

Instead, convert the **entire MapleStory folder** that contains the full set of WZ files. Valhalla supports either:

- a single `Data.nx` file, or
- a folder containing many `.nx` files

## Documentation Map

### Start Here

- **[Installation Guide](Installation.md)** - full setup from MapleStory install to first launch
- **[Local / Dev Mode Guide](Local.md)** - easiest way to run the server on your PC
- **[Configuration Guide](Configuration.md)** - explains `config_dev.toml`, `config_*.toml`, and the NX path option

### Advanced / Optional

- **[Building from Source](Building.md)** - only if you need to compile the server yourself
- **[Docker Setup](Docker.md)** - advanced/containerized setup
- **[Kubernetes Setup](Kubernetes.md)** - advanced production deployment
- **[Admin Commands](Admin-Commands.md)** - GM/admin command reference

## Quick Navigation

| I want to... | Read this |
|---|---|
| Do the easiest Windows setup | [Installation](Installation.md) → [Local](Local.md) |
| Run everything in one process | [Local](Local.md) |
| Fix config values | [Configuration](Configuration.md) |
| Build the executable myself | [Building](Building.md) |
| Use Docker instead | [Docker](Docker.md) |

## Extra Links

- [Main README](../README.md)
- [go-wztonx-converter](https://github.com/ErwinsExpertise/go-wztonx-converter)
- [NX File Format](https://nxformat.github.io/)
