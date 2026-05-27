package main

import "bytes"

//Btree struct
type BTree struct {
	//root pointer
	root uint64
	//callback for managing the on-disk pages
	get func(uint64) []byte
	new func([]byte) uint64 //allocate a new page number with data
	del func(uint64)        //deallocate a page number
}

//Btree node insert
func treeInsert(tree *BTree, node BNode, key []byte, val []byte) BNode {
	//init the BNode
	new := BNode(make([]byte, BTREE_PAGE_SIZE*2))
	//find if the key already exists in the node
	idx := nodeLookupLE(node, key)
	//jude if the node is a leaf node or not
	switch node.btype() {
	case BNODE_LEAF:
		//if the key already exists,update the value
		if bytes.Equal(node.getkey(idx), key) {
			leafUpdate(new, node, idx, key, val)
		} else {
			//if the key does not exist,insert the key and value into the node
			leafInsert(new, node, idx, key, val)
		}
	case BNODE_NODE:
		//recursively insertion to the child node
		kptr := node.getPtr(idx)
		//first insert the key and value into the child node
		knode := treeInsert(tree, tree.get(kptr), key, val)
		//then check if the child node is split or not
		//after insertion,split the node
		nsplit, split := nodeSplit3(knode)
		//deallocation the old kid node
		tree.del(kptr)
		//update the kid links
		nodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
	}
	return new
}

//a function to replace the kid nodes of a node with new kid nodes
func nodeReplaceKidN(tree *BTree, new BNode, old BNode, idx uint16, kid ...BNode) {
	inc := uint16(len(kid))
	new.setHeader(BNODE_NODE, old.nkeys()+inc-1)
	//copy the old node to the new node
	nodeAppendRange(new, old, 0, 0, idx)
	//append the new kid nodes to the new node
	for i, node := range kid {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getkey(0), nil)
	}
	nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))
}

//High level API for BTree
//Btree insert
func (tree *BTree) Insert(key []byte, val []byte) error {
	//check the length limit imposed by the node format
	if err := checkLimit(key, val); err != nil {
		return err
	}

	//if tree is empty
	if tree.root == 0 {
		//create a new leaf node
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		//set the
		tree.root = tree.new(root) //allocate a new page number with data
		return nil
	}

	//get the root node
	node := treeInsert(tree, tree.get(tree.root), key, val)
	//if grow the tree,split the root node
	//nsplit is the number of new nodes
	nsplit, split := nodeSplit3(node)
	//deallocation the old root node
	tree.del(tree.root)
	if nsplit > 1 {
		//create a new root node
		root := BNode(make([]byte, BTREE_PAGE_SIZE*2))

		//set the new root node
		tree.root = tree.new(root)
	} else {
		tree.root = tree.new(split[0])
	}
	return nil
}
