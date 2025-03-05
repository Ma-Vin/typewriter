package common

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	ASTERISK       string = "*"
	CRON_INTERVAL  string = "-"
	CRON_INCREMENT string = "/"
	CRON_LIST      string = ","
)

// Structure which represent a cron format expresion and provides the next point in time
type Crontab struct {
	NextTime            *time.Time
	location            *time.Location
	minutes             []int
	minuteIndex         int
	hours               []int
	hourIndex           int
	daysOfMonth         []int
	months              []int
	monthIndex          int
	daysOfWeek          []int
	year                int
	allDaysOfMonth      []int
	allDaysOfMonthIndex int
}

// Returns the minutes of the next point in time
func (s *Crontab) Minute() int {
	return getTimeElement(&s.minutes, &s.minuteIndex, true)
}

// Returns the hour of the next point in time
func (s *Crontab) Hour() int {
	return getTimeElement(&s.hours, &s.hourIndex, true)
}

// Returns the day of month of the next point in time
func (s *Crontab) DayOfMonth() int {
	return getTimeElement(&s.allDaysOfMonth, &s.allDaysOfMonthIndex, false)
}

// Returns the month of the next point in time
func (s *Crontab) Month() int {
	return getTimeElement(&s.months, &s.monthIndex, false)
}

// Returns the year of the next point in time
func (s *Crontab) Year() int {
	return s.year
}

// determines the time value from a Crontab element
func getTimeElement(timeUnitArray *[]int, index *int, zerobased bool) int {
	if len(*timeUnitArray) == 0 {
		if zerobased {
			return *index
		} else {
			return *index + 1
		}
	}
	return (*timeUnitArray)[*index]
}

// Calculates the next point in time after NextTime
func (s *Crontab) CalculateNextTime() {
	if s.NextTime == nil {
		s.initializeNextTimeFromCurrentDate()
		return
	}
	s.increaseMinute()
	s.setNextTime()
}

// initializes the crontab time inidices and the last NextTime derived from current time
func (s *Crontab) initializeNextTimeFromCurrentDate() {
	base := getNow()

	s.location = base.Location()
	s.year = base.Year()

	refMonth := int(base.Month())
	interateToLastReachedTimeUnitIndex(&s.months, &s.monthIndex, refMonth, false)
	if s.Month() > refMonth {
		s.allDaysOfMonthIndex = 0
		s.hourIndex = 0
		s.minuteIndex = 0
		s.determineAllDaysOfMonth()
		s.setNextTime()
		return
	}
	if s.Month() < refMonth {
		s.monthIndex = 0
		s.allDaysOfMonthIndex = 0
		s.hourIndex = 0
		s.minuteIndex = 0
		s.year++
		s.determineAllDaysOfMonth()
		s.setNextTime()
		return
	}

	refDayOfMonth := base.Day()
	s.determineAllDaysOfMonth()
	interateToLastReachedTimeUnitIndex(&s.allDaysOfMonth, &s.allDaysOfMonthIndex, refDayOfMonth, false)
	if s.DayOfMonth() > refDayOfMonth {
		s.hourIndex = 0
		s.minuteIndex = 0
		s.setNextTime()
		return
	}
	if s.DayOfMonth() < refDayOfMonth {
		s.allDaysOfMonthIndex = 0
		s.hourIndex = 0
		s.minuteIndex = 0
		s.increaseMonth()
		s.setNextTime()
		return
	}

	refHour := base.Hour()
	interateToLastReachedTimeUnitIndex(&s.hours, &s.hourIndex, refHour, true)
	if s.Hour() > refHour {
		s.minuteIndex = 0
		s.setNextTime()
		return
	}
	if s.Hour() < refHour {
		s.hourIndex = 0
		s.minuteIndex = 0
		s.increaseDayOfMonth()
		s.setNextTime()
		return
	}

	refMinute := base.Minute()
	interateToLastReachedTimeUnitIndex(&s.minutes, &s.minuteIndex, refMinute, true)
	if s.Minute() < refMinute {
		s.minuteIndex = 0
		s.increaseHour()
	}
	if s.Minute() == refMinute {
		s.increaseMinute()
	}

	s.setNextTime()
}

