# Cash Shop Set Packet Notes

This documents the v48 `SendChannelSetCashShop` / `packetCashShopSet` layout as traced from IDA and workspace references.

Status

- This document is intended as the reference baseline for rebuilding `packetCashShopSet`.
- It is based on `CStage::OnSetCashShop`, `CCashShop::CCashShop`, `CCashShop::LoadData`, and `CharacterData::Decode`.
- The current server implementation is still in progress and should be validated section-by-section against this document.

Source

- Builder: `cashshop/packets.go` -> `packetCashShopSet`
- Character blob: `channel/player.go` -> `WriteCharacterDataPacket`
- Client entry: `CStage::OnSetCashShop`
- Client character blob: `CharacterData::Decode`
- Client cash shop load: `CCashShop::CCashShop`, `CCashShop::LoadData`

Top-Level Packet Order

1. `CharacterData`
2. `byte bCashShopAuthorized`
3. `string accountName`
4. `CWvsContext::SetSaleInfo`
5. fixed `aBest` block (`0x438` bytes)
6. `DecodeStock` block
7. `DecodeLimitGoods` block
8. `DecodeZeroGoods` block
9. `byte bEventOn`
10. `int32 highestCharacterLevelInAccount`

`CStage::OnSetCashShop` Call Tree

1. `CharacterData::Decode`
2. stage transition helpers
3. `CCashShop::CCashShop`
4. `CCashShop::LoadData`
5. trailing `Decode1()` into `this[316]`

`CharacterData` Section Map

The v48 cash shop path uses a low 16-bit section mask.

Known cash shop mask:

- `0xDF7F`

Enabled sections:

- `0x0001` `CHARACTER`
- `0x0002` `MONEY`
- `0x0004` `ITEMSLOTEQUIP`
- `0x0008` `ITEMSLOTCONSUME`
- `0x0010` `ITEMSLOTINSTALL`
- `0x0020` `ITEMSLOTETC`
- `0x0040` `ITEMSLOTCASH`
- `0x0100` `SKILLRECORD`
- `0x0200` `QUESTRECORD`
- `0x0400` `MINIGAMERECORD`
- `0x0800` `COUPLERECORD`
- `0x1000` `MAPTRANSFER`
- `0x4000` `QUESTCOMPLETE`
- `0x8000` `SKILLCOOLTIME`

Not present in this path:

- `0x0080` `INVENTORYSIZE`
- `0x2000` avatar-related later section

`CharacterData` Order

1. `int16 flags`
2. if `flags & 0x0001`: `GW_CharacterStat::Decode`, then `byte friendMax`
3. if `flags & 0x0002`: `int32 mesos`
4. if `flags & 0x0004`: equip item sections
5. if `flags & 0x0008`: consume inventory
6. if `flags & 0x0010`: install inventory
7. if `flags & 0x0020`: etc inventory
8. if `flags & 0x0040`: cash inventory
9. if `flags & 0x0100`: skill record
10. if `flags & 0x8000`: skill cooldowns
11. if `flags & 0x0200`: active quests
12. if `flags & 0x4000`: completed quests
13. if `flags & 0x0400`: minigame records
14. if `flags & 0x0800`: couple/friend/marriage records
15. if `flags & 0x1000`: map transfer arrays

`CHARACTER` Payload

Write exactly:

1. `int32 charID`
2. `char[13] paddedName`
3. `byte gender`
4. `byte skin`
5. `int32 face`
6. `int32 hair`
7. `int64 petCashID1`
8. `int64 petCashID2`
9. `int64 petCashID3`
10. `byte level`
11. `int16 job`
12. `int16 str`
13. `int16 dex`
14. `int16 int`
15. `int16 luk`
16. `int32 hp`
17. `int32 maxHP`
18. `int32 mp`
19. `int32 maxMP`
20. `int16 ap`
21. `int16 sp`
22. `int32 exp`
23. `int16 fame`
24. `int32 tempExp`
25. `int32 mapID`
26. `byte portal`
27. `int32 playTime`
28. `int16 subJob`
29. `byte friendMax`

Notes:

- HP/MP fields are `int32` here.
- `friendMax` is the byte immediately after the stat block, as confirmed by `CharacterData::Decode`.

Equip and Inventory Sections

`ITEMSLOTEQUIP` is three byte-terminated sublists:

1. equipped visible items
   - repeated: `byte slot`, `item body`
   - terminator: `byte 0`
2. equipped cash items
   - repeated: `byte slot`, `item body`
   - terminator: `byte 0`
3. equip inventory
   - repeated: `byte slot`, `item body`
   - terminator: `byte 0`

Other inventories are also byte-loop sections:

