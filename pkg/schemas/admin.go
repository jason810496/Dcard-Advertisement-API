package schemas

type CreateAdRequest struct {
	Title      string `json:"title" binding:"required" example:"AD 123"`
	StartAt    string `json:"startAt" binding:"required" example:"2024-01-01 10:00:00"`
	EndAt      string `json:"endAt" binding:"required" example:"2024-02-01 10:00:00"`
	Conditions []struct {
		AgeStart int      `json:"ageStart" binding:"required" example:"18"`
		AgeEnd   int      `json:"ageEnd" binding:"required" example:"30"`
		Country  []string `json:"country" binding:"required" example:"TW,JP,KR,US"`
		Platform []string `json:"platform" binding:"required" example:"ios,web"`
	} `json:"conditions" binding:"required"`
}

type CreateAdResponse struct {
	Message string `json:"message" example:"create success"`
}
