if (plr.getEventProperty("r6_rp") !== "1") {
    for (var i = 1; i <= 16; i++) {
        plr.setEventProperty("r6way" + i, String(1 + Math.floor(Math.random() * 4)));
    }
    plr.setEventProperty("r6_rp", "1");
}

var id = portal.id();
var idx = Math.floor((id - 24) / 4) + 1;
var choice = ((id - 24) % 4) + 1;
var target = "np16";

function nextPortal(n) {
    return n < 10 ? "np0" + n : "np" + n;
}

if (idx >= 1 && idx <= 4) {
    target = plr.getEventProperty("r6way" + idx) === String(choice) ? nextPortal(idx - 1) : "np16";
} else if (idx >= 5 && idx <= 8) {
    target = plr.getEventProperty("r6way" + idx) === String(choice) ? nextPortal(idx - 1) : "np03";
} else if (idx >= 9 && idx <= 12) {
    target = plr.getEventProperty("r6way" + idx) === String(choice) ? nextPortal(idx - 1) : "np07";
} else if (idx >= 13 && idx <= 16) {
    target = plr.getEventProperty("r6way" + idx) === String(choice) ? nextPortal(idx - 1) : "np11";
}

portal.warp(plr.mapID(), target);
