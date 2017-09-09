package main

import (
	"blog/config"
	"blog/models"
	"blog/services"
	"context"
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
}

func main() {
	//createBlogIndex()
	syncPosts()
}
