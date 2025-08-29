package nftchainservice

import (
	"context"

	"github.com/ProjectsTask/SwapBase/chain/chainclient"
	"github.com/ProjectsTask/SwapBase/xhttp"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Service struct {
	ctx        context.Context
	
	Abi        *abi.ABI
	HttpClient *xhttp.Client
	NodeClient chainclient.ChainClient

	ChainName      string
	NodeName       string
	NameTags       []string
	ImageTags      []string
	AttributesTags []string
	TraitNameTags  []string
	TraitValueTags []string
}
