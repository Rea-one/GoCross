package main

import (
	"context"
	"runtime"
	"time"

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

	wkr_num_ int
	workers  *mList[*worker]
	release_ chan int
	wkr_map_ *map[int]*mListNode[*worker]
}

func (tar *manager) Start() {
	tickerFast := time.NewTicker(10 * time.Millisecond)  // 高频轮询间隔
	tickerSlow := time.NewTicker(100 * time.Millisecond) // 低频等待间隔
	defer tickerFast.Stop()
	defer tickerSlow.Stop()

	useFastPolling := false

	for {
		if useFastPolling {
			select {
			case IP := <-tar.signal_:
				now := tar.workers.Head()
				now.data.Change((*tar.io_map_.imp_)[IP], (*tar.io_map_.omp_)[IP])
				now.Get().Start()
				tar.workers.Move_tail(now)
				tickerFast.Reset(10 * time.Millisecond) // 保持高频
			case rls := <-tar.release_:
				tar.workers.Move_head((*tar.wkr_map_)[rls])
			case <-tickerFast.C:
				// 继续高频轮询
			}
		} else {
			select {
			case IP := <-tar.signal_:
				useFastPolling = true
				now := tar.workers.Head()
				now.data.Change((*tar.io_map_.imp_)[IP], (*tar.io_map_.omp_)[IP])
				now.Get().Start()
				tar.workers.Move_tail(now)
			case rls := <-tar.release_:
				tar.workers.Move_head((*tar.wkr_map_)[rls])
			case <-tickerSlow.C:
				// 低频等待，防止空转
			}
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
	tar.conn_pool_, _ = pgxpool.NewWithConfig(context.Background(), config)
	tar.wkr_num_ = 8
	tar.workers = new(mList[*worker])
	tar.workers.Init_with_num(tar.wkr_num_)
}
