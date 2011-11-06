/*
Copyright Mathieu Lonjaret

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package histo

import (
	"sort"
)

type Bar struct {
	Value int64
	Count int64
	Min   int64
	Max   int64
	
}

type sortable []int64 

func (s sortable) Swap(i, j int)      {s[i], s[j] = s[j], s[i] }
func (s sortable) Len() int           { return len(s) }
func (s sortable) Less(i, j int) bool { return s[i] < s[j] }

type Histo struct {
	bar [](*Bar)
	unsorted sortable
}

// NewHisto returns an histogram set up with n bins
func NewHisto() *Histo {
	return &Histo{nil, nil}
}


func (h Histo) sort() {
	if h.unsorted != nil {
		sort.Sort(h.unsorted)
	}
}

func (h Histo) Add(v int64) {
	h.unsorted = append(h.unsorted, v)
}

func (h Histo) bin(num int) {
	if h.unsorted == nil {
		return
	}
	h.sort()
	max := h.unsorted[len(h.unsorted) - 1]
	min := h.unsorted[0]
	binWidth := 1 + (max - min)/int64(num)
	np := int64(0)
	average := int64(0)
	sup := min + binWidth
	for _, v := range h.unsorted {
		if v > sup {
			average /= np
			h.bar = append(h.bar, &Bar{average, np, sup, sup + binWidth})
			sup += binWidth
			average = v
			np = 1
		} else {
			average += v
			np++
		}
	}
	h.unsorted = nil
}

func (h Histo) Bin(num int) {
	h.bin(num)
}

func (h Histo) Bars(num int) [](*Bar) {
	h.bin(num)
	return h.bar
}

/*
func (h Histo) At(num int) (hb *Bar) {
	if h.bar != nil {
		if num > len(h.bar - 1) {
			return nil
		}
		return h.bar[num]
	}
	return nil 
}

func (h Histo) Value(value int64) (hb *Bar) {
	for _, v := range h {
		if v.Value == value {
			hb = v
			break
		}
	}
	return hb
}

func (h Histo) MaxCount() (mc int64) {
	for _, v := range h {
		if v.Count > mc {
			mc = v.Count
		}
	}
	return mc
}

func (h Histo) MaxValue() int64 {
	return h[len(h)-1].Value
}

func (h Histo) MinValue() int64 {
	return h[0].Value
}


func (h Histo) ReBin(nb int) *Histo {
	binWidth := 1 + (h.MaxValue()-h.MinValue())/int64(nb)
	nh := make(Histo, 0, 1)
	np := int64(0)
	average := int64(0)
	max := h.MinValue() + binWidth
	for _, v := range h {
		if v.Value > max {
			average /= np
			nh = append(nh, &Bar{average, np, max, max + binWidth})
			max += binWidth
			average = v.Value
			np = 1
		} else {
			average += v.Value
			np++
		}
	}
	return &nh
}
*/