package main

import (
	"fmt"

	"github.com/jfcg/sorty/v2"
)

type tomato struct {
	Quantity int
	Period   int
}

func (t *tomato) Next_day() {
	if t.Period > 0 {
		t.Period -= 1
	} else {
		fmt.Println("??")
	}

}

func remove_tomato(slice *[]tomato, s int) []tomato {
	return append((*slice)[:s], (*slice)[s+1:]...)
}

func Sort_tomato(b *[]tomato) {
	lsw := func(i, k, r, s int) bool {
		if (*b)[i].Period < (*b)[k].Period {
			if r != s {
				(*b)[r], (*b)[s] = (*b)[s], (*b)[r]
			}
			return true
		}
		return false
	}
	sorty.Sort(len((*b)), lsw)
}

func Sales(n int, t *[]tomato) *[]tomato {

	if (*t)[0].Quantity < n {
		rest := n - (*t)[0].Quantity
		(*t) = remove_tomato(t, 0)
		Sales(rest, t)
	} else if (*t)[0].Quantity == n {
		// t = remove_tomato(t, 0)
		return t
	} else {
		(*t)[0].Quantity = (*t)[0].Quantity - n
		return t
	}
	return nil
}

func total_tomato(t *[]tomato) int {
	var total = 0
	for _, v := range *t {
		total += v.Quantity
	}
	return total
}
