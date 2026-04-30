if (plr.getEventProperty("r4_rp") !== "1") {
    plr.setEventProperty("r4way1", String(1 + Math.floor(Math.random() * 3)));
    plr.setEventProperty("r4way2", String(1 + Math.floor(Math.random() * 3)));
    plr.setEventProperty("r4_rp", "1");
}

var current = plr.mapID();
var id = portal.id();
var way1 = plr.getEventProperty("r4way1");
var way2 = plr.getEventProperty("r4way2");

if (id === 11) portal.warp(current, way1 === "1" ? "np00" : "np02");
else if (id === 12) portal.warp(current, way1 === "2" ? "np00" : "np02");
else if (id === 13) portal.warp(current, way1 === "3" ? "np00" : "np02");
else if (id === 14) portal.warp(current, way2 === "1" ? "np01" : "np02");
else if (id === 15) portal.warp(current, way2 === "2" ? "np01" : "np02");
else if (id === 16) portal.warp(current, way2 === "3" ? "np01" : "np02");
