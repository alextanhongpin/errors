package cause_test

import (
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/errors/cause"
)

type Node struct {
	Name  string
	Node  *Node
	Nodes []*Node
}

func (n *Node) Validate() error {
	return cause.NewMapValidator().
		Required("name", n.Name).
		Optional("node", n.Node).
		Optional("nodes", len(n.Nodes), cause.When(len(n.Nodes) > 3, "too many nodes")).
		Optional("nodes", cause.Slice(n.Nodes)).
		Validate()
}

func ExampleFields_node_valid() {
	n := &Node{
		Name: "A",
	}
	validateNode(n)

	// Output:
	// is nil: true
	// err: <nil>
	// null
}

func ExampleFields_node_invalid_name() {
	n := &Node{
		Name: "",
	}
	validateNode(n)

	// Output:
	// is nil: false
	// err: invalid fields: name
	// {
	//   "name": "required"
	// }
}

func ExampleFields_node_invalid_node() {
	n := &Node{
		Node: &Node{Name: ""},
	}
	validateNode(n)

	// Output:
	// is nil: false
	// err: invalid fields: name, node
	// {
	//   "name": "required",
	//   "node": {
	//     "name": "required"
	//   }
	// }
}

func ExampleFields_node_invalid_nested_node() {
	n := &Node{
		Node: &Node{
			Node: &Node{},
		},
	}
	validateNode(n)

	// Output:
	// is nil: false
	// err: invalid fields: name, node
	// {
	//   "name": "required",
	//   "node": {
	//     "name": "required",
	//     "node": {
	//       "name": "required"
	//     }
	//   }
	// }
}

func ExampleFields_nodes_invalid() {
	n := &Node{
		Nodes: []*Node{&Node{}, &Node{Name: "X"}, &Node{Name: "Y"}},
		Node: &Node{
			Node: &Node{},
			Name: "B",
			Nodes: []*Node{
				{Name: "C"},
				{Name: "D"},
				{Name: "E"},
				{Name: "F"}, // This will trigger the "too many nodes" error.
				{},
			},
		},
	}
	validateNode(n)

	// Output:
	// is nil: false
	// err: invalid fields: name, node, nodes[0]
	// {
	//   "name": "required",
	//   "node": {
	//     "node": {
	//       "name": "required"
	//     },
	//     "nodes": "too many nodes",
	//     "nodes[4]": {
	//       "name": "required"
	//     }
	//   },
	//   "nodes[0]": {
	//     "name": "required"
	//   }
	// }
}

func validateNode(n *Node) {
	err := n.Validate()
	fmt.Println("is nil:", err == nil)
	fmt.Println("err:", err)

	b, err := json.MarshalIndent(err, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
