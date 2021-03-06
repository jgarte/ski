package ski

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	iotaK = Apply(Apply(Apply(Apply(newNode(I), newNode(S)), newNode(K)), newNode(S)), newNode(K))
	iotaS = Apply(Apply(iotaK, newNode(S)), newNode(K))
	jotK  = Apply(Apply(Apply(Apply(Apply(newNode(S), Apply(newNode(K), Apply(newNode(S), Apply(newNode(K), Apply(newNode(S), Apply(newNode(K), newNode(I))))))), newNode(S)), newNode(K)), newNode(S)), newNode(K))
	jotS  = Apply(Apply(Apply(Apply(Apply(Apply(Apply(newNode(S), Apply(newNode(K), Apply(newNode(S), Apply(newNode(K), Apply(newNode(S), Apply(newNode(K), Apply(newNode(S), Apply(newNode(K), Apply(newNode(S), Apply(newNode(K), newNode(I))))))))))), newNode(S)), newNode(K)), newNode(S)), newNode(K)), newNode(S)), newNode(K))
)

type skiTest struct {
	fs       string
	s        string
	n        *Node
	simplify string
	reduce   string
	reduceN  int
}

var valid = []skiTest{
	{"I", "I", NewNode(I), "I", "a", 1},
	{"K", "K", NewNode(K), "K", "a", 2},
	{"S", "S", NewNode(S), "S", "ac(bc)", 3},
	{"B", "B", NewNode(B), "B", "a(bc)", 3},
	{"C", "C", NewNode(C), "C", "acb", 3},
	{"W", "W", NewNode(W), "W", "abb", 2},
	{"((SK)K)", "SKK", Apply(Apply(NewNode(S), NewNode(K)), NewNode(K)), "SKK", "a", 1},
	{"(B(CW))", "B(CW)", Apply(NewNode(B), Apply(NewNode(C), NewNode(W))), "B(CW)", "c(ab)(ab)", 3},
	{"((((IS)K)S)K)", "ISKSK", iotaK, "K", "a", 2},
	{"((((((IS)K)S)K)S)K)", "ISKSKSK", iotaS, "S", "ac(bc)", 3},
	{"(((((S(K(S(K(S(KI))))))S)K)S)K)", "S(K(S(K(S(KI)))))SKSK", jotK, "K", "a", 2},
	{"(((((((S(K(S(K(S(K(S(K(S(KI))))))))))S)K)S)K)S)K)", "S(K(S(K(S(K(S(K(S(KI)))))))))SKSKSK", jotS, "S", "ac(bc)", 3},
}

var validWithSpaces = []skiTest{
	{" S ", " S ", NewNode(S), "S", "ac(bc)", 3},
	{" ( K I ) ", " K I ", Apply(NewNode(K), NewNode(I)), "KI", "b", 2},
	{"	(	C	C	)	", "	C	C	", Apply(NewNode(C), NewNode(C)), "CC", "bca", 3},
}

func TestParseValidSKI(t *testing.T) {
	for _, test := range append(valid, validWithSpaces...) {
		if got, err := Parse(test.fs); err != nil || !reflect.DeepEqual(got, test.n) {
			t.Errorf("parseSKI(%v): got %#v, %v; want %#v, nil", test.fs, got, err, test.n)
		}
		if got, err := Parse(test.s); err != nil || !reflect.DeepEqual(got, test.n) {
			t.Errorf("parseSKI(%v): got %#v, %v; want %#v, nil", test.s, got, err, test.n)
		}
	}
}

var invalidSKI = []string{
	"",
	" ",
	"()",
	"(S)",
	"(",
	")",
	"((SK)K",
	"(C(BI)))",
	"Z",
	"s",
}

func TestParseInvalidSKI(t *testing.T) {
	for _, test := range invalidSKI {
		if got, err := Parse(test); err == nil {
			t.Errorf("parseSKI(%v): got %#v, nil; want nil, error", test, got)
		}
	}
}

