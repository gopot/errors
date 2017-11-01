//   Copyright © 2015-2017 Ivan Kostko (github.com/ivan-kostko; github.com/gopot)

//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at

//       http://www.apache.org/licenses/LICENSE-2.0

//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package errors_test

import (
	"fmt"
	"testing"

	e "errors"

	"github.com/gopot/errors"
)

type CriticalError interface {
	error
	IsCritical() bool
}

func NewCriticalError(msg string, isCritical bool) *customErrorOneDetail {
	return &customErrorOneDetail{
		message:    msg,
		isCritical: isCritical,
	}
}

type customErrorOneDetail struct {
	message    string
	isCritical bool
}

func (this *customErrorOneDetail) Error() string {
	ret := this.message
	if this.isCritical {
		ret += " IS CRITICAL"
	}
	return ret
}

func (this *customErrorOneDetail) IsCritical() bool {
	return this.isCritical
}

func NewCriticalErrorTenDetails(msg string, isCritical bool) *customErrorTenDetails {
	return &customErrorTenDetails{
		message:    msg,
		isCritical: isCritical,
	}
}

type customErrorTenDetails struct {
	message    string
	isCritical bool
	one        string
	two        string
	three      string
	four       string
	five       string
	six        string
	seven      string
	eight      string
	nine       string
}

func (this *customErrorTenDetails) Error() string {
	ret := this.message
	if this.isCritical {
		ret += " IS CRITICAL"
	}
	return ret
}

func (this *customErrorTenDetails) IsCritical() bool {
	return this.isCritical
}

func BenchmarkErrorCycle(b *testing.B) {

	details := []struct{ Key, Value interface{} }{
		{"One", 1},
		{"Two", 2},
		{"Three", 3},
		{"Four", 4},
		{"Five", 5},
		{"Six", 6},
		{"Seven", 7},
		{"Eight", 8},
		{"Nine", 9},
		{"Zero", 0},
	}

	benchCases := []struct {
		Alias        string
		ErrorFactory func(msg string) error
	}{
		{
			Alias:        "Golang Standard errors.New",
			ErrorFactory: e.New,
		},
		{
			Alias:        "fmt.Errorf error",
			ErrorFactory: func(msg string) error { return fmt.Errorf(msg) },
		},
		{
			Alias:        "Custom error one detail",
			ErrorFactory: func(msg string) error { return NewCriticalError(msg, true) },
		},
		{
			Alias:        "Custom error ten details",
			ErrorFactory: func(msg string) error { return NewCriticalErrorTenDetails(msg, true) },
		},
		{
			Alias:        "Gopot errors.New",
			ErrorFactory: func(msg string) error { return errors.New(msg) },
		},
		{
			Alias:        "Gopot errors.New With 1 Detail",
			ErrorFactory: func(msg string) error { return errors.NewWithDetails(msg, details[:1]...) },
		},
		{
			Alias:        "Gopot errors.New With 10 Detail",
			ErrorFactory: func(msg string) error { return errors.NewWithDetails(msg, details...) },
		},
	}

	for _, bCase := range benchCases {

		benchmark := func(b *testing.B) {

			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				err := bCase.ErrorFactory("Some error")
				_ = err.Error()
			}
		}
		b.Run(bCase.Alias, benchmark)
	}
}

