package path

import (
	"os"
	"path/filepath"
)

//获取目录路径
func GetDirPath(dirName string) (appDirPath string) {
	var err error
	var appPath string
	if appPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appDirPath = filepath.Join(workPath, dirName)

	if !FileExists(appDirPath) {
		appDirPath = filepath.Join(appPath, dirName)
	}
	return
}

//获取文件路径
func GetFilePath(dirName string, fileName string) (appFilePath string) {
	var err error
	var appPath string
	if appPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		panic(err)
	}
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	appFilePath = filepath.Join(workPath, dirName, fileName)

	if !FileExists(appFilePath) {
		appFilePath = filepath.Join(appPath, dirName, fileName)
	}
	return
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
