package lighting

import (
	"context"
	"log"
	"math/rand"
	"time"

	"brick-smart-template/examples/brick-smart-lighting/pkg/httpclient"
)

type Lighting struct {
	ID         string
	httpClient *httpclient.Client
	isOn       bool
	brightness int
	colorTemp  int
	color      string
	mode       string
	scene      string
	energyUsage float64
	errorCode  string
}

func NewLighting(id string, httpClient *httpclient.Client) *Lighting {
	return &Lighting{
		ID:         id,
		httpClient: httpClient,
		isOn:       true,
		brightness: 80,
		colorTemp:  4000, // 4000K
		color:      "#FFFFFF",
		mode:       "normal",
		scene:      "living_room",
		energyUsage: 0.0,
	}
}

func (l *Lighting) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.simulate()
			l.reportStatus(ctx)
		}
	}
}

func (l *Lighting) simulate() {
	// 模拟开关状态变化
	if rand.Float64() < 0.1 {
		l.isOn = !l.isOn
	}

	if l.isOn {
		// 模拟亮度变化
		l.brightness = 60 + rand.Intn(41) // 60-100

		// 模拟色温变化
		if rand.Float64() < 0.2 {
			colorTemps := []int{2700, 3000, 4000, 5000, 6500}
			l.colorTemp = colorTemps[rand.Intn(len(colorTemps))]
		}

		// 模拟颜色变化
		if rand.Float64() < 0.15 {
			colors := []string{"#FFFFFF", "#FFD700", "#FF6B6B", "#4ECDC4", "#45B7D1"}
			l.color = colors[rand.Intn(len(colors))]
		}

		// 模拟模式变化
		if rand.Float64() < 0.2 {
			modes := []string{"normal", "reading", "relax", "party", "work"}
			l.mode = modes[rand.Intn(len(modes))]
		}

		// 模拟场景变化
		if rand.Float64() < 0.1 {
			scenes := []string{"living_room", "bedroom", "kitchen", "bathroom", "study"}
			l.scene = scenes[rand.Intn(len(scenes))]
		}

		// 模拟能耗
		l.energyUsage += 0.01 + rand.Float64()*0.02
	} else {
		l.brightness = 0
		l.mode = "off"
		l.energyUsage += 0.001 // 待机能耗
	}

	// 模拟故障
	l.errorCode = ""
	if rand.Float64() < 0.02 { // 2%概率出现故障
		errors := []string{"BULB_ERROR", "SENSOR_ERROR", "COMMUNICATION_ERROR"}
		l.errorCode = errors[rand.Intn(len(errors))]
	}
}

func (l *Lighting) reportStatus(ctx context.Context) {
	status := map[string]interface{}{
		"device_id":   l.ID,
		"device_type": "lighting",
		"timestamp":   time.Now().Unix(),
		"data": map[string]interface{}{
			"is_on":        l.isOn,
			"brightness":   l.brightness,
			"color_temp":   l.colorTemp,
			"color":        l.color,
			"mode":         l.mode,
			"scene":        l.scene,
			"energy_usage": l.energyUsage,
		},
	}

	if l.errorCode != "" {
		status["data"].(map[string]interface{})["error_code"] = l.errorCode
	}

	if err := l.httpClient.ReportStatus(ctx, status); err != nil {
		log.Printf("Failed to report status: %v", err)
	}
} 