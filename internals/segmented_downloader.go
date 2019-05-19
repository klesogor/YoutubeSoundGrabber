package internals

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/semaphore"
)

type SegmentDownloadable interface {
	GetDownloadUrl() string
	GetContentLengh() int
}

type SegmentedDownloadConfig struct {
	ConcurrencyLimit int
	ChunkSize        int
}

type DownloadHandler func(url string, start, end int) ([]byte, error)

func DownloadSegmented(download SegmentDownloadable, config SegmentedDownloadConfig, handler DownloadHandler) ([]byte, error) {
	var err error
	var wg sync.WaitGroup
	ctx := context.Background()
	concurrency, chunkSize := getSettings(config)
	totalDownloads := download.GetContentLengh() / chunkSize
	sem := semaphore.NewWeighted(int64(concurrency))
	buffer := make([][]byte, totalDownloads+1)
	url := download.GetDownloadUrl()
	res, err := handler(url, 0, chunkSize)
	if err != nil {
		return nil, err
	}
	buffer[0] = res
	wg.Add(totalDownloads)
	for i := 1; i < totalDownloads+1; i++ {
		sem.Acquire(ctx, 1)
		go func(offset int) {
			start := offset*chunkSize + offset
			fmt.Printf("Aquired sem, range: %d-%d\n", start, start+chunkSize)
			res, downloadError := handler(url, start, start+chunkSize)
			if downloadError != nil {
				err = downloadError
				wg.Done()
				sem.Release(1)
				fmt.Printf("Released 1 concurency due to error, range: %d-%d\n", start, start+chunkSize)
				return
			}
			buffer[offset] = res
			wg.Done()
			sem.Release(1)
			fmt.Printf("Released 1 concurency due to success, range: %d-%d\n", start, start+chunkSize)
		}(i)
	}
	wg.Wait()
	if err != nil {
		return nil, err
	}
	result := make([]byte, 0, download.GetContentLengh())
	for _, v := range buffer {
		result = append(result, v...)
	}
	return result, nil
}

func getSettings(config SegmentedDownloadConfig) (conc int, chunk int) {
	if config.ConcurrencyLimit == 0 {
		conc = 100
	} else {
		conc = config.ConcurrencyLimit
	}

	if config.ChunkSize == 0 {
		chunk = 989898
	} else {
		chunk = config.ChunkSize
	}

	return
}
