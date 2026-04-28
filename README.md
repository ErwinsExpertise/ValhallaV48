![Alt text](img/logo.png?raw=true "Valhalla")

[![Actions Status](https://github.com/Hucaru/Valhalla/workflows/Go/badge.svg)](https://github.com/Hucaru/Valhalla/actions)
[Visit our Discord channel](https://discord.gg/KHky9Qy9jF)
## What is this?

This project exists to preserve and archive an early version of the game (v48 of global).

## Acknowledgements

- Hucaru for the original v28 project
  - [Valhalla](https://github.com/Hucaru/Valhalla) 
- Sunnyboy for providing a [list](http://forum.ragezone.com/f921/library-idbs-versions-named-addresses-987815/) of idbs for which this project would not have started
- The following projects were used to help reverse packet structures that were not clearly shown in the idb
    - [Vana](https://github.com/retep998/Vana)
    - [WvsGlobal](https://github.com/diamondo25/WvsGlobal)
    - [OpenMG](https://github.com/sewil/OpenMG)
- [NX](https://nxformat.github.io/) file format (see acknowledgements at link)

## Getting Started

Valhalla supports multiple deployment methods. Choose the one that best fits your needs:

📚 **[Installation Guide](docs/Installation.md)** - Start here! Covers Data.wz conversion and client setup

### Quick Links by Deployment Method

- 🚀 **[Dev Mode](#dev-mode)** - Run all servers in one process (easiest for solo play and testing)
- 🖥️ **[Local Setup](docs/Local.md)** - Run directly on your machine (best for quick testing)
- 🐳 **[Docker Setup](docs/Docker.md)** - Run using Docker Compose (recommended for most users)
- ☸️ **[Kubernetes Setup](docs/Kubernetes.md)** - Deploy to a Kubernetes cluster (for production)
- 🔨 **[Building from Source](docs/Building.md)** - Build for development work

### Dev Mode
Dev mode is the simplest way to run Valhalla for solo play or testing. It starts all four server types (login, world, channel, and cashshop) in a single process, eliminating the need to manage multiple processes.

**Quick Start:**

```bash
# Using the provided dev config
./Valhalla -type dev -config config_dev.toml

# Or build from source
go build
./Valhalla -type dev -config config_dev.toml
```

The dev mode:
- ✅ Runs all servers (login, world, channel, cashshop) in one process
- ✅ Uses localhost networking for inter-server communication
- ✅ Has auto-register enabled by default for easy testing
- ✅ Configured with 2x EXP/Drop/Mesos rates for faster testing
- ✅ Ideal for solo play and development

**Note:** Dev mode still requires a MySQL database. For production deployments, use separate processes for better reliability and scalability.

### Configuration

⚙️ **[Configuration Guide](docs/Configuration.md)** - Complete reference for all configuration options

All server types support both TOML configuration files and environment variables. See the Configuration Guide for details on:
- Command line flags (`-type`, `-config`, `-metrics-port`)
  - Server types: `login`, `world`, `channel`, `cashshop`, or `dev` (all-in-one)
- Database settings
- Server-specific options (login, world, channel, cashshop)
- Network configuration
- Performance tuning

For dev mode, use `config_dev.toml` which includes all necessary server configurations in a single file.

## Advanced Topics

### NPC Scripting

NPCs are scripted in JavaScript powered by [goja](https://github.com/dop251/goja). For detailed NPC chat formatting codes and scripting information, see the scripts directory and existing NPC implementations.

For NPC chat display formatting reference, see the [NPC Chat Formatting](#npc-chat-formatting) section below.

### Production Deployments

- **[Kubernetes](docs/Kubernetes.md)** - Production-ready deployment with Helm, ingress, scaling, and monitoring
- **[Docker](docs/Docker.md)** - Containerized deployment with Docker Compose

## NPC Chat Formatting

When scripting NPCs in JavaScript, use these formatting codes:

- `#b` - Blue text
- `#c[itemid]#` - Shows how many [itemid] the player has in inventory
- `#d` - Purple text
- `#e` - Bold text
- `#f[imagelocation]#` - Shows an image from .wz files
- `#g` - Green text
- `#h #` - Shows the player's name
- `#i[itemid]#` - Shows a picture of the item
- `#k` - Black text
- `#l` - Selection close
- `#m[mapid]#` - Shows the name of the map
- `#n` - Normal text (removes bold)
- `#o[mobid]#` - Shows the name of the mob
- `#p[npcid]#` - Shows the name of the NPC
- `#q[skillid]#` - Shows the name of the skill
- `#r` - Red text
- `#s[skillid]#` - Shows the image of the skill
- `#t[itemid]#` - Shows the name of the item
- `#v[itemid]#` - Shows a picture of the item
- `#x` - Returns "0%" (usage varies)
- `#z[itemid]#` - Shows the name of the item
- `#B[%]#` - Shows a progress bar
- `#F[imagelocation]#` - Shows an image from .wz files
- `#L[number]#` - Selection open
- `\r\n` - Moves down a line
- `\r` - Return carriage
- `\n` - New line
- `\t` - Tab (4 spaces)
- `\b` - Backwards

Reference from [RageZone forums](http://forum.ragezone.com/f428/add-learning-npcs-start-finish-643364/)
