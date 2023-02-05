package catalog

type Reply struct {
	No       int64  `json:"no"`
	Now      string `json:"now"`
	Name     string `json:"name"`
	Com      string `json:"com"`
	Filename string `json:"filename"`
	Ext      string `json:"ext"`
	W        int    `json:"w"`
	H        int    `json:"h"`
	TnW      int    `json:"tn_w"`
	TnH      int    `json:"tn_h"`
	Tim      int64  `json:"tim"`
	Time     int    `json:"time"`
	Md5      string `json:"md5"`
	Fsize    int    `json:"fsize"`
	Resto    int    `json:"resto"`
	Capcode  string `json:"capcode"`
}
