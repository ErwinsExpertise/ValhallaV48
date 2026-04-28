# NPC Conversion Reference

## Scope

This document maps LeaderMS NPC scripts in `C:\Users\daddy\desktop\MapleDev\LeaderMS-English\scripts\npc` to Valhalla NPC scripts in `scripts/npc`.

Goals:

- preserve behavior, not source syntax
- prefer existing Valhalla logic when it already matches the source
- port only the missing behavior
- add runtime helpers only when script parity actually needs them

## LeaderMS NPC Scripting Model

LeaderMS uses the classic OdinMS style Rhino/Nashorn bridge:

- entry points are `start()` and `action(mode, type, selection)`
- scripts receive `cm`, an instance of `NPCConversationManager`
- script authors usually store conversation position in a global `status`
- `mode` drives forward, backout, and cancel flow
- `selection` carries menu or style answers
- scripts explicitly call `cm.dispose()` to end the conversation

Important runtime pieces:

- `src/scripting/npc/NPCScriptManager.java`
- `src/scripting/npc/NPCConversationManager.java`
- `src/scripting/AbstractPlayerInteraction.java`
- `src/handling/channel/handler/NPCTalkHandler.java`
- `src/handling/channel/handler/NPCMoreTalkHandler.java`

Common LeaderMS APIs:

- dialogue: `sendNext`, `sendPrev`, `sendNextPrev`, `sendOk`, `sendYesNo`, `sendAcceptDecline`, `sendSimple`, `sendStyle`, `sendGetNumber`, `sendGetText`
- shops/storage: `openShop`, direct storage access via `cm.getPlayer().getStorage().sendStorage(...)`
- quests: `getQuestStatus`, `startQuest`, `completeQuest`, `forfeitQuest`
- items: `haveItem`, `gainItem`, `removeAll`, `canHold`
- mesos/exp: `getMeso`, `gainMeso`, `gainExp`
- maps: `warp(map)`, `warp(map, portal)`, `warp(map, portalName)`
- jobs/stats: `getJob`, `changeJob`, direct stat mutation through `getPlayer()` in some job advancers
- party/event helpers: event manager access, PQ helpers, party item/exp helpers

## Valhalla NPC Scripting Model

Valhalla uses Goja and replays the script from the top on each client response.

- scripts are plain JS files executed by `channel/script.go`
- the runtime injects `npc`, `plr`, and `map`
- conversation state is tracked in Go by `npcChatStateTracker`
- helpers like `npc.sendNext()` interrupt execution, send a packet, then resume by re-running the script later
- scripts are written in a blocking style instead of a manual `status` state machine
- the controller exits automatically when the script finishes without another pending prompt

Important runtime pieces:

- `channel/script.go`
- `channel/handlers_client.go` (`npcChatStart`, `npcChatContinue`, `npcShop`)
- `channel/npc.go`

Common Valhalla APIs currently exposed:

