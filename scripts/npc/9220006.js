var total = 0;
var props = map.properties();
if (props["wxmasCount"]) {
    total = parseInt(String(props["wxmasCount"]), 10);
}

if (props["wxmasBoss"] === "1" && map.mobCountByID(9400707) + map.mobCountByID(9400708) + map.mobCountByID(9400709) + map.mobCountByID(9400710) === 0) {
    props["wxmasBoss"] = "0";
    props["wxmasCount"] = "0";
    total = 0;
    map.removeMobsByID(9400714);
    map.removeMobsByID(9400715);
    map.removeMobsByID(9400716);
    map.removeMobsByID(9400717);
    map.removeMobsByID(9400718);
    map.removeMobsByID(9400719);
    map.removeMobsByID(9400720);
    map.removeMobsByID(9400721);
    map.removeMobsByID(9400722);
    map.removeMobsByID(9400723);
    map.removeMobsByID(9400724);
    map.spawnMonster(9400714, 1450, 140);
}

if (map.mobCountByID(9400707) > 0 || map.mobCountByID(9400708) > 0 || map.mobCountByID(9400709) > 0 || map.mobCountByID(9400710) > 0) {
    npc.sendOk("The snow machine has gone wild! Please help deal with the trouble before anything worse happens.");
} else {
    npc.sendNext("Merry Christmas!! Happy is keeping the snow machine running, and Roodolph can take you back to Happyville whenever you are ready.");
    npc.sendOk("Right now the snow machine has #b" + total + " / 50000#k Snow Powder. If you find more #b#t4031875##k, give it to Happy.");
}
