package main

import (
	"fmt"
	"sync"
	"time"
)

type ChannelManager struct {
	dataChan chan string
	wg       sync.WaitGroup
	stopChan chan struct{}
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		dataChan: make(chan string, 10),
		stopChan: make(chan struct{}),
	}
}

// 启动发送者
func (cm *ChannelManager) StartSender(id int) {
	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()
		
		for i := 0; ; i++ {
			msg := fmt.Sprintf("Sender%d-消息%d", id, i)
			
			select {
			case cm.dataChan <- msg:
				fmt.Printf("发送: %s\n", msg)
			case <-cm.stopChan:
				fmt.Printf("发送者%d 退出\n", id)
				return
			}
			
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

// 启动接收者
func (cm *ChannelManager) StartReceiver(id int) {
	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()
		
		for {
			select {
			case msg, ok := <-cm.dataChan:
				if !ok {
					fmt.Printf("接收者%d 退出\n", id)
					return
				}
				fmt.Printf("接收者%d 收到: %s\n", id, msg)
			case <-cm.stopChan:
				fmt.Printf("接收者%d 退出\n", id)
				return
			}
		}
	}()
}

// 优雅关闭
func (cm *ChannelManager) Stop() {
	fmt.Println("开始优雅关闭...")
	
	// 1. 先关闭stopChan，通知所有goroutine停止发送
	close(cm.stopChan)
	
	// 2. 等待所有goroutine退出
	cm.wg.Wait()
	
	// 3. 关闭数据通道
	close(cm.dataChan)
	
	fmt.Println("所有goroutine已退出")
}

func main() {
	manager := NewChannelManager()
	
	// 启动3个发送者
	for i := 1; i <= 3; i++ {
		manager.StartSender(i)
	}
	
	// 启动2个接收者
	for i := 1; i <= 2; i++ {
		manager.StartReceiver(i)
	}
	
	// 运行5秒后关闭
	time.Sleep(5 * time.Second)
	manager.Stop()
	
	fmt.Println("程序退出")
}