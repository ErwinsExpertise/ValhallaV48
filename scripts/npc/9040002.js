npc.sendSelection("Sharenian was once a great kingdom, and the Guild Quest is the path our strongest guilds take to uncover its secrets. What would you like to know?#b\r\n#L0#What's Sharenian?#l\r\n#L1##t4001024#? What's that?#l\r\n#L2#Guild Quest?#l\r\n#L3#No, I'm fine now.#l");
var selection = npc.selection();

if (selection === 0) {
    npc.sendNext("Sharenian was a civilization from the past that once held great influence over Victoria Island. Many of the old ruins still scattered across the island trace back to that age.");
    npc.sendOk("Its final king, Sharen III, was said to be wise and compassionate, but the kingdom collapsed without leaving behind a clear explanation.");
} else if (selection === 1) {
    npc.sendOk("#t4001024# is a legendary jewel said to grant eternal youth. Ironically, everyone who sought to possess it seems to have met a terrible end, and many believe it was tied to Sharenian's downfall.");
} else if (selection === 2) {
    npc.sendNext("We've sent many explorers into the ruins before, but none returned. That's why the Guild Quest exists now: to challenge guilds strong enough to uncover what lies ahead.");
    npc.sendOk("The ultimate goal of the Guild Quest is to explore Sharenian and find #t4001024#. It is not a trial brute force alone can solve. Teamwork matters much more there.");
} else {
    npc.sendOk("If you have anything else to ask, come back and speak with me again.");
}
