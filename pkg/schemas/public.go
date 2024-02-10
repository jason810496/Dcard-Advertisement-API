package schemas

type PublicAdItem struct {
	Title string `json:"title" example:"This is an AD title"`
	EndAt string `json:"endAt" example:"2021-12-31 23:59:59"`
}

type PublicAdResponse struct {
	Items []PublicAdItem `json:"items"`
}
