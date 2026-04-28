function menu(prompt, options) {
    var text = prompt + "#b";
    for (var i = 0; i < options.length; i++) {
        text += "\r\n#L" + i + "# " + options[i] + "#l";
    }
    npc.sendSelection(text);
    return npc.selection();
}

function buildPrompt(recipe) {
    var text = "You want me to make a #t" + recipe.item + "#? In that case, I'm going to need specific items from you in order to make it. Make sure you have room in your inventory, though!#b";
    for (var i = 0; i < recipe.mats.length; i++) {
        text += "\r\n#i" + recipe.mats[i] + "# " + recipe.qty[i] + " #t" + recipe.mats[i] + "#";
    }
    text += "\r\n#i4031138# " + recipe.cost + " meso";
    return text;
}

function hasMaterials(recipe) {
    for (var i = 0; i < recipe.mats.length; i++) {
        if (!plr.haveItem(recipe.mats[i], recipe.qty[i])) {
            return false;
        }
    }
    return true;
}

function takeMaterials(recipe) {
    for (var i = 0; i < recipe.mats.length; i++) {
        plr.gainItem(recipe.mats[i], -recipe.qty[i]);
    }
}

var type = menu("Hello there. I'm Orbis' number one glove maker. Would you like me to make you something?", [
    "Create or upgrade a Warrior glove",
    "Create or upgrade a Bowman glove",
    "Create or upgrade a Magician glove",
    "Create or upgrade a Thief glove"
]);

var recipe;

if (type === 0) {
    var sel0 = menu("Warrior glove? Okay, then which one?", ["#t1082103#", "#t1082104#", "#t1082105#", "#t1082114#", "#t1082115#", "#t1082116#", "#t1082117#"]);
    recipe = [
        { item: 1082103, mats: [4005000, 4011000, 4011006, 4000030, 4003000], qty: [2, 8, 3, 70, 55], cost: 90000 },
        { item: 1082104, mats: [1082103, 4011002, 4021006], qty: [1, 6, 4], cost: 90000 },
        { item: 1082105, mats: [1082103, 4021006, 4021008], qty: [1, 8, 3], cost: 100000 },
        { item: 1082114, mats: [4005000, 4005002, 4021005, 4000030, 4003000], qty: [2, 1, 8, 90, 60], cost: 100000 },
        { item: 1082115, mats: [1082114, 4005000, 4005002, 4021003], qty: [1, 1, 1, 7], cost: 110000 },
        { item: 1082116, mats: [1082114, 4005002, 4021000], qty: [1, 3, 8], cost: 110000 },
        { item: 1082117, mats: [1082114, 4005000, 4005002, 4021008], qty: [1, 2, 1, 4], cost: 120000 }
    ][sel0];
} else if (type === 1) {
    var sel1 = menu("Bowman glove? Okay, then which one?", ["#t1082106#", "#t1082107#", "#t1082108#", "#t1082109#", "#t1082110#", "#t1082111#", "#t1082112#"]);
    recipe = [
        { item: 1082106, mats: [4005002, 4021005, 4011004, 4000030, 4003000], qty: [2, 8, 3, 70, 55], cost: 90000 },
        { item: 1082107, mats: [1082106, 4021006, 4011006], qty: [1, 5, 3], cost: 90000 },
        { item: 1082108, mats: [1082106, 4021007, 4021008], qty: [1, 2, 3], cost: 100000 },
        { item: 1082109, mats: [4005002, 4005000, 4021000, 4000030, 4003000], qty: [2, 1, 8, 90, 60], cost: 100000 },
        { item: 1082110, mats: [1082109, 4005002, 4005000, 4021005], qty: [1, 1, 1, 7], cost: 110000 },
        { item: 1082111, mats: [1082109, 4005002, 4005000, 4021003], qty: [1, 1, 1, 7], cost: 110000 },
        { item: 1082112, mats: [1082109, 4005002, 4005000, 4021008], qty: [1, 2, 1, 4], cost: 120000 }
    ][sel1];
} else if (type === 2) {
    var sel2 = menu("Magician glove? Okay, then which one?", ["#t1082098#", "#t1082099#", "#t1082100#", "#t1082121#", "#t1082122#", "#t1082123#"]);
    recipe = [
        { item: 1082098, mats: [4005001, 4011000, 4011004, 4000030, 4003000], qty: [2, 6, 6, 70, 55], cost: 90000 },
        { item: 1082099, mats: [1082098, 4021002, 4021007], qty: [1, 6, 2], cost: 90000 },
        { item: 1082100, mats: [1082098, 4021008, 4011006], qty: [1, 3, 3], cost: 100000 },
        { item: 1082121, mats: [4005001, 4005003, 4021003, 4000030, 4003000], qty: [2, 1, 8, 90, 60], cost: 100000 },
        { item: 1082122, mats: [1082121, 4005001, 4005003, 4021005], qty: [1, 1, 1, 7], cost: 110000 },
        { item: 1082123, mats: [1082121, 4005001, 4005003, 4021008], qty: [1, 2, 1, 4], cost: 120000 }
    ][sel2];
} else {
    var sel3 = menu("Thief glove? Okay, then which one?", ["#t1082095#", "#t1082096#", "#t1082097#", "#t1082118#", "#t1082119#", "#t1082120#"]);
    recipe = [
        { item: 1082095, mats: [4005003, 4011000, 4011003, 4000030, 4003000], qty: [2, 6, 6, 70, 55], cost: 90000 },
        { item: 1082096, mats: [1082095, 4011004, 4021007], qty: [1, 6, 2], cost: 90000 },
        { item: 1082097, mats: [1082095, 4021007, 4011006], qty: [1, 3, 3], cost: 100000 },
        { item: 1082118, mats: [4005003, 4005002, 4011002, 4000030, 4003000], qty: [2, 1, 8, 90, 60], cost: 100000 },
        { item: 1082119, mats: [1082118, 4005003, 4005002, 4021001], qty: [1, 1, 1, 7], cost: 110000 },
        { item: 1082120, mats: [1082118, 4005003, 4005002, 4021000], qty: [1, 2, 1, 8], cost: 120000 }
    ][sel3];
}

if (!npc.sendYesNo(buildPrompt(recipe))) {
    npc.sendOk("All right. Come back when you're ready.");
} else if (!plr.canHold(recipe.item, 1)) {
    npc.sendOk("Check your inventory for a free slot first.");
} else if (plr.getMesos() < recipe.cost) {
    npc.sendOk("I'm afraid you cannot afford my services.");
} else if (!hasMaterials(recipe)) {
    npc.sendOk("I'm afraid that substitute items are unacceptable, if you want your gloves made properly.");
} else {
    takeMaterials(recipe);
    plr.gainMesos(-recipe.cost);
    plr.gainItem(recipe.item, 1);
    npc.sendOk("Done. If you need anything else, just ask again.");
}
