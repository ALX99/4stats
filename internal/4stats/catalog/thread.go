package catalog

type Thread struct {
	No            int64   `json:"no"`
	Sticky        int     `json:"sticky,omitempty"`
	Closed        int     `json:"closed,omitempty"`
	Now           string  `json:"now"`
	Name          string  `json:"name"`
	Sub           string  `json:"sub,omitempty"`
	Com           string  `json:"com"`
	Filename      string  `json:"filename"`
	Ext           string  `json:"ext"`
	W             int     `json:"w"`
	H             int     `json:"h"`
	TnW           int     `json:"tn_w"`
	TnH           int     `json:"tn_h"`
	Tim           int64   `json:"tim"`
	Time          int     `json:"time"`
	Md5           string  `json:"md5"`
	Fsize         int     `json:"fsize"`
	Resto         int     `json:"resto"`
	Capcode       string  `json:"capcode,omitempty"`
	SemanticURL   string  `json:"semantic_url"`
	Replies       int     `json:"replies"`
	Images        int     `json:"images"`
	OmittedPosts  int     `json:"omitted_posts,omitempty"`
	OmittedImages int     `json:"omitted_images,omitempty"`
	LastReplies   []Reply `json:"last_replies"`
}
