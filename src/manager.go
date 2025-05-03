package main

import (
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager interface {
	Start()
	Stop()
	Wait()
	Init(chan string, *iomap, *cimess)
}

type manager struct {
	signal_ chan string
	io_map_ *iomap
	rec_    *pgxpool.Config

	conn_pool_ *pgxpool.Pool

	workers  *mList[*worker]
	release_ chan int
	wkr_map_ *map[int]*mListNode[*worker]
}

func (tar *manager) Start() {
	// for id := range tar.workers {
	// 	conn, _ := tar.conn_pool_.Acquire(context.Background())
	// 	tar.workers[id].Init(id, conn, tar.imess_, tar.omess_, tar.pick_up_)
	// 	tar.workers[id].Start()
	// }

	// go func() {
	// 	for {
	// 		select {
	// 		case task := <-tar.omess_:
	// 			tar.opasser_ <- task
	// 		default:
	// 			time.Sleep(time.Millisecond * 100)
	// 		}
	// 	}
	// }()
	// for {
	// 	select {
	// 	case task := <-tar.ipasser_:
	// 		tar.imess_ <- task
	// 	default:
	// 		time.Sleep(time.Millisecond * 100)
	// 	}
	// }
	for {
		select {
		case IP := <-tar.signal_:
			now := tar.workers.Head()
			now.data.Change((*tar.io_map_.imp_)[IP],
				(*tar.io_map_.omp_)[IP])
		}
		select {
		case rls := <-tar.release_:
			tar.workers.Move_head((*tar.wkr_map_)[rls])
		}
	}
}

func (tar *manager) Stop() {

}

func (tar *manager) Wait() {

}

func (tar *manager) Init(signal chan string, iom *iomap, mess *cimess) {
	config, _ := pgxpool.ParseConfig(mess.String())

	cpuNum := int32(runtime.NumCPU())

	config.MaxConns = cpuNum * 4 // 动态调整
	config.MinConns = cpuNum * 2
	tar.io_map_ = iom
	tar.signal_ = signal
}
