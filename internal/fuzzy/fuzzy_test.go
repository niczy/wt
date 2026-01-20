package fuzzy

import (
	"testing"
)

func TestFuzzyMatch_ExactMatch(t *testing.T) {
	candidates := []string{"foo", "bar", "baz"}
	matches := FuzzyMatch("foo", candidates)

	if len(matches) == 0 {
		t.Fatal("expected at least one match")
	}
	if matches[0].Text != "foo" {
		t.Errorf("expected first match to be 'foo', got '%s'", matches[0].Text)
	}
}

func TestFuzzyMatch_ContainsPattern(t *testing.T) {
	candidates := []string{"my-feature", "your-feature", "bugfix"}
	matches := FuzzyMatch("feature", candidates)

	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
	// Both should contain "feature"
	for _, m := range matches {
		if m.Text != "my-feature" && m.Text != "your-feature" {
			t.Errorf("unexpected match: %s", m.Text)
		}
	}
}

func TestFuzzyMatch_PartialMatch(t *testing.T) {
	candidates := []string{"repo-feature-x", "repo-bugfix", "repo-feature-y"}
	matches := FuzzyMatch("fx", candidates)

	// Should match "repo-bugfix" (contains f and x in order)
	if len(matches) == 0 {
		t.Fatal("expected at least one match for 'fx'")
	}
}

func TestFuzzyMatch_NoMatch(t *testing.T) {
	candidates := []string{"foo", "bar", "baz"}
	matches := FuzzyMatch("xyz", candidates)

	if len(matches) != 0 {
		t.Errorf("expected no matches, got %d", len(matches))
	}
}

func TestFuzzyMatch_CaseInsensitive(t *testing.T) {
	candidates := []string{"MyFeature", "BUGFIX", "hotfix"}
	matches := FuzzyMatch("myfeature", candidates)

	if len(matches) == 0 {
		t.Fatal("expected case-insensitive match")
	}
	if matches[0].Text != "MyFeature" {
		t.Errorf("expected 'MyFeature', got '%s'", matches[0].Text)
	}
}

func TestFuzzyMatch_EmptyPattern(t *testing.T) {
	candidates := []string{"foo", "bar", "baz"}
	matches := FuzzyMatch("", candidates)

	// Empty pattern should match all
	if len(matches) != 3 {
		t.Errorf("expected 3 matches for empty pattern, got %d", len(matches))
	}
}

func TestFuzzyMatch_EmptyCandidates(t *testing.T) {
	matches := FuzzyMatch("foo", []string{})

	if len(matches) != 0 {
		t.Errorf("expected no matches for empty candidates, got %d", len(matches))
	}
}

func TestFuzzyMatch_SortedByScore(t *testing.T) {
	candidates := []string{"abcdef", "abc", "abcdefghij"}
	matches := FuzzyMatch("abc", candidates)

	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}

	// Exact match should be first
	if matches[0].Text != "abc" {
		t.Errorf("expected exact match 'abc' first, got '%s'", matches[0].Text)
	}
}

func TestBestMatch(t *testing.T) {
	candidates := []string{"feature-a", "feature-b", "bugfix"}
	best := BestMatch("feature", candidates)

	if best != "feature-a" && best != "feature-b" {
		t.Errorf("expected one of the feature matches, got '%s'", best)
	}
}

func TestBestMatch_NoMatch(t *testing.T) {
	candidates := []string{"foo", "bar"}
	best := BestMatch("xyz", candidates)

	if best != "" {
		t.Errorf("expected empty string for no match, got '%s'", best)
	}
}

func TestTopMatches(t *testing.T) {
	candidates := []string{"a1", "a2", "a3", "a4", "a5"}
	top := TopMatches("a", candidates, 3)

	if len(top) != 3 {
		t.Errorf("expected 3 top matches, got %d", len(top))
	}
}

func TestTopMatches_FewerThanN(t *testing.T) {
	candidates := []string{"a1", "a2"}
	top := TopMatches("a", candidates, 5)

	if len(top) != 2 {
		t.Errorf("expected 2 matches (all available), got %d", len(top))
	}
}

func TestCalculateScore_WordBoundaryBonus(t *testing.T) {
	// Test that word boundaries get bonus
	score1 := calculateScore("fb", "foo-bar")   // f at start, b after hyphen
	score2 := calculateScore("fb", "foobxxbar") // f and b not at boundaries

	if score1 <= score2 {
		t.Errorf("expected word boundary match to score higher: boundary=%d, non-boundary=%d", score1, score2)
	}
}

func TestCalculateScore_ConsecutiveBonus(t *testing.T) {
	// Consecutive matches should score higher
	score1 := calculateScore("abc", "abcdef")  // consecutive
	score2 := calculateScore("abc", "axbxcxx") // non-consecutive

	if score1 <= score2 {
		t.Errorf("expected consecutive match to score higher: consecutive=%d, non-consecutive=%d", score1, score2)
	}
}
