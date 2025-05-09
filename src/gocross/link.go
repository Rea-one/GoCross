package gocross

type CIMess interface {
	PGString()
	MNString()
}

type cimess struct {
	dominator_ string
	password_  string
	database_  string
	pg_host_   string
	mn_host_   string
	host_      string
}

func (tar *cimess) PGString() string {
	return "postgresql://" + tar.dominator_ + ":" + tar.password_ +
		"@" + tar.pg_host_ + "/" + tar.database_
}
