package global

import "ihosts/configs"

const IPAddressDomain = "https://www.ipaddress.com/website/"

var (
	EnvConfig                      = configs.Load()
	NumberOfUrlsAfterDeDuplication int
)

type Hosts struct {
	Domain string
	IP     string
}
