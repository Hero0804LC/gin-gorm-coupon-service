package async

// TaskChan 全局任务通道（buffer 10000）
var TaskChan = make(chan *SeckillTask, 10000)

// Dispatch 非阻塞投递
func Dispatch(task *SeckillTask) bool {
	select {
	case TaskChan <- task:
		return true
	default:
		return false
	}
}
