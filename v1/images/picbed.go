package images

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

func DownloadImg(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	img, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	return img
}

func GetBufferFromByte(by []byte) (*bytes.Buffer, string, error) {
	// 创建一个新的 form，然后将文件添加到该 form
	var buffer bytes.Buffer
	w := multipart.NewWriter(&buffer)
	fw, err := w.CreateFormFile("file", strconv.FormatInt(time.Now().Unix(), 10)+".tmp")
	if err != nil {
		return nil, "", err
	}

	_, err = fw.Write(by)
	if err != nil {
		return nil, "", err
	}

	// 关闭 multipart writer，以便写入终止边界
	err = w.Close()
	if err != nil {
		return nil, "", err
	}
	return &buffer, w.FormDataContentType(), nil
}

func UploadImg2PicBed(bys []byte) string {
	if bys == nil {
		return ""
	}
	buffer, contentType, err := GetBufferFromByte(bys)
	if err != nil {
		return ""
	}

	// 创建并发送请求
	req, err := http.NewRequest("POST", "https://cdn.ipfsscan.io/api/v0/add?pin=false", buffer)
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 解析并打印响应
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	type UploadResponse struct {
		Name string `json:"Name"`
		Hash string `json:"Hash"`
		Size string `json:"Size"`
	}

	var uploadResponse UploadResponse
	err = json.Unmarshal(responseData, &uploadResponse)
	if err != nil {
		return ""
	}
	if uploadResponse.Hash == "" {
		return ""
	}
	imgUrl := "https://cdn.ipfsscan.io/ipfs/" + uploadResponse.Hash
	return imgUrl
}
