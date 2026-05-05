package nx

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type StringMatch struct {
	Type  string
	ID    int32
	Name  string
	Extra string
}

var itemSearchAliases = map[string]string{
	"item":        "item",
	"items":       "item",
	"equip":       "equip",
	"equips":      "equip",
	"cash":        "cash",
	"use":         "use",
	"uses":        "use",
	"consume":     "use",
	"consumes":    "use",
	"setup":       "setup",
	"setups":      "setup",
	"install":     "setup",
	"installs":    "setup",
	"etc":         "etc",
	"pet":         "pet",
	"pets":        "pet",
	"accessory":   "accessory",
	"accessories": "accessory",
	"cap":         "cap",
	"caps":        "cap",
	"cape":        "cape",
	"capes":       "cape",
	"coat":        "coat",
	"coats":       "coat",
	"face":        "face",
	"faces":       "face",
	"glove":       "glove",
	"gloves":      "glove",
	"hair":        "hair",
	"hairs":       "hair",
	"longcoat":    "longcoat",
	"longcoats":   "longcoat",
	"overall":     "longcoat",
	"overalls":    "longcoat",
	"robe":        "longcoat",
	"robes":       "longcoat",
	"pants":       "pants",
	"pant":        "pants",
	"petequip":    "petequip",
	"pet-equip":   "petequip",
	"pet_equip":   "petequip",
	"ring":        "ring",
	"rings":       "ring",
	"shield":      "shield",
	"shields":     "shield",
	"shoe":        "shoes",
	"shoes":       "shoes",
	"weapon":      "weapon",
	"weapons":     "weapon",
}

var equipSearchCategories = map[string]bool{
	"accessory": true,
	"cap":       true,
	"cape":      true,
	"coat":      true,
	"face":      true,
	"glove":     true,
	"hair":      true,
	"longcoat":  true,
	"pants":     true,
	"petequip":  true,
	"ring":      true,
	"shield":    true,
	"shoes":     true,
	"weapon":    true,
}

func SearchItemsByName(query string, limit int) []StringMatch {
	return SearchItemsByCategory(query, "item", limit)
}

func SearchItemsByCategory(query, category string, limit int) []StringMatch {
	query = normalizeSearchQuery(query)
	category, ok := NormalizeItemSearchCategory(category)
	if query == "" || limit <= 0 || !ok {
		return nil
	}

	queryTokens := tokenizeSearchQuery(query)
	results := make([]StringMatch, 0, limit)
	for id, item := range items {
		if item.Name == "" || !matchesSearchQuery(item.Name, query, queryTokens) {
			continue
		}
		if !matchesItemSearchCategory(item, category) {
			continue
		}
		results = append(results, StringMatch{Type: "item", ID: id, Name: item.Name, Extra: item.SearchCategory})
	}

	sortStringMatches(results)
	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func NormalizeItemSearchCategory(category string) (string, bool) {
	normalized := normalizeSearchQuery(category)
	if normalized == "" {
		return "item", true
	}
	resolved, ok := itemSearchAliases[normalized]
	return resolved, ok
}

func SearchMapsByName(query string, limit int) []StringMatch {
	query = normalizeSearchQuery(query)
	if query == "" || limit <= 0 {
		return nil
	}

	queryTokens := tokenizeSearchQuery(query)
	results := make([]StringMatch, 0, limit)
	for id, field := range maps {
		name := strings.TrimSpace(field.MapName)
		street := strings.TrimSpace(field.StreetName)
		if name == "" && street == "" {
			continue
		}

		haystack := strings.TrimSpace(street + " " + name)
		if !matchesSearchQuery(haystack, query, queryTokens) {
			continue
		}

		fullName := name
		if street != "" {
			fullName = fmt.Sprintf("%s - %s", street, name)
		}

		results = append(results, StringMatch{Type: "map", ID: id, Name: fullName})
	}

	sortStringMatches(results)
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

	queryTokens := tokenizeSearchQuery(query)
	results := make([]StringMatch, 0, limit)
	for id, quest := range quests {
		if quest.Name == "" || !matchesSearchQuery(quest.Name, query, queryTokens) {
			continue
		}
		results = append(results, StringMatch{Type: "quest", ID: int32(id), Name: quest.Name})
	}

	sortStringMatches(results)
	if len(results) > limit {
		results = results[:limit]
	}

	return results
}

func matchesItemSearchCategory(item Item, category string) bool {
	if category == "item" {
		return true
	}
	if category == "equip" {
		return equipSearchCategories[item.SearchCategory]
	}
	return item.SearchCategory == category
}

func sortStringMatches(results []StringMatch) {
	sort.Slice(results, func(i, j int) bool {
		if results[i].Name == results[j].Name {
			return results[i].ID < results[j].ID
		}
		return results[i].Name < results[j].Name
	})
}

func normalizeSearchQuery(query string) string {
	return strings.ToLower(strings.TrimSpace(query))
}

func tokenizeSearchQuery(query string) []string {
	fields := strings.FieldsFunc(normalizeSearchQuery(query), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	tokens := make([]string, 0, len(fields))
	for _, field := range fields {
		if field != "" {
			tokens = append(tokens, field)
		}
	}
	return tokens
}

func matchesSearchQuery(name, query string, queryTokens []string) bool {
	normalizedName := normalizeSearchQuery(name)
	if strings.Contains(normalizedName, query) {
		return true
	}
	for _, token := range queryTokens {
		if !strings.Contains(normalizedName, token) {
			return false
		}
	}
	return len(queryTokens) > 0
}
