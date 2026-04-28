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

var type = menu("Hm? Who might you be? Oh, you've heard about my forging skills? In that case, I'd be glad to process some of your ores... for a fee.", [
    "Refine a mineral ore",
    "Refine a jewel ore",
    "Upgrade a helmet",
    "Upgrade a shield"
]);

var qty = 1;
var recipe;

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
    var sel2 = menu("Ah, you wish to upgrade a helmet? Then tell me, which one?", [
        "#t1002042#", "#t1002041#", "#t1002002#", "#t1002044#", "#t1002003#", "#t1002040#", "#t1002007#", "#t1002052#", "#t1002011#", "#t1002058#", "#t1002009#", "#t1002056#",
        "#t1002087#", "#t1002088#", "#t1002050#", "#t1002049#", "#t1002047#", "#t1002048#", "#t1002099#", "#t1002098#", "#t1002085#", "#t1002028#", "#t1002022#", "#t1002101#"
    ]);
    recipe = {
        item: [1002042, 1002041, 1002002, 1002044, 1002003, 1002040, 1002007, 1002052, 1002011, 1002058, 1002009, 1002056, 1002087, 1002088, 1002050, 1002049, 1002047, 1002048, 1002099, 1002098, 1002085, 1002028, 1002022, 1002101][sel2],
        mats: [[1002001, 4011002], [1002001, 4021006], [1002043, 4011001], [1002043, 4011002], [1002039, 4011001], [1002039, 4011002], [1002051, 4011001], [1002051, 4011002], [1002059, 4011001], [1002059, 4011002], [1002055, 4011001], [1002055, 4011002], [1002027, 4011002], [1002027, 4011006], [1002005, 4011005], [1002005, 4011006], [1002004, 4021000], [1002004, 4021005], [1002021, 4011002], [1002021, 4011006], [1002086, 4011002], [1002086, 4011004], [1002100, 4011007, 4011001], [1002100, 4011007, 4011002]][sel2],
        qty: [[1, 1], [1, 1], [1, 1], [1, 1], [1, 1], [1, 1], [1, 2], [1, 2], [1, 3], [1, 3], [1, 3], [1, 3], [1, 4], [1, 4], [1, 5], [1, 5], [1, 3], [1, 3], [1, 5], [1, 6], [1, 5], [1, 4], [1, 1, 7], [1, 1, 7]][sel2],
        cost: [500, 300, 500, 800, 500, 800, 1000, 1500, 1500, 2000, 1500, 2000, 2000, 4000, 4000, 5000, 8000, 10000, 12000, 15000, 20000, 25000, 30000, 30000][sel2]
    };
} else {
    var sel3 = menu("Ah, you wish to upgrade a shield? Then tell me, which one?", ["#t1092014#", "#t1092013#", "#t1092010#", "#t1092011#"]);
    recipe = {
        item: [1092014, 1092013, 1092010, 1092011][sel3],
        mats: [[1092012, 4011003], [1092012, 4011002], [1092009, 4011007, 4011004], [1092009, 4011007, 4011003]][sel3],
        qty: [[1, 10], [1, 10], [1, 1, 15], [1, 1, 15]][sel3],
        cost: [100000, 100000, 120000, 120000][sel3]
    };
}

if (!npc.sendYesNo(buildPrompt(recipe, qty))) {
    npc.sendOk("All right. Come back when you're ready.");
} else if (!plr.canHold(recipe.item, qty)) {
    npc.sendOk("Check your inventory for a free slot first.");
} else if (plr.getMesos() < recipe.cost * qty) {
    npc.sendOk("I'm afraid you cannot afford my services.");
} else if (!hasMaterials(recipe, qty)) {
    npc.sendOk("I'm afraid you're missing something for the item you want. See you another time, yes?");
} else {
    takeMaterials(recipe, qty);
    plr.gainMesos(-(recipe.cost * qty));
    plr.gainItem(recipe.item, qty);
    npc.sendOk("There, finished. What do you think, a piece of art, isn't it? Well, if you need anything else, you know where to find me.");
}
