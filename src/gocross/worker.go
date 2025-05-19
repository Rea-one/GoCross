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
	act_config(*sqlmap.Task)
	act_sync(*sqlmap.Task)
	act_addfriend(*sqlmap.Task)
	respond(*sqlmap.Task)
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
			tar.act_task(&task)
			switch task.Ttype {
			case "nomore":
				tar.Stop()
			case "register":
				tar.act_register(&task)
			case "login":
				tar.act_login(&task)
			case "pass":
				tar.act_pass(&task)
			case "config":
				tar.act_config(&task)
			case "addfriend":
				tar.act_addfriend(&task)
			case "sync":
				tar.act_sync(&task)
			default:
				task.State = "rejected"
				task.Feedback = "wrong type"
			}
		case <-time.After(time.Second * 5):
			continue
		}
	}
}

func (tar *worker) registerd(task *sqlmap.Task) bool {
	query := "select * from reg where id=$1"
	rows, err := tar.link_.Exec(context.Background(),
		query, task.At)
	if err != nil {
		task.State = err.Error()
	}
	return rows.RowsAffected() > 0
}

func (tar *worker) act_task(task *sqlmap.Task) {
	tar.ana_.Ana(task)
}

func (tar *worker) act_register(task *sqlmap.Task) {
	if tar.registerd(task) {
		task.State = "rejected"
		task.Feedback = "another using"
	} else {
		query := "insert into reg(id,password)" +
			"values($1, $2)"
		fb, err := tar.link_.Exec(context.Background(),
			query, task.At, task.Password)
		if err != nil {
			task.State = err.Error()
		} else if fb.RowsAffected() > 0 {
			task.State = "registration success"
		}
	}
	task.Message = ""
	task.Password = ""
	tar.respond(task)
}

func (tar *worker) act_login(task *sqlmap.Task) {
	query := "select * from reg where id=$1 and password=$2"
	fb, err := tar.link_.Exec(context.Background(),
		query, task.At, task.Password)
	if err != nil {
		task.State = err.Error()
	} else if fb.RowsAffected() > 0 {
		task.State = "login success"
		tar.online_ = true
		tar.checker_.Link(tar.index_, task.At) // ID热映射
	} else {
		task.State = "rejected"
		task.Feedback = "wrong password or wrong id"
		tar.online_ = false
	}
	task.Message = ""
	task.Password = ""
	tar.respond(task)
}

func (tar *worker) act_pass(task *sqlmap.Task) {
	if !tar.online_ {
		task.State = "rejected"
		task.Feedback = "not logined"
		return
	}
	query := "insert into message(id,receiver,message)" +
		"values($1, $2, $3)"
	fb, err := tar.link_.Exec(context.Background(),
		query, task.Sender, task.Receiver, task.Message)
	if err != nil {
		task.State = err.Error()
	} else if fb.RowsAffected() > 0 {
		task.State = "pass success"
		tar.to_put_ = tar.checker_.GetOut(task.Receiver)
	}
	task.Password = ""
	tar.respond(task)
}

func (tar *worker) act_config(task *sqlmap.Task) {
	if !tar.online_ {
		task.State = "rejected"
		task.Feedback = "not logined"
		return
	}

	query := `
		INSERT INTO config(id, name, image)
		VALUES ($1, $2, $3)
		ON CONFLICT (id)
		DO UPDATE SET name = $2, image = $3`

	fb, err := tar.link_.Exec(context.Background(),
		query, task.At, task.Sender, task.ImageURL)

	if err != nil {
		task.State = err.Error()
	} else if fb.RowsAffected() > 0 {
		task.State = "config success"
		tar.to_put_ = tar.checker_.GetOut(task.Receiver)
	}
	task.Password = ""
	tar.respond(task)
}
func (tar *worker) act_addfriend(task *sqlmap.Task) {
	if !tar.online_ {
		task.State = "rejected"
		task.Feedback = "not logined"
		return
	}

	// 防止重复添加
	query := "SELECT * FROM friends WHERE user_id = $1 AND friend_id = $2"
	rows, err := tar.link_.Exec(context.Background(), query, task.Sender, task.Receiver)
	if err != nil {
		task.State = err.Error()
		return
	}
	if rows.RowsAffected() > 0 {
		task.State = "rejected"
		task.Feedback = "already added"
		return
	}

	// 插入好友关系（双向）
	query = `INSERT INTO friends(user_id, friend_id, signature, regnature)
		VALUES ($1, $2, true, false), ($2, $1, false, true)`
	fb, err := tar.link_.Exec(context.Background(), query, task.Sender, task.Receiver)
	if err != nil {
		task.State = err.Error()
	} else if fb.RowsAffected() > 0 {
		task.State = "add friend success"
		tar.to_put_ = tar.checker_.GetOut(task.Receiver)
	}
	task.Password = ""
	tar.respond(task)
}
func (tar *worker) act_sync(task *sqlmap.Task) {

}

func (tar *worker) respond(task *sqlmap.Task) {
	tar.back_put_ <- *task
	if tar.to_put_ != nil {
		tar.to_put_ <- *task
	}
}
