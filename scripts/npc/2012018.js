var DELIVERY_ITEM = 4031486;
var PRE_KEY = 8841;
var STATE_KEY = 8837;

if (plr.questData(STATE_KEY) === "end") {
    npc.sendOk("My kitty and I already received Maple Claws' present. Thank you.");
} else if (plr.questData(PRE_KEY) !== "ing") {
    npc.sendOk("Look at my kitty! Isn't she cute?");
} else if (plr.questData(STATE_KEY) !== "ing") {
    npc.sendOk("Maple Claws sent you? Then you must be carrying my present. Please speak with me again once you're ready to hand it over.");
    plr.setQuestData(STATE_KEY, "ing");
} else if (!plr.haveItem(DELIVERY_ITEM, 1)) {
    npc.sendOk("You do not seem to have the present anymore. Please ask Maple Claws for another one.");
} else if (Math.random() < 0.5) {
    npc.sendOk("Hmm~ I'm not quite ready for company right now. Please try again later.");
} else if (!plr.gainItem(DELIVERY_ITEM, -1)) {
    npc.sendOk("I couldn't take the present from you. Please try again.");
} else {
    plr.setQuestData(STATE_KEY, "end");
    npc.sendOk("Thank you for bringing Maple Claws' present all the way to Orbis for me.");
}
