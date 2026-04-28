# PQ Remaining Work

## Scope

This document tracks what is still unfinished for:

- GuildPQ
- Orbis PQ

It only covers implementation gaps, not broader migration commentary.

## GuildPQ

### Already in place

- guild-scoped event startup and join flow
- waiting room event shell
- timed entry lock
- earring enforcement timer
- exit NPC
- info / notice NPCs
- finisher / reward NPC
- stage 1 gatekeeper scaffold
- stage 3 fountain puzzle
- stage 4 soul / ghost gate clear
- some reactor scripts for NPC spawn / GP reward / spear gate
- PQ item cleanup on leave / timeout / finish

### Still missing or incomplete

#### Event structure

- fuller map-by-map stage orchestration beyond the waiting-room shell
- more explicit stage progression state in the event script instead of relying on scattered NPC behavior
- audit of all stage clear properties to ensure naming and usage are consistent

#### Stage 1

- confirm reactor-name/object-id mapping against the real stage-1 statue layout in `v48` NX
- verify stage-1 pattern reveal and guess collection from real statue clicks
- verify phase reset and progression across all three rounds
- confirm `statuegate` reactor behavior on full clear

#### Stage 2

- full spear-route logic still needs verification and likely more implementation
- confirm all spear / portal / trap interactions against source behavior
- check if additional portal scripts are required for stage-2 routing

#### Stage 3

- fountain puzzle exists, but needs gameplay verification
- verify item-area definitions against the real map areas
- verify all punishment mob spawns and attempt reset rules
- confirm `watergate` opening behavior and follow-up routing

#### Stage 4

- clothing/soul path needs fuller validation
- verify all required clothing item sources and ordering expectations
- confirm Sharen III soul spawn and ghost gate flow from source reactor/NPC interactions
- verify Ergoth boss room trigger and reward chain

#### Stage 5 / bonus / reward flow

- verify bonus map entry / exit routing
- verify bonus box reactor behavior that is still source-defined but not fully audited in Valhalla
- verify GP reward pacing and final completion timing

#### Reactors

- audit all `9208xxx` source reactor scripts and confirm which ones still need porting
- confirm which reactors are already automatically handled by core reactor data and which truly need scripts
- remove any remaining duplication between reactor data behavior and explicit scripts

#### Portals

- full portal-script audit for GuildPQ is still incomplete
- verify every staged gate / waiting / return portal from source against Valhalla's current portal scripts

## Orbis PQ

### Already in place

- party-quest event shell
- entry NPC (`2013000`)
- main controller NPC (`2013001`) partial implementation
- reward / finish NPC (`2013002`)
- portal-script infrastructure added to runtime
- Orbis PQ room and return portal scripts added
- stage-1 Chamberlain popup restored
- stage-1 reactor-driven Chamberlain spawn restored
- event cleanup for PQ items on leave / timeout / finish
- GM solo-entry testing support added

### Still missing or incomplete

#### Event structure

- final consistency pass to keep OPQ as close as possible to KPQ/LPQ event style
- confirm event callbacks used by OPQ are the minimum needed and not carrying leftover experimental behavior

#### Stage 1

- verify drop-to-reactor behavior on the real `v48` map under normal gameplay
- confirm Chamberlain spawn timing and location against source behavior
- confirm stage clear only happens after the intended interaction path, not too early or too late

#### Center Tower routing

- verify all center-tower room routing and return flow
- confirm leader/non-leader behavior for room selection and re-entry
- confirm stage clear properties align with portal usage exactly

#### Stage 3 walkway

- implemented hand-in path exists, but needs gameplay verification
- verify required item id/count and party reward pacing

#### Stage 4 storage

- implemented hand-in path exists, but needs gameplay verification
- confirm statue piece acquisition path from monsters/reactors is complete

#### Stage 5 lobby music

- current implementation is only partial
- verify the real music/day mechanic against the `v48` reactor state data
- confirm wrong-CD handling and no-backup behavior

#### Stage 6 sealed room

- still largely placeholder
- needs full combination logic implementation
- should likely use map areas and/or reactor/portal state rather than NPC-only text

#### Stage 7 lounge

- implemented hand-in exists, but needs gameplay verification
- confirm item count, reward pacing, and follow-up routing

#### Stage 8 on-the-way-up

- currently simplified
- verify whether this free piece shortcut is acceptable for this version or should be replaced with a fuller source behavior

#### Garden / darkness / bonus side

- verify garden routing and any missing interactions
- verify room of darkness flow
- verify bonus map timing and reward transition path

#### Portals

- portal scripts exist, but still need a map-by-map audit against source portal behavior
- verify blocked-backtracking behavior after clear on every OPQ side room

#### Reactors

- only the stage-1 Chamberlain reactor is explicitly covered right now
- remaining OPQ reactor-side behavior needs a full audit against source reactor scripts and `v48` NX data

## Recommended Next Order

1. Finish GuildPQ stage mechanics map by map
2. Finish Orbis PQ sealed room and lobby/music mechanics
3. Audit all GPQ/OPQ portal scripts against source and `v48` NX
4. Audit all GPQ/OPQ reactor scripts against source and `v48` NX
5. Do a gameplay verification pass on both PQs end to end
