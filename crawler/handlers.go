package crawler

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
	return nil
}

func (h *PdfFileDownloadHandler) save(res *Resource) error {
	return nil
}
