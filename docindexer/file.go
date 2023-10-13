package docindexer

import (
	"os"

	"github.com/curtisnewbie/miso/miso"
)

const (
	PropTempPath    = "docindexer.temp-path"
	DefaultTempPath = "/tmp/docindexer"
)

func init() {
	miso.SetDefProp(PropTempPath, DefaultTempPath)
}

func MakeTempDirs(rail miso.Rail) error {
	dir := miso.GetPropStr(PropTempPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	rail.Infof("MkdirAll %v finished", dir)
	return nil
}

func TempFilePath(tempTkn string) string {
	return miso.GetPropStr(PropTempPath + "/" + tempTkn)
}
