package schemas

type PublicAdItem struct {
	Title string `json:"title" example:"This is an AD title"`
	EndAt string `json:"endAt" example:"2025-01-01T00:00:00+08:00"`
}

type PublicAdResponse struct {
	Items []PublicAdItem `json:"items"`
}

type PublicAdRequest struct {
	Age      *int   `form:"age" binding:"omitempty,numeric,min=1,max=100" example:"18"`
	Country  string `form:"country" binding:"omitempty,oneof=TW HK JP US KR" example:"TW"`
	Platform string `form:"platform" binding:"omitempty,oneof=ios android web" example:"ios"`
	Gender   string `form:"gender" binding:"omitempty,oneof=F M" example:"F"`
	Limit    *int   `form:"limit" binding:"omitempty,numeric,min=1" example:"10"`
	Offset   *int   `form:"offset" binding:"omitempty,numeric,min=0" example:"0"`
}
