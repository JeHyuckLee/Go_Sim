/*	Copyright (c) 2019-present, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

import (
	"sync/atomic"

	"github.com/jfcg/sixb"
)

// isSortedB returns 0 if ar is sorted in ascending lexicographic
// order, otherwise it returns i > 0 with sixb.BtoS(ar[i]) < sixb.BtoS(ar[i-1])
func isSortedB(ar [][]byte) int {
	for i := len(ar) - 1; i > 0; i-- {
		if sixb.BtoS(ar[i]) < sixb.BtoS(ar[i-1]) {
			return i
		}
	}
	return 0
}

// pre-sort
func presortB(ar [][]byte) {
	l, h := 0, (MaxLenInsFC+1)/3
	for h < len(ar) {
		if sixb.BtoS(ar[h]) < sixb.BtoS(ar[l]) {
			ar[h], ar[l] = ar[l], ar[h]
		}
		l++
		h++
	}
}

// insertion sort
func insertionB(ar [][]byte) {
	for h := 0; h < len(ar)-1; {
		l := h
		h++
		v := ar[h]
		if sixb.BtoS(v) < sixb.BtoS(ar[l]) {
			for {
				ar[l+1] = ar[l]
				l--
				if l < 0 || sixb.BtoS(v) >= sixb.BtoS(ar[l]) {
					break
				}
			}
			ar[l+1] = v
		}
	}
}

// pre+insertion sort
func pinsertB(ar [][]byte) {
	presortB(ar)
	insertionB(ar)
}

// pivotB selects 2n equidistant samples from ar that minimizes max distance to any
// non-selected member, calculates median-of-2n pivot from samples. ensures lo/hi ranges
// have at least n elements by moving sorted samples to n positions at lo/hi ends.
// assumes 5 > n > 0, len(ar) > 4n. returns remaining slice,pivot for partitioning.
func pivotB(ar [][]byte, n int) ([][]byte, string) {

	// sample step, first/last sample positions
	d, s, h, l := minmaxSample(len(ar), n)

	var sample [8][]byte
	for i, k := d, h; i >= 0; i-- {
		sample[i] = ar[k]
		k -= s
	}
	insertionB(sample[:d+1]) // sort 2n samples

	lo, hi := 0, len(ar)
	for { // move sorted samples to lo/hi ends
		hi--
		ar[h] = ar[hi]
		ar[hi] = sample[d]
		ar[l] = ar[lo]
		ar[lo] = sample[lo]
		l += s
		h -= s
		lo++
		d--
		if d < lo {
			break
		}
	}
	return ar[lo:hi:hi], sixb.MeanS(sixb.BtoS(sample[n-1]), sixb.BtoS(sample[n]))
}

// partition ar into <= and >= pivot, assumes len(ar) >= 2
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func partition1B(ar [][]byte, pv string) int {
	l, h := 0, len(ar)-1
	for l < h {
		if sixb.BtoS(ar[h]) < pv { // avoid unnecessary comparisons
			for {
				if pv < sixb.BtoS(ar[l]) {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if pv < sixb.BtoS(ar[l]) { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l
				}
				if sixb.BtoS(ar[h]) < pv {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
			}
		}
		l++
		h--
	}
	if l == h && sixb.BtoS(ar[h]) < pv { // classify mid element
		l++
	}
	return l
}

// rearrange ar[:a] and ar[b:] into <= and >= pivot, assumes 0 < a < b < len(ar)
// gap (a,b) expands until one of the intervals is fully consumed
func partition2B(ar [][]byte, a, b int, pv string) (int, int) {
	a--
	for a >= 0 && b < len(ar) {
		if sixb.BtoS(ar[b]) < pv { // avoid unnecessary comparisons
			for {
				if pv < sixb.BtoS(ar[a]) {
					ar[a], ar[b] = ar[b], ar[a]
					break
				}
				a--
				if a < 0 {
					return a, b
				}
			}
		} else if pv < sixb.BtoS(ar[a]) { // extend ranges in balance
			for {
				b++
				if b >= len(ar) {
					return a, b
				}
				if sixb.BtoS(ar[b]) < pv {
					ar[a], ar[b] = ar[b], ar[a]
					break
				}
			}
		}
		a--
		b++
	}
	return a, b
}

// new-goroutine partition
func gpart1B(ar [][]byte, pv string, ch chan int) {
	ch <- partition1B(ar, pv)
}

// concurrent dual partitioning of ar
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func cdualparB(ar [][]byte, ch chan int) int {

	aq, pv := pivotB(ar, 4) // median-of-8 pivot
	k := len(aq) >> 1
	a, b := k>>1, sixb.MeanI(k, len(aq))

	go gpart1B(aq[a:b:b], pv, ch) // mid half range

	k = a
	a, b = partition2B(aq, a, b, pv) // left/right quarter ranges

	k += <-ch // convert returned indice to aq

	// only one gap is possible
	for ; 0 <= a; a-- { // gap left in low range?
		if pv < sixb.BtoS(aq[a]) {
			k--
			aq[a], aq[k] = aq[k], aq[a]
		}
	}
	for ; b < len(aq); b++ { // gap left in high range?
		if sixb.BtoS(aq[b]) < pv {
			aq[b], aq[k] = aq[k], aq[b]
			k++
		}
	}
	return k + 4 // convert k indice to ar
}

// short range sort function, assumes MaxLenInsFC < len(ar) <= MaxLenRec
func shortB(ar [][]byte) {
start:
	aq, pv := pivotB(ar, 2) // median-of-4 pivot
	k := partition1B(aq, pv)

	k += 2 // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenInsFC {
		shortB(aq) // recurse on the shorter range
		goto start
	}
psort:
	presortB(aq)
	insertionB(aq) // at least one insertion range

	if len(ar) > MaxLenInsFC {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto psort // two insertion ranges
	}
}

// new-goroutine sort function
func glongB(ar [][]byte, sv *syncVar) {
	longB(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec
func longB(ar [][]byte, sv *syncVar) {
start:
	aq, pv := pivotB(ar, 3) // median-of-6 pivot
	k := partition1B(aq, pv)

	k += 3 // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	// branches below are optimal for fewer total jumps
	if len(aq) <= MaxLenRec { // at least one not-long range?

		if len(aq) > MaxLenInsFC {
			shortB(aq)
		} else {
			pinsertB(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortB(ar) // we know len(ar) > MaxLenInsFC
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longB(aq, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: longB: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glongB(ar, sv)
	ar = aq
	goto start
}

// sortB concurrently sorts ar in ascending lexicographic order.
func sortB(ar [][]byte) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longB(ar, nil)
		} else if len(ar) > MaxLenInsFC {
			shortB(ar)
		} else {
			pinsertB(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := cdualparB(ar, sv.done)
		var aq [][]byte

		if k < len(ar)-k {
			aq = ar[:k:k]
			ar = ar[k:] // ar is the longer range
		} else {
			aq = ar[k:]
			ar = ar[:k:k]
		}

		// handle shorter range
		if len(aq) > MaxLenRec {
			if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
				panic("sorty: sortB: counter overflow")
			}
			go glongB(aq, &sv)

		} else if len(aq) > MaxLenInsFC {
			shortB(aq)
		} else {
			pinsertB(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longB(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
