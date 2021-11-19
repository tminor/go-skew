package skew

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	allowUnexportedOpts  = cmp.AllowUnexported(value{})
	letters              = []rune("abcdefghijklmnopqrstuvwxyz0123456789/_-;.")
)

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Test_groups(t *testing.T) {
	// len("mississippi") === 2 (mod 3)
	t.Run("Target(\"mississippi\")", func(t *testing.T) {
		x      := "mississippi"
		target := Target(x)

		v0  := &value{Start: 0,  Length: 11, target: &target}
		v1  := &value{Start: 1,  Length: 10, target: &target}
		v2  := &value{Start: 2,  Length: 9,  target: &target}
		v3  := &value{Start: 3,  Length: 8,  target: &target}
		v4  := &value{Start: 4,  Length: 7,  target: &target}
		v5  := &value{Start: 5,  Length: 6,  target: &target}
		v6  := &value{Start: 6,  Length: 5,  target: &target}
		v7  := &value{Start: 7,  Length: 4,  target: &target}
		v8  := &value{Start: 8,  Length: 3,  target: &target}
		v9  := &value{Start: 9,  Length: 2,  target: &target}
		v10 := &value{Start: 10, Length: 1,  target: &target}

		trigrams0 := &value{Start: 1,  Length: 3, target: &target}
		trigrams1 := &value{Start: 4,  Length: 3, target: &target}
		trigrams2 := &value{Start: 7,  Length: 3, target: &target}
		trigrams3 := &value{Start: 10, Length: 3, target: &target}
		trigrams4 := &value{Start: 2,  Length: 3, target: &target}
		trigrams5 := &value{Start: 5,  Length: 3, target: &target}
		trigrams6 := &value{Start: 8,  Length: 3, target: &target}

		want := []SuffixArray{
			{v0, v3, v6, v9},
			{v1, v4, v7, v10},
			{v2, v5, v8},
			{trigrams0, trigrams1, trigrams2, trigrams3, trigrams4, trigrams5, trigrams6},
		}

		s0, s1, s2, s12 := target.groups()

		got := []SuffixArray{s0, s1, s2, s12}

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"mississippi\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	// len("processing") === 1 (mod 3)
	t.Run("Target(\"pascagoula\")", func(t *testing.T) {
		x      := "pascagoula"
		target := Target(x)

		v0  := &value{Start: 0, Length: 10, target: &target}
		v1  := &value{Start: 1, Length: 9,  target: &target}
		v2  := &value{Start: 2, Length: 8,  target: &target}
		v3  := &value{Start: 3, Length: 7,  target: &target}
		v4  := &value{Start: 4, Length: 6,  target: &target}
		v5  := &value{Start: 5, Length: 5,  target: &target}
		v6  := &value{Start: 6, Length: 4,  target: &target}
		v7  := &value{Start: 7, Length: 3,  target: &target}
		v8  := &value{Start: 8, Length: 2,  target: &target}
		v9  := &value{Start: 9, Length: 1,  target: &target}

		trigrams0 := &value{Start: 1, Length: 3, target: &target}
		trigrams1 := &value{Start: 4, Length: 3, target: &target}
		trigrams2 := &value{Start: 7, Length: 3, target: &target}
		trigrams3 := &value{Start: 2, Length: 3, target: &target}
		trigrams4 := &value{Start: 5, Length: 3, target: &target}
		trigrams5 := &value{Start: 8, Length: 3, target: &target}

		want := []SuffixArray{
			{v0, v3, v6, v9},
			{v1, v4, v7},
			{v2, v5, v8},
			{trigrams0, trigrams1, trigrams2, trigrams3, trigrams4, trigrams5},
		}

		s0, s1, s2, s12 := target.groups()

		got := []SuffixArray{s0, s1, s2, s12}

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"pascagoula\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	// len("oxford") === 0 (mod 3)
	t.Run("Target(\"tupelo\")", func(t *testing.T) {
		x      := "tupelo"
		target := Target(x)

		v0  := &value{Start: 0, Length: 6, target: &target}
		v1  := &value{Start: 1, Length: 5, target: &target}
		v2  := &value{Start: 2, Length: 4, target: &target}
		v3  := &value{Start: 3, Length: 3, target: &target}
		v4  := &value{Start: 4, Length: 2, target: &target}
		v5  := &value{Start: 5, Length: 1, target: &target}

		trigrams0 := &value{Start: 1, Length: 3, target: &target}
		trigrams1 := &value{Start: 4, Length: 3, target: &target}
		trigrams2 := &value{Start: 2, Length: 3, target: &target}
		trigrams3 := &value{Start: 5, Length: 3, target: &target}

		want := []SuffixArray{
			{v0, v3},
			{v1, v4},
			{v2, v5},
			{trigrams0, trigrams1, trigrams2, trigrams3},
		}

		s0, s1, s2, s12 := target.groups()

		got := []SuffixArray{s0, s1, s2, s12}

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"tupelo\").groups() mismatch (-want +got):\n%s", diff)
		}
	})


	// D'Lo is a town in Mississippi, and without the apostrophe, len("dlo") == 3
	t.Run("Target(\"dlo\")", func(t *testing.T) {
		x      := "dlo"
		target := Target(x)

		v0 := &value{Start: 0, Length: 3, target: &target}
		v1 := &value{Start: 1, Length: 2, target: &target}
		v2 := &value{Start: 2, Length: 1, target: &target}

		trigrams0 := &value{Start: 1, Length: 3, target: &target}
		trigrams1 := &value{Start: 2, Length: 3, target: &target}

		want := []SuffixArray{
			{v0},
			{v1},
			{v2},
			{trigrams0, trigrams1},
		}

		s0, s1, s2, s12 := target.groups()

		got := []SuffixArray{s0, s1, s2, s12}

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"dlo\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	// There are no cities, towns, counties, or census-designated
	// places in Mississippi whose names are less than three
	// letters long :(
	t.Run("Target(\"ab\")", func(t *testing.T) {
		x      := "ab"
		target := Target(x)

		v0 := &value{Start: 0, Length: 2, target: &target}
		v1 := &value{Start: 1, Length: 1, target: &target}

		trigrams0 := &value{Start: 1, Length: 3, target: &target}

		want := []SuffixArray{
			{v0},
			{v1},
			{},
			{trigrams0},
		}

		s0, s1, s2, s12 := target.groups()

		got := []SuffixArray{s0, s1, s2, s12}

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"ab\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Target(\"a\")", func(t *testing.T) {
		x      := "a"
		target := Target(x)

		v0 := &value{Start: 0, Length: 1, target: &target}

		want := []SuffixArray{
			{v0},
			{},
			{},
			{},
		}

		s0, s1, s2, s12 := target.groups()

		got := []SuffixArray{s0, s1, s2, s12}

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"a\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	// The apostrophe in D'Lo is replaced with \u2018 to test
	// UTF-8 support.
	t.Run("Target(\"d‘lo\")", func(t *testing.T) {
		x      := "d‘lo"
		target := Target(x)

		v0 := &value{Start: 0, Length: 4, target: &target}
		v1 := &value{Start: 1, Length: 3, target: &target}
		v2 := &value{Start: 2, Length: 2, target: &target}
		v3 := &value{Start: 3, Length: 1, target: &target}

		trigrams0 := &value{Start: 1, Length: 3, target: &target}
		trigrams1 := &value{Start: 2, Length: 3, target: &target}

		want := []SuffixArray{
			{v0, v3},
			{v1},
			{v2},
			{trigrams0, trigrams1},
		}

		s0, s1, s2, s12 := target.groups()

		got := []SuffixArray{s0, s1, s2, s12}

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"d‘lo\").groups() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_rank(t *testing.T) {
	t.Run("Target(\"mississippi\")", func(t *testing.T) {
		x      := "mississippi"
		target := Target(x)

		_, _, _, s12 := target.groups()

		s12.rank()

		want := SuffixArray{
			&value{Start: 10, Length: 3, Rank: 0, target: &target},
			&value{Start: 7,  Length: 3, Rank: 1, target: &target},
			&value{Start: 4,  Length: 3, Rank: 2, target: &target},
			&value{Start: 1,  Length: 3, Rank: 3, target: &target},
			&value{Start: 8,  Length: 3, Rank: 4, target: &target},
			&value{Start: 5,  Length: 3, Rank: 5, target: &target},
			&value{Start: 2,  Length: 3, Rank: 6, target: &target},
		}

		got := s12

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("rank() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Target(\"pascagoula\")", func(t *testing.T) {
		x      := "pascagoula"
		target := Target(x)

		want := SuffixArray{
			&value{Start: 4, Length: 3, Rank: 0, target: &target},
			&value{Start: 1, Length: 3, Rank: 1, target: &target},
			&value{Start: 5, Length: 3, Rank: 2, target: &target},
			&value{Start: 8, Length: 3, Rank: 3, target: &target},
			&value{Start: 2, Length: 3, Rank: 4, target: &target},
			&value{Start: 7, Length: 3, Rank: 5, target: &target},
		}

		_, _, _, s12 := target.groups()

		s12.rank()

		got := s12

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"pascagoula\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Target(\"tupelo\")", func(t *testing.T) {
		x      := "tupelo"
		target := Target(x)

		want := SuffixArray{
			&value{Start: 4, Length: 3, Rank: 0, target: &target},
			&value{Start: 5, Length: 3, Rank: 1, target: &target},
			&value{Start: 2, Length: 3, Rank: 2, target: &target},
			&value{Start: 1, Length: 3, Rank: 3, target: &target},
		}

		_, _, _, s12 := target.groups()

		s12.rank()

		got := s12

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"tupelo\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Target(\"dlo\")", func(t *testing.T) {
		x      := "dlo"
		target := Target(x)

		want := SuffixArray{
			&value{Start: 1, Length: 3, Rank: 0, target: &target},
			&value{Start: 2, Length: 3, Rank: 1, target: &target},
		}

		_, _, _, s12 := target.groups()

		s12.rank()

		got := s12

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"dlo\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Target(\"ab\")", func(t *testing.T) {
		x      := "ab"
		target := Target(x)

		want := SuffixArray{
			&value{Start: 1, Length: 3, Rank: 0, target: &target},
		}

		_, _, _, s12 := target.groups()

		s12.rank()

		got := s12

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"ab\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Target(\"a\")", func(t *testing.T) {
		x      := "a"
		target := Target(x)

		want := SuffixArray{}

		_, _, _, s12 := target.groups()

		s12.rank()

		got := s12

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"a\").groups() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("Target(\"d‘lo\")", func(t *testing.T) {
		x      := "d‘lo"
		target := Target(x)

		want := SuffixArray{
			{Start: 2, Length: 3, Rank: 0, target: &target},
			{Start: 1, Length: 3, Rank: 1, target: &target},
		}

		_, _, _, s12 := target.groups()

		s12.rank()

		got := s12

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Target(\"d‘lo\").groups() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_New(t *testing.T) {
	t.Run("Target(\"mississippi\")", func(t *testing.T) {
		sa := New("mississippi")

		want := []string{
			"i",
			"ippi",
			"issippi",
			"ississippi",
			"mississippi",
			"pi",
			"ppi",
			"sippi",
			"sissippi",
			"ssippi",
			"ssissippi",
		}

		got := make([]string, len(sa))

		for i, v := range sa {
			got[i] = string(*v.Get())
		}

		if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
			t.Errorf("Insert() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_Insert(t *testing.T) {
	// t1 := Target("georgia")
	// t2 := Target("mississippi")

	sa := New("mississippi")
	newSa := sa.Insert("georgia")

	// want := &SuffixArray{
	// 	{Start: 7,  Length: 1,  target: &t1}, // a
	// 	{Start: 3,  Length: 6,  target: &t1}, // eorgia
	// 	{Start: 0,  Length: 7,  target: &t1}, // georgia
	// 	{Start: 4,  Length: 3,  target: &t1}, // gia
	// 	{Start: 10, Length: 1,  target: &t2}, // i
	// 	{Start: 5,  Length: 2,  target: &t1}, // ia
	// 	{Start: 7,  Length: 4,  target: &t2}, // ippi
	// 	{Start: 4,  Length: 7,  target: &t2}, // issippi
	// 	{Start: 1,  Length: 10, target: &t2}, // ississippi
	// 	{Start: 0,  Length: 11, target: &t2}, // mississippi
	// 	{Start: 2,  Length: 5,  target: &t1}, // orgia
	// 	{Start: 9,  Length: 2,  target: &t2}, // pi
	// 	{Start: 8,  Length: 3,  target: &t2}, // ppi
	// 	{Start: 3,  Length: 4,  target: &t1}, // rgia
	// 	{Start: 6,  Length: 5,  target: &t2}, // sippi
	// 	{Start: 3,  Length: 8,  target: &t2}, // sissippi
	// 	{Start: 5,  Length: 6,  target: &t2}, // ssippi
	// 	{Start: 2,  Length: 9,  target: &t2}, // ssissippi
	// }

	want := []string{
		"a",
		"eorgia",
		"georgia",
		"gia",
		"i",
		"ia",
		"ippi",
		"issippi",
		"ississippi",
		"mississippi",
		"orgia",
		"pi",
		"ppi",
		"rgia",
		"sippi",
		"sissippi",
		"ssippi",
		"ssissippi",
	}

	got := newSa.ToStrings()

	if diff := cmp.Diff(want, got, allowUnexportedOpts); diff != "" {
		t.Errorf("Insert() mismatch (-want +got):\n%s", diff)
	}
}

func Benchmark_New(b *testing.B) {
	b.Run("1000000", func(b *testing.B) {
		New(randStr(1000000))
	})

	b.Run("2000000", func(b *testing.B) {
		New(randStr(2000000))
	})


	b.Run("4000000", func(b *testing.B) {
		New(randStr(4000000))
	})

	b.Run("8000000", func(b *testing.B) {
		ret := New(randStr(8000000))
		fmt.Printf("%d\n", len(ret))
	})
}
