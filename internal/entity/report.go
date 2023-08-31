package entity

import "time"

type ReportRequest struct {
	Month int
	Year  int
}

type ReportResponse struct {
	History []ReportUserHistory `json:"history"       binding:"required"`
}

type ReportUserHistory struct {
	UserId    string    `json:"user_id"       binding:"required"`
	Segment   string    `json:"segment"       binding:"required"`
	Operation string    `json:"operation"     binding:"required"`
	Date      time.Time `json:"date"          binding:"required"`
}
