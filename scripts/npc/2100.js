var mapId = plr.mapID();

if (mapId === 0 || mapId === 3) {
    if (!npc.sendYesNo("Welcome to the world of MapleStory. The purpose of this training camp is to help beginners. Would you like to enter this training camp? Some people start their journey without taking the training program. But I strongly recommend you take the training program first.")) {
        if (!npc.sendYesNo("Do you really want to start your journey right away?")) {
            npc.sendOk("Please talk to me again when you finally made your decision.");
        } else {
            plr.warp(40000);
        }
    } else {
        plr.warp(1);
    }
} else {
    npc.sendNext("This is the image room where your first training program begins. In this room, you can preview the job of your choice.");
    npc.sendNext("Once you train hard enough, you will be able to choose a job. You can become a Bowman in Henesys, a Magician in Ellinia, a Warrior in Perion, or a Thief in Kerning City.");
}
