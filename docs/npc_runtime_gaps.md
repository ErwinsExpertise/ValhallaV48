# NPC Runtime Gaps

## Added For Conversion Consistency

Implemented in `channel/script.go` during the pilot and tightening passes:

- `plr.getMesos()`
- `plr.gainMesos(amount)`
- `plr.gainItem(id, amount)`
- `plr.haveItem(id, quantity)`
- `plr.canHold(id, amount)`
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
- `plr.getEquipInventoryFreeSlot()`
- `plr.getUseInventoryFreeSlot()`
- `plr.getSetupInventoryFreeSlot()`
- `plr.getEtcInventoryFreeSlot()`
- `plr.getCashInventoryFreeSlot()`
- `plr.isGM()`

Why these were added:

- LeaderMS-style scripts depend heavily on canonical item and mesos helpers
- many exchange flows need consistent positive and negative item semantics
- many purchase or reward flows need inventory-capacity checks before side effects
- quest scripts read better when status checks are wrapped in boolean helpers
- some source scripts preserve GM-specific behavior on cleanup or exit paths
- event-property helpers are wired to field-instance properties to stay consistent with existing Valhalla PQ/event scripts

## Canonical Helper Contracts

### `plr.gainItem(id, qty)`

- positive quantity: add item if capacity allows
- zero quantity: success no-op
- negative quantity: remove exactly `abs(qty)` items by ID
- returns `false` on capacity failure or insufficient items

This is the canonical target-side equivalent for LeaderMS `cm.gainItem(...)`.

Compatibility note:

- existing target scripts that call `plr.giveItem(id, -1)` now behave as expected, but new conversions should still prefer `plr.gainItem(...)`

### `plr.gainMesos(amount)`

- direct alias for mesos delta application
- positive adds mesos
- negative removes mesos

This is the canonical target-side equivalent for LeaderMS `cm.gainMeso(...)`.

### Quest wrappers

- `plr.getQuestStatus(id)` remains the normalized primitive
- `plr.questStarted(id)`, `plr.questCompleted(id)`, and `plr.questNotStarted(id)` are the preferred readability helpers for new conversions

## Remaining Gaps

### High priority

- `npc.sendAcceptDecline(...)`
  - needed by job-advancers and some quest-offer scripts
  - source usage confirmed in `1022000.js` and other LeaderMS flows

- base-stat and AP reset helpers
  - needed for accurate first-job advancement parity
  - source usage confirmed in `1022000.js`

### Medium priority

- `npc.openNpc(id)` style chaining
  - some LeaderMS flows jump into another NPC script instead of staying in one file

- generic NPC-side `EventManager` wrapper
  - some LeaderMS NPCs call `getEventManager(...)` directly for timed transport systems or PQ startup
  - Valhalla currently covers some of that behavior through dedicated Go schedulers (`boats.go`) and some through existing event scripts (`scripts/event`)
  - if we need full parity for the remaining event-backed travel NPCs, exposing a minimal `getEventManager` wrapper may be cleaner than hand-adapting every script

- richer exact text-input parity
  - Valhalla already supports input prompts, but not the exact `sendGetText` plus `getText()` model

### Important behavioral limitation

- plain next and ok prompt cancel callbacks are not script-visible
  - if a player closes a plain `sendNext` or `sendOk` prompt, the conversation ends before the script re-runs
  - this blocks exact reproduction of LeaderMS branches that react to `mode == 0` after those prompt types
  - current known impact: `2030011` can preserve the main exit flow, but not the source farewell text on cancel

### Deferred and special-case

- `createPlayerNPC()`
  - niche feature used in LeaderMS max-level job master scripts
  - should not be added unless those exact scripts are being ported

## Guidance

- do not add runtime helpers preemptively for all possible LeaderMS APIs
- add the smallest script-exposed method that matches an existing Valhalla server primitive
- prefer one canonical helper per concept so later conversions stay uniform
- record every new helper here and in `docs/npc_migration_log.md`

Saved-location note:

- the new saved-location helpers are currently in-memory session state, which is sufficient for the current taxi and warp flows but not yet durable across reconnects

Event-property note:

- `plr.getEventProperty` and `plr.setEventProperty` intentionally use the player's current field-instance property map
- this matches the existing KPQ/LPQ event scripting style in Valhalla and avoids splitting event state across two systems
