package main

type Dir int

const (
	Dir_Nort = iota
	Dir_East
	Dir_West
	Dir_South
	DIR_COUNT
)

type Ahead struct {
	front string
	back  string
	left  string
	right string
}

type pos struct {
	x int
	y int
}

type cell_msg struct {
	dir   Dir
	pos   pos
	block bool
}

type cell_to_player_msg struct {
	dir Dir
	pos pos
}
