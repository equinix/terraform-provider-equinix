package main

import (
	"fmt"
)

func main() {
	//es := ProviderStatusNew("PROVISIONED")
	fmt.Println(getEnum())
	fmt.Println(*getEnum())

}

func getEnum() *ProviderStatusNew {
	es := ProviderStatusNew("PROVISIONED")
	return &es
}

func testErr() {

	callErr()

}

func callErr() error {

	return &MyError{}
}

type MyError struct{}

func (m *MyError) Error() string {
	return "boom"
}
