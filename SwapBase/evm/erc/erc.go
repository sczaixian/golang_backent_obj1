package erc

type Erc interface {
	GetItemOwner(address string, tokenId string) (string, error)
}

type NftErc struct {
	Endpoint string `toml:"endpoint" json:"endpoint"`
	Standard string `toml:"standard" json:"standard"`
}
