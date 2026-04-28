# PQ Implementation Guide

## Purpose

This document defines the preferred way to add Party Quests in Valhalla.

Use this as the default pattern for future PQ work.

Goals:

- keep PQs consistent with KPQ and LPQ
- keep event state simple and centralized
- avoid one-off hacks in unrelated runtime paths
- prefer event scripts, portal scripts, and small NPC scripts over giant all-in-one logic blobs
- use field instance properties as the source of truth for stage state

## Core Rules

### 1. Event script owns the PQ lifecycle

The event script should own:

- map reset
- property reset
- initial warp into the PQ
- countdown start
- portal gating behavior
- timeout behavior
- finish behavior
- leave behavior
- item cleanup

Do not spread lifecycle logic across random NPCs if it can live in the event script.

### 2. NPCs should handle stage checks, not event lifecycle

NPCs are for:

- entry/start prompts
- stage explanations
- leader-only validation
- item hand-ins
- stage completion checks
- reward/finish interactions

NPCs should not become giant state machines for the whole PQ.

### 3. Use field instance properties for stage state

Preferred state location:

- `map.properties()["clear"]`
- `map.properties()["stageXclear"]`
- `map.properties()["leader"]`
- `map.properties()["entryTimestamp"]`
- stage-specific counters/combinations

Do not invent a second state store when field properties are sufficient.

### 4. Portal scripts should gate movement

If a PQ room transition depends on stage clear state, use a portal script.

Do not force all movement through NPCs when the source expects portal-gated progression.

### 5. Reactor behavior should live in reactor scripts

If a reactor:

- spawns an NPC
- spawns a monster
- drops custom items
- opens another gate/reactor

put that behavior in `scripts/reactor/<reactorId>.js`.

Do not hardcode map-specific reactor behavior inside generic inventory or drop logic.

### 6. Cleanup belongs in one place

Preferred cleanup location:

- event `timeout()`
- event `playerLeaveEvent()`
- event `finish()`

NPC exit maps can be simple, but should not duplicate complicated event cleanup unless there is a specific reason.

## Standard File Layout

For a PQ, prefer this structure:

- `scripts/event/<pq_name>.js`
- `scripts/npc/<entry_npc>.js`
- `scripts/npc/<stage_npc>.js`
- `scripts/portal/<portal_name>.js`
- `scripts/reactor/<reactor_id>.js`

Not every PQ needs every type, but this is the default shape.

## Event Script Pattern

Preferred event script structure:

```js
var maps = [...];
var entryMapID = ...;
var exitMapID = ...;
var rewardMapID = ...;
var pqItems = [...];

function clearPQItems(plr) {
    for (var i = 0; i < pqItems.length; i++) {
        plr.removeAll(pqItems[i]);
    }
}

function start() {
    ctrl.setDuration("30m");

    for (var i = 0; i < maps.length; i++) {
        var field = ctrl.getMap(maps[i]);
        field.reset();
        field.clearProperties();
    }

    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (var i = 0; i < players.length; i++) {
        players[i].warp(entryMapID);
        players[i].showCountdown(time);
    }
}

function beforePortal(plr, src, dst) {
    var props = src.properties();
    if (props["clear"]) {
        return true;
    }
    plr.sendMessage("Cannot use portal at the moment");
    return false;
}

function afterPortal(plr, dst) {
    plr.showCountdown(ctrl.remainingTime());
    if (dst.properties()["clear"]) {
        plr.portalEffect("gate");
    }
}

function timeout(plr) {
    clearPQItems(plr);
    plr.warp(exitMapID);
}

function finish() {
    var players = ctrl.players();
    for (var i = 0; i < players.length; i++) {
        clearPQItems(players[i]);
        players[i].warp(rewardMapID);
    }
}

function playerLeaveEvent(plr) {
    ctrl.removePlayer(plr);
    clearPQItems(plr);
    plr.warp(exitMapID);

    if (plr.isPartyLeader() || ctrl.playerCount() < 3) {
        var players = ctrl.players();
        for (var i = 0; i < players.length; i++) {
            clearPQItems(players[i]);
            players[i].warp(exitMapID);
        }
        ctrl.finished();
    }
}
```

## Entry NPC Pattern

Preferred entry NPC responsibilities:

- validate party/guild size
- validate leader authority
- validate level range
- validate same-map participation
- optionally award entry buffs
- call `startPartyQuest(...)` or `startGuildQuest(...)`

Keep entry scripts short.

## Stage NPC Pattern

