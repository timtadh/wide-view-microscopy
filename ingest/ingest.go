package ingest

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)


type Image struct {
	Path string
	Metadata map[string]string
}

func Ingest(dir string, format Format) (paths []*Image, err error) {
	files, err := Files(dir)
	if err != nil {
		return nil, err
	}
	var path string
	for path, err, files = files(); files != nil; path, err, files = files() {
		name := filepath.Base(path)
		meta, err := format.Parse([]byte(name))
		if err != nil {
			log.Println("WARN", "skipping", path, "because", err)
		} else {
			var use string
			jpeg, err := Jpeg(path)
			if err != nil {
				use = path
				log.Println("WARN", "could not convert to jpeg", path, "using tiff. because", err)
			} else {
				use = jpeg
			}
			paths = append(paths, &Image{use, meta})
		}
	}
	if err != nil {
		return nil, err
	}
	return paths, nil
}

type StringIterator func()(string, error, StringIterator)

func Files(dir string) (si StringIterator, err error) {
	type entry struct {
		path string
		fi os.FileInfo
	}
	pop := func(stack []entry) (entry, []entry) {
		return stack[len(stack)-1], stack[:len(stack)-1]
	}
	push := func(stack []entry, dir string) ([]entry, error) {
		fis, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for i := len(fis) - 1; i >= 0; i-- {
			fi := fis[i]
			stack = append(stack, entry{
				path: filepath.Join(dir, fi.Name()),
				fi: fi,
			})
		}
		return stack, nil
	}
	stack, err := push(make([]entry, 0, 10), dir)
	if err != nil {
		return nil, err
	}
	si = func() (path string, err error, _ StringIterator) {
		var e entry
		if len(stack) == 0 {
			return "", nil, nil
		}
		e, stack = pop(stack)
		for e.fi.IsDir() {
			stack, err = push(stack, e.path)
			if err != nil {
				return "", err, nil
			}
			e, stack = pop(stack)
		}
		return e.path, nil, si
	}
	return si, nil
}

