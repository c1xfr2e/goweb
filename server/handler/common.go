package handler

// Date const
const (
	DateFormat       = "2006/01/02"
	DefaultDateRange = 15
)

// ReqFingerprint
type ReqFingerprint struct {
	Hash1 string `json:"b"`
	Hash2 string `json:"c"`
	Hash3 string `json:"g"`
	Hash4 string `json:"d"`
	Plain string `json:"info"`
}
