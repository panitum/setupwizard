package custom

import (
	"os"
	"path/filepath"
)

var root string = "custom"
var scriptExt string = ".sh"

type Custom struct {
	scripts []string
}

func NewCustom() (*Custom, error) {
	c := &Custom{}

	if err := ensureDir(); err != nil {
		return nil, err
	}

	if err := c.collectScripts(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Custom) Count() int {
	return len(c.scripts)
}

func (c *Custom) GetScripts() []string {
	return c.scripts
}

func (c *Custom) collectScripts() error {
	dir, err := os.Open(root)
	if err != nil {
		return err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == scriptExt {
			c.scripts = append(c.scripts, filepath.Join(root, file.Name()))
		}
	}

	return nil
}

func Init() error {
	return ensureDir()
}

func ensureDir() error {
	dirPath := filepath.Join(root)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	return nil
}
