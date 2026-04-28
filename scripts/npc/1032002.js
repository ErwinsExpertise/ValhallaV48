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
    if (recipe.cost > 0) {
        text += "\r\n#i4031138# " + recipe.cost + " meso";
    }
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

var type = menu("Welcome to my eco-safe refining operation! What would you like today?", [
    "Make a glove",
    "Upgrade a glove",
    "Upgrade a hat",
    "Make a wand",
    "Make a staff"
]);

var recipe;

if (type === 0) {
    var sel0 = menu("So, what kind of glove would you like me to make?", ["#t1082019#", "#t1082020#", "#t1082026#", "#t1082051#", "#t1082054#", "#t1082062#", "#t1082081#", "#t1082086#"]);
    recipe = {
        item: [1082019, 1082020, 1082026, 1082051, 1082054, 1082062, 1082081, 1082086][sel0],
        mats: [[4000021], [4000021, 4011001], [4000021, 4011006], [4000021, 4021006, 4021000], [4000021, 4011006, 4011001, 4021000], [4000021, 4021000, 4021006, 4003000], [4021000, 4011006, 4000030, 4003000], [4011007, 4011001, 4021007, 4000030, 4003000]][sel0],
        qty: [[15], [30, 1], [50, 2], [60, 1, 2], [70, 1, 3, 2], [80, 3, 3, 30], [3, 2, 35, 40], [1, 8, 1, 50, 50]][sel0],
        cost: [7000, 15000, 20000, 25000, 30000, 40000, 50000, 70000][sel0]
    };
} else if (type === 1) {
    var sel1 = menu("So, what kind of glove are you looking to upgrade to?", ["#t1082021#", "#t1082022#", "#t1082027#", "#t1082028#", "#t1082052#", "#t1082053#", "#t1082055#", "#t1082056#", "#t1082063#", "#t1082064#", "#t1082082#", "#t1082080#", "#t1082087#", "#t1082088#"]);
    recipe = {
        item: [1082021, 1082022, 1082027, 1082028, 1082052, 1082053, 1082055, 1082056, 1082063, 1082064, 1082082, 1082080, 1082087, 1082088][sel1],
        mats: [[1082020, 4011001], [1082020, 4021001], [1082026, 4021000], [1082026, 4021008], [1082051, 4021005], [1082051, 4021008], [1082054, 4021005], [1082054, 4021008], [1082062, 4021002], [1082062, 4021008], [1082081, 4021002], [1082081, 4021008], [1082086, 4011004, 4011006], [1082086, 4021008, 4011006]][sel1],
        qty: [[1, 1], [1, 2], [1, 3], [1, 1], [1, 3], [1, 1], [1, 3], [1, 1], [1, 4], [1, 2], [1, 5], [1, 3], [1, 3, 5], [1, 2, 3]][sel1],
        cost: [20000, 25000, 30000, 40000, 35000, 40000, 40000, 45000, 45000, 50000, 55000, 60000, 70000, 80000][sel1]
    };
} else if (type === 2) {
    var sel2 = menu("A hat? Which one were you thinking of?", ["#t1002065#", "#t1002013#"]);
    recipe = {
        item: [1002065, 1002013][sel2],
        mats: [[1002064, 4011001], [1002064, 4011006]][sel2],
        qty: [[1, 3], [1, 3]][sel2],
        cost: [40000, 50000][sel2]
    };
} else if (type === 3) {
    var sel3 = menu("A wand, huh? Prefer the smaller weapon that fits in your pocket? Which type are you seeking?", ["#t1372005#", "#t1372006#", "#t1372002#", "#t1372004#", "#t1372003#", "#t1372001#", "#t1372000#", "#t1372007#"]);
    recipe = {
        item: [1372005, 1372006, 1372002, 1372004, 1372003, 1372001, 1372000, 1372007][sel3],
        mats: [[4003001], [4003001, 4000001], [4011001, 4000009, 4003000], [4011002, 4003002, 4003000], [4011002, 4021002, 4003000], [4021006, 4011002, 4011001, 4003000], [4021006, 4021005, 4021007, 4003003, 4003000], [4011006, 4021003, 4021007, 4021002, 4003002, 4003000]][sel3],
        qty: [[5], [10, 50], [1, 30, 5], [2, 1, 10], [3, 1, 10], [5, 3, 1, 15], [5, 5, 1, 1, 20], [4, 3, 2, 1, 1, 30]][sel3],
        cost: [1000, 3000, 5000, 12000, 30000, 60000, 120000, 200000][sel3]
    };
} else {
    var sel4 = menu("Ah, a staff, a great symbol of one's power! Which are you looking to make?", ["#t1382000#", "#t1382003#", "#t1382005#", "#t1382004#", "#t1382002#", "#t1382001#"]);
    recipe = {
        item: [1382000, 1382003, 1382005, 1382004, 1382002, 1382001][sel4],
        mats: [[4003001], [4021005, 4011001, 4003000], [4021003, 4011001, 4003000], [4003001, 4011001, 4003000], [4021006, 4021001, 4011001, 4003000], [4011001, 4021006, 4021001, 4021005, 4003000, 4000010, 4003003]][sel4],
        qty: [[5], [1, 1, 5], [1, 1, 5], [50, 1, 10], [2, 1, 1, 15], [8, 5, 5, 5, 30, 50, 1]][sel4],
        cost: [2000, 2000, 2000, 5000, 12000, 180000][sel4]
    };
}

if (!npc.sendYesNo(buildPrompt(recipe))) {
    npc.sendOk("All right. Come back when you're ready.");
} else if (!plr.canHold(recipe.item, 1)) {
    npc.sendOk("Check your inventory for a free slot first.");
} else if (plr.getMesos() < recipe.cost) {
    npc.sendOk("Sorry, but all of us need money to live. Come back when you can pay my fees, yes?");
} else if (!hasMaterials(recipe)) {
    npc.sendOk("Uhm... I don't keep extra material on me. Sorry.");
} else {
    takeMaterials(recipe);
    if (recipe.cost > 0) {
        plr.gainMesos(-recipe.cost);
    }
    plr.gainItem(recipe.item, 1);
    npc.sendOk("It's a success! Oh, I've never felt so alive! Please come back again!");
}
