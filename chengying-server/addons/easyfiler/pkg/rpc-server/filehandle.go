// Licensed to Apache Software Foundation(ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation(ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package rpc_server

import (
	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/filemeta"
	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/limit"
	"dtstack.com/dtstack/easymatrix/addons/easyfiler/pkg/proto"
	"dtstack.com/dtstack/easymatrix/go-common/log"
	"fmt"
	"golang.org/x/time/rate"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	_LINES_TO_PEEK          = 300
	_INTERVAL               = 20
	_SIZE_KB                = 1024
	_TAR_SUFFIX             = ".tar"
	_TOKEN_BUCKET_CAPATCITY = 1
	_LIMIT                  = 1
)

type FileTransferService struct {
	Root    string
	WithDB  bool
	Rate    int
	IsZiped bool
}

type localStorage struct {
	lock            sync.RWMutex
	lastPreviewTime time.Time
	path            string
	lastAtLine      int
	totalLines      int
}

var local = new(localStorage)

func (fts *FileTransferService) Download(r *proto.DownloadRequest, stream proto.TransferService_DownloadServer) error {
	dst := fts.Root + r.Path
	if fts.IsZiped {
		dir, name := spiteDirAndFilename(dst)
		if err := tar(name, dir); err != nil {
			return err
		}
		dst += _TAR_SUFFIX
		defer os.Remove(dst)
	}
	log.Infof("download file %v", dst)
	file, err := os.Open(dst)
	if err != nil {
		return err
	}

	defer file.Close()
	_rate := fts.Rate * _SIZE_KB
	buf := make([]byte, _rate)
	limiter := rate.NewLimiter(_LIMIT, _TOKEN_BUCKET_CAPATCITY)
	limitedReader := limit.NewLimitedReader(file, limiter)
	for {
		n, err := limitedReader.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			stream.Send(&proto.DownloadResponse{Data: buf[:n]})
			break
		}
		if err = stream.Send(&proto.DownloadResponse{
			Data: buf[:n],
		}); err != nil {
			return err
		}
	}
	return nil
}

func (fts *FileTransferService) Preview(r *proto.PreviewRequest, stream proto.TransferService_PreviewServer) (err error) {
	if local.path != r.Path {
		local = &localStorage{}
	}
	if isPathExist(fts.Root+r.Path) == false {
		stream.Send(&proto.PreviewResponse{Validity: "attempt to preview a non-existent file"})
		return fmt.Errorf("attempt to preview a non-existent file")
	}
	log.Infof("%#v\n", local)
	go func() {
		lines, _ := fileLines(fts.Root + r.Path)
		lastatline := lines
		if local.lastAtLine != 0 {
			switch r.Action {
			case "up":
				if local.lastAtLine <= _LINES_TO_PEEK {
					lastatline = local.lastAtLine
				} else {
					lastatline = local.lastAtLine - _INTERVAL
				}
			case "down":
				if local.lastAtLine+_INTERVAL > lines {
					lastatline = lines
				} else {
					lastatline = local.lastAtLine + _INTERVAL
				}
			case "latest":
				lastatline = lines
			}
		}
		local.lock.Lock()
		local = &localStorage{
			path:            r.Path,
			lastPreviewTime: time.Now(),
			lastAtLine:      lastatline,
			totalLines:      lines,
		}
		local.lock.Unlock()
		return
	}()

	out := make([]byte, 1)
	switch r.Action {
	case "up", "down":
		out, err = peekLines(_LINES_TO_PEEK, local.lastAtLine, fts.Root+r.Path)
		if err != nil {
			return fmt.Errorf("failed to peek up")
		}
	case "latest":
		out, err = peekLines(_LINES_TO_PEEK, local.totalLines, fts.Root+r.Path)
		if err != nil {
			return err
		}
	default:
		stream.Send(&proto.PreviewResponse{Validity: "invalid action param"})
		return fmt.Errorf("invalid action param")
	}

	logRows := strings.Split(string(out), "\n")

	for i := range logRows {
		if err := stream.Send(&proto.PreviewResponse{
			LogRow: logRows[i],
		}); err != nil {
			return err
		}
	}
	return err
}

func (fts *FileTransferService) List(r *proto.ListRequest, stream proto.TransferService_ListServer) error {
	err := filepath.Walk(fts.Root, func(p string, info os.FileInfo, err error) error {
		name, err := filepath.Rel(fts.Root, p)
		if err != nil {
			return err
		}
		name = filepath.ToSlash(name)
		f := &proto.ListResponse{
			Name: filepath.ToSlash(name),
		}
		return stream.Send(f)
	})
	return err
}

func (fts *FileTransferService) Upload(stream proto.TransferService_UploadServer) error {
	r, err := stream.Recv()
	if err != nil {
		fmt.Printf("failed to recive file, err:%v\n", err)
		return err
	}
	fileMeta := filemeta.FileMeta{
		Location: fts.Root + r.FileName,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("failed to create file, err:%v\n", err)
		return err
	}
	defer newFile.Close()

	size, _ := newFile.Write(r.Data)
	if fts.WithDB {
		fileMeta.FileSize = int64(size)
		newFile.Seek(0, 0)
		ok := filemeta.UpdateFileMeta(fileMeta)
		if !ok {
			fmt.Println("already exist")
			return fmt.Errorf("failed to upload")
		}
	}
	stream.SendAndClose(&proto.UploadResponse{
		Finished: true,
	})
	return nil
}

func sed(start, end int, path string) ([]byte, error) {
	cmdStr := `sed -n '%d,%dp' %s`
	cmd := exec.Command("sh", "-c", fmt.Sprintf(cmdStr, start, end, path))
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func peekLines(linesToPeek, end int, path string) ([]byte, error) {
	cmdStr := `head -n %d %s | tail -n %d`
	cmd := exec.Command("sh", "-c", fmt.Sprintf(cmdStr, end, path, linesToPeek))
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, err
}

func fileLines(path string) (int, error) {
	cmdStr := `awk 'END{print NR}' %s | awk '{print $1}'`
	//cmdStr := `wc -l %s | awk '{print $1}'`
	cmd := exec.Command("sh", "-c", fmt.Sprintf(cmdStr, path))
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	lines, _ := strconv.Atoi(strings.Split(string(out), "\n")[0])
	return lines, nil
}

func isPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

func spiteDirAndFilename(path string) (dir, filename string) {
	pathSlice := strings.SplitAfter(path, "/")
	for i := 0; i < len(pathSlice)-1; i++ {
		dir += pathSlice[i]
	}
	filename = pathSlice[len(pathSlice)-1]
	return
}

func tar(name, dir string) error {
	cmdStr := `tar -czf %s%s %s`
	fmt.Println(fmt.Sprintf(cmdStr, name, _TAR_SUFFIX, name))
	cmd := exec.Command("sh", "-c", fmt.Sprintf(cmdStr, name, _TAR_SUFFIX, name))
	cmd.Dir = dir
	return cmd.Run()
}
