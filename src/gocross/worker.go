package gocross

import (
	"context"
	"now/sqlmap"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker interface {
	Start()
	Stop()
	Wait()
	Action(*sqlmap.Task)
	Change(string)
	Init(int, string, *iomap, chan int, *pgxpool.Conn)
	registerd(*sqlmap.Task) bool
	serve()
	act_task(*sqlmap.Task)
	act_register(*sqlmap.Task)
	act_login(*sqlmap.Task)
	act_pass(*sqlmap.Task)
}

type worker struct {
	id_       int
	stop_     bool
	checker_  *Checker
	index_    string
	input_    <-chan sqlmap.Task
	back_put_ chan<- sqlmap.Task
	to_put_   chan<- sqlmap.Task
	link_     *pgxpool.Conn
	release_  chan<- int
	online_   bool
	ana_      sqlmap.SqlMap
	counter   int
}

func (tar *worker) Start() {
	go tar.serve()
}

func (tar *worker) Stop() {
	tar.stop_ = true
	tar.release_ <- tar.id_
}

func (tar *worker) Wait() {

}

func (tar *worker) Init(id int, index string, checker *Checker,
	rsl chan int, link *pgxpool.Conn) {
	tar.id_ = id
	tar.stop_ = false
	tar.link_ = link
	tar.index_ = index
	tar.checker_ = checker
	tar.input_ = tar.checker_.getin(tar.index_)
	tar.back_put_ = tar.checker_.getout(tar.index_)
	tar.to_put_ = nil
	tar.online_ = false
	tar.release_ = rsl
	tar.ana_ = sqlmap.SqlMap{}
	tar.ana_.Init()
	tar.counter = 0
}

// func processRows(rows *pgx.Rows) string {
// 	var result string
// 	for (*rows).Next() {
// 		var row string
// 		(*rows).Scan(row)
// 		result += row
// 	}
// 	return result
// }

func (tar *worker) Change(index string) {
	tar.index_ = index
	tar.input_ = tar.checker_.getin(tar.index_)
	tar.back_put_ = tar.checker_.getout(tar.index_)
	tar.to_put_ = nil
}

func (tar *worker) serve() {
	for !tar.stop_ {
		select {
		case task := <-tar.input_:
			tar.counter++
			if task.State == "nomore" {
				task.Feedback = "nomore"
				tar.Stop()
			} else {
				tar.act_task(&task)
				switch task.Ttype {
				case "register":
					tar.act_register(&task)
				case "login":
					tar.act_login(&task)
				case "pass":
					tar.act_pass(&task)
				case "1":
				case "2":
				case "3":
				case "4":
				case "5":
				case "6":
				case "7":
				case "8":
				case "9":
				}
			}
			tar.back_put_ <- task
			if tar.to_put_ != nil {
				tar.to_put_ <- task
			}
		case <-time.After(time.Second * 5):
			continue
		}
	}
}

// func (tar *worker) to_db(task *sqlmap.Task) {
// 	rows, err := tar.link_.Query(context.Background(), task.Query)
// 	if err != nil {
// 		task.State = err.Error()
// 	} else {
// 		task.SQL_fb = processRows(&rows)
// 	}
// 	rows.Close()
// }

func (tar *worker) act_task(task *sqlmap.Task) {
	tar.ana_.Ana(task)
}

func (tar *worker) act_register(task *sqlmap.Task) {
	if tar.registerd(task) {
		task.State = "rejected"
		task.Feedback = "another using"
	} else {
		task.Query = "insert into reg(id,password)" +
			"values($1, $2)"
		fb, err := tar.link_.Exec(context.Background(),
			task.Query, task.Sender, task.Password)
		if err != nil {
			task.State = err.Error()
		} else if fb.RowsAffected() > 0 {
			task.State = "registered"
			task.Feedback = "registration success"
		}
	}
	task.Query = ""
	task.Password = ""
}

func (tar *worker) act_login(task *sqlmap.Task) {
	task.Query = "select * from reg where id=$1 and password=$2"
	fb, err := tar.link_.Exec(context.Background(),
		task.Query, task.Sender, task.Password)
	if err != nil {
		task.State = err.Error()
	} else if fb.RowsAffected() > 0 {
		task.State = "logined"
		task.Feedback = "login success"
		tar.online_ = true
		tar.checker_.Link(tar.index_, task.Sender)
	} else {
		task.State = "rejected"
		task.Feedback = "wrong password"
		tar.online_ = false
	}
	task.Query = ""
	task.Password = ""
}

func (tar *worker) act_pass(task *sqlmap.Task) {
	if tar.online_ {
		task.Query = "insert into message(id,receiver,message)" +
			"values($1, $2, $3)"
		fb, err := tar.link_.Exec(context.Background(),
			task.Query, task.Sender, task.Receiver, task.Message)
		if err != nil {
			task.State = err.Error()
		} else if fb.RowsAffected() > 0 {
			task.State = "passed"
			task.Feedback = "pass success"
			tar.to_put_ = tar.checker_.GetOut(task.Receiver)
		}
	} else {
		task.State = "rejected"
		task.Feedback = "not logined"
	}
}

func (tar *worker) registerd(task *sqlmap.Task) bool {
	query := "select * from reg where id=$1"
	rows, err := tar.link_.Exec(context.Background(),
		query, task.Sender)
	if err != nil {
		task.State = err.Error()
	}
	return rows.RowsAffected() > 0
}
