package aws

import (
	"github.com/cockroachdb/pebble/vfs"
	"io"
	"os"
)

type CloudFS struct {
	wrapperFs vfs.FS
}

func (c *CloudFS) Create(name string) (vfs.File, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) Link(oldname, newname string) error {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) Open(name string, opts ...vfs.OpenOption) (vfs.File, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) OpenDir(name string) (vfs.File, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) Remove(name string) error {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) RemoveAll(name string) error {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) Rename(oldname, newname string) error {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) ReuseForWrite(oldname, newname string) (vfs.File, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) MkdirAll(dir string, perm os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) Lock(name string) (io.Closer, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) List(dir string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) Stat(name string) (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) PathBase(path string) string {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) PathJoin(elem ...string) string {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) PathDir(path string) string {
	//TODO implement me
	panic("implement me")
}

func (c *CloudFS) GetDiskUsage(path string) (vfs.DiskUsage, error) {
	//TODO implement me
	panic("implement me")
}

func NewCloudFS(fs vfs.FS) vfs.FS {
	cfs := &CloudFS{
		wrapperFs: fs,
	}
	return cfs
}
