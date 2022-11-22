package dag

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestExecute(t *testing.T) {
	g := &Graph{
		Nodes: []Node{0, 1, 2, 3, 4, 5},
		Edges: []Edge{
			{Depender: 1, Dependee: 0},
			{Depender: 2, Dependee: 0},
			{Depender: 3, Dependee: 2},
			{Depender: 4, Dependee: 3},
			{Depender: 5, Dependee: 4},
		},
	}
	_ = Execute(g, func(node Node) error {
		time.Sleep(3 * time.Second)
		fmt.Println(int(node))
		if int(node) == 3 {
			return errors.New("error")
		}
		return nil
	})
	//assert.NotNil(t, err)
}
