package usecase

import (
	"context"
	"github.com/engineerXIII/maiSystemBackend/config"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/internal/product"
	"github.com/engineerXIII/maiSystemBackend/pkg/httpErrors"
	"github.com/engineerXIII/maiSystemBackend/pkg/logger"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// Redis variables
//const (
//	basePrefix    = "api-products"
//	cacheDuration = 3600
//)

type productUC struct {
	cfg         *config.Config
	productRepo product.Repository
	logger      logger.Logger
}

func NewProductUseCase(cfg *config.Config, productRepo product.Repository, logger logger.Logger) product.UseCase {
	return &productUC{cfg: cfg, productRepo: productRepo, logger: logger}
}

func (u *productUC) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Create")
	defer span.Finish()

	//user, err := utils.GetUserFromCtx(ctx)
	//if err != nil {
	//	return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "productUC.Create.GetUserFromCtx"))
	//}

	//if err = utils.ValidateStruct(ctx, product); err != nil {
	if err := utils.ValidateStruct(ctx, product); err != nil {
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "productUC.Create.ValidateStruct"))
	}

	p, err := u.productRepo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return p, err
}

func (u *productUC) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Update")
	defer span.Finish()

	//productID, err := u.productRepo.GetProductByID(ctx, product.ProductID)
	_, err := u.productRepo.GetProductByID(ctx, product.ProductID)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := u.productRepo.Update(ctx, product)
	if err != nil {
		return nil, err
	}

	return updatedProduct, nil
}

func (u *productUC) GetProductByID(ctx context.Context, productUUID uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.GetProductByID")
	defer span.Finish()

	p, err := u.productRepo.GetProductByID(ctx, productUUID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (u *productUC) Delete(ctx context.Context, productUUID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.Delete")
	defer span.Finish()

	_, err := u.productRepo.GetProductByID(ctx, productUUID)
	if err != nil {
		return err
	}

	if err := u.productRepo.Delete(ctx, productUUID); err != nil {
		return err
	}

	return nil
}

func (u *productUC) GetProducts(ctx context.Context, pq *utils.PaginationQuery) (*models.ProductList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.GetProducts")
	defer span.Finish()

	return u.productRepo.GetProducts(ctx, pq)
}

func (u *productUC) SearchByName(ctx context.Context, name string, pq *utils.PaginationQuery) (*models.ProductList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productUC.SearchByName")
	defer span.Finish()

	return u.productRepo.SearchByName(ctx, name, pq)
}
