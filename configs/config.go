package configs

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	StartUrls      []string `json:"start_urls"`
	PingLog        bool     `json:"ping_log"`
	GoroutineCount int8     `json:"goroutine_count"`
}

var (
	seeds *Config
	once  sync.Once
)

func Load() *Config {
	once.Do(func() {
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalln("打开配置文件出错", err)
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		if err = json.NewDecoder(file).Decode(&seeds); err != nil {
			log.Fatalln("读取 json 文件出错", err)
		}
	})

	return seeds
}
