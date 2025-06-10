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
	return cause.Map{
		"name":  cause.Required(n.Name),
		"node":  cause.Optional(n.Node),
		"nodes": cause.Optional(n.Nodes).When(len(n.Nodes) > 3, "too many nodes"),
	}.Err()
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
		Nodes: []*Node{{}, {Name: "X"}, {Name: "Y"}},
		Node: &Node{
			Node: &Node{Nodes: []*Node{{Name: "D"}, {Name: "E"}, {}}},
			Name: "B",
			Nodes: []*Node{
				{Name: "C"},
				{Name: "D"},
				{Name: "E"},
				{Name: "F"}, // This will trigger the "too many nodes" error.
				{},          // This will be ignored
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
	//       "name": "required",
	//       "nodes[2]": {
	//         "name": "required"
	//       }
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
