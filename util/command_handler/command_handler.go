package command_handler

import (
	"errors"
	"fmt"
	"time"

	constants "github.com/jigten/telegram-notion-bot/constants"
	greeting "github.com/jigten/telegram-notion-bot/util/greeting"
)

func HandleCommand(command string) (string, error) {
	if command == constants.GREETING_COMMAND {
		return greeting.ReadGreetingFile(), nil
	}

	if command == constants.COUNTDOWN_COMMAND {
		nextMeet := time.Date(2022, time.November, 24, 14, 0, 0, 0, time.Local)
		year, month, day, hour, min, sec, done := diff(time.Now(), nextMeet)

		if done {
			return fmt.Sprint("Countdown over. Have a nice evening together."), nil
		}

		return fmt.Sprintf("%d years, %d months, %d days, %d hours, %d mins and %d seconds till you two meet.",
			year, month, day, hour, min, sec), nil
	}

	return "", errors.New("unknown command")
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int, done bool) {
	if a.After(b) {
		done = true
		return
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.Local)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
