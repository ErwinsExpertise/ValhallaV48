if (!npc.sendYesNo("Do you want to leave the waiting room? You can, but the ticket is NOT refundable. Are you sure you still want to leave this room?")) {
    npc.sendOk("You'll get to your destination in a moment. Go ahead and talk to other people, and before you know it, you'll be there already.");
} else {
    plr.warp(260000100);
}
