package nodes

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NodeManager struct {
	mu    sync.Mutex
	nodes map[string]*Node
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodes: make(map[string]*Node),
	}
}

func (m *NodeManager) AddNode(ctx context.Context, address string) (*NodeInfo, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	node, err := InitializeNode(ctx, conn)
	if err != nil {
		return nil, err
	}

	info, err := node.GetInfo(ctx)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes[node.NodeName] = node

	return info, nil
}

func (m *NodeManager) ListNodes() []*Node {
	m.mu.Lock()
	defer m.mu.Unlock()

	nodes := make([]*Node, 0, len(m.nodes))
	for _, n := range m.nodes {
		nodes = append(nodes, n)
	}
	return nodes
}

func (m *NodeManager) ListNodesInfo() ([]*NodeInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	infos := make([]*NodeInfo, 0, len(m.nodes))
	for _, node := range m.nodes {
		info, err := node.GetInfo(context.Background())
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}

	return infos, nil
}

func (m *NodeManager) RemoveNode(nodeName string) error {
	m.mu.Lock()
	node, ok := m.nodes[nodeName]
	if ok {
		delete(m.nodes, nodeName)
	}
	m.mu.Unlock()

	if !ok {
		return ErrNodeNotExists
	}

	return node.NodeConn.Close()
}

func (m *NodeManager) GetNodeByName(ctx context.Context, nodeName string) (*Node, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	node, ok := m.nodes[nodeName]
	if !ok {
		return nil, ErrNodeNotExists
	}
	return node, nil
}
