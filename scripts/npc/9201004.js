var rings = [1112803, 1112806, 1112807, 1112809]
var divorceFee = 500000

function currentWeddingRing() {
    for (var i = 0; i < rings.length; i++) {
        if (plr.haveItem(rings[i], 1)) {
            return rings[i]
        }
    }
    return -1
}

var infoOptions = ["How can I engage someone?", "How can I marry?", "How can I divorce?", plr.isMarried() ? "I want a divorce now." : "I want my old wedding ring removed."]
var menuText = "Hello, welcome to #bAmoria#k. Do you have any questions about marriage?#b"
for (var m = 0; m < infoOptions.length; m++) {
    menuText += "\r\n#L" + m + "#" + infoOptions[m] + "#l"
}
npc.sendSelection(menuText)
var choice = npc.selection()

if (choice === 0) {
    npc.sendOk("The #bengagement process#k is simple. Speak with #b#p9201000##k, craft an engagement ring box, and use it on the person you want to propose to while both of you are on the same map.")
} else if (choice === 1) {
    npc.sendOk("Once you are engaged, choose either the #bCathedral#k or the #bChapel#k, buy the proper reservation ticket, and speak with the wedding assistant there to reserve the ceremony. After that, both partners must return to start the wedding together.")
} else if (choice === 2) {
    npc.sendOk("Divorce is permanent. If you truly wish to separate, speak with me again and choose the final divorce option. There is a #r" + divorceFee + " mesos#k fee, and a remarriage cooldown of about one week.")
} else if (choice === 3) {
    var ring = currentWeddingRing()
    if (!plr.isMarried() || ring < 0) {
        npc.sendOk("You do not currently have a wedding ring to remove.")
    } else if (plr.underMarriageCooldown()) {
        npc.sendOk("You are already under a remarriage cooldown. Please wait before attempting another marriage.")
    } else if (plr.getMesos() < divorceFee) {
        npc.sendOk("You need #r" + divorceFee + " mesos#k to pay the divorce fee.")
    } else if (!npc.sendYesNo("Divorce cannot be undone. Do you really want to break your marriage?")) {
        npc.sendOk("Very well. Think it over carefully.")
    } else {
        plr.gainMesos(-divorceFee)
        plr.breakMarriage(ring)
        npc.sendOk("The divorce has been processed.")
    }
}
