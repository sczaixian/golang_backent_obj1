


// 使用了时间轮结构


// Activity ：collection order-activity
type Order struct {
	// order Id
	orderID        string
	CollectionAddr string
	// chain suffix name: ethw/bsc
	ChainSuffix string
	// expireIn - createIn (unit: s)
	CycleCount int64
	// position of the task on the time wheel
	WheelPosition int64

	Next *Order
}

type wheel struct {
	// linked list
	NotifyActivities *Order
}

//  wheel  order 组成的一个 轮子

type OrderManager struct {
	chain string

	// cycle time wheel
	TimeWheel [WheelSize]wheel   
	// current time wheel index
	CurrentIndex int64

	collectionOrders map[string]*collectionTradeInfo

	collectionListedCh chan string
	project            string

	Xkv *xkv.Store
	DB  *gorm.DB
	Ctx context.Context
	Mux *sync.RWMutex
}



ordermanager.New(ctx, db, kvStore, cfg.ChainCfg.Name, cfg.ProjectCfg.Name)























//------------------- 订单簿 -----------------------------------------
// Order 表示一个交易订单的核心信息，包含买卖方向、资产详情、价格及有效期等关键数据。
type Order struct {
	Side     uint8          // 订单方向，例如：0=买入，1=卖出。
	SaleKind uint8          // 销售类型，例如：0=固定价格，1=拍卖。
	Maker    common.Address // 订单创建者的区块链地址。
	
	// Nft 描述订单关联的非同质化代币（NFT）资产详情。
	Nft      struct {
		TokenId        *big.Int         // NFT 的唯一标识符。
		CollectionAddr common.Address   // NFT 所属智能合约的地址。
		Amount         *big.Int         // 交易数量（适用于ERC1155等多数量代币）
	}

	Price  *big.Int // 订单价格（以最小单位表示，如Wei）
	Expiry uint64   // 订单过期时间戳（Unix时间）
	Salt   uint64   // 随机数，用于保证订单哈希的唯一性
}



// IndexedStatus 索引状态记录结构体，用于跟踪区块链数据同步进度。
type IndexedStatus struct {
	Id                // 主键 ID（自增）
	ChainId               // 区块链类型（1: 以太坊），默认值为 1
	LastIndexedBlock  // 最后已索引的区块高度（允许为空）
	LastIndexedTime    // 最后索引完成的时间戳（毫秒，允许为空）
	IndexType         // 索引类型：0=活动索引，1=交易信息索引
	CreateTime       // 记录创建时间（自动填充毫秒时间戳）
	UpdateTime       // 记录更新时间（自动更新毫秒时间戳）
}




// 启动了两个独立的、永不停止的循环协程（Loop），直到服务被关闭（通过 context 取消）
func (s *Service) Start() {
	threading.GoSafe(s.SyncOrderBookEventLoop)
	threading.GoSafe(s.UpKeepingCollectionFloorChangeLoop)
}


// 等待一个  当前最高块 - 阈值  防止主链分叉的情况
// currentBlockNum-MultiChainMaxBlockDifference[s.chain]
if lastSyncBlock > currentBlockNum-MultiChainMaxBlockDifference[s.chain] { // 如果上次同步的区块高度大于当前区块高度，等待一段时间后再次轮询
	time.Sleep(SleepInterval * time.Second)
	continue
}



[contract_cfg]
eth_address = "0x0000000000000000000000000000000000000000"
weth_address = "0x4200000000000000000000000000000000000006"
dex_address = "0x5560e1c2E0260c2274e400d80C30CDC4B92dC8ac" # undeploy
1. `eth_address = "0x0000000000000000000000000000000000000000"`
   - 这通常代表以太坊上的原生代币ETH的地址。在以太坊上，ETH本身没有合约地址，
   所以通常用零地址（0x0）来表示原生货币。
2. `weth_address = "0x4200000000000000000000000000000000000006"`
   - 这是WETH（Wrapped ETH）的合约地址。WETH是以太坊上ETH的封装版本，使其符合ERC-20标准。
   注意这个地址看起来像是在一个Layer2网络（比如Optimism）上，因为0x420...006是Optimism上WETH的常见地址。
3. `dex_address = "0x5560e1c2E0260c2274e400d80C30CDC4B92dC8ac" # undeploy`
   - 这应该是去中心化交易所（DEX）的合约地址。
   但是后面的注释`# undeploy`表明这个合约已经被取消部署（可能不再使用，或者是一个测试用的合约，现在已经撤销了）




// 用于以太坊（或兼容以太坊的区块链）日志过滤查询的结构体。它的主要作用是指定过滤条件，用于检索智能合约产生的事件日志
type FilterQuery struct {
	BlockHash string   // used by eth_getLogs, return logs only from block with this hash
	FromBlock *big.Int // beginning of the queried range, nil means genesis block
	ToBlock   *big.Int // end of the range, nil means latest block
	Addresses []string // restricts matches to events created by specific contracts

	// The Topic list restricts matches to particular event topics. Each event has a list
	// of topics. Topics matches a prefix of that list. An empty element slice matches any
	// topic. Non-empty elements represent an alternative that matches any of the
	// contained topics.
	//
	// Examples:
	// {} or nil          matches any topic list  不过滤任何主题。匹配所有事件。
	// {{A}}              matches topic A in first position 这通常用于匹配特定事件A
	// {{}, {B}}          matches any topic in first position AND B in second position 第二个主题（topic[1]） 必须等于 B 的事件
	// {{A}, {B}}         matches topic A in first position AND B in second position
	// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position
	Topics [][]string   // 根据事件的主题（Topics） 进行过滤。
	// opic[0]: 总是事件签名的 Keccak-256 哈希值（例如 Keccak256("Transfer(address,address,uint256)"）。
	// 这是识别事件类型的唯一标识符。
}

// --------------  应用  ----------------
// 查找特定合约的 Transfer 事件
query := FilterQuery{
    FromBlock: big.NewInt(1000000),
    ToBlock:   big.NewInt(1001000),
    Addresses: []string{"0xContractAddress"},
    Topics: [][]string{
        {"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}, // Transfer 事件签名
        {"0xFromAddress"}, // 可选的发送方地址
        {"0xToAddress"},   // 可选的接收方地址
    },
}
典型应用场景：
	监控特定合约的事件
	追踪代币转账记录
	监听DeFi协议的状态变化
	区块链数据分析