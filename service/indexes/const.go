package indexes

const (
	defaultPartition = 0
)

const (
	sexField = iota
	statusField
	countryField
	cityField
	interestField
	birthYearField
	fnameField
	snameField
	premiumField
)

const (
	NullToken = iota
	NotNullToken
)

// Not used with Null/NotNull, can be started from 0
const (
	MaleToken = iota
	FemaleToken
)

// Not used with Null/NotNull, can be started from 0
const (
	SingleToken = iota
	InRelToken
	ComplToken
	NotSingleToken
	NotInRelToken
	NotComplToken
)

// Used with Null/NotNull, skip first four
const (
	PremiumNowToken = 4
)
