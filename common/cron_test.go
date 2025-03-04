package common

import (
	"strconv"
	"testing"
	"time"

	"github.com/ma-vin/typewriter/testutil"
)

var cronCrontabTestTime = time.Date(2025, time.March, 1, 16, 10, 20, 30, time.UTC)

func checkSliceValues(length int, startValue int, increment int, slice *[]int, t *testing.T, objectName string) {
	testutil.AssertEquals(length, len(*slice), t, "len("+objectName+")")
	value := startValue
	for i := 0; i < length; i++ {
		testutil.AssertEquals(value, (*slice)[i], t, objectName+"["+strconv.Itoa(i)+"]")
		value += increment
	}
}

func checkEmptyMinute(result *Crontab, t *testing.T) {
	testutil.AssertEquals(0, len(result.minutes), t, "len(result.minute)")
}
func checkEmptyHour(result *Crontab, t *testing.T) {
	testutil.AssertEquals(0, len(result.hours), t, "len(result.hour)")
}
func checkEmptyDayOfMonth(result *Crontab, t *testing.T) {
	testutil.AssertEquals(0, len(result.daysOfMonth), t, "len(result.dayOfMonth)")
}
func checkEmptyMonth(result *Crontab, t *testing.T) {
	testutil.AssertEquals(0, len(result.months), t, "len(result.month)")
}
func checkEmptyDayOfWeekh(result *Crontab, t *testing.T) {
	testutil.AssertEquals(0, len(result.daysOfWeek), t, "len(result.dayOfWeek)")
}

