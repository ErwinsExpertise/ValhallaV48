![Alt text](img/logo.png?raw=true "Valhalla")

[![Actions Status](https://github.com/Hucaru/Valhalla/workflows/Go/badge.svg)](https://github.com/Hucaru/Valhalla/actions)
[Visit our Discord channel](https://discord.gg/KHky9Qy9jF)

## What is this?

Valhalla is a MapleStory **v48** server project.

If you just want to get the server running on your own PC, use **dev mode**. It is the easiest setup because login, world, channel, and cash shop all start together in one window.

## Super Quick Start (Windows)

1. **Get a MapleStory v48 install**
   - Use https://msdl.xyz/
2. **Download the localhost-ready client**
   - https://mega.nz/file/EFV1wA4B#Y7oLs0xrRv9bbR7B8slUF-D1Sq0uHb2EWAxN-IeOlW0
3. **Install MySQL**
4. **Import the database**
   - Use `/sql/maplestory.sql`
5. **Convert your v48 WZ files into NX files**
   - Easiest option: run `setup\convert-wz-to-nx.bat`
6. **Edit `config_dev.toml`**
   - Set your MySQL password
7. **Start Valhalla in dev mode**

```powershell
.\Valhalla.exe -type dev -config config_dev.toml
```

After that, launch the client and connect to `127.0.0.1`.

## The Important Part: v48 Uses Multiple WZ Files

Do **not** follow old instructions that only mention `Data.wz`.

For v48, you should convert the **entire MapleStory folder that contains all the WZ files** such as:

- `Base.wz`
- `Character.wz`
- `Effect.wz`
- `Etc.wz`
- `Item.wz`
- `Map.wz`
- `Mob.wz`
- `Npc.wz`
- `Quest.wz`
- `Reactor.wz`
- `Skill.wz`
- `Sound.wz`
- `String.wz`
- `TamingMob.wz`
- `UI.wz`

Valhalla can load either:

- a single `Data.nx` file, or
- a folder full of `.nx` files

The included helper script creates an easy `nx` folder in this repository.

## Dev Mode

This is the recommended way to run Valhalla for almost everyone.

```powershell
.\Valhalla.exe -type dev -config config_dev.toml
```

Why dev mode is recommended:

- ✅ easiest setup
- ✅ one command
- ✅ one window
- ✅ auto-register enabled in the sample config
- ✅ good for solo play, testing, and learning

## Easy WZ Conversion

1. Download the converter from:
   - https://github.com/ErwinsExpertise/go-wztonx-converter/releases/tag/v0.1.1
2. Place the converter `.exe` inside the `setup` folder, or anywhere on your `PATH`
3. Double-click:
   - `setup\convert-wz-to-nx.bat`
4. Pick your MapleStory v48 folder when prompted
5. Wait for conversion to finish

By default, the converted NX files will be placed in:

```text
<repo>\nx
```

## If Your NX Files Are Somewhere Else

You can point Valhalla at them in either of these ways:

### Option 1: config file

Add this to your config file:

```toml
[nx]
path = "C:/path/to/your/nx"
```

### Option 2: command line flag

```powershell
.\Valhalla.exe -type dev -config config_dev.toml -nx "C:\path\to\your\nx"
```

## Recommended Reading Order

- **[docs/Installation.md](docs/Installation.md)** - easiest full setup guide
- **[docs/Local.md](docs/Local.md)** - running on your own machine
- **[docs/Configuration.md](docs/Configuration.md)** - config file help
- **[docs/README.md](docs/README.md)** - docs index

## Other Ways to Run It

These are more advanced and not recommended for most first-time users:

- [docs/Docker.md](docs/Docker.md)
- [docs/Kubernetes.md](docs/Kubernetes.md)
- [docs/Building.md](docs/Building.md)

## Acknowledgements

- Hucaru for the original v28 project
  - [Valhalla](https://github.com/Hucaru/Valhalla)
- Sunnyboy for providing a [list](http://forum.ragezone.com/f921/library-idbs-versions-named-addresses-987815/) of idbs for which this project would not have started
- The following projects were used to help reverse packet structures that were not clearly shown in the idb
    - [Vana](https://github.com/retep998/Vana)
    - [WvsGlobal](https://github.com/diamondo25/WvsGlobal)
    - [OpenMG](https://github.com/sewil/OpenMG)
- [NX](https://nxformat.github.io/) file format (see acknowledgements at link)

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
