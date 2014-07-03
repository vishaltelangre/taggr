package taggr

import (
	"errors"
	"github.com/vishaltelangre/taggr/v1"
	"os"
)

type Taggr interface {
	Title() string
	Artist() string
	Album() string
	Year() string
	Comments() []string
	TrackNumber() uint8
	Genre() string
	SetTitle(string)
	SetArtist(string)
	SetAlbum(string)
	SetYear(string)
	SetGenre(string)
	SetTrackNumber(uint8)
	Size() int
	Buffer() []byte
}

type File struct {
	fileRef *os.File
	Taggr
}

func Open(name string) (*File, error) {
	f, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	file := &File{fileRef: f}
	v1Tag := v1.Parse(f)
	if v1Tag != nil {
		file.Taggr = v1Tag
	} else {
		return nil, errors.New("Unknown file type")
	}

	return file, nil
}

func (f *File) Close() error {
	defer f.fileRef.Close()
	switch f.Taggr.(type) {
	case (*v1.Tag):
		if _, err := f.fileRef.Seek(-v1.MaxTagSize, os.SEEK_END); err != nil {
			return err
		}
	default:
		return errors.New("Unknown tag version")
	}

	if _, err := f.fileRef.Write(f.Taggr.Buffer()); err != nil {
		return err
	}

	return nil
}
