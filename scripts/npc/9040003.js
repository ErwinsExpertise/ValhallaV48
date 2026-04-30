var leader = plr.getEventProperty("leader");
var alreadyClear = plr.getEventProperty("stage4clear");

if (leader == null) {
    plr.warp(990001100);
} else if (leader !== plr.name()) {
    if (alreadyClear === true || alreadyClear === "true") {
        npc.sendOk("The path ahead is already open. Hurry to the throne.");
    } else {
        npc.sendOk("I need the registered leader to speak with me.");
    }
} else if (alreadyClear === true || alreadyClear === "true") {
    npc.sendOk("The path ahead is already open. Hurry to the throne.");
} else {
    plr.setEventProperty("stage4clear", true);
    plr.setEventProperty("ghostgateopen", "yes");
    map.setReactorStateByName("ghostgate", 1);
    map.showEffect("quest/party/clear");
    map.playSound("Party1/Clear");
    plr.logEvent("gpq stage4 clear");
    npc.sendOk("At last... I can rest. The path to the throne is now open.");
}
