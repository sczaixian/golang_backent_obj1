


// 使用CAS操作，真正无锁，适合高性能场景
package ringbuf

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

var (
	ErrFull  = errors.New("ringbuf: buffer is full")
	ErrEmpty = errors.New("ringbuf: buffer is empty")
)

// RingBuf 无锁环形缓冲区
type RingBuf struct {
	buf    []unsafe.Pointer
	head   uint32
	tail   uint32
	mask   uint32
}

// New 创建指定大小的环形缓冲区，大小必须是2的幂
func New(size uint32) *RingBuf {
	if size&(size-1) != 0 {
		panic("ringbuf: size must be a power of 2")
	}
	return &RingBuf{
		buf:  make([]unsafe.Pointer, size),
		mask: size - 1,
	}
}

// Put 向缓冲区写入数据（使用CAS确保线程安全）
func (r *RingBuf) Put(item interface{}) error {
	for {
		tail := atomic.LoadUint32(&r.tail)
		head := atomic.LoadUint32(&r.head)
		
		// 检查是否已满
		if tail-head >= uint32(len(r.buf)) {
			return ErrFull
		}
		
		// 计算写入位置
		idx := tail & r.mask
		
		// 使用CAS确保原子性写入
		ptr := unsafe.Pointer(&item)
		if atomic.CompareAndSwapPointer(&r.buf[idx], nil, ptr) {
			// 成功写入，更新tail指针
			atomic.StoreUint32(&r.tail, tail+1)
			return nil
		}
		
		// CAS失败，重试
	}
}

// Get 从缓冲区读取数据（使用CAS确保线程安全）
func (r *RingBuf) Get() (interface{}, error) {
	for {
		head := atomic.LoadUint32(&r.head)
		tail := atomic.LoadUint32(&r.tail)
		
		// 检查是否为空
		if head == tail {
			return nil, ErrEmpty
		}
		
		// 计算读取位置
		idx := head & r.mask
		
		// 尝试读取数据
		ptr := atomic.LoadPointer(&r.buf[idx])
		if ptr != nil {
			// 使用CAS清空槽位，确保只有一个消费者能读取
			if atomic.CompareAndSwapPointer(&r.buf[idx], ptr, nil) {
				// 成功读取，更新head指针
				atomic.StoreUint32(&r.head, head+1)
				item := *(*interface{})(ptr)
				return item, nil
			}
		}
		
		// 读取失败，重试
	}
}

// Len 返回缓冲区中元素数量
func (r *RingBuf) Len() uint32 {
	tail := atomic.LoadUint32(&r.tail)
	head := atomic.LoadUint32(&r.head)
	if tail >= head {
		return tail - head
	}
	return uint32(len(r.buf)) - head + tail
}

// Capacity 返回缓冲区容量
func (r *RingBuf) Capacity() uint32 {
	return uint32(len(r.buf))
}

// IsEmpty 检查缓冲区是否为空
func (r *RingBuf) IsEmpty() bool {
	return atomic.LoadUint32(&r.head) == atomic.LoadUint32(&r.tail)
}







/*
	------------------------------  简单版不是严格无锁 --------------------------------------
*/

package ringbuf

import (
	"errors"
	"sync/atomic"
)

// AtomicRingBuf 基于原子操作的线程安全环形缓冲区
type AtomicRingBuf struct {
	buf     []interface{}
	head    uint32
	tail    uint32
	mask    uint32
	mu      uint32 // 简单的自旋锁
}

func NewAtomicRingBuf(size uint32) *AtomicRingBuf {
	if size&(size-1) != 0 {
		panic("ringbuf: size must be a power of 2")
	}
	return &AtomicRingBuf{
		buf:  make([]interface{}, size),
		mask: size - 1,
	}
}

func (r *AtomicRingBuf) Put(item interface{}) error {
	// 简单的自旋锁实现
	for !atomic.CompareAndSwapUint32(&r.mu, 0, 1) {
	}
	
	defer atomic.StoreUint32(&r.mu, 0)
	
	tail := r.tail
	head := r.head
	
	if tail-head >= uint32(len(r.buf)) {
		return ErrFull
	}
	
	r.buf[tail&r.mask] = item
	r.tail = tail + 1
	return nil
}

func (r *AtomicRingBuf) Get() (interface{}, error) {
	for !atomic.CompareAndSwapUint32(&r.mu, 0, 1) {
	}
	
	defer atomic.StoreUint32(&r.mu, 0)
	
	head := r.head
	tail := r.tail
	
	if head == tail {
		return nil, ErrEmpty
	}
	
	item := r.buf[head&r.mask]
	r.buf[head&r.mask] = nil // 防止内存泄漏
	r.head = head + 1
	return item, nil
}