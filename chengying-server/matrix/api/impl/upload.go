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

package impl

import (
	sysContext "context"
	apibase "dtstack.com/dtstack/easymatrix/go-common/api-base"
	"dtstack.com/dtstack/easymatrix/matrix/log"
	"dtstack.com/dtstack/easymatrix/matrix/model"
	"fmt"
	"github.com/kataras/iris/context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var downloaderMap = make(map[int]*Downloader, 0)

const (
	STATUS_RUNNING      = "running"
	STATUS_SUCCESS      = "success"
	STATUS_CANCEL       = "cancel"
	STATUS_FAIL         = "fail"
	UPLOAD_TYPE_PRODUCT = 1
)

type uploadAsyncParam struct {
	Name []string `json:"name"`
}

type uploadSyncParam struct {
	Name string `json:"name"`
}

// UploadAsync
// @Description  	upload package async
// @Summary      	异步上传包
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           message body uploadAsyncParam true "主机密码信息"
// @Success         200 {object} string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/uploadAsync [post]
func UploadAsync(ctx context.Context) apibase.Result {
	param := uploadAsyncParam{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("[UploadAsync] invalid param")
		return fmt.Errorf("参数格式非法")
	}
	err := CheckParams(param)
	if err != nil {
		return err
	}

	userId, err := ctx.Values().GetInt("userId")
	type fileIdName struct {
		Id   int
		Name string
	}
	var fileList = make([]fileIdName, 0)
	var canceledList = make([]string, 0)
	var wg sync.WaitGroup
	for _, url := range param.Name {
		// 开始下载流程
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			uploadId, err := model.UploadRecord.AddUploadRecord(url, STATUS_RUNNING, UPLOAD_TYPE_PRODUCT)
			if err == nil {
				file, err := DownloadFile(int(uploadId), url)
				if err != nil {
					log.Errorf("[UploadAsync] download file error: %v", err)
				}
				_, err = model.UploadRecord.GetCancelProductById(int(uploadId))
				if err == nil {
					log.Errorf("[UploadAsync] upload %v canceled", uploadId)
					canceledList = append(canceledList, file)
				} else {
					fileList = append(fileList, fileIdName{int(uploadId), file})
				}
			}
		}(url)
	}
	wg.Wait()
	for _, file := range canceledList {
		_ = os.Remove(file)
	}
	errors := make([]string, 0)
	for _, file := range fileList {
		result := DoUnzipAndParse(file.Name, userId)
		model.UploadRecord.DeleteUploadingProduct(file.Id)
		log.Infof("[UploadAsync] file: %v, result: %v", file, result)
		if _, ok := result.(error); ok {
			errors = append(errors, result.(error).Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

func UploadSync(ctx context.Context) apibase.Result {
	param := uploadSyncParam{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("[UploadSync] invalid param")
		return fmt.Errorf("参数格式非法")
	}
	type fileIdName struct {
		Id   int
		Name string
	}
	url := param.Name
	var fileList = make([]fileIdName, 0)
	var canceledList = make([]string, 0)
	// 开始下载流程
	uploadId, err := model.UploadRecord.AddUploadRecord(url, STATUS_RUNNING, UPLOAD_TYPE_PRODUCT)
	if err == nil {
		file, err := DownloadFile(int(uploadId), url)
		if err != nil {
			log.Errorf("[UploadAsync] download file error: %v", err)
		}
		_, err = model.UploadRecord.GetCancelProductById(int(uploadId))
		if err == nil {
			log.Errorf("[UploadAsync] upload %v canceled", uploadId)
			canceledList = append(canceledList, file)
		} else {
			fileList = append(fileList, fileIdName{int(uploadId), file})
		}
	}
	for _, file := range canceledList {
		_ = os.Remove(file)
	}
	errors := make([]string, 0)
	for _, file := range fileList {
		result := DoUnzipAndParse(file.Name, 1)
		model.UploadRecord.DeleteUploadingProduct(file.Id)
		log.Infof("[UploadAsync] file: %v, result: %v", file, result)
		if _, ok := result.(error); ok {
			errors = append(errors, result.(error).Error())
		} else {
			return result
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}

func DoUnzipAndParse(file string, userId int) apibase.Result {
	uploadLock.Lock()
	defer uploadLock.Unlock()
	f, err := os.Open(file)
	defer func() {
		_ = f.Close()
		os.Remove(file)
	}()
	if err != nil {
		log.Errorf("open file error: %v", err)
	}
	return UnzipAndParse(f, userId)
}

// CheckAvailableLink
// @Description  	异步上传预检查
// @Summary      	upload async check
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           message body uploadAsyncParam true "主机密码信息"
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/check_param [post]
func CheckAvailableLink(ctx context.Context) apibase.Result {
	param := uploadAsyncParam{}
	if err := ctx.ReadJSON(&param); err != nil {
		log.Errorf("[CheckAvailableLink] param error: %v", err)
		return err
	}
	return CheckParams(param)
}

func CheckParams(param uploadAsyncParam) error {
	if len(param.Name) == 0 {
		return fmt.Errorf("请输入地址")
	}
	var respList []*http.Response
	defer func() {
		if len(respList) > 0 {
			for _, resp := range respList {
				if resp != nil {
					_ = resp.Body.Close()
				}
			}
		}
	}()
	var err error
	for _, value := range param.Name {
		resp, err := http.Get(value)
		respList = append(respList, resp)
		if err != nil || (resp != nil && resp.StatusCode != 200) {
			log.Errorf("[UploadAsync] url: %s invalid, err: %v", value, err)
			return fmt.Errorf("格式错误或无效地址")
		}
	}
	return err
}

func DownloadFile(uploadId int, url string) (string, error) {
	out, err := ioutil.TempFile("", "upload-file")
	if err != nil {
		log.Errorf("[Download] create temp file error: %v", err)
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Errorf("[Download] http get url error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	downloader := NewDownloader(uploadId, uint64(resp.ContentLength))
	downloaderMap[uploadId] = downloader
	go downloader.ReportProgress()

	_, err = io.Copy(out, io.TeeReader(resp.Body, downloader))
	if err != nil {
		log.Errorf("[Download] process error: %v", err)
		downloader.failCh <- true
		return "", err
	}
	return out.Name(), nil
}

type SimpleUploadProduct struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
	Status  string  `json:"status"`
}

// GetUploadingProducts
// @Description  	Get Uploading Products
// @Summary      	查询包上传状态
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Success         200 {object}  string "{"msg":"ok","code":0,"data":{"total":"","data":""}}"
// @Router          /api/v2/product/in_progress [get]
func GetUploadingProducts(ctx context.Context) apibase.Result {
	uploadingProducts, err := model.UploadRecord.GetUploadingProduct(STATUS_RUNNING, STATUS_SUCCESS)
	if err != nil {
		log.Errorf("[GetUploadingProducts] query from database error: %v", err)
		return err
	}
	uploadingResult := make([]SimpleUploadProduct, 0)
	var transfer = func(name string) string {
		if parts := strings.Split(name, "/"); len(parts) > 1 {
			return parts[len(parts)-1]
		}
		return name
	}
	if len(uploadingProducts) > 0 {
		for _, uploadingProduct := range uploadingProducts {
			uploadingResult = append(uploadingResult, SimpleUploadProduct{
				Id:      uploadingProduct.Id,
				Name:    transfer(uploadingProduct.Name),
				Percent: uploadingProduct.Progress,
				Status:  uploadingProduct.Status,
			})
		}
	}
	return map[string]interface{}{
		"total": len(uploadingResult),
		"data":  uploadingResult,
	}
}

// CancelUpload
// @Description  	Cancel Upload
// @Summary      	取消上传包
// @Tags         	product
// @Accept          application/json
// @Produce 		application/json
// @Param           record_id formData string  true "上传ID"
// @Success         200 {object} string "{"msg":"ok","code":0,"data":null}"
// @Router          /api/v2/product/cancel_upload [post]
func CancelUpload(ctx context.Context) apibase.Result {
	uploadIdStr := ctx.FormValue("record_id")
	uploadId, _ := strconv.Atoi(uploadIdStr)
	err := model.UploadRecord.DeleteUploadingProduct(uploadId)
	if err != nil {
		log.Errorf("[CancelUpload] delete database error:%v", err)
		return err
	}
	downloader, ok := downloaderMap[uploadId]
	if ok {
		close(downloader.cancelCh)
	} else {
		log.Errorf("[CancelUpload] downloader:%d cache not exist", uploadId)
	}
	return nil
}

type IDownloder interface {
	ReportProgress()
}

type Downloader struct {
	Id         int
	Total      uint64
	Current    uint64
	completeCh chan interface{}
	failCh     chan interface{}
	cancelCh   chan interface{}
}

func NewDownloader(id int, total uint64) *Downloader {
	return &Downloader{
		Id:         id,
		Total:      total,
		Current:    0,
		completeCh: make(chan interface{}, 1),
		failCh:     make(chan interface{}, 1),
		cancelCh:   make(chan interface{}, 1),
	}
}

func (d *Downloader) Write(p []byte) (int, error) {
	n := len(p)
	d.Current += uint64(n)
	if d.Current == d.Total {
		close(d.completeCh)
	}
	return n, nil
}

func (d *Downloader) ReportProgress() {
	defer func() {
		delete(downloaderMap, d.Id)
	}()
	ticker := time.NewTicker(3 * time.Second)
	ctx, cancel := sysContext.WithTimeout(sysContext.Background(), 60*time.Minute)
	defer cancel()
	// report progress
	for {
		select {
		case <-ticker.C:
			// report progress, go on loop
			progress, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float64(d.Current)/float64(d.Total)), 64)
			err := model.UploadRecord.UpdateProgress(d.Id, progress*100, STATUS_RUNNING)
			log.Infof("current: %v, total: %v, progress: %v", d.Current, d.Total, progress)
			if err != nil {
				log.Errorf("[UploadRecord] update progress regularly error: %v", err)
			}
		case <-d.completeCh:
			// report progress, break loop
			err := model.UploadRecord.UpdateProgress(d.Id, 100, STATUS_SUCCESS)
			if err != nil {
				log.Errorf("[UploadRecord] update progress success error: %v", err)
			}
			return
		case <-d.failCh:
			// report progress, break loop
			progress, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", float64(d.Current)/float64(d.Total)), 64)
			err := model.UploadRecord.UpdateProgress(d.Id, progress*100, STATUS_FAIL)
			if err != nil {
				log.Errorf("[UploadRecord] update progress fail error: %v", err)
			}
			return
		case <-d.cancelCh:
			log.Infof("[UploadRecord] cancel upload")
			return
		case <-ctx.Done():
			err := ctx.Err()
			log.Errorf("[UploadRecord] error: %v", err)
			return
		default:
		}
	}
}
