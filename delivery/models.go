package delivery

type NodeInfoResponse struct {
	NodeName    string `json:"node_name"`
	XrayRunning bool   `json:"is_xray_running"`
}

type AddNodeRequest struct {
	NodeAddress string `json:"node_address"`
}

type UpdateConfig struct {
	NodeName  string `json:"node_name"`
	NewConfig string `json:"new_config"`
}

type GetCurrentConfigRequest struct {
	NodeName string `json:"node_name"`
}
type GetCurrentConfigResponse struct {
	NodeName      string `json:"node_name"`
	CurrentConfig string `json:"current_config"`
}
