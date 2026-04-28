var text = "Pelvis Bebop here, baby. What do you need?#b\r\n#L0#What is happening right now?#l\r\n#L1#We are ready for the vows.#l"
npc.sendSelection(text)
var selection = npc.selection()

if (selection === 0) {
    npc.sendOk("This is the Chapel altar. Once the ceremony is underway, guests can offer blessings and the couple can proceed to the party after the vows.")
} else if (selection === 1) {
    if (plr.partnerID() <= 0) {
        npc.sendOk("Only the engaged couple can proceed with the vows.")
    } else if (plr.completeWedding(false)) {
        npc.sendOk("Let's shake this place up, baby. You two are now officially married.")
    } else {
        npc.sendOk("I can't proceed yet. Make sure both partners are here at the altar with their engagement rings.")
    }
}
