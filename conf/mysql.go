package conf

type DbConfig struct {
	DriverName   string
	Dsn          string
	ShowSql      bool
	ShowExecTime bool
	MaxIdle      int
	MaxOpen      int
}

var Db = map[string]DbConfig{
	"db1": {
		DriverName:   "mysql",
		Dsn:          "root:root@tcp(118.25.27.160:33306)/ssodb?charset=utf8mb4&parseTime=true&loc=Local",
		ShowSql:      true,
		ShowExecTime: false,
		MaxIdle:      10,
		MaxOpen:      200,
	},
}
