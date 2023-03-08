package index

type Node struct {
	id     int64
	offset int64
	left   *Node
	right  *Node
}

type BinaryTreeIndex struct {
	root *Node
	size int64
}

func NewBinaryTreeIndex() *BinaryTreeIndex {
	return &BinaryTreeIndex{
		root: nil,
		size: 0,
	}
}

func (tree *BinaryTreeIndex) Contains(id int64) bool {
	var found bool = false

	tmp := &Node{
		id: id,
	}

	ptr := tree.root
	for {
		if ptr.id == id {
			found = true
			break
		}

		if tmp.compare(ptr) == 1 {
			if ptr.right == nil {
				break
			}

			ptr = ptr.right
		} else {
			if ptr.left == nil {
				break
			}

			ptr = ptr.left
		}
	}

	return found
}

// op methods: insert, find, compare
func (tree *BinaryTreeIndex) Find(id int64) int64 {
	var found bool = false
	var value int64

	tmp := &Node{
		id: id,
	}

	ptr := tree.root
	for {
		if ptr.id == id {
			found = true
			value = ptr.offset
			break
		}

		if tmp.compare(ptr) == 1 {
			if ptr.right == nil {
				break
			}

			ptr = ptr.right
		} else {
			if ptr.left == nil {
				break
			}

			ptr = ptr.left
		}
	}

	if found {
		return value
	}

	return -1
}

func (tree *BinaryTreeIndex) Insert(id, offset int64) {
	nn := &Node{
		id:     id,
		offset: offset,
		left:   nil,
		right:  nil,
	}

	if tree.root == nil {
		tree.root = nn
		tree.size++

		return
	}

	ptr := tree.root
	for {
		if nn.compare(ptr) == 1 {
			if ptr.right == nil {
				break
			}

			ptr = ptr.right
		} else {
			if ptr.left == nil {
				break
			}

			ptr = ptr.left
		}
	}

	if nn.compare(ptr) == 1 {
		ptr.right = nn
	} else {
		ptr.left = nn
	}

	tree.size++
}

// Compare returns:
// 1 if node.id > node2.id
// 0 if node.id < node2.id
func (node *Node) compare(node2 *Node) int {
	if node.id > node2.id {
		return 1
	}

	return 0
}

// implement index
func (tree *BinaryTreeIndex) Size() int64 {
	return tree.size
}

func (tree *BinaryTreeIndex) Has(id int64) bool {
	return tree.Contains(id)
}

func (tree *BinaryTreeIndex) Get(id int64) int64 {
	return tree.Find(id)
}

func (tree *BinaryTreeIndex) Put(id, offset int64) {
	tree.Insert(id, offset)
}

func (tree *BinaryTreeIndex) Type() string {
	return "binary-search-tree"
}
