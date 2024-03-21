package isoeditor

import "io"

type ignitionImageReader struct {
	io.Reader
	io.Closer
}

// NewIgnitionImageReader returns the filename of the ignition image in the ISO,
// along with a stream of the ignition image with ignition content embedded.
// This can be used to overwrite the ignition image file of an ISO previously
// unpacked by Extract() in order to embed ignition data.
func NewIgnitionImageReader(isoPath string, ignitionContent *IgnitionContent) (string, io.ReadCloser, error) {
	info, iso, err := ignitionOverlay(isoPath, ignitionContent, true)
	if err != nil {
		return "", nil, err
	}
	imageOffset, imageLength, err := GetISOFileInfo(info.File, isoPath)
	if err != nil {
		return "", nil, err
	}

	length := info.Offset + info.Length
	// include any trailing data
	if imageLength > length {
		length = imageLength
	}

	if _, err := iso.Seek(imageOffset, io.SeekStart); err != nil {
		iso.Close()
		return "", nil, err
	}
	reader := ignitionImageReader{
		Reader: io.LimitReader(iso, length),
		Closer: iso,
	}
	return info.File, &reader, nil
}
