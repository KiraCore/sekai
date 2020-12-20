package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	common "github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
)

// RefMeta is a struct to be used for reference metadata
type RefMeta struct {
	ContentLength int64     `json:"content_length"`
	LastModified  time.Time `json:"last_modified"`
}

// RefCache is a struct to be used for saving reference metadata
type RefCache struct {
	Path   string      `json:"path"`
	Header http.Header `json:"header"`
}

var references = make(map[string]RefMeta)

func getMeta(url string) (*RefMeta, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	contentLength, err := strconv.Atoi(resp.Header["Content-Length"][0])
	if err != nil {
		return nil, err
	}
	lastModified, err := time.Parse(time.RFC1123, resp.Header["Last-Modified"][0])
	if err != nil {
		return nil, err
	}
	return &RefMeta{
		ContentLength: int64(contentLength),
		LastModified:  lastModified,
	}, nil
}

func saveReference(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		common.Mutex.Lock()

		err = ioutil.WriteFile(path, bodyBytes, 0644)
		if err != nil {
			return err
		}

		common.Mutex.Unlock()

		cache := RefCache{
			Path:   path,
			Header: resp.Header,
		}

		data, err := json.Marshal(cache)
		if err != nil {
			return err
		}

		common.Mutex.Lock()

		err = ioutil.WriteFile(path+".meta", data, 0644)
		if err != nil {
			return err
		}

		common.Mutex.Unlock()
	}

	return nil
}

// LoadRefCacheMeta is a function to load reference cache
func LoadRefCacheMeta(key string) (RefCache, error) {
	filePath := interx.Config.CacheDir + "/reference/" + key + ".meta"

	common.Mutex.Lock()
	data, err := ioutil.ReadFile(filePath)
	common.Mutex.Unlock()

	cache := RefCache{}

	if err == nil {
		err = json.Unmarshal([]byte(data), &cache)
	}

	return cache, err
}

// DataReferenceCheck is a function to check cache data for data references.
func DataReferenceCheck(isLog bool) {
	for {
		for k, v := range common.DataRefs {
			ref, err := getMeta(v.Reference)
			if err != nil {
				continue
			}

			// Check if reference has changed (check length and last modified)
			if references[k].ContentLength == ref.ContentLength && ref.LastModified.Equal(references[k].LastModified) {
				continue
			}

			fmt.Println(ref, interx.Config.DownloadFileSizeLimitation)
			// Check the download file size limitation
			if ref.ContentLength > interx.Config.DownloadFileSizeLimitation {
				continue
			}

			err = saveReference(v.Reference, interx.Config.CacheDir+"/reference/"+k)
			if err != nil {
				continue
			}

			if isLog {
				fmt.Println("save response: (key - " + k + ", ref - " + v.Reference + ")")
			}

			references[k] = RefMeta{
				ContentLength: ref.ContentLength,
				LastModified:  ref.LastModified,
			}
		}
	}
}
