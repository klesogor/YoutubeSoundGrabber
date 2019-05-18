package internals

import (
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
	concurrency, chunkSize := getSettings(config)
	totalDownloads := download.GetContentLengh() / chunkSize
	sem := semaphore.NewWeighted(int64(concurrency))
	buffer := make([][]byte, totalDownloads)
	url := download.GetDownloadUrl()
	res, err := handler(url, 0, chunkSize)
	if err != nil {
		return nil, err
	}
	buffer[0] = res
	wg.Add(totalDownloads - 1)
	for i := 1; i < totalDownloads; i++ {
		sem.TryAcquire(1)
		go func(offset int) {
			start := offset*chunkSize + 1
			res, downloadError := handler(url, start, start+chunkSize-1)
			if downloadError != nil {
				err = downloadError
				wg.Done()
				sem.Release(1)
				return
			}
			buffer[offset] = res
			wg.Done()
			sem.Release(1)
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
