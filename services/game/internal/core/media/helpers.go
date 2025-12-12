package media

import "github.com/google/uuid"

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

