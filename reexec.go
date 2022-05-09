package reexec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	registeredInitializers = make(map[string]func())
	// exec command default runtime path env name
	execRuntimePathEnvName = "EXEC_DEFAULT_RUNTIME_PATH"
)

// SetExecRuntimePathEnvName change exec command default runtime path
func SetExecRuntimePathEnvName(name string) {
	execRuntimePathEnvName = name
}

// Register adds an initialization func under the specified name
func Register(name string, initializer func()) {
	if _, exists := registeredInitializers[name]; exists {
		panic(fmt.Sprintf("reexec func already registered under name %q", name))
	}

	registeredInitializers[name] = initializer
}

// Init is called as the first part of the exec process and returns true if an
// initialization function was called.
func Init() bool {
	if envValue := os.Getenv(execRuntimePathEnvName); envValue == "" {
		_ = os.Setenv(execRuntimePathEnvName, naiveSelf())
	}
	initializer, exists := registeredInitializers[os.Args[0]]
	if exists {
		initializer()

		return true
	}
	return false
}

func naiveSelf() string {
	// use global exec runtime path
	if envValue := os.Getenv(execRuntimePathEnvName); envValue != "" {
		return envValue
	}
	name := os.Args[0]
	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err == nil {
			return lp
		}
	}
	// handle conversion of relative paths to absolute
	if absName, err := filepath.Abs(name); err == nil {
		return absName
	}
	// if we couldn't get absolute name, return original
	// (NOTE: Go only errors on Abs() if os.Getwd fails)
	return name
}
