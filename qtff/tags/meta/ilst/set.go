package ilst

import (
	"encoding/binary"
)

func (il *ItemList) Set(id string, value []byte) (err error) {
	var internationalText internationalText
	str := string(value)

	switch id {
	// international text
	case "aART":
		internationalText, err = decodeInternationalText(value)
		il.AlbumArtist = &internationalText
	case "cprt":
		internationalText, err = decodeInternationalText(value)
		il.Copyright = &internationalText
	case "desc":
		internationalText, err = decodeInternationalText(value)
		il.Description = &internationalText
	case "soaa":
		internationalText, err = decodeInternationalText(value)
		il.SortAlbumArtist = &internationalText
	case "soal":
		internationalText, err = decodeInternationalText(value)
		il.SortAlbum = &internationalText
	case "soar":
		internationalText, err = decodeInternationalText(value)
		il.SortArtist = &internationalText
	case "soco":
		internationalText, err = decodeInternationalText(value)
		il.SortComposer = &internationalText
	case "sonm":
		internationalText, err = decodeInternationalText(value)
		il.SortName = &internationalText
	case "sosn":
		internationalText, err = decodeInternationalText(value)
		il.SortShow = &internationalText
	case "(c)ART":
		internationalText, err = decodeInternationalText(value)
		il.Artist = &internationalText
	case "(c)alb":
		internationalText, err = decodeInternationalText(value)
		il.AlbumC = &internationalText
	case "(c)ard":
		internationalText, err = decodeInternationalText(value)
		il.ArtDirector = &internationalText
	case "(c)arg":
		internationalText, err = decodeInternationalText(value)
		il.Arranger = &internationalText
	case "(c)aut":
		internationalText, err = decodeInternationalText(value)
		il.AuthorC = &internationalText
	case "(c)cmt":
		internationalText, err = decodeInternationalText(value)
		il.Comment = &internationalText
	case "(c)com":
		internationalText, err = decodeInternationalText(value)
		il.ComposerC = &internationalText
	case "(c)con":
		internationalText, err = decodeInternationalText(value)
		il.Conductor = &internationalText
	case "(c)cpy":
		internationalText, err = decodeInternationalText(value)
		il.CopyrightC = &internationalText
	case "(c)day":
		internationalText, err = decodeInternationalText(value)
		il.ContentCreateDate = &internationalText
	case "(c)des":
		internationalText, err = decodeInternationalText(value)
		il.DescriptionC = &internationalText
	case "(c)dir":
		internationalText, err = decodeInternationalText(value)
		il.Director = &internationalText
	case "(c)enc":
		internationalText, err = decodeInternationalText(value)
		il.EncodedBy = &internationalText
	case "(c)gen":
		internationalText, err = decodeInternationalText(value)
		il.GenreC = &internationalText
	case "(c)grp":
		internationalText, err = decodeInternationalText(value)
		il.GroupingC = &internationalText
	case "(c)lyr":
		internationalText, err = decodeInternationalText(value)
		il.Lyrics = &internationalText
	case "(c)nam":
		internationalText, err = decodeInternationalText(value)
		il.TitleC = &internationalText
	case "(c)nrt":
		internationalText, err = decodeInternationalText(value)
		il.Narrator = &internationalText
	case "(c)ope":
		internationalText, err = decodeInternationalText(value)
		il.OriginalArtist = &internationalText
	case "(c)prd":
		internationalText, err = decodeInternationalText(value)
		il.Producer = &internationalText
	case "(c)pub":
		internationalText, err = decodeInternationalText(value)
		il.Publisher = &internationalText
	case "(c)sne":
		internationalText, err = decodeInternationalText(value)
		il.SoundEngineer = &internationalText
	case "(c)sol":
		internationalText, err = decodeInternationalText(value)
		il.Soloist = &internationalText
	case "(c)st3":
		internationalText, err = decodeInternationalText(value)
		il.Subtitle = &internationalText
	case "(c)too":
		internationalText, err = decodeInternationalText(value)
		il.Encoder = &internationalText
	case "(c)trk":
		internationalText, err = decodeInternationalText(value)
		il.Track = &internationalText
	case "(c)wrk":
		internationalText, err = decodeInternationalText(value)
		il.Work = &internationalText
	case "(c)wrt":
		internationalText, err = decodeInternationalText(value)
		il.ComposerCWRT = &internationalText
	case "(c)xpd":
		internationalText, err = decodeInternationalText(value)
		il.ExecutiveProducer = &internationalText
	case "(c)xyz":
		internationalText, err = decodeInternationalText(value)
		il.GPSCoordinates = &internationalText
	case "(c)mvn":
		internationalText, err = decodeInternationalText(value)
		il.MovementName = &internationalText

	// string
	case "@PST":
		il.ParentShortTitle = &str
	case "@ppi":
		il.ParentProductID = &str
	case "@pti":
		il.ParentTitle = &str
	case "@sti":
		il.ShortTitle = &str
	case "AACR":
		il.UnknownAACR = &str
	case "CDEK":
		il.UnknownCDEK = &str
	case "CDET":
		il.UnknownCDET = &str
	case "GUID":
		il.GUID = &str
	case "VERS":
		il.ProductVersion = &str
	case "albm":
		il.Album = &str
	case "apID":
		il.AppleStoreAccount = &str
	case "auth":
		il.Author = &str
	case "catg":
		il.Category = &str
	case "cmID":
		il.ComposerID = &str
	case "covr":
		il.CoverArt = &str
	case "grup":
		il.Grouping = &str
	case "gshh":
		il.GoogleHostHeader = &str
	case "gspm":
		il.GooglePingMessage = &str
	case "gspu":
		il.GooglePingURL = &str
	case "gssd":
		il.GoogleSourceData = &str
	case "gsst":
		il.GoogleStartTime = &str
	case "gstd":
		il.GoogleTrackDuration = &str
	case "keyw":
		il.Keyword = &str
	case "ldes":
		il.LongDescription = &str
	case "ownr":
		il.Owner = &str
	case "egid":
		il.EpisodeGlobalUniqueID = &str
	case "perf":
		il.Performer = &str
	case "prID":
		il.ProductID = &str
	case "purd":
		il.PurchaseDate = &str
	case "purl":
		il.PodcastURL = &str
	case "rldt":
		il.ReleaseDate = &str
	case "snal":
		il.PreviewImage = &str
	case "titl":
		il.Title = &str
	case "tven":
		il.TVEpisodeID = &str
	case "tvsh":
		il.TVShow = &str
	case "xid ":
		il.ISRC = &str
	case "yrrc":
		il.Year = &str

	// int8
	case "akID":
		il.AppleStoreAccountType = new(int8)
		_, err = binary.Encode(value, binary.BigEndian, il.AppleStoreAccountType)
	case "rate":
		il.Rating = new(int8)
		_, err = binary.Encode(value, binary.BigEndian, il.Rating)
	case "stik":
		il.MediaType = new(int8)
		_, err = binary.Encode(value, binary.BigEndian, il.MediaType)

	// int16
	case "tmpo":
		il.BeatsPerMinute = new(int16)
		_, err = binary.Encode(value, binary.BigEndian, il.BeatsPerMinute)
	case "(c)mvc":
		il.MovementCount = new(int16)
		_, err = binary.Encode(value, binary.BigEndian, il.MovementCount)
	case "(c)mvi":
		il.MovementNumber = new(int16)
		_, err = binary.Encode(value, binary.BigEndian, il.MovementNumber)

	// int32
	case "atID":
		il.ArtistID = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.ArtistID)
	case "cnID":
		il.AppleStoreCatalogID = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.AppleStoreCatalogID)
	case "disk":
		il.DiskNumber = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.DiskNumber)
	case "geID":
		il.GenreID = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.GenreID)
	case "gnre":
		il.Genre = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.Genre)
	case "sfID":
		il.AppleStoreCountry = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.AppleStoreCountry)
	case "trkn":
		il.TrackNumber = new(int32)
		il.TrackNumber = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.TrackNumber)
	case "tves":
		il.TVEpisode = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.TVEpisode)
	case "tvsn":
		il.TVSeason = new(int32)
		_, err = binary.Encode(value, binary.BigEndian, il.TVSeason)

	// bool
	case "cpil":
		il.Compilation = new(bool)
		_, err = binary.Encode(value, binary.BigEndian, il.Compilation)
	case "hdvd":
		il.HDVideo = new(bool)
		_, err = binary.Encode(value, binary.BigEndian, il.HDVideo)
	case "itnu":
		il.ITunesU = new(bool)
		_, err = binary.Encode(value, binary.BigEndian, il.ITunesU)
	case "pcst":
		il.Podcast = new(bool)
		_, err = binary.Encode(value, binary.BigEndian, il.Podcast)
	case "pgap":
		il.DisableInsertPlayGap = new(bool)
		_, err = binary.Encode(value, binary.BigEndian, il.DisableInsertPlayGap)
	case "shwm":
		il.ShowMovement = new(bool)
		_, err = binary.Encode(value, binary.BigEndian, il.ShowMovement)
	}
	return err
}
