package main

type tomato struct {
	Quantity int
	Period   int
}

func (t *tomato) Next_day() {
	t.Period--
}

func Sales(n int, t []tomato) []tomato {
	if t[0].Quantity < n {
		rest := n - t[0].Quantity
		t = remove(t, 0)
		Sales(rest, t)
	} else if t[0].Quantity == n {
		t = remove(t, 0)
		return t
	} else {
		t[0].Quantity = t[0].Quantity - n
		return t
	}
	return nil
}

func remove(slice []tomato, s int) []tomato {
	return append(slice[:s], slice[s+1:]...)
}
