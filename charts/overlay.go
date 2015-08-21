package charts

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/timtadh/wide-view-microscopy/ingest"
	"github.com/disintegration/imaging"
)


func Overlay(images []*ingest.Image) (*ingest.Image, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("empty slice was passed in")
	} else if len(images) == 1 {
		return images[0], nil
	}
	meta := CommonMeta(images)
	path := OverlayName(images)
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		// ok we will make it below
	} else if err != nil {
		return nil, err
	} else if fi.Size() > 0 {
		return &ingest.Image{path, meta}, nil
	}
	imgs := make([]image.Image, 0, len(images))
	for _, i := range images {
		img, err := ingest.LoadImage(i.Path)
		if err != nil {
			return nil, err
		}
		imgs = append(imgs, img)
	}
	overlay := imgs[0]
	for i := 1; i < len(imgs); i++ {
		overlay = imaging.Overlay(overlay, imgs[i], image.Pt(0,0), .50)
	}
	err = ingest.WriteJpeg(path, overlay)
	if err != nil {
		return nil, err
	}
	return &ingest.Image{path, meta}, nil
}

func OverlayName(images []*ingest.Image) string {
	dir := filepath.Dir(images[0].Path)
	names := make([]string, 0, len(images))
	for _, img := range images {
		name := filepath.Base(img.Path)
		ext := filepath.Ext(name)
		name = strings.TrimSuffix(name, ext)
		names = append(names, name)
	}
	name := "overlay::" + strings.Join(names, ":") + ".jpeg"
	return filepath.Join(dir, name)
}

func CommonMeta(images []*ingest.Image) ingest.Metadata {
	meta := make(ingest.Metadata, len(images[0].Meta()))
	kv: for k, v := range images[0].Meta() {
		for _, img := range images {
			if img.Meta()[k] != v {
				continue kv
			}
		}
		meta[k] = v
	}
	return meta
}

