package ilst

import (
	"bytes"
	"errors"

	"github.com/tingtt/qtffilst/internal/binary"
)

var (
	ErrInvalidLength = errors.New("invalid length")
)

// https://exiftool.org/TagNames/QuickTime.html#ItemList
// https://developer.apple.com/documentation/quicktime-file-format/user_data_atoms#Media-characteristic-tags
// Commented out fields are not supported
type ItemList struct {
	// iTunesInfo	 `id:"----"` // QuickTime iTunesInfo Tags //? Unsupported
	// ParentShortTitle      *string                `id:"@PST"`
	// ParentProductID       *string                `id:"@ppi"`
	// ParentTitle           *string                `id:"@pti"`
	// ShortTitle            *string                `id:"@sti"`
	// UnknownAACR           *string                `id:"AACR"`
	// UnknownCDEK           *string                `id:"CDEK"`
	// UnknownCDET           *string                `id:"CDET"`
	// GUID                  *string                `id:"GUID"`
	// ProductVersion        *string                `id:"VERS"`
	AlbumArtist *internationalText `id:"aART"`
	// AppleStoreAccountType *AppleStoreAccountType `id:"akID"`
	// Album                 *string                `id:"albm"`
	// AppleStoreAccount     *string                `id:"apID"`
	ArtistID *Int16WithHeader0x15_0 `id:"atID"`
	// Author                *string                `id:"auth"`
	// Category              *string                `id:"catg"`
	// ComposerID            *string                `id:"cmID"`
	// AppleStoreCatalogID   *int32                 `id:"cnID"`
	// CoverArt              *string                `id:"covr"`
	Compilation *BoolWithHeader0x15_0 `id:"cpil"`
	Copyright   *internationalText    `id:"cprt"`
	Description *internationalText    `id:"desc"`
	DiskNumber  *DiskNumber           `id:"disk"`
	// Description           string `id:"dscp"` //? Unsupported
	// EpisodeGlobalUniqueID *string               `id:"egid"`
	// GenreID               *int32                `id:"geID"` // QuickTime GenreID Values
	Genre *Genre `id:"gnre"`
	// Grouping              *string               `id:"grup"`
	// GoogleHostHeader      *string               `id:"gshh"`
	// GooglePingMessage     *string               `id:"gspm"`
	// GooglePingURL         *string               `id:"gspu"`
	// GoogleSourceData      *string               `id:"gssd"`
	// GoogleStartTime       *string               `id:"gsst"`
	// GoogleTrackDuration   *string               `id:"gstd"`
	// HDVideo               *bool                 `id:"hdvd"`
	// ITunesU               *bool                 `id:"itnu"`
	// Keyword               *string               `id:"keyw"`
	// LongDescription       *string               `id:"ldes"`
	// Owner                 *string               `id:"ownr"`
	// Podcast               *bool                 `id:"pcst"`
	// Performer             *string               `id:"perf"`
	DisableInsertPlayGap *BoolWithHeader0x15_0 `id:"pgap"`
	// AlbumID	int32[2]	  `id:"plID"` //? Unsupported, because I can’t understand document.
	// ProductID    *string `id:"prID"`
	// PurchaseDate *string `id:"purd"`
	// PodcastURL   *string `id:"purl"`
	// RatingPercent     *string    `id:"rate"`  //? Unsupported
	ReleaseDate *internationalText `id:"rldt"`
	// Rating            *Rating                `id:"rate"`
	// StoreDescription  *string                `id:"sdes"`
	// AppleStoreCountry *int32                 `id:"sfID"` // QuickTime AppleStoreCountry Values
	// ShowMovement      *bool                  `id:"shwm"`
	// PreviewImage      *string                `id:"snal"`
	SortAlbumArtist *internationalText `id:"soaa"`
	SortAlbum       *internationalText `id:"soal"`
	SortArtist      *internationalText `id:"soar"`
	SortComposer    *internationalText `id:"soco"`
	SortName        *internationalText `id:"sonm"`
	SortShow        *internationalText `id:"sosn"`
	// MediaType         *MediaType             `id:"stik"`
	// Title             *string                `id:"titl"`
	BeatsPerMinute *Int16WithHeader0x15_0 `id:"tmpo"`
	// ThumbnailImage    *string                `id:"tnal"`
	TrackNumber *TrackNumber `id:"trkn"`
	// TVEpisodeID       *string            `id:"tven"`
	// TVEpisode         *int32             `id:"tves"`
	// TVNetworkName     *string            `id:"tvnn"`
	// TVShow            *string            `id:"tvsh"`
	// TVSeason          *int32             `id:"tvsn"`
	// ISRC              *string            `id:"xid "`
	// Year              *string            `id:"yrrc"`
	Artist            *internationalText `id:"(c)ART"`
	AlbumC            *internationalText `id:"(c)alb"`
	ArtDirector       *internationalText `id:"(c)ard"`
	Arranger          *internationalText `id:"(c)arg"`
	AuthorC           *internationalText `id:"(c)aut"`
	Comment           *internationalText `id:"(c)cmt"`
	ComposerC         *internationalText `id:"(c)com"`
	Conductor         *internationalText `id:"(c)con"`
	CopyrightC        *internationalText `id:"(c)cpy"`
	ContentCreateDate *internationalText `id:"(c)day"`
	DescriptionC      *internationalText `id:"(c)des"`
	Director          *internationalText `id:"(c)dir"`
	EncodedBy         *internationalText `id:"(c)enc"`
	GenreC            *internationalText `id:"(c)gen"`
	GroupingC         *internationalText `id:"(c)grp"`
	Lyrics            *internationalText `id:"(c)lyr"`
	// MovementCount     *int16             `id:"(c)mvc"`
	// MovementNumber    *int16             `id:"(c)mvi"`
	MovementName      *internationalText `id:"(c)mvn"`
	TitleC            *internationalText `id:"(c)nam"`
	Narrator          *internationalText `id:"(c)nrt"`
	OriginalArtist    *internationalText `id:"(c)ope"`
	Producer          *internationalText `id:"(c)prd"`
	Publisher         *internationalText `id:"(c)pub"`
	SoundEngineer     *internationalText `id:"(c)sne"`
	Soloist           *internationalText `id:"(c)sol"`
	Subtitle          *internationalText `id:"(c)st3"`
	Encoder           *internationalText `id:"(c)too"`
	Track             *internationalText `id:"(c)trk"`
	Work              *internationalText `id:"(c)wrk"`
	ComposerCWRT      *internationalText `id:"(c)wrt"`
	ExecutiveProducer *internationalText `id:"(c)xpd"`
	GPSCoordinates    *internationalText `id:"(c)xyz"`
}

