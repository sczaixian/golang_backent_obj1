




每个事件的日志都被永久固定在它发生时那个区块的 交易收据 里

错误（Errors）的情况
该交易的收据中会有一个状态字段标记为 0（失败），并且不会包含任何事件日志（因为执行回滚了）。
错误信息（如 require 中的提示字符串）通常不存储在链上。
要获取错误信息，客户端需要本地模拟交易执行或通过节点的调试接口来获取回退原因。




EVM 的设计首先服务于以太坊的核心愿景：创建一个去中心化的、可编程的区块链平台
确定性：给定相同的初始状态和交易输入，任何节点上的EVM执行结果必须完全一致。这是达成共识的基础。
安全性与隔离性：智能合约代码必须在一个沙盒环境中运行。一个合约的错误或恶意行为不能导致整个网络崩溃或影响其他合约的执行。
可终止性：必须有一种机制防止代码陷入无限循环，耗尽所有网络资源。gas
平台无关性：EVM字节码可以在任何操作系统（Windows, Linux, macOS）上的任何以太坊节点中运行，只要该节点遵循协议规则。



以太坊虚拟机（EVM）选择基于栈（Stack-based）的设计，而不是基于寄存器（Register-based）（如JVM、Lua VM），
    是经过深思熟虑的，主要基于简单性、确定性和指令集紧凑性的考量。
简单性与实现简便、编译器更容易设计、解释器更容易实现
指令集紧凑，字节码更小、部署成本更低、传输效率更高
确定性与一致性，执行路径唯一


数据加载基于栈
临时变量存memory，一个易失性的、可线性寻址的字节数组。合约每次调用时，内存都是空的。
智能合约的 Storage（存储） 的数据最终被持久化到了全球所有以太坊全节点和归档节点的本地数据库里





归档节点（Archive Node）	完整的一切：所有历史区块+所有历史状态	极高（数TB级别）	服务于区块链浏览器、高级数据分析工具等需要查询历史状态的特定应用。
全节点（Full Node）	最新区块历史 + 最新世界状态	高（数百GB级别）	网络的中坚力量，独立验证所有交易和区块，但不存储所有历史状态。
轻节点（Light Node）	仅区块头（Block Headers）	极低（数GB级别）	快速同步，依赖全节点提供数据验证，适合手机钱包等轻量级应用。


无状态客户端（Stateless Clients）：这是更受青睐的未来方案。在这个模型中，验证者（矿工/验证节点）不再需要存储完整的状态。相反，他们通过一种叫做 “见证数据（Witness）” 的加密证明来验证交易。

工作原理：当您发送一笔交易时，您需要附带一个小的“见证”（如Merkle证明），来证明您的账户状态和交易所需的合约状态。验证者只需要验证这个证明是否正确即可，而无需自己存储所有这些状态。

效果：这将极大地降低运行一个全节点的资源要求，从根本上解决状态爆炸问题。Vitalik Buterin（以太坊创始人）多次强调这是以太坊未来的关键升级之一。




在 Solidity 中，变量可以存储在以下位置：
Storage（存储）：持久化存储在区块链上，是合约状态的一部分，消耗大量 Gas。
Memory（内存）：临时存储，仅在函数执行期间存在，Gas 成本低。
Calldata：不可修改的临时存储，用于存储函数调用的原始输入数据。
Stack：用于存储小型的局部变量，由 EVM 直接管理。
代码/常量：在合约部署时确定并嵌入字节码，不占用存储空间。



