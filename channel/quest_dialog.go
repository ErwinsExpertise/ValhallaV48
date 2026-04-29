package channel

import (
	"fmt"
	"strings"
)

type questDialogKind byte

const (
	questDialogStart questDialogKind = iota
	questDialogComplete
	questDialogProgress
)

type questDialogState byte

const (
	questDialogStateMenu questDialogState = iota
	questDialogStatePrompt
	questDialogStateRewardSelection
	questDialogStateOk
)

type questDialogEntry struct {
	kind    questDialogKind
	questID int16
}

type questDialog struct {
	npcID          int32
	entries        []questDialogEntry
	active         *questDialogEntry
	state          questDialogState
	rewardChoices  []string
	rewardSelected bool
}

func openQuestDialog(plr *Player, npcID int32) *questDialog {
	available, inProgress, completable := plr.questIDsForNPC(npcID)
	entries := make([]questDialogEntry, 0, len(available)+len(inProgress)+len(completable))

	for _, questID := range completable {
		entries = append(entries, questDialogEntry{kind: questDialogComplete, questID: questID})
	}
	for _, questID := range available {
		entries = append(entries, questDialogEntry{kind: questDialogStart, questID: questID})
	}
	for _, questID := range inProgress {
		entries = append(entries, questDialogEntry{kind: questDialogProgress, questID: questID})
	}

	if len(entries) == 0 {
		return nil
	}

	dialog := &questDialog{npcID: npcID, entries: entries}
	if len(entries) == 1 {
		dialog.active = &dialog.entries[0]
	}
	return dialog
}

func (d *questDialog) start(plr *Player) {
	if d.active == nil {
		d.state = questDialogStateMenu
		plr.Send(packetNpcChatSelection(d.npcID, d.buildEntryMenu(plr)))
		return
	}
	d.sendActivePrompt(plr)
}

func (d *questDialog) handle(plr *Player, action byte, selection int) bool {
	if action == 0 {
		return true
	}

	switch d.state {
	case questDialogStateMenu:
		if selection < 0 || selection >= len(d.entries) {
			return true
		}
		d.active = &d.entries[selection]
		d.sendActivePrompt(plr)
		return false
	case questDialogStatePrompt:
		if d.active == nil {
			return true
		}
		switch d.active.kind {
		case questDialogStart:
			_ = plr.tryStartQuest(d.active.questID)
			return true
		case questDialogComplete:
			d.rewardChoices = plr.questSelectableRewards(d.active.questID)
			if len(d.rewardChoices) > 0 {
				d.state = questDialogStateRewardSelection
				plr.Send(packetNpcChatSelection(d.npcID, d.buildRewardMenu()))
				return false
			}
			_ = plr.tryCompleteQuest(d.active.questID)
			return true
		default:
			return true
		}
	case questDialogStateRewardSelection:
		if d.active == nil || selection < 0 || selection >= len(d.rewardChoices) {
			return true
		}
		_ = plr.tryCompleteQuestSelection(d.active.questID, selection)
		return true
	case questDialogStateOk:
		return true
	default:
		return true
	}
}

func (d *questDialog) sendActivePrompt(plr *Player) {
	if d.active == nil {
		d.state = questDialogStateOk
		plr.Send(packetNpcChatOk(d.npcID, "No quests are available."))
		return
	}

	questID := d.active.questID
	switch d.active.kind {
	case questDialogStart:
		d.state = questDialogStatePrompt
		plr.Send(packetNpcChatYesNo(d.npcID, joinQuestLines(plr.questSayLines(questID, "start.0"), plr.questSayLines(questID, "start.yes"), fmt.Sprintf("Would you like to start %s?", plr.questDisplayName(questID)))))
	case questDialogComplete:
		d.state = questDialogStatePrompt
		plr.Send(packetNpcChatYesNo(d.npcID, joinQuestLines(plr.questSayLines(questID, "complete.0"), plr.questSayLines(questID, "complete.yes"), fmt.Sprintf("Would you like to complete %s?", plr.questDisplayName(questID)))))
	default:
		d.state = questDialogStateOk
		plr.Send(packetNpcChatOk(d.npcID, strings.Join(plr.questIncompleteLines(questID), "\n\n")))
	}
}

func (d *questDialog) buildEntryMenu(plr *Player) string {
	var b strings.Builder
	b.WriteString("What would you like to do?\n")
	for i, entry := range d.entries {
		prefix := "In Progress"
		switch entry.kind {
		case questDialogComplete:
			prefix = "Complete"
		case questDialogStart:
			prefix = "Start"
		}
		fmt.Fprintf(&b, "#L%d#%s: %s#l\n", i, prefix, plr.questDisplayName(entry.questID))
	}
	return b.String()
}

func (d *questDialog) buildRewardMenu() string {
	var b strings.Builder
	b.WriteString("Choose your reward.\n")
	for i, reward := range d.rewardChoices {
		fmt.Fprintf(&b, "#L%d#%s#l\n", i, reward)
	}
	return b.String()
}

func joinQuestLines(primary, secondary []string, fallback string) string {
	parts := make([]string, 0, len(primary)+len(secondary)+1)
	parts = append(parts, primary...)
	parts = append(parts, secondary...)
	parts = compactQuestLines(parts)
	if len(parts) == 0 {
		return fallback
	}
	return strings.Join(parts, "\n\n")
}

func compactQuestLines(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}
