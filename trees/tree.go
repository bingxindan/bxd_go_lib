package trees

// Tree 统一定义菜单树的数据结构，也可以自定义添加其他字段
type Tree struct {
	Title           string      `json:"title"`            //节点名字
	Data            interface{} `json:"data"`             //自定义对象
	Leaf            bool        `json:"leaf"`             //叶子节点
	//Selected        bool        `json:"checked"`          //选中
	//PartialSelected bool        `json:"partial_selected"` //部分选中
	Children        []Tree      `json:"children"`         //子节点
}

// ConvertToINodeArray 其他的结构体想要生成菜单树，直接实现这个接口
type INode interface {
	// GetTitle 获取显示名字
	GetTitle() string
	// GetId获取id
	GetId() int
	// GetFatherId 获取父id
	GetFatherId() int
	// GetData 获取附加数据
	GetData() interface{}
	// IsRoot 判断当前节点是否是顶层根节点
	IsRoot() bool
}

type INodes []INode

func (nodes INodes) Len() int {
	return len(nodes)
}

func (nodes INodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}

func (nodes INodes) Less(i, j int) bool {
	return nodes[i].GetId() < nodes[j].GetId()
}

// GenerateTree 自定义的结构体实现 INode 接口后调用此方法生成树结构
// nodes 需要生成树的节点
// selectedNode 生成树后选中的节点
// menuTrees 生成成功后的树结构对象
func CustomTree(nodes []INode) (trees []Tree) {
	var roots, childes []INode

	for _, v := range nodes {
		if v.IsRoot() {
			roots = append(roots, v)
		}
		childes = append(childes, v)
	}

	for _, v := range roots {
		childTree := &Tree{
			Title: v.GetTitle(),
			Data: v.GetData(),
		}
		// 递归
		recursiveTree(childTree, childes)
		// 递归之后，根据子节点确认是否是叶子节点
		childTree.Leaf = len(childTree.Children) == 0
		trees = append(trees, *childTree)
	}

	return
}

// recursiveTree 递归生成树结构
// tree 递归的树对象
// nodes 递归的节点
// selectedNodes 选中的节点
func recursiveTree(tree *Tree, nodes []INode) {
	data := tree.Data.(INode)

	for _, v := range nodes {
		if v.IsRoot() {
			continue
		}
		if data.GetId() == v.GetFatherId() {
			childTree := &Tree{
				Title: v.GetTitle(),
				Data: v.GetData(),
			}
			recursiveTree(childTree, nodes)
			// 递归之后，根据子节确认是否是叶子节点
			childTree.Leaf = len(childTree.Children) == 0
			tree.Children = append(tree.Children, *childTree)
		}
	}
}