// determines all days of month by using all days defined by daysOfWeek, by using all days defined by daysOfMonth or merging daysOfWeek and daysOfMonth into allDaysOfMonth
func (s *Crontab) determineAllDaysOfMonth() {
	if len(s.daysOfWeek) == 0 && len(s.daysOfMonth) == 0 {
		return
	}

	if len(s.daysOfWeek) == 0 && len(s.daysOfMonth) > 0 && len(s.allDaysOfMonth) == 0 {
		s.allDaysOfMonth = make([]int, len(s.daysOfMonth))
		copy(s.allDaysOfMonth, s.daysOfMonth)
		return
	}

	if len(s.daysOfWeek) > 0 && len(s.daysOfMonth) == 0 {
		s.determineAllDaysOfMonthOnlyByDaysOfWeek()
		return
	}

	if len(s.daysOfWeek) > 0 && len(s.daysOfMonth) > 0 {
		s.determineAllDaysOfMonthByMergeOfDaysOfWeekAndMonth()
	}
}

// determines all days of month by using all days defined by daysOfWeek
func (s *Crontab) determineAllDaysOfMonthOnlyByDaysOfWeek() {
	weekday := int(time.Date(s.Year(), time.Month(s.Month()), 1, 0, 0, 0, 0, s.location).Weekday())
	maxDaysInMonth := s.getMaxIndexDayOfMonth()
	for i := 0; i <= maxDaysInMonth; i++ {
		if slices.Contains(s.daysOfWeek, weekday) {
			s.allDaysOfMonth = append(s.allDaysOfMonth, i+1)
		}
		weekday = (weekday + 1) % 7
	}
}

// determines all days of month by merging daysOfWeek and daysOfMonth into allDaysOfMonth
func (s *Crontab) determineAllDaysOfMonthByMergeOfDaysOfWeekAndMonth() {
	s.allDaysOfMonth = make([]int, len(s.daysOfMonth))
	copy(s.allDaysOfMonth, s.daysOfMonth)
	weekday := int(time.Date(s.Year(), time.Month(s.Month()), 1, 0, 0, 0, 0, s.location).Weekday())
	maxDaysInMonth := s.getMaxIndexDayOfMonth()
	for i := 0; i <= maxDaysInMonth; i++ {
		if !slices.Contains(s.daysOfMonth, i+1) && slices.Contains(s.daysOfWeek, weekday) {
			s.allDaysOfMonth = append(s.allDaysOfMonth, i+1)
		}
		weekday = (weekday + 1) % 7
	}
	slices.Sort(s.allDaysOfMonth)
}

// adjust the index of a time elements to the last valid value compared to targetValue
func interateToLastReachedTimeUnitIndex(timeUnitArray *[]int, index *int, targetValue int, zerobased bool) {
	if len(*timeUnitArray) == 0 {
		if zerobased {
			*index = targetValue
		} else {
			*index = targetValue - 1
		}
		return
	} else {
		for *index = 0; *index < len(*timeUnitArray); (*index)++ {
			if (*timeUnitArray)[*index] >= targetValue {
				break
			}
		}
		if len(*timeUnitArray) <= *index {
			*index = len(*timeUnitArray) - 1
		}
	}
	if *index < 0 {
		*index = 0
	}
}

// returns a [time.Time] representation of the current indices of the crontab
func (s *Crontab) setNextTime() {
	nextTime := time.Date(s.Year(), time.Month(s.Month()), s.DayOfMonth(), s.Hour(), s.Minute(), 0, 0, s.location)
	s.NextTime = &nextTime
}

// increases the minute by one unit. If the last element was reached already, the hour will be increased
func (s *Crontab) increaseMinute() {
	increaseTimeUnit(&s.minutes, &s.minuteIndex, 59, true, s.increaseHour)
}

// increases the hour by one unit. If the last element was reached already, the day of month will be increased
func (s *Crontab) increaseHour() {
	increaseTimeUnit(&s.hours, &s.hourIndex, 23, true, s.increaseDayOfMonth)
}

