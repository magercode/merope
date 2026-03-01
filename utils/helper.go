package utils

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func GetOsType() string {
	info, err := host.Info()
	if err != nil {
		return "Unknown"
	}
	return fmt.Sprintf("%s %s (%s)", info.Platform, info.PlatformVersion, info.KernelArch)
}

func GetTotalCore() int {
	count, _ := cpu.Counts(true)
	return count
}

func GetTotalMemory() uint64 {
	v, _ := mem.VirtualMemory()
	return v.Total
}

func GetTotalDisk() uint64 {
	parts, _ := disk.Partitions(false)
	var total uint64
	for _, p := range parts {
		usage, err := disk.Usage(p.Mountpoint)
		if err == nil {
			total += usage.Total
		}
	}
	return total
}