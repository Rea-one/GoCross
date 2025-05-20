package gocross

import (
	"context"
	"now/sqlmap"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// 缺陷： 	没有使用ORM，容易sql注入 很多方法未验证
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
	act_addFriend(*sqlmap.Task)
	act_resAddFriend(*sqlmap.Task)
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
	// 输入通道
	tar.input_ = tar.checker_.getin(tar.index_)
	// 返回通道，默认与输入通道相同
	tar.back_put_ = tar.checker_.getout(tar.index_)
	// 传递通道，默认为空
	tar.to_put_ = nil
	// 在线状态
	tar.online_ = false
	tar.release_ = rsl
	// sql语句表
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
	// 与init()的机制类似
	tar.input_ = tar.checker_.getin(tar.index_)
	tar.back_put_ = tar.checker_.getout(tar.index_)
	tar.to_put_ = nil
}

func (tar *worker) serve() {
	for !tar.stop_ {
		select {
		// 任务处理
		case task := <-tar.input_:
			tar.counter++

			// 转化语句
			tar.act_task(&task)

			// 根据类型执行对应操作
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
			case "add friend":
				tar.act_addFriend(&task)
			case "response add friend":
				tar.act_resAddFriend(&task)
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
		task.Feedback = "已有登录"
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
	tar.back_put_ <- *task
}

func (tar *worker) act_login(task *sqlmap.Task) {
	query := `
		SELECT *
		FROM reg
		WHERE id = $1
		AND password = $2`
	fb, err := tar.link_.Exec(context.Background(), query, task.At, task.Password)
	if err != nil {
		task.State = "error"
		task.Feedback = err.Error()
	} else if fb.RowsAffected() > 0 {
		task.State = "login_success"
		task.Feedback = "登录成功"
		tar.online_ = true
		tar.checker_.Link(tar.index_, task.At) // ID热映射
	} else {
		task.State = "rejected"
		task.Feedback = "密码错误或ID不存在"
	}
	task.Message = ""
	task.Password = ""
	tar.back_put_ <- *task
}

func (tar *worker) act_pass(task *sqlmap.Task) {
	if !tar.online_ {
		task.State = "rejected"
		task.Feedback = "未登录"
		tar.back_put_ <- *task
		return
	}

	query := `
		INSERT INTO
		message(id, receiver, message)
		VALUES ($1, $2, $3)`
	fb, err := tar.link_.Exec(context.Background(), query, task.Sender, task.Receiver, task.Message)
	if err != nil {
		task.State = "error"
		task.Feedback = err.Error()
		tar.back_put_ <- *task
	} else if fb.RowsAffected() > 0 {
		task.State = "pass_success"
		task.Feedback = "消息发送成功"
		tar.back_put_ <- *task
		task.Feedback = "待接收消息"
		tar.to_put_ <- *task
		tar.to_put_ = tar.checker_.GetOut(task.Receiver)
	}
}

func (tar *worker) act_config(task *sqlmap.Task) {
	if !tar.online_ {
		task.State = "rejected"
		task.Feedback = "未登录"
		tar.back_put_ <- *task
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
	tar.back_put_ <- *task
}

func (tar *worker) act_addFriend(task *sqlmap.Task) {
	if !tar.online_ {
		task.State = "rejected"
		task.Feedback = "未登录"
		tar.back_put_ <- *task
		return
	}

	// 防止重复添加
	var query = `
		SELECT *
		FROM friends
		WHERE user_id = $1
		AND friend_id = $2`

	rows, err := tar.link_.Exec(context.Background(), query, task.Sender, task.Receiver)

	if err != nil {
		task.State = "error"
		task.Feedback = "数据库错误：" + err.Error()
		tar.back_put_ <- *task
		return
	}

	if rows.RowsAffected() > 0 {
		task.State = "already"
		task.Feedback = "已经是好友"
		tar.back_put_ <- *task
		return
	}

	// 插入双向好友关系
	query = `
		INSERT INTO friends(user_id, friend_id, signature, regnature)
		VALUES ($1, $2, $3, $4), ($2, $1, $4, $3)`
	fb, err := tar.link_.Exec(context.Background(), query, task.Sender, task.Receiver, "agree", "null")
	if err != nil {
		task.State = "error"
		task.Feedback = "数据库错误：" + err.Error()
		tar.back_put_ <- *task
		return
	}

	if fb.RowsAffected() > 0 {
		task.State = "friend_adding"
		task.Feedback = "好友请求已发送"
		tar.to_put_ <- *task
		task.Feedback = "待接收好友请求"
		tar.to_put_ <- *task
		tar.to_put_ = tar.checker_.GetOut(task.Receiver)
	}
}

func (tar *worker) act_resAddFriend(task *sqlmap.Task) {
	if !tar.online_ {
		task.State = "rejected"
		task.Feedback = "请先登录以执行此操作"
		tar.back_put_ <- *task
		return
	}

	// 查询是否有待确认的好友请求
	query := `
		SELECT *
		FROM friends
		WHERE user_id = $1
		AND friend_id = $2
		AND signature IS false`
	rows, err := tar.link_.Exec(context.Background(), query, task.Receiver, task.Sender)
	if err != nil {
		task.State = "database error"
		task.Feedback = err.Error()
		tar.back_put_ <- *task
		return
	}
	if rows.RowsAffected() == 0 {
		task.State = "rejected"
		task.Feedback = "没有待确认的好友请求"
		tar.back_put_ <- *task
		return
	}

	switch task.Message {
	case "accept":
		updateQuery := `
			UPDATE friends
			SET signature = true, regenature = true
			WHERE user_id = $1 AND friend_id = $2`
		_, err = tar.link_.Exec(context.Background(), updateQuery, task.Receiver, task.Sender)
		if err != nil {
			task.State = "database error"
			task.Feedback = err.Error()
		} else {
			task.State = "accepted"
			task.Feedback = "已接受好友请求"
			tar.back_put_ <- *task
			task.Feedback = "已添加为好友"
			tar.to_put_ <- *task
			tar.to_put_ = tar.checker_.GetOut(task.Receiver)

		}

	case "reject":
		deleteQuery := `
			DELETE
			FROM friends
			WHERE user_id = $1
			AND friend_id = $2`
		_, err = tar.link_.Exec(context.Background(), deleteQuery, task.Receiver, task.Sender)
		if err != nil {
			task.State = "database error"
			task.Feedback = err.Error()
		} else {
			task.State = "rejected"
			task.Feedback = "已拒绝好友请求"
			tar.back_put_ <- *task
			task.Feedback = "已拒绝好友请求"
			tar.to_put_ <- *task
			tar.to_put_ = tar.checker_.GetOut(task.Receiver)
		}
	}
}

func (tar *worker) act_sync(task *sqlmap.Task) {

}
