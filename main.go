package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"google.golang.org/protobuf/proto"
)

var logger = getLogger()

func getFeedMessage() *FeedMessage {
	ts := time.Now().Unix()
	resp, err := http.NewRequest(http.MethodGet, os.Getenv("GTFS-RT_ENDPOINT"), nil)
	resp.Header.Add("Authorization", "Bearer "+os.Getenv("GTFS_RT_API_KEY"))
	resp.Header.Add("User-Agent", "TrailateService")

	if err != nil {
		logger.Error("Error creating HTTP request", slog.String("error", err.Error()))
	}

	client := &http.Client{}
	response, err := client.Do(resp)

	if err != nil {
		logger.Error("Error creating HTTP request", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if response.StatusCode != http.StatusOK {
		logger.Error("Non-200 response: " + response.Status)
		os.Exit(1)
	}

	// Ensure the response body is closed after reading
	defer response.Body.Close()

	// Read response body into byte array
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error creating HTTP request", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Compress and write to file
	compBodBytes := compressData(bodyBytes)
	compressedFileName := fmt.Sprintf("%d_gtfs-rt.pb.gz", ts)

	err = os.WriteFile(compressedFileName, compBodBytes, 0644)

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Wrote %d bytes to file %s", len(compBodBytes), compressedFileName))

	// Create a new FeedMessage instance
	feedMessage := &FeedMessage{}

	// Unmarshal the protobuf data
	if err := proto.Unmarshal(bodyBytes, feedMessage); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	return feedMessage
}

func getNumberOfTripUpdates(feedMessage *FeedMessage) int {
	sum_updates := 0

	for _, entity := range feedMessage.Entity {
		if entity.TripUpdate == nil {
			continue
		}
		sum_updates += len(entity.TripUpdate.StopTimeUpdate)
	}

	return sum_updates
}

func getNumberOfEntities(feedMessage *FeedMessage) int {
	return len(feedMessage.Entity)
}

func main() {
	feedMessage := getFeedMessage()

	updateCount := getNumberOfTripUpdates(feedMessage)
	entityCount := getNumberOfEntities(feedMessage)

	logger.Info("Number of entities:", slog.Int("entityCount", entityCount))
	logger.Info("Number of trip updates:", slog.Int("updateCount", updateCount))
}
