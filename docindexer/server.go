package docindexer

import (
	"github.com/curtisnewbie/miso/middleware/user-vault/common"
	"github.com/curtisnewbie/miso/miso"
)

func ServerRun(args []string) error {
	common.LoadBuiltinPropagationKeys()
	miso.PreServerBootstrap(MakeTempDirs)
	miso.PreServerBootstrap(RegisterRoutes)
	miso.BootstrapServer(args)
	return nil
}
