var mapId = plr.mapID();

if (mapId === 108000502 && plr.haveItem(4031856, 15)) {
    npc.sendNext("Ohhh... So you managed to gather up #r15 Potent Power Crystals#k! Wasn't it tough? That's amazing... Alright then, now let's talk about The Nautilus.");
    npc.sendNext("These crystals can only be used here, so I'll just take them back.");
    plr.warp(120000101);
} else if (mapId === 108000502) {
    npc.sendOk("You will have to collect me #v4031856##r15 Potent Power Crystals#k. Good luck.");
} else if (mapId === 108000500 && plr.haveItem(4031857, 15)) {
    npc.sendNext("Ohhh... So you managed to gather up #b15 Potent Wind Crystals#k! Wasn't it tough? That's amazing... Alright then, now let's talk about The Nautilus.");
    npc.sendNext("These crystals can only be used here, so I'll just take them back.");
    plr.warp(120000101);
} else if (mapId === 108000500) {
    npc.sendOk("You will have to collect me #v4031857##b15 Potent Wind Crystals#k. Good luck.");
} else {
    npc.sendOk("Something went wrong.");
}
