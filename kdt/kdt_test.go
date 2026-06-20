package kdt

import (
	"slices"
	"testing"
)

func TestNewKdTree(t *testing.T) {
	k := 2
	kdt := NewKdTree(k)

	if kdt.root != nil {
		t.Error("expected kdt.root to be nil")
	}

	if kdt.k != 2 {
		t.Error("expected kdt.k to be 2")
	}
}

func TestInsertWithK2(t *testing.T) {
	k := 2
	kdt := NewKdTree(k)

	root := []int{7, 8}
	node1 := []int{12, 3}
	node2 := []int{14, 1}
	node3 := []int{4, 12}
	node4 := []int{9, 1}
	node5 := []int{2, 7}
	node6 := []int{10, 9}

	kdt.Insert(root)
	kdt.Insert(node1)
	kdt.Insert(node2)
	kdt.Insert(node3)
	kdt.Insert(node4)
	kdt.Insert(node5)
	kdt.Insert(node6)

	if !slices.Equal(kdt.root.data, root) {
		t.Error("[depth=0, k=1] Invalid root received")
	}

	if !slices.Equal(kdt.root.right.data, node1) {
		t.Error("[depth=1, k=2] Invalid right node")
	}

	if !slices.Equal(kdt.root.right.left.data, node2) {
		t.Error("[depth=2, k=1] Invalid left node")
	}

	if !slices.Equal(kdt.root.left.data, node3) {
		t.Error("[depth=1, k=2] Invalid left node")
	}

	if !slices.Equal(kdt.root.right.left.left.data, node4) {
		t.Error("[depth=3, k=2] Invalid left node")
	}

	if !slices.Equal(kdt.root.left.left.data, node5) {
		t.Error("[depth=2, k=1] Invalid left node")
	}

	if !slices.Equal(kdt.root.right.right.data, node6) {
		t.Error("[depth=2, k=1] Invalid right node")
	}
}

func TestInsertWithK3(t *testing.T) {
	k := 3
	kdt := NewKdTree(k)

	node1 := []int{5, 1, 1}
	node2 := []int{1, 5, 1}
	node3 := []int{4, 4, 4}
	node4 := []int{2, 2, 5}
	node5 := []int{6, 6, 2}

	kdt.Insert([]int{3, 3, 3})
	kdt.Insert(node1)
	kdt.Insert(node2)
	kdt.Insert(node3)
	kdt.Insert(node4)
	kdt.Insert(node5)

	root := kdt.root
	if !slices.Equal(root.data, []int{3, 3, 3}) {
		t.Error("[depth=0, dim=0] Invalid root")
	}

	if !slices.Equal(root.right.data, node1) {
		t.Error("[depth=1, dim=1] Invalid right child")
	}

	if !slices.Equal(root.left.data, node2) {
		t.Error("[depth=1, dim=1] Invalid left child")
	}

	if !slices.Equal(root.right.right.data, node3) {
		t.Error("[depth=2, dim=2] Invalid right.right child")
	}

	if !slices.Equal(root.left.left.data, node4) {
		t.Error("[depth=2, dim=2] Invalid left.left child")
	}

	if !slices.Equal(root.right.right.left.data, node5) {
		t.Error("[depth=3, dim=0] Invalid right.right.left child")
	}
}
