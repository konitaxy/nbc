package utils

// CardProviderErrorMessage 从 cardbin / gzy 上游 API 错误中提取可展示文案。
func CardProviderErrorMessage(err error) string {
	return ProviderUserMessage(err)
}
