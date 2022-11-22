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

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dtstack.com/dtstack/easymatrix/matrix/util"
	"dtstack.com/dtstack/easymatrix/schema"
)

var (
	pkgDir, oldPkg string
	all            bool
)

func init() {
	//flag.StringVar(&pkgDir, "pkg-dir", "", "product package directory")
	//flag.StringVar(&oldPkg, "old-pkg", "", "old product package tar file")
	flag.BoolVar(&all, "all", false, "generate the complete package")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-all] pkg-dir [old-pkg]\n", os.Args[0])
	}
}

func verifyPackageDir(pkgDir string) (sc *schema.SchemaConfig, err error) {
	sc, err = schema.ParseSchemaConfigFile(filepath.Join(pkgDir, schema.SCHEMA_FILE))
	if err != nil {
		return
	}
	if err = sc.ParseVariable(); err != nil {
		return
	}

	baseDir, err := ioutil.TempDir("", "mero")
	if err != nil {
		return
	}
	defer os.RemoveAll(baseDir)

	for name, svc := range sc.Service {
		// set fake ip just for check ConfigFiles
		sc.SetServiceAddr(name, []string{"127.0.0.1"}, []string{"localhost"})

		if svc.Instance != nil {
			for _, configPath := range svc.Instance.ConfigPaths {
				configPath = filepath.Join(name, configPath)
				absConfigFile := filepath.Join(baseDir, configPath)
				if err = os.MkdirAll(filepath.Dir(absConfigFile), 0755); err != nil {
					return
				}
				if _, err = util.CopyFile(filepath.Join(pkgDir, configPath), absConfigFile); err != nil {
					return
				}
			}
		}
	}
	if err = sc.ParseVariable(); err != nil {
		return
	}
	// check all ConfigFiles once
	_, err = sc.ParseConfigFiles(baseDir)
	return
}

func createZip(dir string) (zipH *os.File, err error) {
	if zipH, err = ioutil.TempFile("", "mero_zip"); err != nil {
		return
	}
	zw := zip.NewWriter(zipH)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if dir == path {
			return nil
		}
		if baseName := filepath.Base(path); baseName[0] == '.' {
			//skip hidden file
			return nil
		}

		hdr, _ := zip.FileInfoHeader(info)
		relName, _ := filepath.Rel(dir, path)
		hdr.Name = filepath.ToSlash(relName)
		if info.IsDir() {
			hdr.Name += "/"
		}
		hdr.Method = zip.Deflate
		zf, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if info.Mode()&os.ModeSymlink == 0 {
				var f *os.File
				if f, err = os.Open(path); err != nil {
					return err
				}
				_, err = io.Copy(zf, f)
				f.Close()
			} else {
				var l string
				if l, err = os.Readlink(path); err == nil {
					_, err = zf.Write([]byte(l))
				}
			}
		}
		return err
	})
	if err != nil {
		zw.Close()
		zipH.Close()
		os.Remove(zipH.Name())
		return
	}

	if err = zw.Close(); err != nil {
		zipH.Close()
		os.Remove(zipH.Name())
		return
	}
	zipH.Seek(0, io.SeekStart)

	return
}

func createPackage(pkgDir string, sc *schema.SchemaConfig) (err error) {
	pkgName := sc.ProductName + "_" + sc.ProductVersion + ".tar"
	pkgFile, err := os.Create(pkgName)
	if err != nil {
		return
	}
	tw := tar.NewWriter(pkgFile)
	defer func() {
		if err == nil {
			err = tw.Close()
		}
		pkgFile.Close()
		if err != nil {
			os.Remove(pkgName)
		} else {
			fmt.Printf("Product package create success: %q\n", pkgName)
		}
	}()

	if err = insertToTar(pkgDir, schema.SCHEMA_FILE, tw); err != nil {
		return
	}

	for name, svc := range sc.Service {
		if svc.Instance == nil {
			continue
		}

		var f *os.File
		var info os.FileInfo
		if f, err = createZip(filepath.Join(pkgDir, name)); err != nil {
			return
		}
		zipFile := f.Name()
		if info, err = f.Stat(); err != nil {
			f.Close()
			os.Remove(zipFile)
			return
		}
		if err = tw.WriteHeader(&tar.Header{Name: name + ".zip", Mode: 0600, Size: info.Size()}); err != nil {
			f.Close()
			os.Remove(zipFile)
			return
		}
		_, err = io.Copy(tw, f)
		f.Close()
		os.Remove(zipFile)
		if err != nil {
			return
		}
	}

	return
}

