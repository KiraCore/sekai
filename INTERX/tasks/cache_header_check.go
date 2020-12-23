package tasks

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	common "github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
)

// CacheHeaderCheck is a function to check cache headers if it's expired.
func CacheHeaderCheck(rpcAddr string, isLog bool) {
	for {
		err := filepath.Walk(interx.GetResponseCacheDir(),
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// check file header, delete empty directory or expired cache
				delete := false

				if info.IsDir() {
					files, err := ioutil.ReadDir(path)
					if err == nil && len(files) == 0 {
						delete = true
					}
				} else if info.Size() == 0 || info.ModTime().Add(time.Duration(interx.Config.CachingDuration)*time.Second).Before(time.Now()) {
					delete = true
				}

				if path != interx.GetResponseCacheDir() && delete {
					if isLog {
						fmt.Println("deleting file ... ", path)
					}

					common.Mutex.Lock()
					err := os.Remove(path)
					common.Mutex.Unlock()

					if err != nil {
						if isLog {
							fmt.Println("Error deleting file: ", err)
						}
						return err
					}

					return nil
				}

				return nil
			})

		if err != nil {
			log.Println(err)
		}
	}
}
