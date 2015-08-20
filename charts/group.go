package charts

import (
	"sort"
)

import (
	"github.com/timtadh/wide-view-microscopy/ingest"
)

type Chart struct {
	meta ingest.Metadata
	rows []*Row
}

type Row struct {
	meta ingest.Metadata
	images []*ingest.Image
}

type Images interface {
	Meta() ingest.Metadata
	Images() []*ingest.Image
}

type Grouped interface {
	Subgroups() []Images
}

func Submeta(meta ingest.Metadata, keys []string) ingest.Metadata {
	sub := make(ingest.Metadata, len(keys))
	for _, key := range keys {
		sub[key] = meta[key]
	}
	return sub
}

type Sortable interface {
	Sort()
}

type sortableImages struct {
	Images []Images
	On []string
}

func OrderBy(images []Images, on []string) []Images {
	list := make([]Images, len(images))
	copy(list, images)
	s := &sortableImages{
		Images: list,
		On: on,
	}
	sort.Sort(s)
	return list
}

func (s *sortableImages) Len() int {
	return len(s.Images)
}

func (s *sortableImages) Swap(i, j int) {
	s.Images[i], s.Images[j] = s.Images[j], s.Images[i]
}

func (s *sortableImages) Less(i, j int) bool {
	a := s.Images[i].Meta()
	b := s.Images[j].Meta()
	for i := 0; i < len(s.On) - 1; i++ {
		key := s.On[i]
		if a[key] < b[key] {
			return true
		} else if a[key] > b[key] {
			return false
		}
	}
	key := s.On[len(s.On)-1]
	return a[key] < b[key]
}

func Group(images []Images, on []string) ([][]Images, []ingest.Metadata) {
	if len(images) <= 0 {
		return nil, nil
	}
	groups := make([][]Images, 0, len(images))
	metas := make([]ingest.Metadata, 0, len(images))
	if len(on) <= 0 {
		for _, img := range images {
			groups = append(groups, []Images{img})
			metas = append(metas, make(ingest.Metadata))
		}
		return groups, metas
	}
	images = OrderBy(images, on)
	cur := Submeta(images[0].Meta(), on)
	group := make([]Images, 0, 10)
	for _, img := range images {
		sm := Submeta(img.Meta(), on)
		if !cur.Equal(sm) {
			groups = append(groups, group)
			metas = append(metas, cur)
			group = make([]Images, 0, 10)
			cur = sm
		}
		group = append(group, img)
	}
	groups = append(groups, group)
	metas = append(metas, cur)
	return groups, metas
}

func MakeRows(images []*ingest.Image, on []string) []*Row {
	groups, metas := Group(imageListAsImages(images), on)
	rows := make([]*Row, 0, len(groups))
	for i := 0; i < len(groups); i++ {
		rows = append(rows, &Row{meta: metas[i], images: imagesAsImageList(groups[i])})
	}
	return rows
}

func MakeCharts(rows []*Row, on []string) []*Chart {
	groups, metas := Group(rowsAsImages(rows), on)
	charts := make([]*Chart, 0, len(groups))
	for i := 0; i < len(groups); i++ {
		charts = append(charts, &Chart{meta: metas[i], rows: imagesAsRows(groups[i])})
	}
	return charts
}

func imagesAsImageList(images []Images) []*ingest.Image {
	list := make([]*ingest.Image, 0, len(images))
	for _, img := range images {
		list = append(list, img.(*ingest.Image))
	}
	return list
}

func imagesAsRows(images []Images) []*Row {
	list := make([]*Row, 0, len(images))
	for _, img := range images {
		list = append(list, img.(*Row))
	}
	return list
}

func imageListAsImages(images []*ingest.Image) []Images {
	list := make([]Images, 0, len(images))
	for _, img := range images {
		list = append(list, img)
	}
	return list
}

func rowsAsImages(rows []*Row) []Images {
	list := make([]Images, 0, len(rows))
	for _, row := range rows {
		list = append(list, row)
	}
	return list
}

func (c *Chart) Meta() ingest.Metadata {
	return c.meta
}

func (c *Chart) Images() []*ingest.Image {
	list := make([][]*ingest.Image, 0, len(c.rows))
	size := 0
	for _, row := range c.rows {
		imgs := row.Images()
		list = append(list, imgs)
		size += len(imgs)
	}
	images := make([]*ingest.Image, 0, size)
	for _, imgs := range list {
		images = append(images, imgs...)
	}
	return images
}

func (c *Chart) Subgroups() []Images {
	list := make([]Images, 0, len(c.rows))
	for _, row := range c.rows {
		list = append(list, row)
	}
	return list
}

func (r *Row) Meta() ingest.Metadata {
	return r.meta
}

func (r *Row) Images() []*ingest.Image {
	return r.images
}

func (r *Row) Subgroups() []Images {
	return imageListAsImages(r.images)
}

