package ans

import (
	"testing"
)

func TestAns(t *testing.T) {
	ans := NewAns()
	wordCounts := ans.CalcWordCounts([]string{"我golang用的好好的呀？今天什么事情也没有发生发生发生"})
	t.Log(wordCounts)
	t.Log(ans.CalcDailyWordsTrend(wordCounts))
}

func TestJieba(t *testing.T) {
	ans := NewAns()
	wordList := ans.Jieba.Tag("貨拉拉拉不拉拉布拉多")
	t.Log(wordList)
	t.Log(ans.FilterTags(wordList))
}
