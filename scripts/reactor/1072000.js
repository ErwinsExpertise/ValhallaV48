function act() {
    if (!plr) {
        return;
    }

    if (plr.questData(7500) !== "p1" || (plr.job() !== 410 && plr.job() !== 420)) {
        plr.sendMessage("There is a crack here that leads to another dimension, but you can't enter it right now.");
        return;
    }

    if (map.playerCount(108010401, 0) > 0) {
        plr.sendMessage("Someone else is already fighting Dark Lord's clone. Please come back later.");
        return;
    }

    plr.warp(108010401);
}
