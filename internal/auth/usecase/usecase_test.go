package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	//"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/auth/mock"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
)

func TestAuthUC_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, nil, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}

	ctx := context.Background()
	//span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.UploadAvatar")
	//defer span.Finish()

	//mockAuthRepo.EXPECT().FindByEmail(ctxWithTrace, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	mockAuthRepo.EXPECT().FindByEmail(ctx, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	//mockAuthRepo.EXPECT().Register(ctxWitWithTrace, gomock.Eq(user)).Return(nil, sql.ErrNoRows)
	mockAuthRepo.EXPECT().Register(ctx, gomock.Eq(user)).Return(user, nil)

	createdUSer, err := authUC.Register(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, createdUSer)
	require.Nil(t, err)
}

func TestAuthUC_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()
	//span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.Update")
	//defer span.Finish()

	//mockAuthRepo.EXPECT().Update(ctxWithTrace, gomock.Eq(user)).Return(user, nil)
	//mockRedisRepo.EXPECT().DeleteUserCtx(ctxWithTrace, key).Return(nil)
	mockAuthRepo.EXPECT().Update(ctx, gomock.Eq(user)).Return(user, nil)
	mockRedisRepo.EXPECT().DeleteUserCtx(ctx, key).Return(nil)

	updatedUser, err := authUC.Update(ctx, user)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Nil(t, err)
}

func TestAuthUC_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)

	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()
	//span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.Delete")
	//defer span.Finish()

	//mockAuthRepo.EXPECT().Delete(ctxWithTrace, gomock.Eq(user.UserID)).Return(nil)
	//mockRedisRepo.EXPECT().DeleteUserCtx(ctxWithTrace, key).Return(nil)
	mockAuthRepo.EXPECT().Delete(ctx, gomock.Eq(user.UserID)).Return(nil)
	mockRedisRepo.EXPECT().DeleteUserCtx(ctx, key).Return(nil)

	err := authUC.Delete(ctx, user.UserID)
	require.NoError(t, err)
	require.Nil(t, err)
}

func TestAuthUC_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}
	key := fmt.Sprintf("%s: %s", basePrefix, user.UserID)

	ctx := context.Background()
	//span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.GetByID")
	//defer span.Finish()

	//mockRedisRepo.EXPECT().GetByIDCtx(ctxWithTrace, key).Return(nil, nil)
	//mockAuthRepo.EXPECT().GetByID(ctxWithTrace, gomock.Eq(user.UserID)).Return(user, nil)
	//mockRedisRepo.EXPECT().SetUserCtx(ctxWithTrace, key, cacheDuration, user).Return(nil)
	mockRedisRepo.EXPECT().GetByIDCtx(ctx, key).Return(nil, nil)
	mockAuthRepo.EXPECT().GetByID(ctx, gomock.Eq(user.UserID)).Return(user, nil)
	mockRedisRepo.EXPECT().SetUserCtx(ctx, key, cacheDuration, user).Return(nil)

	u, err := authUC.GetByID(ctx, user.UserID)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, u)
}

func TestAuthUC_FindByName(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	userName := "name"
	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}
	ctx := context.Background()
	//span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.FindByName")
	//defer span.Finish()

	usersList := &models.UsersList{}

	//mockAuthRepo.EXPECT().FindByName(ctxWithTrace, gomock.Eq(userName), query).Return(usersList, nil)
	mockAuthRepo.EXPECT().FindByName(ctx, gomock.Eq(userName), query).Return(usersList, nil)

	userList, err := authUC.FindByName(ctx, userName, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userList)
}

func TestAuthUC_GetUsers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	query := &utils.PaginationQuery{
		Size:    10,
		Page:    1,
		OrderBy: "",
	}
	ctx := context.Background()
	//span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.GetUsers")
	//defer span.Finish()

	usersList := &models.UsersList{}

	//mockAuthRepo.EXPECT().GetUsers(ctxWithTrace, query).Return(usersList, nil)
	mockAuthRepo.EXPECT().GetUsers(ctx, query).Return(usersList, nil)

	users, err := authUC.GetUsers(ctx, query)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, users)
}

func TestAuthUC_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthRepo := mock.NewMockRepository(ctrl)
	mockRedisRepo := mock.NewMockRedisRepository(ctrl)
	authUC := NewAuthUseCase(cfg, mockAuthRepo, mockRedisRepo, apiLogger)

	ctx := context.Background()
	//span, ctxWithTrace := opentracing.StartSpanFromContext(ctx, "authUC.Login")
	//defer span.Finish()

	user := &models.User{
		Password: "123456",
		Email:    "email@gmail.com",
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	mockUser := &models.User{
		Email:    "email@gmail.com",
		Password: string(hashPassword),
	}

	//mockAuthRepo.EXPECT().FindByEmail(ctxWithTrace, gomock.Eq(user)).Return(mockUser, nil)
	mockAuthRepo.EXPECT().FindByEmail(ctx, gomock.Eq(user)).Return(mockUser, nil)

	userWithToken, err := authUC.Login(ctx, user)
	require.NoError(t, err)
	require.Nil(t, err)
	require.NotNil(t, userWithToken)
}
