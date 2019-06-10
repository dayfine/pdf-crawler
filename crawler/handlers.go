package crawler

import (
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Handler interface {
	Handle(res *Resource) error
}

type PdfFileDownloadHandler struct {
	saveDir string
}

func NewPdfFileDownloadHandler(saveDir string) Handler {
	return &PdfFileDownloadHandler{saveDir}
}

func (h *PdfFileDownloadHandler) Handle(res *Resource) error {
	url, _ := url.Parse(res.url)
	directory := filepath.Join(h.saveDir, url.Host)
	os.MkdirAll(directory, os.ModePerm)
	filename := filepath.Join(directory, uuid.New().String()+".pdf")
	err := h.save(filename, res.resp.Body)
	return err
}

func (h *PdfFileDownloadHandler) save(filename string, content io.ReadCloser) error {
	dest, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, content)
	return err
}
