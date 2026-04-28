# SetField Packet Notes

This documents the confirmed working `SendChannelWarpToMap` / `packetPlayerEnterGame` layout for v48.

Status

- This packet is confirmed to get the client in-game.
- Equipped items render correctly.
- Mesos and inventory slot sizes are included in the working `SetField` packet.

Source

- Builder: `channel/player.go` -> `packetPlayerEnterGame`
- Client handler: `CStage::OnSetField`
- Supporting decode: `GW_CharacterStat::Decode`

Known-good packet layout

1. `int32 channelIDMinusOne`
2. `byte 0`
3. `byte 1`
4. `12 bytes` random seed data (`3 * int32`)
5. `int16 0xD7FF` section mask
6. character stat block:
   - `int32 charID`
   - `13-byte padded name`
   - `byte gender`
   - `byte skin`
   - `int32 face`
   - `int32 hair`
   - `8 bytes zero`
   - `byte level`
   - `int16 job`
   - `int16 str`
   - `int16 dex`
   - `int16 int`
   - `int16 luk`
   - `int16 hp`
   - `int16 maxHP`
   - `int16 mp`
   - `int16 maxMP`
   - `int16 ap`
   - `int16 sp`
   - `int32 exp`
   - `int16 fame`
   - `int32 mapID`
   - `byte mapPos`
7. `byte postStatByte`
8. economy and inventory metadata:
   - `int32 mesos`
   - `byte equipSlots`
   - `byte useSlots`
   - `byte setupSlots`
   - `byte etcSlots`
   - `byte cashSlots`
9. inventory section using the login item-entry format:
   - equipped visible items
   - `byte 0`
   - equipped cash items
   - `byte 0`
   - equip inventory tab
   - `byte 0`
   - use inventory tab
   - `byte 0`
   - setup inventory tab
   - `byte 0`
   - etc inventory tab
   - `byte 0`
   - cash inventory tab
   - `byte 0`
10. optional sections:
   - skills section:
     - `int16 skillCount`
     - repeat count times:
       - `int32 skillID`
       - `int32 level`
       - `int32 masterLevel` for fourth-job skills only
   - cooldown section:
     - `int16 cooldownCount`
     - repeat count times:
       - `int32 skillID`
       - `int16 cooldown`
   - active quest section via `writeActiveQuests`
   - completed quest section via `writeCompletedQuests`
   - `int16 0` minigames
11. teleport rock section:
   - `5 * int32 InvalidMap`
   - `10 * int32 InvalidMap`

Section Mask Map

- `0x0001` stats

Decoded by `GW_CharacterStat::Decode` followed by one extra byte immediately after the stat block.

- `0x0002` mesos/meta int

Decoded by `sub_49B813`.
Reads exactly one `int32`, then passes it through `sub_41298F`.
In the working packet this corresponds to mesos.

- `0x0004` equipped visible items

Client reads a `byte` slot followed by an item entry repeatedly until `byte 0`.
Slots map to negative equipped slots.

- `0x0008` equip inventory tab

Client optionally reads one slot-count byte first when `0x0080` is set, then reads item entries until `byte 0`.

- `0x0010` use inventory tab

Client optionally reads one slot-count byte first when `0x0080` is set, then reads item entries until `byte 0`.

- `0x0020` setup inventory tab

Client optionally reads one slot-count byte first when `0x0080` is set, then reads item entries until `byte 0`.

- `0x0040` etc inventory tab

Client optionally reads one slot-count byte first when `0x0080` is set, then reads item entries until `byte 0`.

- `0x0080` inventory slot sizes

This is not a standalone payload block.
It changes inventory decoding so the client reads one extra `byte` count for each of the five inventory tabs.
In the working packet those bytes are:
`equipSlots`, `useSlots`, `setupSlots`, `etcSlots`, `cashSlots`.

- `0x0100` skills

Decoded as:
`int16 skillCount`, then for each skill: `int32 skillID`, `int32 level`, and `int32 masterLevel` for fourth-job skills.

- `0x0200` active quests

Decoded as `int16 questCount`, then active quest records.
Current server writer uses `writeActiveQuests`.

- `0x0400` minigames

Decoded as `int16 count`, then repeated minigame records.
Current server sends zero count.

- `0x0800` rings

Decoded as multiple ring-related lists/counts.
This section is still omitted from the working packet.

- `0x1000` teleport rocks

Decoded as `5 * int32` regular rock maps, then `10 * int32` VIP rock maps.

- `0x4000` completed quests

Decoded as `int16 questCount`, then completed quest records.
Current server writer uses `writeCompletedQuests`.

- `0x8000` cooldowns

Decoded as `int16 cooldownCount`, then repeated `int32 skillID` + `int16 cooldown` pairs.

Notes

- Earlier attempts that crashed the client included:
  - omitting the section mask entirely
  - using `-1` as the section mask while not serializing the matching sections
  - writing more than `8` bytes after `hair` before `level`
  - appending v55-style trailers after the minimal v48 structure
- Inventory in `SetField` must use the dedicated login serializer `Item.setFieldBytes()`, not `Item.InventoryBytes()`.
- `0x0002` reads a single `int32` and corresponds to mesos.
- `0x0080` causes the client to read one slot-size byte for each of the five inventory tabs.
- The single byte after the stat block is currently treated as an unknown post-stat field.
- The single byte after the stat block is followed immediately by mesos and slot sizes when `0x0002` and `0x0080` are set.
- Cooldowns are now serialized for skills that currently have an active cooldown.
- Skills are now serialized with real learned-skill entries.
- Active and completed quests are now serialized with the existing login quest writers.
- The current baseline still does not fully initialize the server-side field instance flow; it only gets the client into the map.

Expansion strategy

1. Keep this documented layout as the recovery baseline.
2. Only change one `SetField` section at a time.
3. Only set a mask bit once the exact payload for that section is written.
4. After each section change, verify:
   - client still enters the map
   - first movement packet does not crash the client
   - server-side state still remains consistent
