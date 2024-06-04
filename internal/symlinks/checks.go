package symlinks

import (
	"os"
	"strings"
)

func SymlinkExists(destination string) bool {
	if _, err := os.Readlink(destination); err == nil {
		return true
	}
	return false
}

func PathMatchesSymlink(rootName string, nodeName string, fileSource string, flagPath string) bool {
	ps := string(os.PathSeparator)
	expSource := os.ExpandEnv(rootName + ps + nodeName + ps + fileSource)
	return strings.HasPrefix(expSource, os.ExpandEnv(flagPath))
}

func SymlinkPointsToSource(link string, source string) bool {
	if l, err := os.Readlink(link); err == nil {
		return l == source
	}
	return false
}