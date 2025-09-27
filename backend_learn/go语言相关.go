
// ----------------------- 内存对齐 ---------------------------------
type Inefficient struct {
    a bool    // 1
    // 7字节padding（因为int64需要8字节对齐）
    b int64   // 8
    c int32   // 4
    // 4字节padding（因为结构体需要8字节对齐，所以总大小必须是8的倍数）
}
// 总大小：1+7+8+4+4=24
type Efficient struct {
    b int64   // 8
    c int32   // 4
    a bool    // 1
    // 3字节padding（因为结构体需要8字节对齐）
}
// 总大小：8+4+1+3=16
// 按字段大小降序排列：将较大的字段放在前面，较小的字段放在后面
// 相关字段放在一起：在考虑内存对齐的同时，也要考虑代码可读性



// -----------------------DAO---------------------------------
DAO层位于分层架构的数据访问层，介于业务逻辑层（Service）与数据库之间
它抽象了数据操作细节，使业务逻辑无需关心数据如何存储或检索，仅通过接口调用DAO提供的方法
数据访问层的核心组件，负责封装所有与数据库或其他持久化存储的交互逻辑。
其核心目标是将业务逻辑与底层数据存储细节解耦，提升代码的可维护性、可测试性和灵活性


核心职责
1. 数据操作封装: 实现CRUD（增删查改）操作
2. 数据库连接管理: 负责数据库连接的建立、维护与释放（如通过*sql.DB管理连接池）。
3. 数据模型转换: 将数据库查询结果（如SQL行数据）转换为Go结构体（Model），供业务层使用。
        例如将SELECT结果映射到User结构体


核心优势与价值

1. 解耦业务与数据逻辑
    业务层（Service）仅依赖DAO接口而非具体数据库实现，
    更换数据库（如MySQL→PostgreSQL）时只需调整DAO实现，无需修改业务代码。

2. 提升可测试性
    通过接口抽象，可轻松创建Mock DAO（如内存模拟数据库），实现业务逻辑的单元测试，无需真实数据库。

3. 代码复用与维护
    数据访问逻辑集中管理，避免SQL语句分散在业务代码中，降低重复并简化维护。


// Service层调用DAO
func (s *UserService) CreateUser(user *model.User) error    

// 依赖注入
// 3. Service 层定义
type UserService struct {
    userDAO UserDAO // 依赖接口
}    
// 单元测试时可注入 Mock DAO，无需连接真实数据库
s := &UserService{userDAO:userDAO}  // 数据库改变的时候只需要换注入的 数据库 DAO 就行了
s.getUser(UserID)



----------------------- 闭包 ---------------------------------
闭包（Closure）是 函数与其引用环境的组合体，允许内部函数捕获并持久化外部作用域的变量。
即使外部函数执行结束，闭包仍能访问和修改这些变量，实现状态的跨调用保存

闭包通过引用捕获外部变量，而非值拷贝

生命周期延长：闭包引用的变量会持续存在，直到闭包不再被引用

引用捕获：Go 闭包通过指针持有外部变量的引用。若多个闭包共享同一变量，修改会相互影响。
循环陷阱：在循环中直接使用闭包可能导致意外结果（所有闭包共享循环变量的最终值）
for i := 0; i < 3; i++ {
    i := i  // 创建副本  <------- 加上这一行解决
    go func() { fmt.Println(i) }() // 输出 0,1,2
}
应用场景：
状态封装，延迟计算，资源管理，装饰器，异步回调

多个 goroutine 并发修改闭包变量，使用互斥锁（sync.Mutex）同步访问




----------------------- context ---------------------------------
优雅地通知所有正在进行的 Goroutine（协程）：“任务取消了，请立刻停止手上的工作，释放资源并退出”。

Context（上下文）就是为了解决这个问题而生的。
它的主要目的就是在 API 和 Goroutine 之间传递 
    deadlines（截止时间）、cancellation signals（取消信号）以及其他请求范围的值。

type Context interface {
    // 返回此 Context 被取消的截止时间。如果没有设置截止时间，则 ok 为 false。
    Deadline() (deadline time.Time, ok bool)

    // 返回一个 Channel。当 Context 被取消或超时后，这个 Channel 会被关闭。
    // 可以通过监听这个 Channel 来收到取消信号。
    Done() <-chan struct{}

    // 返回 Context 被取消的原因。在 Done() 的 Channel 关闭前调用会阻塞。
    Err() error

    // 用于从 Context 中获取关联的键值对数据。
    Value(key any) any
}


context.Background(): 通常用作根节点，是所有派生 Context 的源头。
    它永远不会被取消，没有超时时间，也不携带值。
    通常在 main 函数、初始化函数或测试中使用。

