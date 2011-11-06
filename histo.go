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
	Points []int64
}

type sortable []int64 

func (s sortable) Swap(i, j int)      {s[i], s[j] = s[j], s[i] }
func (s sortable) Len() int           { return len(s) }
func (s sortable) Less(i, j int) bool { return s[i] < s[j] }

type Histo struct {
	nb int // number of bins/bars
	np int // number of points
	bar [](*Bar)
	unsorted sortable // pool of points yet to be sorted
}

// NewHisto returns an histogram set up with n bins
func NewHisto(num int) *Histo {
//	bar := make([](*Bar), 0, 1)
//	unsorted := make([]int64, 0, 1)
	return &Histo{num, 0, nil, nil}
}

func (h *Histo) sort() {
	if h.unsorted != nil {
		sort.Sort(h.unsorted)
	}
}

func (h *Histo) Add(v int64) {
	h.unsorted = append(h.unsorted, v)
}

func (h *Histo) resize(min int64, max in64) {
	binWidth := 1 + (max - min)/int64(h.num)
	np := int64(h.np)
	average := int64(0)
	sup := min + binWidth
	for _, v := range h.bar {
		
	}
}

func (h *Histo) distribute() {
	if h.unsorted == nil {
		// no new points; nothing to do
		return
	}

	rebin := false
	max := h.unsorted[len(h.unsorted) - 1]
	min := h.unsorted[0]
	if h.bar != nil {
		if h.bar[len(h.bar) - 1].Max > max {
			max = h.bar[len(h.bar) - 1].Max
		} else {
			rebin = true
		}
		if h.bar[0].Min < min {
			min = h.bar[0].Min
		} else {
			rebin = true
		}
	} 
	binWidth := 1 + (max - min)/int64(h.num)
	np := int64(h.np)
	average := int64(0)
	sup := min + binWidth
	var points []int64
	if h.bar == nil {
		// brand new histo
	for _, v := range h.unsorted {
		if v > sup {
			average /= np
			h.bar = append(h.bar, &Bar{average, np, sup, sup + binWidth, points})
			sup += binWidth
			average = v
			np = 1
		} else {
			points = append(make([]int64, 0, 1), v)
			average += v
			np++
		}
	}
	} else {
		if rebin {

		} else {
			// histo already exists and does not need to be extended for new points
			i := 0
			for _, v := range h.unsorted {
				for i < len(h.bar) {
					bar := h.bar[i]
					if v <= bar.Max {
						bar.Value = (bar.Value * bar.Count + v) / (bar.Count + 1)
						bar.Count++
						bar.Points = append(bar.Points, v)
						break
					}
					i++
				}
			}
		}
	}
	h.unsorted = nil	
}

func (h *Histo) bin(num int) {
	if h.unsorted == nil && num == h.num {
		// no new points and no change in nb of bars; nothing to do.
		return
	}

	max := int64(0)
	min := int64(0)
	if h.unsorted != nil {
		h.sort()
		max = h.unsorted[len(h.unsorted) - 1]
		min = h.unsorted[0]
	}
	if h.bar != nil {
		max = h.bar[len(h.bar) - 1].Max
		min = h.bar[0].Min
		if h.unsorted != nil {
			if h.unsorted[len(h.unsorted) - 1] > max {
				max = h.unsorted[len(h.unsorted) - 1]
			}
			if h.unsorted[0] < min {
				min = h.unsorted[0]
			}
		}
	} 
	binWidth := 1 + (max - min)/int64(num)
	np := int64(h.np)
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
	for _,v := range h.bar {
		println(v.Value, " ", v.Count)
	}
}

func (h *Histo) ReBin(num int) {
	h.bin(num)
}

func (h *Histo) Bars() [](*Bar) {
	h.bin(h.num)
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