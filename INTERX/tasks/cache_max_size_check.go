package tasks

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"

	common "github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
	"github.com/KiraCore/sekai/INTERX/global"
)

// CacheMaxSizeCheck is a function to check if cache reached the maximum size.
func CacheMaxSizeCheck(isLog bool) {
	for {
		var cacheSize int64 = 0
		_ = filepath.Walk(config.GetResponseCacheDir(), func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				cacheSize += info.Size()
			}
			return err
		})

		if cacheSize >= config.Config.Cache.MaxCacheSize {
			if isLog {
				common.GetLogger().Info("[cache] Reached the maximum size")
			}

			for {
				_ = filepath.Walk(config.GetResponseCacheDir(), func(path string, info os.FileInfo, err error) error {
					if _, err := os.Stat(path); os.IsNotExist(err) {
						return nil
					}

					if err != nil || cacheSize*10 < config.Config.Cache.MaxCacheSize*9 { // current size < 90% of max cache size
						return err
					}
					if !info.IsDir() && rand.Intn(5) == 0 {
						cacheSize -= info.Size()

						if isLog {
							common.GetLogger().Info("[cache] Deleting file: ", path)
						}

						global.Mutex.Lock()
						// check if file or path exists
						if _, err := os.Stat(path); os.IsNotExist(err) {
							global.Mutex.Unlock()
							return nil
						}
						err := os.Remove(path)
						global.Mutex.Unlock()

						if err != nil {
							if isLog {
								common.GetLogger().Error("[cache] Error deleting file: ", err)
							}
							return err
						}

						return nil
					}
					return err
				})

				if cacheSize*10 < config.Config.Cache.MaxCacheSize*9 {
					break
				}
			}
		}

		time.Sleep(2 * time.Second)
	}
}
