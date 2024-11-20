package ext

import "cmp"

const _Tiers = 4

type _TreeNode[K cmp.Ordered, V any] struct {
	kvs  [_Tiers - 1]KV[K, V]
	next [_Tiers]*_TreeNode[K, V]
}

func (n *_TreeNode[K, V]) isLeaf() bool {
	return n.next[0] == nil
}
