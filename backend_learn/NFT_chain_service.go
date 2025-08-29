package backendlearn





type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event
	Errors      map[string]Error

	// Additional "special" functions introduced in solidity v0.6.0.
	// It's separated from the original default fallback. Each contract
	// can only define one fallback and receive function.
	Fallback Method // Note it's also used to represent legacy fallback before v0.6.0
	Receive  Method
}


type ChainClient interface {
	FilterLogs(ctx context.Context, q logTypes.FilterQuery) ([]interface{}, error)
	BlockTimeByNumber(context.Context, *big.Int) (uint64, error)
	Client() interface{}
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	CallContractByChain(ctx context.Context, param logTypes.CallParam) (interface{}, error)
	BlockNumber() (uint64, error)
	BlockWithTxs(ctx context.Context, blockNumber uint64) (interface{}, error)
}


type Service struct {
	ctx context.Context

	Abi            *abi.ABI
	HttpClient     *xhttp.Client
	NodeClient     chainclient.ChainClient
	ChainName      string
	NodeName       string
	NameTags       []string
	ImageTags      []string
	AttributesTags []string
	TraitNameTags  []string
	TraitValueTags []string
}

// 会通过配置 [[chain_supported]] 中的 chain_id 做 key 存多个连
nodeSrvs := make(map[int64]*nftchainservice.Service)