package config

type AdditionalFileResp struct {
	Canceled bool
	FilePath string
}

type AdditionalFileProvider interface {
	Watch() <-chan AdditionalFileResp
}
