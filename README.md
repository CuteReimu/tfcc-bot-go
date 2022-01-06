# 东方Project沙包聚集地机器人

这是东方Project沙包聚集地（以下简称“红群”）的机器人，基于[MiraiGo-Template](https://github.com/Logiase/MiraiGo-Template)编写

## 配置文件

第一次运行会自动生成配置文件`application.yaml`，如下：

```yaml
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