func createPackagePatch(pkgDir, oldPkg string, sc *schema.SchemaConfig) (err error) {
	var oldPkgFile *os.File
	if oldPkgFile, err = os.Open(oldPkg); err != nil {
		return
	}
	defer oldPkgFile.Close()

	var oldsc *schema.SchemaConfig
	patch := &schema.Patch{}
	tr := tar.NewReader(oldPkgFile)
	for {
		var hdr *tar.Header
		hdr, err = tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if oldsc == nil {
			if hdr.Name != schema.SCHEMA_FILE {
				return fmt.Errorf("old-pkg package first file is %v, not %v", hdr.Name, schema.SCHEMA_FILE)
			}
			buf := bytes.NewBuffer(make([]byte, 0, hdr.Size))
			if _, err = io.Copy(buf, tr); err != nil {
				return err
			}
			if oldsc, err = schema.ParseSchemaConfigBytes(buf.Bytes()); err != nil {
				return err
			}
			if sc.ProductName != oldsc.ProductName {
				return fmt.Errorf("`%v`-->`%v` product name not same", oldsc.ProductName, sc.ProductName)
			}
			if sc.ProductVersion == oldsc.ProductVersion {
				return fmt.Errorf("product `%v` have same version %v, not need patch", sc.ProductName, sc.ProductVersion)
			}
			patch.ProductName = sc.ProductName
			patch.NewProductVersion = sc.ProductVersion
			patch.OldProductVersion = oldsc.ProductVersion

			fmt.Printf("Product: `%v`, upgrade version %v to %v\n", patch.ProductName, patch.OldProductVersion, patch.NewProductVersion)
			continue
		}
		oldServiceName := strings.TrimSuffix(hdr.Name, ".zip")
		oldsvc, exist := oldsc.Service[oldServiceName]
		if !exist || oldsvc.Instance == nil {
			fmt.Printf("service %v.zip in old-pkg is unnecessary", oldServiceName)
			continue
		}
		newsvc, exist := sc.Service[oldServiceName]
		if !exist {
			patch.DeletedServices = append(patch.DeletedServices, oldServiceName)

			fmt.Printf("Delete Service: %v\n", oldServiceName)
			continue
		}
		if newsvc.Version != oldsvc.Version {
			fmt.Printf("Service: %v, upgrade version %v to %v\n", oldServiceName, oldsvc.Version, newsvc.Version)
			if newsvc.Instance != nil {
				var zipH *os.File
				if zipH, err = ioutil.TempFile("", "mero_zip"); err != nil {
					return err
				}
				if _, err = io.Copy(zipH, tr); err == nil {
					zipH.Seek(0, io.SeekStart)
					diff, err := getDiffServiceFiles(oldServiceName, pkgDir, zipH)
					if err != nil {
						zipH.Close()
						os.Remove(zipH.Name())
						return err
					}
					if diff != nil {
						patch.DiffServices = append(patch.DiffServices, diff)
					}
				}
				zipH.Close()
				os.Remove(zipH.Name())
			} else {
				patch.DeletedServices = append(patch.DeletedServices, oldServiceName)
				fmt.Printf("Delete Virtual Service: %v\n", oldServiceName)
			}
		}
		delete(sc.Service, oldServiceName)
	}
	if oldsc == nil {
		return fmt.Errorf("can't get %v in old-pkg package", schema.SCHEMA_FILE)
	}
	for newServiceName, newsvc := range sc.Service {
		if newsvc.Instance != nil {
			patch.NewServices = append(patch.NewServices, newServiceName)
			fmt.Printf("New Service: %v\n", newServiceName)
		} else {
			// virtual service recheck because old-pkg not have zip
			if oldsvc, exist := oldsc.Service[newServiceName]; !exist {
				fmt.Printf("New Virtual Service: %v\n", newServiceName)
			} else if newsvc.Version != oldsvc.Version {
				fmt.Printf("Service: %v, upgrade version %v to %v\n", newServiceName, oldsvc.Version, newsvc.Version)
			}
		}
	}

	return generatePatch(patch, pkgDir)
}

func generatePatch(patch *schema.Patch, pkgDir string) (err error) {
	var patchBuf bytes.Buffer
	if err = gob.NewEncoder(&patchBuf).Encode(patch); err != nil {
		return
	}

	patchName := patch.ProductName + "_" + patch.OldProductVersion + "_" + patch.NewProductVersion + ".patch"
	pf, err := os.Create(patchName)
	if err != nil {
		return
	}

	tw := tar.NewWriter(pf)
	defer func() {
		if err == nil {
			err = tw.Close()
		}
		pf.Close()
		if err != nil {
			os.Remove(patchName)
		} else {
			fmt.Printf("Product patch create success: %q\n", patchName)
		}
	}()

	if err = insertToTar(pkgDir, schema.SCHEMA_FILE, tw); err != nil {
		return
	}

	if err = tw.WriteHeader(&tar.Header{Name: schema.PATCH_FILE, Mode: 0600, Size: int64(patchBuf.Len())}); err != nil {
		return
	}
	if _, err = tw.Write(patchBuf.Bytes()); err != nil {
		return
	}

	for _, newServiceName := range patch.NewServices {
		var f *os.File
		var info os.FileInfo
		if f, err = createZip(filepath.Join(pkgDir, newServiceName)); err != nil {
			return
		}
		zipFile := f.Name()
		if info, err = f.Stat(); err != nil {
			f.Close()
			os.Remove(zipFile)
			return
		}
		if err = tw.WriteHeader(&tar.Header{Name: newServiceName + ".zip", Mode: 0600, Size: info.Size()}); err != nil {
			f.Close()
			os.Remove(zipFile)
			return
		}
		_, err = io.Copy(tw, f)
		f.Close()
		os.Remove(zipFile)
		if err != nil {
			return
		}
	}

	for _, diff := range patch.DiffServices {
		for _, file := range diff.NewFiles {
			if err = insertToTar(pkgDir, filepath.Join(diff.ServiceName, file), tw); err != nil {
				return
			}
		}
		for _, file := range diff.DiffFiles {
			if err = insertToTar(pkgDir, filepath.Join(diff.ServiceName, file), tw); err != nil {
				return
			}
		}
	}

	return
}

