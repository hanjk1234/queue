package queue

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnqueueDequeue(t *testing.T) {
	q := NewSliceQueue(1)
	v := q.Get().(int)
	assert.Equal(t, 1, v)
	q.Enqueue(2)
	assert.Equal(t, 1, v)
	// will replace 2
	q.Enqueue(3)

	result := q.Dequeue()
	if result != 1 {
		t.Errorf("Expected dequeue result to be 1, got %v", result)
	}

	result = q.Dequeue()
	if result != 3 {
		t.Errorf("Expected dequeue result to be 3 got %v", result)
	}

	// if queue empty，Dequeue is nil
	result = q.Dequeue()
	if result != 3 {
		t.Errorf("Expected dequeue result to be nil, got %v", result)
	}
}

func TestGet(t *testing.T) {
	q := NewSliceQueue(1)

	q.Enqueue(2)

	result := q.Get()
	if result != 1 {
		t.Errorf("Expected Get(0) result to be 1, got %v", result)
	}

}

// 测试并发 Get 操作
func TestConcurrentGet(t *testing.T) {
	q := NewSliceQueue(1)
	maxGoroutines := 20
	guard := make(chan struct{}, maxGoroutines)
	// 使用 WaitGroup 等待所有 Get 操作完成
	var wg sync.WaitGroup
	num := 1000
	results := make([]interface{}, num)
	wg.Add(num)
	go func() {
		for i := 0; i < num; i++ {
			guard <- struct{}{} // 当guard channel已满时会阻塞
			go func(index int) {
				defer wg.Done()
				results[index] = q.Get()
				time.Sleep(20 * time.Millisecond)
				<-guard
			}(i)
		}
	}()
	//time.Sleep(20 * time.Millisecond)
	// mock insert
	wg.Add(10)
	go func() {
		time.Sleep(1 * time.Millisecond)
		for i := 1; i < 11; i++ {
			signalChannel := make(chan interface{})
			go func(idx int) {
				q.Enqueue(idx)
				close(signalChannel)
			}(i)

			time.Sleep(80 * time.Millisecond)
			go func() {
				select {
				case _, ok := <-signalChannel:
					if !ok {
						v := q.Dequeue()
						fmt.Println(v)
						wg.Done()
					}
				}
			}()
		}
	}()

	wg.Wait()
	if results[0] != 1 {
		t.Errorf("Expected concurrent Get(0) result to be 1, got %v", results[0])
	}

	for i := 1; i < num; i++ {
		// 验证 Get 的结果
		if results[i].(int) < 1 {
			t.Errorf("Expected concurrent Get(%d) result to be 2 or 1, got %v", i, results[i])
		} else {
			//fmt.Printf("concurrent Get(%d)  got %v\n", i, results[i])
		}
	}

}
