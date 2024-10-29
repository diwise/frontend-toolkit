package frontendtoolkit

type Asset interface {
	Body() []byte
	ContentLength() int
	ContentType() string
	Path() string
	SHA256() string
}

type AssetLoaderFunc func(name string) Asset

type AssetLoader interface {
	Load(name string) Asset
	LoadFromSha256(sha string) (Asset, error)
}

type Bundle interface {
	For(acceptLanguage string) Localizer
}

type Localizer interface {
	Get(string) string
	GetWithData(string, map[string]any) string
}
