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

// isSortedI8 returns 0 if ar is sorted in ascending
// order, otherwise it returns i > 0 with ar[i] < ar[i-1]
func isSortedI8(ar []int64) int {
	for i := len(ar) - 1; i > 0; i-- {
		if ar[i] < ar[i-1] {
			return i
		}
	}
	return 0
}

// pre-sort
func presortI8(ar []int64) {
	l, h := 0, (MaxLenIns+1)/3
	for h < len(ar) {
		if ar[h] < ar[l] {
			ar[h], ar[l] = ar[l], ar[h]
		}
		l++
		h++
	}
}

// insertion sort
func insertionI8(ar []int64) {
	for h := 0; h < len(ar)-1; {
		l := h
		h++
		v := ar[h]
		if v < ar[l] {
			for {
				ar[l+1] = ar[l]
				l--
				if l < 0 || v >= ar[l] {
					break
				}
			}
			ar[l+1] = v
		}
	}
}

// pre+insertion sort
func pinsertI8(ar []int64) {
	presortI8(ar)
	insertionI8(ar)
}

// pivotI8 selects 2n equidistant samples from ar that minimizes max distance to any
// non-selected member, calculates median-of-2n pivot from samples. ensures lo/hi ranges
// have at least n elements by moving sorted samples to n positions at lo/hi ends.
// assumes 5 > n > 0, len(ar) > 4n. returns remaining slice,pivot for partitioning.
func pivotI8(ar []int64, n int) ([]int64, int64) {

	// sample step, first/last sample positions
	d, s, h, l := minmaxSample(len(ar), n)

	var sample [8]int64
	for i, k := d, h; i >= 0; i-- {
		sample[i] = ar[k]
		k -= s
	}
	insertionI8(sample[:d+1]) // sort 2n samples

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
	return ar[lo:hi:hi], sixb.MeanI8(sample[n-1], sample[n])
}

