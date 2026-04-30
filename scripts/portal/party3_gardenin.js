if (!plr.isLeader()) {
    portal.block("Only the party leader may decide when to leave this room.");
} else if (!plr.haveItem(4001055, 1)) {
    portal.block("We need the power of the Grass of Life.");
} else {
    plr.warpEventMembersToPortal(920010100, "st02");
}
