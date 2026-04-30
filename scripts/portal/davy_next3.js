var count = map.mobCountByID(9300123) + map.mobCountByID(9300124) + map.mobCountByID(9300125);
if (count > 0) {
    portal.block("The portal is still sealed shut.");
} else {
    plr.setEventProperty("clear_4", true);
    plr.warp(925100400);
}
