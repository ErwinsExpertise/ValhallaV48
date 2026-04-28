var leader = plr.getEventProperty("leader");
var alreadyClear = plr.getEventProperty("stage4clear");

if (leader == null) {
    plr.warp(990001100);
} else if (leader !== plr.name()) {
    if (alreadyClear === true || alreadyClear === "true") {
        npc.sendOk("The path ahead is already open. Your hardest test still lies before you.");
    } else {
        npc.sendOk("I need the leader of your guild's group to speak with me.");
    }
} else if (alreadyClear === true || alreadyClear === "true") {
    npc.sendOk("The path ahead is already open. Your hardest test still lies before you.");
} else {
    plr.setEventProperty("stage4clear", true);
    plr.gainGuildPoints(180);
    map.showEffect("quest/party/clear");
    map.playSound("Party1/Clear");
    map.hitReactorByName("ghostgate");
    npc.sendOk("I have opened the path for you. Go now, and face the evil waiting ahead.");
}
