package v1

// Reference: http://id3.org/ID3v1

import (
	"io"
	"os"
)

const (
	MaxTagSize = 128
)

type Tag struct {
	title, artist, album, year, comment string
	track                               uint8
	genre                               byte
}

// Eric Kemp's Genre List
var Genres = []string{
	"Blues", "Classic Rock", "Country", "Dance", "Disco", "Funk", "Grunge",
	"Hip-Hop", "Jazz", "Metal", "New Age", "Oldies", "Other", "Pop",
	"Rhythm and Blues", "Rap", "Reggae", "Rock", "Techno", "Industrial",
	"Alternative", "Ska", "Death Metal", "Pranks", "Soundtrack", "Euro-Techno",
	"Ambient", "Trip-Hop", "Vocal", "Jazz & Funk", "Fusion", "Trance",
	"Classical", "Instrumental", "Acid", "House", "Game", "Sound Clip",
	"Gospel", "Noise", "Alternative Rock", "Bass", "Soul", "Punk", "Space",
	"Meditative", "Instrumental Pop", "Instrumental Rock", "Ethnic", "Gothic",
	"Darkwave", "Techno-Industrial", "Electronic", "Pop-Folk", "Eurodance",
	"Dream", "Southern Rock", "Comedy", "Cult", "Gangsta", "Top 40",
	"Christian Rap", "Pop/Funk", "Jungle", "Native American", "Cabaret",
	"New Wave", "Psychedelic", "Rave", "Showtunes", "Trailer", "Lo-Fi",
	"Tribal", "Acid Punk", "Acid Jazz", "Polka", "Retro", "Musical",
	"Rock & Roll", "Hard Rock",
}

func Parse(f io.ReadSeeker) *Tag {
	// Read from end character to last 128 chars
	f.Seek(-MaxTagSize, os.SEEK_END)

	stuff := make([]byte, MaxTagSize)

	n, err := io.ReadFull(f, stuff)
	// Check if it contains valid ID3v1/1.1 data, which can be done by looking
	// at first 3 characters out of 128 characters, which should be "TAG"
	if n < MaxTagSize || err != nil || string(stuff[:3]) != "TAG" {
		return nil
	}

	var tag = new(Tag)
	tag.title = string(stuff[3:33])
	tag.artist = string(stuff[33:63])
	tag.album = string(stuff[63:93])
	tag.year = string(stuff[93:97])
	tag.comment = string(stuff[97:127])
	tag.genre = stuff[127]

	// If a track number is stored, then 125 byte contains a binary 0, and
	// comment is 28 characters long only.
	if stuff[125] == '0' {
		tag.comment = string(stuff[97:125])
		tag.track = uint8(stuff[126])
	}

	return tag
}

func (t Tag) Title() string      { return t.title }
func (t Tag) Artist() string     { return t.artist }
func (t Tag) Album() string      { return t.album }
func (t Tag) Year() string       { return t.year }
func (t Tag) TrackNumber() uint8 { return t.track }
func (t Tag) Size() int          { return MaxTagSize }

func (t Tag) Comments() []string {
	return []string{t.comment}
}

func (t Tag) Genre() string {
	if int(t.genre) < len(Genres) {
		return Genres[t.genre]
	}
	return ""
}

func (t *Tag) SetTitle(txt string)      { t.title = txt }
func (t *Tag) SetArtist(a string)       { t.artist = a }
func (t *Tag) SetAlbum(a string)        { t.album = a }
func (t *Tag) SetYear(y string)         { t.year = y }
func (t *Tag) SetTrackNumber(tno uint8) { t.track = tno }

func (t *Tag) SetGenre(g string) {
	// Index in a list of genres, or 255
	t.genre = 255
	for i, genre := range Genres {
		if g == genre {
			t.genre = uint8(i)
			break
		}
	}
}

func (t Tag) Buffer() []byte {
	stuff := make([]byte, MaxTagSize)

	copy(stuff[:3], []byte("TAG"))
	copy(stuff[3:33], []byte(t.title))
	copy(stuff[33:63], []byte(t.artist))
	copy(stuff[63:93], []byte(t.album))
	copy(stuff[93:97], []byte(t.year))
	copy(stuff[97:125], []byte(t.comment))
	stuff[125] = '0'
	stuff[126] = t.track
	stuff[127] = t.genre

	return stuff
}
