package response

import "gitlab.com/ucard/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
