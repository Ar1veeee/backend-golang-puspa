package helpers

import (
	"time"
)

func CalculateAge(birthDate time.Time) int {
	today := time.Now()
	age := today.Year() - birthDate.Year()

	if int(today.Month()) < int(birthDate.Month()) ||
		(int(today.Month()) == int(birthDate.Month()) && int(today.Day()) < int(birthDate.Day())) {
		age--
	}

	return age
}
