


// 信号量
func rateLimitWithSemaphore() {
    // 限制同时只能有5个进行中的生产操作
    semaphore := make(chan struct{}, 5)
    dataCh := make(chan string, 100)
    
    // 生产者
    for i := 0; i < 20; i++ {
        go func(id int) {
            semaphore <- struct{}{} // 获取信号量
            
            // 生产数据
            data := fmt.Sprintf("data-%d", id)
            dataCh <- data
            fmt.Printf("Producer %d produced\n", id)
            
            <-semaphore // 释放信号量
        }(i)
    }
}