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

var type = menu("Pst... If you have the right goods, I can turn it into something niice...", [
    "Create a glove",
    "Upgrade a glove",
    "Create a claw",
    "Upgrade a claw",
    "Create materials"
]);

var recipe;
var qty = 1;

if (type === 0) {
    var sel0 = menu("So, what kind of glove would you like me to make?", ["#t1082002#", "#t1082029#", "#t1082030#", "#t1082031#", "#t1082032#", "#t1082037#", "#t1082042#", "#t1082046#", "#t1082075#", "#t1082065#", "#t1082092#"]);
    recipe = {
        item: [1082002, 1082029, 1082030, 1082031, 1082032, 1082037, 1082042, 1082046, 1082075, 1082065, 1082092][sel0],
        mats: [[4000021], [4000021, 4000018], [4000021, 4000015], [4000021, 4000020], [4011000, 4000021], [4011000, 4011001, 4000021], [4011001, 4000021, 4003000], [4011001, 4011000, 4000021, 4003000], [4021000, 4000014, 4000021, 4003000], [4021005, 4021008, 4000030, 4003000], [4011007, 4011000, 4021007, 4000030, 4003000]][sel0],
        qty: [[15], [30, 20], [30, 20], [30, 20], [2, 40], [2, 1, 10], [2, 50, 10], [3, 1, 60, 15], [3, 200, 80, 30], [3, 1, 40, 30], [1, 8, 1, 50, 50]][sel0],
        cost: [1000, 7000, 7000, 7000, 10000, 15000, 25000, 30000, 40000, 50000, 70000][sel0]
    };
} else if (type === 1) {
    var sel1 = menu("An upgraded glove? Sure thing, but note that upgrades won't carry over to the new item...", ["#t1082033#", "#t1082034#", "#t1082038#", "#t1082039#", "#t1082043#", "#t1082044#", "#t1082047#", "#t1082045#", "#t1082076#", "#t1082074#", "#t1082067#", "#t1082066#", "#t1082093#", "#t1082094#"]);
    recipe = {
        item: [1082033, 1082034, 1082038, 1082039, 1082043, 1082044, 1082047, 1082045, 1082076, 1082074, 1082067, 1082066, 1082093, 1082094][sel1],
        mats: [[1082032, 4011002], [1082032, 4021004], [1082037, 4011002], [1082037, 4021004], [1082042, 4011004], [1082042, 4011006], [1082046, 4011005], [1082046, 4011006], [1082075, 4011006], [1082075, 4021008], [1082065, 4021000], [1082065, 4011006, 4021008], [1082092, 4011001, 4000014], [1082092, 4011006, 4000027]][sel1],
        qty: [[1, 1], [1, 1], [1, 2], [1, 2], [1, 2], [1, 1], [1, 3], [1, 2], [1, 4], [1, 2], [1, 5], [1, 2, 1], [1, 7, 200], [1, 7, 150]][sel1],
        cost: [5000, 7000, 10000, 12000, 15000, 20000, 22000, 25000, 40000, 50000, 55000, 60000, 70000, 80000][sel1]
    };
} else if (type === 2) {
    var sel2 = menu("So, what kind of claw would you like me to make?", ["#t1472001#", "#t1472004#", "#t1472007#", "#t1472008#", "#t1472011#", "#t1472014#", "#t1472018#"]);
    recipe = {
        item: [1472001, 1472004, 1472007, 1472008, 1472011, 1472014, 1472018][sel2],
        mats: [[4011001, 4000021, 4003000], [4011000, 4011001, 4000021, 4003000], [1472000, 4011001, 4000021, 4003001], [4011000, 4011001, 4000021, 4003000], [4011000, 4011001, 4000021, 4003000], [4011000, 4011001, 4000021, 4003000], [4011000, 4011001, 4000030, 4003000]][sel2],
        qty: [[1, 20, 5], [2, 1, 30, 10], [1, 3, 20, 30], [3, 2, 50, 20], [4, 2, 80, 25], [3, 2, 100, 30], [4, 2, 40, 35]][sel2],
        cost: [2000, 3000, 5000, 15000, 30000, 40000, 50000][sel2]
    };
} else if (type === 3) {
    var sel3 = menu("An upgraded claw? Sure thing, but note that upgrades won't carry over to the new item...", ["#t1472002#", "#t1472003#", "#t1472005#", "#t1472006#", "#t1472009#", "#t1472010#", "#t1472012#", "#t1472013#", "#t1472015#", "#t1472016#", "#t1472017#", "#t1472019#", "#t1472020#"]);
    recipe = {
        item: [1472002, 1472003, 1472005, 1472006, 1472009, 1472010, 1472012, 1472013, 1472015, 1472016, 1472017, 1472019, 1472020][sel3],
        mats: [[1472001, 4011002], [1472001, 4011006], [1472004, 4011001], [1472004, 4011003], [1472008, 4011002], [1472008, 4011003], [1472011, 4011004], [1472011, 4021008], [1472014, 4021000], [1472014, 4011003], [1472014, 4021008], [1472018, 4021000], [1472018, 4021005]][sel3],
        qty: [[1, 1], [1, 1], [1, 2], [1, 2], [1, 3], [1, 3], [1, 4], [1, 1], [1, 5], [1, 5], [1, 2], [1, 6], [1, 6]][sel3],
        cost: [1000, 2000, 3000, 5000, 10000, 15000, 20000, 25000, 30000, 30000, 35000, 40000, 40000][sel3]
    };
} else {
    var sel4 = menu("Materials? I know of a few materials that I can make for you...", ["Make Processed Wood with Tree Branch", "Make Processed Wood with Firewood", "Make Screws (packs of 15)"]);
    recipe = [
        { item: 4003001, mats: [4000003], qty: [10], cost: 0 },
        { item: 4003001, mats: [4000018], qty: [5], cost: 0 },
        { item: 4003000, mats: [4011000, 4011001], qty: [1, 1], cost: 0, outputQty: 15 }
    ][sel4];
    qty = npc.sendNumber("So, you want me to make some #t" + recipe.item + "#s? In that case, how many do you want me to make?", 1, 1, 100);
}

var finalQty = recipe.outputQty ? recipe.outputQty * qty : qty;

if (!npc.sendYesNo(buildPrompt(recipe, qty))) {
    npc.sendOk("All right. Come back when you're ready.");
} else if (!plr.canHold(recipe.item, finalQty)) {
    npc.sendOk("Check your inventory for a free slot first.");
} else if (plr.getMesos() < recipe.cost * qty) {
    npc.sendOk("I'm afraid you cannot afford my services.");
} else if (!hasMaterials(recipe, qty)) {
    npc.sendOk("What are you trying to pull? I can't make anything unless you bring me what I ask for.");
} else {
    takeMaterials(recipe, qty);
    if (recipe.cost > 0) {
        plr.gainMesos(-(recipe.cost * qty));
    }
    plr.gainItem(recipe.item, finalQty);
    npc.sendOk("All done. If you need anything else... Well, I'm not going anywhere.");
}
