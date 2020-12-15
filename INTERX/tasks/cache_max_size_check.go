package tasks

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	common "github.com/KiraCore/sekai/INTERX/common"
	interx "github.com/KiraCore/sekai/INTERX/config"
)

// CacheMaxSizeCheck is a function to check if cache reached the maximum size.
func CacheMaxSizeCheck(isLog bool) {
	for {
		var cacheSize int64 = 0
		_ = filepath.Walk(interx.Config.CacheDir+"/response", func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				cacheSize += info.Size()
			}
			return err
		})

		if cacheSize >= interx.Config.MaxCacheSize {
			if isLog {
				fmt.Println("Cache reached the maximum size")
			}

			for {
				_ = filepath.Walk(interx.Config.CacheDir+"/response", func(path string, info os.FileInfo, err error) error {
					if err != nil || cacheSize*10 < interx.Config.MaxCacheSize*9 { // current size < 90% of max cache size
						return err
					}
					if !info.IsDir() && rand.Intn(5) == 0 {
						cacheSize -= info.Size()

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
					return err
				})

				if cacheSize*10 < interx.Config.MaxCacheSize*9 {
					break
				}
			}
		}
	}
}
