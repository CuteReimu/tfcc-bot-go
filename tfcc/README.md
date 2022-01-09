# 查询TFCC分数表的功能

想要启用TFCC分数表的功能，请在运行程序之前在`assets/`下增加`score.yaml`。如果没有这个文件，则不启用查询分数表的功能。

文件格式如下：

```yaml
"6":  # 作品
  - work: "6"         # 作品
    rank: Lunatic     # 难度
    route: ""         # 路线（永夜抄6A或者6B）
    character: Reimu  # 机体
    ctype: A          # 子机
    allspell: false   # 是否是全卡
    jf: 10.26         # 分数
  - work: "6"
    rank: Lunatic
    route: ""
    character: Reimu
    ctype: A
    allspell: false
    jf: 10.26
"7":
  - work: "7"
    rank: Lunatic
    route: ""
    character: Reimu
    ctype: A
    allspell: false
    jf: 10.17
```

rank 难度枚举（如果这一作分数表不区分难度则为空）

| 枚举 | 备注 |
| --- | --- |
| Easy | |
| Normal | |
| Hard | |
| Lunatic | |
| Extra | 妖妖梦PH和EX使用同一套分数表，所以都用Extra |

character 机体枚举

| 枚举 | 含义 | 备注 |
| --- | --- | --- |
| Reimu | 灵梦 | |
| Marisa | 魔理沙 | |
| Sakuya | 咲夜 | |
| Sanae | 早苗 | |
| Youmu | 妖梦 | |
| RY | 结界组 | 仅在永夜抄中有 |
| MA | 咏唱组 | 仅在永夜抄中有 |
| SR | 红魔组 | 仅在永夜抄中有 |
| YY | 幽冥组 | 仅在永夜抄中有 |
| Yukari | 紫 | 仅在永夜抄中有 |
| Alice | 爱丽丝 | 仅在永夜抄中有 |
| Remilia | 蕾米莉亚 | 仅在永夜抄中有 |
| Yuyuko | 幽幽子 | 仅在永夜抄中有 |
| Reisen | 灵仙 | 仅在绀珠传中有 |
| Cirno | 琪露诺 | 仅在天空璋中有。大战争为空 |
| Aya | 射命丸文 | 仅在天空璋中有 |

ctype 子机枚举

| 枚举 | 含义 | 作品 |
| --- | --- | --- |
| Spring | 春 | 天空璋 |
| Summer | 夏 | 天空璋 |
| Autumn | 秋 | 天空璋 |
| Winter | 冬 | 天空璋 |
| Wolf | 狼 | 鬼形兽 |
| Otter | 獭 | 鬼形兽 |
| Eagle | 鹰 | 鬼形兽 |

