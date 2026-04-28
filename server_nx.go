package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hucaru/Valhalla/nx"
)

func loadNXData(configFile string) {
	nxConfig := nxConfigFromFile(configFile)
	nxPath := resolveNXPath(nxConfig.Path)

	start := time.Now()
	nx.LoadFile(nxPath)
	log.Printf("Loaded and parsed Wizet data (NX) from %q in %s", nxPath, time.Since(start))
}

func resolveNXPath(configPath string) string {
	if path := strings.TrimSpace(*nxPtr); path != "" {
		return path
	}

	if path := strings.TrimSpace(configPath); path != "" {
		return path
	}

	candidates := []string{
		"Data.nx",
		"nx",
	}

	if executablePath, err := os.Executable(); err == nil {
		executableDir := filepath.Dir(executablePath)
		candidates = append(candidates,
			filepath.Join(executableDir, "Data.nx"),
			filepath.Join(executableDir, "nx"),
		)
	}

	// Keep the old Hucaru-derived local layout as a final fallback so older checkouts still work.
	candidates = append(candidates, filepath.Join("..", "v48", "wz", "nx"))

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	return candidates[0]
}
