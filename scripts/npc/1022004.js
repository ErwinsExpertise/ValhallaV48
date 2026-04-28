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

var type = menu("Um... Hi, I'm Mr. Thunder's apprentice. He's getting up there in age, so he handles most of the heavy-duty work while I handle some of the lighter jobs. What can I do for you?", [
    "Make a glove",
    "Upgrade a glove",
    "Create materials"
]);

var recipe;
var qty = 1;

if (type === 0) {
    var sel0 = menu("Okay, so which glove do you want me to make?", [
        "#t1082003#",
        "#t1082000#",
        "#t1082004#",
        "#t1082001#",
        "#t1082007#",
        "#t1082008#",
        "#t1082023#",
        "#t1082009#",
        "#t1082059#"
    ]);
    recipe = {
        item: [1082003, 1082000, 1082004, 1082001, 1082007, 1082008, 1082023, 1082009, 1082059][sel0],
        mats: [[4000021, 4011001], [4011001], [4000021, 4011000], [4011001], [4011000, 4011001, 4003000], [4000021, 4011001, 4003000], [4000021, 4011001, 4003000], [4011001, 4021007, 4000030, 4003000], [4011007, 4011000, 4011006, 4000030, 4003000]][sel0],
        qty: [[15, 1], [2], [40, 2], [2], [3, 2, 15], [30, 4, 15], [50, 5, 40], [3, 2, 30, 45], [1, 8, 2, 50, 50]][sel0],
        cost: [1000, 2000, 5000, 10000, 20000, 30000, 40000, 50000, 70000][sel0]
    };
} else if (type === 1) {
    var sel1 = menu("Upgrade a glove? That shouldn't be too difficult. Which did you have in mind?", [
        "#t1082005#",
        "#t1082006#",
        "#t1082035#",
        "#t1082036#",
        "#t1082024#",
        "#t1082025#",
        "#t1082010#",
        "#t1082011#",
        "#t1082060#",
        "#t1082061#"
    ]);
    recipe = {
        item: [1082005, 1082006, 1082035, 1082036, 1082024, 1082025, 1082010, 1082011, 1082060, 1082061][sel1],
        mats: [[1082007, 4011001], [1082007, 4011005], [1082008, 4021006], [1082008, 4021008], [1082023, 4011003], [1082023, 4021008], [1082009, 4011002], [1082009, 4011006], [1082059, 4011002, 4021005], [1082059, 4021007, 4021008]][sel1],
        qty: [[1, 1], [1, 2], [1, 3], [1, 1], [1, 4], [1, 2], [1, 5], [1, 4], [1, 3, 5], [1, 2, 2]][sel1],
        cost: [20000, 25000, 30000, 40000, 45000, 50000, 55000, 60000, 70000, 80000][sel1]
    };
} else {
    var sel2 = menu("Materials? I know of a few materials that I can make for you...", [
        "Make Processed Wood with Tree Branch",
        "Make Processed Wood with Firewood",
        "Make Screws (packs of 15)"
    ]);
    recipe = [
        { item: 4003001, mats: [4000003], qty: [10], cost: 0 },
        { item: 4003001, mats: [4000018], qty: [5], cost: 0 },
        { item: 4003000, mats: [4011000, 4011001], qty: [1, 1], cost: 0, outputQty: 15 }
    ][sel2];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
}

var finalQty = recipe.outputQty ? recipe.outputQty * qty : qty;

if (!npc.sendYesNo(buildPrompt(recipe, qty))) {
    npc.sendOk("Talk to me again anytime.");
} else if (!plr.canHold(recipe.item, finalQty)) {
    npc.sendOk("Check your inventory for a free slot first.");
} else if (plr.getMesos() < recipe.cost * qty) {
    npc.sendOk("I may still be an apprentice, but I do need to earn a living.");
} else if (!hasMaterials(recipe, qty)) {
    npc.sendOk("I'm still an apprentice, I don't know if I can substitute other items in yet... Can you please bring what the recipe calls for?");
} else {
    takeMaterials(recipe, qty);
    if (recipe.cost > 0) {
        plr.gainMesos(-(recipe.cost * qty));
    }
    plr.gainItem(recipe.item, finalQty);
    npc.sendOk("Did that come out right? Come by me again if you have anything for me to practice on.");
}
