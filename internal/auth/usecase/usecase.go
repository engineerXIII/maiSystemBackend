package usecase

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/auth"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/pkg/httpErrors"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
)

const (
	basePrefix    = "api-auth"
	cacheDuration = 3600
)

// Auth UseCase
type authUC struct {
	cfg       *config.Config
	authRepo  auth.Repository
	redisRepo auth.RedisRepository
	logger    logger.Logger
}

// Auth UseCase constructor
func NewAuthUseCase(cfg *config.Config, authRepo auth.Repository, redisRepo auth.RedisRepository /*, awsRepo auth.AWSRepository*/, log logger.Logger) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo, redisRepo: redisRepo /* awsRepo: awsRepo,*/, logger: log}
}

// Create new user
func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Register")
	defer span.Finish()

	existsUser, err := u.authRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareCreate"))
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Update existing user
func (u *authUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Update")
	defer span.Finish()

	if err := user.PrepareUpdate(); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "authUC.Register.PrepareUpdate"))
	}

	updatedUser, err := u.authRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	if err = u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(user.UserID.String())); err != nil {
		u.logger.Errorf("AuthUC.Update.DeleteUserCtx: %s", err)
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

// Delete new user
func (u *authUC) Delete(ctx context.Context, userID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Delete")
	defer span.Finish()

	if err := u.authRepo.Delete(ctx, userID); err != nil {
		return err
	}

	if err := u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(userID.String())); err != nil {
		u.logger.Errorf("AuthUC.Delete.DeleteUserCtx: %s", err)
	}

	return nil
}

// Get user by id
func (u *authUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	//span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetByID")
	//defer span.Finish()

	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.GenerateUserKey(userID.String()))
	if err != nil {
		u.logger.Errorf("authUC.GetByID.GetByIDCtx: %v", err)
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := u.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err = u.redisRepo.SetUserCtx(ctx, u.GenerateUserKey(userID.String()), cacheDuration, user); err != nil {
		u.logger.Errorf("authUC.GetByID.SetUserCtx: %v", err)
	}

	user.SanitizePassword()

	return user, nil
}

// Find users by name
func (u *authUC) FindByName(ctx context.Context, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.FindByName")
	defer span.Finish()

	return u.authRepo.FindByName(ctx, name, query)
}

// Get users with pagination
func (u *authUC) GetUsers(ctx context.Context, pq *utils.PaginationQuery) (*models.UsersList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.GetUsers")
	defer span.Finish()

	return u.authRepo.GetUsers(ctx, pq)
}

// Login user, returns user model with jwt token
func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "authUC.Login")
	defer span.Finish()

	foundUser, err := u.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "authUC.GetUsers.ComparePasswords"))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "authUC.GetUsers.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (u *authUC) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}
