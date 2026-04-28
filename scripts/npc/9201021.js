var mapId = plr.mapID()

if (mapId === 680000401) {
    if (npc.sendYesNo("Do you want to go back to the hunting grounds? Returning here again will cost another 7 keys.")) {
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
        text += "\r\n#L1#I have 7 keys. Bring me to smash boxes.#l"
    }
    text += "\r\n#L2#Please warp me out.#l"
    npc.sendSelection(text)
    var selection = npc.selection()
    if (selection === 0) {
        if (!plr.haveItem(4000313, 1)) {
            npc.sendOk("It seems like you lost your #b#t4000313##k. I can't let you proceed without it.")
        } else {
            plr.warp(680000400)
        }
    } else if (selection === 1) {
        if (plr.haveItem(4031217, 7)) {
            plr.gainItem(4031217, -7)
            plr.warp(680000401)
        } else {
            npc.sendOk("You need 7 keys. Hunt the cakes and candles in the hunting grounds first.")
        }
    } else if (selection === 2) {
        plr.warp(680000500)
    }
}
