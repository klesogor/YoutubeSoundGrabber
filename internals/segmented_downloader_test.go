package internals

import (
	"errors"
	"testing"
)

var mockData = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
var downloadConfig = SegmentedDownloadConfig{ConcurrencyLimit: 5, ChunkSize: 4}

type downloadMock struct {
}

func (d downloadMock) GetDownloadUrl() string {
	return ""
}
func (d downloadMock) GetContentLengh() int {
	return len(mockData)
}

func mockDownloaderSuccess(url string, start, end int) ([]byte, error) {
	if start > len(mockData) {
		return make([]byte, 0), nil
	}
	end++
	if end > len(mockData) {
		end = len(mockData)
	}
	return mockData[start:end], nil
}

func TestSegmentedDownloadPositive(t *testing.T) {
	data, err := DownloadSegmented(downloadMock{}, downloadConfig, mockDownloaderSuccess)
	if err != nil {
		t.Error(err.Error())
	}
	for k, _ := range data {
		if data[k] != mockData[k] {
			t.Errorf("Expected data: %v\n Received data: %v", mockData, err)
		}
	}
}

func TestSegmentedDownloadFailFirstByte(t *testing.T) {
	callCount := 0
	_, err := DownloadSegmented(downloadMock{}, downloadConfig, func(url string, start, end int) ([]byte, error) {
		callCount++
		return nil, errors.New("Failed to load data!")
	})

	if err == nil {
		t.Error("Err is nil!!!")
	}

	if callCount > 1 {
		t.Error("Called download function more than once! Wasted resources!")
	}
}

func TestSegmentedDownloadFailSecondByte(t *testing.T) {
	callCount := 0
	_, err := DownloadSegmented(downloadMock{}, downloadConfig, func(url string, start, end int) ([]byte, error) {
		callCount++
		if callCount > 2 {
			return nil, errors.New("Failed to load data!")
		}
		return make([]byte, 1), nil
	})

	if err == nil {
		t.Error("Err is nil!!!")
	}
}
