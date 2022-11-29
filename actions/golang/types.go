package golang

type Platform struct {
	Goos   string `required:"true" json:"goos"`
	Goarch string `required:"true" json:"goarch"`
}
