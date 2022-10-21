package greeting

import (
	"os"
)

func ReadGreetingFile() string {
	greeting, err := os.ReadFile("./static/greeting.txt")
	if err != nil {
		panic(err)
	}

	return string(greeting)
}

func SetGreeting(greeting string) error {
	if err := os.Truncate("./static/greeting.txt", 0); err != nil {
		return err
	}
	data := []byte(greeting)
	err := os.WriteFile("./static/greeting.txt", data, 0644)
	if err != nil {
		panic(err)
	}
	return nil
}
