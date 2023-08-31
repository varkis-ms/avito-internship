package entity

type UserAddToSegmentRequest struct {
	UserId   int      `json:"user_id"       binding:"required"  example:"1000"`
	Segments []string `json:"segments" binding:"required"  example:"AVITO_VOICE_MESSAGES,AVITO_PERFORMANCE_VAS"`
	Ttl      int      `json:"ttl"                          example:"2"`
}

type UserRemoveFromSegmentRequest struct {
	UserId   int      `json:"user_id"       binding:"required"  example:"1000"`
	Segments []string `json:"segments" binding:"required"  example:"AVITO_VOICE_MESSAGES,AVITO_PERFORMANCE_VAS"`
}

type UserActiveSegmentRequest struct {
	UserId int
}
