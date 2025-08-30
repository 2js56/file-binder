package loader

import (
	"embed"
	"path"
)

//go:embed "module"
var moduleFolder embed.FS

// Modules 存储模板名称到模板内容的映射
var Modules = make(map[string][]byte)

// LoadModules 加载嵌入的模块模板
func LoadModules() error {
	entries, err := moduleFolder.ReadDir("module")
	if err != nil {
		return err
	}

	// 预分配足够的容量
	if len(Modules) == 0 {
		Modules = make(map[string][]byte, len(entries))
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		entryInfo, err := entry.Info()
		if err != nil {
			continue
		}

		mainPath := path.Join("module", entryInfo.Name(), "main.go")
		loaderFileContent, err := moduleFolder.ReadFile(mainPath)
		if err != nil {
			continue
		}

		Modules[entryInfo.Name()] = loaderFileContent
	}

	return nil
}

// 通过embed将模块的loader装载进程序，不再依赖本地文件
func init() {
	if err := LoadModules(); err != nil {
		panic("Failed to load modules: " + err.Error())
	}
}
