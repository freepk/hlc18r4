package db

//import (
//	"io"
//)

/*
func (a *account) reset() {
	a.loginSize = 0
	a.domain = 0
}

func (a *account) setEmail(b []byte) error {
	login, domain, ok := splitEmail(b)
	if !ok {
		return DefaultError
	}
	loginSize := len(login)
	if loginSize > loginMaxSize {
		return DefaultError
	}
	a.loginSize = uint8(loginSize)
	for i := 0; i < loginSize; i++ {
		a.login[i] = login[i]
	}
	a.domain = uint8(domainLookup.GetKeyOrSet(domain))
	return nil
}

func (a *account) putEmail(w io.Writer) error {
	return nil
}
*/
