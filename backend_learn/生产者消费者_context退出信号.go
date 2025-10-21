package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Processor struct {
	ch     chan int
	wg     sync.WaitGroup
	quit   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
}

func NewProcessor() *Processor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Processor{
		ch:     make(chan int, 10),
		quit:   make(chan struct{}),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Processor) Start() {
	p.wg.Add(2)
	go p.producer()
	go p.consumer()
}

func (p *Processor) producer() {
	defer p.wg.Done()

	for i := 0; i < 100; i++ {
		select {
		case <-p.ctx.Done():
			fmt.Println("producer: 收到停止信号，退出")
			return
		case p.ch <- i:
			fmt.Println("producer: 生产数据:", i)
		}
	}
	close(p.ch) // 生产完成
}

func (p *Processor) consumer() {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			fmt.Println("consumer: 收到停止信号，退出")
			return // 直接退出，不消费剩余数据
		case v, ok := <-p.ch:
			if !ok {
				fmt.Println("consumer: channel 已关闭，退出")
				return
			}
			fmt.Println("consumer: 消费数据:", v)
		}
	}
}

func (p *Processor) Stop() {
	p.cancel()    // 通知停止
	p.wg.Wait()   // 等待完成
	close(p.quit) // 完全退出
}

func main() {
	processor := NewProcessor()
	processor.Start()

	// 模拟运行
	time.Sleep(20)

	// 优雅退出
	processor.Stop()
	fmt.Println("程序优雅退出")
}
