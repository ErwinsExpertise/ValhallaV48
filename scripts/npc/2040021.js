function menu(prompt, options) {
    var text = prompt + "#b";
    for (var i = 0; i < options.length; i++) {
        text += "\r\n#L" + i + "# " + options[i] + "#l";
    }
    npc.sendSelection(text);
    return npc.selection();
}

function buildPrompt(recipe, stimId, stimulator) {
    var text = "You want me to make a #t" + recipe.item + "#? In that case, I'm going to need specific items from you in order to make it. Make sure you have room in your inventory, though!#b";
    if (stimulator) {
        text += "\r\n\r\n#rStimulator forging here keeps the original 10% destruction risk, but if it succeeds you'll receive the normal item rather than a random-stat reroll.#k";
        text += "\r\n#i" + stimId + "# 1 #t" + stimId + "#";
    }
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

var stimId = 4130001;
var rawType = menu("Hello, and welcome to the Ludibrium Shoe Store. How can I help you today?", [
    "What's a stimulator?",
    "Create Warrior shoes",
    "Create Bowman shoes",
    "Create Magician shoes",
    "Create Thief shoes",
    "Create Warrior shoes with a Stimulator",
    "Create Bowman shoes with a Stimulator",
    "Create Magician shoes with a Stimulator",
    "Create Thief shoes with a Stimulator"
]);

if (rawType === 0) {
    npc.sendOk("A stimulator is a special potion that I can add into the process of creating certain items. It gives it stats as though it had dropped from a monster. However, it is possible to have no change, and it is also possible for the item to be below average. There's also a 10% chance of not getting any item when using a stimulator, so please choose wisely.");
} else {
    var stimulator = rawType > 4;
    var type = stimulator ? rawType - 4 : rawType;
    var recipe;
    var sel;

    if (type === 1) {
        sel = menu("Warrior shoes? Sure thing, which kind?", ["#t1072003#", "#t1072039#", "#t1072040#", "#t1072041#", "#t1072002#", "#t1072112#", "#t1072113#", "#t1072000#", "#t1072126#", "#t1072127#", "#t1072132#", "#t1072133#", "#t1072134#", "#t1072135#"]);
        recipe = {
            item: [1072003, 1072039, 1072040, 1072041, 1072002, 1072112, 1072113, 1072000, 1072126, 1072127, 1072132, 1072133, 1072134, 1072135][sel],
            mats: [[4021003, 4011001, 4000021, 4003000], [4011002, 4011001, 4000021, 4003000], [4011004, 4011001, 4000021, 4003000], [4021000, 4011001, 4000021, 4003000], [4011001, 4021004, 4000021, 4000030, 4003000], [4011002, 4021004, 4000021, 4000030, 4003000], [4021008, 4021004, 4000021, 4000030, 4003000], [4011003, 4000021, 4000030, 4003000, 4000103], [4011005, 4021007, 4000030, 4003000, 4000104], [4011002, 4021007, 4000030, 4003000, 4000105], [4021008, 4011001, 4021003, 4000030, 4003000], [4021008, 4011001, 4011002, 4000030, 4003000], [4021008, 4011001, 4011005, 4000030, 4003000], [4021008, 4011001, 4011006, 4000030, 4003000]][sel],
            qty: [[4, 2, 45, 15], [4, 2, 45, 15], [4, 2, 45, 15], [4, 2, 45, 15], [3, 1, 30, 20, 25], [3, 1, 30, 20, 25], [2, 1, 30, 20, 25], [4, 100, 40, 30, 100], [4, 1, 40, 30, 100], [4, 1, 40, 30, 100], [1, 3, 6, 65, 45], [1, 3, 6, 65, 45], [1, 3, 6, 65, 45], [1, 3, 6, 65, 45]][sel],
            cost: [20000, 20000, 20000, 20000, 22000, 22000, 25000, 38000, 38000, 38000, 50000, 50000, 50000, 50000][sel] * 0.9
        };
    } else if (type === 2) {
        sel = menu("Bowman shoes? Sure thing, which kind?", ["#t1072079#", "#t1072080#", "#t1072081#", "#t1072082#", "#t1072083#", "#t1072101#", "#t1072102#", "#t1072103#", "#t1072118#", "#t1072119#", "#t1072120#", "#t1072121#", "#t1072122#", "#t1072123#", "#t1072124#", "#t1072125#"]);
        recipe = {
            item: [1072079, 1072080, 1072081, 1072082, 1072083, 1072101, 1072102, 1072103, 1072118, 1072119, 1072120, 1072121, 1072122, 1072123, 1072124, 1072125][sel],
            mats: [[4000021, 4021000, 4003000], [4000021, 4021005, 4003000], [4000021, 4021003, 4003000], [4000021, 4021004, 4003000], [4000021, 4021006, 4003000], [4021002, 4021006, 4000030, 4000021, 4003000], [4021003, 4021006, 4000030, 4000021, 4003000], [4021000, 4021006, 4000030, 4000021, 4003000], [4021000, 4003000, 4000030, 4000106], [4021006, 4003000, 4000030, 4000107], [4011003, 4003000, 4000030, 4000108], [4021002, 4003000, 4000030, 4000099], [4011001, 4021006, 4021008, 4000030, 4003000, 4000033], [4011001, 4021006, 4021008, 4000030, 4003000, 4000032], [4011001, 4021006, 4021008, 4000030, 4003000, 4000041], [4011001, 4021006, 4021008, 4000030, 4003000, 4000042]][sel],
            qty: [[50, 2, 15], [50, 2, 15], [50, 2, 15], [50, 2, 15], [50, 2, 15], [3, 1, 15, 30, 20], [3, 1, 15, 30, 20], [3, 1, 15, 30, 20], [4, 30, 45, 100], [4, 30, 45, 100], [5, 30, 45, 100], [5, 30, 45, 100], [3, 3, 1, 60, 35, 80], [3, 3, 1, 60, 35, 150], [3, 3, 1, 60, 35, 100], [3, 3, 1, 60, 35, 250]][sel],
            cost: [19000, 19000, 19000, 19000, 19000, 19000, 20000, 20000, 20000, 32000, 32000, 40000, 40000, 50000, 50000, 50000][sel] * 0.9
        };
    } else if (type === 3) {
        sel = menu("Magician shoes? Sure thing, which kind?", ["#t1072075#", "#t1072076#", "#t1072077#", "#t1072078#", "#t1072089#", "#t1072090#", "#t1072091#", "#t1072114#", "#t1072115#", "#t1072116#", "#t1072117#", "#t1072140#", "#t1072141#", "#t1072142#", "#t1072143#"]);
        recipe = {
            item: [1072075, 1072076, 1072077, 1072078, 1072089, 1072090, 1072091, 1072114, 1072115, 1072116, 1072117, 1072140, 1072141, 1072142, 1072143][sel],
            mats: [[4021000, 4000021, 4003000], [4021002, 4000021, 4003000], [4011004, 4000021, 4003000], [4021008, 4000021, 4003000], [4021001, 4021006, 4000021, 4000030, 4003000], [4021000, 4021006, 4000021, 4000030, 4003000], [4021008, 4021006, 4000021, 4000030, 4003000], [4021000, 4000030, 4000110, 4003000], [4021005, 4000030, 4000111, 4003000], [4011006, 4021007, 4000030, 4000100, 4003000], [4021008, 4021007, 4000030, 4000112, 4003000], [4021009, 4011006, 4021000, 4000030, 4003000], [4021009, 4011006, 4021005, 4000030, 4003000], [4021009, 4011006, 4021001, 4000030, 4003000], [4021009, 4011006, 4021003, 4000030, 4003000]][sel],
            qty: [[2, 50, 15], [2, 50, 15], [2, 50, 15], [1, 50, 15], [3, 1, 30, 15, 20], [3, 1, 30, 15, 20], [2, 1, 40, 25, 20], [4, 40, 100, 25], [4, 40, 100, 25], [2, 1, 40, 100, 25], [2, 1, 40, 100, 30], [1, 3, 3, 60, 40], [1, 3, 3, 60, 40], [1, 3, 3, 60, 40], [1, 3, 3, 60, 40]][sel],
            cost: [18000, 18000, 18000, 18000, 20000, 20000, 22000, 30000, 30000, 35000, 40000, 50000, 50000, 50000, 50000][sel] * 0.9
        };
    } else {
        sel = menu("Thief shoes? Sure thing, which kind?", ["#t1072032#", "#t1072033#", "#t1072035#", "#t1072036#", "#t1072104#", "#t1072105#", "#t1072106#", "#t1072107#", "#t1072108#", "#t1072109#", "#t1072110#", "#t1072128#", "#t1072130#", "#t1072129#", "#t1072131#"]);
        recipe = {
            item: [1072032, 1072033, 1072035, 1072036, 1072104, 1072105, 1072106, 1072107, 1072108, 1072109, 1072110, 1072128, 1072130, 1072129, 1072131][sel],
            mats: [[4011000, 4000021, 4003000], [4011001, 4000021, 4003000], [4011004, 4000021, 4003000], [4011006, 4000021, 4003000], [4021000, 4021004, 4000021, 4000030, 4003000], [4021003, 4021004, 4000021, 4000030, 4003000], [4021002, 4021004, 4000021, 4000030, 4003000], [4021000, 4000030, 4000113, 4003000], [4021003, 4000030, 4000095, 4003000], [4021006, 4000030, 4000096, 4003000], [4021005, 4000030, 4000097, 4003000], [4011007, 4021005, 4000030, 4000114, 4003000], [4011007, 4021000, 4000030, 4000115, 4003000], [4011007, 4021003, 4000030, 4000109, 4003000], [4011007, 4021001, 4000030, 4000036, 4003000]][sel],
            qty: [[3, 50, 15], [3, 50, 15], [2, 50, 15], [2, 50, 15], [3, 1, 30, 15, 20], [3, 1, 30, 15, 20], [3, 1, 30, 15, 20], [5, 45, 100, 30], [4, 45, 100, 30], [4, 45, 100, 30], [4, 45, 100, 30], [2, 3, 50, 100, 35], [2, 3, 50, 100, 35], [2, 3, 50, 100, 35], [2, 3, 50, 80, 35]][sel],
            cost: [19000, 19000, 19000, 21000, 20000, 20000, 20000, 40000, 32000, 35000, 35000, 50000, 50000, 50000, 50000][sel] * 0.9
        };
    }

    if (!npc.sendYesNo(buildPrompt(recipe, stimId, stimulator))) {
        npc.sendOk("All right. Come back when you're ready.");
    } else if (!plr.canHold(recipe.item, 1)) {
        npc.sendOk("Check your inventory for a free slot first.");
    } else if (plr.getMesos() < recipe.cost) {
        npc.sendOk("Sorry, we only accept meso.");
    } else if (!hasMaterials(recipe) || (stimulator && !plr.haveItem(stimId, 1))) {
        npc.sendOk("Sorry, but I have to have those items to get this exactly right. Perhaps next time.");
    } else {
        takeMaterials(recipe);
        plr.gainMesos(-recipe.cost);
        if (stimulator) {
            plr.gainItem(stimId, -1);
            if (Math.floor(Math.random() * 10) === 0) {
                npc.sendOk("Eek! I think I accidently added too much stimulator and, well, the whole thing is unusable now... Sorry, but I can't offer a refund.");
            } else {
                plr.gainItem(recipe.item, 1);
                npc.sendOk("There, the shoes are ready. Be careful, they're still hot. The stimulator kept its 10% destruction risk here, but successful crafts use the base item stats.");
            }
        } else {
            plr.gainItem(recipe.item, 1);
            npc.sendOk("There, the shoes are ready. Be careful, they're still hot.");
        }
    }
}