// partition ar into <= and >= pivot, assumes len(ar) >= 2
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func partition1I8(ar []int64, pv int64) int {
	l, h := 0, len(ar)-1
	for l < h {
		if ar[h] < pv { // avoid unnecessary comparisons
			for {
				if pv < ar[l] {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
				l++
				if l >= h {
					return l + 1
				}
			}
		} else if pv < ar[l] { // extend ranges in balance
			for {
				h--
				if l >= h {
					return l
				}
				if ar[h] < pv {
					ar[l], ar[h] = ar[h], ar[l]
					break
				}
			}
		}
		l++
		h--
	}
	if l == h && ar[h] < pv { // classify mid element
		l++
	}
	return l
}

// rearrange ar[:a] and ar[b:] into <= and >= pivot, assumes 0 < a < b < len(ar)
// gap (a,b) expands until one of the intervals is fully consumed
func partition2I8(ar []int64, a, b int, pv int64) (int, int) {
	a--
	for a >= 0 && b < len(ar) {
		if ar[b] < pv { // avoid unnecessary comparisons
			for {
				if pv < ar[a] {
					ar[a], ar[b] = ar[b], ar[a]
					break
				}
				a--
				if a < 0 {
					return a, b
				}
			}
		} else if pv < ar[a] { // extend ranges in balance
			for {
				b++
				if b >= len(ar) {
					return a, b
				}
				if ar[b] < pv {
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
func gpart1I8(ar []int64, pv int64, ch chan int) {
	ch <- partition1I8(ar, pv)
}

// concurrent dual partitioning of ar
// returns k with ar[:k] <= pivot, ar[k:] >= pivot
func cdualparI8(ar []int64, ch chan int) int {

	aq, pv := pivotI8(ar, 4) // median-of-8 pivot
	k := len(aq) >> 1
	a, b := k>>1, sixb.MeanI(k, len(aq))

	go gpart1I8(aq[a:b:b], pv, ch) // mid half range

	k = a
	a, b = partition2I8(aq, a, b, pv) // left/right quarter ranges

	k += <-ch // convert returned indice to aq

	// only one gap is possible
	for ; 0 <= a; a-- { // gap left in low range?
		if pv < aq[a] {
			k--
			aq[a], aq[k] = aq[k], aq[a]
		}
	}
	for ; b < len(aq); b++ { // gap left in high range?
		if aq[b] < pv {
			aq[b], aq[k] = aq[k], aq[b]
			k++
		}
	}
	return k + 4 // convert k indice to ar
}

// short range sort function, assumes MaxLenIns < len(ar) <= MaxLenRec
func shortI8(ar []int64) {
start:
	aq, pv := pivotI8(ar, 2) // median-of-4 pivot
	k := partition1I8(aq, pv)

	k += 2 // convert k indice from aq to ar

	if k < len(ar)-k {
		aq = ar[:k:k]
		ar = ar[k:] // ar is the longer range
	} else {
		aq = ar[k:]
		ar = ar[:k:k]
	}

	if len(aq) > MaxLenIns {
		shortI8(aq) // recurse on the shorter range
		goto start
	}
psort:
	presortI8(aq)
	insertionI8(aq) // at least one insertion range

	if len(ar) > MaxLenIns {
		goto start
	}
	if &ar[0] != &aq[0] {
		aq = ar
		goto psort // two insertion ranges
	}
}

// new-goroutine sort function
func glongI8(ar []int64, sv *syncVar) {
	longI8(ar, sv)

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) == 0 { // decrease goroutine counter
		sv.done <- 0 // we are the last, all done
	}
}

// long range sort function, assumes len(ar) > MaxLenRec
func longI8(ar []int64, sv *syncVar) {
start:
	aq, pv := pivotI8(ar, 3) // median-of-6 pivot
	k := partition1I8(aq, pv)

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

		if len(aq) > MaxLenIns {
			shortI8(aq)
		} else {
			pinsertI8(aq)
		}

		if len(ar) > MaxLenRec { // two not-long ranges?
			goto start
		}
		shortI8(ar) // we know len(ar) > MaxLenIns
		return
	}

	// max goroutines? not atomic but good enough
	if sv == nil || gorFull(sv) {
		longI8(aq, sv) // recurse on the shorter range
		goto start
	}

	if atomic.AddUint32(&sv.ngr, 1) == 0 { // increase goroutine counter
		panic("sorty: longI8: counter overflow")
	}
	// new-goroutine sort on the longer range only when
	// both ranges are big and max goroutines is not exceeded
	go glongI8(ar, sv)
	ar = aq
	goto start
}

// sortI8 concurrently sorts ar in ascending order.
func sortI8(ar []int64) {

	if len(ar) < 2*(MaxLenRec+1) || MaxGor <= 1 {

		if len(ar) > MaxLenRec { // single-goroutine sorting
			longI8(ar, nil)
		} else if len(ar) > MaxLenIns {
			shortI8(ar)
		} else {
			pinsertI8(ar)
		}
		return
	}

	// create channel only when concurrent partitioning & sorting
	sv := syncVar{1, // number of goroutines including this
		make(chan int)} // end signal
	for {
		// concurrent dual partitioning with done
		k := cdualparI8(ar, sv.done)
		var aq []int64

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
				panic("sorty: sortI8: counter overflow")
			}
			go glongI8(aq, &sv)

		} else if len(aq) > MaxLenIns {
			shortI8(aq)
		} else {
			pinsertI8(aq)
		}

		// longer range big enough? max goroutines?
		if len(ar) < 2*(MaxLenRec+1) || gorFull(&sv) {
			break
		}
		// dual partition longer range
	}

	longI8(ar, &sv) // we know len(ar) > MaxLenRec

	if atomic.AddUint32(&sv.ngr, ^uint32(0)) != 0 { // decrease goroutine counter
		<-sv.done // we are not the last, wait
	}
}
