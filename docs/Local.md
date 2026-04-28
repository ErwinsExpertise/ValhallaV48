# Local / Dev Mode Guide

This is the **recommended** way to run Valhalla.

If you are playing or testing on your own PC, use **dev mode**.

## Why Dev Mode?

Dev mode is best for most users because:

- it starts everything together
- it only needs one command
- it is easier to troubleshoot
- it is perfect for solo testing

## Before You Start

Make sure you already finished:

1. [Installation Guide](Installation.md)
2. MySQL setup
3. WZ to NX conversion
4. `config_dev.toml` database password update

## Start the Server

From the Valhalla folder, run:

```powershell
.\Valhalla.exe -type dev -config config_dev.toml
```

If your NX files are in a custom location:

```powershell
.\Valhalla.exe -type dev -config config_dev.toml -nx "C:\path\to\your\nx"
```

## What Dev Mode Starts

Dev mode starts these together:

- Login server
- World server
- Channel server(s)
- Cash shop server

## Default Local Addresses

- Client login: `127.0.0.1:8484`
- World internal port: `127.0.0.1:8584`
- Channel 1 internal port: `127.0.0.1:8685`
- Cash shop internal port: `127.0.0.1:8600`

## First Login

The sample `config_dev.toml` uses:

```toml
autoRegister = true
```

So the first username/password you type in the client will create an account automatically.

## Easiest Folder Layout

If you use the included conversion helper, a simple layout is:

```text
ValhallaV48/
├── Valhalla.exe
├── config_dev.toml
├── drops.json
├── reactors.json
├── reactor_drops.json
├── scripts/
└── nx/
    ├── Base.nx
    ├── Character.nx
    ├── Effect.nx
    └── ...
```

Valhalla also supports a single `Data.nx` file.

## Custom NX Location

You have 2 easy options.

### Option 1: Put it in the config

```toml
[nx]
path = "C:/path/to/your/nx"
```

### Option 2: Pass it on the command line

```powershell
.\Valhalla.exe -type dev -config config_dev.toml -nx "C:\path\to\your\nx"
```

The command line flag wins if both are set.

## If You Really Want Separate Processes

Most people should skip this section.

If you want to start each server separately, you can still use:

- `config_login.toml`
- `config_world.toml`
- `config_channel_1.toml`
- `config_cashshop.toml`

But this is more work and easier to misconfigure.

## Troubleshooting

### It says it cannot load NX data

- Make sure the conversion completed successfully
- Make sure you converted the **entire v48 WZ folder**, not just one file
- Make sure the NX files are in `nx/`, `Data.nx`, or the path you passed with `-nx`

### It cannot connect to MySQL

- Check the password in `config_dev.toml`
- Make sure MySQL is running
- Make sure the database is named `maplestory`

### The client says unable to connect

- Make sure Valhalla is still running
- Make sure you launched the localhost-ready client
- Make sure port `8484` is not blocked or already in use

## Next

- [Configuration Guide](Configuration.md)
- [Building from Source](Building.md)
