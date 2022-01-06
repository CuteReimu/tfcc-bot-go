<div align="center">

# 东方Project沙包聚集地机器人

![](https://img.shields.io/github/languages/top/Touhou-Freshman-Camp/tfcc-bot-go "语言")
[![](https://img.shields.io/github/workflow/status/Touhou-Freshman-Camp/tfcc-bot-go/Go)](https://github.com/Touhou-Freshman-Camp/tfcc-bot-go/actions/workflows/golangci-lint.yml "代码分析")
[![](https://img.shields.io/github/contributors/Touhou-Freshman-Camp/tfcc-bot-go)](https://github.com/Touhou-Freshman-Camp/tfcc-bot-go/graphs/contributors "贡献者")
[![](https://img.shields.io/github/license/Touhou-Freshman-Camp/tfcc-bot-go)](https://github.com/Touhou-Freshman-Camp/tfcc-bot-go/blob/master/LICENSE "许可协议")
</div>

这是东方Project沙包聚集地（以下简称“红群”）的机器人，基于[MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)编写

## 配置文件

第一次运行会自动生成配置文件`application.yaml`，如下：

```yaml
bilibili:
  area_v2: "236"  # 直播分区，236-主机游戏
  mid: "12345678"  # B站ID
  password: "12345678"  # 密码
  room_id: "12345678"  # B站直播间房间号
  username: "13888888888"  # B站用户名
bot:
  account: 0  # 机器人QQ号
  password: ""  # 机器人密码（不填就是扫码登录）
qq:
  super_admin_qq: 12345678  # 主管理员QQ号
repeater_interruption:
  allowance: 5  # 打断复读功能限制的复读次数
  cool_down: 3  # 打断复读冷却时间（秒）
  qq_group:  # 打断复读的Q群
  - 12345678
schedule:
  before:  # 预约功能提前提醒时间
  - 10800
  - 21600
  qq_group:  # 预约功能提前QQ群
  - 12345678
  video_push_delay: 600  # 视频推送间隔
```

修改配置文件后重新启动即可。

## 模块

- bilibili 和B站相关的代码
- commandHandler QQ聊天中输入的命令。想要新增命令，继承`cmdHandler`接口并在`init()`中调用`register()`即可
- db 一个嵌入式Key-Value型数据库，使用这个模块存储的数据会被存在硬盘里，下次重启后仍然保留。
- perm 权限管理，管理员和白名单
- main.go 程序入口

## 运行时生成的文件

以下文件会在运行时自动生成

- assets/database/ 是db模块的缓存文件
- log/ 日志文件
- application.yaml 是配置文件
- device.json 是设备信息文件。不要删除，否则会被QQ认为你换了一台设备登录
- session.token 是会话信息文件，用于重启时快速恢复登录状态。想要更换QQ号登录请删除这个文件。