// increases the day of month by one unit. If the last element was reached already, the month will be increased
func (s *Crontab) increaseDayOfMonth() {
	increaseTimeUnit(&s.allDaysOfMonth, &s.allDaysOfMonthIndex, s.getMaxIndexDayOfMonth(), false, s.increaseMonth)
}

// increases the month by one unit. If the last element was reached already, the year will be increased
func (s *Crontab) increaseMonth() {
	increaseTimeUnit(&s.months, &s.monthIndex, 11, false, func() { s.year++ })
	s.determineAllDaysOfMonth()
}

// increases a time unit of a Crontab element  If the last element was reached already, the following unit will be increased by given nextIncrease
func increaseTimeUnit(timeUnitArray *[]int, index *int, maxTimeUnitIndex int, zerobased bool, nextIncrease func()) {
	(*index)++
	if isUpperBoundExceededAllElements(timeUnitArray, index, maxTimeUnitIndex) || isUpperBoundExceededLengthlements(timeUnitArray, index) ||
		isUpperBoundExceededMaxIndexlements(timeUnitArray, index, maxTimeUnitIndex, zerobased) {

		*index = 0
		nextIncrease()
	}
}

// Checks if the maximum value is exceeded by the index of an empty timeUnitArray
func isUpperBoundExceededAllElements(timeUnitArray *[]int, index *int, maxTimeUnitIndex int) bool {
	return len(*timeUnitArray) == 0 && *index >= maxTimeUnitIndex
}

// Checks if the index is greater or equal compared to the length of a filled timeUnitArray
func isUpperBoundExceededLengthlements(timeUnitArray *[]int, index *int) bool {
	return len(*timeUnitArray) > 0 && len(*timeUnitArray) <= *index
}

// Checks if the maximum value defined by timeUnitArray is exceeded by the value at the index of a filled timeUnitArray
func isUpperBoundExceededMaxIndexlements(timeUnitArray *[]int, index *int, maxTimeUnitIndex int, zerobased bool) bool {
	var refMaxValue int
	if zerobased {
		refMaxValue = maxTimeUnitIndex
	} else {
		refMaxValue = maxTimeUnitIndex + 1
	}
	return len(*timeUnitArray) > 0 && len(*timeUnitArray) > *index && (*timeUnitArray)[*index] > refMaxValue
}

// returns the maximal index of a day of month
func (s *Crontab) getMaxIndexDayOfMonth() int {
	var monthToCheck int
	if len(s.months) == 0 {
		monthToCheck = s.monthIndex + 1
	} else {
		monthToCheck = s.months[s.monthIndex]
	}
	switch monthToCheck {
	case 2:
		if s.year%4 == 0 && s.year%100 != 0 {
			return 28
		} else {
			return 27
		}
	case 4, 6, 9, 11:
		return 29
	default:
		return 30
	}
}

// Parses the giuven cron expresion and creates a new Crontab element with initialized NextTime compared to [time.Now]
func CreateCrontab(cronExpression string) *Crontab {
	cronParts := strings.Split(cronExpression, " ")
	result := Crontab{nil, nil, []int{}, -1, []int{}, -1, []int{}, []int{}, -1, []int{}, 0, []int{}, -1}

	if len(cronParts) > 0 {
		determineCronElement(cronParts[0], 0, 59, &result.minutes)
	}
	if len(cronParts) > 1 {
		determineCronElement(cronParts[1], 0, 23, &result.hours)
	}
	if len(cronParts) > 2 {
		determineCronElement(cronParts[2], 1, 31, &result.daysOfMonth)
	}
	if len(cronParts) > 3 {
		determineCronElement(cronParts[3], 1, 12, &result.months)
	}
	if len(cronParts) > 4 {
		determineCronElement(cronParts[4], 0, 6, &result.daysOfWeek)
	}

	result.CalculateNextTime()

	return &result
}

