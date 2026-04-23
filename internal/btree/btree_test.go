package btree

import (
	"slices"
	"testing"
)

// diagramKeys are all the keys visible in the simplified PostgreSQL
// B-Tree index diagram, inserted in a representative order.
var diagramKeys = []int{
	10, 30, 5, 8, 1, 3, 4, 6, 7, 9,
	12, 18, 11, 13, 15, 19, 25, 28,
	31, 42, 32, 35, 39, 43, 45, 48,
}

func buildDiagramTree(t *testing.T) *BTree {
	t.Helper()
	tree := New(2)
	for _, k := range diagramKeys {
		tree.Insert(k)
	}
	return tree
}

func TestInsertAndSearch(t *testing.T) {
	tree := buildDiagramTree(t)

	for _, k := range diagramKeys {
		if !tree.Search(k) {
			t.Errorf("expected key %d to be found", k)
		}
	}

	missing := []int{0, 2, 14, 20, 50, 100}
	for _, k := range missing {
		if tree.Search(k) {
			t.Errorf("did not expect key %d to be found", k)
		}
	}
}

func TestInOrder(t *testing.T) {
	tree := buildDiagramTree(t)
	got := tree.InOrder()

	want := make([]int, len(diagramKeys))
	copy(want, diagramKeys)
	slices.Sort(want)

	if len(got) != len(want) {
		t.Fatalf("in-order length = %d; want %d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("in-order[%d] = %d; want %d", i, got[i], want[i])
		}
	}
}

func TestHeight(t *testing.T) {
	tree := buildDiagramTree(t)
	h := tree.Height()
	if h < 2 {
		t.Errorf("height = %d; expected at least 2 for %d keys", h, len(diagramKeys))
	}
}

func TestDelete(t *testing.T) {
	tree := buildDiagramTree(t)

	toDelete := []int{6, 13, 30, 1, 48}
	for _, k := range toDelete {
		if !tree.Delete(k) {
			t.Errorf("Delete(%d) returned false; expected true", k)
		}
		if tree.Search(k) {
			t.Errorf("key %d still found after deletion", k)
		}
	}

	remaining := make([]int, len(diagramKeys))
	copy(remaining, diagramKeys)
	for _, d := range toDelete {
		for i, v := range remaining {
			if v == d {
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}

	for _, k := range remaining {
		if !tree.Search(k) {
			t.Errorf("expected remaining key %d to still be found", k)
		}
	}
}

func TestDeleteNonExistent(t *testing.T) {
	tree := buildDiagramTree(t)
	if tree.Delete(999) {
		t.Error("Delete(999) returned true; expected false for non-existent key")
	}
}

func TestDeleteAllKeys(t *testing.T) {
	tree := buildDiagramTree(t)

	sorted := make([]int, len(diagramKeys))
	copy(sorted, diagramKeys)
	slices.Sort(sorted)

	for _, k := range sorted {
		tree.Delete(k)
	}

	if len(tree.InOrder()) != 0 {
		t.Error("expected empty tree after deleting all keys")
	}
}

func TestEmptyTree(t *testing.T) {
	tree := New(2)

	if tree.Search(42) {
		t.Error("Search on empty tree should return false")
	}
	if tree.Delete(42) {
		t.Error("Delete on empty tree should return false")
	}
	if len(tree.InOrder()) != 0 {
		t.Error("InOrder on empty tree should return empty slice")
	}
	if tree.Height() != 1 {
		t.Errorf("Height of empty tree = %d; want 1", tree.Height())
	}
}

func TestString(t *testing.T) {
	tree := buildDiagramTree(t)
	s := tree.String()
	if s == "" || s == "(empty)" {
		t.Error("String() should produce a non-empty representation")
	}
}

func TestMinDegreeFloor(t *testing.T) {
	tree := New(0)
	tree.Insert(5)
	tree.Insert(3)
	tree.Insert(7)
	if !tree.Search(5) || !tree.Search(3) || !tree.Search(7) {
		t.Error("tree with clamped min degree should still function")
	}
}

func TestLargerMinDegree(t *testing.T) {
	tree := New(4)
	for i := range 100 {
		tree.Insert(i)
	}
	for i := range 100 {
		if !tree.Search(i) {
			t.Errorf("expected key %d to be found in degree-4 tree", i)
		}
	}
	got := tree.InOrder()
	for i := range 100 {
		if got[i] != i {
			t.Fatalf("in-order[%d] = %d; want %d", i, got[i], i)
		}
	}
}