// 状态变量 - 全部存储在 Storage 中
contract FundMe{
    address public owner; // -> Storage (槽位 0)
    mapping(address => uint256) public fundersToAmount; // -> Storage (整个映射分散在多个槽位，通过 keccak 哈希计算地址)

    uint256 constant MINIMUM_VALUE = 1 * 10 ** 17; // -> 常量，直接嵌入字节码，不占用 Storage
    uint256 constant TARGET = 1000 * 10 ** 18; // -> 常量，直接嵌入字节码，不占用 Storage
    uint256 deploymentTimestamp; // -> Storage (槽位 1)
    uint256 lockTime; // -> Storage (槽位 2)
    address erc20Addr; // -> Storage (槽位 3)
    bool public getFundSuccess = false;  // -> Storage (槽位 4)

    AggregatorV3Interface public dataFeed; // -> Storage (槽位 5)

    // 事件定义不占用存储空间
    event FundWithdrawByOwner(uint256);
    event RefundByFunder(address, uint256);

    // 修饰器定义不占用存储空间
    modifier onlyOwner(){
        // msg.sender -> 特殊全局变量，来自交易上下文
        require(msg.sender == owner, "this function can only called by owner");
        _;
    }

    modifier windowClosed(){
        // block.timestamp -> 特殊全局变量，来自区块上下文
        require(block.timestamp >= deploymentTimestamp + lockTime, "window is not closed");
        _;
    }

    constructor (uint256 _lockTime, address dataFeedAddr){
        // _lockTime, dataFeedAddr -> 函数参数，存储在 Calldata 或 Stack
        dataFeed = AggregatorV3Interface(dataFeedAddr); // -> 赋值给 Storage 变量
        owner = msg.sender; // -> 赋值给 Storage 变量 (msg.sender 来自交易上下文)
        deploymentTimestamp = block.timestamp; // -> 赋值给 Storage 变量 (block.timestamp 来自区块上下文)
        lockTime = _lockTime; // -> 赋值给 Storage 变量
    }

    function fund() external payable {
        // msg.value -> 特殊全局变量，来自交易上下文
        // convertEthToUsd(msg.value) -> 函数调用，参数和返回值在 Memory 中处理
        require(convertEthToUsd(msg.value) >= MINIMUM_VALUE, "send more eth");
        require(block.timestamp <= deploymentTimestamp + lockTime, "window is closed");
        fundersToAmount[msg.sender] += convertEthToUsd(msg.value); // -> 访问并修改 Storage 映射
    }

    function getFund() external payable onlyOwner windowClosed {
        // address(this).balance -> 特殊全局变量，访问当前合约的余额
        require(convertEthToUsd(address(this).balance) >= TARGET, "target is not reached");

        // 局部变量 -> 存储在 Stack 或 Memory
        bool success;
        uint256 balance = address(this).balance; // -> balance 存储在 Stack/Memory

        // 低级调用，参数和返回值在 Memory 中处理
        (success, ) = payable(msg.sender).call{value: address(this).balance}("");
        require(success, "");
        fundersToAmount[msg.sender] = 0; // -> 修改 Storage 映射
        getFundSuccess = true; // -> 修改 Storage 变量
        emit FundWithdrawByOwner(balance); // -> 事件参数在 Memory 中
    }

    function refund() external windowClosed {
        require(convertEthToUsd(address(this).balance) < TARGET, "target is reached");
        require(fundersToAmount[msg.sender] != 0, "there is no fund for you");
        // 局部变量 -> 存储在 Stack 或 Memory
        bool success;
        uint256 balance = fundersToAmount[msg.sender]; // -> 从 Storage 读取到 Stack/Memory
        (success,) = payable(msg.sender).call{value: fundersToAmount[msg.sender]}(""); // -> 访问 Storage 映射
        require(success, "transfer tx failed");
        fundersToAmount[msg.sender] = 0; // -> 修改 Storage 映射
        emit RefundByFunder(msg.sender, balance); // -> 事件参数在 Memory 中
    }

    function setFuntToAmount(address funder, uint256 amountToUpdate) external {
        // funder, amountToUpdate -> 函数参数，存储在 Calldata
        require(msg.sender == erc20Addr, "you do not have permission to call this function");
        fundersToAmount[funder] = amountToUpdate; // -> 修改 Storage 映射
    }

    function transferOwnership(address newOwner) external windowClosed onlyOwner{
        // newOwner -> 函数参数，存储在 Calldata
        owner = newOwner; // -> 修改 Storage 变量
    }

    function getChainlinkDataFeedLatestAnswer() public view returns (int) {
        // 函数返回值变量 -> 存储在 Memory
        // 解构赋值中的临时变量 -> 存储在 Stack
        (
            /* uint80 roundId */,
            int256 answer,
            /*uint256 startedAt*/,
            /*uint256 updatedAt*/,
            /*uint80 answeredInRound*/
        ) = dataFeed.latestRoundData(); // -> 调用外部合约，返回值在 Memory 中
        return answer; // -> 从 Memory 返回
    }

    function convertEthToUsd(uint256 ethAmount) internal view returns(uint256) {
        // ethAmount -> 函数参数，存储在 Calldata 或 Stack
        uint256 ethPrice = uint256(getChainlinkDataFeedLatestAnswer()); // -> 函数返回值存储在 Stack/Memory
        return ethAmount * ethPrice / (10 ** 8); // -> 计算在 Stack 中进行，结果返回
    }
}