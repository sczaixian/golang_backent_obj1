package syncpart



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
