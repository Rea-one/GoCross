package gocross

import "now/sqlmap"

type IOMap interface {
	Init()
	Register(string)
	Erase(string)
	GetIn(string)
	GetOut(string)
}

type ActChecker interface {
	Init()
	Link(string, string)
	Register(string)
	Erase(string)
	GetIn(string)
	getin(string)
	GetOut(string)
	getout(string)
}
type iomap struct {
	imp_ *map[string]chan sqlmap.Task
	omp_ *map[string]chan sqlmap.Task
}

type Checker struct {
	iom_   *iomap
	token_ map[string]string
}

func (tar *iomap) Init() {
	tar.imp_ = new(map[string]chan sqlmap.Task)
	tar.omp_ = new(map[string]chan sqlmap.Task)
	tar.Register("default")
}

func (tar *iomap) Register(id string) {
	if _, ok := (*tar.imp_)[id]; !ok {
		(*tar.imp_)[id] = make(chan sqlmap.Task)
		(*tar.omp_)[id] = make(chan sqlmap.Task)
	}
}

func (tar *iomap) Erase(id string) {
	delete(*tar.imp_, id)
	delete(*tar.omp_, id)
}

func (tar *iomap) GetIn(id string) chan sqlmap.Task {
	return (*tar.imp_)[id]
}

func (tar *iomap) GetOut(id string) chan sqlmap.Task {
	return (*tar.omp_)[id]
}

func (tar *Checker) Init() {
	tar.iom_ = new(iomap)
	tar.iom_.Init()
	tar.token_ = make(map[string]string)
}

func (tar *Checker) Link(id string, token string) {
	tar.iom_.Register(id)
	tar.token_[token] = id
}

func (tar *Checker) Register(token string) {
	tar.iom_.Register(token)
}

func (tar *Checker) Erase(token string) {
	id, OK := tar.token_[token]
	if OK {
		tar.iom_.Erase(id)
		delete(tar.token_, token)
	}
}

func (tar *Checker) GetIn(token string) chan sqlmap.Task {
	id, OK := tar.token_[token]
	if OK {
		return tar.iom_.GetIn(id)
	}
	return nil
}

func (tar *Checker) getin(id string) chan sqlmap.Task {
	return tar.iom_.GetIn(id)
}

func (tar *Checker) GetOut(token string) chan sqlmap.Task {
	id, OK := tar.token_[token]
	if OK {
		return tar.iom_.GetOut(id)
	}
	return nil
}

func (tar *Checker) getout(id string) chan sqlmap.Task {
	return tar.iom_.GetOut(id)
}
