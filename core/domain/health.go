package domain

// Health is the domain model for the health.
type Health struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Version  string `json:"version"`
	Date     string `json:"date"`
	Datetime string `json:"datetime"`
	Time     string `json:"time"`
	Timezone string `json:"timezone"`
	Revision string `json:"revision"`
	Redis    string `json:"redis"`
	Mongo    string `json:"mongo"`
}

// NewHealth returns a new Health.
func NewHealth(status, datetime, message, version, date, time, timezone, revision, redisStatus, mongoStatus string) *Health {
	return &Health{
		Status:   status,
		Message:  message,
		Version:  version,
		Date:     date,
		Datetime: datetime,
		Time:     time,
		Timezone: timezone,
		Revision: revision,
		Redis:    redisStatus,
		Mongo:    mongoStatus,
	}
}
