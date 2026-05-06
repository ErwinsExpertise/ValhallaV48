var DELIVERY_ITEM = 4031486;
var PRE_KEY = 8839;
var STATE_KEY = 8835;

if (plr.questData(STATE_KEY) === "end") {
    npc.sendOk("I already received Maple Claws' present. Thank you again.");
} else if (plr.questData(PRE_KEY) !== "ing") {
    npc.sendOk("Please don't bother me unless you need me right this minute.");
} else if (plr.questData(STATE_KEY) !== "ing") {
    npc.sendOk("Maple Claws sent you? Then you must be carrying my present. Please speak with me again once you're ready to hand it over.");
    plr.setQuestData(STATE_KEY, "ing");
} else if (!plr.haveItem(DELIVERY_ITEM, 1)) {
    npc.sendOk("You do not seem to have the present anymore. Please ask Maple Claws for another one.");
} else if (Math.random() < 0.5) {
    npc.sendOk("I'm not in the mood right now. Please come back and try again in a little while.");
} else if (!plr.gainItem(DELIVERY_ITEM, -1)) {
    npc.sendOk("I couldn't take the present from you. Please try again.");
} else {
    plr.setQuestData(STATE_KEY, "end");
    npc.sendOk("Thank you for delivering Maple Claws' present to me. Please give Maple Claws my thanks.");
}
