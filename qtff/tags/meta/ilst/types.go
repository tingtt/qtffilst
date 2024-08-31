package ilst

import (
	"bytes"
	"errors"
	"qtffilst/internal/binary"
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
	AlbumArtist *InternationalText `id:"aART"`
	// AppleStoreAccountType *AppleStoreAccountType `id:"akID"`
	// Album                 *string                `id:"albm"`
	// AppleStoreAccount     *string                `id:"apID"`
	ArtistID *int16WithHeader0x15_0 `id:"atID"`
	// Author                *string                `id:"auth"`
	// Category              *string                `id:"catg"`
	// ComposerID            *string                `id:"cmID"`
	// AppleStoreCatalogID   *int32                 `id:"cnID"`
	// CoverArt              *string                `id:"covr"`
	Compilation *boolWithHeader0x15_0 `id:"cpil"`
	Copyright   *InternationalText    `id:"cprt"`
	Description *InternationalText    `id:"desc"`
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
	DisableInsertPlayGap *boolWithHeader0x15_0 `id:"pgap"`
	// AlbumID	int32[2]	  `id:"plID"` //? Unsupported, because I can’t understand document.
	// ProductID    *string `id:"prID"`
	// PurchaseDate *string `id:"purd"`
	// PodcastURL   *string `id:"purl"`
	// RatingPercent     *string    `id:"rate"`  //? Unsupported
	// ReleaseDate       *string                `id:"rldt"`
	// Rating            *Rating                `id:"rate"`
	// StoreDescription  *string                `id:"sdes"`
	// AppleStoreCountry *int32                 `id:"sfID"` // QuickTime AppleStoreCountry Values
	// ShowMovement      *bool                  `id:"shwm"`
	// PreviewImage      *string                `id:"snal"`
	SortAlbumArtist *InternationalText `id:"soaa"`
	SortAlbum       *InternationalText `id:"soal"`
	SortArtist      *InternationalText `id:"soar"`
	SortComposer    *InternationalText `id:"soco"`
	SortName        *InternationalText `id:"sonm"`
	SortShow        *InternationalText `id:"sosn"`
	// MediaType         *MediaType             `id:"stik"`
	// Title             *string                `id:"titl"`
	BeatsPerMinute *int16WithHeader0x15_0 `id:"tmpo"`
	// ThumbnailImage    *string                `id:"tnal"`
	TrackNumber *TrackNumber `id:"trkn"`
	// TVEpisodeID       *string            `id:"tven"`
	// TVEpisode         *int32             `id:"tves"`
	// TVNetworkName     *string            `id:"tvnn"`
	// TVShow            *string            `id:"tvsh"`
	// TVSeason          *int32             `id:"tvsn"`
	// ISRC              *string            `id:"xid "`
	// Year              *string            `id:"yrrc"`
	Artist            *InternationalText `id:"(c)ART"`
	AlbumC            *InternationalText `id:"(c)alb"`
	ArtDirector       *InternationalText `id:"(c)ard"`
	Arranger          *InternationalText `id:"(c)arg"`
	AuthorC           *InternationalText `id:"(c)aut"`
	Comment           *InternationalText `id:"(c)cmt"`
	ComposerC         *InternationalText `id:"(c)com"`
	Conductor         *InternationalText `id:"(c)con"`
	CopyrightC        *InternationalText `id:"(c)cpy"`
	ContentCreateDate *InternationalText `id:"(c)day"`
	DescriptionC      *InternationalText `id:"(c)des"`
	Director          *InternationalText `id:"(c)dir"`
	EncodedBy         *InternationalText `id:"(c)enc"`
	GenreC            *InternationalText `id:"(c)gen"`
	GroupingC         *InternationalText `id:"(c)grp"`
	Lyrics            *InternationalText `id:"(c)lyr"`
	// MovementCount     *int16             `id:"(c)mvc"`
	// MovementNumber    *int16             `id:"(c)mvi"`
	MovementName      *InternationalText `id:"(c)mvn"`
	TitleC            *InternationalText `id:"(c)nam"`
	Narrator          *InternationalText `id:"(c)nrt"`
	OriginalArtist    *InternationalText `id:"(c)ope"`
	Producer          *InternationalText `id:"(c)prd"`
	Publisher         *InternationalText `id:"(c)pub"`
	SoundEngineer     *InternationalText `id:"(c)sne"`
	Soloist           *InternationalText `id:"(c)sol"`
	Subtitle          *InternationalText `id:"(c)st3"`
	Encoder           *InternationalText `id:"(c)too"`
	Track             *InternationalText `id:"(c)trk"`
	Work              *InternationalText `id:"(c)wrk"`
	ComposerCWRT      *InternationalText `id:"(c)wrt"`
	ExecutiveProducer *InternationalText `id:"(c)xpd"`
	GPSCoordinates    *InternationalText `id:"(c)xyz"`
}

// https://developer.apple.com/documentation/quicktime-file-format/user_data_atoms#User-data-text-strings-and-language-codes
type InternationalText struct {
	size         int32
	Text         string
	LanguageCode int32
}

func decodeInternationalText(data []byte) (InternationalText, error) {
	if len(data) < 8 {
		return InternationalText{}, ErrInvalidLength
	}
	size, err := binary.BigEdian.ReadI32(bytes.NewBuffer(data[:4]))
	if err != nil {
		return InternationalText{}, err
	}
	langCode, err := binary.BigEdian.ReadI32(bytes.NewBuffer(data[4:8]))
	if err != nil {
		return InternationalText{}, err
	}

	return InternationalText{
		size:         size,
		Text:         string(data[8:]),
		LanguageCode: langCode,
	}, nil
}

func (it InternationalText) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
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

type int16WithHeader0x15_0 struct {
	Value int16
}

func decodeInt16WithHeader0x15_0(data []byte) (int16WithHeader0x15_0, error) {
	if len(data) < 10 {
		return int16WithHeader0x15_0{}, ErrInvalidLength
	}
	value, err := binary.BigEdian.ReadI16(bytes.NewBuffer(data[8:]))
	if err != nil {
		return int16WithHeader0x15_0{}, err
	}

	return int16WithHeader0x15_0{value}, nil
}

func (i int16WithHeader0x15_0) Bytes() []byte {
	HEADER := []byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x0, 0x0, 0x0}
	valueBuf := binary.BigEdian.BytesI16(i.Value)
	return append(HEADER, valueBuf...)
}

type boolWithHeader0x15_0 struct {
	Value bool
}

func decodeBoolWithHeader0x15_0(data []byte) (boolWithHeader0x15_0, error) {
	if len(data) < 9 {
		return boolWithHeader0x15_0{}, ErrInvalidLength
	}
	value, err := binary.BigEdian.ReadI8(bytes.NewBuffer(append([]byte{0x0}, data[8:]...)))
	if err != nil {
		return boolWithHeader0x15_0{}, err
	}
	return boolWithHeader0x15_0{value != 0}, nil
}

func (i boolWithHeader0x15_0) Bytes() []byte {
	HEADER := []byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x0, 0x0, 0x0}
	intBool := int8(0)
	if i.Value {
		intBool = 1
	}
	return append(HEADER, binary.BigEdian.BytesI8(intBool)...)
}
