package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello,B+tree!")
	node := BNode(make([]byte, BTREE_PAGE_SIZE))
	node.setHeader(BNODE_LEAF, 2)
	nodeAppendKV(node, 0, 0, []byte("key1"), []byte("Hello"))
	nodeAppendKV(node, 1, 0, []byte("key2"), []byte("World"))
	idxkey := nodeLookupLE(node, []byte("key1"))
	fmt.Printf("idxkey: %d\n", idxkey)
	fmt.Printf("key: %s, val: %s\n", node.getkey(idxkey), node.getVal(idxkey))
	//print every key and value in the node
	for i := uint16(0); i < node.nkeys(); i++ {
		fmt.Printf("key: %s, val: %s\n", node.getkey(i), node.getVal(i))
	}
	//find and update
	key := []byte("key1")
	newVal := []byte("Hi")
	idx := nodeLookupLE(node, key)
	if string(node.getkey(idx)) == string(key) {
		newNode := BNode(make([]byte, BTREE_PAGE_SIZE))
		leafUpdate(newNode, node, idx, key, newVal)
	}
}
