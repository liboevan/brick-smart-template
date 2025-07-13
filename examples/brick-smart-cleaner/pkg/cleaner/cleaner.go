package cleaner

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"brick-smart-template/examples/brick-smart-cleaner/pkg/httpclient"
)

// Cleaner 扫地机实例
type Cleaner struct {
	ID          string
	httpClient  *httpclient.Client
	cycleCount  int
	totalTime   time.Duration
	startTime   time.Time
	batteryLevel int
	dustLevel    int
}

// Room 房间信息
type Room struct {
	ID   int
	Name string
}

// CleaningStatus 清理状态
type CleaningStatus struct {
	CycleCount    int           `json:"cycle_count"`
	RoomID        int           `json:"room_id"`
	RoomName      string        `json:"room_name"`
	Progress      int           `json:"progress"`      // 0-100
	TotalTime     time.Duration `json:"total_time"`
	CurrentTime   time.Duration `json:"current_time"`
	Status        string        `json:"status"`
	BatteryLevel  int           `json:"battery_level"`
	DustLevel     int           `json:"dust_level"`
	ErrorCode     string        `json:"error_code,omitempty"`
}

// NewCleaner 创建新的扫地机实例
func NewCleaner(id string, httpClient *httpclient.Client) *Cleaner {
	return &Cleaner{
		ID:          id,
		httpClient:  httpClient,
		startTime:   time.Now(),
		batteryLevel: 100,
		dustLevel:    0,
	}
}

// StartCleaning 开始清理任务
func (c *Cleaner) StartCleaning(ctx context.Context) {
	rooms := []Room{
		{ID: 1, Name: "客厅"},
		{ID: 2, Name: "卧室"},
		{ID: 3, Name: "厨房"},
		{ID: 4, Name: "卫生间"},
		{ID: 5, Name: "书房"},
		{ID: 6, Name: "阳台"},
	}

	log.Printf("Cleaner %s started cleaning", c.ID)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Cleaner %s received shutdown signal", c.ID)
			return
		default:
			c.cycleCount++
			log.Printf("Cleaner %s starting cycle %d", c.ID, c.cycleCount)

			// 清理每个房间
			for _, room := range rooms {
				select {
				case <-ctx.Done():
					return
				default:
					c.cleanRoom(ctx, room)
				}
			}

			log.Printf("Cleaner %s completed cycle %d", c.ID, c.cycleCount)
		}
	}
}

// cleanRoom 清理单个房间
func (c *Cleaner) cleanRoom(ctx context.Context, room Room) {
	log.Printf("Cleaner %s starting to clean room %d (%s)", c.ID, room.ID, room.Name)

	roomStartTime := time.Now()
	progressTicker := time.NewTicker(1 * time.Second) // 每秒更新进度
	defer progressTicker.Stop()

	// 模拟清理过程（100秒）
	cleanDuration := 100 * time.Second
	roomEndTime := roomStartTime.Add(cleanDuration)

	for {
		select {
		case <-ctx.Done():
			return
		case <-progressTicker.C:
			now := time.Now()
			if now.After(roomEndTime) {
				// 房间清理完成
				c.reportStatus(ctx, CleaningStatus{
					CycleCount:  c.cycleCount,
					RoomID:      room.ID,
					RoomName:    room.Name,
					Progress:    100,
					TotalTime:   time.Since(c.startTime),
					CurrentTime: time.Since(roomStartTime),
					Status:      "completed",
					BatteryLevel: c.batteryLevel,
					DustLevel:    c.dustLevel,
				})
				log.Printf("Cleaner %s completed cleaning room %d (%s)", c.ID, room.ID, room.Name)
				return
			}

			// 计算进度
			elapsed := now.Sub(roomStartTime)
			progress := int((elapsed.Seconds() / cleanDuration.Seconds()) * 100)
			if progress > 100 {
				progress = 100
			}

			// 模拟电池消耗和灰尘积累
			c.batteryLevel = max(0, c.batteryLevel - rand.Intn(2))
			c.dustLevel = min(100, c.dustLevel + rand.Intn(3))

			// 添加一些随机性，模拟真实的清理过程
			if rand.Float64() < 0.1 { // 10%概率遇到障碍物
				progress = progress - 1
				if progress < 0 {
					progress = 0
				}
			}

			// 模拟错误情况
			errorCode := ""
			if rand.Float64() < 0.05 { // 5%概率出现错误
				errorCode = "OBSTACLE_DETECTED"
			}

			// 上报状态
			c.reportStatus(ctx, CleaningStatus{
				CycleCount:  c.cycleCount,
				RoomID:      room.ID,
				RoomName:    room.Name,
				Progress:    progress,
				TotalTime:   time.Since(c.startTime),
				CurrentTime: elapsed,
				Status:      "cleaning",
				BatteryLevel: c.batteryLevel,
				DustLevel:    c.dustLevel,
				ErrorCode:    errorCode,
			})

			// 打印进度
			if progress%10 == 0 { // 每10%打印一次
				log.Printf("Cleaner %s cleaning room %d (%s): %d%%, Battery: %d%%, Dust: %d%%", 
					c.ID, room.ID, room.Name, progress, c.batteryLevel, c.dustLevel)
			}
		}
	}
}

// reportStatus 通过HTTP上报状态
func (c *Cleaner) reportStatus(ctx context.Context, status CleaningStatus) {
	// 构建通用的设备状态消息
	deviceStatus := map[string]interface{}{
		"device_id":   c.ID,
		"device_type": "cleaner",
		"timestamp":   time.Now().Unix(),
		"data": map[string]interface{}{
			"cycle_count":   status.CycleCount,
			"room_id":       status.RoomID,
			"room_name":     status.RoomName,
			"progress":      status.Progress,
			"total_time":    status.TotalTime.Seconds(),
			"current_time":  status.CurrentTime.Seconds(),
			"status":        status.Status,
			"battery_level": status.BatteryLevel,
			"dust_level":    status.DustLevel,
		},
	}

	if status.ErrorCode != "" {
		deviceStatus["data"].(map[string]interface{})["error_code"] = status.ErrorCode
	}

	// 通过HTTP客户端发送状态
	if err := c.httpClient.ReportStatus(ctx, deviceStatus); err != nil {
		log.Printf("Failed to report status: %v", err)
	}
}

// GetStatus 获取当前状态
func (c *Cleaner) GetStatus() CleaningStatus {
	return CleaningStatus{
		CycleCount:  c.cycleCount,
		TotalTime:   time.Since(c.startTime),
		Status:      "running",
		BatteryLevel: c.batteryLevel,
		DustLevel:    c.dustLevel,
	}
}

// String 返回扫地机的字符串表示
func (c *Cleaner) String() string {
	return fmt.Sprintf("Cleaner{ID: %s, Cycles: %d, Running: %v, Battery: %d%%, Dust: %d%%}", 
		c.ID, c.cycleCount, time.Since(c.startTime), c.batteryLevel, c.dustLevel)
}

// 辅助函数
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 