- dialogue: `npc.sendNext`, `npc.sendBackNext`, `npc.sendBack`, `npc.sendOk`, `npc.sendYesNo`, `npc.sendSelection`, `npc.sendMenu`, `npc.sendAvatar`, `npc.sendStyles`, `npc.sendNumber`, `npc.sendBoxText`, `npc.sendQuiz`, `npc.sendInputText`, `npc.sendInputNumber`, `npc.sendSlideMenu`
- shops/storage: `npc.sendShop`, `npc.sendStorage`
- player state: `plr.job`, `plr.setJob`, `plr.level`, `plr.getLevel`, `plr.name`, `plr.gender`, `plr.hair`, `plr.setHair`, `plr.face`, `plr.setFace`, `plr.skin`, `plr.setSkinColor`
- items/mesos: `plr.giveItem`, `plr.gainItem`, `plr.itemCount`, `plr.haveItem`, `plr.removeItemsByID`, `plr.removeItemsByIDSilent`, `plr.removeAll`, `plr.canHold`, `plr.canHoldAll`, `plr.mesos`, `plr.getMesos`, `plr.giveMesos`, `plr.gainMesos`, `plr.takeMesos`
- saved locations: `plr.saveLocation`, `plr.getSavedLocation`, `plr.clearSavedLocation`
- inventory capacity: `plr.getEquipInventoryFreeSlot`, `plr.getUseInventoryFreeSlot`, `plr.getSetupInventoryFreeSlot`, `plr.getEtcInventoryFreeSlot`, `plr.getCashInventoryFreeSlot`
- quests: `plr.getQuestStatus`, `plr.checkQuestStatus`, `plr.quest`, `plr.questData`, `plr.setQuestData`, `plr.startQuest`, `plr.completeQuest`, `plr.forfeitQuest`, `plr.questStarted`, `plr.questCompleted`, `plr.questNotStarted`
- maps/events: `plr.warp`, `plr.warpToPortalName`, `map.getMap`, `map.playerCount`, `map.portalEnabled`, `map.isPortalEnabled`, `map.properties`
- party/event helpers: `plr.inParty`, `plr.isPartyLeader`, `plr.partyMembersOnMapCount`, `plr.partyGiveExp`, `plr.startPartyQuest`, `plr.leavePartyQuest`
- event-instance helpers: `plr.getEventProperty`, `plr.setEventProperty`, `plr.finishEvent`, `plr.leaveEvent`, `plr.countMonster`, `plr.spawnMonster`, `plr.gainGuildPoints`
- misc: `plr.isGM`, fame/AP/SP/HP/MP helpers, guild creation/emblem prompts

## API Mapping Table

| LeaderMS | Valhalla | Notes |
| --- | --- | --- |
| `start()` + `action(mode,type,selection)` | top-level linear JS | Rewrite flow, do not emulate `status` blindly |
| `cm.dispose()` | usually implicit end | Only needed conceptually; Valhalla ends when script finishes |
| `cm.sendSimple(text)` | `npc.sendSelection(text)` or `npc.sendMenu(text, ...labels)` | Use `sendSelection` for raw `#L...#` text, `sendMenu` for rebuilt sequential menus |
| `cm.sendStyle(text, styles)` | `npc.sendAvatar(text, ...styles)` or `npc.sendStyles(text, styles)` | Returned value is the chosen index |
| `cm.sendGetText` + `cm.getText()` | `npc.sendBoxText(...)` or `npc.sendInputText(...)` | Pull-style getter becomes a return value |
| `cm.sendGetNumber` | `npc.sendNumber(...)` | Same pattern |
| `cm.sendYesNo` | `npc.sendYesNo(...)` | Returns boolean directly |
| `cm.sendNextPrev` chain | repeated `npc.sendNext`, `npc.sendBackNext`, `npc.sendOk` | Use linear script order |
| `cm.openShop(id)` | `npc.sendShop(items)` | Leader shop IDs must be mapped to explicit goods lists in target |
| `cm.getPlayer().getStorage().sendStorage(...)` | `npc.sendStorage(npcId)` | Runtime-level abstraction already exists in Valhalla |
| `cm.haveItem(id, qty)` | `plr.haveItem(id, qty)` | Canonical helper for new conversions |
| `cm.gainItem(id, qty)` | `plr.gainItem(id, qty)` | Canonical helper for adds and removals |
| `cm.removeAll(id)` | `plr.removeAll(id)` | Script-exposed during pilot tightening |
| `cm.canHold(id)` | `plr.canHold(id, qty)` | Script-exposed during pilot tightening |
| exchange preflight on several outputs | `plr.canHoldAll([[id, qty], ...])` | Canonical helper for multi-item exchanges |
| `cm.getMeso()` / `cm.gainMeso(x)` | `plr.getMesos()` / `plr.gainMesos(x)` | `plr.mesos()` and `plr.giveMesos()` remain valid aliases |
| `cm.warp(map)` | `plr.warp(map)` | Portal-specific variants are also available |
| `cm.changeJob(job)` | `plr.setJob(jobId)` | Full job-advance parity still needs stat-reset support |
| `cm.getQuestStatus(id)` | `plr.getQuestStatus(id)` | Normalized to `0 = not started`, `1 = started`, `2 = completed` |
| `cm.startQuest/completeQuest` | `plr.startQuest/completeQuest` | Semantics are close but should still be verified per NPC |

## Exact Helper Semantics

