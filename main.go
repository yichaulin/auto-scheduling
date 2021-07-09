package main

import (
	"fmt"
	"os"

	"auto-scheduling/schedule"

	"gopkg.in/yaml.v3"
)

func main() {
	file, err := os.ReadFile("./shedule.yml")
	if err != nil {
		fmt.Println(err)
	}

	var schedule schedule.Schedule
	yaml.Unmarshal(file, &schedule)
	err = schedule.PrintShifts()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\nAuto Scheduling")
}
