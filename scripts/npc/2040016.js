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

function outputQty(recipe, qty) {
    return recipe.outputQty ? recipe.outputQty * qty : qty;
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

var type = menu("Hm? Who might you be? Oh, you've heard about my forging skills? In that case, I'd be glad to process some of your ores... for a fee.", [
    "Refine a mineral ore",
    "Refine a jewel ore",
    "Refine a rare jewel",
    "Refine a crystal ore",
    "Create materials",
    "Create Arrows"
]);

var recipe;
var qty = 1;

if (type === 0) {
    var sel0 = menu("So, what kind of mineral ore would you like to refine?", ["Bronze", "Steel", "Mithril", "Adamantium", "Silver", "Orihalcon", "Gold"]);
    recipe = [
        { item: 4011000, mats: [4010000], qty: [10], cost: 270 },
        { item: 4011001, mats: [4010001], qty: [10], cost: 270 },
        { item: 4011002, mats: [4010002], qty: [10], cost: 270 },
        { item: 4011003, mats: [4010003], qty: [10], cost: 450 },
        { item: 4011004, mats: [4010004], qty: [10], cost: 450 },
        { item: 4011005, mats: [4010005], qty: [10], cost: 450 },
        { item: 4011006, mats: [4010006], qty: [10], cost: 720 }
    ][sel0];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
} else if (type === 1) {
    var sel1 = menu("So, what kind of jewel ore would you like to refine?", ["Garnet", "Amethyst", "Aquamarine", "Emerald", "Opal", "Sapphire", "Topaz", "Diamond", "Black Crystal"]);
    recipe = [
        { item: 4021000, mats: [4020000], qty: [10], cost: 450 },
        { item: 4021001, mats: [4020001], qty: [10], cost: 450 },
        { item: 4021002, mats: [4020002], qty: [10], cost: 450 },
        { item: 4021003, mats: [4020003], qty: [10], cost: 450 },
        { item: 4021004, mats: [4020004], qty: [10], cost: 450 },
        { item: 4021005, mats: [4020005], qty: [10], cost: 450 },
        { item: 4021006, mats: [4020006], qty: [10], cost: 450 },
        { item: 4021007, mats: [4020007], qty: [10], cost: 900 },
        { item: 4021008, mats: [4020008], qty: [10], cost: 2700 }
    ][sel1];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
} else if (type === 2) {
    var sel2 = menu("A rare jewel? Which one were you thinking of?", ["Moon Rock", "Star Rock"]);
    recipe = [
        { item: 4011007, mats: [4011000, 4011001, 4011002, 4011003, 4011004, 4011005, 4011006], qty: [1, 1, 1, 1, 1, 1, 1], cost: 9000 },
        { item: 4021009, mats: [4021000, 4021001, 4021002, 4021003, 4021004, 4021005, 4021006, 4021007, 4021008], qty: [1, 1, 1, 1, 1, 1, 1, 1, 1], cost: 13500 }
    ][sel2];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
} else if (type === 3) {
    var sel3 = menu("Crystal ore? I love refining those!", ["Power Crystal", "Wisdom Crystal", "DEX Crystal", "LUK Crystal"]);
    recipe = [
        { item: 4005000, mats: [4004000], qty: [10], cost: 4500 },
        { item: 4005001, mats: [4004001], qty: [10], cost: 4500 },
        { item: 4005002, mats: [4004002], qty: [10], cost: 4500 },
        { item: 4005003, mats: [4004003], qty: [10], cost: 4500 }
    ][sel3];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
} else if (type === 4) {
    var sel4 = menu("Materials? I know of a few materials that I can make for you...", ["Make Processed Wood with Tree Branch", "Make Processed Wood with Firewood", "Make Screws (packs of 15)"]);
    recipe = [
        { item: 4003001, mats: [4000003], qty: [10], cost: 0 },
        { item: 4003001, mats: [4000018], qty: [5], cost: 0 },
        { item: 4003000, mats: [4011000, 4011001], qty: [1, 1], cost: 0, outputQty: 15 }
    ][sel4];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
} else {
    var sel5 = menu("Arrows? Not a problem at all.", ["#t2060000#", "#t2061000#", "#t2060001#", "#t2061001#", "#t2060002#", "#t2061002#"]);
    recipe = [
        { item: 2060000, mats: [4003001, 4003004], qty: [1, 1], cost: 0, outputQty: 1000 },
        { item: 2061000, mats: [4003001, 4003004], qty: [1, 1], cost: 0, outputQty: 1000 },
        { item: 2060001, mats: [4011000, 4003001, 4003004], qty: [1, 3, 10], cost: 0, outputQty: 900 },
        { item: 2061001, mats: [4011000, 4003001, 4003004], qty: [1, 3, 10], cost: 0, outputQty: 900 },
        { item: 2060002, mats: [4011001, 4003001, 4003005], qty: [1, 5, 15], cost: 0, outputQty: 800 },
        { item: 2061002, mats: [4011001, 4003001, 4003005], qty: [1, 5, 15], cost: 0, outputQty: 800 }
    ][sel5];
}

var finalQty = outputQty(recipe, qty);

if (!npc.sendYesNo(buildPrompt(recipe, qty))) {
    npc.sendOk("All right. Come back when you're ready.");
} else if (!plr.canHold(recipe.item, finalQty)) {
    npc.sendOk("I'm afraid you have no slots available for this transaction.");
} else if (plr.getMesos() < recipe.cost * qty) {
    npc.sendOk("I'm afraid you cannot afford my services.");
} else if (!hasMaterials(recipe, qty)) {
    npc.sendOk("Hold it, I can't finish that without all of the proper materials. Bring them first, then we'll talk.");
} else {
    takeMaterials(recipe, qty);
    if (recipe.cost > 0) {
        plr.gainMesos(-(recipe.cost * qty));
    }
    plr.gainItem(recipe.item, finalQty);
    npc.sendOk("All done. If you need anything else, you know where to find me.");
}
