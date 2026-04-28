var options = [
    { label: "Moonstone", item: 2240000, mats: [4011007, 4021007], qty: [1, 1], cost: 30000 },
    { label: "Star Gem", item: 2240001, mats: [4021009, 4021007], qty: [1, 1], cost: 20000 },
    { label: "Golden Heart", item: 2240002, mats: [4011006, 4021007], qty: [1, 1], cost: 10000 },
    { label: "Silver Swan", item: 2240003, mats: [4011004, 4021007], qty: [1, 1], cost: 5000 }
]

if (!plr.hasEngagementBox()) {
    var text = "I'm #p9201000#, the #bengagement ring maker#k. What kind of engagement ring box do you want me to craft?#b"
    for (var i = 0; i < options.length; i++) {
        text += "\r\n#L" + i + "#" + options[i].label + "#l"
    }
    npc.sendSelection(text)
    var ringSel = npc.selection()
    var ring = options[ringSel]
    var prompt = "Then I'm going to craft you a #b#t" + ring.item + "##k. In that case, I'm going to need specific items from you in order to make it. Make sure you have room in your inventory, though!#b"
    for (var j = 0; j < ring.mats.length; j++) {
        prompt += "\r\n#i" + ring.mats[j] + "# " + ring.qty[j] + " #t" + ring.mats[j] + "#"
    }
    prompt += "\r\n#i4031138# " + ring.cost + " meso"

    if (!npc.sendYesNo(prompt)) {
        npc.sendOk("Changed your mind? Come back if you decide you want a ring box.")
    } else if (!plr.canHold(ring.item, 1)) {
        npc.sendOk("Check your inventory for a free ETC slot first.")
    } else if (plr.getMesos() < ring.cost) {
        npc.sendOk("I'm sorry but there's a fee for my services. Please bring me the right amount of mesos first.")
    } else if (!plr.haveItem(ring.mats[0], ring.qty[0]) || !plr.haveItem(ring.mats[1], ring.qty[1])) {
        npc.sendOk("Hm, it seems you're lacking some ingredients for that engagement ring box.")
    } else {
        plr.gainItem(ring.mats[0], -ring.qty[0])
        plr.gainItem(ring.mats[1], -ring.qty[1])
        plr.gainMesos(-ring.cost)
        plr.gainItem(ring.item, 1)
        npc.sendOk("All done. Use the ring box on the person you wish to propose to while both of you are on the same map.")
    }
} else {
    if (!npc.sendYesNo("You already have an engagement ring box. Do you want me to discard it for you?")) {
        npc.sendOk("All right. Come back if you need another ring box later.")
    } else {
        for (var k = 2240000; k <= 2240003; k++) {
            if (plr.haveItem(k, 1)) {
                plr.removeAll(k)
            }
        }
        npc.sendOk("Your ring box has been discarded.")
    }
}
