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

## CLI Usage

```sh
make build
```

### Read

```sh
./probe -f /path/to/music.m4a
```

## References

- [QuickTime File Format | Apple Developer Documentation](https://developer.apple.com/documentation/quicktime-file-format)
- [QuickTime Tags (ItemList)](https://exiftool.org/TagNames/QuickTime.html#ItemList)
