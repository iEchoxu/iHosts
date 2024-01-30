# What is iHosts
A tool for updating the hosts file

## How To Use ?
- `go mod tidy`
- `go run main.go configs/template/github.json` (Run with root privileges)

## FAQ
- 连接 ipaddress.com 异常：%!s(<nil>)
  - 检查网络,另外可能时触发了 ipaddress.com 的反爬机制,退出 powershell 然后再次执行此程序
- Win11 报病毒错误
  - 因为该程序会修改 hosts 文件，所以会被报病毒错误
  - 解决办法: `设置---隐私和安全性---打开 windows 安全中心---病毒和威胁防护---管理设置---排除项---添加或删除排除项`,选择解压后的 iHosts 文件夹
## 效果图
![效果图](./screenhot/ihosts.gif)