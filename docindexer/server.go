package docindexer

import (
	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

func ServerRun(args []string) error {
	common.LoadBuiltinPropagationKeys()
	miso.PreServerBootstrap(MakeTempDirs)
	miso.PreServerBootstrap(RegisterRoutes)
	miso.BootstrapServer(args)
	return nil
}
