package main

import (
	// "errors"

	"plugin"

	// "strings"

	"encoding/gob"
	// "errors"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	pak "github.com/septianw/jas/common"
	ty "github.com/septianw/jas/types"
	"github.com/spf13/viper"
)

type Module struct {
	Bootstrap func()
	Router    func(*gin.Engine)
}

func GetConfig(key string) interface{} {
	return viper.Get(key)
}

func GetAllConfig() map[string]interface{} {
	return viper.AllSettings()
}

// This function will load *.so library without parsing its function.
// After load library with this function you need to lookup your function.
func LoadSo(path string) *plugin.Plugin {
	plug, err := plugin.Open(path)
	pak.ErrHandler(err)

	return plug
}

func LoadCoreModule(moduleName string) (*Module, error) {
	var mod Module
	var modpath = filepath.Join(Modloc, "core", moduleName, moduleName+".so")
	_, ferr := os.Stat(modpath)
	if os.IsExist(ferr) {
		lib := LoadSo(modpath)
		bootsym, err := lib.Lookup("Bootstrap")
		pak.ErrHandler(err)

		routersym, err := lib.Lookup("Routers")
		pak.ErrHandler(err)

		mod.Bootstrap = bootsym.(func())
		mod.Router = routersym.(func(*gin.Engine))
		return &mod, nil
	}

	return nil, ferr
}

func LoadContribModule(moduleName string) (*Module, error) {
	var mod Module
	var modpath = filepath.Join(Modloc, "contrib", moduleName, moduleName+".so")
	_, ferr := os.Stat(modpath)
	if os.IsExist(ferr) {
		lib := LoadSo(modpath)
		bootsym, err := lib.Lookup("Bootstrap")
		pak.ErrHandler(err)

		routersym, err := lib.Lookup("Routers")
		pak.ErrHandler(err)

		mod.Bootstrap = bootsym.(func())
		mod.Router = routersym.(func(*gin.Engine))
		return &mod, nil
	}

	return nil, ferr
}

func WriteRuntime(rt ty.Runtime) {
	RuntimeFile, err := os.OpenFile("/tmp/shinyRuntimeFile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	pak.ErrHandler(err)

	enc := gob.NewEncoder(RuntimeFile)
	err = enc.Encode(rt)
	pak.ErrHandler(err)

	err = RuntimeFile.Close()
	pak.ErrHandler(err)
}

func LoadDatabase(libpath string, d ty.Dbconf) ty.Database {
	// rt := pak.ReadRuntime()

	// pak.ErrHandler(errors.New(rt.Libloc))
	// pak.ErrHandler(errors.New(filepath.Join(rt.Libloc, "database.so")))

	plug := LoadSo(libpath)
	symd, err := plug.Lookup("Database")
	pak.ErrHandler(err)
	sd := symd.(ty.Database)

	return sd
}
