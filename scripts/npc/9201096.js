var craftedItem = 2022284;
var materials = [4032010, 4032011, 4032012];
var counts = [60, 60, 45];
var cost = 75000;

npc.sendNext("Hey, are you aware about the expeditions running right now at the Crimsonwood Keep? So, there is a great opportunity for one to improve themselves, one can rack up experience and loot pretty fast there.");

if (!npc.sendYesNo("Said so, methinks making use of some strong utility potions can potentially create some differential on the front, and by this I mean to start crafting #b#t2022284##k's to help on the efforts. So, getting right down to business, I'm currently pursuing #rplenty#k of those items: #r#t4032010##k, #r#t4032011##k, #r#t4032012##k, and some funds to support the cause. Would you want to get some of these boosters?")) {
    npc.sendOk("Very well, see you around.");
} else {
    var qty = npc.sendNumber("Ok, I'll be crafting some #t" + craftedItem + "#. In that case, how many of those do you want me to make?", 1, 1, 100);
    var totalCost = cost * qty;
    var prompt = "So, you want me to make " + (qty === 1 ? "a #t" + craftedItem + "#?" : qty + " #t" + craftedItem + "#?") + " In that case, I'm going to need specific items from you in order to make it. And make sure you have room in your inventory!#b";

    for (var i = 0; i < materials.length; i++) {
        prompt += "\r\n#i" + materials[i] + "# " + (counts[i] * qty) + " #t" + materials[i] + "#";
    }

    prompt += "\r\n#i4031138# " + totalCost + " meso";

    if (!npc.sendYesNo(prompt)) {
        npc.sendOk("Very well, see you around.");
    } else if (plr.getMesos() < totalCost) {
        npc.sendOk("Well, I DID say I would be needing some funds to craft it, wasn't it?");
    } else if (!plr.canHold(craftedItem, qty)) {
        npc.sendOk("You didn't check if you got a slot to spare on your inventory before crafting, right?");
    } else if (!plr.haveItem(materials[0], counts[0] * qty) || !plr.haveItem(materials[1], counts[1] * qty) || !plr.haveItem(materials[2], counts[2] * qty)) {
        npc.sendOk("There are not enough resources on your inventory. Please check it again.");
    } else {
        for (var j = 0; j < materials.length; j++) {
            plr.gainItem(materials[j], -(counts[j] * qty));
        }
        plr.gainMesos(-totalCost);
        plr.gainItem(craftedItem, qty);
        npc.sendOk("There it is! Thanks for your cooperation.");
    }
}