### `cm.gainItem(id, qty)` -> `plr.gainItem(id, qty)`

Canonical rule for all future conversions:

- positive quantity: create the item through the target runtime and add it through normal Valhalla inventory rules
- zero quantity: no-op, returns success
- negative quantity: remove exactly `abs(qty)` items by ID from normal inventories
- equip vs stackable behavior: delegated to Valhalla item creation and `Player.GiveItem`, which already routes equips and stackables differently
- inventory validation: positive grants fail if the target inventory cannot receive the item
- removal validation: negative quantities fail if the player does not own enough items

Implications:

- use `plr.gainItem(...)` for translated LeaderMS `gainItem(...)`
- do not keep using raw `plr.giveItem(id, -1)` in new conversions even though the runtime now supports it for compatibility
- for large exchanges, preflight with `plr.canHold(...)` or `plr.canHoldAll(...)` before removing source materials

### `cm.sendSimple(text)` exact rule

Use `npc.sendSelection(text)` when:

- the source already builds its own `#L...#...#l` markup
- option IDs are non-sequential or embedded into a larger dynamic text block
- you want the target script to preserve the exact menu text layout

Use `npc.sendMenu(text, ...labels)` when:

- you are rebuilding the menu idiomatically in Valhalla
- the option list is sequential `0..n-1`
- the labels are naturally represented as a flat JS list rather than a prebuilt markup string

Do not mix these styles within one NPC without a reason.

### Quest status normalization

LeaderMS quest checks commonly compare against `MapleQuestStatus.Status` values whose IDs are effectively:

- `0` not started
- `1` started
- `2` completed

Valhalla now documents and exposes the same normalized values:

- `plr.getQuestStatus(id)` -> `0/1/2`
- `plr.questNotStarted(id)` -> boolean
- `plr.questStarted(id)` -> boolean
- `plr.questCompleted(id)` -> boolean

Use the boolean helpers in new conversions when they improve readability.

### Event-backed NPC mapping

LeaderMS commonly splits PQ and transport logic across two layers:

- NPC scripts call `getEventManager(...)` or `getEventInstance()`
- event scripts own timers, properties, map flow, and completion callbacks

Current Valhalla mapping is split differently:

- map-wide transport loops like Ellinia/Orbis/Ludibrium boats are implemented in Go in `channel/boats.go`
- PQ instances are implemented as Goja event scripts under `scripts/event`
- NPC scripts can participate in an active instance through `plr` event helpers, not a full LeaderMS-style `EventManager`

Canonical mapping rules:

| LeaderMS concept | Valhalla equivalent | Notes |
| --- | --- | --- |
| `cm.getPlayer().getEventInstance()` | active event on `plr` via `plr.getEventProperty`, `plr.setEventProperty`, `plr.finishEvent`, `plr.leaveEvent` | Use these for active-instance NPCs |
| `eim.getProperty/setProperty` | `plr.getEventProperty` / `plr.setEventProperty` | In Valhalla these map to the current field-instance properties, which is the established target-side event-state pattern |
| `eim.finishPQ()` | `plr.finishEvent()` | Ends the active event instance |
| `eim.unregisterPlayer(...)` | `plr.leaveEvent()` or explicit event-player warp logic | Prefer `leaveEvent` when the existing event script already defines leave behavior |
| `cm.spawnMonster` | `plr.spawnMonster(id, x, y)` | Spawns on the player's current map instance |
| `cm.countMonster()` | `plr.countMonster()` | Counts mobs on the player's current instance |
| `cm.isLeader()` | `plr.isLeader()` | Alias to party leader check |
| guild PQ reward points | `plr.gainGuildPoints(points)` | Added for event-backed reward flows |

Important limitation:

- Valhalla does not yet expose a full generic LeaderMS-style `EventManager` object to NPC scripts
- Valhalla event-state properties should stay on field instances, not in a separate parallel event-property store
- for transport systems already implemented in Go, prefer the existing Go scheduler and map properties instead of adding a second scheduler in JS
- use KPQ, LPQ, and the boat system as the primary target-side references for event-backed travel conversions
- only add missing PQ families that are actually wanted in this version; currently that means GuildPQ and Orbis PQ, not every LeaderMS PQ/event family

