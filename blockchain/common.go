package blockchain

type IO interface {
	ForgeData(address []byte, data []byte) error
	FetchData(address []byte) ([]byte, error)
}

var (
	bdtp_env_seed = "BDTP_SEED"
)
