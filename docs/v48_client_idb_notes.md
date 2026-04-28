# v48 Client IDB Notes

These notes document packet and item-format findings derived from the v48 client IDB.

## Rechargeable / Star Inventory Encoding

- `CWvsContext::OnInventoryOperation` at `0x71A4F8`
  - New-item inventory operations decode an item object via `sub_49BC56`.
  - For item operation type `0` (new item), the client decodes the full item body after the server-sent `invID` and `slot`.

- `sub_49BC56` at `0x49BC56`
  - Reads the item kind byte, constructs the corresponding client item object with `sub_49BCD0`, then dispatches to the object decode method.

- `sub_49BCD0` at `0x49BCD0`
  - Item kind `2` constructs the bundle-item class used for normal stackable items.
  - The bundle-item vtable is `off_79DB70`.

- Bundle-item base encode/decode:
  - `sub_49BE79` at `0x49BE79`
    - Decodes base item header:
      - `int32 itemID`
      - `byte hasCashID`
      - optional `8` bytes cash ID if present
      - `8` bytes expire time / base trailing qword
  - `sub_49BE21` at `0x49BE21`
    - Encodes the same base header.

- Bundle-item stack encode/decode:
  - `sub_49C5C0` at `0x49C5C0`
    - Decodes bundle-item payload:
      - `int16 amount`
      - `string creator`
      - `int16 flag`
      - if `itemID / 10000 == 207`, decode an additional raw `8` bytes into object offsets `+0x38/+0x3C`
  - `sub_49C53F` at `0x49C53F`
    - Encodes the same payload:
      - `amount`
      - `creator`
      - `flag`
      - if `itemID / 10000 == 207`, encode raw `8` bytes from object offsets `+0x38/+0x3C`

- Bundle-item amount accessor used by client logic:
  - `sub_4A066A` at `0x4A066A`
    - Bundle-item vtable getter at offset `+0x10`
    - Returns the decoded `amount` field (`this + 40/44`), not the star-only 8-byte tail
  - `sub_4A067C` at `0x4A067C`
    - Setter for the same `amount` field

- Client ranged-attack projectile validation:
  - `sub_64CF57` at `0x64CF57`
    - weapon type `45` / `1472063` => `2060xxx`
    - weapon type `46` => `2061xxx`
    - weapon type `47` => `207xxxx`

- Client local projectile-count checks:
  - `sub_69D19D` at `0x69D19D`
    - Finds a valid projectile item by `itemID` + weapon/projectile compatibility via `sub_64CF57`
    - Uses the item virtual count getter to compare against required consumption
  - `sub_69D4EC` at `0x69D4EC`
    - Similar local ammo scan for ranged attack paths

### Why this matters

- The v48 client definitely requires the extra `8` bytes for `207xxxx` bundle items. Omitting them desynchronizes the item decode and can crash the client.
- The client's local ranged-attack ammo checks use the decoded bundle-item `amount` getter, so a star item being treated as empty is most directly explained by the client seeing a bad `amount` or bad local item state.
- The star-only `8` bytes are real and version-specific, but the IDB evidence above shows they are not the `amount` field itself.

## Throwing-Star Empty Message Path

- String id `2543` references:
  - `sub_6AD035` at `0x6AD035`
  - `sub_6ADD4C` at `0x6ADD4C`

- `sub_6AD035` at `0x6AD035`
  - This is one local ranged-attack path.
  - It sets `v23 = skillID / 10000 / 100` to determine projectile family.
  - If the local projectile lookup helper `sub_69D19D(...)` fails and no bypass buff/state is active:
    - family `3` -> `GetString(2542)`
    - family `4` -> `GetString(2543)`
  - Therefore `2543` is the local claw / throwing-star failure message.

- `sub_6ADD4C` at `0x6ADD4C`
  - Another local ranged / special ranged path.
  - Uses the same `sub_69D19D(...)` helper and the same family split:
    - family `3` -> string `2542`
    - family `4` -> string `2543`

- `sub_69D19D` at `0x69D19D`
  - This is the local client-side “find usable projectile” helper.
  - It checks:
    - current weapon type via `sub_4411BF(this[191])`
    - projectile item compatibility via `sub_64CF57`
    - projectile count via the item virtual count getter (`vtable + 0x10`, bundle implementation `sub_4A066A`)
    - an additional item-data requirement via `sub_4F42FD(itemID)`
  - Full behavior confirmed from decompilation:
    - Reads character data via `CWvsContext::GetCharacterData`
    - Computes required projectile count from the skill data block when a skill data pointer is provided
    - If a client state at `dword_80C8A0 + 10036/10044` is non-zero, enters an exact-id mode and looks only for item id `2069999 + state`
    - Otherwise scans USE inventory tab (`sub_44D10D(..., 2, slot)`) for the first item satisfying all of:
      - `sub_64CF57(equippedWeaponID, projectileItemID)`
      - item count getter `>= requiredCount`
      - character level byte at `GW_CharacterStat + 35` `>= sub_4F42FD(projectileItemID)`
    - After that, optionally scans CASH inventory tab for `5020xxx` / `5021xxx` cosmetic projectile items, but this does not gate success of finding the usable projectile itself
  - If it cannot find a valid projectile, the caller surfaces `2542` or `2543` before any attack packet is sent.

### Why this matters

- The client does not show string `2543` because of a server-side attack decode failure.
- It shows `2543` because the local client failed `sub_69D19D`, which means the local item state for the star is not passing the client’s own usability checks.

### Runtime correlation

- Server-side encoder log from repro:
  - `[StarEncode] itemID=2070000 slot=10 invID=2 amount=500 flag=0 cash=false expire=150842304000000000 footer=00 00 00 00 00 00 00 00`
- User confirmed the client inventory UI displays the stack count correctly as `500`.
- User repro case:
  - character level `200`
  - skill `Lucky Seven` level `18` (`4001344`)

Given the client logic above, that rules out:

- server-side attack decode failure
- zero stack count on the wire
- character level requirement failure for `2070000`

and narrows the remaining likely failure causes to:

- the client’s exact-id projectile selection state at `dword_80C8A0 + 10036/10044`
- or some local inventory reconstruction state that still leaves the star item visible in UI but not accepted by `sub_69D19D`
