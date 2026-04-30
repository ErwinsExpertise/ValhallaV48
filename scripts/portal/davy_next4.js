var count = map.mobCountByID(9300120) + map.mobCountByID(9300121) + map.mobCountByID(9300122) + map.mobCountByID(9300126);
var sealed = map.reactorStateByName("sMob1") === 1 && map.reactorStateByName("sMob2") === 1 && map.reactorStateByName("sMob3") === 1 && map.reactorStateByName("sMob4") === 1;
if (!sealed || count > 0) {
    portal.block("The portal is still sealed shut.");
} else {
    plr.setEventProperty("clear_5", true);
    plr.warp(925100500);
}
