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
	front Dir
	back  Dir
	left  Dir
	right Dir
}

type pos struct {
	x int
	y int
}

type cell_msg struct {
	dir   Dir
	pos   pos
	block int
}

type cell_to_player_msg struct {
	dir Dir
	pos pos
}

func create_map(width int, heigth int) [][]int {

	cell_map := make([][]int, heigth)

	for i := 0; i < heigth; i++ {

		cell_map[i] = make([]int, width)

		for j := 0; j < width; j++ {
			//cell의 원자모델 들 생성

			cell_map[i][j] = 0
		}

	}

	return cell_map
}
