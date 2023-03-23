package aws

import (
	"github.com/cockroachdb/pebble/cloud/common"
	"github.com/cockroachdb/pebble/vfs"
	"github.com/devlibx/gox-base/errors"
	"io"
	"os"
)

type S3Helper interface {
	SyncFileToS3(file vfs.File, name string) error
	DeleteS3File(name string) error
}

type CloudFS struct {
	wrapperFs vfs.FS
	options   common.CloudFsOption
	s3Helper  S3Helper
}

func NewCloudFS(fs vfs.FS, options common.CloudFsOption) (vfs.FS, error) {

	s3Helper, err := common.NewS3Helper(options)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create S3 helper")
	}

	cfs := &CloudFS{
		wrapperFs: fs,
		options:   options,
		s3Helper:  s3Helper,
	}
	return cfs, nil
}

func (c *CloudFS) Create(name string) (vfs.File, error) {
	if f, err := c.wrapperFs.Create(name); err == nil {
		return NewCloudFile(f, name, c.s3Helper, c.options)
	} else {
		return nil, errors.Wrap(err, "failed to create file: name=%s", name)
	}
}

func (c *CloudFS) Remove(name string) error {
	if err := c.s3Helper.DeleteS3File(name); err == nil {
		return c.wrapperFs.Remove(name)
	} else {
		return errors.Wrap(err, "failed to delete S3 file: name=%s", name)
	}
}

func (c *CloudFS) Rename(oldName, newName string) error {
	if baseFile, err := c.wrapperFs.Create(oldName); err == nil {
		if oldFile, err := NewCloudFile(baseFile, oldName, c.s3Helper, c.options); err == nil {
			if err = c.s3Helper.SyncFileToS3(oldFile, newName); err == nil {
				return c.wrapperFs.Rename(oldName, newName)
			} else {
				return errors.Wrap(err, "failed to sync file for rename to s3: oldName=%s, newName=%s", oldName, newName)
			}
		} else {
			return errors.Wrap(err, "failed to create a cloud file to rename to s3: oldName=%s, newName=%s", oldName, newName)
		}
	} else {
		return errors.Wrap(err, "failed to create a wrapper file to rename to s3: oldName=%s, newName=%s", oldName, newName)
	}
}

func (c *CloudFS) Link(oldname, newname string) error {
	return c.wrapperFs.Link(oldname, newname)
}

func (c *CloudFS) Open(name string, opts ...vfs.OpenOption) (vfs.File, error) {
	return c.wrapperFs.Open(name, opts...)
}

func (c *CloudFS) OpenDir(name string) (vfs.File, error) {
	return c.wrapperFs.OpenDir(name)
}

func (c *CloudFS) RemoveAll(name string) error {
	return c.wrapperFs.RemoveAll(name)
}

func (c *CloudFS) ReuseForWrite(oldname, newname string) (vfs.File, error) {
	return c.wrapperFs.ReuseForWrite(oldname, newname)
}

func (c *CloudFS) MkdirAll(dir string, perm os.FileMode) error {
	return c.wrapperFs.MkdirAll(dir, perm)
}

func (c *CloudFS) Lock(name string) (io.Closer, error) {
	return c.wrapperFs.Lock(name)
}

func (c *CloudFS) List(dir string) ([]string, error) {
	return c.wrapperFs.List(dir)
}

func (c *CloudFS) Stat(name string) (os.FileInfo, error) {
	return c.wrapperFs.Stat(name)
}

func (c *CloudFS) PathBase(path string) string {
	return c.wrapperFs.PathBase(path)
}

func (c *CloudFS) PathJoin(elem ...string) string {
	return c.wrapperFs.PathJoin(elem...)
}

func (c *CloudFS) PathDir(path string) string {
	return c.wrapperFs.PathDir(path)
}

func (c *CloudFS) GetDiskUsage(path string) (vfs.DiskUsage, error) {
	return c.wrapperFs.GetDiskUsage(path)
}