Preferred stage NPC responsibilities:

- leader-only stage clear interactions
- item hand-in checks
- stage-specific explanation text
- stage clear reward/EXP
- set `map.properties()["clear"] = true` and any stage-specific keys

Example:

```js
if (!plr.isLeader()) {
    npc.sendOk("Please ask your party leader to speak with me.");
} else if (map.properties()["clear"]) {
    npc.sendOk("Please continue to the next stage.");
} else if (plr.haveItem(ITEM_ID, COUNT)) {
    plr.gainItem(ITEM_ID, -COUNT);
    map.properties()["clear"] = true;
    map.showEffect("quest/party/clear");
    map.playSound("Party1/Clear");
    plr.partyGiveExp(7500);
    npc.sendOk("Well done. You may proceed.");
} else {
    npc.sendOk("Bring me the required items first.");
}
```

## Portal Script Pattern

Portal scripts should be very small.

Example:

```js
if (plr.getEventProperty("3stageclear") == null) {
    portal.warp(920010200, "st00");
} else {
    portal.block("You may not go back in this room.");
}
```

Preferred uses:

- block forward movement until clear
- block re-entry to previous rooms after clear
- route center tower / hub room branches

## Reactor Script Pattern

Use reactor scripts for actual reactor-driven behavior.

Supported useful helpers currently include:

- `rm.spawnNpcAtReactor(id)`
- `rm.spawnMonster(id, x, y)`
- `rm.dropItems()`
- `rm.gainGuildPoints(points)`
- `rm.hitMapReactorByName(mapID, name)`
- `rm.mapMessage(type, msg)`
- `rm.showEffect(path)`
- `rm.playSound(path)`

Example:

```js
function act() {
    rm.showEffect("quest/party/clear");
    rm.playSound("Party1/Clear");
    rm.spawnNpcAtReactor(2013001);
}
```

## Party vs Guild PQs

### Party PQs

Use:

- `plr.startPartyQuest(name, instID)`
- `plr.leavePartyQuest()`

### Guild PQs

Use:

- `plr.startGuildQuest(name, instID)`
- `plr.joinGuildQuest()`

Guild-specific details:

- entry NPC should validate guild rank/authority
- event can use guild-specific rewards via `plr.gainGuildPoints(...)`

## Current Good References

Use these as the primary target-side references:

- KPQ event: `scripts/event/kerning_pq.js`
- LPQ event: `scripts/event/ludibrium_pq.js`
- OPQ event: `scripts/event/orbis_pq.js`
- GPQ event: `scripts/event/guild_pq.js`

Use these as supporting reference layers:

- OPQ controller NPC: `scripts/npc/2013001.js`
- GPQ entry/finisher NPCs: `scripts/npc/9040000.js`, `scripts/npc/9040010.js`
- GPQ reactor scripts: `scripts/reactor/9201001.js`, `scripts/reactor/9208004.js`, `scripts/reactor/9208007.js`

## Common Mistakes To Avoid

### 1. Do not put PQ-specific hacks in unrelated systems

Bad:

- hardcoding one PQ's item behavior directly in `moveItem`

Good:

- use reactor scripts
- use event scripts
- use NPC scripts

### 2. Do not duplicate cleanup across many places

Put PQ item cleanup in event lifecycle functions.

### 3. Do not overuse NPCs for map progression

If a portal gate is the real mechanic, use a portal script.

### 4. Do not split event state across multiple systems

Use field-instance properties for stage state.

### 5. Do not make exit-map NPCs run full PQ logic

Exit-map NPCs should stay light.

## Recommended Workflow For New PQs

1. Add the event script first.
2. Add entry NPC.
3. Add core item cleanup.
4. Add portal scripts for room gating.
5. Add reactor scripts for reactor-driven behavior.
6. Add stage NPCs one stage at a time.
7. Keep each stage self-contained.
8. Validate leave/timeout/finish before deeper polish.

## Minimum Validation Checklist

Before calling a PQ "added", verify:

- entry works
- timer starts
- players warp into the right first map
- leave path works
- timeout path works
- finish path works
- PQ items are cleaned up on leave/timeout/finish
- stage clear sets properties consistently
- portal gating respects `clear`
- reactor-driven stages actually respond
- re-entering a fresh instance resets maps, props, and spawned NPC state

## Current Scope Reminder

For this repo right now, keep PQ additions scoped to:

- GuildPQ
- Orbis PQ

Do not broaden into other PQ families unless requested.
