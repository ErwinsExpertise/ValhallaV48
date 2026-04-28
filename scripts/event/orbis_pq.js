var maps = [920010000, 920010100, 920010200, 920010300, 920010400, 920010500, 920010600, 920010700, 920010800, 920010900, 920011000, 920011100, 920011300];
var centerMapID = 920010100;
var entryMapID = 920010000;
var rewardMapID = 920011300;
var exitMapID = 920011200;
var chamberlainID = 2013001;
var entryPrompt = "Hi, my name is Eak, the Chamberlain of the Goddess. Don't be alarmed; you won't be able to see me right now. Back when the Goddess turned into a block of stone, I simultaneously lost my own power. If you gather up the power of the Magic Cloud of Orbis, however, then I'll be able to recover my body and re-transform back to my original self. Please collect #b20#k Magic Clouds and bring them back to me. Right now, you'll only see me as a tiny, flickering light.";
var pqItems = [4001063, 4001044, 4001045, 4001046, 4001047, 4001048, 4001049, 4001050, 4001051, 4001052];

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
        if (maps[i] === entryMapID) {
            field.removeNpcByTemplate(chamberlainID);
        }
    }

    var centerMap = ctrl.getMap(centerMapID);
    centerMap.setPortalScriptByID(13, "orbisPQSealedRoom");
    centerMap.setPortalScriptByID(4, "orbisPQWalkway");
    centerMap.setPortalScriptByID(12, "orbisPQStorage");
    centerMap.setPortalScriptByID(5, "orbisPQLobby");
    centerMap.setPortalScriptByID(14, "orbisPQOnTheWayUp");
    centerMap.setPortalScriptByID(15, "orbisPQLounge");
    centerMap.setPortalScriptByID(16, "orbisPQRoomOfDarkness");

    ctrl.getMap(920010200).setPortalScriptByID(13, "orbisPQWalkwayExit");
    ctrl.getMap(920010300).setPortalScriptByID(1, "orbisPQStorageExit");
    ctrl.getMap(920010400).setPortalScriptByID(8, "orbisPQLobbyExit");
    ctrl.getMap(920010500).setPortalScriptByID(3, "orbisPQSRExit");
    ctrl.getMap(920010600).setPortalScriptByID(17, "orbisPQLoungeExit");
    ctrl.getMap(920010700).setPortalScriptByID(23, "orbisPQOnTheWayUpExit");
    ctrl.getMap(920010800).setPortalScriptByID(1, "orbisPQGardenExit");
    ctrl.getMap(920011000).setPortalScriptByID(1, "orbisPQRoomOfDarknessExit");

    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (var j = 0; j < players.length; j++) {
        players[j].warp(entryMapID);
        players[j].showCountdown(time);
        players[j].showNpcOk(chamberlainID, entryPrompt);
    }
}

function beforePortal(plr, src, dst) {
    if (src.getMapID() === centerMapID) {
        return true;
    }
    var props = src.properties();
    if (props["clear"]) {
        return true;
    }
    plr.sendMessage("You cannot use the portal yet.");
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
    clearPQItems(plr);
    ctrl.removePlayer(plr);
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
