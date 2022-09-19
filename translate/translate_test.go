package translate

import "testing"

func TestWR(t *testing.T) {
	s := `AAAA beat the WR in Hollow Knight - Pantheon of Hallownest: Level - All Bindings (1.4.3.2+). The new WR is 35m 21s.`
	s2 := Translate(s)
	if s2 != `AAAA 打破了世界纪录：圣巢万神殿 - 四锁 (1.4.3.2+).新的世界纪录是35m 21s.` {
		t.Error(s2)
	}
}

func TestWR2(t *testing.T) {
	s := `CCCC beat the WR in Hollow Knight Category Extensions - King's Pass: Level - Slower. The new WR is 0m 55s 900ms.`
	s2 := Translate(s)
	if s2 != `CCCC 打破了世界纪录：国王山道 - Slower.新的世界纪录是0m 55s 900ms.` {
		t.Error(s2)
	}
}

func TestWR3(t *testing.T) {
	s := `AAAA beat the WR in Hollow Knight Category Extensions - Save Myla - 1.4.3.2+ NMG. The new WR is 34m 57s.`
	s2 := Translate(s)
	if s2 != `AAAA 打破了世界纪录：拯救米拉 - 1.4.3.2+无主要邪道.新的世界纪录是34m 57s.` {
		t.Error(s2)
	}
}

func TestWR4(t *testing.T) {
	s := `VVVV beat the WR in Hollow Knight - 112% APB - No Major Glitches. The new WR is 3h 09m 43s.`
	s2 := Translate(s)
	if s2 != `VVVV 打破了世界纪录：112% 全万神殿BOSS - 无主要邪道.新的世界纪录是3h 09m 43s.` {
		t.Error(s2)
	}
}

func TestWR5(t *testing.T) {
	s := `Gusten13 beat the WR in Hollow Knight - All Achievements - No Major Glitches. The new WR is 6h 38m 52s.`
	s2 := Translate(s)
	if s2 != `Gusten13打破了世界纪录：全成就 - 无主要邪道.新的世界纪录是6h 38m 52s.` {
		t.Error(s2)
	}
}

func TestWR6(t *testing.T) {
	s := `Deites beat the WR in Hollow Knight - Console Runs - Any%, PS5/XSeries. The new WR is 59m 45s.`
	s2 := Translate(s)
	if s2 != `Deites 打破了世界纪录：主机速通 - Any%, PS5/XSeries.新的世界纪录是59m 45s.` {
		t.Error(s2)
	}
}

func TestWR7(t *testing.T) {
	s := `YouYu beat the WR in Hollow Knight - Pantheon of the Artist: Level - Any Bindings. The new WR is 3m 04s 360ms.`
	s2 := Translate(s)
	if s2 != `YouYu 打破了世界纪录：艺术家万神殿 - 任意锁.新的世界纪录是3m 04s 360ms.` {
		t.Error(s2)
	}
}

func TestWR8(t *testing.T) {
	s := `Skitter_ beat the WR in Hollow Knight - Console Runs - 112% APB, Switch. The new WR is 4h 23m 21s.`
	s2 := Translate(s)
	if s2 != `Skitter_打破了世界纪录：主机速通 - 112% 全万神殿BOSS, Switch.新的世界纪录是4h 23m 21s.` {
		t.Error(s2)
	}
}

func TestWR9(t *testing.T) {
	s := `fatmoose beat the WR in Hollow Knight - Trial of the Warrior: Level. The new WR is 1m 55s 720ms.`
	s2 := Translate(s)
	if s2 != `fatmoose 打破了世界纪录：勇士的试炼.新的世界纪录是1m 55s 720ms.` {
		t.Error(s2)
	}
}

func TestTop3(t *testing.T) {
	s := `BBBB got a new top 3 PB in Hollow Knight Category Extensions - 0 Geo - All Glitches. Their time is 15m 31s.`
	s2 := Translate(s)
	if s2 != `BBBB 获得了前三：0吉欧 - 允许所有邪道.时间是15m 31s.` {
		t.Error(s2)
	}
}
