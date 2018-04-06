package file

const (
	Video           = "video"
	Image           = "image"
	Content         = "content"
	MaxFileSize     = int64(100 * 1000 * 1024)
	MaxBodySize     = int64(100 * 1000 * 4096)
	defaultImageExt = "jpg"
)

var (
	Format = map[string]string{
		"jpeg": Image, "jpg": Image, "png": Image, "avi": Video,
		"mkv": Video, "mp4": Video, "ogg": Video, "webm": Video,
		"doc": Content, "docx": Content, "xls": Content, "xlsx": Content,
		"pdf": Content, "zip": Content, "tar.gz": Content,
	}
)
