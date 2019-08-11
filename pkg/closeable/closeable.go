package closeable

import (
	"io"

	"go.uber.org/zap"
)

type CloseLogger interface {
	With(args ...interface{}) *zap.SugaredLogger
}

func CloseObjects(l CloseLogger, cls ...io.Closer) {
	for _, c := range cls {
		CloseWithErrorLogging(l, c)
	}
}

func CloseWithErrorLogging(l CloseLogger, cl io.Closer) {
	if cl != nil {
		if err := cl.Close(); err != nil {
			l.With("error", err).Error("error while closing")
		}
	}
}
