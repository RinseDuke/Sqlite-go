package main

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
func treeInsert(tree *BTree, node BNode, key []byte, val []byte) {
	//find if the key already exists in the node
	idx := nodeLookupLE(node, key)
	//jude if the node is a leaf node or not
	switch node.btype() {

	}
}
