package translate

import "testing"

func TestWR(t *testing.T) {
	s := `AAAA beat the WR in Hollow Knight - Pantheon of Hallownest: Level - All Bindings (1.4.3.2+). The new WR is 35m 21s.`
	s2 := Translate(s)
	if s2 != `AAAA 打破了世界纪录：圣巢万神殿-四锁 (1.4.3.2+).新的世界纪录是35m 21s.` {
		t.Error(s2)
	}
}

func TestTop3(t *testing.T) {
	s := `BBBB got a new top 3 PB in Hollow Knight Category Extensions - 0 Geo - All Glitches. Their time is 15m 31s.`
	s2 := Translate(s)
	if s2 != `BBBB 获得了前三：0吉欧-允许所有邪道.时间是15m 31s.` {
		t.Error(s2)
	}
}
