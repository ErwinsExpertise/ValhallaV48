var mapId = plr.mapID()

if (mapId === 680000401) {
    if (npc.sendYesNo("Do you want to go back to the hunting grounds?")) {
        plr.warp(680000400)
    } else {
        npc.sendOk("Take your time.")
    }
} else {
    var text = "Hello, where would you like to go?#b"
    if (mapId !== 680000400) {
        text += "\r\n#L0#Untamed Hearts Hunting Ground#l"
    }
    if (mapId === 680000400) {
        text += "\r\n#L1#I have 5 keys. Bring me to smash boxes.#l"
    }
    text += "\r\n#L2#Please warp me out.#l"
    npc.sendSelection(text)
    var selection = npc.selection()
    if (selection === 0) {
        if (!plr.weddingIsPremium()) {
            npc.sendOk("Only Premium wedding parties may continue to the hunting grounds.")
        } else {
            plr.warp(680000400)
        }
    } else if (selection === 1) {
        if (plr.haveItem(4031409, 5)) {
            plr.gainItem(4031409, -5)
            plr.warp(680000401)
        } else {
            npc.sendOk("You need 5 #b#t4031409##k. Hunt the wedding monsters in the hunting grounds first.")
        }
    } else if (selection === 2) {
        plr.warp(680000500)
    }
}
