package repo

type Like struct {
	ID uint32
	TS uint32
}

type StatusEnum uint8

const (
	_          = iota
	FreeStatus = StatusEnum(iota)
	BusyStatus
	ComplicatedStatus
)

type PeriodEnum uint8

const (
	_           = iota
	MonthPeriod = PeriodEnum(iota)
	QuarterPeriod
	HalfYearPeriod
)

type Account struct {
	Email         string
	Birth         uint32
	Joined        uint32
	Fname         string
	Sname         string
	Country       string
	City          string
	Status        StatusEnum
	PremiumFinish uint32
	PremiumPeriod PeriodEnum
	Interests     []string
	LikesTo       []Like
}
