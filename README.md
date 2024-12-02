# qtffilst

QuickTime File Format `moov.udta.meta.ilst` tag library.

## Library Usage

```sh
go get github.com/tingtt/qtffilst
```

```go
import "github.com/tingtt/qtffilst"
```

### Read

Opening is omitted from the examples.

```go
file, _ := os.Open("/path/to/track.m4a")
defer file.Close()

r, err := qtffilst.NewReader(file)
if err != nil {
	return err
}

itemListTag, err := r.Read()
if err != nil {
	return err
}
```

#### Read album title

```go
itemListTag, err := r.Read()
if err != nil {
	panic(err)
}
fmt.Println(itemListTag.AlbumC.Text)
```

### Write

```go
file, _ := os.Open("/path/to/track.m4a")
defer file.Close()
tmp1, _ := os.Create("tmp1.m4a")
defer tmp1.Close()
defer func() { os.Remove(tmp1.Name()) } ()
tmp2, _ := os.Create("tmp2.m4a")
defer tmp2.Close()
defer func() { os.Remove(tmp2.Name()) } ()
dest, _ := os.Create("dest.m4a")
defer dest.Close()

rw, err := qtffilst.ParseReadWriter(file)
if err != nil {
	return err
}

// Sample: Set new title and remove subtitle.
err = rw.Write(dest, tmp1, tmp2,
	ilst.ItemList{TitleC: ilst.NewInternationalText("New title")},
	/* delete ilst */ []string{ /* subtitle */ "(c)st3"},
)
if err != nil {
	return err
}
```

## CLI Usage

```sh
make build
```

### probe

```sh
qtffprobe -f /path/to/music.m4a
```

### edit

```sh
# Edit compilation title
qtffilst -f /path/to/music.m4a -o out.m4a -d "(c)nam=Title"

# Remove compilation title
qtffilst -f /path/to/music.m4a -o out.m4a -r "(c)nam"
```

## References

- [QuickTime File Format | Apple Developer Documentation](https://developer.apple.com/documentation/quicktime-file-format)
- [QuickTime Tags (ItemList)](https://exiftool.org/TagNames/QuickTime.html#ItemList)
