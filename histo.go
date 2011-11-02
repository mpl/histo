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
	"image"
	"image/draw"
	"os"
	"sort"
)

type Bar struct {
	Value int64
	Count int64
	Min   int64
	Max   int64
}

type Histo [](*Bar)

func (h Histo) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Histo) Len() int           { return len(h) }
func (h Histo) Less(i, j int) bool { return h[i].Value < h[j].Value }

func NewHisto(m map[int64]int64, nb int) *Histo {
	h := make(Histo, len(m))
	i := 0
	for k, v := range m {
		h[i] = &Bar{Value: k, Count: v}
		i++
	}
	sort.Sort(h)
	return h.ReBin(nb)
}

func (h Histo) At(value int64) (hb *Bar) {
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

// TODO(mpl): variable width bins
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

type Params struct {
	fg draw.Image
	bg draw.Image
	op draw.Op
}

func NewParams(fg draw.Image, bg draw.Image, op draw.Op) (p *Params) {
	return &Params{fg: fg, bg: bg, op: op}
}

func (histo Histo) Draw(params *Params) (err os.Error) {
	if params.bg == nil {
		return os.NewError("Need a background image to draw on")
	}
	bg := params.bg
	w := bg.Bounds().Dx()
	h := bg.Bounds().Dy()
	var fg draw.Image
	if params.fg == nil {
		// defaults to plain green bars
		green := image.RGBAColor{0, 255, 0, 255}
		fg = image.NewNRGBA(w, h)
		draw.Draw(fg, fg.Bounds(), &image.ColorImage{green}, image.ZP, draw.Src)
	} else {
		fg = params.fg
	}

	barWidth := w / (2*len(histo) + 1)
	spacing := barWidth
	barHeight := 0
	pos := spacing
	highestCount := histo.MaxCount()

	for _, hb := range histo {
		barHeight = int(hb.Count * int64(h) / highestCount)
		r := image.Rect(pos, h-barHeight, pos+barWidth, h)
		draw.Draw(bg, r, fg, r.Min, params.op)
		pos += barWidth + spacing
	}
	return err
}
