


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
	TimeWheel [WheelSize]wheel   // 多个个轮子  TODO:看看是怎么用的
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