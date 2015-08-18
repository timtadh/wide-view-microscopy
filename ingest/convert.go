package ingest

import (
	"image"
	"image/jpeg"
	_ "image/png"
	_ "image/gif"
	"io"
	"os"
	"path/filepath"
	"strings"
)

import (
	_ "golang.org/x/image/tiff"
)


func EncodeJpeg(img image.Image, to io.Writer) error {
	return jpeg.Encode(to, img, &jpeg.Options{Quality: 75})
}

func Decode(r io.Reader) (image.Image, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func ConvertToJpeg(from io.Reader, to io.Writer) error {
	img, err := Decode(from)
	if err != nil {
		return err
	}
	return EncodeJpeg(img, to)
}

func Jpeg(path string) (jpegPath string, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}
	dir := filepath.Dir(path)
	name := filepath.Base(path)
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)
	if ext == ".jpeg" || ext == ".jpg" {
		return path, nil
	}
	jpegPath = filepath.Join(dir, name + ".jpeg")
	fi, err := os.Stat(jpegPath)
	if err != nil && os.IsNotExist(err) {
		// its ok the path isn't there
	} else if err != nil {
		return "", err
	} else if fi.Size() > 0 {
		return jpegPath, nil
	}

	from, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer from.Close()

	to, err := os.Create(jpegPath)
	if err != nil {
		return "", err
	}
	defer to.Close()

	err = ConvertToJpeg(from, to)
	if err != nil {
		os.Remove(jpegPath)
		return "", err
	}
	return jpegPath, nil
}

