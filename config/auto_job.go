package config

// AutoJobEnabled 渠道自动任务开关。未配置 enable 时默认开启。
func AutoJobEnabled(enable *bool) bool {
	if enable == nil {
		return true
	}
	return *enable
}
