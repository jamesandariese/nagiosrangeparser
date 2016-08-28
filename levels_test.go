package naglevelparse

import (
	"math"
	"testing"
)

type NagLevelTest struct {
	pattern      string
	neginf       bool
	negseven     bool
	negone       bool
	zero         bool
	pointfive    bool
	one          bool
	onepointfive bool
	two          bool
	seventyseven bool
	posinf       bool
}

func TestNormalRanges(t *testing.T) {
	tests := []NagLevelTest{
		NagLevelTest{
			pattern:      "1:2",
			neginf:       true,
			negseven:     true,
			negone:       true,
			zero:         true,
			pointfive:    true,
			one:          false,
			onepointfive: false,
			two:          false,
			seventyseven: true,
			posinf:       true,
		},
		NagLevelTest{
			pattern:      "-1:1",
			neginf:       true,
			negseven:     true,
			negone:       false,
			zero:         false,
			pointfive:    false,
			one:          false,
			onepointfive: true,
			two:          true,
			seventyseven: true,
			posinf:       true,
		},
		NagLevelTest{
			pattern:      ":77",
			neginf:       true,
			negseven:     true,
			negone:       true,
			zero:         false,
			pointfive:    false,
			one:          false,
			onepointfive: false,
			two:          false,
			seventyseven: false,
			posinf:       true,
		},
		NagLevelTest{
			pattern:      "0:77",
			neginf:       true,
			negseven:     true,
			negone:       true,
			zero:         false,
			pointfive:    false,
			one:          false,
			onepointfive: false,
			two:          false,
			seventyseven: false,
			posinf:       true,
		},
		NagLevelTest{
			pattern:      "~:",
			neginf:       false,
			negseven:     false,
			negone:       false,
			zero:         false,
			pointfive:    false,
			one:          false,
			onepointfive: false,
			two:          false,
			seventyseven: false,
			posinf:       false,
		},
		NagLevelTest{
			pattern:      "~:-7",
			neginf:       false,
			negseven:     false,
			negone:       true,
			zero:         true,
			pointfive:    true,
			one:          true,
			onepointfive: true,
			two:          true,
			seventyseven: true,
			posinf:       true,
		},
		NagLevelTest{
			pattern:      "-7:77.5",
			neginf:       true,
			negseven:     false,
			negone:       false,
			zero:         false,
			pointfive:    false,
			one:          false,
			onepointfive: false,
			two:          false,
			seventyseven: false,
			posinf:       true,
		},
		NagLevelTest{
			pattern:      "-7:76.5",
			neginf:       true,
			negseven:     false,
			negone:       false,
			zero:         false,
			pointfive:    false,
			one:          false,
			onepointfive: false,
			two:          false,
			seventyseven: true,
			posinf:       true,
		},
	}

	for _, test := range tests {
		if pp, err := Compile(test.pattern); err != nil {
			t.Errorf("Couldn't compile %v: %#v", test.pattern, err)
		} else {
			runOne := func(cv float64, expected bool) {
				if r := pp.Compare(cv); r != expected {
					t.Errorf("'%s' Compare(%v) should return %v but returned %v instead", test.pattern, cv, expected, r)
				} else {
					t.Logf("'%s' Compare(%v) returned %v as expected", test.pattern, cv, expected)
				}
			}
			runOne(math.Inf(-1), test.neginf)
			runOne(-7.0, test.negseven)
			runOne(-1.0, test.negone)
			runOne(0.0, test.zero)
			runOne(0.5, test.pointfive)
			runOne(1.0, test.one)
			runOne(1.5, test.onepointfive)
			runOne(2.0, test.two)
			runOne(77.0, test.seventyseven)
			runOne(math.Inf(1), test.posinf)
		}

		if pp, err := Compile("@" + test.pattern); err != nil {
			t.Errorf("Couldn't compile %v: %#v", test.pattern, err)
		} else {
			runOne := func(cv float64, expected bool) {
				if r := pp.Compare(cv); r != expected {
					t.Errorf("'%s' Compare(%v) should return %v but returned %v instead", test.pattern, cv, expected, r)
				} else {
					t.Logf("'%s' Compare(%v) returned %v as expected", "@" + test.pattern, cv, expected)
				}
			}
			runOne(math.Inf(-1), !test.neginf)
			runOne(-7.0, !test.negseven)
			runOne(-1.0, !test.negone)
			runOne(0.0, !test.zero)
			runOne(0.5, !test.pointfive)
			runOne(1.0, !test.one)
			runOne(1.5, !test.onepointfive)
			runOne(2.0, !test.two)
			runOne(77.0, !test.seventyseven)
			runOne(math.Inf(1), !test.posinf)
		}
	}
}

func TestCompileErrors(t *testing.T) {
	runOne := func(pattern string) {
		t.Logf("Testing if %v is a syntax error", pattern)
		if _, err := Compile(pattern); err == nil {
			t.Errorf("%v is a parse error but instead it compiled cleanly", pattern)
		}
	}
	runOne("xyz")
	runOne("@@@")
	runOne("~~")
}
