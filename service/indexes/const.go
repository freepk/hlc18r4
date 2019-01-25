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
)

const (
	NullToken = iota
	NotNullToken
)

const (
	MaleToken = iota
	FemaleToken
)

const (
	SingleToken = iota
	InRelToken
	ComplToken
	NotSingleToken
	NotInRelToken
	NotComplToken
)
