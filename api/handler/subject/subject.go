package subject

import (
	"github.com/gin-gonic/gin"
	"go-rest-api-boilerplate/api/contextutil"
	"go-rest-api-boilerplate/api/handler"
	"go-rest-api-boilerplate/api/middleware"
	"go-rest-api-boilerplate/api/response"
	db "go-rest-api-boilerplate/internal/db/sqlc"
	"go-rest-api-boilerplate/internal/services/subject"
	"log/slog"
	"strconv"
)

func Routes(h *handler.Handler) {
	subGroup := h.App.Group("/api/v1/web")

	subSGroup := subGroup.Use(middleware.AuthMiddleware(h))
	{
		subSGroup.POST("/subject/create", Add(h))
		subSGroup.GET("/subjects", GetSubjects(h))
		subSGroup.GET("/subject/:slug", GetSubjectByID(h))
		subSGroup.PUT("/subject", Update(h))
	}
}

func Update(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data subject.AddSubjectParams
		if err := ctx.ShouldBindJSON(&data); err != nil {
			slog.ErrorContext(ctx, "Error binding JSON", "error", err)
			response.Error(ctx, nil)
			return
		}

		sub, err := h.SubjectService.Update(ctx, data.ID, data.Name, data.Description, data.Tags)
		if err != nil {
			slog.ErrorContext(ctx, "Error updating subject", "error", err)
			response.Error(ctx, nil)
			return
		}
		subInfo := subject.NewSfSubject(sub.ID, sub.UserID, sub.Name, sub.Description, sub.Tags)
		response.Success(ctx, subInfo)
	}
}

func GetSubjectByID(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slug := ctx.Param("slug")
		id, err := strconv.Atoi(slug)
		if err != nil {
			slog.ErrorContext(ctx, "Error parsing slug to int", "error", err)
			response.Error(ctx, nil)
			return
		}
		sub, err := h.SubjectRepo.GetSubjectByID(ctx, int64(id))
		if err != nil {
			slog.ErrorContext(ctx, "Error getting subject", "error", err)
			response.Error(ctx, nil)
			return
		}

		response.Success(ctx, sub)

	}
}

func GetSubjects(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subs, err := h.SubjectRepo.GetAllSubjects(ctx, db.GetAllSubjectsParams{Limit: 10})
		if err != nil {
			slog.ErrorContext(ctx, "Error getting subjects", "error", err)
			response.Error(ctx, nil)
			return
		}
		response.Success(ctx, subs)
	}
}

func Add(h *handler.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data subject.AddSubjectParams
		if err := ctx.ShouldBindJSON(&data); err != nil {
			response.Error(ctx, nil)
			return
		}
		userID, err := contextutil.GetUser(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "Error getting user from context", "error", err)
			response.Error(ctx, err)
			return
		}
		sub, err := h.SubjectService.AddSubject(ctx, userID, data.Name, data.Description, data.Tags)
		if err != nil {
			slog.ErrorContext(ctx, "Error adding subject", "error", err)
			response.Error(ctx, err)
			return
		}
		response.Success(ctx, sub)
	}
}
