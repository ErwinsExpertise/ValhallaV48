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

if (!plr.questCompleted(8225)) {
    npc.sendOk("Step aside, novice, we're doing business here.");
} else {
    var forgeRecipes = [
        { item: 2070018, mats: [4032015, 4032016, 4032017, 4021008, 4032005], qty: [1, 1, 1, 100, 30], cost: 70000 },
        { item: 1382060, mats: [4032016, 4032017, 4032004, 4032005, 4032012, 4005001], qty: [1, 1, 400, 10, 30, 4], cost: 70000 },
        { item: 1442068, mats: [4032015, 4032017, 4032004, 4032005, 4032012, 4005000], qty: [1, 1, 500, 40, 20, 4], cost: 70000 },
        { item: 1452060, mats: [4032015, 4032016, 4032004, 4032005, 4032012, 4005002], qty: [1, 1, 300, 75, 10, 4], cost: 70000 }
    ];
    var upgradeRecipes = [
        { item: 1472074, mats: [4032017, 4005001, 4021008], qty: [1, 10, 20], cost: 75000 },
        { item: 1472073, mats: [4032015, 4005002, 4021008], qty: [1, 10, 30], cost: 50000 },
        { item: 1472075, mats: [4032016, 4005000, 4021008], qty: [1, 5, 20], cost: 50000 },
        { item: 1332079, mats: [4032017, 4005001, 4021008], qty: [1, 10, 20], cost: 75000 },
        { item: 1332078, mats: [4032015, 4005002, 4021008], qty: [1, 10, 30], cost: 50000 },
        { item: 1332080, mats: [4032016, 4005000, 4021008], qty: [1, 5, 20], cost: 50000 },
        { item: 1462054, mats: [4032017, 4005001, 4021008], qty: [1, 10, 20], cost: 75000 },
        { item: 1462053, mats: [4032015, 4005002, 4021008], qty: [1, 10, 30], cost: 50000 },
        { item: 1462055, mats: [4032016, 4005000, 4021008], qty: [1, 5, 20], cost: 50000 },
        { item: 1402050, mats: [4032017, 4005001, 4021008], qty: [1, 10, 20], cost: 75000 },
        { item: 1402049, mats: [4032015, 4005002, 4021008], qty: [1, 10, 30], cost: 50000 },
        { item: 1402051, mats: [4032016, 4005000, 4021008], qty: [1, 5, 20], cost: 50000 }
    ];

    var type = npc.sendMenu("Hey, partner! If you have the right goods, I can turn it into something very nice...", "Weapon Forging", "Weapon Upgrading");
    var recipes = type === 0 ? forgeRecipes : upgradeRecipes;
    var labels = [];

    for (var i = 0; i < recipes.length; i++) {
        labels.push("#t" + recipes[i].item + "#");
    }

    var selection = menu(type === 0 ? "So, what kind of weapon would you like me to forge?" : "An upgraded weapon? Of course, but note that upgrades won't carry over to the new item...", labels);
    var recipe = recipes[selection];

    if (!npc.sendYesNo(buildPrompt(recipe))) {
        npc.sendOk("All right. Come back when you're ready.");
    } else if (!plr.canHold(recipe.item, 1)) {
        npc.sendOk("Check your inventory for a free slot first.");
    } else if (plr.getMesos() < recipe.cost) {
        npc.sendOk("I am afraid you don't have enough to pay me, partner. Please check this out first, ok?");
    } else if (!hasMaterials(recipe)) {
        npc.sendOk("Hey, I need those items to craft properly, you know?");
    } else {
        for (var j = 0; j < recipe.mats.length; j++) {
            plr.gainItem(recipe.mats[j], -recipe.qty[j]);
        }
        plr.gainMesos(-recipe.cost);
        plr.gainItem(recipe.item, 1);
        npc.sendOk("All done. If you need anything else... Well, I'm not going anywhere.");
    }
}

function menu(prompt, options) {
    var text = prompt + "#b";
    for (var i = 0; i < options.length; i++) {
        text += "\r\n#L" + i + "# " + options[i] + "#l";
    }
    npc.sendSelection(text);
    return npc.selection();
}
