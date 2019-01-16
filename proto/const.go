package proto

const (
	BirthKey             = `"birth":`
	BirthLen             = len(BirthKey)
	CityKey              = `"city":`
	CityLen              = len(CityKey)
	CountryKey           = `"country":`
	CountryLen           = len(CountryKey)
	EmailKey             = `"email":`
	EmailLen             = len(EmailKey)
	IdKey                = `"id":`
	IdLen                = len(IdKey)
	JoinedKey            = `"joined":`
	JoinedLen            = len(JoinedKey)
	FnameKey             = `"fname":`
	FnameLen             = len(FnameKey)
	InterestsKey         = `"interests":`
	InterestsLen         = len(InterestsKey)
	LikesKey             = `"likes":`
	LikesLen             = len(LikesKey)
	PhoneKey             = `"phone":`
	PhoneLen             = len(PhoneKey)
	PremiumKey           = `"premium":`
	PremiumLen           = len(PremiumKey)
	SexKey               = `"sex":`
	SexLen               = len(SexKey)
	SnameKey             = `"sname":`
	SnameLen             = len(SnameKey)
	StatusKey            = `"status":`
	StatusLen            = len(StatusKey)
	TsKey                = `"ts":`
	TsLen                = len(TsKey)
	StartKey             = `"start":`
	StartLen             = len(StartKey)
	FinishKey            = `"finish":`
	FinishLen            = len(FinishKey)
	BusyStatusStr        = `"\u0437\u0430\u043d\u044f\u0442\u044b"`
	BusyStatusLen        = len(BusyStatusStr)
	FreeStatusStr        = `"\u0441\u0432\u043e\u0431\u043e\u0434\u043d\u044b"`
	FreeStatusLen        = len(FreeStatusStr)
	ComplicatedStatusStr = `"\u0432\u0441\u0451 \u0441\u043b\u043e\u0436\u043d\u043e"`
	ComplicatedStatusLen = len(ComplicatedStatusStr)
)

const (
	MaleSexStr   = `"m"`
	FemaleSexStr = `"f"`
)

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

type SexEnum byte

const (
	_       = iota
	MaleSex = SexEnum(iota)
	FemaleSex
)
