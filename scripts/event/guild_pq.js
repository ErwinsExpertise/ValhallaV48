var maps = [990000000, 990000100, 990000200, 990000300, 990000400, 990000500, 990000600, 990000700, 990000800, 990000900, 990001000];
var waitingMapID = 990000000;
var exitMapID = 990001100;
var GQItems = [1032033, 4001024, 4001025, 4001026, 4001027, 4001028, 4001031, 4001032, 4001033, 4001034, 4001035, 4001037];

function clearGQItems(plr) {
    for (var i = 0; i < GQItems.length; i++) {
        plr.removeAll(GQItems[i]);
    }
}

function start() {
    ctrl.setDuration("60m");

    for (var i = 0; i < maps.length; i++) {
        var field = ctrl.getMap(maps[i]);
        field.reset();
        field.clearProperties();
    }

    var waitingField = ctrl.getMap(waitingMapID);
    waitingField.properties()["canEnter"] = true;
	waitingField.properties()["leader"] = ctrl.players()[0].name();
	waitingField.properties()["entryTimestamp"] = Date.now().toString();
	ctrl.schedule("begin", "1m");

    var players = ctrl.players();
    var time = ctrl.remainingTime();
    for (var j = 0; j < players.length; j++) {
        players[j].warp(waitingMapID);
        players[j].showCountdown(time);
    }
}

function begin() {
    var waitingField = ctrl.getMap(waitingMapID);
    waitingField.properties()["canEnter"] = false;
	ctrl.schedule("earringcheck", "15s");

    var players = ctrl.players();
    for (var i = 0; i < players.length; i++) {
        players[i].sendMessage("[Guild Quest] The quest has begun!");
    }
}

function earringcheck() {
	var players = ctrl.players();
	for (var i = 0; i < players.length; i++) {
		var plr = players[i];
		if (plr.mapID() > 990000200 && plr.itemCount(1032033) < 1) {
			plr.setHP(1);
			plr.sendMessage("[Guild Quest] You were struck down for entering without the proper earrings.");
		}
	}
	ctrl.schedule("earringcheck", "15s");
}

function beforePortal(plr, src, dst) {
	var props = src.properties();
	if (props["clear"]) {
		return true;
	}
	return true;
}

function afterPortal(plr, dst) {
	plr.showCountdown(ctrl.remainingTime());
	if (dst.properties()["clear"]) {
		plr.portalEffect("gate");
	}
}

function timeout(plr) {
	clearGQItems(plr);
	plr.warp(exitMapID);
}

function playerLeaveEvent(plr) {
	ctrl.removePlayer(plr);
	clearGQItems(plr);
	plr.warp(exitMapID);

	if (plr.isLeader() || ctrl.playerCount() < 6) {
		var players = ctrl.players();
		for (var i = 0; i < players.length; i++) {
			clearGQItems(players[i]);
			players[i].warp(exitMapID);
		}
		ctrl.finished();
	}
}
