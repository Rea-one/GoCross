package gocross

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager interface {
	Start()
	Stop()
	Wait()
	serve()
	Init(chan string, *iomap, *cimess)
}

type manager struct {
	signal_ chan string
	io_map_ *iomap

	conn_pool_ *pgxpool.Pool
	stop_      bool

	wkr_num_ int
	workers  *mList[*worker]
	release_ chan int
	wkr_map_ map[int]*mListNode[*worker]
}

func (tar *manager) Start() {
	go tar.serve()
	log.Printf("manager 开始运行")
}

func (tar *manager) Stop() {

}

func (tar *manager) Wait() {

}

func (tar *manager) Init(signal chan string, iom *iomap, mess *cimess) {
	config, _ := pgxpool.ParseConfig(mess.String())

	config.MinConns = 4
	config.MaxConns = 8

	tar.io_map_ = iom
	tar.signal_ = signal
	tar.conn_pool_, _ = pgxpool.NewWithConfig(context.Background(), config)
	tar.wkr_num_ = 8
	tar.release_ = make(chan int, tar.wkr_num_)
	tar.wkr_map_ = make(map[int]*mListNode[*worker])
	for count := range tar.wkr_num_ {
		conn, _ := tar.conn_pool_.Acquire(context.Background())
		if tar.workers == nil {
			tar.workers = new(mList[*worker])

		}
		tar.workers.Push_tail(&mListNode[*worker]{})
		tar.workers.Tail().Init(new(worker))
		tar.workers.Tail().data.Init(count, conn,
			tar.io_map_.GetIn("default"), tar.io_map_.GetOut("default"), tar.release_)
		tar.wkr_map_[count] = tar.workers.Tail()
	}
	log.Printf("manager 初始化成功")
}

func (tar *manager) serve() {
	tickerFast := time.NewTicker(10 * time.Millisecond)  // 高频轮询间隔
	tickerSlow := time.NewTicker(100 * time.Millisecond) // 低频等待间隔
	defer tickerFast.Stop()
	defer tickerSlow.Stop()

	useFastPolling := false

	for !tar.stop_ {
		if useFastPolling {
			select {
			case IP := <-tar.signal_:
				now := tar.workers.Head()
				now.data.Change((*tar.io_map_.imp_)[IP], (*tar.io_map_.omp_)[IP])
				now.Get().Start()
				tar.workers.Move_tail(now)
				tickerFast.Reset(10 * time.Millisecond) // 保持高频
			case rls := <-tar.release_:
				tar.workers.Move_head(tar.wkr_map_[rls])
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
				tar.workers.Move_head(tar.wkr_map_[rls])
			case <-tickerSlow.C:
				// 低频等待，防止空转
			}
		}
	}
}
