package orm

import (
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"github.com/go-xorm/xorm"
)

// 全局默认Orm引擎, GoRoutine安全
// 无需手动关闭, 程序退出时自动关闭
var Engine *xorm.Engine

// Orm引擎初始化
func init() {
	var err error
	Engine, err = xorm.NewEngine("mysql", "root:123456@tcp(10.141.12.65:3306)/test?charset=utf8")
	if err != nil {
		fmt.Println("mysql connection error")
	}
	// 最大打开的连接数, 默认值为0: 不限制
	Engine.SetMaxOpenConns(2000)
	// 闲置的连接数
	Engine.SetMaxIdleConns(100)
}