## Formal State-Machine Conversion Rules

### Rule 1: Convert `status` branches into explicit linear checkpoints

LeaderMS `status` is not ported literally. Each `status == N` branch becomes one explicit checkpoint in script order.

Safe translation rule:

- every `sendX(...)` that waits for input becomes one top-level Valhalla statement
- every side effect that originally happened only after the next `action(...)` callback must remain after the waiting call, never before it
- do not move warps, rewards, removals, quest mutations, or job changes above the prompt that gated them

LeaderMS pattern:

```js
var status = -1;
function start() { action(1, 0, 0); }
function action(mode, type, selection) {
  if (mode == 1) status++; else status--;
  if (status == 0) cm.sendNext("...");
  else if (status == 1) cm.sendYesNo("...");
  else if (status == 2) cm.warp(100000000);
}
```

Valhalla pattern:

```js
npc.sendNext("...");
if (npc.sendYesNo("...")) {
    plr.warp(100000000);
}
```

### Rule 2: Nested menus become nested return-value branches

LeaderMS stores `selection` from `action(...)` and later branches on both `status` and `selection`.

Valhalla translation rule:

- call the first menu, branch on its returned value
- inside each branch, call the next menu only when that branch was actually selected
- if a later menu is only reachable from one branch in LeaderMS, it must also be physically nested in that branch in Valhalla

Valhalla example:

```js
var first = npc.sendMenu("Choose.", "Explain", "Travel");
if (first === 0) {
    var second = npc.sendMenu("What should I explain?", "A", "B");
    // explanation branch only
} else if (first === 1) {
    var destination = npc.sendMenu("Where to?", "Town 1", "Town 2");
    // travel branch only
}
```

### Rule 3: `mode == 0` backout and cancel handling must be translated by prompt type

#### 3a. `sendYesNo` / confirmation prompts

LeaderMS pattern:

- `mode == 0` usually means explicit “No” or cancel
- source often sends alternate dialogue on the false path

Valhalla rule:

- translate to `if (!npc.sendYesNo(...)) { ... } else { ... }`
- put the original “declined” branch in the false path

#### 3b. `sendSimple` / menu prompts

LeaderMS pattern:

- `mode == 0` usually disposes or jumps back to an earlier branch

Valhalla rule:

- if LeaderMS cancel simply disposes, rely on normal conversation termination and do not fabricate extra dialogue
- if LeaderMS cancel jumps back to an earlier logical branch, rewrite as an explicit loop or helper function in JS
- if the source reopens the same choice repeatedly until valid, use a local loop

#### 3c. `sendNext` / `sendOk` prompts

Important limitation:

- in Valhalla's current runtime, closing a plain next/ok prompt terminates the session before the script re-runs
- therefore a LeaderMS branch that reacts to `mode == 0` after `sendNext` cannot be reproduced exactly without extra runtime support

Conversion rule:

- if cancel on a next/ok prompt only changes farewell text, accept the silent-close difference and document it
- if cancel controls important gameplay state, stop and treat it as a runtime gap or manual-review case

### Rule 4: Repeated menu loops must be written explicitly

LeaderMS often re-enters a prior branch by setting `status` backward or by sending the same menu again after a rejection.

Valhalla rule:

- use `while` loops only when the original script truly repeats until valid input
- otherwise prefer one-shot branching with an explicit “come back later” exit

Example pattern:

```js
var picked = -1;
while (picked === -1) {
    var sel = npc.sendMenu("Choose.", "A", "B");
    if (sel < 0) {
        break;
    }
    if (valid(sel)) {
        picked = sel;
    } else {
        npc.sendOk("That option is not available right now.");
    }
}
```

### Rule 5: Re-entry into prior logical branches should use functions or local blocks, not fake `status`

If the source has a tree like:

- explanation menu
- choose destination
- cancel returns to explanation menu

then write small local functions or nested loops in Valhalla. Do not recreate `status = 1`, `status = 2`, etc.

### Rule 6: Side effects remain server-authoritative and last

LeaderMS often does cleanup in the final `status` branch right before `warp()`.

Valhalla should keep the same sequencing explicitly:

