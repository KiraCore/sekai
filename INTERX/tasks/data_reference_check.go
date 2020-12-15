package tasks

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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

func saveReference(url string, dir string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	common.Mutex.Lock()

	file, err := os.Create(dir + path.Base(url))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	common.Mutex.Unlock()

	return nil
}

// DataReferenceCheck is a function to check cache data for data references.
func DataReferenceCheck(isLog bool) {
	for {
		for k, v := range common.DataRefs {
			ref, err := getMeta(v.Reference)
			if err != nil {
				continue
			}

			if references[k].ContentLength == ref.ContentLength && ref.LastModified.Equal(references[k].LastModified) {
				continue
			}

			if isLog {
				fmt.Println(k, "(", v.Reference, ") - ContentLength: ", ref.ContentLength, " Last Modified", ref.LastModified)
			}

			err = saveReference(v.Reference, interx.Config.CacheDir+"/reference/")
			if err != nil {
				continue
			}

			references[k] = RefMeta{
				ContentLength: ref.ContentLength,
				LastModified:  ref.LastModified,
			}
		}
	}
}
