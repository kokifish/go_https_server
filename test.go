package main

import (
	"fmt"
)

func test() {
	// %d 表示整型数字，%s 表示字符串
	// const identifier [type] = value
	const i_const int = 999
	var ii int = 123
	var iii *int = &ii
	var enddate = "2020-12-31"
	var target_url = fmt.Sprintf("Code=%d&endDate=%s", ii, enddate)
	fmt.Println(target_url)
	fmt.Println(iii, i_const)
	var a int
	var f float64
	var b bool
	var s string
	fmt.Printf("%v %v %v %q\n", a, f, b, s)
	const (
		i = 1 << iota
		j = 3 << iota // 0b11 << 1 = 0b110
		k             //3 << iota // 0b11 << 2 = 0b1100
		l             //3 << iota // 0b11 << 3 = 0b11000 //24
	)

	fmt.Println("i=", i)
	fmt.Println("j=", j)
	fmt.Println("k=", k)
	fmt.Println("l=", l)
	//go f(x, y, z) //goroutine
}
