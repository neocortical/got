package index

import (
	"testing"
)

func TestCalculatePathnameNullsDoRead(t *testing.T) {
	var tests = []struct {
		input    uint16
		expected int
	}{
		{1, 1},
		{2, 8},
		{3, 7},
		{4, 6},
		{5, 5},
		{6, 4},
		{7, 3},
		{8, 2},
		{9, 1},
		{15, 3},
		{17, 1},
	}

	for i, test := range tests {
		actual := calculatePathnameNullsDoRead(test.input)
		if actual != test.expected {
			t.Errorf("test %d failed: expected %d but got %d", i, test.expected, actual)
		}
	}
}

func TestDirectoryReplacesFile(t *testing.T) {
	idx := &index{entryMap: map[string]*Entry{}, parentMap: map[string][]string{}}
	idx.Add(&Entry{name: "alice.txt", pathname: "alice.txt"})
	idx.Add(&Entry{name: "bob.txt", pathname: "bob.txt"})
	idx.Add(&Entry{name: "nested.txt", pathname: "alice.txt/nested.txt"})

	expectedNames := []string{"nested.txt", "bob.txt"}
	expectedPathnames := []string{"alice.txt/nested.txt", "bob.txt"}

	entries := idx.Entries()
	if len(entries) != len(expectedNames) {
		t.Fatalf("expected %d entries but got %d", len(expectedNames), len(entries))
	}

	for i, e := range idx.Entries() {
		if e.name != expectedNames[i] {
			t.Errorf("expected name %s but got %s", expectedNames[i], e.name)
		}
		if e.pathname != expectedPathnames[i] {
			t.Errorf("expected pathname %s but got %s", expectedPathnames[i], e.pathname)
		}
	}
}

func TestFileReplacesDirectory(t *testing.T) {
	idx := &index{entryMap: map[string]*Entry{}, parentMap: map[string][]string{}}
	idx.Add(&Entry{name: "alice.txt", pathname: "alice.txt"})
	idx.Add(&Entry{name: "bob.txt", pathname: "nested/bob.txt"})
	idx.Add(&Entry{name: "carol.txt", pathname: "nested/extranested/carol.txt"})
	idx.Add(&Entry{name: "nested", pathname: "nested"})

	expectedNames := []string{"alice.txt", "nested"}
	expectedPathnames := []string{"alice.txt", "nested"}

	entries := idx.Entries()
	if len(entries) != len(expectedNames) {
		t.Fatalf("expected %d entries but got %d", len(expectedNames), len(entries))
	}

	for i, e := range idx.Entries() {
		if e.name != expectedNames[i] {
			t.Errorf("expected name %s but got %s", expectedNames[i], e.name)
		}
		if e.pathname != expectedPathnames[i] {
			t.Errorf("expected pathname %s but got %s", expectedPathnames[i], e.pathname)
		}
	}
}
