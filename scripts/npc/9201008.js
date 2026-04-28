var commonTicket = 5251001
var premiumTicket = 5251002
var mapId = plr.mapID()
var indoors = mapId >= 680000100 && mapId <= 680000500

if (indoors) {
    var stage = plr.weddingStage(false)
    if (stage < 0) {
        npc.sendOk("There is no active Chapel wedding session right now.")
    } else if (stage === 0) {
        npc.sendOk("Welcome to the Chapel waiting hall. Please remain here while the couple greets their guests.")
    } else if (npc.sendYesNo("The bride and groom are on their way to the altar. Would you like to proceed to the ceremony now?")) {
        if (plr.enterWeddingAsGuest(false)) {
            npc.sendOk("Please enjoy the Chapel wedding.")
        } else {
            npc.sendOk("I cannot move you to the ceremony right now.")
        }
    } else {
        npc.sendOk("Please wait until you are ready to head to the altar.")
    }
} else {
    var chapelOptions = ["How do I prepare a wedding?", "I have an engagement and want to arrange the wedding.", "I am a guest and want to enter the wedding.", "Make additional invitation cards."]
    var chapelText = "Welcome to the #bChapel#k. How can I help you?#b"
    for (var ch = 0; ch < chapelOptions.length; ch++) {
        chapelText += "\r\n#L" + ch + "#" + chapelOptions[ch] + "#l"
    }
    npc.sendSelection(chapelText)
    var choice = npc.selection()

    if (choice === 0) {
        npc.sendOk("To marry in the Chapel, you must first be engaged. Then one partner needs either a #b#t5251001##k or a #b#t5251002##k. After the reservation, both partners receive stacked invitation cards for their guests.")
    } else if (choice === 1) {
        var premium = plr.haveItem(premiumTicket, 1)
        var regular = plr.haveItem(commonTicket, 1)
        if (!premium && !regular) {
            npc.sendOk("You need a Chapel wedding reservation ticket first.")
        } else if (plr.reserveWedding(false, premium)) {
            npc.sendOk("Your Chapel wedding has been reserved. Both partners have received 15 invitation cards. Speak with #b#p9201012##k when you are ready to begin.")
        } else {
            npc.sendOk("I could not reserve the wedding. Make sure your partner is here, you are engaged, and both inventories have room for invitation cards.")
        }
    } else if (choice === 2) {
        if (plr.enterWeddingAsGuest(false)) {
            npc.sendOk("Please proceed inside and enjoy the ceremony.")
        } else {
            npc.sendOk("I cannot admit you right now. Make sure there is an active Chapel wedding and that you have the correct guest ticket.")
        }
    } else if (choice === 3) {
        if (!plr.hasWeddingReservation(false)) {
            npc.sendOk("Your couple does not currently have a Chapel reservation.")
        } else if (!plr.haveItem(5251100, 1)) {
            npc.sendOk("You need #b#t5251100##k if you want additional invitation cards.")
        } else if (!plr.canHold(plr.weddingInviteItem(false), 3)) {
            npc.sendOk("Please free an ETC slot before receiving more invitation cards.")
        } else {
            plr.gainItem(5251100, -1)
            plr.gainItem(plr.weddingInviteItem(false), 3)
            npc.sendOk("Here are three additional Chapel invitation cards.")
        }
    }
}
