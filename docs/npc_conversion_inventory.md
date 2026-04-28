# NPC Conversion Inventory

Generated source of truth:

- `tools/npc_conversion_inventory.json`
- generator: `go run .\tools\compare_npc_inventory.go`

## Current Snapshot

LeaderMS NPC scripts scanned: 370

Valhalla numeric NPC scripts scanned: 285

Valhalla total NPC scripts on disk: 286 including `default.js`

Combined unique NPC IDs: 470

Category summary:

- `behaviorally_equivalent_but_structurally_different`: 72
- `behaviorally_different`: 56
- `exists_only_in_leaderms`: 152
- `exists_only_in_valhalla`: 105
- `ambiguous_manual_review`: 85

Family summary from the generated inventory groups:

- appearance
- branching_dialogue
- event_minigame
- exchange_reward
- job_advancement
- quest
- shop
- simple_dialogue
- storage
- taxi_warp
- warp

## Meaning Of The Refined Categories

### `behaviorally_equivalent_but_structurally_different`

- both repos appear to implement the same gameplay outcome
- script structure differs because Valhalla is already written idiomatically
- default action: keep Valhalla and verify, do not rewrite just to match source shape

### `behaviorally_different`

- both repos have the NPC ID, but behavior differs enough to justify patching or rewriting
- default action: patch or rewrite the Valhalla script to restore source behavior, unless the target behavior is intentionally preferred

### `ambiguous_manual_review`

- automation found a concrete blocker that makes blind conversion risky
- each such entry now records a blocker and reason codes in `tools/npc_conversion_inventory.json`

Common reason codes:

- `behavior_mismatch`
- `complex_branching`
- `event_specific`
- `source_target_mismatch`
- `runtime_gap_send_accept_decline`
- `runtime_gap_base_stat_mutation`
- `runtime_gap_create_player_npc`
- `runtime_gap_open_npc`
- `runtime_gap_text_input`

## Pilot Status

Verified unchanged:

- `1002005` storage NPC: existing Valhalla `npc.sendStorage(1002005)` is the correct target-side abstraction for the LeaderMS direct storage call

Converted in pilot:

- `1012000`: Henesys taxi updated to restore the missing Nautilus destination and to align pricing and confirmation flow with LeaderMS
- `2030011`: PQ exit NPC updated to restore PQ item cleanup before exit warp; remaining mismatch is documented because next-prompt cancel callbacks are not script-visible in the current runtime
- `2040000`: Ludibrium ticket seller updated to restore the original 6000 mesos cost and proper inventory-capacity validation

Converted in taxi/warp batch:

- `22000`
- `2100`
- `2101`
- `1022001`
- `1022101`
- `1032000`
- `1032004`
- `1032008`
- `1032009`
- `1052016`
- `1002002`
- `1002000`
- `1081001`
- `2002000`
- `9000020`
- `2010005`
- `2012001`
- `2012013`
- `1002004`
- `1032005`
- `1061100`
- `2040048`
- `2041000`
- `2081009`
- `9201010`
- `9201057`
- `9270038`
- `9270041`
- `2200002`
- `9000002`
- `9000010`
- `9101001`
- `9201049`
- `2082003`
- `2101018`
- `2082001`
- `2012021`
- `2012025`
- `2102000`
- `9103002`
- `1072008`
- `9201006`
- `2101016`
- `2012002`
- `2101013`
- `9060000`
- `9120200`
- `9270017`
- `9270018`

Pilot re-check outcome:

- `1012000`: behavior now matches the source's taxi branch closely enough to use as the canonical taxi template
- `2030011`: main success path matches, but source-side cancel-after-next farewell text still cannot be reproduced exactly without extra runtime support
- `2040000`: price, capacity check, mesos check, and sale flow now match the source behavior

## Batch Guidance

Immediate execution order remains constrained to taxi and warp work only.

Recommended next taxi and warp review set:

1. Remaining rewrite-or-patch set
   - none currently required from the already-selected core set
2. Manual-review taxi or warp outliers
   - `11100`
   - `1052012`
3. Follow-up checks
    - confirm whether any additional Florina-side travel NPCs should use the same saved-location helpers
    - continue through the remaining unresolved travel list in `tools/npc_conversion_inventory.json`
    - continue through the remaining source-only boarding and travel ushers before touching more manual-review target-overlap cases
    - reserve event-instance and guild-finisher travel scripts for a dedicated runtime-gap pass where needed
    - the remaining source-only travel set is now increasingly event-instance heavy

## Notes

- the JSON inventory is intentionally machine-readable and is the canonical per-NPC list
- this markdown file is the human summary layer, not a duplicate of all 470 entries
- manual review should now always include both a blocker string and reason codes in the JSON output
- when overrides are discovered, update the generator instead of hand-editing the generated JSON
