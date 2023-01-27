package board

import (
	"encoding/json"
	"io"
	"net/http"
)

type Boards struct {
	Boards []struct {
		Board           string `json:"board"`
		Title           string `json:"title"`
		WsBoard         int    `json:"ws_board"`
		PerPage         int    `json:"per_page"`
		Pages           int    `json:"pages"`
		MaxFilesize     int    `json:"max_filesize"`
		MaxWebmFilesize int    `json:"max_webm_filesize"`
		MaxCommentChars int    `json:"max_comment_chars"`
		MaxWebmDuration int    `json:"max_webm_duration"`
		BumpLimit       int    `json:"bump_limit"`
		ImageLimit      int    `json:"image_limit"`
		Cooldowns       struct {
			Threads int `json:"threads"`
			Replies int `json:"replies"`
			Images  int `json:"images"`
		} `json:"cooldowns"`
		MetaDescription string `json:"meta_description"`
		IsArchived      int    `json:"is_archived,omitempty"`
		Spoilers        int    `json:"spoilers,omitempty"`
		CustomSpoilers  int    `json:"custom_spoilers,omitempty"`
		UserIds         int    `json:"user_ids,omitempty"`
		CountryFlags    int    `json:"country_flags,omitempty"`
		CodeTags        int    `json:"code_tags,omitempty"`
		WebmAudio       int    `json:"webm_audio,omitempty"`
		MinImageWidth   int    `json:"min_image_width,omitempty"`
		MinImageHeight  int    `json:"min_image_height,omitempty"`
		Oekaki          int    `json:"oekaki,omitempty"`
		SjisTags        int    `json:"sjis_tags,omitempty"`
		BoardFlags      struct {
			FourCC string `json:"4CC"`
			Ada    string `json:"ADA"`
			An     string `json:"AN"`
			Anf    string `json:"ANF"`
			Apb    string `json:"APB"`
			Aj     string `json:"AJ"`
			Ab     string `json:"AB"`
			Au     string `json:"AU"`
			Bb     string `json:"BB"`
			Bm     string `json:"BM"`
			Bp     string `json:"BP"`
			Bs     string `json:"BS"`
			Cl     string `json:"CL"`
			Co     string `json:"CO"`
			Cg     string `json:"CG"`
			Che    string `json:"CHE"`
			Cb     string `json:"CB"`
			Day    string `json:"DAY"`
			Dd     string `json:"DD"`
			Der    string `json:"DER"`
			Dt     string `json:"DT"`
			Dis    string `json:"DIS"`
			Eqa    string `json:"EQA"`
			Eqf    string `json:"EQF"`
			Eqp    string `json:"EQP"`
			Eqr    string `json:"EQR"`
			Eqt    string `json:"EQT"`
			Eqi    string `json:"EQI"`
			Eqs    string `json:"EQS"`
			Era    string `json:"ERA"`
			Fau    string `json:"FAU"`
			Fle    string `json:"FLE"`
			Fl     string `json:"FL"`
			Gi     string `json:"GI"`
			Ht     string `json:"HT"`
			Iz     string `json:"IZ"`
			Li     string `json:"LI"`
			Lt     string `json:"LT"`
			Ly     string `json:"LY"`
			Ma     string `json:"MA"`
			Mau    string `json:"MAU"`
			Min    string `json:"MIN"`
			Ni     string `json:"NI"`
			Nur    string `json:"NUR"`
			Oct    string `json:"OCT"`
			Par    string `json:"PAR"`
			Pc     string `json:"PC"`
			Pce    string `json:"PCE"`
			Pi     string `json:"PI"`
			Plu    string `json:"PLU"`
			Pm     string `json:"PM"`
			Pp     string `json:"PP"`
			Qc     string `json:"QC"`
			Rar    string `json:"RAR"`
			Rd     string `json:"RD"`
			Rlu    string `json:"RLU"`
			S1L    string `json:"S1L"`
			Sco    string `json:"SCO"`
			Shi    string `json:"SHI"`
			Sil    string `json:"SIL"`
			Son    string `json:"SON"`
			Sp     string `json:"SP"`
			Spi    string `json:"SPI"`
			Ss     string `json:"SS"`
			Sta    string `json:"STA"`
			Stl    string `json:"STL"`
			Spt    string `json:"SPT"`
			Sun    string `json:"SUN"`
			Sus    string `json:"SUS"`
			Swb    string `json:"SWB"`
			Tfa    string `json:"TFA"`
			Tfo    string `json:"TFO"`
			Tfp    string `json:"TFP"`
			Tfs    string `json:"TFS"`
			Tft    string `json:"TFT"`
			Tfv    string `json:"TFV"`
			Tp     string `json:"TP"`
			Ts     string `json:"TS"`
			Twi    string `json:"TWI"`
			Tx     string `json:"TX"`
			Vs     string `json:"VS"`
			Ze     string `json:"ZE"`
			Zs     string `json:"ZS"`
		} `json:"board_flags,omitempty"`
		TextOnly       int `json:"text_only,omitempty"`
		RequireSubject int `json:"require_subject,omitempty"`
		BoardFlags0    struct {
			Ac string `json:"AC"`
			An string `json:"AN"`
			Bl string `json:"BL"`
			Cf string `json:"CF"`
			Cm string `json:"CM"`
			Ct string `json:"CT"`
			Dm string `json:"DM"`
			Eu string `json:"EU"`
			Fc string `json:"FC"`
			Gn string `json:"GN"`
			Gy string `json:"GY"`
			Jh string `json:"JH"`
			Kn string `json:"KN"`
			Mf string `json:"MF"`
			Nb string `json:"NB"`
			Nt string `json:"NT"`
			Nz string `json:"NZ"`
			Pc string `json:"PC"`
			Pr string `json:"PR"`
			Re string `json:"RE"`
			Mz string `json:"MZ"`
			Tm string `json:"TM"`
			Tr string `json:"TR"`
			Un string `json:"UN"`
			Wp string `json:"WP"`
		} `json:"board_flags,omitempty"`
		MathTags int `json:"math_tags,omitempty"`
	} `json:"boards"`
}

func GetBoards() (b Boards, err error) {
	resp, err := http.Get("https://a.4cdn.org/boards.json")
	if err != nil {
		return Boards{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Boards{}, err
	}
	resp.Body.Close()

	boards := Boards{}
	if err = json.Unmarshal(body, &boards); err != nil {
		return Boards{}, err
	}

	return boards, nil
}
