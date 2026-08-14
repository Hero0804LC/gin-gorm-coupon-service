package utils

import (
	"fmt"
	"sync"
	"time"
)

const (
	defaultEpoch  = int64(1640995200000) // 2022-01-01 00:00:00 UTC（毫秒）
	defaultNodeID = int64(1)
)

type Snowflake struct {
	mu       sync.Mutex
	epoch    int64
	nodeID   int64
	sequence int64
	maxSeq   int64
	lastTs   int64
}

var (
	sf   *Snowflake
	once sync.Once
)

// InitSnowflake 手动初始化（可选，不调用也会自动用默认值）
func InitSnowflake(nodeID int64) {
	once.Do(func() {
		sf = &Snowflake{
			epoch:  defaultEpoch,
			nodeID: nodeID,
			maxSeq: 4095, // 2^12 - 1
			lastTs: -1,
		}
	})
}

// GenerateOrderNo 生成订单号（雪花算法，带时间回拨保护，不会死锁）
func GenerateOrderNo() string {
	once.Do(func() {
		if sf == nil {
			sf = &Snowflake{
				epoch:  defaultEpoch,
				nodeID: defaultNodeID,
				maxSeq: 4095,
				lastTs: -1,
			}
		}
	})

	sf.mu.Lock()
	defer sf.mu.Unlock()

	now := time.Now().UnixMilli()

	// 时间回拨检测：如果当前时间小于上次生成时间，直接 panic 暴露问题
	if now < sf.lastTs {
		sf.mu.Unlock()
		panic(fmt.Sprintf("snowflake: clock moved backwards, lastTs=%d now=%d", sf.lastTs, now))
	}

	if now == sf.lastTs {
		// 同一毫秒内，序列号自增
		sf.sequence = (sf.sequence + 1) & sf.maxSeq

		if sf.sequence == 0 {
			// 序列号用尽，等到下一毫秒
			deadline := time.After(10 * time.Millisecond)
			ticker := time.NewTicker(50 * time.Microsecond)
			defer ticker.Stop()

			for {
				select {
				case <-deadline:
					// 10ms 还没到下一毫秒，说明时钟有问题，直接 panic
					sf.mu.Unlock()
					panic("snowflake: wait for next millisecond timeout, clock may be stalled")
				case <-ticker.C:
					now = time.Now().UnixMilli()
					if now > sf.lastTs {
						goto next
					}
				}
			}
		}
	} else {
		// 新的毫秒，序列号重置
		sf.sequence = 0
	}

next:
	sf.lastTs = now

	ts := uint64(now - sf.epoch)
	node := uint64(sf.nodeID)
	seq := uint64(sf.sequence)

	// 41 bits timestamp | 10 bits node | 12 bits sequence
	id := (ts << 22) | (node << 12) | seq

	return fmt.Sprintf("%d", id)
}
