var count = map.playerCountInMap(920010930) + map.playerCountInMap(920010931) + map.playerCountInMap(920010932);
if (count > 0) {
    portal.block("Someone is already inside.");
} else {
    portal.warp(920010930, "out00");
}
