package schedule

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

type Slot struct {
	OpeningTimeIndex int      `yaml:"opening_time_index"`
	Workers          []string `yaml:"workers"`
}

type Shift struct {
	Date  int    `yaml:"date"`
	Slots []Slot `yaml:"slots"`
}

type Shifts []Shift

type Schedule struct {
	Year        int      `yaml:"year"`
	Month       int      `yaml:"month"`
	OpeningTime []string `yaml:"opening_time"`
	Shifts      Shifts   `yaml:"shifts"`
}

func (schedule Schedule) PrintShifts() error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"opening_time", "Sun", "Mon", "Tue", "Web", "Thu", "Fri", "Sat"})
	table.SetRowLine(true)
	table.SetAutoWrapText(false)

	year := schedule.Year
	month := schedule.Month
	openingTime := schedule.OpeningTime
	shifts := schedule.Shifts

	loc, _ := time.LoadLocation("Asia/Singapore")
	firstSlotDatetime, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%d-%02d-%02d", year, month, shifts[0].Date), loc)

	if err != nil {
		return err
	}

	weekDay := int(firstSlotDatetime.Weekday())
	shifts = append(make([]Shift, weekDay), shifts...)

	var chunkedShifts [][]Shift

	chunkedCount := 0
	chunkNum := 7
	for chunkedCount < len(shifts) {
		rest := len(shifts) - chunkedCount
		if len(shifts)-chunkedCount < chunkNum {
			chunkNum = rest
		}
		chunkedShifts = append(chunkedShifts, shifts[chunkedCount:chunkedCount+chunkNum])
		chunkedCount += chunkNum
	}

	for _, shifts := range chunkedShifts {
		for i, openingTime := range append([]string{""}, openingTime...) {
			weekData := []string{openingTime}
			for _, shift := range shifts {
				if openingTime == "" {
					weekDay := ""
					if shift.Date > 0 {
						weekDay = strconv.Itoa(shift.Date)
					}

					weekData = append(weekData, weekDay)
					continue
				}

				workers := shift.getWorkersByOpeningTimeIndex(i - 1)
				weekData = append(weekData, strings.Join(workers, "\n"))
			}
			table.Append(weekData)
		}
	}

	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetCaption(true, fmt.Sprintf("%d-%d", year, month))
	table.Render()

	return nil
}

func (shift Shift) getWorkersByOpeningTimeIndex(index int) []string {
	for _, slot := range shift.Slots {
		if slot.OpeningTimeIndex == index {
			return slot.Workers
		}
	}

	return make([]string, 0)
}
