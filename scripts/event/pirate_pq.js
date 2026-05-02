var entryMapID = 925100000;
var exitMapID = 925100700;
var rewardMapID = 925100600;
var bossMapID = 925100500;
var maps = [925100000, 925100100, 925100200, 925100201, 925100202, 925100300, 925100301, 925100302, 925100400, 925100500, 925100600, 925100700];
var pqItems = [4001117, 4001119, 4001120, 4001121, 4001122, 4001123];

var stageTimers = {
    925100000: "240s",
    925100100: "360s",
    925100200: "360s",
    925100300: "360s",
    925100400: "360s",
    925100500: "480s",
    925100600: "300s"
};

var deckSpawns = [
    [9300123, 300, 220], [9300123, 0, 220], [9300123, 300, 220], [9300123, 600, 220], [9300123, 200, 9], [9300123, 500, 9], [9300123, 800, 9],
    [9300124, 150, 220], [9300124, 450, 220], [9300124, 750, 220], [9300124, 1200, 9], [9300124, 1500, 9], [9300124, 1800, 9],
    [9300125, 1050, 220], [9300125, 1350, 220], [9300125, 1650, 220], [9300125, 1950, 220], [9300125, 1350, 9], [9300125, 1650, 9]
];

function clearPQItems(plr) {
    for (var i = 0; i < pqItems.length; i++) {
        plr.removeAll(pqItems[i]);
    }
}

function clearPQItemsForParty() {
    var players = ctrl.players();
    for (var i = 0; i < players.length; i++) {
        clearPQItems(players[i]);
    }
}

function computePartyStats() {
    var players = ctrl.players();
    var count = 0;
    var totalLevel = 0;
    var over70 = 0;

    for (var i = 0; i < players.length; i++) {
        var level = players[i].level();
        if (level <= 0) {
            continue;
        }
        count++;
        totalLevel += level;
        if (level > 70) {
            over70++;
        }
    }

    ctrl.setProperty("partySize", count);
    ctrl.setProperty("avgLevel", count > 0 ? Math.floor(totalLevel / count) : 0);
    ctrl.setProperty("over70", over70);
}

function setStageTimer(mapID) {
    var timer = stageTimers[mapID];
    if (timer) {
        ctrl.setDuration(timer);
    }
}

function spawnDeckMobs(mapID) {
    var field = ctrl.getMap(mapID);
    for (var i = 0; i < deckSpawns.length; i++) {
        field.spawnMonster(deckSpawns[i][0], deckSpawns[i][1], deckSpawns[i][2]);
    }
}

function resetPQMaps() {
    for (var i = 0; i < maps.length; i++) {
        var field = ctrl.getMap(maps[i]);
        field.reset();
        field.clearProperties();
    }
}

function initStage(mapID) {
    if (mapID === 925100100 && !ctrl.getProperty("stage2Initialized")) {
        var stage2 = ctrl.getMap(925100100);
        stage2.setMobSpawnEnabled(9300114, false);
        stage2.setMobSpawnEnabled(9300115, false);
        stage2.setMobSpawnEnabled(9300116, false);
        stage2.removeAllMobs();
        ctrl.setProperty("stage2Initialized", true);
        return;
    }

    if (mapID === 925100200 && !ctrl.getProperty("stage3Initialized")) {
        spawnDeckMobs(925100200);
        ctrl.setProperty("stage3Initialized", true);
        return;
    }

    if (mapID === 925100300 && !ctrl.getProperty("stage4Initialized")) {
        spawnDeckMobs(925100300);
        ctrl.setProperty("stage4Initialized", true);
        return;
    }

    if (mapID === bossMapID && !ctrl.getProperty("bossSpawned")) {
        ctrl.getMap(bossMapID).spawnMonster(9300119, 630, 213);
        ctrl.setProperty("bossSpawned", true);
        ctrl.schedule("checkBossClear", "1s");
    }
}

function checkBossClear() {
    if (!ctrl.getProperty("bossSpawned") || ctrl.getProperty("bossNpcSpawned")) {
        return;
    }

    var bossMap = ctrl.getMap(bossMapID);
    if (bossMap.mobCountByID(9300119) > 0) {
        ctrl.schedule("checkBossClear", "1s");
        return;
    }

    bossMap.hitReactorByTemplate(2516000);
    bossMap.showEffect("quest/party/clear");
    bossMap.playSound("Party1/Clear");
    ctrl.setProperty("bossNpcSpawned", true);
}

function start() {
    ctrl.setDuration(stageTimers[entryMapID]);

    resetPQMaps();

    ctrl.getMap(bossMapID).removeNpcByTemplate(2094001);

    computePartyStats();
    ctrl.setProperty("mobGen", "0");
    ctrl.setProperty("completed", false);
    ctrl.setProperty("bossSpawned", false);
    ctrl.setProperty("bossNpcSpawned", false);
    ctrl.setProperty("bossRewarded", false);
    ctrl.setProperty("stage2Initialized", false);
    ctrl.setProperty("stage3Initialized", false);
    ctrl.setProperty("stage4Initialized", false);

    clearPQItemsForParty();

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
    var mapID = dst.getMapID();
    setStageTimer(mapID);
    initStage(mapID);
    plr.showCountdown(ctrl.remainingTime());
}

function timeout(plr) {
    clearPQItems(plr);
    resetPQMaps();
    plr.warp(exitMapID);
}

function playerLeaveEvent(plr) {
    clearPQItems(plr);
    ctrl.removePlayer(plr);
    plr.warp(exitMapID);

    if (ctrl.getProperty("completed")) {
        if (ctrl.playerCount() < 1) {
            resetPQMaps();
            ctrl.finished();
        }
        return;
    }

    if (plr.isPartyLeader() || ctrl.playerCount() < 3) {
        var players = ctrl.players();
        for (var i = 0; i < players.length; i++) {
            clearPQItems(players[i]);
            players[i].warp(exitMapID);
        }
        resetPQMaps();
        ctrl.finished();
    }
}
