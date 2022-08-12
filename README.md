<div align="center">

# 东方Project沙包聚集地机器人

![](https://img.shields.io/github/languages/top/Touhou-Freshman-Camp/tfcc-bot-go "语言")
[![](https://img.shields.io/github/workflow/status/Touhou-Freshman-Camp/tfcc-bot-go/Go)](https://github.com/Touhou-Freshman-Camp/tfcc-bot-go/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/Touhou-Freshman-Camp/tfcc-bot-go)](https://github.com/Touhou-Freshman-Camp/tfcc-bot-go/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/Touhou-Freshman-Camp/tfcc-bot-go)](https://github.com/Touhou-Freshman-Camp/tfcc-bot-go/blob/master/LICENSE "许可协议")
</div>

这是东方Project沙包聚集地（以下简称“红群”）的机器人，基于[MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)编写。

## 声明

* 本项目采用`AGPLv3`协议开源。同时**强烈建议**各位开发者遵循以下原则：
  * **任何间接接触本项目的软件也要求使用`AGPLv3`协议开源**
  * **不鼓励，不支持一切商业用途**
* **由于使用本项目提供的接口、文档等造成的不良影响和后果与本人和红群无关**
* 由于本项目的特殊性，可能随时停止开发或删档
* 本项目为开源项目，不接受任何的催单和索取行为

## 使用方法

建议在编译前更新一下依赖：

```bash
go get -u
```

编译：

```bash
go build -o tfcc-bot.exe
```

然后双击运行生成出来的`tfcc-bot.exe`即可。比较建议在cmd窗口中输入`tfcc-bot.exe`运行，以防panic后报错信息无法看到。

在功能完善后，会将编译好的包放在Release中供大家下载。

关闭程序时，请使用ctrl+C关闭，以确保db、log、bot等模块正常关闭。如果强制退出，可能会导致部分数据未写入硬盘，下次启动时丢失数据。

## 配置文件

第一次运行会自动生成配置文件`application.yaml`，如下：

```yaml
bilibili:
  area_v2: "236"           # 直播分区，236-主机游戏
  mid: "12345678"          # B站ID
  password: "12345678"     # 密码
  room_id: "12345678"      # B站直播间房间号
  username: "13888888888"  # B站用户名
bot:
  loginmethod: qrcode  # 登录方式
  account: 0           # 机器人QQ号
  password: ""         # 机器人密码
qq:
  super_admin_qq: 12345678  # 主管理员QQ号
  qq_group: # 主要功能的QQ群
    - 12345678
```

修改配置文件后重新启动即可。

## 模块

- chatPipeline 非命令式的QQ聊天消息处理。想要新增，实现`pipelineHandler`接口并在`init()`中调用`register()`即可
- commandHandler QQ聊天中输入的命令。想要新增命令，实现`cmdHandler`接口并在`init()`中调用`register()`即可
- db 一个嵌入式Key-Value型数据库，使用这个模块存储的数据会被存在硬盘里，下次重启后仍然保留
- perm 权限管理，管理员和白名单
- main.go 程序入口

## 运行时生成的文件

以下文件会在运行时自动生成

- assets/database/ 是db模块的数据文件
- log/ 日志文件
- application.yaml 是配置文件
- device.json 是设备信息文件。不要删除，否则会被QQ认为你换了一台设备登录
- session.token 是会话信息文件，用于重启时快速恢复登录状态。想要更换QQ号登录请删除这个文件

## 功能一览

**简化版只启用部分功能**

- [x] 管理员、白名单
- [x] B站开播、修改直播标题、查询直播状态
- [ ] ~~随作品、随机体~~
- [ ] ~~B站视频解析~~
- [ ] ~~B站视频推送~~
- [ ] ~~投票~~
- [ ] ~~查新闻~~
- [ ] ~~增加预约功能~~
- [ ] ~~查询分数表~~
- [x] 打断复读
- [ ] ~~随符卡~~
- [ ] ~~rep解析~~

## 第三方库的使用

- github.com/Mrs4s/MiraiGo 一个移植于mirai的golang实现的库
- github.com/Logiase/MiraiGo-Template 基于MiraiGo的多模块设计组合
- github.com/dgraph-io/badger 一个强大的内嵌的数据库系统
- github.com/sirupsen/logrus 一个强大的日志库
- github.com/go-resty/resty 强大的Http Client库
- github.com/tidwall/gjson 易用的json解析库
- github.com/spf13/viper 强大的config解析库
- github.com/ozgio/strutil 一个字符串处理库
- github.com/araddon/dateparse 一个可以识别任意格式日期的库
- github.com/dlclark/regexp2 完整版的正则表达式库
- github.com/pkg/errors 可以将golang本身的error包装的库
- github.com/CuteReimu/threp 一个东方replay文件的解析库
- github.com/CuteReimu/bilibili 一个B站API的golang版sdk
