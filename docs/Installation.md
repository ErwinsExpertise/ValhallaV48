# Installation Guide

This is the easiest setup guide for **Valhalla v48**.

If you are on Windows and just want the server running, this is the guide you want.

## What You Need

- A MapleStory **v48** install
  - You can use https://msdl.xyz/
- The localhost-ready client from the main README
- MySQL
- This Valhalla folder
- The WZ to NX converter
  - https://github.com/ErwinsExpertise/go-wztonx-converter/releases/tag/v0.1.1

## Recommended Setup Path

Use **dev mode** unless you have a specific reason not to.

Dev mode starts login, world, channel, and cash shop together with one command.

## Step 1: Get MapleStory v48

Install or extract MapleStory v48.

After that, you should have a folder that contains many WZ files like:

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

> v48 does **not** use only one `Data.wz` file for this setup flow. Convert the whole folder.

## Step 2: Get the Localhost Client

Use the localhost-ready client from the main README:

- https://mega.nz/file/EFV1wA4B#Y7oLs0xrRv9bbR7B8slUF-D1Sq0uHb2EWAxN-IeOlW0

This is the client you should launch after the server is running.

## Step 3: Install MySQL and Import the Database

1. Install MySQL
2. Create a database named `maplestory`
3. Import:

```text
sql/maplestory.sql
```

### Simple import example

If you already have the MySQL command line available:

```powershell
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS maplestory;"
mysql -u root -p maplestory < sql\maplestory.sql
```

If you prefer a GUI, you can also import `sql/maplestory.sql` with MySQL Workbench.

## Step 4: Convert the WZ Files to NX

### Easiest way

1. Download `go-wztonx-converter.exe`
2. Put it in the `setup` folder of this repo
3. Double-click:

```text
setup\convert-wz-to-nx.bat
```

4. Choose your MapleStory v48 folder
5. Let it finish

By default, the converted output is placed in:

```text
nx\
```

inside this repository.

### Manual converter command

If you want to run the converter yourself, point it at the **full v48 folder**:

```powershell
.\go-wztonx-converter.exe --server "C:\path\to\MapleStoryV48"
```

## Step 5: Edit `config_dev.toml`

Open `config_dev.toml` and check these values:

### Database section

```toml
[database]
address = "127.0.0.1"
port = "3306"
user = "root"
password = "your_mysql_password"
database = "maplestory"
```

### Optional NX section

If you used the helper script and left the output in the default `nx` folder, you usually do not need to change anything.

If your NX files are stored somewhere else, set:

```toml
[nx]
path = "C:/path/to/your/nx"
```

## Step 6: Start Valhalla in Dev Mode

```powershell
.\Valhalla.exe -type dev -config config_dev.toml
```

If your NX files are in a custom location, you can also do this:

```powershell
.\Valhalla.exe -type dev -config config_dev.toml -nx "C:\path\to\your\nx"
```

## Step 7: Launch the Client

When the server is up, launch the localhost-ready client and connect to:

```text
127.0.0.1
```

The login server listens on port `8484`.

## First Login

The sample dev config has `autoRegister = true`.

That means you can type a username and password and the account will be created automatically.

## If Something Goes Wrong

### Database errors

- Make sure MySQL is running
- Make sure the password in `config_dev.toml` matches MySQL
- Make sure the `maplestory` database exists
- Make sure `sql/maplestory.sql` was imported

### NX/WZ errors

- Make sure you converted the **entire v48 WZ folder**
- Make sure the converter finished successfully
- Make sure your NX output exists as either:
  - `Data.nx`, or
  - a folder containing `.nx` files
- If needed, use the `-nx` flag or `[nx].path`

### Client cannot connect

- Make sure Valhalla is running in dev mode
- Make sure the localhost-ready client is the one you launched
- Make sure nothing else is already using port `8484`

## Next

- [Local / Dev Mode Guide](Local.md)
- [Configuration Guide](Configuration.md)
