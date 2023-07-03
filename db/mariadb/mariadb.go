package mariadb

import (
	"github.com/Nerinyan/Nerinyan-APIV2/config"
	"github.com/pterm/pterm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var Mariadb *gorm.DB

func Connect() {

	orm, err := gorm.Open(
		mysql.Open(config.Config.Sql.Url), &gorm.Config{
			AllowGlobalUpdate: true,
			//                                        config.Env.Debug ? debug : error
			Logger:                                   logger.Default.LogMode(logger.Info),
			CreateBatchSize:                          100,
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)

	if Mariadb = orm; err != nil || orm == nil {
		pterm.Fatal.WithShowLineNumber().Println("Gorm Connect Fail", err)
		panic(err)
	}
	var one int
	if err = orm.Raw("SELECT 1").Scan(&one).Error; err != nil || one != 1 {
		pterm.Error.WithShowLineNumber().Println("Gorm Connect Fail", err)
		panic(err)
	}
	db, err := orm.DB()
	if err == nil {
		db.SetMaxIdleConns(3)
		db.SetConnMaxLifetime(time.Second * 30)
		db.SetMaxOpenConns(100)
	} else {
		pterm.Fatal.WithShowLineNumber().Println("Failed to get gorm database", err)
	}
	pterm.Success.Println("GORM Connected")
}