var validIota = []struct {
	s string
	n *Node
}{
	{"i", leftIota(newNode(I))},
	{"*ii", newNode(I)},
	{"*i*i*ii", iotaK},
	{"*i*i*i*ii", iotaS},
	{" *    i  i   ", newNode(I)},
}

var invalidIota = []string{
	"*",
	"ii",
	"*i",
	"*i*",
	"i*ii",
	"*i*ii*",
	"iiii***",
	"****iiii",
	"*i*i*i*i*",
}

func TestParseValidIota(t *testing.T) {
	for _, test := range validIota {
		if got, err := Parse(test.s); err != nil || !reflect.DeepEqual(got, test.n) {
			t.Errorf("parseIota(%v): got %#v, %v; want %#v, nil", test.s, got, err, test.n)
		}
	}
}

func TestParseInvalidIota(t *testing.T) {
	for _, test := range invalidIota {
		if got, err := Parse(test); err == nil {
			t.Errorf("parseIota(%v): got %#v, nil; want nil, error", test, got)
		}
	}
}

var validJot = []struct {
	s string
	n *Node
}{
	{"11100", jotK},
	{"11111000", jotS},
	{"	1  1 1110	0		0  ", jotS},
}

func TestParseJot(t *testing.T) {
	for _, test := range validJot {
		if got, err := Parse(test.s); err != nil || !reflect.DeepEqual(got, test.n) {
			t.Errorf("parseJot(%v): got %#v, %v; want %#v, nil", test.s, got, err, test.n)
		}
	}
}

func ExampleParse() {
	for _, s := range []string{
		"((BC)(SC))",
		"WBWB",
		"***i*i*i*ii*iii",
		"110110010",
	} {
		n, err := Parse(s)
		if err != nil {
			panic(err)
		}
		fmt.Println(s, n.String())
	}
	// Output:
	// ((BC)(SC)) BC(SC)
	// WBWB WBWB
	// ***i*i*i*ii*iii S(K(ISKSKSKI))
	// 110110010 S(K(S(K(S(K(S(K(S(KI)))SK))))SKSK))SK
}

func TestFullString(t *testing.T) {
	for _, test := range valid {
		if got := test.n.FullString(); got != test.fs {
			t.Errorf("%#v.FullString(): got %v, want %v", test.n, got, test.fs)
		}
	}
}

func TestString(t *testing.T) {
	for _, test := range valid {
		if got := test.n.String(); got != test.s {
			t.Errorf("%#v.String(): got %v, want %v", test.n, got, test.s)
		}
	}
}

// TestSimplify must be executed after the Iota and Jot tests because it mutates iotaK, iotaS, jotK, and jotS.
func TestSimplify(t *testing.T) {
	for _, test := range valid {
		if got := Simplify(test.n).String(); got != test.simplify {
			t.Errorf("Simplify(%#v): got %v, want %v", test.n, got, test.simplify)
		}
	}
}

func ExampleSimplify() {
	for _, n := range []*Node{
		Apply(NewNode(I), NewNode(C)),
		Apply(Apply(Apply(NewNode(S), NewNode(K)), NewNode(S)), NewNode(K)),
		Apply(Apply(NewNode(W), Apply(NewNode(B), NewNode(S))), NewNode(C)),
	} {
		s := Simplify(n)
		fmt.Println(n, s)
	}
	// Output:
	// IC C
	// SKSK K
	// W(BS)C S(CC)
}

func TestReduce(t *testing.T) {
	for _, test := range valid {
		if got, n := Reduce(test.n); got.String() != test.reduce || n != test.reduceN {
			t.Errorf("Reduce(%#v): got %v, %v; want %v, %v", test.n, got.String(), n, test.reduce, test.reduceN)
		}
	}
}

func ExampleReduce() {
	for _, c := range []Comb{I, K, S} {
		r, n := Reduce(NewNode(c))
		fmt.Println(c, r, n)
	}
	// Output:
	// I a 1
	// K a 2
	// S ac(bc) 3
}
