package gocross

import (
	"context"
	"now/sqlmap"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Worker interface {
	Start()
	Stop()
	Wait()
	Action(*sqlmap.Task)
	Change(string)
	Init(int, string, *iomap, chan int, *pgxpool.Conn)
	serve()
	to_db(*sqlmap.Task)
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
}

func processRows(rows *pgx.Rows) string {
	var result string
	for (*rows).Next() {
		var row string
		(*rows).Scan(row)
		result += row
	}
	return result
}

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
			if task.State == "nomore" {
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

func (tar *worker) to_db(task *sqlmap.Task) {
	rows, err := tar.link_.Query(context.Background(), task.Query)
	if err != nil {
		task.State = err.Error()
	} else {
		task.SQL_fb = processRows(&rows)
	}
	rows.Close()
}

func (tar *worker) act_task(task *sqlmap.Task) {
	tar.ana_.Ana(task)
}

func (tar *worker) act_register(task *sqlmap.Task) {
	task.Query = "insert into user(username,password)" +
		"values('" + task.Sender + "','" + task.Password + "')"
	tar.to_db(task)
	task.Query = ""
	task.Password = ""
	if task.State == "" {
		task.State = "success"
		tar.checker_.Link(tar.index_, task.Sender)
	}
}

func (tar *worker) act_login(task *sqlmap.Task) {
	task.Query = "select * from user where username='" + task.Sender + "'" +
		"and password='" + task.Password + "'"
	tar.to_db(task)
	task.Query = ""
	task.Password = ""
	if task.State == "" {
		task.State = "success"
		tar.online_ = true
		tar.checker_.Link(tar.index_, task.Sender)
	}
}

func (tar *worker) act_pass(task *sqlmap.Task) {
	if tar.online_ {
		task.Query = "insert * from user where username='" + task.Sender + "'" +
			""
		tar.to_db(task)
		if task.State == "" {
			task.State = "success"
			tar.to_put_ = tar.checker_.GetOut(task.Receiver)
			task.Feedback = "s " + task.Sender +
				" r " + task.Receiver + " m " + task.Message
		}
	}
}
