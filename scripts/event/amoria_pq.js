var entryMapID = 670010200;
var exitMapID = 670011000;
var bonusMapID = 670010800;
var maps = [670010200, 670010300, 670010301, 670010302, 670010400, 670010500, 670010600, 670010700, 670010750, 670010800, 670011000];
var pqItems = [4031592, 4031593, 4031594, 4031595, 4031596, 4031597];
var bonusStarted = false;

function clearPQItems(plr) {
    for (var i = 0; i < pqItems.length; i++) {
        plr.removeAll(pqItems[i]);
    }
}

function start() {
    ctrl.setDuration("60m");

    for (var i = 0; i < maps.length; i++) {
        var field = ctrl.getMap(maps[i]);
        field.reset();
        field.clearProperties();
        if (maps[i] === 670010300 || maps[i] === 670010301 || maps[i] === 670010302 || maps[i] === 670010400) {
            field.portalEnabled(false, "next00");
        }
        if (maps[i] === 670010200) {
            field.portalEnabled(false, "go00");
            field.portalEnabled(false, "go01");
            field.portalEnabled(false, "go02");
        }
    }

    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (var j = 0; j < players.length; j++) {
        players[j].warp(entryMapID);
        players[j].showCountdown(time);
    }
}

function afterPortal(plr, dst) {
    plr.showCountdown(ctrl.remainingTime());
}

function onMapChange(plr, dst) {
    if (!bonusStarted && dst.getMapID() === bonusMapID) {
        bonusStarted = true;
        ctrl.setDuration("1m");
        var players = ctrl.players();
        var time = ctrl.remainingTime();
        for (var i = 0; i < players.length; i++) {
            players[i].showCountdown(time);
        }
        return;
    }

    plr.showCountdown(ctrl.remainingTime());
}

function timeout(plr) {
    clearPQItems(plr);
    plr.warp(exitMapID);
}

function playerLeaveEvent(plr) {
    clearPQItems(plr);
    ctrl.removePlayer(plr);
    plr.warp(exitMapID);

    if (plr.isPartyLeader() || ctrl.playerCount() < 6) {
        var players = ctrl.players();
        for (var i = 0; i < players.length; i++) {
            clearPQItems(players[i]);
            players[i].warp(exitMapID);
        }
        ctrl.finished();
    }
}
