package main

import (
	"blog/config"
	"blog/models"
	"blog/services"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	index = "blog"
)

func getOrm() *gorm.DB {
	opts := &config.Options{}
	opts.Database = &services.DBOptions{}
	opts.Elastic = &services.ElasticOption{}
	config.GetConfig("../config/config.yaml", opts)
	dbService := services.NewDBService(opts.Database.DBType, opts.Database)

	return dbService.ORM
}

func createBlogIndex() {
	// Create a client
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	// Create an index
	_, err = client.CreateIndex(index).Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println("创建索引成功")
}

func syncPosts() {
	orm := getOrm()
	posts := []*models.Post{}
	err := orm.Select("id, title, created_at, author_id").Find(&posts).Error
	if err != nil {
		panic(err)
	}

	// Create a client
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
	defer cancel()

	for _, post := range posts {
		_, err := client.Index().
			Index(index).
			Type("posts").
			Id(strconv.Itoa(int(post.ID))).
			BodyJson(post).
			Refresh("true").
			Do(ctx)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("同步文章到elastic成功")
}

func clearTodayIncome() {
	db := getOrm()
	err := db.Table("accounts").Where("today_income > ?", 0).Updates(map[string]interface{}{
		"today_income": 0,
	}).Error

	if err != nil {
		panic(err)
	}
	fmt.Println("清除今日活跃点数成功")
}

func awardActiveUser() {
	db := getOrm()
	var incomes []float64
	db.Model(&models.User{}).Order("today_income desc").Limit(10).Pluck("today_income", &incomes)
	lastUserIncome := incomes[len(incomes)-1]
	if lastUserIncome == 0 {
		return
	}
	err := db.Table("accounts").Where("today_income >= ?", lastUserIncome).Updates(map[string]interface{}{"total": gorm.Expr("total + ?", 5)}).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("奖励前十名用户成功")
}

func main() {
	//createBlogIndex()
	if len(os.Args) < 2 {
		panic("usage:\n\t elastic act \n act: syncpost-sync posts to elastic; index: index blog index; scoreclear: clear today's active scores")
	}

	switch os.Args[1] {
	case "syncpost":
		syncPosts()
	case "index":
		createBlogIndex()
	case "scoreclear":
		awardActiveUser()
		clearTodayIncome()
	default:
		panic("act invalid, just support syncposts and index act")
	}
}