context.TODO(): 功能与 Background 完全相同。
    通常在你不确定该使用哪个 Context，或者暂时还没有可用的 Context 传入时，作为一个占位符使用。
    静态分析工具可以使用它来提醒你后续需要传入一个真正的 Context。

简单来说：一切 Context 的源头都是 Background() 或 TODO()。    


通过“派生”，我们可以给父 Context 附加额外的控制能力。主要有四种函数

1. WithCancel：创建一个可以手动取消的 Context。（需要主动发出取消信号）
ctx, cancel := context.WithCancel(parentCtx)
defer cancel() // 通常用 defer 保证在任何情况下都能取消，释放资源
// 调用 cancel() 函数时，ctx.Done() 返回的 Channel 会被关闭。    


2. WithTimeout：创建一个会超时自动取消的 Context。（网络请求，数据库操作）
// 1 秒后超时自动取消
ctx, cancel := context.WithTimeout(parentCtx, 1*time.Second)
defer cancel() // 即使提前完成，也调用 cancel 释放资源

3. WithDeadline：创建一个在指定时间点自动取消的 Context
// 在 2023-10-01 12:00:00 自动取消
d := time.Date(2023, time.October, 1, 12, 0, 0, 0, time.UTC)
ctx, cancel := context.WithDeadline(parentCtx, d)
defer cancel()

4. WithValue：创建一个可以携带键值对数据的 Context（在请求链路中传递 traceId、用户认证信息等）
// 用于在流程间传递请求范围的数据，如 traceId、用户认证token等。
ctx := context.WithValue(parentCtx, key, "value")
v := ctx.Value(key)

重要规则：当一个父 Context 被取消时，所有由它派生的子 Context 也会被自动取消。这是一个级联效应。


最佳实践和注意事项
    第一个参数显式传递
    不要存储 Context：不要将一个 Context 保存在一个结构体（struct）中。应该显式地传递它
    Context 是线程安全的：你可以安全地在多个 Goroutine 中同时使用同一个 Context
    总是调用 cancel()：只要调用了 WithCancel, WithTimeout, WithDeadline，
        就必须在函数退出前调用返回的 cancel 函数（通常用 defer cancel()），
        以确保及时释放相关资源。即使操作提前完成也要调用，这是一种良好的习惯。

    谨慎使用 WithValue：Context 的值传递机制应该仅用于传递请求范围的进程和 
        API 边界的数据（如请求ID、用户令牌等），而不应被用作函数的可选参数来传递。
        键的类型最好是自定义的，避免冲突



----------------------- goroutine ---------------------------------
用户态线程； 创建、切换开销小、内存占用低（kb级别）
使用 go 关键字创建
使用 channel 通信（有无缓冲区）
有共享内存时使用 互斥锁 保护临界区
无法从外部强制杀死一个 Goroutine，只能等它自己结束。通常通过 context.Context 来传递取消信号。
避免在 Goroutine 中使用外部的循环变量
---  deepseek 搜索 ------

几乎所有的代码都在 goroutine 上运行
对于你编写的业务逻辑和并发任务，它们 100% 都运行在 Goroutine 中

调度器运行在线程上
实现调度器、内存分配器、垃圾回收器等核心功能的“元代码直接运行在操作系统线程上



select语法类似于 switch，但每个 case必须是通道操作（发送或接收）
若多个 case同时就绪，随机执行一个（避免优先级问题）
无 default分支时，select会阻塞直到至少一个 case就绪
有 default时，无就绪操作则立即执行 default（实现非阻塞）
空 select：无任何 case时（如 select{}），会永久阻塞，常用于防止主协程退出

for {
    select {
    case <-s.ctx.Done(): // 检查服务是否被通知关闭
        ... // 优雅退出
        return
    case <-time.After(3 * time.Second):  防止因通道无响应导致程序卡死，适用于网络请求或任务执行
        // fmt.Println("超时！")
    default:  <----  实现非阻塞
    }
    ... // 主要工作逻辑在这里
}


优先级处理：嵌套 select实现高优先级通道优先处理
select {
case urgentMsg := <-urgentChan: // 紧急事件
    handleUrgent(urgentMsg)
default:
    select { // 无紧急事件时处理普通事件
    case normalMsg := <-normalChan:
        handleNormal(normalMsg)
    }
}


多通道监听
超时控制
非阻塞操作
任务取消与优雅退出
循环监听与事件驱动
嵌套select实现优先级处理



time.After的问题与根源
内存泄漏机制
每次调用 time.After(d)会创建新的 Timer，底层通过 NewTimer(d).C实现。
若循环频率高（如每秒百万次），而超时时间较长（如 3 分钟），
    大量未触发的 Timer会累积内存（每个 Timer约 200 字节），直至超时后才会释放。