// https://developer.apple.com/documentation/quicktime-file-format/user_data_atoms#User-data-text-strings-and-language-codes
type internationalText struct {
	size         int32
	Text         string
	LanguageCode int32
}

func NewInternationalText(text string) *internationalText {
	return &internationalText{
		size:         1,
		Text:         text,
		LanguageCode: 0,
	}
}

func decodeInternationalText(data []byte) (internationalText, error) {
	if len(data) < 8 {
		return internationalText{}, ErrInvalidLength
	}
	size, err := binary.BigEdian.ReadI32(bytes.NewBuffer(data[:4]))
	if err != nil {
		return internationalText{}, err
	}
	langCode, err := binary.BigEdian.ReadI32(bytes.NewBuffer(data[4:8]))
	if err != nil {
		return internationalText{}, err
	}

	return internationalText{
		size:         size,
		Text:         string(data[8:]),
		LanguageCode: langCode,
	}, nil
}

func (it internationalText) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	if it.size == 0 {
		it.size = 1
	}
	_, err := buf.Write(binary.BigEdian.BytesI32(it.size))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(binary.BigEdian.BytesI32(it.LanguageCode))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(([]byte)(it.Text))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type (
	AppleStoreAccountType = int8
	Rating                = int8
	MediaType             = int8
)

const (
	AppleStoreAccountTypeITunes AppleStoreAccountType = iota
	AppleStoreAccountTypeAOL
)

const (
	RatingNone Rating = iota
	RatingExplicit
	RatingClean
	RatingExplicitOld = 4
)

const (
	MediaTypeMovieOld MediaType = iota
	MediaTypeNormalMusic
	MediaTypeAudiobook
	MediaTypeWhackedBookmark = 5
	MediaTypeMusicVideo      = 6
	MediaTypeMovie           = 9
	MediaTypeTVShow          = 10
	MediaTypeBooklet         = 11
	MediaTypeRingtone        = 14
	MediaTypePodcast         = 21
	MediaTypeiTunesU         = 23
)

type Genre int8

func decodeGenre(data []byte) (Genre, error) {
	if len(data) < 11 {
		return 0, ErrInvalidLength
	}
	buf := &bytes.Buffer{}
	_, err := buf.Write([]byte{0x0})
	if err != nil {
		return 0, err
	}
	_, err = buf.Write(data[10:])
	if err != nil {
		return 0, err
	}
	value, err := binary.BigEdian.ReadI8(buf)
	return Genre(value), err
}

