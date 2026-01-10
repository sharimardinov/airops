package models

const (
	// Flight statuses
	FlightStatusScheduled = "Scheduled"
	FlightStatusOnTime    = "On Time"
	FlightStatusDelayed   = "Delayed"
	FlightStatusCancelled = "Cancelled"
	FlightStatusDeparted  = "Departed"
	FlightStatusArrived   = "Arrived"

	// Fare classes
	FareClassEconomy  = "Economy"
	FareClassComfort  = "Comfort"
	FareClassBusiness = "Business"

	// Booking statuses
	BookingStatusActive    = "Active"
	BookingStatusCancelled = "Cancelled"
)
