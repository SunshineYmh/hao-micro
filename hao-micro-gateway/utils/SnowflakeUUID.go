package utils

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// Snowflake 结构体
type Snowflake struct {
	mu          sync.Mutex
	startTime   int64 // 起始时间戳（毫秒）
	machineID   int64 // 机器ID
	sequence    int64 // 序列号
	lastGenTime int64 // 上次生成ID的时间戳（毫秒）
}

var snowflake *Snowflake

func IntoSnowflake() error {
	sf, err := NewSnowflake(1) // 传入机器ID
	if err != nil {
		fmt.Println("Failed to create Snowflake:", err)
		return err
	}
	snowflake = sf
	return nil
}

func GetUUID() string {
	uuid := snowflake.Generate()
	return strconv.FormatInt(uuid, 10)
}

// NewSnowflake 创建一个Snowflake实例
func NewSnowflake(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID >= (1<<10) {
		return nil, errors.New("machine ID must be between 0 and 1023")
	}
	return &Snowflake{
		startTime:   1635000000000, // 设置起始时间戳，适当调整为合适的值
		machineID:   machineID,
		sequence:    0,
		lastGenTime: -1,
	}, nil
}

// Generate 生成一个新的雪花ID
func (sf *Snowflake) Generate() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	currentTime := time.Now().UnixNano() / 1e6 // 当前时间戳（毫秒）

	if currentTime < sf.lastGenTime {
		panic("Clock moved backwards. Refusing to generate ID.")
	}

	if currentTime == sf.lastGenTime {
		sf.sequence = (sf.sequence + 1) & 4095
		if sf.sequence == 0 {
			for currentTime <= sf.lastGenTime {
				currentTime = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		sf.sequence = 0
	}

	sf.lastGenTime = currentTime

	id := (currentTime-sf.startTime)<<22 | (sf.machineID << 12) | sf.sequence
	return id
}

func main2() {
	sf, err := NewSnowflake(1) // 传入机器ID
	if err != nil {
		fmt.Println("Failed to create Snowflake:", err)
		return
	}

	var wg sync.WaitGroup
	numWorkers := 10 // 线程数量

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()

			// 每个线程生成10个雪花ID并打印
			for j := 0; j < 10; j++ {
				id := sf.Generate()
				fmt.Println("Generated ID:", id)
			}
		}()
	}

	wg.Wait()
}
