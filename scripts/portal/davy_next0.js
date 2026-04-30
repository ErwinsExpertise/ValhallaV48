var count = map.mobCountByID(9300108) + map.mobCountByID(9300109) + map.mobCountByID(9300110) + map.mobCountByID(9300111);
if (count > 0) {
    portal.block("You must defeat every monster here before moving on.");
} else {
    plr.setEventProperty("clear_1", true);
    plr.warp(925100100);
}
