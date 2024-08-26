package ilst

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
	AlbumArtist           *string                `id:"aART"`
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
	Copyright             *string                `id:"cprt"`
	Description           *string                `id:"desc"`
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
	ReleaseDate       *string    `id:"rldt"`
	Rating            *Rating    `id:"rate"`
	StoreDescription  *string    `id:"sdes"`
	AppleStoreCountry *int32     `id:"sfID"` // QuickTime AppleStoreCountry Values
	ShowMovement      *bool      `id:"shwm"`
	PreviewImage      *string    `id:"snal"`
	SortAlbumArtist   *string    `id:"soaa"`
	SortAlbum         *string    `id:"soal"`
	SortArtist        *string    `id:"soar"`
	SortComposer      *string    `id:"soco"`
	SortName          *string    `id:"sonm"`
	SortShow          *string    `id:"sosn"`
	MediaType         *MediaType `id:"stik"`
	Title             *string    `id:"titl"`
	BeatsPerMinute    *int16     `id:"tmpo"`
	ThumbnailImage    *string    `id:"tnal"`
	TrackNumber       *int32     `id:"trkn"`
	TVEpisodeID       *string    `id:"tven"`
	TVEpisode         *int32     `id:"tves"`
	TVNetworkName     *string    `id:"tvnn"`
	TVShow            *string    `id:"tvsh"`
	TVSeason          *int32     `id:"tvsn"`
	ISRC              *string    `id:"xid "`
	Year              *string    `id:"yrrc"`
	Artist            *string    `id:"(c)ART"`
	AlbumC            *string    `id:"(c)alb"`
	ArtDirector       *string    `id:"(c)ard"`
	Arranger          *string    `id:"(c)arg"`
	AuthorC           *string    `id:"(c)aut"`
	Comment           *string    `id:"(c)cmt"`
	ComposerC         *string    `id:"(c)com"`
	Conductor         *string    `id:"(c)con"`
	CopyrightC        *string    `id:"(c)cpy"`
	ContentCreateDate *string    `id:"(c)day"`
	DescriptionC      *string    `id:"(c)des"`
	Director          *string    `id:"(c)dir"`
	EncodedBy         *string    `id:"(c)enc"`
	GenreC            *string    `id:"(c)gen"`
	GroupingC         *string    `id:"(c)grp"`
	Lyrics            *string    `id:"(c)lyr"`
	MovementCount     *int16     `id:"(c)mvc"`
	MovementNumber    *int16     `id:"(c)mvi"`
	MovementName      *string    `id:"(c)mvn"`
	TitleC            *string    `id:"(c)nam"`
	Narrator          *string    `id:"(c)nrt"`
	OriginalArtist    *string    `id:"(c)ope"`
	Producer          *string    `id:"(c)prd"`
	Publisher         *string    `id:"(c)pub"`
	SoundEngineer     *string    `id:"(c)sne"`
	Soloist           *string    `id:"(c)sol"`
	Subtitle          *string    `id:"(c)st3"`
	Encoder           *string    `id:"(c)too"`
	Track             *string    `id:"(c)trk"`
	Work              *string    `id:"(c)wrk"`
	ComposerCWRT      *string    `id:"(c)wrt"`
	ExecutiveProducer *string    `id:"(c)xpd"`
	GPSCoordinates    *string    `id:"(c)xyz"`
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
