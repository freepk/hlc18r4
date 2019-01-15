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
	Fname         uint8
	Sname         uint16
	Country       uint8
	City          uint16
	Status        StatusEnum
	PremiumFinish uint32
	PremiumPeriod PeriodEnum
	Interests     []uint8
	LikesTo       []Like
}
