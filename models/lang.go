package models

type Language string

const (
	EN Language = "en"
	ID Language = "id"
)

type Message struct {
	CPUHigh        string `json:"cpu_high"`
	MemoryHigh     string `json:"memory_high"`
	DiskHigh       string `json:"disk_high"`
	SystemInfo     string `json:"system_info"`
	AlertTitle     string `json:"alert_title"`
	AlertMessage   string `json:"alert_message"`
	Time           string `json:"time"`
	Level          string `json:"level"`
	OK             string `json:"ok"`
	Warning        string `json:"warning"`
	Critical       string `json:"critical"`
	ServiceStarted string `json:"service_started"`
}

type LanguagePack struct {
	EN Message
	ID Message
}