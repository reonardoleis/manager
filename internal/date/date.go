package date

import "time"

func Today() []string {
	return []string{time.Now().Format("2006-01-02")}
}

func Yesterday() []string {
	return []string{time.Now().AddDate(0, 0, -1).Format("2006-01-02")}
}

func DaysAgo(n int) []string {
	return []string{time.Now().AddDate(0, 0, -n).Format("2006-01-02")}
}

func Days(daysAgo int) []string {
	days := []string{}
	for i := 0; i < daysAgo; i++ {
		days = append(days, DaysAgo(i)[0])
	}
	return days
}
