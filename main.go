package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	unix      = flag.Int("unix", 0, "Unix Timestamp in seconds")
	unixMilli = flag.Int("unixmilli", 0, "Unix Timestamp in milli seconds")
	unixMicro = flag.Int("unixmicro", 0, "Unix Timestamp in micro seconds")
	in        = flag.String("in", "", "YYYY-MM-DD HH:MM:SS styled timestamp")
	timezone  = flag.String("tz", "UTC", "Timezone currently UTC|JST|PT|PST|PDT|IST are available")
)

func getLocationFromAbbreviation(abbr string) (*time.Location, error) {
	timezoneMap := map[string]string{
		"JST": "Asia/Tokyo",
		"UTC": "UTC",
		"PT":  "America/Los_Angeles",
		"PST": "America/Los_Angeles",
		"PDT": "America/Los_Angeles",
		"IST": "Asia/Kolkata",
	}
	abbr = strings.ToUpper(abbr)
	tz, ok := timezoneMap[abbr]
	if !ok {
		return nil, fmt.Errorf("timezone abbreviation '%s' not found", abbr)
	}
	location, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("failed to load location for '%s': %w", tz, err)
	}
	return location, nil
}

func parsePartialDate(s string, location *time.Location) (int64, error) {
	formats := []string{
		"2006",
		"2006-01",
		"2006-01-02",
		"2006-01-02 15",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.000",
	}
	for _, format := range formats {
		if len(s) == len(format) {
			t, err := time.ParseInLocation(format, s, location)
			if err == nil {
				return t.Unix(), nil
			}
		}
	}
	return 0, fmt.Errorf("invalid date format: '%s'", s)
}

// parses inputs then returns unix timestamp in second
func normalize() (int64, error) {
	if *in == "" {
		return time.Now().Unix(), nil
	}
	if *unix != 0 {
		return int64(*unix), nil
	}
	if *unixMilli != 0 {
		return int64(*unixMilli) / 1000, nil
	}
	if *unixMicro != 0 {
		return int64(*unixMicro) / 1000000, nil
	}
	// Parse datetime from given string
	location, err := getLocationFromAbbreviation(*timezone)
	if err != nil {
		return 0, err
	}
	t, err := parsePartialDate(*in, location)
	if err != nil {
		return 0, err
	}
	return t, nil
}

func buildOutput(unixTimestamp int64) (string, error) {
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return "", err
	}
	pt, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return "", err
	}
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return "", err
	}
	t := time.Unix(unixTimestamp, 0)
	res := fmt.Sprintf("Normalized: %d\n%s\n%s\n%s\n", unixTimestamp, t.In(utc), t.In(jst), t.In(pt))
	return res, nil
}

func main() {
	flag.Parse()
	unixTimestamp, err := normalize()
	if err != nil {
		log.Fatalln(err)
	}
	output, err := buildOutput(unixTimestamp)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Print(output)
}
