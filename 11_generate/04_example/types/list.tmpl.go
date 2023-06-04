package types

type DUMMYTYPEList []*DUMMYTYPE

type DUMMYTYPEToBool func(*DUMMYTYPE) bool

func (al DUMMYTYPEList) Filter(f DUMMYTYPEToBool) DUMMYTYPEList {
	var ret DUMMYTYPEList
	for _, a := range al {
		if f(a) {
			ret = append(ret, a)
		}
	}
	return ret
}