func (g Genre) Bytes() ([]byte, error) {
	HEADER := bytes.Repeat([]byte{0x0}, 9)
	valueBuf := &bytes.Buffer{}
	_, err := valueBuf.Write(binary.BigEdian.BytesI8(int8(g)))
	return append(HEADER, valueBuf.Bytes()...), err
}

type TrackNumber struct {
	Number int16
	Total  int16
}

func decodeTrackNumber(data []byte) (TrackNumber, error) {
	if len(data) < 14 {
		return TrackNumber{}, ErrInvalidLength
	}
	buf := &bytes.Buffer{}
	_, err := buf.Write(data[10:12])
	if err != nil {
		return TrackNumber{}, err
	}
	number, err := binary.BigEdian.ReadI16(buf)
	if err != nil {
		return TrackNumber{}, err
	}
	buf2 := &bytes.Buffer{}
	_, err = buf2.Write(data[12:14])
	if err != nil {
		return TrackNumber{}, err
	}
	total, err := binary.BigEdian.ReadI16(buf2)
	if err != nil {
		return TrackNumber{}, err
	}
	return TrackNumber{
		Number: number,
		Total:  total,
	}, nil
}

func (tn TrackNumber) Bytes() ([]byte, error) {
	HEADER := bytes.Repeat([]byte{0x0}, 10)
	valueBuf := &bytes.Buffer{}
	_, err := valueBuf.Write(binary.BigEdian.BytesI16(tn.Number))
	if err != nil {
		return nil, err
	}
	_, err = valueBuf.Write(binary.BigEdian.BytesI16(tn.Total))
	if err != nil {
		return nil, err
	}
	return append(HEADER, valueBuf.Bytes()...), nil
}

type DiskNumber struct {
	Number int16
	Total  int16
}

func decodeDiskNumber(data []byte) (DiskNumber, error) {
	if len(data) < 14 {
		return DiskNumber{}, ErrInvalidLength
	}
	buf := &bytes.Buffer{}
	_, err := buf.Write(data[10:12])
	if err != nil {
		return DiskNumber{}, err
	}
	number, err := binary.BigEdian.ReadI16(buf)
	if err != nil {
		return DiskNumber{}, err
	}
	buf2 := &bytes.Buffer{}
	_, err = buf2.Write(data[12:14])
	if err != nil {
		return DiskNumber{}, err
	}
	total, err := binary.BigEdian.ReadI16(buf2)
	if err != nil {
		return DiskNumber{}, err
	}
	return DiskNumber{
		Number: number,
		Total:  total,
	}, nil
}

func (tn DiskNumber) Bytes() ([]byte, error) {
	HEADER := bytes.Repeat([]byte{0x0}, 10)
	FOOTER := []byte{0x0, 0x0}
	valueBuf := &bytes.Buffer{}
	_, err := valueBuf.Write(binary.BigEdian.BytesI16(tn.Number))
	if err != nil {
		return nil, err
	}
	_, err = valueBuf.Write(binary.BigEdian.BytesI16(tn.Total))
	if err != nil {
		return nil, err
	}
	_, err = valueBuf.Write(FOOTER)
	if err != nil {
		return nil, err
	}
	return append(HEADER, valueBuf.Bytes()...), nil
}

type Int16WithHeader0x15_0 struct {
	Value int16
}

func decodeInt16WithHeader0x15_0(data []byte) (Int16WithHeader0x15_0, error) {
	if len(data) < 10 {
		return Int16WithHeader0x15_0{}, ErrInvalidLength
	}
	value, err := binary.BigEdian.ReadI16(bytes.NewBuffer(data[8:]))
	if err != nil {
		return Int16WithHeader0x15_0{}, err
	}

	return Int16WithHeader0x15_0{value}, nil
}

func (i Int16WithHeader0x15_0) Bytes() []byte {
	HEADER := []byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x0, 0x0, 0x0}
	valueBuf := binary.BigEdian.BytesI16(i.Value)
	return append(HEADER, valueBuf...)
}

type BoolWithHeader0x15_0 struct {
	Value bool
}

func decodeBoolWithHeader0x15_0(data []byte) (BoolWithHeader0x15_0, error) {
	if len(data) < 9 {
		return BoolWithHeader0x15_0{}, ErrInvalidLength
	}
	value, err := binary.BigEdian.ReadI8(bytes.NewBuffer(append([]byte{0x0}, data[8:]...)))
	if err != nil {
		return BoolWithHeader0x15_0{}, err
	}
	return BoolWithHeader0x15_0{value != 0}, nil
}

func (i BoolWithHeader0x15_0) Bytes() []byte {
	HEADER := []byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x0, 0x0, 0x0}
	intBool := int8(0)
	if i.Value {
		intBool = 1
	}
	return append(HEADER, binary.BigEdian.BytesI8(intBool)...)
}
