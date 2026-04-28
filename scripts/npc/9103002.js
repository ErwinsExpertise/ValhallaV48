var itemArray = [
    2000004, 2000005, 2000006, 2002020, 2002021, 2002022, 2002023, 2002024, 2002025, 2002026,
    2001000, 2001002, 2002015, 2050005, 2020014, 2020015,
    2100001, 2100002, 2100003, 2100004, 2100005,
    2061003, 2060003, 2060004, 2061004,
];

var itemQuan = [
    50, 20, 200, 200, 200, 200, 200, 200, 200, 200,
    200, 200, 5, 30, 2, 100,
    50, 1, 1, 1, 1,
    1, 2000, 2000, 2000, 2000,
];

npc.sendNext("Wow, you did it! Here is a reward for your work. Well done!");

var rand = Math.floor(Math.random() * itemArray.length);
plr.gainItem(itemArray[rand], itemQuan[rand]);

if (rand > 3) {
    rand = 3;
}

plr.warp(220000000);
