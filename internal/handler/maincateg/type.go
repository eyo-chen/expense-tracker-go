package maincateg

type icon struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

type mainCateg struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Icon icon   `json:"icon"`
}

type getAllMainCategResp struct {
	Categories []mainCateg `json:"categories"`
}