后果：内存持续增长，可能引发 OOM（Out of Memory）

优化方案：单次创建，周期性触发，无额外内存分配，适用于心跳检测、状态上报等场景
ticker := time.NewTicker(1 * time.Second) // 创建 Ticker，周期 1 秒
defer ticker.Stop()                      // 必须调用 Stop() 释放资源

for {
    select {
    case <-ticker.C:      // 周期性触发
        doTask()
    case <-stopChan:      // 退出信号
        return
    }
}














func (om *OrderManager) orderExpiryProcess() {
	// 1. 使用 defer recover 来捕获可能的 panic,防止主协程死掉
	defer func() {
		if r := recover(); r != nil {
			xzap.WithContext(om.Ctx).Error("[Order Manage] dq process recovered: " + fmt.Sprintf("%v", r))
		}
	}()
    for condition {
        // ....
    }
}

如果在函数进入循环 以前 发生 panic， recover函数会捕获，然后xzap写日志，orderExpiryProcess函数 返回 
如果在函数进入循环 以后 发生 panic， 会立即进入下一次循环（因为有个无限for循环），但是 发生点以后得代码不会执行








channel  用在协程间通信，如果有缓冲可以看成2条方向的管道，如果没有缓冲等于协程间同步，
这样避免了共享内存的一系列复杂加锁设计，带缓冲区可以实现异步通信
多线程并发避免竞争锁开销，保证安全，较少锁设计的产生的性能问题，降低了程序设计的难度


一个 channle 包含一个 ring buffer 和 两个阻塞队列

type hchan struct {
    /*--------- ring buffer ------------*/
    qcount   uint           // 队列中的元素数量
    dataqsiz uint           // 缓冲区大小   无缓冲区 0
    buf      unsafe.Pointer // 指向底层循环数组的指针（有缓冲区的channel）  无缓冲区 nil
    elemsize uint16         // 元素的大小（类型宽度）
    elemtype *_type         // chan 中的元素类型
    sendx    uint           // 已发送元素环形数组中的下标索引
    recvx    uint           // 已接收元素环形数组中的下标索引
    /*--------- queue ------------*/
    // 无法立即完成的放等待队列
    recvq    waitq  // 等待接收的协程队列
    sendq    waitq  // 等待发送的协程队列

    closed   uint32   //channel是否关闭的标志
    lock mutex //互斥锁
}

ch := make(chan int, 1)
ch <- 1           // 缓冲区满
go func() { ch <- 2 }() // 发送者阻塞，加入 sendq
<-ch              // 接收者取出 1，并唤醒 sendq 中的发送者写入 2


panic 出现的场景还有：
关闭值为 nil 的 channel
关闭已经关闭的 channel
向已经关闭的 channel 中写数据





5. Channel有哪些常见的使用场景   ------  deepseek
任务分发和处理：可以通过Channel将任务分发给多个goroutine进行处理，并将处理结果发送回主goroutine进行汇总和处理。
并发控制：可以通过Channel来进行信号量控制，限制并发的数量，避免资源竞争和死锁等问题。
数据流处理：可以通过Channel实现数据流的处理，将数据按照一定的规则传递给不同的goroutine进行处理，提高并发处理效率。
事件通知和处理：可以通过Channel来实现事件的通知和处理，将事件发送到Channel中，让订阅了该Channel的goroutine进行相应的处理。
异步处理：可以通过Channel实现异步的处理，将任务交给其他goroutine处理，自己继续执行其他任务，等待处理结果时再从Channel中获取。




文章：https://juejin.cn/post/7541297378125414441



go 与其他语言进行对比，？





//------------------------ json -------------------------
json.Unmarshal([]byte(result), &listing) 
[]byte(result) 将字符串 result 转换为字节切片
 Unmarshal 仅接受字节类型输入














//------------------------ time -------------------------
time.Now().Unix() 是 Go 语言中用于获取 当前时间的 Unix 时间戳（秒级） 的标准方法
time.Sleep(1 * time.Second)       <----------- 阻塞当前协程，不影响其他任务
time.Nanosecond（纳秒）
time.Microsecond（微秒）
time.Millisecond（毫秒）
time.Second（秒）
time.Minute（分钟）
time.Hour（小时）






















----------------------- 切片 ---------------------------------
slice := make([]int, 5)           // 长度和容量都是 5
slice := make([]int, 3, 5)                 // 三个元素 容量5
slice := []string{"a", "b", "c"}  // 长度和容量都是3
slice := []int{10,20,30,40}       // 长度和容量都是4
slice := []string{99:""}          // 使用空字符串初始化第100个元素

