package entity

type SegmentRequest struct {
	Segment string  `json:"segment"       binding:"required"  example:"AVITO_VOICE_MESSAGES"`
	Percent float32 `json:"percent"       example:"0.5"`
}
