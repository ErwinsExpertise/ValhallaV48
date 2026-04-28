# Building from Source

Only use this guide if you need to compile Valhalla yourself.

If you already have a working executable, go back to:

- [Installation Guide](Installation.md)
- [Local / Dev Mode Guide](Local.md)

## Requirements

- Go 1.25+
- Git
- MySQL
- NX data (`Data.nx` or a folder of `.nx` files)

## Build

```powershell
go build .
```

This creates:

- `Valhalla.exe` on Windows
- `Valhalla` on Linux/macOS

## Run After Building

The recommended command is still dev mode:

```powershell
.\Valhalla.exe -type dev -config config_dev.toml
```

If your NX files are somewhere else:

```powershell
.\Valhalla.exe -type dev -config config_dev.toml -nx "C:\path\to\your\nx"
```

## Tests

```powershell
go test ./...
```

## Build Check

```powershell
go build ./...
```

## Notes About NX Data

For v48, convert the **full MapleStory WZ folder**, not just one file.

Valhalla can read either:

- `Data.nx`
- or a folder like `nx\` containing many `.nx` files
