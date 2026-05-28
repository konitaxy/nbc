package constant

// ChainType 链类型（链上收款监听）。
type ChainType string

const (
	ChainType_TRON ChainType = "TRON"
)

func (c ChainType) String() string {
	return string(c)
}
