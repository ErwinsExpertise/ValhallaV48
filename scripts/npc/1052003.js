function menu(prompt, options) {
    var text = prompt + "#b";
    for (var i = 0; i < options.length; i++) {
        text += "\r\n#L" + i + "# " + options[i] + "#l";
    }
    npc.sendSelection(text);
    return npc.selection();
}

function buildPrompt(recipe, qty) {
    var text = "You want me to make " + (qty === 1 ? "a #t" + recipe.item + "#?" : qty + " #t" + recipe.item + "#?") + " In that case, I'm going to need specific items from you in order to make it. Make sure you have room in your inventory, though!#b";
    for (var i = 0; i < recipe.mats.length; i++) {
        text += "\r\n#i" + recipe.mats[i] + "# " + (recipe.qty[i] * qty) + " #t" + recipe.mats[i] + "#";
    }
    if (recipe.cost > 0) {
        text += "\r\n#i4031138# " + (recipe.cost * qty) + " meso";
    }
    return text;
}

function hasMaterials(recipe, qty) {
    for (var i = 0; i < recipe.mats.length; i++) {
        if (!plr.haveItem(recipe.mats[i], recipe.qty[i] * qty)) {
            return false;
        }
    }
    return true;
}

function takeMaterials(recipe, qty) {
    for (var i = 0; i < recipe.mats.length; i++) {
        plr.gainItem(recipe.mats[i], -(recipe.qty[i] * qty));
    }
}

var type = menu("Yes, I do own this forge. If you're willing to pay, I can offer you some of my services.", [
    "Refine a mineral ore",
    "Refine a jewel ore",
    "I have Iron Hog's Metal Hoof...",
    "Upgrade a claw"
]);

var recipe;
var qty = 1;

if (type === 0) {
    var sel0 = menu("So, what kind of mineral ore would you like to refine?", ["Bronze", "Steel", "Mithril", "Adamantium", "Silver", "Orihalcon", "Gold"]);
    recipe = [
        { item: 4011000, mats: [4010000], qty: [10], cost: 300 },
        { item: 4011001, mats: [4010001], qty: [10], cost: 300 },
        { item: 4011002, mats: [4010002], qty: [10], cost: 300 },
        { item: 4011003, mats: [4010003], qty: [10], cost: 500 },
        { item: 4011004, mats: [4010004], qty: [10], cost: 500 },
        { item: 4011005, mats: [4010005], qty: [10], cost: 500 },
        { item: 4011006, mats: [4010006], qty: [10], cost: 800 }
    ][sel0];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
} else if (type === 1) {
    var sel1 = menu("So, what kind of jewel ore would you like to refine?", ["Garnet", "Amethyst", "Aquamarine", "Emerald", "Opal", "Sapphire", "Topaz", "Diamond", "Black Crystal"]);
    recipe = [
        { item: 4021000, mats: [4020000], qty: [10], cost: 500 },
        { item: 4021001, mats: [4020001], qty: [10], cost: 500 },
        { item: 4021002, mats: [4020002], qty: [10], cost: 500 },
        { item: 4021003, mats: [4020003], qty: [10], cost: 500 },
        { item: 4021004, mats: [4020004], qty: [10], cost: 500 },
        { item: 4021005, mats: [4020005], qty: [10], cost: 500 },
        { item: 4021006, mats: [4020006], qty: [10], cost: 500 },
        { item: 4021007, mats: [4020007], qty: [10], cost: 1000 },
        { item: 4021008, mats: [4020008], qty: [10], cost: 3000 }
    ][sel1];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
} else if (type === 2) {
    if (!npc.sendYesNo("You know about that? Not many people realize the potential in the Iron Hog's Metal Hoof... I can make this into something special, if you want me to.")) {
        npc.sendOk("Come back if you change your mind.");
    } else {
        recipe = { item: 4011001, mats: [4000039], qty: [100], cost: 1000 };
        qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
        if (!npc.sendYesNo(buildPrompt(recipe, qty))) {
            npc.sendOk("Come back if you change your mind.");
        } else if (!plr.canHold(recipe.item, qty)) {
            npc.sendOk("Check your inventory for a free slot first.");
        } else if (plr.getMesos() < recipe.cost * qty) {
            npc.sendOk("Cash only, no credit.");
        } else if (!hasMaterials(recipe, qty)) {
            npc.sendOk("I cannot accept substitutes. If you don't have what I need, then I won't be able to help you.");
        } else {
            takeMaterials(recipe, qty);
            plr.gainMesos(-(recipe.cost * qty));
            plr.gainItem(recipe.item, qty);
            npc.sendOk("Phew... I almost didn't think that would work for a second... Well, I hope you enjoy it, anyway.");
        }
    }
} else {
    var sel3 = menu("Ah, you wish to upgrade a claw? Then tell me, which one?", ["#t1472023#", "#t1472024#", "#t1472025#"]);
    recipe = {
        item: [1472023, 1472024, 1472025][sel3],
        mats: [[1472022, 4011007, 4021000, 2012000], [1472022, 4011007, 4021005, 2012002], [1472022, 4011007, 4021008, 4000046]][sel3],
        qty: [[1, 1, 8, 10], [1, 1, 8, 10], [1, 1, 3, 5]][sel3],
        cost: [80000, 80000, 100000][sel3]
    };
    if (!npc.sendYesNo(buildPrompt(recipe, qty))) {
        npc.sendOk("Come back if you change your mind.");
    } else if (!plr.canHold(recipe.item, 1)) {
        npc.sendOk("Check your inventory for a free slot first.");
    } else if (plr.getMesos() < recipe.cost) {
        npc.sendOk("Cash only, no credit.");
    } else if (!hasMaterials(recipe, 1)) {
        npc.sendOk("I cannot accept substitutes. If you don't have what I need, then I won't be able to help you.");
    } else {
        takeMaterials(recipe, 1);
        plr.gainMesos(-recipe.cost);
        plr.gainItem(recipe.item, 1);
        npc.sendOk("Phew... I almost didn't think that would work for a second... Well, I hope you enjoy it, anyway.");
    }
}