func insertToTar(pkgDir, filename string, tw *tar.Writer) error {
	f, err := os.Open(filepath.Join(pkgDir, filename))
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.IsDir() && filepath.Dir(filename) != "." {
		// two layer file need add parent dir
		dirInfo, err := os.Lstat(filepath.Dir(f.Name()))
		if err != nil {
			return err
		}
		dirHdr, err := tar.FileInfoHeader(dirInfo, "")
		if err != nil {
			return err
		}
		dirHdr.Name = filepath.Dir(filename)
		if err = tw.WriteHeader(dirHdr); err != nil {
			return err
		}
	}
	hdr, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	hdr.Name = filename
	if err = tw.WriteHeader(hdr); err != nil {
		return err
	}
	if !info.IsDir() {
		if info.Mode()&os.ModeSymlink == 0 {
			_, err = io.Copy(tw, f)
		} else {
			var l string
			if l, err = os.Readlink(filepath.Join(pkgDir, filename)); err == nil {
				_, err = tw.Write([]byte(l))
			}
		}
	}

	return err
}

func getDiffServiceFiles(serviceName, pkgDir string, zipH *os.File) (*schema.DiffService, error) {
	zipInfo, err := zipH.Stat()
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(zipH, zipInfo.Size())
	if err != nil {
		return nil, err
	}

	allFile, err := walkFileMap(filepath.Join(pkgDir, serviceName))
	if err != nil {
		return nil, err
	}

	diff := &schema.DiffService{ServiceName: serviceName}
	for _, file := range zr.File {
		filePath := filepath.Join(pkgDir, serviceName, file.Name)
		fileInfo, exist := allFile[file.Name]
		if !exist {
			diff.DeletedFiles = append(diff.DeletedFiles, file.Name)

			fmt.Printf("........ delete file: %v\n", file.Name)
			continue
		}
		delete(allFile, file.Name)

		if file.FileInfo().IsDir() {
			continue
		}

		if file.FileInfo().Mode()&os.ModeSymlink != 0 {
			lnew, err := os.Readlink(filePath)
			if err != nil {
				return nil, err
			}
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			lold, err := ioutil.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, err
			}
			if lnew != string(lold) {
				diff.DiffFiles = append(diff.DiffFiles, file.Name)

				fmt.Printf("........ diff file(link): %v\n", file.Name)
			}
			continue
		}

		if fileInfo.size != file.FileInfo().Size() {
			diff.DiffFiles = append(diff.DiffFiles, file.Name)

			fmt.Printf("........ diff file: %v\n", file.Name)
			continue
		}

		if !fileInfo.modeTime.Equal(file.Modified) {
			fileH, err := os.Open(filePath)
			if err != nil {
				return nil, err
			}
			h32 := crc32.NewIEEE()
			io.Copy(h32, fileH)
			if h32.Sum32() != file.CRC32 {
				diff.DiffFiles = append(diff.DiffFiles, file.Name)

				fmt.Printf("........ diff file: %v\n", file.Name)
			}
			fileH.Close()
		}
	}
	for newFile := range allFile {
		diff.NewFiles = append(diff.NewFiles, newFile)

		fmt.Printf("........ new file: %v\n", newFile)
	}

	return diff, nil
}

type fileInfo struct {
	size     int64
	modeTime time.Time
}

func walkFileMap(dir string) (map[string]fileInfo, error) {
	fileMap := make(map[string]fileInfo)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if dir == path {
			return nil
		}

		if baseName := filepath.Base(path); baseName[0] == '.' {
			//skip hidden file
			return nil
		}
		relName, _ := filepath.Rel(dir, path)
		if info.IsDir() {
			relName += "/"
		}
		fileMap[filepath.ToSlash(relName)] = fileInfo{info.Size(), info.ModTime()}
		return nil
	})

	return fileMap, err
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		flag.Usage()
		os.Exit(1)
	}
	pkgDir = args[0]
	if len(args) == 2 {
		oldPkg = args[1]
	}
	if pkgDir == "" {
		fmt.Fprintln(os.Stderr, "pkg-dir is empty!")
		flag.Usage()
		os.Exit(1)
	}

	var err error
	var sc *schema.SchemaConfig
	if sc, err = verifyPackageDir(pkgDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if risks := sc.CheckRisks(pkgDir); len(risks) > 0 {
		for name, list := range risks {
			fmt.Printf("[告警] 服务%s包含危险脚本: %v\n", name, list)
		}
	}

	if oldPkg != "" {
		err = createPackagePatch(pkgDir, oldPkg, sc)
	}
	if oldPkg == "" || (err == nil && all) {
		err = createPackage(pkgDir, sc)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
