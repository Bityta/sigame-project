package media

import (
	"github.com/google/uuid"
)

const (
	MediaTypeText  = "text"
	MediaTypeImage = "image"
	MediaTypeAudio = "audio"
	MediaTypeVideo = "video"

	MediaSizeImage   = 500_000
	MediaSizeAudio   = 3_000_000
	MediaSizeVideo   = 10_000_000
	MediaSizeDefault = 100_000
)

func buildMediaID(round, themeIndex, price int) string {
	return uuid.NewString()
}

func estimateMediaSize(mediaType string) int64 {
	switch mediaType {
	case MediaTypeImage:
		return MediaSizeImage
	case MediaTypeAudio:
		return MediaSizeAudio
	case MediaTypeVideo:
		return MediaSizeVideo
	default:
		return MediaSizeDefault
	}
}

