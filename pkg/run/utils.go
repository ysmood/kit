package run

import (
	os_path "path"

	"github.com/ysmood/gokit/pkg/os"
	"github.com/ysmood/gokit/pkg/utils"
)

// MustGoTool try to find executable bin under GOPATH, if not exists run go get to download it
func MustGoTool(path string) {
	if !os.Exists(os.GoPath() + "/bin/" + os_path.Base(path)) {
		utils.E(Exec("go", "get", path).Dir(os.HomeDir()).Do())
	}
}
