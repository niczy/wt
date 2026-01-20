package fuzzy

import (
	"sort"
	"strings"
)

// Match represents a fuzzy match result
type Match struct {
	Text  string
	Score int
}

// FuzzyMatch performs fuzzy matching on a list of strings
// Returns matches sorted by score (higher is better)
func FuzzyMatch(pattern string, candidates []string) []Match {
	pattern = strings.ToLower(pattern)
	var matches []Match

	for _, candidate := range candidates {
		score := calculateScore(pattern, strings.ToLower(candidate))
		if score > 0 {
			matches = append(matches, Match{
				Text:  candidate,
				Score: score,
			})
		}
	}

	// Sort by score descending
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches
}

// calculateScore computes a fuzzy match score
// Returns 0 if pattern doesn't match
func calculateScore(pattern, text string) int {
	if pattern == "" {
		return 1
	}

	// Exact match gets highest score
	if text == pattern {
		return 1000
	}

	// Contains exact pattern gets high score
	if strings.Contains(text, pattern) {
		return 500 + (100 - len(text)) // Prefer shorter matches
	}

	// Fuzzy match: all characters must appear in order
	patternIdx := 0
	score := 0
	lastMatchIdx := -1
	consecutiveBonus := 0

	for i := 0; i < len(text) && patternIdx < len(pattern); i++ {
		if text[i] == pattern[patternIdx] {
			// Bonus for consecutive matches
			if lastMatchIdx == i-1 {
				consecutiveBonus += 10
			}
			// Bonus for matching at word boundaries
			if i == 0 || text[i-1] == '-' || text[i-1] == '_' || text[i-1] == '/' {
				score += 20
			}
			score += 10 + consecutiveBonus
			lastMatchIdx = i
			patternIdx++
		}
	}

	// Only return score if all pattern characters were matched
	if patternIdx == len(pattern) {
		return score
	}

	return 0
}

// BestMatch returns the single best match, or empty string if none
func BestMatch(pattern string, candidates []string) string {
	matches := FuzzyMatch(pattern, candidates)
	if len(matches) == 0 {
		return ""
	}
	return matches[0].Text
}

// TopMatches returns up to n best matches
func TopMatches(pattern string, candidates []string, n int) []string {
	matches := FuzzyMatch(pattern, candidates)
	result := make([]string, 0, n)
	for i := 0; i < len(matches) && i < n; i++ {
		result = append(result, matches[i].Text)
	}
	return result
}
