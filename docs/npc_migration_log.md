# NPC Migration Log

## 2026-04-21

### Framework tightening

- expanded the conversion reference with formal state-machine translation rules instead of relying on generic "blocking style" guidance
- documented exact translation rules for nested menus, repeated menus, yes-no branches, backout behavior, and re-entry into prior logical branches
- documented the current runtime limitation around cancel handling for plain next and ok prompts

### API tightening

- standardized item translation on `plr.gainItem(...)` as the canonical target equivalent of `cm.gainItem(...)`
- documented exact rules for when `cm.sendSimple(...)` maps to `npc.sendSelection(...)` versus `npc.sendMenu(...)`
- documented and normalized quest status semantics to `0 = not started`, `1 = started`, `2 = completed`

### Runtime additions

- exposed `plr.getMesos()`
- exposed `plr.gainMesos(amount)`
- exposed `plr.gainItem(id, amount)`
- exposed `plr.haveItem(id, quantity)`
- exposed `plr.canHoldAll([[id, qty], ...])`
- exposed `plr.questStarted(id)`
- exposed `plr.questCompleted(id)`
- exposed `plr.questNotStarted(id)`
- exposed `plr.saveLocation(slot)`
- exposed `plr.getSavedLocation(slot)`
- exposed `plr.clearSavedLocation(slot)`
- exposed `plr.getEventProperty(key)`
- exposed `plr.setEventProperty(key, value)`
- exposed `plr.finishEvent()`
- exposed `plr.leaveEvent()`
- exposed `plr.countMonster()`
- exposed `plr.spawnMonster(id, x, y)`
- exposed `plr.gainGuildPoints(points)`
- retained previously added helpers: `plr.removeAll`, inventory free-slot helpers, `plr.canHold`, `plr.isGM`

### Tooling upgrades

- refined `tools/compare_npc_inventory.go` categories from the coarse shared-script bucket into:
  - `behaviorally_equivalent_but_structurally_different`
  - `behaviorally_different`
- added similarity scoring and a substantial-difference flag
- added manual-review blocker strings and reason codes to inventory entries
- kept manual overrides for known false positives and reviewed pilot cases

### Pilot re-check

- `1012000`
  - rechecked beginner discount, mesos check, and confirmation flow
  - aligned script to canonical helper names `getMesos` and `gainMesos`
  - result: acceptable canonical taxi template

- `2030011`
  - corrected the pilot away from a yes-no prompt back to a plain next prompt so the main success path matches the source better
  - documented the remaining cancel-path mismatch caused by current runtime limitations

- `2040000`
  - rechecked cost, mesos check, capacity check, and sale flow
  - aligned script to canonical helper names `getMesos`, `canHold`, `gainItem`, and `gainMesos`
  - result: acceptable canonical purchase and exchange template

### Inventory snapshot after tightening

- 470 unique NPC IDs
- 72 `behaviorally_equivalent_but_structurally_different`
- 56 `behaviorally_different`
- 152 `exists_only_in_leaderms`
- 105 `exists_only_in_valhalla`
- 85 `ambiguous_manual_review`

### Next constrained batch

Taxi and warp family only.

Selected for next batch review or conversion:

1. Verify-equivalent candidates
   - `22000`
   - `1022001`
   - `1032000`
   - `1052016`
2. Patch or rewrite candidates
   - `1002000`
   - `1002002`
   - `1002004`
   - `1032005`
   - `1061100`
3. Manual-review taxi or warp outliers
   - `11100`
   - `1052012`

## 2026-04-21 Taxi Batch Continued

### Converted taxi and warp scripts

- `22000`
  - restored recommendation-letter handling, level gate, and exact 150 mesos departure flow
- `1022001`
  - restored the missing Nautilus destination and aligned Henesys-side pricing and failure text
- `1032000`
  - restored Ellinia destination list, pricing, and insufficient-mesos text
- `1052016`
  - restored Kerning destination ordering, Nautilus destination, and source-side pricing
- `1002002`
  - restored the Florina 1500 mesos fare and the non-consuming VIP-ticket path
- `1002000`
  - restored the long Victoria Island explanation branch alongside the travel branch
- `1081001`
  - restored Florina return-map handling with a dedicated saved-location slot and Lith Harbor fallback
- `1002004`
  - aligned Lith Harbor VIP cab wording and beginner and non-beginner fee handling
- `1032005`
  - aligned Ellinia VIP cab wording and beginner and non-beginner fee handling
- `1061100`
  - corrected Sleepywood hotel destination maps and aligned room-selection wording with the source

## 2026-04-21 Travel Batch Continued Again

### Additional converted travel scripts

- `2100`
  - ported the beginner training-camp entry and skip-training branch
- `2101`
  - ported the training-camp exit NPC
- `1022101`
  - ported Happy Village entry and stores the return map in a dedicated saved-location slot
- `2002000`
  - ported Happy Village exit using the saved return-map slot
- `9000020`
  - ported Mushroom Shrine world-tour travel and return using the `WORLDTOUR` saved-location slot
- `2010005`
  - ported the Orbis-to-Florina travel guide using the canonical Florina saved-location pattern
- `2040048`
  - ported the Ludibrium-to-Florina travel guide using the canonical Florina saved-location pattern
- `2081009`
  - ported the direct Moose exit warp to the correct destination portal
