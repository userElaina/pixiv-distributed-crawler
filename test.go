package main

import (
	"./pic"
)

func main() {
	ch := make(chan int)
	// mian := pic.DefaultMuti("./tmp", ch)
	mian := pic.NewMuti(2, 4, 5, "./tmp", ch)
	mian.MutiThread()
	ch <- 86498318
	ch <- 87011701
	ch <- 100480461
	ch <- 101034134
	ch <- -1
	mian.Join()
}

/*
86498318 mange with series
87011701 muti page illustration
100480461 manga without series
101034134 one page illustration
*/
