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
		nodeRepalceKidN(tree, new, node, idx, split[:nsplit]...)
	}
	return new
}