- `9201010`
  - ported the simple leave-room menu warp
- `9201057`
  - ported the Kerning and NLC subway ticket purchase and early-exit flow
- `9270038`
  - ported the Singapore Airport return ticket purchase and boarding flow
- `9270041`
  - ported the Singapore Airport outbound ticket purchase and boarding flow

### Additional patched travel scripts

- `1032004`
  - restored the simple Ellinia return prompt and warp outcome
- `1032008`
  - aligned the Ellinia-to-Orbis boarding usher to the source-side ticket check and direct warp flow
- `1032009`
  - aligned the pre-takeoff waiting-room exit script to the source wording and return warp
- `2012001`
  - aligned the Orbis-to-Ellinia boarding usher to the simpler source-side immediate boarding flow
- `2012013`
  - aligned the Orbis-to-Ludibrium boarding usher to the simpler source-side immediate boarding flow
- `2041000`
  - aligned the Ludibrium-to-Orbis boarding usher to the source-side ticket check and direct warp flow

## 2026-04-21 Travel Batch Continued Further

### Additional converted direct-warp scripts

- `2200002`
  - ported the direct LPQ trap warp box
- `9000002`
  - ported the mission-top completion exit warp
- `9000010`
  - ported the direct event-room exit warp
- `9101001`
  - ported the Mushroom Town training-camp exit sequence
- `9201049`
  - ported the direct event exit warp
- `2082003`
  - ported the two-way direct map toggle warp

### Additional patched existing scripts

- `1061010`
  - restored the source-side direct warp behavior for this NPC id

### Remaining travel scope

- 68 source-backed taxi and warp entries still remain unresolved
- the remaining set is increasingly weighted toward event-specific, ticket-gated, and branching travel flows

## 2026-04-21 Source-Only Boarding Batch

### Additional converted source-only travel scripts

- `2101018`
  - ported the Ariant PQ entry NPC and stores the Ariant PQ return location
- `2082001`
  - ported the Leafre-to-Orbis boarding usher
- `2012021`
  - ported the Orbis-to-Leafre boarding usher
- `2012025`
  - ported the Orbis-to-Ariant boarding usher
- `2102000`
  - ported the Ariant-to-Orbis boarding usher

### Constraint reaffirmed

- preserve existing Valhalla instance or event logic when an NPC already has complex target-specific behavior
- use LeaderMS primarily for source-only scripts or for narrow corrections on otherwise matching target scripts

## 2026-04-21 Complex-Safe Travel Batch

### Additional converted source-only travel scripts

- `9103002`
  - ported the reward-and-exit NPC with random consumable reward logic
- `1072008`
  - ported the pirate test crystal check and Nautilus return flow
- `9201006`
  - ported the wedding-room question and exit menu flow
- `2101016`
  - ported the APQ jewelry hand-in and exit flow using current runtime item and EXP helpers

### Deferred for runtime-aware pass

- `9040010`
  - still depends on event-instance completion and guild-point behavior that should be handled in a dedicated runtime-gap pass instead of a partial port

## 2026-04-21 Broader Source-Only Travel Pass

### Additional converted source-only travel scripts

- `2012002`
  - ported the leave-the-boat NPC with map-dependent return warp
- `2101013`
  - ported the simple return-to-Ellinia NPC
- `9060000`
  - ported the item-gated zoo return NPC
- `9120200`
  - ported the simple hideout return NPC
- `9270017`
  - ported the airport waiting-room exit pilot based on the source's intended return destination
- `9270018`
  - ported the airport waiting-room exit pilot

### Batch note

- this pass intentionally favored clearing multiple source-only travel IDs at once
- target-overlap instance and event NPCs are still being treated conservatively after the earlier regressions

## 2026-04-21 Event-Backed Travel Mapping

### Event-runtime alignment

- reviewed LeaderMS event scripts in `scripts/event`, especially `AirPlane.js`, `Subway.js`, and `Cabin.js`
- reviewed current Valhalla event references in `scripts/event/kerning_pq.js`, `scripts/event/ludibrium_pq.js`, and `channel/boats.go`
- documented that Valhalla already has strong reference implementations for KPQ, LPQ, and the Ellinia/Orbis/Ludibrium boat system

### New event-backed conversions

- `9201047`
  - ported using `plr.isLeader()`, `plr.countMonster()`, and `plr.spawnMonster(...)`
- `9040010`
  - ported using event property helpers and `plr.finishEvent()` / `plr.gainGuildPoints(...)`

### Event-model conclusion

- remaining event-backed travel NPCs should be mapped against existing Valhalla event scripts and `boats.go` first
- a generic `EventManager` wrapper may still be useful later, but it is no longer required for the first event-backed NPC conversions
- event-property helpers were aligned to field-instance properties to avoid crossing wires with the established Valhalla PQ/event pattern
- PQ scope was narrowed again: add GuildPQ and Orbis PQ, but do not carry over unsupported PQ families like Henesys PQ for this version

### Deferred within taxi and warp family

- `11100`
  - still manual review due source-target mismatch
- `1052012`
  - still manual review due family mismatch and branching complexity

### Next taxi and warp review set

1. Review the remaining taxi or warp manual-review outliers
2. Verify whether any remaining taxi scripts already became equivalent after the canonical helper cleanup
3. Revisit saved-location persistence only if reconnect behavior becomes important for travel flows
