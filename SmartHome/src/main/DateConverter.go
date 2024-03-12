package main

var nonLeapYearMonths = [12]int{
	31, //Jan
	28, //Feb
	31, //Mar
	30, //Apr
	31, //May
	30, //Jun
	31, //Jul
	31, //Aug
	30, //Sep
	31, //Oct
	30, //Nov
	31, //Dec
}

var nonLeapYearMonthsTotals = [12]int{
	0,   //Jan
	31,  //Feb
	59,  //Mar
	90,  //Apr
	120, //May
	151, //Jun
	181, //Jul
	212, //Aug
	243, //Sep
	273, //Oct
	304, //Nov
	334, //Dec
}

var leapYearMonths = [12]int{
	31, //Jan
	29, //Feb
	31, //Mar
	30, //Apr
	31, //May
	30, //Jun
	31, //Jul
	31, //Aug
	30, //Sep
	31, //Oct
	30, //Nov
	31, //Dec
}

var leapYearMonthsTotals = [12]int{
	0,   //Jan
	31,  //Feb
	60,  //Mar
	91,  //Apr
	121, //May
	152, //Jun
	182, //Jul
	213, //Aug
	244, //Sep
	274, //Oct
	305, //Nov
	335, //Dec
}

type yearData struct {
	months        [12]int
	monthsTotals  [12]int
	days          int
	daysFromStart int
}

func (year *yearData) calculateMonthAndDay(y, index int) int {
	daysFromStart := index - year.daysFromStart
	month := daysFromStart / 30
	if month > 11 {
		month = 11
	}
	for {
		daysFromStartMonth := year.monthsTotals[month]
		if daysFromStart < daysFromStartMonth {
			month -= 1
		} else if daysFromStart >= daysFromStartMonth+year.months[month] {
			month += 1
		} else {
			day := daysFromStart - daysFromStartMonth + 1
			return y*10000 + (month+1)*100 + day
		}
	}
}

type dateConverter struct {
	minYear int
	years   []yearData
}

func newDateConverter(minYear, minMonth, yearsToCreate int) dateConverter {
	return dateConverter{
		minYear: minYear,
		years:   buildYears(minYear, minMonth, yearsToCreate),
	}
}

func (c *dateConverter) toDate(index int) int {
	y := index / 365
	l := len(c.years)
	if y >= l {
		y = l - 1
	}
	for {
		year := c.years[y]
		if index < year.daysFromStart {
			y -= 1
		} else if index > year.daysFromStart+year.days {
			y += 1
		} else {
			return year.calculateMonthAndDay(y+c.minYear, index)
		}
	}
}

func (c *dateConverter) fromDate(date int) int {
	year := date / 10000
	month := (date / 100) % 100
	day := date % 100
	y := c.years[year-c.minYear]
	return y.daysFromStart + y.monthsTotals[month-1] + day - 1
}

func buildYears(minYear, minMonth, yearsToCreate int) []yearData {
	var result []yearData
	var daysFromStart int
	year := minYear
	for yearsToCreate > 0 {
		var y yearData
		if (year % 4) == 0 {
			if year == minYear {
				daysFromStart = -leapYearMonthsTotals[minMonth-1]
			}
			y = yearData{
				months:        leapYearMonths,
				monthsTotals:  leapYearMonthsTotals,
				days:          366,
				daysFromStart: daysFromStart,
			}
			daysFromStart += 366
		} else {
			if year == minYear {
				daysFromStart = -nonLeapYearMonthsTotals[minMonth-1]
			}
			y = yearData{
				months:        nonLeapYearMonths,
				monthsTotals:  nonLeapYearMonthsTotals,
				days:          365,
				daysFromStart: daysFromStart,
			}
			daysFromStart += 365
		}
		result = append(result, y)
		yearsToCreate--
	}
	return result
}
