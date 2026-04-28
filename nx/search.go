package nx

import (
	"fmt"
	"sort"
	"strings"
)

type StringMatch struct {
	Type  string
	ID    int32
	Name  string
	Extra string
}

func SearchItemsByName(query string, limit int) []StringMatch {
	query = normalizeSearchQuery(query)
	if query == "" || limit <= 0 {
		return nil
	}

	results := make([]StringMatch, 0, limit)
	for id, item := range items {
		if item.Name == "" || !strings.Contains(strings.ToLower(item.Name), query) {
			continue
		}
		results = append(results, StringMatch{Type: "item", ID: id, Name: item.Name})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Name == results[j].Name {
			return results[i].ID < results[j].ID
		}
		return results[i].Name < results[j].Name
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func SearchMapsByName(query string, limit int) []StringMatch {
	query = normalizeSearchQuery(query)
	if query == "" || limit <= 0 {
		return nil
	}

	results := make([]StringMatch, 0, limit)
	for id, field := range maps {
		name := strings.TrimSpace(field.MapName)
		street := strings.TrimSpace(field.StreetName)
		if name == "" && street == "" {
			continue
		}

		haystack := strings.ToLower(strings.TrimSpace(street + " " + name))
		if !strings.Contains(haystack, query) {
			continue
		}

		fullName := name
		if street != "" {
			fullName = fmt.Sprintf("%s - %s", street, name)
		}

		results = append(results, StringMatch{Type: "map", ID: id, Name: fullName})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Name == results[j].Name {
			return results[i].ID < results[j].ID
		}
		return results[i].Name < results[j].Name
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func SearchQuestsByName(query string, limit int) []StringMatch {
	query = normalizeSearchQuery(query)
	if query == "" || limit <= 0 {
		return nil
	}

	results := make([]StringMatch, 0, limit)
	for id, quest := range quests {
		if quest.Name == "" || !strings.Contains(strings.ToLower(quest.Name), query) {
			continue
		}
		results = append(results, StringMatch{Type: "quest", ID: int32(id), Name: quest.Name})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Name == results[j].Name {
			return results[i].ID < results[j].ID
		}
		return results[i].Name < results[j].Name
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func normalizeSearchQuery(query string) string {
	return strings.ToLower(strings.TrimSpace(query))
}
