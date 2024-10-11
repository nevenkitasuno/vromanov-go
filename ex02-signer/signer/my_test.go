package main

import (
	"testing"
)

func TestSingleHash(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{"Hello", "4157704578~306818087"},
		{"world", "980881731~980913668"},
	}

	for _, testCase := range testCases {
		in := make(chan interface{})
		out := make(chan interface{})

		go func() {
			SingleHash(in, out)
		}()

		in <- testCase.input

		result := <-out
		if result != testCase.expect {
			t.Errorf("Expected: %s, Got: %s", testCase.expect, result)
		}
	}
}

func TestSingleHashZero(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{"0", "4108050209~502633748"},
	}

	for _, testCase := range testCases {
		in := make(chan interface{})
		out := make(chan interface{})

		go func() {
			SingleHash(in, out)
		}()

		in <- testCase.input

		result := <-out
		if result != testCase.expect {
			t.Errorf("Expected: %s, Got: %s", testCase.expect, result)
		}
	}
}

func TestMultiHash(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{"Hello_world", "1977180353400111531025817903023825310519943495023984201681"},
		{"1234567890_abcdefghij", "300603355113423131106916207568068828312967784591516445723"},
	}

	for _, testCase := range testCases {
		in := make(chan interface{})
		out := make(chan interface{})

		go func() {
			MultiHash(in, out)
		}()

		in <- testCase.input

		result := <-out
		if result != testCase.expect {
			t.Errorf("Expected: %s, Got: %s", testCase.expect, result)
		}
	}
}

func TestMultiHashFromHash(t *testing.T) {
	testCases := []struct {
		input  string
		expect string
	}{
		{"4108050209~502633748", "29568666068035183841425683795340791879727309630931025356555"},
	}

	for _, testCase := range testCases {
		in := make(chan interface{})
		out := make(chan interface{})

		go func() {
			MultiHash(in, out)
		}()

		in <- testCase.input

		result := <-out
		if result != testCase.expect {
			t.Errorf("Expected: %s, Got: %s", testCase.expect, result)
		}
	}
}

func TestCombineResults(t *testing.T) {
	testCases := []struct {
		input1 string
		input2 string
		expect string
	}{
		{"Hello", "world", "Hello_world"},
		{"abc", "123", "123_abc"},
	}

	for _, testCase := range testCases {
		in := make(chan interface{}, 1)
		out := make(chan interface{})

		go func(out chan<- interface{}) {
			out <- testCase.input1
			out <- testCase.input2
			close(out)
		}(in)

		go func() {
			CombineResults(in, out)
		}()

		result := <-out
		if result != testCase.expect {
			t.Errorf("Expected: %s, Got: %s", testCase.expect, result)
		}
	}
}
