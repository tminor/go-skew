package skew

import (
	"math"
	"strconv"
)

type Target []rune

type value struct {
	Start, Length, Rank  int
	target               *Target
}

type SuffixArray []*value

type SearchResult struct {
	Index int
	Value string
}

func (s SuffixArray) ToStrings() []string {
	ret := make([]string, len(s))

	for i, v := range s {
		ret[i] = string(*(*v).Get())
	}

	return ret
}

func (s SuffixArray) trigrams(s2 SuffixArray, target *Target) (s12 SuffixArray) {
	sa := SuffixArray(make([]*value, len(s) + len(s2)))
	s12 = sa

	for i , j := 0, 0; i < len(s) + len(s2); i++ {
		if i < len(s) {
			v := &value{Start: s[i].Start, Length: 3}
			v.target = target
			s12[i] = v
		} else {
			v := &value{Start: s2[j].Start, Length: 3}
			v.target = target
			s12[i] = v
			j++
		}
	}

	return s12
}

func (v *value) Get() *[]rune {
	val := make([]rune, v.Length)

	start := v.Start
	end   := start + v.Length
	for i, j := start, 0; i <= end && j < v.Length; i++ {
		if i >= len(*v.target) {
			nul := rune(0)
			val[j] = nul
		} else {
			b := (*v.target)[i]
			val[j] = b
		}
		j++
	}

	return &val
}

// func (s *SuffixArray) Search(pattern string, bounds ...int) []SearchResult {
// 	ret := make([]SearchResult, 0)

// 	max := *(*s)[len(*s) - 1].Get()
// 	min := *(*s)[0].Get()
// 	mid := *(*s)[(len(*s) - 1) / 2].Get()

// 	pat := []rune(pattern)

// 	return ret
// }

func (s SuffixArray) Insert(strs ...string) SuffixArray {
	for _, str := range strs {
		sa := New(str)
		for j, k := 0, 0; j < len(s) && k < len(sa); {
			s1 := string(*s[j].Get())
			s2 := string(*sa[k].Get())

			if s1 < s2 {
				j++
			} else if s2 <= s1 {
				s = append(s[:j+1], s[j:]...)
				s[j] = sa[k]
				j++
				k++
			} else if j == len(s) - 1 {
				s = append(s, sa[k:]...)
			}
		}
	}

	return s
}

func (s SuffixArray) rank() {
	rankMap  := make(map[*[]rune]int)
	valueMap := make(map[*value]*[]rune)
	values   := make([]*[]rune, len(s))
	tmp      := make(SuffixArray, len(s))

	copy(tmp, s)

	for i, v := range s {
		val := v.Get()
		values[i] = val
		valueMap[v] = val
	}

	sorted := radixSort(values)

	testEq := func(a, b []rune) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	for i, v := range sorted {
		if _, ok := rankMap[v]; !ok {
			if i < len(sorted)-1 && testEq(*v, *(sorted[i+1])) {
				rankMap[v] = i + 1
			} else if i > 0 && testEq(*v, *(sorted[i-1])){
				rankMap[v] = i - 1
			} else {
				rankMap[v] = i
			}
		}
	}

	for _, v := range s {
		v.Rank = rankMap[valueMap[v]]
	}

	for _, v := range tmp {
		s[v.Rank] = v
	}
}

