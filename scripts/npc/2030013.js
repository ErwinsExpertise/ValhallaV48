function renderExpeditionList(state) {
    var lines = []
    lines.push("#eZakum Expedition Signup List#n")
    lines.push("Leader: #b" + state.leaderName + "#k")
    lines.push("Members: #b" + state.members.length + " / 20#k")

    if (state.open) {
        var minutes = Math.floor(state.remainingSeconds / 60)
        var seconds = state.remainingSeconds % 60
        var secondsText = seconds < 10 ? "0" + seconds : "" + seconds
        lines.push("Time remaining: #b" + minutes + ":" + secondsText + "#k")
    } else if (state.started) {
        lines.push("The expedition has already departed.")
    } else {
        lines.push("Registrations are closed.")
    }

    lines.push("")
    for (var i = 0; i < state.members.length; i++) {
        var member = state.members[i]
        var prefix = member.isLeader ? "#r[Leader]#k " : ""
        lines.push(prefix + member.name)
    }

    return lines.join("\r\n")
}

function sendMenuSelection(text, labels) {
    var menu = text
    for (var i = 0; i < labels.length; i++) {
        menu += "\r\n#L" + i + "#" + labels[i] + "#l"
    }
    npc.sendSelection(menu)
    return npc.selection()
}

function handleResult(result) {
    if (result.ok) {
        if (result.message && result.message.length > 0) {
            npc.sendOk(result.message)
        }
        return true
    }

    if (!result.handled && result.message && result.message.length > 0) {
        npc.sendOk(result.message)
    }

    return false
}

var state = plr.getZakumExpedition()
var myName = plr.name()

if (!state.exists) {
    if (!npc.sendYesNo("No Zakum expedition has been registered yet. Would you like to become the expedition leader?")) {
        npc.sendOk("Come back when you are ready to assemble an expedition.")
    } else {
        handleResult(plr.createZakumExpedition())
    }
} else {
    var isLeader = state.leaderName === myName
    var isMember = false
    for (var i = 0; i < state.members.length; i++) {
        if (state.members[i].name === myName) {
            isMember = true
            break
        }
    }

    if (isLeader) {
        var options = ["View signup list", "Start expedition now", "Disband expedition"]
        var choice = sendMenuSelection("What would you like to do?\r\n", options)

        if (choice === 0) {
            var listState = plr.getZakumExpedition()
            if (!listState.exists) {
                npc.sendOk("There is no active Zakum expedition right now.")
            } else {
                var kickOptions = []
                var kickMembers = []
                for (var j = 0; j < listState.members.length; j++) {
                    if (!listState.members[j].isLeader) {
                        kickMembers.push(listState.members[j])
                        kickOptions.push("Kick " + listState.members[j].name)
                    }
                }

                if (kickOptions.length === 0) {
                    npc.sendOk(renderExpeditionList(listState))
                } else {
                    var kickMenu = renderExpeditionList(listState) + "\r\n\r\nSelect a member below to remove them from the expedition.\r\n#L0#Back#l"
                    for (var k = 0; k < kickOptions.length; k++) {
                        kickMenu += "\r\n#L" + (k + 1) + "#" + kickOptions[k] + "#l"
                    }

                    npc.sendSelection(kickMenu)
                    var kickChoice = npc.selection() - 1
                    if (kickChoice >= 0 && kickChoice < kickMembers.length) {
                        handleResult(plr.kickZakumExpeditionMember(kickMembers[kickChoice].id))
                    }
                }
            }
        } else if (choice === 1) {
            var startResult = plr.startZakumExpedition()
            if (!startResult.ok && !startResult.handled && startResult.message && startResult.message.length > 0) {
                npc.sendOk(startResult.message)
            }
        } else if (choice === 2) {
            if (npc.sendYesNo("Would you like to disband the current Zakum expedition?")) {
                handleResult(plr.terminateZakumExpedition())
            } else {
                npc.sendOk("The expedition will remain registered.")
            }
        }
    } else if (state.open) {
        if (isMember) {
            var memberChoice = sendMenuSelection("What would you like to do?\r\n", ["View signup list", "Leave expedition"])
            if (memberChoice === 0) {
                npc.sendOk(renderExpeditionList(plr.getZakumExpedition()))
            } else if (memberChoice === 1) {
                handleResult(plr.leaveZakumExpedition())
            }
        } else {
            var guestChoice = sendMenuSelection("What would you like to do?\r\n", ["Join expedition", "View signup list"])
            if (guestChoice === 0) {
                handleResult(plr.joinZakumExpedition())
            } else if (guestChoice === 1) {
                npc.sendOk(renderExpeditionList(plr.getZakumExpedition()))
            }
        }
    } else {
        npc.sendOk(state.started
            ? "The current Zakum expedition has already departed. Please wait for the next expedition."
            : "Registrations for the current Zakum expedition are closed. Please wait for the next expedition.")
    }
}
