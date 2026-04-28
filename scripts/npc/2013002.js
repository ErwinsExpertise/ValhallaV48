var itemArray = [
    2000004, 2000005, 2000006, 2002020, 2002021, 2002022, 2002023, 2002024, 2002025, 2002026,
    2001000, 2001002, 2002015, 2050005, 2022179, 2020014, 2020015,
    2100000, 2100001, 2100002, 2100003, 2100004, 2100005,
    2061003, 2060003, 2060004, 2061004,
    2070006, 2070005, 2070007, 2070004,
    2210000, 2210001, 2210002,
];

var itemQuan = [
    50, 20, 200, 200, 200, 200, 200, 200, 200, 200, 200, 200, 5, 30, 2, 100, 50,
    1, 1, 1, 1, 1, 1,
    2000, 2000, 2000, 2000,
    1, 1, 1, 1,
    5, 5, 5,
];

var pqItems = [4001022, 4001023];

if (plr.mapID() === 920010100) {
    npc.sendNext("Thank you for not only restoring the statue, but rescuing me, Minerva, from the entrapment. May the blessing of the goddess be with you till the end...");
    plr.finishEvent();
} else if (plr.mapID() === 920011300) {
    npc.sendNext("Thank you for not only restoring the statue, but rescuing me, Minerva, from the entrapment. May the blessing of the goddess be with you till the end...");
    for (var i = 0; i < pqItems.length; i++) {
        plr.removeAll(pqItems[i]);
    }
    var rand = Math.floor(Math.random() * itemArray.length);
    plr.gainItem(itemArray[rand], itemQuan[rand]);
    plr.warp(200080101);
}
