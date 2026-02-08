package nodes

import (
	"context"
	apiv1 "gofronet-foundation/gofro-control/gen/api/v1"

	"google.golang.org/grpc"
)

type Node struct {
	NodeName string
	NodeConn *grpc.ClientConn
}

type NodeData struct {
	NodeName    string `json:"node_name"`
	NodeAddress string `json:"node_address"`
}

type NodeInfo struct {
	NodeName      string
	IsXrayRunning bool
}

func InitializeNode(ctx context.Context, conn *grpc.ClientConn) (*Node, error) {
	node := &Node{
		NodeConn: conn,
	}

	info, err := node.GetInfo(ctx)
	if err != nil {
		return nil, err
	}

	node.NodeName = info.NodeName
	return node, nil

}

func (node *Node) GetInfo(ctx context.Context) (*NodeInfo, error) {
	client := apiv1.NewXrayServiceClient(node.NodeConn)
	resp, err := client.GetNodeInfo(ctx, &apiv1.GetNodeInfoRequest{})
	if err != nil {
		return nil, err
	}

	return &NodeInfo{
		NodeName:      resp.NodeName,
		IsXrayRunning: resp.XrayRunning,
	}, nil
}

func (node *Node) UpdateConfig(ctx context.Context, newConfig string) error {
	client := apiv1.NewXrayServiceClient(node.NodeConn)
	_, err := client.UpdateXrayConfig(ctx, &apiv1.UpdateXrayConfigRequest{
		NewConfig: newConfig,
	})
	if err != nil {
		return err
	}
	return nil
}

func (node *Node) GetCurrentConfig(ctx context.Context) (string, error) {
	client := apiv1.NewXrayServiceClient(node.NodeConn)
	resp, err := client.GetCurrentConfig(ctx, &apiv1.GetCurrentConfigRequest{})
	if err != nil {
		return "", err
	}
	return resp.CurrentConfig, nil
}

func (node *Node) Restart(ctx context.Context) error {
	client := apiv1.NewXrayServiceClient(node.NodeConn)
	_, err := client.RestartXray(ctx, &apiv1.RestartXrayRequest{})
	return err
}

func (node *Node) Start(ctx context.Context) error {
	client := apiv1.NewXrayServiceClient(node.NodeConn)
	_, err := client.StartXray(ctx, &apiv1.StartXrayRequest{})
	return err
}

func (node *Node) Stop(ctx context.Context) error {
	client := apiv1.NewXrayServiceClient(node.NodeConn)
	_, err := client.StopXray(ctx, &apiv1.StopXrayRequest{})
	return err
}
