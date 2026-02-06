package nodes_test

import (
	"context"
	"gofronet-foundation/gofro-control/nodes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodesManager(t *testing.T) {

	t.Run("Add Node", func(t *testing.T) {
		nodesManager := nodes.NewNodeManager()

		result, err := nodesManager.AddNode(context.Background(), "147.45.214.213:50051")
		require.NoError(t, err)

		t.Log(result)
	})

	t.Run("Get Nodes", func(t *testing.T) {
		nodesManager := nodes.NewNodeManager()

		_, err := nodesManager.AddNode(context.Background(), "147.45.214.213:50051")
		require.NoError(t, err)

		nodes := nodesManager.ListNodes()
		for _, node := range nodes {
			info, err := node.GetInfo(context.Background())
			require.NoError(t, err)

			t.Log(info)
		}

	})

	t.Run("Get List Nodes Info", func(t *testing.T) {
		nodesManager := nodes.NewNodeManager()

		_, err := nodesManager.AddNode(context.Background(), "147.45.214.213:50051")
		require.NoError(t, err)

		nodes, err := nodesManager.ListNodesInfo()
		require.NoError(t, err)
		for _, nodeInfo := range nodes {
			t.Log(nodeInfo)
		}
	})

}
