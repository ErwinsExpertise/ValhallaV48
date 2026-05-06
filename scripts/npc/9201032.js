var choice = npc.askMenu(
    "Confounded monsters! They keep running off with the supplies for my holiday sweets. What would you like to ask?#b",
    "What happened to your graham-cracker house?",
    "What can I do in Happyville?"
);

if (choice === 0) {
    npc.sendNext("I was building a proper graham-cracker holiday house when monsters made off with half the supplies. Until I replace them, all I can do is complain about it.");
    npc.sendOk("If you see me pacing around town, that is why. I am determined to rebuild it before the season is over.");
} else {
    npc.sendNext("Happyville has plenty to do. Maple Claws is collecting stolen presents, Torr is still looking for his horn, and the snowmen can send you to the giant Christmas tree rooms.");
    npc.sendOk("If you want something a little rougher, Roodolph can send you to the Extra Frosty Snow Zone once you have the right mittens.");
}
