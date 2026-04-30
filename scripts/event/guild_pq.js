var maps = [
    990000000, 990000100, 990000200, 990000300, 990000301,
    990000400, 990000401, 990000410, 990000420, 990000430, 990000431, 990000440,
    990000500, 990000501, 990000502,
    990000600, 990000610, 990000611, 990000620, 990000630, 990000631, 990000640, 990000641,
    990000700, 990000800, 990000900, 990001000
];
var waitingMapID = 990000000;
var exitMapID = 990001100;
var GQItems = [1032033, 4001024, 4001025, 4001026, 4001027, 4001028, 4001029, 4001030, 4001031, 4001032, 4001033, 4001034, 4001035, 4001036, 4001037];

function clearGQItems(plr) {
    for (var i = 0; i < GQItems.length; i++) {
        plr.removeAll(GQItems[i]);
    }
}

function players() {
    return ctrl.players().filter(function(plr) { return plr != null; });
}

function broadcast(msg) {
    var list = players();
    for (var i = 0; i < list.length; i++) {
        list[i].sendMessage(msg);
    }
}

function warpOutAll(reason) {
    var list = players();
    for (var i = 0; i < list.length; i++) {
        clearGQItems(list[i]);
        if (reason) {
            list[i].sendMessage(reason);
        }
        list[i].warp(exitMapID);
    }
}

function failEvent(reason) {
    if (reason) {
        broadcast(reason);
        ctrl.log("gpq fail: " + reason);
    }
    warpOutAll("");
    ctrl.finished();
}

function start() {
    ctrl.setDuration("93m");

    for (var i = 0; i < maps.length; i++) {
        var field = ctrl.getMap(maps[i]);
        field.reset();
        field.clearProperties();
    }

    var waitingField = ctrl.getMap(waitingMapID);
    waitingField.properties()["canEnter"] = true;
    waitingField.properties()["leader"] = ctrl.players()[0].name();
    waitingField.properties()["leaderID"] = "" + ctrl.players()[0].id();
    waitingField.properties()["entryTimestamp"] = Date.now().toString();

    ctrl.setProperty("completed", false);
    ctrl.setProperty("bonusStage", false);

    ctrl.schedule("oneMinuteWarning", "2m");
    ctrl.schedule("thirtySecondWarning", "2m30s");
    ctrl.schedule("begin", "3m");

    var list = players();
    var time = ctrl.remainingTime();
    for (var j = 0; j < list.length; j++) {
        clearGQItems(list[j]);
        list[j].warp(waitingMapID);
        list[j].showCountdown(time);
    }

    broadcast("[Guild Quest] The door to Sharenian will open in 3 minutes.");
    ctrl.log("gpq start leader=" + ctrl.players()[0].name() + " players=" + list.length);
}

function oneMinuteWarning() {
    if (ctrl.getMap(waitingMapID).properties()["canEnter"]) {
        broadcast("[Guild Quest] The door will open in 1 minute.");
    }
}

function thirtySecondWarning() {
    if (ctrl.getMap(waitingMapID).properties()["canEnter"]) {
        broadcast("[Guild Quest] The door will open in 30 seconds.");
    }
}

function begin() {
    var waitingField = ctrl.getMap(waitingMapID);
    waitingField.properties()["canEnter"] = false;

    if (ctrl.playerCount() < 6) {
        failEvent("[Guild Quest] Fewer than 6 guild members entered before the door opened. The run has been cancelled.");
        return;
    }

    broadcast("[Guild Quest] The quest has begun!");
    ctrl.schedule("earringcheck", "15s");
    ctrl.log("gpq begin players=" + ctrl.playerCount());
}

function earringcheck() {
    if (ctrl.getProperty("completed") === true || ctrl.getProperty("completed") === "true") {
        return;
    }

    var list = players();
    for (var i = 0; i < list.length; i++) {
        var plr = list[i];
        if (plr.mapID() > 990000200 && plr.itemCount(1032033) < 1) {
            plr.setHP(1);
            plr.sendMessage("[Guild Quest] You were struck down for entering without the Protector Rock.");
        }
    }

    ctrl.schedule("earringcheck", "15s");
}

function beforePortal(plr, src, dst) {
    return true;
}

function afterPortal(plr, dst) {
    plr.showCountdown(ctrl.remainingTime());
}

function timeout(plr) {
    clearGQItems(plr);
    plr.warp(exitMapID);
}

function playerLeaveEvent(plr) {
    var completed = ctrl.getProperty("completed") === true || ctrl.getProperty("completed") === "true";
    var leader = plr.getEventProperty("leader");

    ctrl.removePlayer(plr);
    clearGQItems(plr);
    plr.warp(exitMapID);

    if (completed) {
        ctrl.log("gpq exit after clear player=" + plr.name() + " remaining=" + ctrl.playerCount());
        if (ctrl.playerCount() === 0) {
            ctrl.finished();
        }
        return;
    }

    ctrl.log("gpq leave player=" + plr.name() + " remaining=" + ctrl.playerCount());
    if (leader === plr.name() || ctrl.playerCount() < 6) {
        failEvent("[Guild Quest] The run has ended because the registered leader left or fewer than 6 guild members remain.");
    }
}
