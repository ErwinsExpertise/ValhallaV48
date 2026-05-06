var state = plr.questData(7500);
var job = plr.job();
var currentMap = plr.mapID();
var targetMap = 0;
var blockedMessage = "There is a door here that leads to another dimension, but you can't enter it right now.";
var occupiedMessage = "Someone else is already fighting in there. Please come back later.";

if (currentMap === 105070001 && state === "p1" && (job === 110 || job === 120 || job === 130)) {
    targetMap = 108010301;
} else if (currentMap === 100040106 && state === "p1" && (job === 210 || job === 220 || job === 230)) {
    targetMap = 108010201;
} else if (currentMap === 105040305 && state === "p1" && (job === 310 || job === 320)) {
    targetMap = 108010101;
}

if (targetMap === 0) {
    portal.block(blockedMessage);
} else if (map.playerCount(targetMap, 0) > 0) {
    portal.block(occupiedMessage);
} else {
    portal.warp(targetMap, "");
}
