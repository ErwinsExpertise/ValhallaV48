package nx

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadV48DirectorySmoke(t *testing.T) {
	v48Path := filepath.Join("..", "..", "v48", "wz", "nx")
	if _, err := os.Stat(v48Path); err != nil {
		t.Skipf("v48 NX directory unavailable: %v", err)
	}

	LoadFile(v48Path)

	item, err := GetItem(2000000)
	if err != nil {
		t.Fatalf("GetItem(2000000): %v", err)
	}
	if item.Name == "" {
		t.Fatalf("GetItem(2000000): empty name")
	}

	cider, err := GetItem(2022002)
	if err != nil {
		t.Fatalf("GetItem(2022002): %v", err)
	}
	if cider.Time != 180000 {
		t.Fatalf("GetItem(2022002): expected 180000ms duration, got %d", cider.Time)
	}
	if cider.PAD != 20 || cider.ACC != -5 {
		t.Fatalf("GetItem(2022002): unexpected stat change payload: PAD=%d ACC=%d", cider.PAD, cider.ACC)
	}

	morph, err := GetItem(2210000)
	if err != nil {
		t.Fatalf("GetItem(2210000): %v", err)
	}
	if morph.Morph != 1 || morph.Time != 3600000 {
		t.Fatalf("GetItem(2210000): unexpected morph data: morph=%d time=%d", morph.Morph, morph.Time)
	}

	field, err := GetMap(100000000)
	if err != nil {
		t.Fatalf("GetMap(100000000): %v", err)
	}
	if field.MapName == "" || field.StreetName == "" {
		t.Fatalf("GetMap(100000000): missing map names: map=%q street=%q", field.MapName, field.StreetName)
	}

	mob, err := GetMob(100100)
	if err != nil {
		t.Fatalf("GetMob(100100): %v", err)
	}
	if mob.MaxHP <= 0 {
		t.Fatalf("GetMob(100100): invalid hp %d", mob.MaxHP)
	}

	quest, err := GetQuest(1000)
	if err != nil {
		t.Fatalf("GetQuest(1000): %v", err)
	}
	if quest.Name == "" {
		t.Fatalf("GetQuest(1000): empty name")
	}

	if len(GetMaps()) < 900 {
		t.Fatalf("unexpected map count: %d", len(GetMaps()))
	}
	if len(GetCommodities()) < 4000 {
		t.Fatalf("unexpected commodity count: %d", len(GetCommodities()))
	}
	if len(GetPackages()) < 200 {
		t.Fatalf("unexpected package count: %d", len(GetPackages()))
	}
	if len(GetReactorInfoList()) < 200 {
		t.Fatalf("unexpected reactor count: %d", len(GetReactorInfoList()))
	}
	if len(GetQuests()) < 500 {
		t.Fatalf("unexpected quest count: %d", len(GetQuests()))
	}
}
