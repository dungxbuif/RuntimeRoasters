package http

import (
	"net/http"

	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/auth-service/internal/usecase"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: u,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	if id, ok := identity.FromContext(ctx); ok {
		if id.Role != "farm_admin" && id.Role != "admin" {
			log.Warn("forbidden attempt to create user", zap.String("user_id", id.Subject), zap.String("role", id.Role))
			c.Error(errs.ErrForbidden)
			c.Abort()
			return
		}
	} else {
		c.Error(errs.ErrUnauthorized)
		c.Abort()
		return
	}

	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errs.ErrValidation).SetMeta(err.Error())
		c.Abort()
		return
	}

	user, err := h.userUsecase.CreateUser(ctx, req)
	if err != nil {
		log.Error("failed to create user", zap.Error(err))
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, domain.CreateUserResponse{
		User: user,
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.FromContext(ctx)

	if id, ok := identity.FromContext(ctx); ok {
		if id.Role != "farm_admin" && id.Role != "admin" {
			c.Error(errs.ErrForbidden)
			c.Abort()
			return
		}
	} else {
		c.Error(errs.ErrUnauthorized)
		c.Abort()
		return
	}

	users, err := h.userUsecase.ListUsers(ctx)
	if err != nil {
		log.Error("failed to list users", zap.Error(err))
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
