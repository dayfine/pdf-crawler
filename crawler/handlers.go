package crawler

type UrlHandler interface {
	Handle(url string) error
}
