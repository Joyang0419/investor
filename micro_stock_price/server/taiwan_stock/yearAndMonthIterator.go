package taiwan_stock

type YearAndMonthIterator struct {
	currentYear  uint32
	currentMonth uint32
	endYear      uint32
	endMonth     uint32
}

func NewYearAndMonthIterator(startYearAndMonth, endYearAndMonth YearAndMonth) *YearAndMonthIterator {
	return &YearAndMonthIterator{
		currentYear:  startYearAndMonth.Year,
		currentMonth: startYearAndMonth.Month - 1, // 看Next()的實作, 這裡要減1, 因為Next()會先加1
		endYear:      endYearAndMonth.Year,
		endMonth:     endYearAndMonth.Month,
	}
}

func (it *YearAndMonthIterator) Next() bool {
	if it.currentYear > it.endYear || (it.currentYear == it.endYear && it.currentMonth > it.endMonth) {
		return false
	}

	it.currentMonth++
	if it.currentMonth > 12 {
		it.currentYear++
		it.currentMonth = 1
	}

	return true
}

func (it *YearAndMonthIterator) Current() YearAndMonth {
	return YearAndMonth{
		Year:  it.currentYear,
		Month: it.currentMonth,
	}
}
