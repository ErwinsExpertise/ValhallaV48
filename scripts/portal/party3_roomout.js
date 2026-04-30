if (!plr.isLeader()) {
    portal.block("Only the party leader may decide when to leave this room.");
} else {
    var field = plr.mapID();
    var retportal = "in00";

    if (field === 920010300) retportal = "in01";
    else if (field === 920010400) retportal = "in02";
    else if (field === 920010500) retportal = "in03";
    else if (field === 920010600) {
        if (map.playerCountInMap(920010601) !== 0 || map.playerCountInMap(920010602) !== 0 || map.playerCountInMap(920010603) !== 0 || map.playerCountInMap(920010604) !== 0) {
            portal.block("You cannot leave while someone in your party is still inside one of the bedrooms.");
        } else {
            retportal = "in04";
            plr.warpEventMembersToPortal(920010100, retportal);
        }
    } else if (field === 920010700) retportal = "in05";
    else if (field === 920011000) retportal = "in06";

    if (field !== 920010600) {
        plr.warpEventMembersToPortal(920010100, retportal);
    }
}
