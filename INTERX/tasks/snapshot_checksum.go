package tasks

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"time"

	"github.com/KiraCore/sekai/INTERX/common"
	"github.com/KiraCore/sekai/INTERX/config"
)

var (
	SnapshotChecksumAvailable bool = false
	SnapshotLength            int64
	SnapshotChecksum          string
	SnapshotModTime           time.Time
)

func calcChecksum(isLog bool) {
	SnapshotChecksumAvailable = false
	SnapshotLength = 0

	common.GetLogger().Info("[cache] calculating snapshot checksum: ")

	f, err := os.Open(config.SnapshotPath())
	if err != nil {
		if isLog {
			common.GetLogger().Error("[cache] can't read snapshot file: ", err)
		}

		return
	}

	defer func() {
		_ = f.Close()
	}()

	buf := make([]byte, 1024*1024)
	h := sha256.New()

	var totalRead int64 = 0
	for {
		bytesRead, err := f.Read(buf)
		if err != nil {
			if err != io.EOF {
				if isLog {
					common.GetLogger().Error("[cache] failed to read snapshot: ", err)
				}
				return
			}

			break
		}

		// do some other work with buf before adding it to the hasher
		// processBuffer(buf)

		h.Write(buf[:bytesRead])
		totalRead += int64(bytesRead)
	}

	SnapshotChecksumAvailable = true
	SnapshotChecksum = hex.EncodeToString(h.Sum(nil))
	SnapshotLength = totalRead
	common.GetLogger().Info("[cache] snapshot checksum: ", SnapshotChecksum)
}

// CalcSnapshotChecksum is a function for syncing sekaid status.
func CalcSnapshotChecksum(isLog bool) {
	available := 0
	for {
		file, err := os.Stat(config.SnapshotPath())

		if err != nil {
			if available != 1 && isLog {
				common.GetLogger().Error("[cache] can't read snapshot file: ", err)
				available = 1
			}

			continue
		}

		available = 0

		if file.ModTime().Equal(SnapshotModTime) {
			continue
		}

		SnapshotModTime = file.ModTime()

		calcChecksum(isLog)
	}
}
