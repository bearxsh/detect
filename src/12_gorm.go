package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type DetectTask struct {
	Id          int
	TaskId      string
	Name        string
	Action      string `gorm:"-"`
	Type        string
	Param       string
	Interval    int
	Status      int
	ReportAddr  string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName 指定表名
func (task *DetectTask) TableName() string {
	return "detect_task"
}

func main() {
	db, err := gorm.Open(sqlite.Open("product.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	/*err = db.AutoMigrate(&Product{})
	if err != nil {
		panic(err)
	}
	// Create
	db.Create(&Product{Code: "D42", Price: 100})
	// Read
	var product Product
	db.First(&product, "code = ?", "D43") // find product with code D42
	fmt.Println(product)*/

	sqlStmt := `CREATE TABLE IF NOT EXISTS "detect_task" (
					"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
					"task_id" TEXT NOT NULL,
					"name" TEXT NOT NULL,
					"type" TEXT NOT NULL,
					"param" TEXT NOT NULL,
					"interval" INTEGER NOT NULL,
					"status" INTEGER NOT NULL,
					"report_addr" TEXT NOT NULL,
					"description" TEXT NOT NULL,
					"created_at" datetime NOT NULL,
					"updated_at" datetime
				);`
	if err = db.Exec(sqlStmt).Error; err != nil {
		panic(err)
	}

	createIndexStmt := `CREATE UNIQUE INDEX IF NOT EXISTS "unique_index_task_id" ON "detect_task" ("task_id" ASC)`
	if err = db.Exec(createIndexStmt).Error; err != nil {
		panic(err)
	}

	task := DetectTask{
		TaskId:      "105",
		Name:        "任务1",
		Type:        "ping",
		Param:       "www.baidu.com",
		Interval:    3,
		Status:      0,
		ReportAddr:  "localhost:8888",
		Description: "我是描述3",
	}
	err = db.Create(&task).Error
	//db.Where("task_id=?", "100").Delete(&DetectTask{})
	//err = db.Model(&DetectTask{}).Where("task_id=?", "103").Select("Name", "Type", "Param", "Interval", "Status", "ReportAddr", "Description").Updates(&task).Error
	/*	var task1 DetectTask
		err = db.Where("task_id=?", "103").First(&task1).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("数据不存在")
			return
		}*/

/*	var tasks []DetectTask
	if err = db.Where("status=?", 1).Find(&tasks).Error; err != nil {
		panic(err)
	} else {
		for _, task := range tasks {
			fmt.Println(task)
		}
	}*/
	db.Model(&DetectTask{}).Where("task_id=?", "105").Update("status", 0)
}
