package domain

// TimeRangeType is an enumeration of time range types
type TimeRangeType int8

const (
	// TimeRangeTypeUnSpecified is an enumeration of unspecified time range type
	TimeRangeTypeUnSpecified TimeRangeType = iota

	// TimeRangeTypeOneDay is an enumeration of one week time range type(from sunday to saturday)
	TimeRangeTypeOneWeekDay

	// TimeRangeTypeOneWeek is an enumeration of one week time range type
	TimeRangeTypeOneWeek

	// TimeRangeTypeTwoWeeks is an enumeration of two weeks time range type
	TimeRangeTypeTwoWeeks

	// TimeRangeTypeOneMonth is an enumeration of one month time range type
	TimeRangeTypeOneMonth

	// TimeRangeTypeThreeMonths is an enumeration of three months time range type
	TimeRangeTypeThreeMonths

	// TimeRangeTypeSixMonths is an enumeration of six months time range type
	TimeRangeTypeSixMonths

	// TimeRangeTypeOneYear is an enumeration of one year time range type
	TimeRangeTypeOneYear
)

func (t TimeRangeType) IsValid() bool {
	switch t {
	case TimeRangeTypeOneWeekDay, TimeRangeTypeOneWeek, TimeRangeTypeTwoWeeks, TimeRangeTypeOneMonth, TimeRangeTypeThreeMonths, TimeRangeTypeSixMonths, TimeRangeTypeOneYear:
		return true
	}
	return false
}

func CvtToTimeRangeType(s string) TimeRangeType {
	switch s {
	case "one_week_day":
		return TimeRangeTypeOneWeekDay
	case "one_week":
		return TimeRangeTypeOneWeek
	case "two_weeks":
		return TimeRangeTypeTwoWeeks
	case "one_month":
		return TimeRangeTypeOneMonth
	case "three_months":
		return TimeRangeTypeThreeMonths
	case "six_months":
		return TimeRangeTypeSixMonths
	case "one_year":
		return TimeRangeTypeOneYear
	}
	return TimeRangeTypeUnSpecified
}
