package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	DB_TYPE           = "mysql"
	DB_SOURCE         = "root:Zsba@mysql2018*@tcp(127.0.0.1:3306)/switch?charset=utf8&parseTime=true"
	DB_MAXIDLE        = 30
	DB_MAXOPEN        = 30
	DB_ISLOGER        = true
	DB_TLS_ENABLED    = false
	DB_TLS_CLIENTCERT = "db_client.pem"
	DB_TLS_CLIENTKEY  = "db_client.key"
	DB_TLS_CACERT     = "db_ca.pem"
)

type ServerTLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
	CaCert   string
}

type DataBaseConfig struct {
	Type       string
	DataSource string
	MaxIdle    int
	MaxOpen    int
	IsLoger    bool
	TLS        ServerTLSConfig
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func LoadDBConfig(dir string) (*DataBaseConfig, error) {
	path := filepath.Join(dir, "./configs/dao.toml")
	filePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	config := new(DataBaseConfig)
	if CheckFileIsExist(filePath) { //文件存在
		if _, err := toml.DecodeFile(filePath, config); err != nil {
			return nil, err
		}
	} else {
		config.Type = DB_TYPE
		config.DataSource = DB_SOURCE
		config.MaxIdle = DB_MAXIDLE
		config.MaxOpen = DB_MAXOPEN
		config.IsLoger = DB_ISLOGER
		config.TLS.Enabled = DB_TLS_ENABLED
		config.TLS.CertFile = DB_TLS_CLIENTCERT
		config.TLS.KeyFile = DB_TLS_CLIENTKEY
		config.TLS.CaCert = DB_TLS_CACERT

		configBuf := new(bytes.Buffer)
		if err := toml.NewEncoder(configBuf).Encode(config); err != nil {
			return nil, err
		}
		err := ioutil.WriteFile(filePath, configBuf.Bytes(), 0666)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
