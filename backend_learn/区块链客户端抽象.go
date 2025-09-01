




// 定义了一个区块链客户端抽象层（ChainClient接口）及其工厂函数（New），主要用于与以太坊及兼容链（如Optimism、Sepolia）进行交互。
// 其核心功能是封装底层链的差异，提供统一的区块链数据访问接口
type ChainClient interface {
	// 查询符合过滤条件的交易日志（如智能合约事件）。
	// 类似以太坊的 eth_getLogs JSON-RPC 方法，用于监听合约事件或分析链上活动
	FilterLogs(ctx context.Context, q logTypes.FilterQuery) ([]interface{}, error)
	// 根据区块号获取区块的时间戳。
	BlockTimeByNumber(context.Context, *big.Int) (uint64, error)
	// 返回底层客户端实例（如以太坊的 ethclient.Client），用于直接调用链特定功能。
	Client() interface{}
	// 执行只读的合约调用（不消耗Gas），返回合约函数结果
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	// 扩展的合约调用方法，支持自定义参数（如指定链类型）
	CallContractByChain(ctx context.Context, param logTypes.CallParam) (interface{}, error)
	// 获取当前最新区块号
	BlockNumber() (uint64, error)
	// 获取包含完整交易列表的区块数据。
	BlockWithTxs(ctx context.Context, blockNumber uint64) (interface{}, error)
}


// 通过统一接口屏蔽不同EVM链（以太坊、Optimism等）的底层实现差异，便于业务代码跨链复用
func New(chainID int, nodeUrl string) (ChainClient, error) {
	switch chainID {
	case chain.EthChainID, chain.OptimismChainID, chain.SepoliaChainID:
		return evmclient.New(nodeUrl)
	default:
		return nil, errors.New("unsupported chain id")
	}
}

/*
核心使用场景​​：
1. 监听合约事件​​：如DeFi应用通过 FilterLogs 实时追踪代币转账或交易完成事件。
2.查询链上数据​​：获取区块时间、交易详情或合约状态（如代币余额）。
3.合约交互​​：执行无状态的合约查询（CallContract）或发送交易（需结合签名逻辑）。
扩展性​​：当前仅支持EVM链，但工厂模式便于未来扩展至非EVM链（如通过新增实现类）。
调用方只需要关心 ChainClient 接口，不需要知道具体实现，将对象创建与使用分离，客户端代码与具体实现解耦

是个工厂模式 ， `New` 函数是一个工厂函数   为什么：
	工厂模式是一种创建型设计模式，它提供了一种创建对象的最佳方式，而无需暴露创建逻辑
	在这个例子中，工厂模式被用来根据不同的链ID创建相应的区块链客户端。
	抽象接口：工厂返回的是 `ChainClient` 接口类型，而不是具体的实现。这意味着客户端代码只依赖于接口，而不是具体实现，从而实现了松耦合
*/