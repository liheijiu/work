package app_mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"work/base/app_conf"
)

var (
	db *sql.DB //数据库连接池对象
)

type app_user struct {
	id   int
	name string
	age  int
}

//mysql  init
func InitDb(cfg *app_conf.AppConf) (err error) {
	address := cfg.AppMysql.Address
	database := cfg.AppMysql.Database
	userName := cfg.AppMysql.Username
	password := cfg.AppMysql.Password
	dsn := fmt.Sprintf(userName + ":" + password + address + "/" + database)
	db, err = sql.Open("mysql", dsn) //校验格式是否正确
	if err != nil {
		return
	}
	err = db.Ping() //尝试连接数据库，校验用户名与密码是否正确
	if err != nil {
		return
	}
	//设置数据库连接池的数额
	db.SetMaxOpenConns(10)
	//设置数据库连接池的空闲数额
	db.SetMaxIdleConns(5)

	return
}

//单条查询
func MysqlQuire(id int) {
	//声明一个对象
	var u1 app_user

	//1.查询单条记录的sql语句
	sqlStr := `SELECT id,name,age from user  where id=?`

	//2.执行
	err := db.QueryRow(sqlStr, id).Scan(&u1.id, &u1.name, &u1.age)
	if err != nil {
		return
	}

	//3.拿到结果

	//4.打印
	fmt.Printf("%#v\n", u1)
}

//多条查询
func MysqlQuires(n int) {
	//声明一个对象
	var u1 app_user
	//sql语句
	sqlStr := `SELECT id,name,age from user  where id > ?;`
	//执行
	rows, err := db.Query(sqlStr, n)
	if err != nil {
		fmt.Printf("exec %s query failed, err%v\n", sqlStr, err)
		return
	}
	//关闭rows
	defer rows.Close()
	//循环取值
	for rows.Next() {
		err := rows.Scan(&u1.id, &u1.name, &u1.age)
		if err != nil {
			return
		}
		fmt.Printf("%#v\n", u1)
	}
}

//新增、插入
func MysqInsert() {
	//写sql语句
	sqlStr := `insert into user (name,age) VALUES("云朵", 18)`

	//Exec插入数据
	ret, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	//如果是插入数据，将会得到id
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get id failed, err:%v\n", err)
		return
	}
	fmt.Println("id:", id)
}

//更新及修改
func MysqUpdate(newAge, id int) {
	//写sql语句
	sqlStr := `update  user set age=? where id=?`
	ret, err := db.Exec(sqlStr, newAge, id)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	affected, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get id failed, err:%v\n", err)
		return
	}
	fmt.Printf("update line %d\n", affected)
}

//删除
func MysqDelete(id int) {
	//写sql语句
	sqlStr := `delete  from user where id=?`
	ret, err := db.Exec(sqlStr, id)
	if err != nil {
		return
	}
	affected, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get id failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete %d line\n", affected)
}

//预处理-插入
func PrepareInsertDemo() {
	sqlStr := `insert into user(name,age) values(?,?)`
	prepare, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("Prepare failed, err:%v\n", err)
		return
	}
	defer prepare.Close()

	//后续只需要拿到prepare去执行操作
	var m = map[string]int{
		"云朵": 18,
		"雨":  3,
		"雷电": 100,
	}
	for k, v := range m {
		ret, err := prepare.Exec(k, v)
		if err != nil {
			fmt.Printf("insert failed, err:%v\n", err)
			return
		}
		id, err := ret.LastInsertId()
		if err != nil {
			fmt.Printf("get  id failed, err:%v\n", err)
			return
		}
		fmt.Printf("id: %d\n", id)
	}
}

//事务操作
func TransactionDemo() {

	//开启事务
	bg, err := db.Begin()
	if err != nil {
		fmt.Printf("Begin failed, err:%v\n", err)
		return
	}

	//执行多个sql操作
	sqlStr1 := `update  user set age=age-2  where id=1`
	sqlStr2 := `update  users set age=age-2  where id=5`

	_, err = bg.Exec(sqlStr1)
	if err != nil {
		//需要回归
		bg.Rollback()
		fmt.Println("执行事务1出错，需要回滚")
		return
	}

	_, err = bg.Exec(sqlStr2)
	if err != nil {
		//需要回归
		bg.Rollback()
		fmt.Println("执行事务2出错，需要回滚")
		return
	}

	//上面执行如果成功就提交事务
	err = bg.Commit()
	if err != nil {
		fmt.Println("提交出错，需要回滚")
		return
	}
	fmt.Println("事务执行成功")
}
