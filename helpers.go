package main

import (
	"bytes"
	"compress/gzip"
	"log/slog"
	"os"
)

func getLogger() *slog.Logger {
	// Custom handler to format output
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug, // minimum level
		AddSource: true,            // include file + line
	})
	logger := slog.New(handler)

	return logger
}

func compressData(data []byte) []byte {
	// Create a buffer to hold the compressed data
	var buf bytes.Buffer

	// Create a new gzip writer
	gz := gzip.NewWriter(&buf)
	defer gz.Close()

	// Write the data to the gzip writer
	if _, err := gz.Write(data); err != nil {
		logger.Error(err.Error())
	}
	// Close the gzip writer to flush any remaining data
	if err := gz.Close(); err != nil {
		logger.Error(err.Error())
	}

	// Return the compressed data
	return buf.Bytes()
}
