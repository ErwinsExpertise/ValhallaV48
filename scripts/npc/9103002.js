var COUPON = 4001106;
var EXIT_MAP = 220000000;
var rewards = [
    { max: 300, id: 2040001, amount: 1 },
    { max: 600, id: 2040002, amount: 1 },
    { max: 900, id: 2040401, amount: 1 },
    { max: 1200, id: 2040402, amount: 1 },
    { max: 1500, id: 2040504, amount: 1 },
    { max: 1800, id: 2040505, amount: 1 },
    { max: 2100, id: 2040601, amount: 1 },
    { max: 2400, id: 2040602, amount: 1 },
    { max: 2700, id: 2040901, amount: 1 },
    { max: 3000, id: 2040902, amount: 1 },
    { max: 3200, id: 2041017, amount: 1 },
    { max: 3400, id: 2041020, amount: 1 },
    { max: 3700, id: 2041004, amount: 1 },
    { max: 4000, id: 2041005, amount: 1 },
    { max: 4300, id: 2040008, amount: 1 },
    { max: 4600, id: 2040009, amount: 1 },
    { max: 4900, id: 2040404, amount: 1 },
    { max: 5200, id: 2040405, amount: 1 },
    { max: 5500, id: 2040510, amount: 1 },
    { max: 5800, id: 2040511, amount: 1 },
    { max: 6100, id: 2040604, amount: 1 },
    { max: 6400, id: 2040605, amount: 1 },
    { max: 6700, id: 2040904, amount: 1 },
    { max: 7000, id: 2040905, amount: 1 },
    { max: 7300, id: 2041026, amount: 1 },
    { max: 7600, id: 2041027, amount: 1 },
    { max: 7900, id: 2041028, amount: 1 },
    { max: 8200, id: 2041029, amount: 1 },
    { max: 10200, id: 2020006, amount: 100 },
    { max: 13200, id: 2020007, amount: 100 },
    { max: 18200, id: 4031562, amount: 1 },
    { max: 23200, id: 2022019, amount: 50 },
    { max: 28200, id: 2020008, amount: 20 },
    { max: 33200, id: 2001001, amount: 5 },
    { max: 38200, id: 2000006, amount: 100 },
    { max: 43200, id: 2020009, amount: 100 },
    { max: 48990, id: 2022000, amount: 50 },
    { max: 53990, id: 2020010, amount: 20 },
    { max: 58990, id: 2001002, amount: 5 },
    { max: 63990, id: 2001000, amount: 50 },
    { max: 68990, id: 2000004, amount: 5 },
    { max: 73990, id: 2000005, amount: 1 },
    { max: 78990, id: 2030008, amount: 20 },
    { max: 83990, id: 2030009, amount: 20 },
    { max: 88990, id: 2000006, amount: 100 },
    { max: 88991, id: 1072263, amount: 1 },
    { max: 89991, id: 1032013, amount: 1 },
    { max: 89999, id: 1302016, amount: 1 },
    { max: 90000, id: 1332030, amount: 1 },
    { max: 90500, id: 1442017, amount: 1 },
    { max: 91000, id: 1322025, amount: 1 },
];

function rewardRoll() {
    var roll = Math.floor(Math.random() * 91000) + 1;

    for (let i = 0; i < rewards.length; i++) {
        if (roll <= rewards[i].max) {
            return rewards[i];
        }
    }

    return rewards[rewards.length - 1];
}

if (!npc.sendYesNo("Your party gave a stellar effort and gathered up at least 30 coupons. For that, I have a present for each and every one of you. After receiving the present, you will be sent back to Ludibrium. Now, would you like to receive the present right now?")) {
    npc.sendOk("If you wish to receive your rewards and return to Ludibrium, please let me know!");
} else {
    var reward = rewardRoll();

    if (!plr.canHold(reward.id, reward.amount)) {
        npc.sendOk("Please make sure your inventory has at least one spot available.");
    } else if (!plr.gainItem(reward.id, reward.amount)) {
        npc.sendOk("There seems to be a problem here. Please try again.");
    } else {
        plr.logEvent("lmpq: reward granted item=" + reward.id + " amount=" + reward.amount);
        plr.removeAll(COUPON);
        plr.warp(EXIT_MAP);
    }
}
