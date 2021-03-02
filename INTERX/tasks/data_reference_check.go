package tasks

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	common "github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	database "github.com/KiraCore/sekai/INTERX/database"
	"github.com/KiraCore/sekai/INTERX/global"
)

// RefMeta is a struct to be used for reference metadata
type RefMeta struct {
	ContentLength int64     `json:"content_length"`
	LastModified  time.Time `json:"last_modified"`
}

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
	path = config.GetReferenceCacheDir() + "/" + path
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

		global.Mutex.Lock()

		if _, err := os.Stat(filepath.Dir(path)); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(path), 0700)
		}

		err = ioutil.WriteFile(path, bodyBytes, 0644)
		if err != nil {
			global.Mutex.Unlock()
			return err
		}

		global.Mutex.Unlock()
	}

	return nil
}

// DataReferenceCheck is a function to check cache data for data references.
func DataReferenceCheck(isLog bool) {
	for {
		references, err := database.GetAllReferences()
		if err == nil {
			for _, v := range references {
				ref, err := getMeta(v.URL)
				if err != nil {
					continue
				}

				// Check if reference has changed (check length and last modified)
				if v.ContentLength == ref.ContentLength && ref.LastModified.Equal(v.LastModified) {
					continue
				}

				// Check the download file size limitation
				if ref.ContentLength > config.Config.Cache.DownloadFileSizeLimitation {
					continue
				}

				err = saveReference(v.URL, v.FilePath)
				if err != nil {
					continue
				}

				if isLog {
					common.GetLogger().Info("[cache] Data reference updated")
					common.GetLogger().Info("[cache] Key = ", v.Key)
					common.GetLogger().Info("[cache] Ref = ", v.URL)
				}

				database.AddReference(v.Key, v.URL, ref.ContentLength, ref.LastModified, v.FilePath)
			}
		}

		time.Sleep(2 * time.Second)
	}
}
