package hop

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"go.hop.io/sdk/types"
)

type dirItem struct {
	Name        string          `json:"name"`
	Directory   bool            `json:"directory"`
	Permissions int             `json:"permissions"`
	UpdatedAt   types.Timestamp `json:"updated_at"`
	Size        int64           `json:"size"`
}

type volumeFile struct {
	// Defines the read closer.
	rcLock sync.Mutex
	rc     io.ReadCloser

	// Defines if this is a file.
	isFile bool

	// Defines the path and client information.
	path string
	c    clientDoer
	ctx  context.Context
	opts []ClientOption

	// Defines the file information.
	filename   string
	dirListing []*dirItem
}

func (f *volumeFile) getReadStream() error {
	var resp *http.Response
	err := f.c.do(f.ctx, ClientArgs{
		Method:      "GET",
		Path:        f.path,
		Query:       map[string]string{"stream": "true"},
		Ignore404:   false,
		PassRequest: func(r *http.Response) { resp = r },
	}, f.opts)
	if err != nil {
		return err
	}
	f.rc = resp.Body
	return nil
}

func (f *volumeFile) Read(b []byte) (int, error) {
	f.rcLock.Lock()
	defer f.rcLock.Unlock()

	if f.rc != nil {
		// Reader already here meaning checks have been passed.
		return f.rc.Read(b)
	}

	if !f.isFile {
		// This is unintuitive, but this is the right error to emulate a bad path.
		return 0, &fs.PathError{
			Op:   "read",
			Path: f.filename,

			// this feels *close enough* since libraries probably should not be relying
			// on OS specific error codes and this is the error content you'd expect on
			// anything POSIX based.
			Err: errors.New("is a directory"),
		}
	}

	err := f.getReadStream()
	var n int
	if err == nil {
		n, err = f.rc.Read(b)
	}
	return n, err
}

func (f *volumeFile) Close() error {
	f.rcLock.Lock()
	defer f.rcLock.Unlock()

	if f.rc != nil {
		return f.rc.Close()
	}
	return nil
}

func (f *volumeFile) Mode() fs.FileMode {
	if f.isFile {
		return fs.FileMode(f.dirListing[0].Permissions)
	}
	return 0o777
}

func (f *volumeFile) Stat() (fs.FileInfo, error) {
	// Return ourselves since we implement both.
	return f, nil
}

func (f *volumeFile) Name() string {
	_, file := path.Split(f.filename)
	return file
}

func (f *volumeFile) Size() int64 {
	if f.isFile {
		return f.dirListing[0].Size
	}

	var size int64
	for _, v := range f.dirListing {
		size += v.Size
	}
	return size
}

func (f *volumeFile) ModTime() time.Time {
	if f.isFile {
		t, _ := f.dirListing[0].UpdatedAt.Time()
		return t
	}

	// Find the latest time.
	var latestTime time.Time
	for _, v := range f.dirListing {
		t, _ := v.UpdatedAt.Time()
		if t.After(latestTime) {
			latestTime = t
		}
	}
	return latestTime
}

func (f *volumeFile) IsDir() bool {
	return !f.isFile
}

func (f *volumeFile) Sys() any {
	// It is within spec to not implement this, and due to the emulation
	// nature of this, we do not.
	return nil
}

var (
	_ fs.File     = (*volumeFile)(nil)
	_ fs.FileInfo = (*volumeFile)(nil)
)

type volumeFs struct {
	c            clientDoer
	ctx          context.Context
	deploymentId string
	volumeId     string
	opts         []ClientOption
}

const fileNotFoundCode = "file_not_found"

func (f volumeFs) Open(name string) (fs.File, error) {
	// Run fs package validations.
	v := fs.ValidPath(name)
	if !v {
		return nil, fs.ErrInvalid
	}

	// Get the path.
	urlChunk := ""
	if name != "" && name != "." {
		urlChunk = "/" + url.PathEscape(name)
	}
	urlChunk = "/ignite/deployments/" +
		url.PathEscape(f.deploymentId) + "/volumes/" +
		url.PathEscape(f.volumeId) + "/files" + urlChunk

	// Make the network request.
	var resp struct {
		Folder bool            `json:"folder"`
		File   json.RawMessage `json:"file"`
	}
	err := f.c.do(f.ctx, ClientArgs{
		Method:    "GET",
		Path:      urlChunk,
		Result:    &resp,
		Ignore404: false,
	}, f.opts)

	// Handle any errors.
	if err != nil {
		notFound, ok := err.(types.NotFound)
		if !ok || notFound.Code != fileNotFoundCode {
			// A 404 somewhere else in the chain.
			return nil, err
		}

		// Throw a not found error.
		return nil, fs.ErrNotExist
	}

	// Get the files.
	var files []*dirItem
	if resp.Folder {
		if err = json.Unmarshal(resp.File, &files); err != nil {
			return nil, err
		}
	} else {
		var file dirItem
		if err = json.Unmarshal(resp.File, &file); err != nil {
			return nil, err
		}
		files = []*dirItem{&file}
	}

	// Treat this as an actual user file.
	return &volumeFile{
		isFile:     !resp.Folder,
		path:       urlChunk,
		c:          f.c,
		ctx:        f.ctx,
		opts:       f.opts,
		filename:   name,
		dirListing: files,
	}, nil
}

var _ fs.FS = volumeFs{}

// VolumeVirtualFS is used to make a fs.FS compatible virtual filesystem.
// Note that the context should live as long as the filesystem in this instance.
func (c ClientCategoryIgniteDeployments) VolumeVirtualFS(
	ctx context.Context, deploymentId, volumeId string, opts ...ClientOption,
) fs.FS {
	return volumeFs{
		c:            c.c,
		ctx:          ctx,
		deploymentId: deploymentId,
		volumeId:     volumeId,
		opts:         opts,
	}
}
