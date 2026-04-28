if (!npc.sendYesNo("We're just about to depart. Are you sure you want to get off the train? You may do so, but then you'll have to wait until the next available train. Do you still wish to get off?")) {
    npc.sendOk("You'll get to your destination in a short while. Talk to other passengers and share your stories with them, and you'll be there before you know it.");
} else if (plr.mapID() === 220000111) {
    plr.warp(220000110);
} else {
    plr.warp(220000121);
}
