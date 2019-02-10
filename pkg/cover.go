package pkg

import (
	"io"
	"sync"
	"sync/atomic"
)

type tCover struct {
	fileName string
	counter  []uint32
	pos      []uint32
	numStmts []uint16
}

type tCovers []*tCover

var covers_ = make(tCovers, 0)
var coversMu sync.Mutex

func RegisterCover(fileName string, counter []uint32, pos []uint32, numStmts []uint16) {
	coversMu.Lock()
	defer coversMu.Unlock()
	covers_ = append(covers_, &tCover{
		fileName: fileName,
		counter:  counter,
		pos:      pos,
		numStmts: numStmts,
	})
}

func makeCoverProfile(covers tCovers) ([]*Profile, error) {
	profiles := make([]*Profile, 0)
	for _, cover := range covers {
		profile := &Profile{
			FileName: cover.fileName,
			Mode:     "set",
			Blocks:   make([]ProfileBlock, 0),
		}
		for i := range cover.counter {
			block := ProfileBlock{
				StartLine: int(cover.pos[3*i+0]),
				StartCol:  int(uint16(cover.pos[3*i+2])),
				EndLine:   int(cover.pos[3*i+1]),
				EndCol:    int(uint16(cover.pos[3*i+2] >> 16)),
				NumStmt:   int(cover.numStmts[i]),
				Count:     int(atomic.LoadUint32(&cover.counter[i])),
			}
			profile.Blocks = append(profile.Blocks, block)
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

func GenerateHtmlReport(out io.Writer) error {
	coversMu.Lock()
	covers := covers_
	coversMu.Unlock()
	profiles, err := makeCoverProfile(covers)
	if nil != err {
		return err
	}
	return htmlOutput(profiles, out)
}
