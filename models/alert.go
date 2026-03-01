package models

type AlertLevel string

const (
	INFO     AlertLevel = "INFO"
	WARNING  AlertLevel = "WARNING"
	CRITICAL AlertLevel = "CRITICAL"
)

type Alert struct {
	Title          string
	Message        string
	Recommendation string
	Level          AlertLevel
	Time           string
}