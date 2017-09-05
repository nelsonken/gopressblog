package services

import (
	"fmt"
	"github.com/fpay/gopress"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	// DBServerName is the identity of cache service
	DBServerName = "GormDB"
)

// Dialects interface to getdialects
type Dialects interface {
	GetDialects() string
}

// DBOptions dboptions
type DBOptions struct {
	Name      string `yaml:"name"`
	DBType    string `yaml:"dbtype"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	IP        string `yaml:"ip"`
	Port      string `yaml:"port"`
	Charset   string `yaml:"charset"`
	ParseTime string `yaml:"parsetime"`
}

// GetDialects 获取连接字符串
func (dbc *DBOptions) GetDialects() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=Local",
		dbc.User,
		dbc.Password,
		dbc.IP,
		dbc.Port,
		dbc.Name,
		dbc.Charset,
		dbc.ParseTime)
}

// DBService type
type DBService struct {
	ORM *gorm.DB
}

// NewDBService returns instance of cache service
func NewDBService(dbtype string, dia Dialects) *DBService {
	var err error
	s := new(DBService)
	s.ORM, err = gorm.Open(dbtype, dia.GetDialects())
	if err != nil {
		panic(err)
	}

	return s
}

// ServiceName is used to implements gopress.Service
func (s *DBService) ServiceName() string {
	return DBServerName
}

// RegisterContainer is used to implements gopress.Service
func (s *DBService) RegisterContainer(c *gopress.Container) {
}
