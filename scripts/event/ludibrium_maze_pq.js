var START_ROOMS = [809050001, 809050002, 809050003, 809050004, 809050005, 809050006, 809050007, 809050008, 809050009, 809050010, 809050011, 809050012, 809050013, 809050014];
var MAZE_MAPS = [809050000, 809050001, 809050002, 809050003, 809050004, 809050005, 809050006, 809050007, 809050008, 809050009, 809050010, 809050011, 809050012, 809050013, 809050014, 809050015];
var REWARD_MAP = 809050016;
var FAIL_MAP = 809050017;
var COUPON = 4001106;

function eventPlayers() {
    var players = ctrl.players();
    var live = [];

    for (let i = 0; i < players.length; i++) {
        if (players[i] && players[i].id() > 0) {
            live.push(players[i]);
        }
    }

    return live;
}

function clearMazeState() {
    for (let i = 0; i < MAZE_MAPS.length; i++) {
        var field = ctrl.getMap(MAZE_MAPS[i]);
        field.reset();
        field.clearProperties();
        field.removeDrops();
    }

    ctrl.getMap(REWARD_MAP).removeDrops();
    ctrl.getMap(REWARD_MAP).clearProperties();
    ctrl.getMap(FAIL_MAP).removeDrops();
    ctrl.getMap(FAIL_MAP).clearProperties();
}

function warpPlayers(players, mapID) {
    for (let i = 0; i < players.length; i++) {
        players[i].warp(mapID);
        players[i].showCountdown(ctrl.remainingTime());
    }
}

function finishRun(reason, mapID) {
    if (ctrl.getProperty("finished")) {
        return;
    }

    ctrl.setProperty("finished", reason);
    ctrl.log("lmpq: finishing run reason=" + reason + " remaining=" + ctrl.playerCount());

    var players = eventPlayers();
    for (let i = 0; i < players.length; i++) {
        players[i].removeAll(COUPON);
        players[i].hideCountdown();
        players[i].warp(mapID);
    }

    ctrl.finished();
}

function start() {
    clearMazeState();
    ctrl.setDuration("15m");

    var startMap = START_ROOMS[Math.floor(Math.random() * START_ROOMS.length)];
    ctrl.setProperty("startMap", startMap);
    ctrl.log("lmpq: start startMap=" + startMap + " players=" + ctrl.playerCount());

    warpPlayers(eventPlayers(), startMap);
}

function afterPortal(plr, dst) {
    ctrl.log("lmpq: player " + plr.id() + " moved to map=" + dst.getMapID());
    plr.showCountdown(ctrl.remainingTime());
}

function onMapChange(plr, dst) {
    if (!ctrl.getProperty("finished")) {
        plr.showCountdown(ctrl.remainingTime());
    }
}

function timeout(plr) {
    finishRun("timeout", FAIL_MAP);
}

function playerLeaveEvent(plr) {
    plr.removeAll(COUPON);
    plr.hideCountdown();
    ctrl.removePlayer(plr);
    ctrl.log("lmpq: player left id=" + plr.id() + " remaining=" + ctrl.playerCount());

    if (plr.mapID() != REWARD_MAP && plr.mapID() != FAIL_MAP) {
        plr.warp(FAIL_MAP);
    }

    if (ctrl.playerCount() === 0) {
        ctrl.setProperty("finished", "empty");
        ctrl.finished();
    }
}
