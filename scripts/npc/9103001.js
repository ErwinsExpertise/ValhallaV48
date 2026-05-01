var INSTANCE = 1;
var COUPON = 4001106;
var ENTRY_MAP = 809050000;
var PQ_MAPS = [809050000, 809050001, 809050002, 809050003, 809050004, 809050005, 809050006, 809050007, 809050008, 809050009, 809050010, 809050011, 809050012, 809050013, 809050014, 809050015, 809050016, 809050017];

function allPresentHere() {
    return plr.partyMembersOnMapCount() === plr.partyMembersCount();
}

function hasCoupons(members) {
    for (let i = 0; i < members.length; i++) {
        if (members[i].itemCount(COUPON) > 0) {
            return true;
        }
    }

    return false;
}

function badLevel(members) {
    for (let i = 0; i < members.length; i++) {
        var level = members[i].level();
        if (level < 51 || level > 70) {
            return true;
        }
    }

    return false;
}

function mazeOccupied() {
    for (let i = 0; i < PQ_MAPS.length; i++) {
        if (map.getMap(PQ_MAPS[i], INSTANCE).playerCount() > 0) {
            return true;
        }
    }

    return false;
}

var action = npc.sendSelection("This is the entrance to the Ludibrium Maze. Enjoy!#b\r\n#L0#Enter the Ludibrium Maze.#l\r\n#L1#What is the Ludibrium Maze?#l#k");

if (action == 1) {
    npc.sendOk("This maze is available to all parties of 3 or more members, and all participants must be between Level 51~70. You will be given 15 minutes to escape the maze. At the center of the room, there will be a Warp Portal set up to transport you to a different room. These portals will transport you to other rooms where you'll hopefully find the exit. Pietri will be waiting at the exit, so all you need to do is talk to him, and he'll let you out. Break all the boxes located in the room, and a monster inside the box will drop a coupon. After escaping the maze, you will be awarded with EXP based on the coupons collected. Additionally, if the leader possesses at least 30 coupons, then a special gift will be presented to the party. If you cannot escape the maze within the allotted 15 minutes, you will receive 0 EXP for your time in the maze. If you decide to log off while you're in the maze, you will be automatically kicked out of the maze. Even if members of the party leave in the middle of the quest, the remaining members will be able to continue on with the quest. Your fighting spirit and wits will be tested! Good luck!");
} else if (!plr.inParty()) {
    npc.sendOk("Hmm...you're currently not affiliated with any party. You need to be in a party in order to tackle this maze.");
} else if (!plr.isLeader()) {
    npc.sendOk("Try taking on the Maze Quest with your party. If you do decide to tackle it, please have your Party Leader notify me!");
} else if (plr.partyQuestActive()) {
    npc.sendOk("Your party is already in the middle of a party quest.");
} else if (plr.partyMembersCount() < 3) {
    npc.sendOk("Your party needs to consist of at least 3 members in order to tackle this maze.");
} else if (!allPresentHere()) {
    npc.sendOk("I don't think all your party members are here right now. Please let me know when everyone is ready.");
} else {
    var members = plr.partyMembersOnMap();

    if (badLevel(members)) {
        npc.sendOk("One of the members of your party is not the required Level of 51 ~ 70. Please organize your party to match the required level.");
    } else if (hasCoupons(members)) {
        npc.sendOk("Someone in your party seems to be carrying a coupon. You may not enter the maze with a coupon.");
    } else if (mazeOccupied()) {
        npc.sendOk("A different party is currently exploring the maze. Please try again later!");
    } else {
        plr.startPartyQuestAt("ludibrium_maze_pq", ENTRY_MAP, INSTANCE);
    }
}
