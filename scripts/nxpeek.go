package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Hucaru/gonx"
)

func main() {
	file := flag.String("file", "", "nx file to inspect")
	path := flag.String("path", "", "node path to inspect")
	depth := flag.Int("depth", 1, "depth from selected node")
	flag.Parse()

	if *file == "" {
		log.Fatal("-file is required")
	}

	nodes, text, _, _, err := gonx.Parse(*file)
	if err != nil {
		log.Fatal(err)
	}

	if *path == "" {
		printNode(&nodes[0], nodes, text, 0, *depth)
		return
	}

	var target *gonx.Node
	if !gonx.FindNode(*path, nodes, text, func(n *gonx.Node) { target = n }) {
		log.Fatalf("path not found: %s", *path)
	}

	printNode(target, nodes, text, 0, *depth)
}

func printNode(node *gonx.Node, nodes []gonx.Node, text []string, level, maxDepth int) {
	name := text[node.NameID]
	if name == "" && level == 0 {
		name = filepath.Base("/")
	}
	fmt.Printf("%s%s children=%d type=%d data=%v\n", strings.Repeat("  ", level), name, node.ChildCount, node.Type, node.Data)

	if level >= maxDepth {
		return
	}

	children := make([]gonx.Node, 0, node.ChildCount)
	for i := uint32(0); i < uint32(node.ChildCount); i++ {
		children = append(children, nodes[node.ChildID+i])
	}

	sort.Slice(children, func(i, j int) bool {
		return text[children[i].NameID] < text[children[j].NameID]
	})

	for i := range children {
		child := children[i]
		printNode(&child, nodes, text, level+1, maxDepth)
	}
}
