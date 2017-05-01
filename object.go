package sysfs

import (
	// "os"

	"path/filepath"
	"strings"
)

type Object string

func objFullPath(path string) Object {
	path, err := filepath.EvalSymlinks(path)
	if err != nil {
		return Object(path)
	}
	return Object(path)
}

func (obj Object) Exists() bool {
	d, _ := dirExists(string(obj))
	return d
}

func (obj Object) Name() string {
	return string(obj)[strings.LastIndex(string(obj), "/")+1:]
}

func (obj Object) SubObjects() []Object {
	path := string(obj) + "/"
	objects := make([]Object, 0)
	lsDirs(path, func(name string) {
		objects = append(objects, objFullPath(path+name))
	})
	return objects
}

func (obj Object) SubObjectsFilter(filter string) []Object {
	path := string(obj) + "/"
	objects := make([]Object, 0)

	lsDirs(path, func(name string) {
		match, err := filepath.Match(filter, name)
		if match && err == nil {
			objects = append(objects, objFullPath(path+name))
		}
	})
	return objects
}

func (obj Object) SubObject(name string) (Object, error) {
	d, err := dirExists(string(obj) + "/" + name)
	if d {
		return objFullPath(string(obj) + "/" + name), nil
	}
	return "", err

}

func (obj Object) Attributes() []Attribute {
	path := string(obj) + "/"
	attribs := make([]Attribute, 0)
	lsFiles(path, func(name string) {
		attribs = append(attribs, Attribute{Path: path + name})
	})
	return attribs
}

func (obj Object) Attribute(name string) *Attribute {
	return &Attribute{Path: string(obj) + "/" + name}
}

func (obj Object) Parent(count int) Object {
	if count < 0 {
		p := strings.Split(string(obj), "/")
		return objFullPath(string(obj) + strings.Repeat("/..", len(p)+count))
	}
	return objFullPath(string(obj) + strings.Repeat("/..", count))

}
