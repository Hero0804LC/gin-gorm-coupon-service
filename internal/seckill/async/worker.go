package async

import (
	"context"
	"fmt"
)

// TaskHandler 任务处理函数类型（依赖倒置，不依赖 service 包）
type TaskHandler func(ctx context.Context, task *SeckillTask) error

type Worker struct {
	taskChan <-chan *SeckillTask
	handler  TaskHandler
}

// NewWorker 只需要 Channel + 处理函数
func NewWorker(taskChan <-chan *SeckillTask, handler TaskHandler) *Worker {
	return &Worker{
		taskChan: taskChan,
		handler:  handler,
	}
}

func (w *Worker) Start(ctx context.Context) {
	go func() {
		fmt.Println("[SeckillWorker] started, waiting for tasks...")
		for {
			select {
			case <-ctx.Done():
				fmt.Println("[SeckillWorker] stopped")
				return
			case task, ok := <-w.taskChan:
				if !ok {
					fmt.Println("[SeckillWorker] channel closed, exiting")
					return
				}
				w.process(task)
			}
		}
	}()
}

func (w *Worker) process(task *SeckillTask) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[SeckillWorker] panic: %v, user=%d coupon=%d\n", r, task.UserID, task.CouponID)
		}
	}()

	fmt.Printf("[SeckillWorker] processing user=%d coupon=%d\n", task.UserID, task.CouponID)

	if err := w.handler(context.Background(), task); err != nil {
		fmt.Printf("[SeckillWorker] FAILED user=%d coupon=%d err=%v\n", task.UserID, task.CouponID, err)
		return
	}

	fmt.Printf("[SeckillWorker] SUCCESS user=%d coupon=%d\n", task.UserID, task.CouponID)
}
