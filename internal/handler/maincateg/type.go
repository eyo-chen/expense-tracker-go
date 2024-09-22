package maincateg

type mainCateg struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IconType string `json:"icon_type"`
	IconData string `json:"icon_data"`
}

type getAllMainCategResp struct {
	Categories []mainCateg `json:"categories"`
}
