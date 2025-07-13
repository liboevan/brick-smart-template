package thermostat

import (
	"context"
	"log"
	"math/rand"
	"time"

	"brick-smart-template/examples/brick-smart-thermostat/pkg/httpclient"
)

type Thermostat struct {
	ID         string
	httpClient *httpclient.Client
	mode       string
	targetTemp float64
	roomTemp   float64
	humidity   float64
	energyUsage float64
	errorCode  string
	fanSpeed   int
}

func NewThermostat(id string, httpClient *httpclient.Client) *Thermostat {
	return &Thermostat{
		ID:         id,
		httpClient: httpClient,
		mode:       "auto",
		targetTemp: 22.0,
		roomTemp:   20.0 + rand.Float64()*4,
		humidity:   45.0 + rand.Float64()*20,
		energyUsage: 0.0,
		fanSpeed:   2,
	}
}

func (t *Thermostat) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.simulate()
			t.reportStatus(ctx)
		}
	}
}

func (t *Thermostat) simulate() {
	// 模拟温度变化
	if t.roomTemp < t.targetTemp {
		t.roomTemp += 0.2 + rand.Float64()*0.2
		t.energyUsage += 0.1 + rand.Float64()*0.2
	} else if t.roomTemp > t.targetTemp {
		t.roomTemp -= 0.2 + rand.Float64()*0.2
		t.energyUsage += 0.05 + rand.Float64()*0.1
	}

	// 模拟湿度变化
	t.humidity += (rand.Float64() - 0.5) * 2
	if t.humidity < 30 {
		t.humidity = 30
	} else if t.humidity > 70 {
		t.humidity = 70
	}

	// 模拟模式变化
	if rand.Float64() < 0.1 {
		modes := []string{"auto", "eco", "comfort", "sleep"}
		t.mode = modes[rand.Intn(len(modes))]
	}

	// 模拟风扇速度变化
	if rand.Float64() < 0.15 {
		t.fanSpeed = 1 + rand.Intn(3) // 1-3档
	}

	// 模拟故障
	t.errorCode = ""
	if rand.Float64() < 0.03 { // 3%概率出现故障
		errors := []string{"SENSOR_ERROR", "COMMUNICATION_ERROR", "SYSTEM_ERROR"}
		t.errorCode = errors[rand.Intn(len(errors))]
	}
}

func (t *Thermostat) reportStatus(ctx context.Context) {
	status := map[string]interface{}{
		"device_id":   t.ID,
		"device_type": "thermostat",
		"timestamp":   time.Now().Unix(),
		"data": map[string]interface{}{
			"room_temp":    t.roomTemp,
			"target_temp":  t.targetTemp,
			"humidity":     t.humidity,
			"mode":         t.mode,
			"energy_usage": t.energyUsage,
			"fan_speed":    t.fanSpeed,
		},
	}

	if t.errorCode != "" {
		status["data"].(map[string]interface{})["error_code"] = t.errorCode
	}

	if err := t.httpClient.ReportStatus(ctx, status); err != nil {
		log.Printf("Failed to report status: %v", err)
	}
} 