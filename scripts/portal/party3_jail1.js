var count = map.playerCountInMap(920010910) + map.playerCountInMap(920010911) + map.playerCountInMap(920010912);
if (count > 0) {
    portal.block("Someone is already inside.");
} else {
    portal.warp(920010910, "out00");
}
