var count = map.playerCountInMap(920010920) + map.playerCountInMap(920010921) + map.playerCountInMap(920010922);
if (count > 0) {
    portal.block("Someone is already inside.");
} else {
    portal.warp(920010920, "out00");
}