func BenchmarkErrorIsCriticalCycle(b *testing.B) {

	details := []struct{ Key, Value interface{} }{
		{"One", 1},
		{"Two", 2},
		{"Three", 3},
		{"Four", 4},
		{"Five", 5},
		{"Six", 6},
		{"Seven", 7},
		{"Eight", 8},
		{"Nine", 9},
		{"Zero", 0},
	}

	type errorIsCritical string

	const ErrorISCritical = errorIsCritical("I Critical")

	benchCases := []struct {
		Alias        string
		ErrorFactory func(msg string) error
		IsCritical   func(err error) bool
	}{
		{
			Alias: "Custom error	one detail          	                         ",
			ErrorFactory: func(msg string) error { return NewCriticalError(msg, true) },
			IsCritical:   func(err error) bool { return (err.(CriticalError)).IsCritical() },
		},
		{
			Alias: "Custom error	ten details         	                         ",
			ErrorFactory: func(msg string) error { return NewCriticalErrorTenDetails(msg, true) },
			IsCritical:   func(err error) bool { return (err.(CriticalError)).IsCritical() },
		},
		{
			Alias: "Gopot errors.New	                  	                     ",
			ErrorFactory: func(msg string) error { return errors.New(msg) },
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
		{
			Alias: "Gopot errors.New With 1 Detail	                          ",
			ErrorFactory: func(msg string) error {
				return errors.NewWithDetails(msg, struct{ Key, Value interface{} }{Key: ErrorISCritical})
			},
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
		{
			Alias:        "Gopot errors.New With 10 Detail No Critical key          ",
			ErrorFactory: func(msg string) error { return errors.NewWithDetails(msg, details...) },
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
		{
			Alias: "Gopot errors.New With 10 Detail the first is Critical key",
			ErrorFactory: func(msg string) error {
				return errors.NewWithDetails(msg, append([]struct{ Key, Value interface{} }{{Key: ErrorISCritical}}, details[:9]...)...)
			},
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
		{
			Alias: "Gopot errors.New With 10 Detail the last is Critical key ",
			ErrorFactory: func(msg string) error {
				return errors.NewWithDetails(msg, append(details[:9], struct{ Key, Value interface{} }{Key: ErrorISCritical})...)
			},
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
	}

	for _, bCase := range benchCases {

		benchmark := func(b *testing.B) {

			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				err := bCase.ErrorFactory("Some error")
				_ = bCase.IsCritical(err)
			}
		}
		b.Run(bCase.Alias, benchmark)
	}
}

func BenchmarkErrorIsCriticalLookup(b *testing.B) {

	details := []struct{ Key, Value interface{} }{
		{"One", 1},
		{"Two", 2},
		{"Three", 3},
		{"Four", 4},
		{"Five", 5},
		{"Six", 6},
		{"Seven", 7},
		{"Eight", 8},
		{"Nine", 9},
		{"Zero", 0},
	}

	type errorIsCritical string

	const ErrorISCritical = errorIsCritical("I Critical")

	msg := "Some error"

	benchCases := []struct {
		Alias      string
		Error      error
		IsCritical func(err error) bool
	}{
		{
			Alias: "Custom error	 one detail       	                         ",
			Error:      NewCriticalError(msg, true),
			IsCritical: func(err error) bool { return (err.(CriticalError)).IsCritical() },
		},
		{
			Alias: "Custom error	 ten details       	                        ",
			Error:      NewCriticalErrorTenDetails(msg, true),
			IsCritical: func(err error) bool { return (err.(CriticalError)).IsCritical() },
		},
		{
			Alias: "Gopot errors.New	                  	                     ",
			Error: errors.New(msg),
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
		{
			Alias: "Gopot errors.New With 1 Detail	                          ",
			Error: errors.NewWithDetails(msg, struct{ Key, Value interface{} }{Key: ErrorISCritical}),
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
		{
			Alias: "Gopot errors.New With 10 Detail No Critical key          ",
			Error: errors.NewWithDetails(msg, details...),
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
		{
			Alias: "Gopot errors.New With 10 Detail the first is Critical key",
			Error: errors.NewWithDetails(msg, append([]struct{ Key, Value interface{} }{{Key: ErrorISCritical}}, details[:9]...)...),
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
		{
			Alias: "Gopot errors.New With 10 Detail the last is Critical key ",
			Error: errors.NewWithDetails(msg, append(details[:9], struct{ Key, Value interface{} }{Key: ErrorISCritical})...),
			IsCritical: func(err error) bool {
				_, found := (err.(errors.DetailedError)).Get(ErrorISCritical)
				return found
			},
		},
	}

	for _, bCase := range benchCases {

		benchmark := func(b *testing.B) {

			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_ = bCase.IsCritical(bCase.Error)
			}
		}
		b.Run(bCase.Alias, benchmark)
	}
}
