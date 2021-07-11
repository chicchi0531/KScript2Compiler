package ast

import(
	"ks2/compiler/vm"
	cm "ks2/compiler/common"
)

type NFunction struct{
	Node
	name string
	args []vm.INode
}

func MakeFunctionNode(lineno int, name string, args []vm.INode, driver *vm.Driver) *NFunction{
	n := new(NFunction)
	n.name = name
	n.args = args
	n.Lineno = lineno
	n.Driver = driver
	return n
}

func (n *NFunction) Push() int{
	f := n.Driver.FunctionTable.Find(n.name)
	if f != nil{
		// 引数の数チェック
		if len(f.Args) != len(n.args){
			n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0020,"関数："+f.Name)
			return f.RetrunType
		}

		// 引数逆積み
		for i:=len(n.args)-1; i>=0; i--{
			argType := n.args[i].Push()
			// 引数型チェック
			if argType != f.Args[i].VarType{
				n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0021, "関数："+f.Name)
				return f.RetrunType
			}
		}
		// 引数の数積み
		n.Driver.OpPushInteger(len(f.Args))
		// call
		n.Driver.OpCall(f.Address)

		return f.RetrunType
	}
	//関数が見つからなかったらエラー
	n.Driver.Err.LogError(n.Driver.Filename, n.Lineno, cm.ERR_0022, "関数："+n.name)
	return -1
}
