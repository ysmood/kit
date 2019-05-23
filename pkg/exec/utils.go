package exec

import (
	os_path "path"

	. "github.com/ysmood/gokit/pkg/os"
	. "github.com/ysmood/gokit/pkg/utils"
)

func EnsureGoTool(path string) {
	if !Exists(GoPath() + "/bin/" + os_path.Base(path)) {
		E(Exec("go", "get", path).Dir(HomeDir()).Do())
	}
}
