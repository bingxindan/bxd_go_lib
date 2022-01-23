package trees

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"testing"
)

// 定义我们自己的菜单对象
type SystemMenu struct {
	CollectionId int `json:"collectionId"`
	ParentId     int `json:"parentId"`
}

// region 实现ITree 所有接口
func (s SystemMenu) GetTitle() string {
	return ""
}

func (s SystemMenu) GetId() int {
	return s.CollectionId
}

func (s SystemMenu) GetParentId() int {
	return s.ParentId
}

func (s SystemMenu) GetData() interface{} {
	return s
}

func (s SystemMenu) IsRoot() bool {
	// 这里通过ParentId等于0 或者 ParentId等于自身Id表示顶层根节点
	return s.ParentId == 0 || s.ParentId == s.CollectionId
}

// endregion

type SystemMenus []SystemMenu

// ConvertToINodeArray 将当前数组转换成父类 INode 接口 数组
func (s SystemMenus) ConvertToINodeArray() (nodes []INode) {
	for _, v := range s {
		nodes = append(nodes, v)
	}
	return
}

func TestGenerateTree(t *testing.T) {
	// 模拟获取数据库中所有菜单，在其它所有的查询中，也是首先将数据库中所有数据查询出来放到数组中，
	// 后面的遍历递归，都在这个 allMenu中进行，而不是在数据库中进行递归查询，减小数据库压力。
	allMenu := []SystemMenu{
		{CollectionId: 1, ParentId: 0},
		{CollectionId: 2, ParentId: 0},

		{CollectionId: 3, ParentId: 1},
		{CollectionId: 4, ParentId: 1},

		{CollectionId: 5, ParentId: 2},
		{CollectionId: 6, ParentId: 3},
		{CollectionId: 7, ParentId: 3},
	}

	fmt.Printf("all: %+v\n", allMenu)

	// 生成完全树
	resp := CustomTree(SystemMenus.ConvertToINodeArray(allMenu))
	bytes, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Println(string(pretty.Color(pretty.PrettyOptions(bytes, pretty.DefaultOptions), nil)))
}
