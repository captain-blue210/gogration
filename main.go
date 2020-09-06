package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	HEADER = []string{"User", "Email", "Client", "Project", "Task", "Description", "Billable", "Start date", "Start time", "End date", "End time", "Duration", "Tags", "Amount ()"}
	USER   = "ぶるー"
	EMAIL  = "j5hfca7pm@gmail.com"
	SUFFIX = ":00"
	FORMAT = ""
)

func readCSV() [][]string {
	// CSVを読み込む
	file, err := os.Open("detailed.csv")
	if err != nil {
		log.Fatalf("File open error. %v", err)
	}

	r := csv.NewReader(file)

	var timecampData [][]string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("CSV read error. %v", err)
		}

		timecampData = append(timecampData, record)
	}
	return timecampData
}

func converttoTogglData(timecampData [][]string) [][]string {
	// toggl移行用のCSVに加工する
	var togglData [][]string
	// ヘッダーの設定
	togglData = append(togglData, HEADER)
	for i, row := range timecampData {
		if i < 2 || i == len(timecampData)-1 {
			// 最初の2行と最後の1行は不要
			continue
		}

		startDate := row[0]
		endDate := row[0]
		description := row[2]
		timestamp := strings.Split(row[4], "-")
		startTime := timestamp[0] + SUFFIX
		endTime := timestamp[1] + SUFFIX

		// durationを計算する
		start, _ := time.Parse("2006-01-02 15:04:05", startDate+" "+startTime)
		end, _ := time.Parse("2006-01-02 15:04:05", endDate+" "+endTime)
		tduration := end.Sub(start)
		// durationをtogglの形式合わせて加工
		duration := formatDuration(tduration)

		project := row[7]

		var newRow []string
		newRow = append(newRow, USER)
		newRow = append(newRow, EMAIL)
		newRow = append(newRow, "") //client
		newRow = append(newRow, project)
		newRow = append(newRow, "") // task
		newRow = append(newRow, description)
		newRow = append(newRow, "No") // billable
		newRow = append(newRow, startDate)
		newRow = append(newRow, startTime)
		newRow = append(newRow, endDate)
		newRow = append(newRow, endTime)
		newRow = append(newRow, duration)
		newRow = append(newRow, "") //tags
		newRow = append(newRow, "") // amount
		togglData = append(togglData, newRow)
	}
	return togglData
}

func formatDuration(beforeDur time.Duration) string {
	re := regexp.MustCompile("[hm]")
	tmp := re.Split(strings.TrimRight(beforeDur.String(), "s"), -1)

	durslice := []string{}
	for _, elm := range tmp {
		if len(tmp) == 3 {
			elm = paddingZeroToElm(elm)
		} else if len(tmp) == 2 && len(durslice) == 0 {
			durslice = paddingZeroToDurSlice(durslice, 1)
		} else if len(tmp) == 2 && len(elm) < 2 {
			elm = "0" + elm
		} else if len(tmp) == 1 {
			durslice = paddingZeroToDurSlice(durslice, 2)
			elm = paddingZeroToElm(elm)
		}
		durslice = append(durslice, elm)
	}

	return strings.Join(durslice, ":")
}

func paddingZeroToElm(elm string) string {
	if len(elm) < 2 {
		elm = "0" + elm
	}
	return elm
}

func paddingZeroToDurSlice(durslice []string, n int) []string {
	for i := 0; i < n; i++ {
		durslice = append(durslice, "00")
	}
	return durslice
}

func writeTogglCSV(togglData [][]string) error {
	file, err := os.Create("toggl.csv")
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)

	for _, record := range togglData {
		if err := w.Write(record); err != nil {
			return err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

func main() {
	// timeCampのデータを読み取る
	timecampData := readCSV()

	// toggl移行用のCSVに加工する
	togglData := converttoTogglData(timecampData)

	// toggl用CSVとして出力
	err := writeTogglCSV(togglData)
	if err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
}
