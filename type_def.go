package main

type Dir int

const (
	Dir_Nort = iota
	Dir_East
	Dir_West
	Dir_South
	DIR_COUNT
)

type pos struct {
	x int
	y int
}