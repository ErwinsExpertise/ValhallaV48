var allMaps = [
    920010000, 920010100, 920010200, 920010300, 920010400, 920010500, 920010600,
    920010601, 920010602, 920010603, 920010604, 920010700, 920010800, 920010900,
    920010910, 920010911, 920010912, 920010920, 920010921, 920010922, 920010930,
    920010931, 920010932, 920011000, 920011100, 920011200, 920011300
];
var entryMapID = 920010000;
var bonusMapID = 920011100;
var rewardMapID = 920011300;
var exitMapID = 920011200;
var chamberlainID = 2013001;
var entryPrompt = "Hi, my name is Eak, the Chamberlain of the Goddess. Collect #b20 Cloud Pieces#k to restore my form, then gather the statue pieces and the #bGrass of Life#k to save Goddess Minerva.";
var cleanupItems = [
    4001044, 4001045, 4001046, 4001047, 4001048, 4001049,
    4001050, 4001051, 4001052, 4001053, 4001054, 4001055,
    4001056, 4001057, 4001058, 4001059, 4001060, 4001061,
    4001062, 4001063, 4001074
];

function clearPQItems(plr) {
    for (var i = 0; i < cleanupItems.length; i++) {
        plr.removeAll(cleanupItems[i]);
    }
}

function setRewardBoxes() {
    if (Math.floor(Math.random() * 10) < 2) {
        ctrl.getMap(920010912).setReactorStateByName("party3_box01", 1);
    }
    if (Math.floor(Math.random() * 10) >= 2 && Math.floor(Math.random() * 10) < 4) {
        ctrl.getMap(920010922).setReactorStateByName("party3_box02", 1);
    }
    if (Math.floor(Math.random() * 10) >= 4 && Math.floor(Math.random() * 10) < 6) {
        ctrl.getMap(920010932).setReactorStateByName("party3_box03", 1);
    }
}

function start() {
    ctrl.setDuration("60m");
    ctrl.setProperty("completed", false);
    ctrl.setProperty("bonusStarted", false);
    ctrl.setProperty("bonusTimerActive", false);
    ctrl.setProperty("rewardPhase", false);
    ctrl.setProperty("rewardWarped", false);

    for (var i = 0; i < allMaps.length; i++) {
        var field = ctrl.getMap(allMaps[i]);
        field.reset();
        field.clearProperties();
        if (allMaps[i] === entryMapID) {
            field.removeNpcByTemplate(chamberlainID);
        }
    }

    setRewardBoxes();
    ctrl.getMap(entryMapID).hitReactorByTemplate(2006000);

    ctrl.log("OPQ instance created for " + ctrl.playerCount() + " player(s)");
    ctrl.warpPlayersToPortal(entryMapID, "sp");

    var players = ctrl.players();
    for (var j = 0; j < players.length; j++) {
        players[j].showCountdown(ctrl.remainingTime());
        players[j].showNpcOk(chamberlainID, entryPrompt);
    }
}

function onMapChange(plr, dst) {
    var mapID = dst.getMapID();

    if (ctrl.getProperty("bonusStarted") && !ctrl.getProperty("bonusTimerActive") && mapID === bonusMapID) {
        ctrl.setProperty("bonusTimerActive", true);
        ctrl.setDuration("60s");
        ctrl.log("OPQ bonus room started");
    }

    if (mapID === rewardMapID || mapID === exitMapID) {
        plr.hideCountdown();
        return;
    }

    plr.showCountdown(ctrl.remainingTime());
    if (dst.properties()["clear"]) {
        plr.portalEffect("gate");
    }
}

function timeout(plr) {
    if (ctrl.getProperty("completed") && ctrl.getProperty("bonusStarted") && !ctrl.getProperty("rewardPhase")) {
        if (!ctrl.getProperty("rewardWarped")) {
            ctrl.setProperty("rewardPhase", true);
            ctrl.setProperty("rewardWarped", true);
            ctrl.setDuration("10m");
            ctrl.log("OPQ bonus room completed, warping party to reward map");
            ctrl.warpPlayersToPortal(rewardMapID, "sp");
        }
        return;
    }

    clearPQItems(plr);
    plr.hideCountdown();
    plr.warpToPortalNameInInstance(exitMapID, "sp", 0);
}

function playerLeaveEvent(plr) {
    clearPQItems(plr);
    plr.hideCountdown();
    ctrl.removePlayer(plr);
    plr.warpToPortalNameInInstance(exitMapID, "sp", 0);

    if (ctrl.getProperty("completed")) {
        ctrl.log("OPQ participant left after completion; remaining=" + ctrl.playerCount());
        if (ctrl.playerCount() < 1) {
            ctrl.log("OPQ instance cleaned up after completion");
            ctrl.finished();
        }
        return;
    }

    ctrl.log("OPQ participant left before completion; remaining=" + ctrl.playerCount());
    if (plr.isPartyLeader() || ctrl.playerCount() < 3) {
        var players = ctrl.players();
        for (var i = 0; i < players.length; i++) {
            clearPQItems(players[i]);
            players[i].hideCountdown();
            players[i].warpToPortalNameInInstance(exitMapID, "sp", 0);
        }
        ctrl.log("OPQ failed and cleaned up");
        ctrl.finished();
    }
}
