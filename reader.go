package pgzip

import "io"

func CompressingReader(r io.Reader) io.Reader {
	return CompressingReaderLevel(r, DefaultCompression)
}

func CompressingReaderLevel(r io.Reader, level int) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()

		gzw, err := NewWriterLevel(pw, level)
		defer gzw.Close()
		if err != nil {
			pw.CloseWithError(err)
		}

		_, err = io.Copy(gzw, r)
		if err != nil {
			pw.CloseWithError(err)
		}
	}()
	return pr
}
