package ans

import (
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/yanyiwu/gojieba"
)

var PartOfSpeechNeeded map[string]bool = map[string]bool {
	"a": true,
	"ad": true,
	"ag": true,
	"an": true,
	"e": true,
	"g": true,
	"i": true,
	"l": true,
	"n": true,
	"ng": true,
	"nr": true,
	"nrfg": true,
	"nrt": true,
	"ns": true,
	"nt": true,
	"nz": true,
	"o": true,
	"t": true,
	"tg": true,
	"vg": true,
	"vi": true,
	"vn": true,
	"vq": true,
}

type Ans struct {
	TmpDir string
	Jieba *gojieba.Jieba
}

func NewAns() (*Ans) {
	new_ans := new(Ans)
	new_ans.Jieba = gojieba.NewJieba()
	return new_ans
}

func (ans *Ans) FilterTags(tags []string) ([]string) {
	ret := make([]string, 0, 10)
	for _, tag := range tags {
		wordAndPs := strings.Split(tag, "/")
		if len(wordAndPs) != 2 {
			continue
		}
		if !ans.IsPartOfSpeechNeeded(wordAndPs[1]) {
			continue
		}
		ret = append(ret, wordAndPs[0])
	}
	return ret
}

func (ans *Ans) IsPartOfSpeechNeeded(word string) (bool) {
	_, ok := PartOfSpeechNeeded[word]
	return ok
}

func (ans *Ans) CalcDailyWordsTrend(wordCounts map[string]int) ([]string) {
	pl := rankByWordCount(wordCounts)
	wordsToday := make([]string, 0)
	i := 0
	for _, p := range pl {
		if i >= 10 {
			break
		}
		if utf8.RuneCountInString(p.Key) <= 1 {
			continue
		}
		wordsToday = append(wordsToday, p.Key)
		i++
	}
	return wordsToday
}

func rankByWordCount(wordFrequencies map[string]int) PairList{
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}
type Pair struct {
	Key string
	Value int
}
type PairList []Pair
func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

func (ans *Ans) CalcWordCounts(words []string) (map[string]int) {
	ret := make(map[string]int)
	for _, word := range words {
		x := ans.Jieba.Tag(word)
		x = ans.FilterTags(x)
		for _, i := range x {
			count, ok := ret[i]
			if ok {
				ret[i] = count + 1
			} else {
				ret[i] = 1
			}
		}
	}
	return ret
}
