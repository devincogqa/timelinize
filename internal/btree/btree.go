package btree

import (
	"fmt"
	"strconv"
	"strings"
)

// BTree represents a B-Tree data structure, modeled after the simplified
// PostgreSQL B-Tree index. It supports generic integer keys with a
// configurable minimum degree.
type BTree struct {
	root      *node
	minDegree int // minimum degree (t): every node except root must have at least t-1 keys
}

type node struct {
	keys     []int
	children []*node
	leaf     bool
}

func maxKeysForDegree(t int) int {
	return 2*t - 1
}

// New creates a new B-Tree with the given minimum degree t.
// The minimum degree must be at least 2.
func New(minDegree int) *BTree {
	const minAllowedDegree = 2
	if minDegree < minAllowedDegree {
		minDegree = minAllowedDegree
	}
	return &BTree{
		root: &node{
			leaf: true,
		},
		minDegree: minDegree,
	}
}

// Search returns true if the key exists in the B-Tree.
func (t *BTree) Search(key int) bool {
	return t.root.search(key)
}

func (n *node) search(key int) bool {
	i := 0
	for i < len(n.keys) && key > n.keys[i] {
		i++
	}
	if i < len(n.keys) && n.keys[i] == key {
		return true
	}
	if n.leaf {
		return false
	}
	return n.children[i].search(key)
}

// Insert adds a key to the B-Tree.
func (t *BTree) Insert(key int) {
	root := t.root
	if len(root.keys) == maxKeysForDegree(t.minDegree) {
		newRoot := &node{
			children: []*node{root},
		}
		newRoot.splitChild(0, t.minDegree)
		t.root = newRoot

		i := 0
		if newRoot.keys[0] < key {
			i++
		}
		newRoot.children[i].insertNonFull(key, t.minDegree)
	} else {
		root.insertNonFull(key, t.minDegree)
	}
}

func (n *node) insertNonFull(key int, minDegree int) {
	i := len(n.keys) - 1
	maxKeys := maxKeysForDegree(minDegree)

	if n.leaf {
		n.keys = append(n.keys, 0)
		for i >= 0 && n.keys[i] > key {
			n.keys[i+1] = n.keys[i]
			i--
		}
		n.keys[i+1] = key
		return
	}

	for i >= 0 && n.keys[i] > key {
		i--
	}
	i++

	if len(n.children[i].keys) == maxKeys {
		n.splitChild(i, minDegree)
		if n.keys[i] < key {
			i++
		}
	}
	n.children[i].insertNonFull(key, minDegree)
}

func (n *node) splitChild(i int, minDegree int) {
	child := n.children[i]
	sibling := &node{
		leaf: child.leaf,
	}

	midIndex := minDegree - 1
	midKey := child.keys[midIndex]

	sibling.keys = make([]int, len(child.keys[midIndex+1:]))
	copy(sibling.keys, child.keys[midIndex+1:])

	if !child.leaf {
		sibling.children = make([]*node, len(child.children[midIndex+1:]))
		copy(sibling.children, child.children[midIndex+1:])
	}

	child.keys = child.keys[:midIndex]
	if !child.leaf {
		child.children = child.children[:midIndex+1]
	}

	n.children = append(n.children, nil)
	copy(n.children[i+2:], n.children[i+1:])
	n.children[i+1] = sibling

	n.keys = append(n.keys, 0)
	copy(n.keys[i+1:], n.keys[i:])
	n.keys[i] = midKey
}

// Delete removes a key from the B-Tree. Returns true if the key was found
// and removed.
func (t *BTree) Delete(key int) bool {
	if t.root == nil {
		return false
	}
	found := t.root.delete(key, t.minDegree)
	if len(t.root.keys) == 0 && !t.root.leaf {
		t.root = t.root.children[0]
	}
	return found
}

func (n *node) delete(key int, minDeg int) bool {
	idx := n.findKeyIndex(key)

	if idx < len(n.keys) && n.keys[idx] == key {
		if n.leaf {
			n.keys = append(n.keys[:idx], n.keys[idx+1:]...)
			return true
		}
		return n.deleteInternal(idx, minDeg)
	}

	if n.leaf {
		return false
	}

	child := n.children[idx]
	if len(child.keys) < minDeg {
		n.fill(idx, minDeg)
	}

	if idx > len(n.keys) {
		return n.children[idx-1].delete(key, minDeg)
	}
	return n.children[idx].delete(key, minDeg)
}

func (n *node) findKeyIndex(key int) int {
	idx := 0
	for idx < len(n.keys) && n.keys[idx] < key {
		idx++
	}
	return idx
}

