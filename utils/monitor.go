package utils

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"merope/models"
)

func CheckSystem() (*models.Alert, error) {
	cpuPercent, _ := cpu.Percent(time.Second, false)
	vmStat, _ := mem.VirtualMemory()

	if cpuPercent[0] > 80 {
		return &models.Alert{
			Title:   "CPU Usage High",
			Message: fmt.Sprintf("CPU Usage: %.2f%%", cpuPercent[0]),
			Level:   models.CRITICAL,
			Time:    time.Now().Format(time.RFC3339),
		}, nil
	}

	if vmStat.UsedPercent > 80 {
		return &models.Alert{
			Title:   "Memory Usage High",
			Message: fmt.Sprintf("RAM Usage: %.2f%%", vmStat.UsedPercent),
			Level:   models.WARNING,
			Time:    time.Now().Format(time.RFC3339),
		}, nil
	}

	return nil, nil
}