func TestCreateCronCrontabAllAsterisk(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * * *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 11, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabAllConstatnts(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("1 2 3 4 5")

	checkSliceValues(1, 1, 1, &result.minutes, t, "result.minute")
	checkSliceValues(1, 2, 1, &result.hours, t, "result.hour")
	checkSliceValues(1, 3, 1, &result.daysOfMonth, t, "result.dayOfMonth")
	checkSliceValues(1, 4, 1, &result.months, t, "result.month")
	checkSliceValues(1, 5, 1, &result.daysOfWeek, t, "result.dayOfWeek")

	testutil.AssertEquals(time.Date(2025, time.April, 3, 2, 1, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabAllDifferent(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("11,15,19 */4 5-10 3-11/3 *")

	checkSliceValues(3, 11, 4, &result.minutes, t, "result.minute")
	checkSliceValues(6, 0, 4, &result.hours, t, "result.hour")
	checkSliceValues(6, 5, 1, &result.daysOfMonth, t, "result.dayOfMonth")
	checkSliceValues(3, 3, 3, &result.months, t, "result.month")
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 5, 0, 11, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

// Minute
func TestCreateCronCrontabMinuteRange(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("30-45 * * * *")

	checkSliceValues(16, 30, 1, &result.minutes, t, "result.minute")
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 30, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabMinuteRangeIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("30-45/5 * * * *")

	checkSliceValues(4, 30, 5, &result.minutes, t, "result.minute")
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 30, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabMinuteList(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("30,32,34,36 * * * *")

	checkSliceValues(4, 30, 2, &result.minutes, t, "result.minute")
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)


	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 30, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabMinuteAsteriskIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("*/20 * * * *")

	checkSliceValues(3, 0, 20, &result.minutes, t, "result.minute")
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 20, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

// Hour
func TestCreateCronCrontabHourRange(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* 11-15 * * *")

	checkEmptyMinute(result, t)
	checkSliceValues(5, 11, 1, &result.hours, t, "result.hour")
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 2, 11, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabHourRangeIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* 8-15/3 * * *")

	checkEmptyMinute(result, t)
	checkSliceValues(3, 8, 3, &result.hours, t, "result.hour")
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 2, 8, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabHourList(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* 10,12,14,16 * * *")

	checkEmptyMinute(result, t)
	checkSliceValues(4, 10, 2, &result.hours, t, "result.hour")
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 11, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabHourAsteriskIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* */3 * * *")

	checkEmptyMinute(result, t)
	checkSliceValues(8, 0, 3, &result.hours, t, "result.hour")
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 1, 18, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

// Day of month
func TestCreateCronCrontabDayOfMonthRange(t *testing.T) {
	var testTime = time.Date(2025, time.March, 16, 16, 10, 20, 30, time.UTC)
	SetLogValuesMockTime(&testTime)
	result := CreateCrontab("* * 5-15 * *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkSliceValues(11, 5, 1, &result.daysOfMonth, t, "result.dayOfMonth")
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.April, 5, 0, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabDayOfMonthRangeIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * 5-23/5 * *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkSliceValues(4, 5, 5, &result.daysOfMonth, t, "result.dayOfMonth")
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 5, 0, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabDayOfMonthList(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * 10,12,14,16 * *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkSliceValues(4, 10, 2, &result.daysOfMonth, t, "result.dayOfMonth")
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 10, 0, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabDayOfMonthAsteriskIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * */4 * *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkSliceValues(8, 1, 4, &result.daysOfMonth, t, "result.dayOfMonth")
	checkEmptyMonth(result, t)
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 11, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

// Month
func TestCreateCronCrontabMonthRange(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * 3-8 *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkSliceValues(6, 3, 1, &result.months, t, "result.month")
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 11, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabMonthRangeIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * 5-11/3 *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkSliceValues(3, 5, 3, &result.months, t, "result.month")
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabMonthList(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * 8,10,12 *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkSliceValues(3, 8, 2, &result.months, t, "result.month")
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

func TestCreateCronCrontabMonthAsteriskIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * */4 *")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkSliceValues(3, 1, 4, &result.months, t, "result.month")
	checkEmptyDayOfWeekh(result, t)

	testutil.AssertEquals(time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC), *result.NextTime, t, "init result.NextTime")
}

// Day of week
func TestCreateCronCrontabDayOfWeekRange(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * * 3-5")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkSliceValues(3, 3, 1, &result.daysOfWeek, t, "result.dayOfWeek")
}

func TestCreateCronCrontabDayOfWeekRangeIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * * 3-6/2")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkSliceValues(2, 3, 2, &result.daysOfWeek, t, "result.dayOfWeek")
}

func TestCreateCronCrontabDayOfWeekList(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * * 2,4,6")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkSliceValues(3, 2, 2, &result.daysOfWeek, t, "result.dayOfWeek")
}

func TestCreateCronCrontabDayOfWeekAsteriskIncrement(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("* * * * */3")

	checkEmptyMinute(result, t)
	checkEmptyHour(result, t)
	checkEmptyDayOfMonth(result, t)
	checkEmptyMonth(result, t)
	checkSliceValues(3, 0, 3, &result.daysOfWeek, t, "result.dayOfWeek")
}

// upper and lower bound
func TestCreateCronCrontabRangeUpperBound(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("1-61 1-25 2-32 2-13 1-7")

	checkSliceValues(59, 1, 1, &result.minutes, t, "result.minute")
	checkSliceValues(23, 1, 1, &result.hours, t, "result.hour")
	checkSliceValues(30, 2, 1, &result.daysOfMonth, t, "result.dayOfMonth")
	checkSliceValues(11, 2, 1, &result.months, t, "result.month")
	checkSliceValues(6, 1, 1, &result.daysOfWeek, t, "result.dayOfWeek")
}

func TestCreateCronCrontabRangeLowerBound(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("0-58 0-22 0-30 0-11 0-5")

	checkSliceValues(59, 0, 1, &result.minutes, t, "result.minute")
	checkSliceValues(23, 0, 1, &result.hours, t, "result.hour")
	checkSliceValues(30, 1, 1, &result.daysOfMonth, t, "result.dayOfMonth")
	checkSliceValues(11, 1, 1, &result.months, t, "result.month")
	checkSliceValues(6, 0, 1, &result.daysOfWeek, t, "result.dayOfWeek")
}

func TestCreateCronScheduleListUpperBound(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("59,60 23,24 31,32 12,13 6,7")

	checkSliceValues(1, 59, 1, &result.minutes, t, "result.minute")
	checkSliceValues(1, 23, 1, &result.hours, t, "result.hour")
	checkSliceValues(1, 31, 1, &result.daysOfMonth, t, "result.dayOfMonth")
	checkSliceValues(1, 12, 1, &result.months, t, "result.month")
	checkSliceValues(1, 6, 1, &result.daysOfWeek, t, "result.dayOfWeek")
}

func TestCreateCronCrontabListLowerBound(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	result := CreateCrontab("0,1 0,1 0,1 0,1 0,1")

	checkSliceValues(2, 0, 1, &result.minutes, t, "result.minute")
	checkSliceValues(2, 0, 1, &result.hours, t, "result.hour")
	checkSliceValues(1, 1, 1, &result.daysOfMonth, t, "result.dayOfMonth")
	checkSliceValues(1, 1, 1, &result.months, t, "result.month")
	checkSliceValues(2, 0, 1, &result.daysOfWeek, t, "result.dayOfWeek")
}

func TestCalculateNextTimeEveryMinute(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	crontab := CreateCrontab("* * * * *")

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 11, 0, 0, time.UTC), *crontab.NextTime, t, "init crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 12, 0, 0, time.UTC), *crontab.NextTime, t, "calculate crontab.NextTime")
}

func TestCalculateNextTimeEveryQuarterHour(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	crontab := CreateCrontab("*/15 * * * *")

	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 15, 0, 0, time.UTC), *crontab.NextTime, t, "init crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2025, time.March, 1, 16, 30, 0, 0, time.UTC), *crontab.NextTime, t, "calculate crontab.NextTime")
}

func TestCalculateNextTimeEveryHour(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	crontab := CreateCrontab("0 * * * *")

	testutil.AssertEquals(time.Date(2025, time.March, 1, 17, 0, 0, 0, time.UTC), *crontab.NextTime, t, "init crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2025, time.March, 1, 18, 0, 0, 0, time.UTC), *crontab.NextTime, t, "calculate crontab.NextTime")
}

func TestCalculateNextTimeEveryThirdHour(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	crontab := CreateCrontab("0 */3 * * *")

	testutil.AssertEquals(time.Date(2025, time.March, 1, 18, 0, 0, 0, time.UTC), *crontab.NextTime, t, "init crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2025, time.March, 1, 21, 0, 0, 0, time.UTC), *crontab.NextTime, t, "calculate crontab.NextTime")
}

func TestCalculateNextTimeEveryMonth(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	crontab := CreateCrontab("0 0 1 * *")

	testutil.AssertEquals(time.Date(2025, time.April, 1, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "init crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "calculate crontab.NextTime")
}

func TestCalculateNextTimeEverySecondMonth(t *testing.T) {
	SetLogValuesMockTime(&cronCrontabTestTime)
	crontab := CreateCrontab("0 0 1 */2 *")

	testutil.AssertEquals(time.Date(2025, time.May, 1, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "init crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2025, time.July, 1, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "calculate crontab.NextTime")
}

func TestCalculateNextTimeFebNonLeapyear(t *testing.T) {
	var testTime = time.Date(2025, time.February, 27, 16, 10, 20, 30, time.UTC)
	SetLogValuesMockTime(&testTime)
	crontab := CreateCrontab("0 0 2-31 * *")

	testutil.AssertEquals(time.Date(2025, time.February, 28, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "init crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2025, time.March, 2, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "calculate crontab.NextTime")
}

func TestCalculateNextTimeFebLeapyear(t *testing.T) {
	var testTime = time.Date(2024, time.February, 27, 16, 10, 20, 30, time.UTC)
	SetLogValuesMockTime(&testTime)
	crontab := CreateCrontab("0 0 2-31 * *")

	testutil.AssertEquals(time.Date(2024, time.February, 28, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "init crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "1. calculate crontab.NextTime")
	crontab.CalculateNextTime()
	testutil.AssertEquals(time.Date(2024, time.March, 2, 0, 0, 0, 0, time.UTC), *crontab.NextTime, t, "2. calculate crontab.NextTime")
}
