package multiregex

import (
	"context"
	"io"
)

// multiplyStream returns multiple readers that can read from the same reader
// This function will probably deadlock if all new streams are not read asynchronously
func multiplyStream(ctx context.Context, reader io.ReadCloser, numStreams int) []io.ReadCloser {
	workerWriters := []*io.PipeWriter{}
	workerReaders := []*io.PipeReader{}
	for i := 0; i < numStreams; i++ {
		r, w := io.Pipe()
		workerWriters = append(workerWriters, w)
		workerReaders = append(workerReaders, r)
	}

	// Read stream duplicator to write to all workerWriters
	go func() {
		// Defer closing all workers
		closeAll := func() {
			for _, writer := range workerWriters {
				writer.Close()
			}
		}
		defer closeAll()

		for {
			// Read
			buf := make([]byte, 1024)
			if _, err := reader.Read(buf); err != nil {
				return
			}

			// Writer to each worker stream
			for _, writer := range workerWriters {
				writer.Write(buf)
			}

			// Check if we are canceled
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	// Cast to io.ReadClosers
	readers := make([]io.ReadCloser, numStreams)
	for i := range readers {
		readers[i] = workerReaders[i]
	}

	return readers
}
