package define

// 文件MIME类型对应扩展名
var MIME_EXT = map[string]string{
	// 文档文件类型的
	"application/postscript":                  ".ai",
	"application/octet-stream":                ".exe",
	"application/vnd.ms-word":                 ".doc",
	"application/vnd.ms-excel":                ".xls",
	"application/vnd.ms-powerpoint":           ".ppt",
	"application/pdf":                         ".pdf",
	"application/xml":                         ".xml",
	"application/vnd.oasis.opendocument.text": ".odt",
	"application/x-shockwave-flash":           ".swf",

	// 压缩文件类型的
	"application/x-gzip":          ".gz",
	"application/x-bzip2":         ".bz",
	"application/zip":             ".zip",
	"application/x-rar":           ".rar",
	"application/x-tar":           ".tar",
	"application/x-7z-compressed": ".7",

	// 文字类型
	"text/plain":         ".txt",
	"text/x-php":         ".php",
	"text/html":          ".html",
	"text/javascript":    ".js",
	"text/css":           ".css",
	"text/rtf":           ".rtf",
	"text/rtfd":          ".rtfd",
	"text/x-python":      ".py",
	"text/x-java-source": ".java",
	"text/x-ruby":        ".rb",
	"text/x-shellscript": ".sh",
	"text/x-perl":        ".pl",
	"text/x-sql":         ".sql",

	// 图片类型的
	"image/x-ms-bmp":            ".bmp",
	"image/jpeg":                ".jpg",
	"image/gif":                 ".gif",
	"image/png":                 ".png",
	"image/tiff":                ".tif",
	"image/x-targa":             ".tga",
	"image/vnd.adobe.photoshop": ".psd",
	"text/html; charset=utf-8":  ".svg",
	"text/plain; charset=utf-8": ".svg",
	"image/svg+xml":             ".svg",

	// 音频文件类型的
	"audio/mpeg":     ".mp3",
	"audio/midi":     ".mid",
	"audio/ogg":      ".ogg",
	"audio/mp4":      ".mp4a",
	"audio/wav":      ".wav",
	"audio/x-ms-wma": ".wma",

	// 视频文件类型的
	"video/x-msvideo":  ".avi",
	"video/x-dv":       ".dv",
	"video/mp4":        ".mp4",
	"video/mpeg":       ".mpg",
	"video/quicktime":  ".mov",
	"video/x-ms-wmv":   ".wm",
	"video/x-flv":      ".flv",
	"video/x-matroska": ".mkv",
}
