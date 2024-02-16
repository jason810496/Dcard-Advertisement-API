package schemas

import "time"

type CreateAdConditions struct {
	AgeStart int      `json:"ageStart" binding:"omitempty,numeric,min=1,max=100" example:"18"`
	AgeEnd   int      `json:"ageEnd" binding:"omitempty,numeric,min=1,max=100,gtefield=AgeStart" example:"30"`
	Gender   []string `json:"gender" binding:"omitempty,dive,oneof=F M" example:"F,M"`
	Country  []string `json:"country" binding:"omitempty,dive,oneof=TW HK JP US KR" example:"TW,JP"`
	Platform []string `json:"platform" binding:"omitempty,dive,oneof=ios android web" example:"ios,android"`
}

type CreateAdRequest struct {
	Title      string             `json:"title" binding:"required" example:"AD 123"`
	StartAt    time.Time          `json:"startAt" binding:"required" example:"2023-12-10T03:00:00.000Z"`
	EndAt      time.Time          `json:"endAt" binding:"required,gtefield=StartAt" example:"2024-12-31T16:00:00.000Z"`
	Conditions CreateAdConditions `json:"conditions" binding:"required"`
}

type CreateAdResponse struct {
	Message string `json:"message" example:"create success"`
}

func NewCreateAdRequest() CreateAdRequest {
	instance := CreateAdRequest{}
	instance.Conditions.AgeStart = 1
	instance.Conditions.AgeEnd = 100
	return instance
}
