package models

type AirplaneStats struct {
	AirplaneCode    string  `json:"airplane_code"`
	TotalFlights    int     `json:"total_flights"`
	TotalPassengers int     `json:"total_passengers"`
	AvgLoadFactor   float64 `json:"avg_load_factor"` // процент загрузки
}
