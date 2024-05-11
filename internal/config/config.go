package config

import (
	"fmt"
	_ "fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type fommuConfig struct {
    SubDomain   string
    Domain      string
    URL         string
    Port        int
    FileHost    string
}

type dbConfig struct {
    DBHost      string
    DBPort      int
    DBName      string
    DBUser      string
    DBPass      string
}

var DB *dbConfig
var Fommu *fommuConfig

func Init() {
    initFommuConfig()
    initDBConfig()
}

func initFommuConfig() {
    viper.AddConfigPath(".")
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    if err := viper.ReadInConfig(); err != nil {
        panic(fmt.Errorf("fatal error config file: %v", err))
    }
    Fommu = &fommuConfig{}
    Fommu.SubDomain = viper.GetString("fommu.subdomain")
    Fommu.Domain  = viper.GetString("fommu.domain")
    Fommu.FileHost = viper.GetString("fommu.filehost")
    url := "https://"
    if Fommu.SubDomain != "" {
        url += Fommu.SubDomain + "."
    }
    url += Fommu.Domain
    Fommu.URL = url

    port, err := strconv.Atoi(os.Getenv("port"))
    if err != nil {
        panic(err.Error())
    }
    Fommu.Port = port
}

func initDBConfig() {
    DB = &dbConfig{}
    
    DB.DBHost = os.Getenv("dbhost")
    if port, err := strconv.Atoi(os.Getenv("dbport")); err != nil {
        DB.DBPort = 5432
    } else {
        DB.DBPort = port
    }
    DB.DBName = os.Getenv("dbname")
    DB.DBUser = os.Getenv("dbuser")
    DB.DBPass = os.Getenv("dbpass")
}
