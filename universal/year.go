package universal

import "time"

type Year struct {
	Year int
}

func NewYear(year int) Year {
	return Year{Year: year}
}

func CurrentYear() Year {
	return Year{Year: time.Now().Year()}
}

func (y Year) From(loc *time.Location) time.Time {
	if loc == nil {
		loc = time.Local
	}
	return time.Date(y.Year, 1, 1, 0, 0, 0, 0, loc)
}

func (y Year) To(loc *time.Location) time.Time {
	if loc == nil {
		loc = time.Local
	}
	return time.Date(y.Year, 12, 31, 0, 0, 0, 0, loc)
}

func (y Year) FromString() string {
	return y.From(time.Local).Format("2006-01-02")
}

func (y Year) ToString() string {
	return y.To(time.Local).Format("2006-01-02")
}
