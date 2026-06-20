package kdt

type KdTreeNode struct {
	data  []int
	right *KdTreeNode
	left  *KdTreeNode
}

type KdTree struct {
	k    int
	root *KdTreeNode
}

func NewKdTree(k int) *KdTree {
	return &KdTree{k: k, root: nil}
}

func (kdt *KdTree) Insert(v []int) {
	isTreeEmpty := kdt.root == nil
	if isTreeEmpty {
		kdt.root = &KdTreeNode{
			data:  v,
			right: nil,
			left:  nil,
		}
		return
	}

	current := kdt.root
	depth := 0

	for current != nil {
		dimension := depth % kdt.k
		if v[dimension] > current.data[dimension] {
			if current.right == nil {
				current.right = &KdTreeNode{
					data:  v,
					right: nil,
					left:  nil,
				}
				return
			}
			current = current.right
		} else {
			if current.left == nil {
				current.left = &KdTreeNode{
					data:  v,
					right: nil,
					left:  nil,
				}
				return
			}
			current = current.left
		}
		depth++
	}
}
