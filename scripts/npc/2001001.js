var room = npc.askMenu(
    "I can send you to one of the giant Christmas tree rooms. Which room would you like to enter?#b",
    "The first tree room",
    "The second tree room",
    "The third tree room",
    "The fourth tree room",
    "The fifth tree room"
);

var target = 209000001 + room;
if (map.playerCount(target, 0) >= 6) {
    npc.sendOk("That room is already full right now. Please try one of the other tree rooms.");
} else {
    plr.warp(target);
}
