package blockchain

type IO interface {
	ForgeData(address []byte, data []byte) error
	FetchData(address []byte) ([]byte, error)
}

func Factory(prefix string) IO {
	switch prefix {
	case Waves_prefix:
		return Waves{}
	default:
		return nil
	}
}