```js
if (!plr.isGM()) {
    plr.removeAll(4001015);
}
plr.warp(211042300);
```

## Common Conversion Recipes

### Single-step dialogue

- LeaderMS: `start() -> sendOk -> dispose()`
- Valhalla: `npc.sendOk("...")`

### Multi-step `status` conversation

- flatten into normal JS order
- convert `mode == 0` exits into explicit false branches or loops based on the prompt type
- do not preserve numeric `status` values unless absolutely necessary for readability
- if the source relies on `mode == 0` after `sendNext` or `sendOk`, document whether the target can preserve that behavior exactly

### Selection and menu branching

- `sendSimple` maps to `sendSelection` when keeping raw menu text
- `sendSimple` maps to `sendMenu` when rebuilding a sequential list idiomatically
- use returned integers directly instead of storing a global `selection`

### Item exchange

- check mesos and source items first
- check inventory capacity before deducting resources
- use `plr.gainItem(...)` for both additions and removals
- use `plr.canHoldAll(...)` when several result items may be granted together

### Mesos exchange

- use `plr.getMesos()` for checks
- use `plr.gainMesos(-amount)` for payment
- keep source-side pricing exactly unless target is already known-correct and intentionally custom

### Quest start, progress, and complete

- map quest enum checks to normalized `0/1/2`
- prefer `plr.questStarted(...)` and `plr.questCompleted(...)` when they make the script clearer
- preserve ordering of quest mutation, item mutation, and dialogue

### Taxi and warp NPCs

- preserve destination lists, beginner discounts, ticket checks, and cancel text
- many target scripts already have the same skeleton but are missing one destination or branch
- default taxi conversion rule: prefer `npc.sendSelection(...)` when preserving source menu text, and `plr.getMesos()` / `plr.gainMesos(...)` / `plr.warp(...)` for the outcome branch

### Job advancement

- keep for manual review unless target already matches closely
- source often depends on direct base stat/AP mutation and custom quest chains
- avoid partial conversion without runtime support for the full advancement path

### Open shop / list-based interactions

- LeaderMS shop IDs are not directly portable
- prefer existing Valhalla `npc.sendShop([...])` lists where already present
- if porting a shop-only NPC, derive the actual goods list instead of inventing one

## Unsupported Or Missing Valhalla Features

Remaining confirmed gaps after this tightening pass:

- `npc.sendAcceptDecline(...)`
- `openNpc(id)` style script-to-script chaining
- direct base-stat reset or mutation helpers needed by first-job parity
- `createPlayerNPC()` parity
- exact LeaderMS `getText()` model does not exist; Valhalla uses return-value input prompts instead
- plain next/ok cancel callbacks are not script-visible in the current runtime

## Runtime Additions Needed For Parity

Implemented in `channel/script.go` during the pilot and tightening passes:

- `plr.getMesos()`
- `plr.gainMesos(amount)`
- `plr.gainItem(id, amount)`
- `plr.haveItem(id, quantity)`
- `plr.canHoldAll([[id, qty], ...])`
- `plr.questStarted(id)`
- `plr.questCompleted(id)`
- `plr.questNotStarted(id)`
- `plr.saveLocation(slot)`
- `plr.getSavedLocation(slot)`
- `plr.clearSavedLocation(slot)`
- `plr.getEventProperty(key)`
- `plr.setEventProperty(key, value)`
- `plr.finishEvent()`
- `plr.leaveEvent()`
- `plr.countMonster()`
- `plr.spawnMonster(id, x, y)`
- `plr.gainGuildPoints(points)`
- `plr.removeAll(id)`
- `plr.canHold(id, amount)`
- `plr.getEquipInventoryFreeSlot()`
- `plr.getUseInventoryFreeSlot()`
- `plr.getSetupInventoryFreeSlot()`
- `plr.getEtcInventoryFreeSlot()`
- `plr.getCashInventoryFreeSlot()`
- `plr.isGM()`

Still likely needed for later batches:

