/*
 * Author: fasion
 * Created time: 2023-04-18 14:02:09
 * Last Modified by: fasion
 * Last Modified time: 2023-04-20 17:04:47
 */

package jobutils

import (
	"context"
	"sync"
	"time"

	"github.com/fasionchan/goutils/stl"
)

type SmartHandlerPtr = *SmartHandler

type SmartHandler struct {
	callback func()
	ident    string
	async    bool            // 异步执行
	ctx      context.Context // 上下文对象，等待超时
	tickets  JobTokens       // 控制并发数
	delay    time.Duration   // 延迟执行
	interval time.Duration   // 控制执行间隔
	deadline time.Duration   // 控制截止时间
	merge    bool            // 合并执行

	lastCallingTime time.Time
	lastCalledTime  time.Time
}

func (handler *SmartHandler) Dup() *SmartHandler {
	return stl.Dup(handler)
}

func (handler *SmartHandler) WithIdent(ident string) *SmartHandler {
	handler.ident = ident
	return handler
}

func (handler *SmartHandler) WithAsync(async bool) *SmartHandler {
	handler.async = async
	return handler
}

func (handler *SmartHandler) WithCtx(ctx context.Context) *SmartHandler {
	handler.ctx = ctx
	return handler
}

func (handler *SmartHandler) WithConcurrentcy(n int) *SmartHandler {
	handler.tickets = NewJobTokens(n)
	return handler
}

func (handler *SmartHandler) WithDelay(delay time.Duration) *SmartHandler {
	handler.delay = delay
	return handler
}

func (handler *SmartHandler) WithInterval(interval time.Duration) *SmartHandler {
	handler.interval = interval
	return handler
}

func (handler *SmartHandler) WithDeadline(deadline time.Duration) *SmartHandler {
	handler.deadline = deadline
	return handler
}

func (handler *SmartHandler) WithMerge(merge bool) *SmartHandler {
	handler.merge = merge
	return handler
}

func (handler *SmartHandler) Call() {
	// 同步执行
	if !handler.async {
		handler.call(time.Now())
		return
	}

	go handler.callAsync()
}

func (handler *SmartHandler) call(callingTime time.Time) {
	if callingTime.IsZero() {
		callingTime = time.Now()
	}

	// 逾期
	if handler.overdue(callingTime) {
		return
	}

	// 等待：延迟执行、执行间隔
	if ok := handler.doWait(callingTime); !ok {
		return
	}

	// 合并
	if handler.merged(callingTime) {
		return
	}

	handler.lastCallingTime = time.Now()
	handler.callback()
	handler.lastCalledTime = time.Now()
}

func (handler *SmartHandler) callAsync() {
	// 准备调用的时间
	callingTime := time.Now()

	// 获取令牌（用于控制并发数）
	for !handler.tickets.Acquire(handler.ctx, time.Second) {
		// 调用取消：逾期、合并
		if handler.callCanceled(callingTime) {
			return
		}
	}

	// 释放令牌
	defer handler.tickets.Release()

	handler.call(callingTime)
}

func (handler *SmartHandler) overdue(callingTime time.Time) bool {
	if handler.deadline <= 0 {
		return false
	}

	return time.Now().Sub(callingTime) > handler.deadline
}

func (handler *SmartHandler) merged(callingTime time.Time) bool {
	if !handler.merge {
		return false
	}

	return handler.lastCallingTime.After(callingTime)
}

func (handler *SmartHandler) callCanceled(callingTime time.Time) bool {
	return handler.overdue(callingTime) || handler.merged(callingTime)
}

func (handler *SmartHandler) doWait(callingTime time.Time) bool {
	now := time.Now()

	var waitDuration time.Duration

	// 延迟调用
	if left := handler.delay - now.Sub(callingTime); left > waitDuration {
		waitDuration = left
	}

	// 确保间隔
	if left := handler.interval - now.Sub(handler.lastCalledTime); left > waitDuration {
		waitDuration = left
	}

	if waitDuration <= 0 {
		return true
	}

	select {
	case <-time.After(waitDuration):
		return true
	case <-handler.ctx.Done():
		return false
	}
}

type SmartHandlers []*SmartHandler

func (handlers SmartHandlers) Append(more ...*SmartHandler) SmartHandlers {
	return append(handlers, more...)
}

func (handlers SmartHandlers) Purge(f func(*SmartHandler) bool) SmartHandlers {
	return stl.Purge(handlers, f)
}

func (handlers SmartHandlers) PurgeByIdent(ident string) SmartHandlers {
	return handlers.Purge(func(handler *SmartHandler) bool {
		return handler.ident == ident
	})
}

func (handlers SmartHandlers) Call() {
	stl.ForEach(handlers, SmartHandlerPtr.Call)
}

type SmartHandlersMappingByString map[string]SmartHandlers

type TopicBroker struct {
	mutex    sync.RWMutex
	handlers SmartHandlersMappingByString
}

func (broker *TopicBroker) Publish(topic string) {
	broker.mutex.RLock()
	defer broker.mutex.RUnlock()

	broker.handlers[topic].Call()
}

func (broker *TopicBroker) Subscribe(topic string, callback func()) {

}

func (broker *TopicBroker) subscribe(topic string, handler *SmartHandler) *TopicBroker {
	broker.mutex.Lock()
	broker.mutex.Unlock()

	broker.handlers[topic] = broker.handlers[topic].Append(handler)

	return broker
}

func (broker *TopicBroker) UnsubscribeByIdent(topic string, ident string) {
	broker.mutex.Lock()
	broker.mutex.Unlock()

	broker.handlers[topic] = broker.handlers[topic].PurgeByIdent(ident)
}

func (broker *TopicBroker) Unsubscribe(topic string) {
	broker.mutex.Lock()
	broker.mutex.Unlock()

	if _, ok := broker.handlers[topic]; ok {
		delete(broker.handlers, topic)
	}
}

type TopicBrokerSubscribing struct {
	SmartHandler

	broker *TopicBroker
	topic  string
}

func (subscribing *TopicBrokerSubscribing) Done() *TopicBroker {
	return subscribing.broker.subscribe(subscribing.topic, &subscribing.SmartHandler)
}
