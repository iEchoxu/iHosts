package global

import "ihosts/configs"

var (
	EnvConfig                      = configs.Load()
	NumberOfUrlsAfterDeDuplication int
)

type Hosts struct {
	Domain string
	IP     string
}