func merge(s0, s1, s2, s12 SuffixArray) SuffixArray {
	suffixes := make(SuffixArray, len(*(*(s0)[0]).target))

	startRankMap := make(map[int]int)

	ngramsMap := make(map[string]*value)

	s0RankNgrams := make([]string, 0)
	s1RankNgrams := make([]string, 0)
	s2RankNgrams := make([]string, 0)

	for _, v := range s12 {
		startRankMap[v.Start] = v.Rank
	}

	insertNgram := func(dst []string, val string) []string {
		for min, max := 0, len(dst)-1; true; {
			if max < 0 {
				dst = []string{val}
				break
			} else if min == max || min+1 == max {
				if dst[max] < val {
					dst = append(dst, val)
				} else {
					dst = append(dst[:min+1], dst[min:]...)
					if val > dst[min] {
						dst[min+1] = val
					} else {
						dst[min] = val
					}
				}

				break
			} else {
				mid := int(math.Floor(float64(max)/2.0))
				if val > dst[max] {
					dst = append(dst, val)
					break
				} else if val <= dst[mid] {
					max = mid
				} else if val > dst[mid] {
					min = mid
				}
			}
		}

		return dst
	}

	for _, v := range s0 {
		s := string(*(*v).Get())
		padding := len(s) % 3

		for i := padding; i > 0; i-- {
			s = s + "\x00"
		}

		double := string(s[0]) + strconv.Itoa(startRankMap[v.Start + 1])
		triple := string(s[0:2] + strconv.Itoa(startRankMap[v.Start + 2]))

		if double <= triple {
			s0RankNgrams = insertNgram(s0RankNgrams, double)
			ngramsMap[double] = v
		} else {
			s0RankNgrams = insertNgram(s0RankNgrams, triple)
			ngramsMap[triple] = v
		}
	}

	for _, v := range s1 {
		s := string(*(*v).Get())
		padding := len(s) % 3

		for i := padding; i > 0; i-- {
			s = s + "\x00"
		}

		double := string(s[0]) + strconv.Itoa(startRankMap[v.Start + 1])
		s1RankNgrams = insertNgram(s1RankNgrams, double)
		ngramsMap[double] = v
	}

	for _, v := range s2 {
		s := string(*(*v).Get())
		padding := len(s) % 3

		for i := padding; i > 0; i-- {
			s = s + "\x00"
		}

		triple := string(s[0:2] + strconv.Itoa(startRankMap[v.Start + 2]))
		s2RankNgrams = insertNgram(s2RankNgrams, triple)
		ngramsMap[triple] = v
	}

	s12RankNgrams := append(s1RankNgrams, s2RankNgrams...)

	for i, j, k := 0, 0, 0; i < len(*(*(s0)[0]).target); i++ {
		if j < len(s0RankNgrams) {
			suf1 := s0RankNgrams[j]

			suf2 := ""
			if k >= len(s12RankNgrams) {
				suf2 = "zzzzzzzzzzzzzz"
			} else {
				suf2 = s12RankNgrams[k]
			}

			if suf1 <= suf2 {
				suffixes[i] = ngramsMap[s0RankNgrams[j]]
				j++
			} else {
				suffixes[i] = ngramsMap[s12RankNgrams[k]]
				k++
			}
		} else {
			suffixes[i] = ngramsMap[s12RankNgrams[k]]
			k++
		}
	}

	return suffixes
}

func (t *Target) groups() (s0, s1, s2, s12 SuffixArray) {
	max := len(*t)

	largeLen       := int(math.Ceil(float64(len(*t)) / 3.0))
	smallLen       := largeLen - 1
	numLargeGroups := len(*t) % 3

	sa0 := SuffixArray(make([]*value, 0))
	sa1 := SuffixArray(make([]*value, 0))
	sa2 := SuffixArray(make([]*value, 0))
	switch numLargeGroups {
	case 0: // Equivalent to 3 evenly sized groups
		sa0 = SuffixArray(make([]*value, largeLen))
		sa1 = SuffixArray(make([]*value, largeLen))
		sa2 = SuffixArray(make([]*value, largeLen))
	case 2:
		sa0 = SuffixArray(make([]*value, largeLen))
		sa1 = SuffixArray(make([]*value, largeLen))
		sa2 = SuffixArray(make([]*value, smallLen))
	case 1:
		sa0 = SuffixArray(make([]*value, largeLen))
		sa1 = SuffixArray(make([]*value, smallLen))
		sa2 = SuffixArray(make([]*value, smallLen))
	}

	s0 = sa0
	s1 = sa1
	s2 = sa2

	i := 0
	for j := 0; i <= len(s0) - 1; j++ {
		(s0)[i] = &value{Length: max - j, Start: j, target: t}

		j++
		if i < len(s1) {
			s1[i] = &value{Length: max - j, Start: j, target: t}
		}

		j++
		if i < len(s2) {
			s2[i] = &value{Length: max - j, Start: j, target: t}
		}

		i++
	}

	s12 = s1.trigrams(s2, t)

	return s0, s1, s2, s12
}

func New(s string) SuffixArray {
	t := Target(s)

	s0, s1, s2, s12 := t.groups()

	s12.rank()

	return merge(s0, s1, s2, s12)
}
