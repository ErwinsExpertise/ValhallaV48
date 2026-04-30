package channel

import (
	"testing"

	"github.com/Hucaru/Valhalla/constant"
)

func TestEngagementItemsForProposalBox(t *testing.T) {
	tests := []struct {
		box          int32
		partnerRing  int32
		emptyBoxItem int32
	}{
		{constant.ItemEngagementBoxMoonstone, constant.ItemEngagementRingMoonstone, constant.ItemEmptyEngagementBoxMoonstone},
		{constant.ItemEngagementBoxStar, constant.ItemEngagementRingStar, constant.ItemEmptyEngagementBoxStar},
		{constant.ItemEngagementBoxGolden, constant.ItemEngagementRingGolden, constant.ItemEmptyEngagementBoxGolden},
		{constant.ItemEngagementBoxSilver, constant.ItemEngagementRingSilver, constant.ItemEmptyEngagementBoxSilver},
	}

	for _, tc := range tests {
		ring, empty := engagementItemsForProposalBox(tc.box)
		if ring != tc.partnerRing || empty != tc.emptyBoxItem {
			t.Fatalf("unexpected engagement pair for %d: got ring=%d empty=%d", tc.box, ring, empty)
		}
	}
}

func TestSameWeddingPartyRequiresSharedTwoPersonParty(t *testing.T) {
	partyA := &party{}
	plr := &Player{ID: 1}
	partner := &Player{ID: 2}
	stranger := &Player{ID: 3}

	partyA.players[0] = plr
	partyA.players[1] = partner
	plr.party = partyA
	partner.party = partyA

	if !sameWeddingParty(plr, partner) {
		t.Fatal("expected shared two-person party to be valid")
	}

	partyA.players[2] = stranger
	stranger.party = partyA
	if sameWeddingParty(plr, partner) {
		t.Fatal("expected party with more than two members to be invalid")
	}
}
