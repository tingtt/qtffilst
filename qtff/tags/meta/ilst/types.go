package ilst

import (
	"bytes"
	"qtffilst/internal/binary"
)

// https://exiftool.org/TagNames/QuickTime.html#ItemList
// https://developer.apple.com/documentation/quicktime-file-format/user_data_atoms#Media-characteristic-tags
type ItemList struct {
	// iTunesInfo	 `id:"----"` // QuickTime iTunesInfo Tags //? Unsupported
	ParentShortTitle      *string                `id:"@PST"`
	ParentProductID       *string                `id:"@ppi"`
	ParentTitle           *string                `id:"@pti"`
	ShortTitle            *string                `id:"@sti"`
	UnknownAACR           *string                `id:"AACR"`
	UnknownCDEK           *string                `id:"CDEK"`
	UnknownCDET           *string                `id:"CDET"`
	GUID                  *string                `id:"GUID"`
	ProductVersion        *string                `id:"VERS"`
	AlbumArtist           *InternationalText     `id:"aART"`
	AppleStoreAccountType *AppleStoreAccountType `id:"akID"`
	Album                 *string                `id:"albm"`
	AppleStoreAccount     *string                `id:"apID"`
	ArtistID              *int32                 `id:"atID"`
	Author                *string                `id:"auth"`
	Category              *string                `id:"catg"`
	ComposerID            *string                `id:"cmID"`
	AppleStoreCatalogID   *int32                 `id:"cnID"`
	CoverArt              *string                `id:"covr"`
	Compilation           *bool                  `id:"cpil"`
	Copyright             *InternationalText     `id:"cprt"`
	Description           *InternationalText     `id:"desc"`
	DiskNumber            *int32                 `id:"disk"`
	// Description           string `id:"dscp"` //? Unsupported
	EpisodeGlobalUniqueID *string `id:"egid"`
	GenreID               *int32  `id:"geID"` // QuickTime GenreID Values
	Genre                 *int32  `id:"gnre"`
	Grouping              *string `id:"grup"`
	GoogleHostHeader      *string `id:"gshh"`
	GooglePingMessage     *string `id:"gspm"`
	GooglePingURL         *string `id:"gspu"`
	GoogleSourceData      *string `id:"gssd"`
	GoogleStartTime       *string `id:"gsst"`
	GoogleTrackDuration   *string `id:"gstd"`
	HDVideo               *bool   `id:"hdvd"`
	ITunesU               *bool   `id:"itnu"`
	Keyword               *string `id:"keyw"`
	LongDescription       *string `id:"ldes"`
	Owner                 *string `id:"ownr"`
	Podcast               *bool   `id:"pcst"`
	Performer             *string `id:"perf"`
	DisableInsertPlayGap  *bool   `id:"pgap"`
	// AlbumID	int32[2]	  `id:"plID"` //? Unsupported, because I can’t understand document.
	ProductID    *string `id:"prID"`
	PurchaseDate *string `id:"purd"`
	PodcastURL   *string `id:"purl"`
	// RatingPercent     *string    `id:"rate"`  //? Unsupported
	ReleaseDate       *string            `id:"rldt"`
	Rating            *Rating            `id:"rate"`
	StoreDescription  *string            `id:"sdes"`
	AppleStoreCountry *int32             `id:"sfID"` // QuickTime AppleStoreCountry Values
	ShowMovement      *bool              `id:"shwm"`
	PreviewImage      *string            `id:"snal"`
	SortAlbumArtist   *InternationalText `id:"soaa"`
	SortAlbum         *InternationalText `id:"soal"`
	SortArtist        *InternationalText `id:"soar"`
	SortComposer      *InternationalText `id:"soco"`
	SortName          *InternationalText `id:"sonm"`
	SortShow          *InternationalText `id:"sosn"`
	MediaType         *MediaType         `id:"stik"`
	Title             *string            `id:"titl"`
	BeatsPerMinute    *int16             `id:"tmpo"`
	ThumbnailImage    *string            `id:"tnal"`
	TrackNumber       *int32             `id:"trkn"`
	TVEpisodeID       *string            `id:"tven"`
	TVEpisode         *int32             `id:"tves"`
	TVNetworkName     *string            `id:"tvnn"`
	TVShow            *string            `id:"tvsh"`
	TVSeason          *int32             `id:"tvsn"`
	ISRC              *string            `id:"xid "`
	Year              *string            `id:"yrrc"`
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
	MovementCount     *int16             `id:"(c)mvc"`
	MovementNumber    *int16             `id:"(c)mvi"`
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
	Text         string
	LanguageCode int32
}

func decodeInternationalText(data []byte) (InternationalText, error) {
	langCode, err := binary.BigEdian.ReadI32(bytes.NewBuffer(data[:8]))
	if err != nil {
		return InternationalText{}, err
	}

	return InternationalText{
		Text:         string(data[8:]),
		LanguageCode: langCode,
	}, nil
}

func (it InternationalText) Bytes() []byte {
	langCodeBuf := binary.BigEdian.BytesI32(it.LanguageCode)
	textBuf := ([]byte)(it.Text)
	return append(langCodeBuf, textBuf...)
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
