package ilst

import "encoding/binary"

func (il *ItemList) Set(id string, value []byte) (err error) {
	str := string(value)
	switch id {
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
	case "aART":
		il.AlbumArtist = &str
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
	case "cprt":
		il.Copyright = &str
	case "desc":
		il.Description = &str
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
	case "soaa":
		il.SortAlbumArtist = &str
	case "soal":
		il.SortAlbum = &str
	case "soar":
		il.SortArtist = &str
	case "soco":
		il.SortComposer = &str
	case "sonm":
		il.SortName = &str
	case "sosn":
		il.SortShow = &str
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
	case "(c)ART":
		il.Artist = &str
	case "(c)alb":
		il.AlbumC = &str
	case "(c)ard":
		il.ArtDirector = &str
	case "(c)arg":
		il.Arranger = &str
	case "(c)aut":
		il.AuthorC = &str
	case "(c)cmt":
		il.Comment = &str
	case "(c)com":
		il.ComposerC = &str
	case "(c)con":
		il.Conductor = &str
	case "(c)cpy":
		il.CopyrightC = &str
	case "(c)day":
		il.ContentCreateDate = &str
	case "(c)des":
		il.DescriptionC = &str
	case "(c)dir":
		il.Director = &str
	case "(c)enc":
		il.EncodedBy = &str
	case "(c)gen":
		il.GenreC = &str
	case "(c)grp":
		il.GroupingC = &str
	case "(c)lyr":
		il.Lyrics = &str
	case "sdes":
		il.StoreDescription = &str
	case "(c)nam":
		il.TitleC = &str
	case "(c)nrt":
		il.Narrator = &str
	case "(c)ope":
		il.OriginalArtist = &str
	case "(c)prd":
		il.Producer = &str
	case "(c)pub":
		il.Publisher = &str
	case "(c)sne":
		il.SoundEngineer = &str
	case "(c)sol":
		il.Soloist = &str
	case "(c)st3":
		il.Subtitle = &str
	case "(c)too":
		il.Encoder = &str
	case "(c)trk":
		il.Track = &str
	case "(c)wrk":
		il.Work = &str
	case "(c)wrt":
		il.ComposerCWRT = &str
	case "(c)xpd":
		il.ExecutiveProducer = &str
	case "(c)xyz":
		il.GPSCoordinates = &str
	case "tnal":
		il.ThumbnailImage = &str
	case "tvnn":
		il.TVNetworkName = &str
	case "(c)mvn":
		il.MovementName = &str

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