- `npc.sendAcceptDecline(...)`
- base stat/AP reset helpers for job-advancer scripts
- optional `npc.openNpc(...)` if chained NPC flows prove common enough
- prompt-cancel callbacks for plain `sendNext` / `sendOk` flows if exact cancel text parity becomes important
- a true NPC-side `EventManager` wrapper if we decide to port generic timed transport loops directly from LeaderMS event scripts instead of reusing existing Go schedulers

## Canonical Examples

### Canonical Taxi / Warp Example

- source: `C:\Users\daddy\desktop\MapleDev\LeaderMS-English\scripts\npc\1012000.js`
- target: `scripts/npc/1012000.js`
- why canonical: standard beginner-discount town taxi with one destination menu and one yes/no confirmation
- key translation rules used:
  - preserved full destination list and prices
  - translated `sendSimple` destination choice into one Valhalla menu checkpoint
  - translated the final confirmation into one yes/no branch
  - kept mesos check and warp side effects after confirmation only
- runtime helpers required:
- `plr.getMesos()`
- `plr.gainMesos(...)`
- `plr.warp(...)`

### Canonical Item / Reward / Exchange Example

- source: `C:\Users\daddy\desktop\MapleDev\LeaderMS-English\scripts\npc\2040000.js`
- target: `scripts/npc/2040000.js`
- why canonical: simple purchase flow with mesos validation and inventory-capacity validation before item grant
- key translation rules used:
  - confirmation prompt stays a yes/no branch
  - failure path merges mesos-full and inventory-full exactly like the source message
  - item grant happens only after capacity preflight succeeds
- runtime helpers required:
  - `plr.getMesos()`
  - `plr.canHold(...)`
  - `plr.gainItem(...)`
  - `plr.gainMesos(...)`

### Canonical Quest / Event Exit Example

- source: `C:\Users\daddy\desktop\MapleDev\LeaderMS-English\scripts\npc\2030011.js`
- target: `scripts/npc/2030011.js`
- why canonical: cleanup-before-warp exit flow for PQ or event content, including GM exemption and multi-item removal
- key translation rules used:
  - preserved the cleanup-before-warp ordering
  - preserved the GM exemption
  - documented the remaining limitation: the source sends a farewell message if the player cancels the initial next prompt, while Valhalla currently closes silently on next-prompt cancel
- runtime helpers required:
- `plr.isGM()`
- `plr.removeAll(...)`
- `plr.warp(...)`

### Canonical Travel Return Example

- source: `C:\Users\daddy\desktop\MapleDev\LeaderMS-English\scripts\npc\1002002.js` and `C:\Users\daddy\desktop\MapleDev\LeaderMS-English\scripts\npc\1081001.js`
- target: `scripts/npc/1002002.js` and `scripts/npc/1081001.js`
- why canonical: travel-out plus travel-back pair that requires preserving a dedicated return destination instead of relying on raw previous-map behavior
- key translation rules used:
  - save the origin map before outbound warp
  - use the saved return map on the return NPC with a safe fallback
  - keep ticket possession semantics identical to the source
- runtime helpers required:
  - `plr.saveLocation(...)`
  - `plr.getSavedLocation(...)`
  - `plr.getMesos()`
  - `plr.gainMesos(...)`
  - `plr.haveItem(...)`
  - `plr.warp(...)`

## Known Ambiguities And Edge Cases

- same NPC ID does not guarantee the same behavior; some Valhalla scripts are custom rewrites
- the inventory classifier is heuristic, not authoritative
- known manual-review examples:
  - `1002000`: Valhalla preserves the taxi half but not the long Victoria town explanation branch
  - `9200100`: Valhalla preserves lens usage but omits the source coupon-purchase branch
  - `1022000`: first and second warrior advancement logic diverges materially and needs runtime support review
- LeaderMS frequently mixes engine helpers with direct `getPlayer()` access; those cases need extra scrutiny during conversion

## Proposed Migration Plan By Batch

1. Verify existing target-only matches for storage and shop scripts.
2. Convert taxi and warp scripts where the runtime is already sufficient.
3. Convert small reward and exchange flows using the canonical item helpers.
4. Convert straightforward quest scripts that only need normalized quest checks, item checks, and warps.
5. Revisit appearance scripts with coupon-purchase branches.
6. Tackle job advancers and PQ or event NPCs only after the remaining runtime gaps are resolved.
