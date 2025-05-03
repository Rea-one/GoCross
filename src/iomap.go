package main

type IOMap interface {
	Init()
	Register(string)
	Erase(string)
	GetIn(string)
	GetOut(string)
}

type iomap struct {
	imp_ *map[string]chan task
	omp_ *map[string]chan task
}

func (tar *iomap) Init() {
	tar.imp_ = &map[string]chan task{}
	tar.omp_ = &map[string]chan task{}
	tar.Register("default")
}

func (tar *iomap) Register(id string) {
	if _, ok := (*tar.imp_)[id]; !ok {
		(*tar.imp_)[id] = make(chan task)
		(*tar.omp_)[id] = make(chan task)
	}
}

func (tar *iomap) Erase(id string) {
	delete(*tar.imp_, id)
	delete(*tar.omp_, id)
}

func (tar *iomap) GetIn(id string) chan task {
	return (*tar.imp_)[id]
}

func (tar *iomap) GetOut(id string) chan task {
	return (*tar.omp_)[id]
}
