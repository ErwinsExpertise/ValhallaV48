var wishPrizes = [2000000, 2010004, 2020011, 2000004, 2000006, 2022015, 2000005, 1082174, 1002579, 1032039, 1002578, 1002580, 1002577, 1102078]
var wishPrizesQty = [10, 10, 5, 5, 5, 5, 10, 1, 1, 1, 1, 1, 1, 1]
var wishPrizesCost = [10, 15, 20, 30, 30, 50, 100, 400, 450, 500, 500, 530, 550, 600]

function getTierTicket(level) {
    if (level < 50) return 4031543
    if (level < 120) return 4031544
    return 4031545
}

var ticketId = getTierTicket(plr.level())
var ticketCount = plr.itemCount(ticketId)

npc.sendNext("Hi there, how is it going? Since you're passing by Amoria, have you heard about the #bAmorian Challenge#k? You can bring me #bWish Tickets#k from there and trade them for prizes.")

var menuText = "You currently have #b" + ticketCount + " #i" + ticketId + "# #t" + ticketId + "##k.\r\n\r\nPurchase a prize:#b"
for (var i = 0; i < wishPrizes.length; i++) {
    menuText += "\r\n#L" + i + "#" + wishPrizesQty[i] + " #z" + wishPrizes[i] + "##k - " + wishPrizesCost[i] + " wish tickets#l"
}

npc.sendSelection(menuText)
var sel = npc.selection()

if (sel < 0 || sel >= wishPrizes.length) {
    npc.sendOk("Changed your mind? That's fine. Come back any time.")
} else if (ticketCount < wishPrizesCost[sel]) {
    npc.sendOk("You will need #b" + wishPrizesCost[sel] + " #t" + ticketId + "##k to purchase that. Come back when you have enough tickets.")
} else if (!plr.canHold(wishPrizes[sel], wishPrizesQty[sel])) {
    npc.sendOk("Please have a free inventory slot available before claiming that item.")
} else if (!npc.sendYesNo("You have selected #b" + wishPrizesQty[sel] + " #z" + wishPrizes[sel] + "##k, which will require #b" + wishPrizesCost[sel] + " #t" + ticketId + "##k. Will you purchase it?")) {
    npc.sendOk("Changed your mind? That's fine. Come back any time.")
} else {
    plr.gainItem(ticketId, -wishPrizesCost[sel])
    plr.gainItem(wishPrizes[sel], wishPrizesQty[sel])
    npc.sendOk("There you go, have a good day!")
}
