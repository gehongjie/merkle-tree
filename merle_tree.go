package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
)

//要求使用golang实现一颗默克尔树，叶子结点数不限，实现思路不限
//要求：
//1. 哈希加密采用md5算法
//2. 打印出默克尔根，以及所有叶子结点对应的原始数据
//3. 实现一个接口，输入任意一个节点的原始数据，返回默克尔证明路径
//4. 将实现后的代码提交到github，给出对应可访问的项目地址

type MerkleTree struct {
	Root       *MerkleNode
	merkleRoot []byte
	Leafs      []*MerkleNode
}
type MerkleNode struct {
	Parent   *MerkleNode
	Left     *MerkleNode
	Right    *MerkleNode
	isParent bool
	isRight  bool
	Leaf     bool
	Content  []byte
	Hash     []byte
}

// NewMerkleNode 创建Merkle Tree的节点

func getMerkleTree(contents [][]byte) MerkleTree {
	var merkleTree MerkleTree
	//生成叶节点
	var leafs []*MerkleNode
	// 内容导入到叶节点中
	var hash []byte
	for _, content := range contents {
		//按照要求1，使用md5算法加密
		hash16 := md5.Sum(content)
		hash = hash16[:]
		fmt.Println(hash)
		//构造叶节点，并加入到叶节点集合中
		leafs = append(leafs, &MerkleNode{
			Leaf:    true,
			Content: content,
			Hash:    hash,
		})
	}
	merkleTree.Leafs = leafs
	fmt.Println(leafs)
	//如果叶节点的个数为偶数，就复制最后一个输入的内容的节点，并表明该节点为复制节点
	if len(leafs)%2 == 1 {
		duplicate := &MerkleNode{
			Leaf:    true,
			Content: leafs[len(leafs)-1].Content,
			Hash:    leafs[len(leafs)-1].Hash,
		}
		leafs = append(leafs, duplicate)
	}
	//处理完叶节点以后用递归函数返回根值
	root := findRoot(leafs)
	merkleTree.Root = root
	merkleTree.merkleRoot = root.Hash

	return merkleTree
}

//利用递归函数返回根值
func findRoot(level []*MerkleNode) *MerkleNode {
	//设置节点，用于放置上层的节点
	var nodes []*MerkleNode
	var hash []byte
	for i := 0; i < len(level); i += 2 {
		var left, right int = i, i + 1
		//如果最后只剩一个左节点，在生成父节点时，右节点值=左节点值
		if i+1 == len(level) {
			right = i
		}
		hash16 := md5.Sum(append(level[left].Hash, level[right].Hash...))
		hash = hash16[:]
		n := &MerkleNode{
			Left:  level[left],
			Right: level[right],
			Hash:  hash,
		}
		nodes = append(nodes, n)
		//确定选定的左右节点是否有父节点，并指向父节点，而且确定为左节点还是右节点
		level[left].Parent = n
		level[left].isParent = true
		level[left].isRight = false
		level[right].Parent = n
		level[right].isParent = true
		level[right].isRight = true
		if len(level) == 2 {
			return n
		}
	}
	return findRoot(nodes)

}

//返回验证路径
func prove(tree *MerkleTree, nF []byte) ([][]byte, error) {
	//确定所查内容是否存在
	var notExit bool = true
	var route [][]byte
	var leafTarget *MerkleNode
	for _, leaf := range tree.Leafs {
		if bytes.Equal(leaf.Content, nF) {
			//如果存在则确定目标页节点
			notExit = false
			leafTarget = leaf
		}
	}
	if notExit {
		return nil, errors.New("该树中没找到相关路径")
	}
	//查看当前节点是否有父节点，一直到没有就是最终路径
	route = append(route, leafTarget.Hash)
	//根据左右节点返回另一个节点的hash值
	for leafTarget.isParent {
		if leafTarget.isRight {
			route = append(route, leafTarget.Parent.Left.Hash)

		} else {
			route = append(route, leafTarget.Parent.Right.Hash)
		}
		leafTarget = leafTarget.Parent

	}
	//最后加上根节点
	route = append(route, tree.merkleRoot)
	return route, nil
}

func main() {
	//设定输入的内容
	example := [][]byte{{1, 2, 3}, {4}, {1, 23, 4}, {2, 13, 4}}
	var route [][]byte
	tree := getMerkleTree(example)
	//var aa = MerkleNode{}
	//aa = *tree.Root
	//println(aa.Hash)
	//实现要求
	//要求1已经在具体代码中实现
	//要求2
	fmt.Println("根节点的哈希值：", tree.merkleRoot)
	for _, leaf := range tree.Leafs {
		fmt.Println(leaf.Content)
	}
	//要求3,根据某个内容，给出默克尔证明路径
	//设定查询内容
	exampleFind := []byte{1, 2, 3}

	route, err := prove(&tree, exampleFind)
	if err != nil {
		println(err)
	}
	fmt.Println("证明路径：", route)

}
