package asdf

type Std struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Errno  int    `json:"errno"`
}
