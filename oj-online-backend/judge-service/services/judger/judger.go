package judger

// #cgo LDFLAGS: -L./core/ -lcore
// #include "./core/src/rules/runner.h"

import "C"

type Compiler struct {
}

func (receiver Compiler) Compile(compileConfig *CompileConfig, srcPath string, outputDir string) {
	C.run()
}
