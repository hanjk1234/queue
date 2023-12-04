package queue

import (
	"sync/atomic"
)

type SliceQueue struct {
	data [2]interface{}
	head uint32
	tail uint32
}

func NewSliceQueue(v interface{}) (q *SliceQueue) {
	q = &SliceQueue{head: 0, tail: 1}
	q.data[q.head] = v

	return q
}

func (q *SliceQueue) Enqueue(v interface{}) {
	t := atomic.LoadUint32(&q.tail)
	q.data[t] = v
}

// Dequeue 弹出队头元素
func (q *SliceQueue) Dequeue() interface{} {
	h := atomic.LoadUint32(&q.head)
	t := atomic.LoadUint32(&q.tail)

	v := q.data[h]
	if q.data[t] == nil {
		return v
	}

	atomic.StoreUint32(&q.head, t)
	atomic.StoreUint32(&q.tail, h)
	q.data[h] = nil
	return v

}

// Get 获取队列中的第一个元素，不弹出
func (q *SliceQueue) Get() interface{} {
	h := atomic.LoadUint32(&q.head)
	v := q.data[h]
	return v
}