// Determines the array representation of a single element of a cron expression with respecet to its given lower and upper bounds
func determineCronElement(cronElementExpression string, lowerBound int, upperBound int, cronRepresentation *[]int) {
	if cronElementExpression == ASTERISK {
		return
	}

	if strings.Contains(cronElementExpression, CRON_INTERVAL) {
		determineCronIntervalElement(&cronElementExpression, lowerBound, upperBound, cronRepresentation)
		return
	}

	if strings.Contains(cronElementExpression, CRON_LIST) {
		determineCronListElement(&cronElementExpression, lowerBound, upperBound, cronRepresentation)
		return
	}

	if strings.Contains(cronElementExpression, CRON_INCREMENT) {
		determineCronIncrementElement(&cronElementExpression, lowerBound, upperBound, cronRepresentation)
		return
	}

	determineCronConstantElement(&cronElementExpression, lowerBound, upperBound, cronRepresentation)
}

// Determines the array representation of a single element of a cron interval expression with respecet to its given lower and upper bounds
func determineCronIntervalElement(cronElementExpression *string, lowerBound int, upperBound int, cronRepresentation *[]int) {
	interval := strings.Split(*cronElementExpression, CRON_INTERVAL)
	from := getIntOrDefault(interval[0], lowerBound, true)
	to, increment := getIntervalBoundAndIncrement(interval[1], upperBound)
	if from == lowerBound && to == upperBound && increment == 1 {
		return
	}
	for ; from <= to; from += increment {
		*cronRepresentation = append(*cronRepresentation, from)
	}
}

// Determines the array representation of a single element of a cron list expression with respecet to its given lower and upper bounds
func determineCronListElement(cronElementExpression *string, lowerBound int, upperBound int, cronRepresentation *[]int) {
	for _, entry := range strings.Split(*cronElementExpression, CRON_LIST) {
		entryAsInt, err := strconv.Atoi(entry)
		if err == nil && lowerBound <= entryAsInt && entryAsInt <= upperBound {
			*cronRepresentation = append(*cronRepresentation, entryAsInt)
		}
	}
}

// Determines the array representation of a single element of a cron increment expression with respecet to its given lower and upper bounds
func determineCronIncrementElement(cronElementExpression *string, lowerBound int, upperBound int, cronRepresentation *[]int) {
	incrementSplit := strings.Split(*cronElementExpression, CRON_INCREMENT)
	increment := getIntOrDefault(incrementSplit[1], 1, true)
	if increment != 1 {
		for i := lowerBound; i <= upperBound; i += increment {
			*cronRepresentation = append(*cronRepresentation, i)
		}
	}
}

// Determines the array representation of a single element of a cron constant expression with respecet to its given lower and upper bounds
func determineCronConstantElement(cronElementExpression *string, lowerBound int, upperBound int, cronRepresentation *[]int) {
	entryAsInt, err := strconv.Atoi(*cronElementExpression)
	if err == nil && lowerBound <= entryAsInt && entryAsInt <= upperBound {
		*cronRepresentation = append(*cronRepresentation, entryAsInt)
	}
}

// Returns the upper bound and the increament for an intervall for given values of rightside of ”-”
func getIntervalBoundAndIncrement(boundAndIncrementExpression string, boundDefault int) (int, int) {
	if strings.Contains(boundAndIncrementExpression, CRON_INCREMENT) {
		boundAndIncrement := strings.Split(boundAndIncrementExpression, CRON_INCREMENT)
		return getIntOrDefault(boundAndIncrement[0], boundDefault, false), getIntOrDefault(boundAndIncrement[1], 1, true)
	}
	return getIntOrDefault(boundAndIncrementExpression, boundDefault, false), 1
}

// Returns the number representation of given boundExpression.
// If it's not a number, the default value will be returned.
// If the value is lower than the default and its a lower bound, the default value will be returned.
// If the value is grater than the default and its not a lower bound, the default value will be returned
func getIntOrDefault(boundExpression string, boundDefault int, isLowerBound bool) int {
	if len(boundExpression) == 0 {
		return boundDefault
	}
	bound, err := strconv.Atoi(boundExpression)
	if err != nil || (isLowerBound && bound < boundDefault) || (!isLowerBound && boundDefault < bound) {
		return boundDefault
	}
	return bound
}
