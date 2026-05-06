var choice = npc.askMenu(
    "What would you like to know about Happyville?#b",
    "Tell me about the giant Christmas trees",
    "Tell me about the Extra Frosty Snow Zone"
);

if (choice === 0) {
    npc.sendNext("See the group of snowmen around town? Speak with one of them and they will send you to one of the giant tree rooms.");
    npc.sendNext("Each tree room only holds a handful of people at a time, and once you go in you should finish up before leaving. Speak with the scarf snowman inside when you are ready to come back.");
    npc.sendNext("Rudi also sells Christmas ornaments nearby. Most of the decorations are easy to buy, but the biggest and prettiest ornament is said to be carried off by monsters somewhere in the world.");
} else {
    npc.sendNext("Roodolph can take you to the Extra Frosty Snow Zone, but you need to equip #b#t1472063##k first. The gift boxes around Happyville may contain a pair.");
    npc.sendNext("Once you reach the snow zone, gather #b#t4031875##k and bring it to Happy so the snow machine can be refilled.");
}