- consume: repeated `byte slot`, `item body`, then `byte 0`
- install: repeated `byte slot`, `item body`, then `byte 0`
- etc: repeated `byte slot`, `item body`, then `byte 0`
- cash: repeated `byte slot`, `item body`, then `byte 0`

Item Body Notes

Every item body begins with:

1. `byte itemType`
2. `int32 itemID`
3. `byte isCash`
4. if cash: `int64 cashSerial`
5. `int64 expireTime`

Then type-specific data:

- equip items: stats/creator/flags plus equip tail fields
- bundle items: amount/creator/flag
- pet items: pet name/level/closeness/fullness/etc

Important note:

- This path does **not** use the same item-entry format as `SetField`.
- Slot prefixes and item bodies must be encoded exactly for `CharacterData`, not reused blindly from login/map entry packets.

`SKILLRECORD`

1. `int16 count`
2. repeated:
   - `int32 skillID`
   - `int32 skillLevel`
   - `int64 expiration`
   - `int32 masterLevel` only when `skillID / 10000 % 100 != 0 && skillID / 10000 % 10 == 2`

`SKILLCOOLTIME`

1. `int16 count`
2. repeated:
   - `int32 skillID`
   - `int16 remainSeconds`

`QUESTRECORD`

1. `int16 count`
2. repeated:
   - `int16 questID`
   - `string value`

`QUESTCOMPLETE`

1. `int16 count`
2. repeated:
   - `int16 questID`
   - `int64 completedFT`

`MINIGAMERECORD`

The v48 client decodes:

1. `int16 count`
2. repeated records of:
   - `int32 gameID`
   - `int32 wins`
   - `int32 draws`
   - `int32 losses`
   - `int32 score`

Current safe baseline is a zero count until exact values are emitted.

`COUPLERECORD`

The v48 client decodes three counts:

1. `int16 coupleCount`
2. repeated couple records
3. `int16 friendCount`
4. repeated friend records
5. `int16 marriageCount`

Current safe baseline is:

- `int16 0`
- `int16 0`
- `int16 0`

`MAPTRANSFER`

1. `5 * int32` regular rock maps
2. `10 * int32` VIP rock maps

Unused slots are written as `999999999`.

`CWvsContext::SetSaleInfo`

Order:

1. `int32 nNotSaleCount`
2. `nNotSaleCount * int32`
3. `int16 modifiedCommodityCount`
4. repeated modified commodity records:
   - `int32 commoditySN`
   - `int32 flags`
   - conditional fields driven by flags
5. `byte discountRateCount`
6. repeated `3-byte` discount rate entries

Modified commodity conditional field order:

- `ITEMID` -> `int32`
- `COUNT` -> `int16`
- `PRIORITY` -> `byte`
- `PRICE` -> `int32`
- `BONUS` -> `byte`
- `PERIOD` -> `int16`
- `REQPOP` -> `int16`
- `REQLEV` -> `int16`
- `MAPLEPOINT` -> `int32`
- `MESO` -> `int32`
- `FORPREMIUMUSER` -> `byte`
- `COMMODITYGENDER` -> `byte`
- `ONSALE` -> `byte`
- `CLASS` -> `byte`
- `LIMIT` -> `byte`
- `PBCASH` -> `int16`
- `PBPOINT` -> `int16`
- `PBGIFT` -> `int16`
- `PACKAGESN` -> `byte count` + `count * int32`

Fixed `aBest` Block

After `SetSaleInfo`, `CCashShop::LoadData` does:

- `DecodeBuffer(..., 0x438)`

That fixed block is the `CS_BEST` table:

- `9 categories * 2 genders * 5 ranks = 90 entries`
- each entry is:
  - `int32 category`
  - `int32 gender`
  - `int32 commoditySN`

Stock / Limit / Zero Goods

After the fixed block:

1. `uint16 stockCount`
   - `count * 8 bytes`
2. `uint16 limitGoodsCount`
   - `count * 104 bytes`
3. `uint16 zeroGoodsCount`
   - `count * 68 bytes`

Stage Trailer

After `CCashShop::LoadData`, `CCashShop::CCashShop` reads one more byte.
Current stage/open reference also indicates a trailing stage footer:

1. `byte bEventOn`
2. `int32 highestCharacterLevelInAccount`

Known High-Risk Areas

- `CharacterData` item body format mismatches
- wrong master-level conditional in `SKILLRECORD`
- wrong section mask width/order
- wrong minigame/couple counts
- malformed `SetSaleInfo` modified commodity records

Rebuild Strategy

1. Keep this document as the rebuild source of truth.
2. Only change one `CharacterData` or cash-shop section at a time.
3. Confirm each change in-game before touching the next section.
4. Treat `go test` as compile safety only, not packet validation.
