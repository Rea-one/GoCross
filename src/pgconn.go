package main

type CIMess interface {
	String()
}

type cimess struct {
	dominator_ string
	password_  string
	host_      string
	database_  string
}

func (tar *cimess) String() string {
	return "postgresql://" + tar.dominator_ + ":" + tar.password_ +
		"@" + tar.host_ + "/" + tar.database_
}
