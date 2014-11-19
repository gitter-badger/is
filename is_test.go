package is

import (
	"errors"
	"strings"
	"testing"
)

type customErr struct{}

func (e *customErr) Error() string {
	return "Oops"
}

type mockT struct {
	failed bool
}

func (m *mockT) FailNow() {
	m.failed = true
}
func (m *mockT) Failed() bool {
	return m.failed
}

func TestIs(t *testing.T) {

	for _, test := range []struct {
		N     string
		F     func(is I)
		Fails []string
	}{
		// is.Nil
		{
			N: "Nil(nil)",
			F: func(is I) {
				is.Nil(nil)
			},
		},
		{
			N: "Nil(\"nope\")",
			F: func(is I) {
				is.Nil("nope")
			},
			Fails: []string{"expected nil: \"nope\""},
		},
		// is.OK
		{
			N: "OK(false)",
			F: func(is I) {
				is.OK(false)
			},
			Fails: []string{"unexpected false"},
		}, {
			N: "OK(true)",
			F: func(is I) {
				is.OK(true)
			},
		}, {
			N: "OK(nil)",
			F: func(is I) {
				is.OK(nil)
			},
			Fails: []string{"unexpected nil"},
		}, {
			N: "OK(1,2,3)",
			F: func(is I) {
				is.OK(1, 2, 3)
			},
		}, {
			N: "OK(0)",
			F: func(is I) {
				is.OK(0)
			},
			Fails: []string{"unexpected zero"},
		}, {
			N: "OK(1)",
			F: func(is I) {
				is.OK(1)
			},
		}, {
			N: "OK(\"\")",
			F: func(is I) {
				is.OK("")
			},
			Fails: []string{"unexpected \"\""},
		},
		// NoErr
		{
			N: "NoErr(errors.New(\"an error\"))",
			F: func(is I) {
				is.NoErr(errors.New("an error"))
			},
			Fails: []string{"unexpected error: an error"},
		}, {
			N: "NoErr(&customErr{})",
			F: func(is I) {
				is.NoErr(&customErr{})
			},
			Fails: []string{"unexpected error: Oops"},
		}, {
			N: "NoErr(error(nil))",
			F: func(is I) {
				var err error
				is.NoErr(err)
			},
		},
		{
			N: "NoErr(err1, err2, err3)",
			F: func(is I) {
				is.NoErr(&customErr{}, &customErr{}, &customErr{})
			},
			Fails: []string{"unexpected error: Oops"},
		},
		{
			N: "NoErr(err1, err2, err3)",
			F: func(is I) {
				var err1 error
				var err2 error
				var err3 error
				is.NoErr(err1, err2, err3)
			},
		},
		// OK
		{
			N: "OK(customErr(nil))",
			F: func(is I) {
				var err *customErr
				is.NoErr(err)
			},
		}, {
			N: "OK(func) panic",
			F: func(is I) {
				is.OK(func() {
					panic("panic message")
				})
			},
			Fails: []string{"unexpected panic: panic message"},
		}, {
			N: "OK(func) no panic",
			F: func(is I) {
				is.OK(func() {})
			},
		},
		// is.Panic
		{
			N: "PanicWith(\"panic message\", func(){ panic() })",
			F: func(is I) {
				is.PanicWith("panic message", func() {
					panic("panic message")
				})
			},
		},
		{
			N: "PanicWith(\"panic message\", func(){ /* no panic */ })",
			F: func(is I) {
				is.PanicWith("panic message", func() {
				})
			},
			Fails: []string{"expected panic: \"panic message\""},
		},
		{
			N: "Panic(func(){ panic() })",
			F: func(is I) {
				is.Panic(func() {
					panic("panic message")
				})
			},
		},
		{
			N: "Panic(func(){ /* no panic */ })",
			F: func(is I) {
				is.Panic(func() {
				})
			},
			Fails: []string{"expected panic"},
		},
		// is.Equal
		{
			N: "Equal(1,1)",
			F: func(is I) {
				is.Equal(1, 1)
			},
		}, {
			N: "Equal(1,2)",
			F: func(is I) {
				is.Equal(1, 2)
			},
			Fails: []string{"1 != 2"},
		}, {
			N: "Equal(1,nil)",
			F: func(is I) {
				is.Equal(1, nil)
			},
			Fails: []string{"1 != <nil>"},
		}, {
			N: "Equal(nil,1)",
			F: func(is I) {
				is.Equal(nil, 1)
			},
			Fails: []string{"<nil> != 1"},
		}, {
			N: "Equal(false,false)",
			F: func(is I) {
				is.Equal(false, false)
			},
		}, {
			N: "Equal(map1,map2)",
			F: func(is I) {
				is.Equal(
					map[string]interface{}{"package": "is"},
					map[string]interface{}{"package": "is"},
				)
			},
		}} {

		tt := new(mockT)
		is := New(tt)
		var rec interface{}

		func() {
			defer func() {
				rec = recover()
			}()
			test.F(is)
		}()

		if len(test.Fails) > 0 {
			for n, fail := range test.Fails {
				if !tt.Failed() {
					t.Errorf("%s should fail", test.N)
				}
				if test.Fails[n] != fail {
					t.Errorf("expected fail \"%s\" but was \"%s\".", test.Fails[n], fail)
				}
			}
		} else {
			if tt.Failed() {
				t.Errorf("%s shouldn't fail but: %s", test.N, strings.Join(test.Fails, ", "))
			}
		}

	}

}

func TestNewStrict(t *testing.T) {
	tt := new(mockT)
	is := Relaxed(tt)

	is.OK(nil)
	is.Equal(1, 2)
	is.NoErr(errors.New("nope"))

	if tt.Failed() {
		t.Error("Relaxed should not call FailNow")
	}

}
