var bgPrizes = [[2022241, 5], [2022242, 5], [2022273, 3], [2022274, 3], [2022275, 3]]
var cmPrizes = [[2000005, 20], [2020012, 10], [2022015, 5], [2022027, 5], [2049100, 1]]

var action = (plr.haveItem(4031424, 1) ? 1 : (plr.haveItem(4031423, 1) ? 2 : 0))

if (action === 0) {
    npc.sendOk("Welcome to Amoria's Wedding Gift Registry. Onyx Chest exchanges are available here. Full wedding gift registry and wishlist support are still limited on this server.")
} else if (action === 1) {
    if (!plr.isMarried()) {
        npc.sendOk("You must be married to claim the Bride and Groom Onyx Chest prize.")
    } else if (!plr.canHold(bgPrizes[0][0], bgPrizes[0][1])) {
        npc.sendOk("Please free a USE slot before opening the chest.")
    } else {
        var rand = Math.floor(Math.random() * bgPrizes.length)
        plr.gainItem(4031424, -1)
        plr.gainItem(bgPrizes[rand][0], bgPrizes[rand][1])
        npc.sendOk("Enjoy your prize.")
    }
} else {
    if (!plr.canHold(cmPrizes[0][0], cmPrizes[0][1])) {
        npc.sendOk("Please free a USE slot before opening the chest.")
    } else {
        var rand2 = Math.floor(Math.random() * cmPrizes.length)
        plr.gainItem(4031423, -1)
        plr.gainItem(cmPrizes[rand2][0], cmPrizes[rand2][1])
        npc.sendOk("Enjoy your prize.")
    }
}
