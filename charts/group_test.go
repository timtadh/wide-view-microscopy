package charts

import "testing"

import (
	"github.com/timtadh/wide-view-microscopy/ingest"
)

func TestGroup(t *testing.T) {
	eatError := func(m ingest.Metadata, err error) ingest.Metadata {
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	format, err := ingest.ParseFormatString("$(slide) $(sample) $(region) $(stain).tif")
	if err != nil {
		t.Fatal(err)
	}
	images := []*ingest.Image{
		{"path/a", eatError(format.Parse([]byte("slide-1 sample-1 L1 FFa.tif")))},
		{"path/b", eatError(format.Parse([]byte("slide-2 sample-1 L2 FFa.tif")))},
		{"path/c", eatError(format.Parse([]byte("slide-1 sample-1 L3 FFa.tif")))},
		{"path/d", eatError(format.Parse([]byte("slide-2 sample-1 L1 FFb.tif")))},
		{"path/e", eatError(format.Parse([]byte("slide-1 sample-1 L2 FFb.tif")))},
		{"path/f", eatError(format.Parse([]byte("slide-2 sample-1 L3 FFb.tif")))},
		{"path/g", eatError(format.Parse([]byte("slide-1 sample-1 L1 FFc.tif")))},
		{"path/h", eatError(format.Parse([]byte("slide-2 sample-1 L2 FFc.tif")))},
		{"path/i", eatError(format.Parse([]byte("slide-1 sample-1 L3 FFc.tif")))},
	}
	rows := MakeRows(images, []string{"slide", "region"}, []string{})
	for _, row := range rows {
		t.Log("row", row.Meta())
		for _, img := range row.Images() {
			t.Log(img)
		}
		t.Log()
	}
	paths := []string{
		"path/a",
		"path/g",
		"path/e",
		"path/c",
		"path/i",
		"path/d",
		"path/b",
		"path/h",
		"path/f",
	}
	c := &Chart{rows: rows}
	for i, img := range c.Images() {
		if img.Path != paths[i] {
			t.Fatal("img != paths[i]", img, paths[i])
		}
	}
}


func TestMakeCharts(t *testing.T) {
	eatError := func(m ingest.Metadata, err error) ingest.Metadata {
		if err != nil {
			t.Fatal(err)
		}
		return m
	}
	format, err := ingest.ParseFormatString("$(slide) $(sample) $(region) $(stain).tif")
	if err != nil {
		t.Fatal(err)
	}
	images := []*ingest.Image{
		{"path/a-1", eatError(format.Parse([]byte("slide-1 sample-1 L1 FFa.tif")))},
		{"path/b-1", eatError(format.Parse([]byte("slide-1 sample-1 L2 FFa.tif")))},
		{"path/c-1", eatError(format.Parse([]byte("slide-1 sample-1 L3 FFa.tif")))},
		{"path/d-1", eatError(format.Parse([]byte("slide-1 sample-1 L1 FFb.tif")))},
		{"path/e-1", eatError(format.Parse([]byte("slide-1 sample-1 L2 FFb.tif")))},
		{"path/f-1", eatError(format.Parse([]byte("slide-1 sample-1 L3 FFb.tif")))},
		{"path/g-1", eatError(format.Parse([]byte("slide-1 sample-1 L1 FFc.tif")))},
		{"path/h-1", eatError(format.Parse([]byte("slide-1 sample-1 L2 FFc.tif")))},
		{"path/i-1", eatError(format.Parse([]byte("slide-1 sample-1 L3 FFc.tif")))},
		{"path/a-2", eatError(format.Parse([]byte("slide-2 sample-1 L1 FFa.tif")))},
		{"path/b-2", eatError(format.Parse([]byte("slide-2 sample-1 L2 FFa.tif")))},
		{"path/c-2", eatError(format.Parse([]byte("slide-2 sample-1 L3 FFa.tif")))},
		{"path/d-2", eatError(format.Parse([]byte("slide-2 sample-1 L1 FFb.tif")))},
		{"path/e-2", eatError(format.Parse([]byte("slide-2 sample-1 L2 FFb.tif")))},
		{"path/f-2", eatError(format.Parse([]byte("slide-2 sample-1 L3 FFb.tif")))},
		{"path/g-2", eatError(format.Parse([]byte("slide-2 sample-1 L1 FFc.tif")))},
		{"path/h-2", eatError(format.Parse([]byte("slide-2 sample-1 L2 FFc.tif")))},
		{"path/i-2", eatError(format.Parse([]byte("slide-2 sample-1 L3 FFc.tif")))},
		{"path/a-3", eatError(format.Parse([]byte("slide-1 sample-2 L1 FFa.tif")))},
		{"path/b-3", eatError(format.Parse([]byte("slide-1 sample-2 L2 FFa.tif")))},
		{"path/c-3", eatError(format.Parse([]byte("slide-1 sample-2 L3 FFa.tif")))},
		{"path/d-3", eatError(format.Parse([]byte("slide-1 sample-2 L1 FFb.tif")))},
		{"path/e-3", eatError(format.Parse([]byte("slide-1 sample-2 L2 FFb.tif")))},
		{"path/f-3", eatError(format.Parse([]byte("slide-1 sample-2 L3 FFb.tif")))},
		{"path/g-3", eatError(format.Parse([]byte("slide-1 sample-2 L1 FFc.tif")))},
		{"path/h-3", eatError(format.Parse([]byte("slide-1 sample-2 L2 FFc.tif")))},
		{"path/i-3", eatError(format.Parse([]byte("slide-1 sample-2 L3 FFc.tif")))},
		{"path/a-4", eatError(format.Parse([]byte("slide-2 sample-2 L1 FFa.tif")))},
		{"path/b-4", eatError(format.Parse([]byte("slide-2 sample-2 L2 FFa.tif")))},
		{"path/c-4", eatError(format.Parse([]byte("slide-2 sample-2 L3 FFa.tif")))},
		{"path/d-4", eatError(format.Parse([]byte("slide-2 sample-2 L1 FFb.tif")))},
		{"path/e-4", eatError(format.Parse([]byte("slide-2 sample-2 L2 FFb.tif")))},
		{"path/f-4", eatError(format.Parse([]byte("slide-2 sample-2 L3 FFb.tif")))},
		{"path/g-4", eatError(format.Parse([]byte("slide-2 sample-2 L1 FFc.tif")))},
		{"path/h-4", eatError(format.Parse([]byte("slide-2 sample-2 L2 FFc.tif")))},
		{"path/i-4", eatError(format.Parse([]byte("slide-2 sample-2 L3 FFc.tif")))},
	}
	charts := MakeCharts(images, []string{"sample", "slide"}, []string{"region"}, []string{"stain"})
	for _, chart := range charts {
		t.Log("chart", chart.Meta())
		for _, row := range chart.rows {
			t.Log("row", row.Meta())
			for _, img := range row.Images() {
				t.Log(img)
			}
			t.Log()
		}
		t.Log()
	}
}
