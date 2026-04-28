var hasItems = plr.haveItem(4031508, 5) && plr.haveItem(4031507, 5);

if (hasItems) {
    npc.sendNext("Wow! You succeeded in collecting 5 of both #b#t4031508##k and #b#t4031507##k. Okay then, I will send you to the Zoo. Please talk to me again when you get there.");
    plr.warp(230000003);
} else if (!npc.sendYesNo("You haven't completed the requirements. Are you sure you want to leave?")) {
    npc.sendOk("Stay here a little longer if you want to keep trying.");
} else {
    npc.sendOk("Well okay, I will send you back.");
    plr.warp(230000003);
}