func (n *node) deleteInternal(idx int, minDeg int) bool {
	key := n.keys[idx]

	if len(n.children[idx].keys) >= minDeg {
		pred := n.getPredecessor(idx)
		n.keys[idx] = pred
		return n.children[idx].delete(pred, minDeg)
	}

	if len(n.children[idx+1].keys) >= minDeg {
		succ := n.getSuccessor(idx)
		n.keys[idx] = succ
		return n.children[idx+1].delete(succ, minDeg)
	}

	n.merge(idx, minDeg)
	return n.children[idx].delete(key, minDeg)
}

func (n *node) getPredecessor(idx int) int {
	cur := n.children[idx]
	for !cur.leaf {
		cur = cur.children[len(cur.children)-1]
	}
	return cur.keys[len(cur.keys)-1]
}

func (n *node) getSuccessor(idx int) int {
	cur := n.children[idx+1]
	for !cur.leaf {
		cur = cur.children[0]
	}
	return cur.keys[0]
}

func (n *node) fill(idx int, minDeg int) {
	if idx > 0 && len(n.children[idx-1].keys) >= minDeg {
		n.borrowFromPrev(idx)
		return
	}
	if idx < len(n.children)-1 && len(n.children[idx+1].keys) >= minDeg {
		n.borrowFromNext(idx)
		return
	}
	if idx < len(n.children)-1 {
		n.merge(idx, minDeg)
	} else {
		n.merge(idx-1, minDeg)
	}
}

func (n *node) borrowFromPrev(idx int) {
	child := n.children[idx]
	sibling := n.children[idx-1]

	child.keys = append([]int{n.keys[idx-1]}, child.keys...)
	if !child.leaf {
		child.children = append([]*node{sibling.children[len(sibling.children)-1]}, child.children...)
		sibling.children = sibling.children[:len(sibling.children)-1]
	}

	n.keys[idx-1] = sibling.keys[len(sibling.keys)-1]
	sibling.keys = sibling.keys[:len(sibling.keys)-1]
}

func (n *node) borrowFromNext(idx int) {
	child := n.children[idx]
	sibling := n.children[idx+1]

	child.keys = append(child.keys, n.keys[idx])
	if !child.leaf {
		child.children = append(child.children, sibling.children[0])
		sibling.children = sibling.children[1:]
	}

	n.keys[idx] = sibling.keys[0]
	sibling.keys = sibling.keys[1:]
}

func (n *node) merge(idx int, _ int) {
	child := n.children[idx]
	sibling := n.children[idx+1]

	child.keys = append(child.keys, n.keys[idx])
	child.keys = append(child.keys, sibling.keys...)

	if !child.leaf {
		child.children = append(child.children, sibling.children...)
	}

	n.keys = append(n.keys[:idx], n.keys[idx+1:]...)
	n.children = append(n.children[:idx+1], n.children[idx+2:]...)
}

// InOrder returns all keys in sorted order.
func (t *BTree) InOrder() []int {
	var result []int
	t.root.inOrder(&result)
	return result
}

func (n *node) inOrder(result *[]int) {
	if n == nil {
		return
	}
	for i, key := range n.keys {
		if !n.leaf {
			n.children[i].inOrder(result)
		}
		*result = append(*result, key)
	}
	if !n.leaf {
		n.children[len(n.keys)].inOrder(result)
	}
}

// Height returns the height of the B-Tree (number of levels).
func (t *BTree) Height() int {
	h := 0
	cur := t.root
	for cur != nil {
		h++
		if cur.leaf || len(cur.children) == 0 {
			break
		}
		cur = cur.children[0]
	}
	return h
}

// String returns a human-readable representation of the tree structure,
// showing each level with its nodes.
func (t *BTree) String() string {
	if t.root == nil {
		return "(empty)"
	}
	var sb strings.Builder
	type queued struct {
		n     *node
		level int
	}
	queue := []queued{{n: t.root, level: 0}}
	prevLevel := -1

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		if item.level != prevLevel {
			if prevLevel >= 0 {
				sb.WriteString("\n")
			}
			fmt.Fprintf(&sb, "Level %d: ", item.level)
			prevLevel = item.level
		} else {
			sb.WriteString("  ")
		}

		sb.WriteString(formatKeys(item.n.keys))

		for _, child := range item.n.children {
			queue = append(queue, queued{n: child, level: item.level + 1})
		}
	}
	return sb.String()
}

func formatKeys(keys []int) string {
	strs := make([]string, len(keys))
	for i, k := range keys {
		strs[i] = strconv.Itoa(k)
	}
	return "[" + strings.Join(strs, ",") + "]"
}
