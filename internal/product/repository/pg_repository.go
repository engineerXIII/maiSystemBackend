package repository

import (
	"context"
	"database/sql"
	"github.com/engineerXIII/maiSystemBackend/internal/models"
	"github.com/engineerXIII/maiSystemBackend/internal/product"
	"github.com/engineerXIII/maiSystemBackend/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) product.Repository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepo.Create")
	defer span.Finish()

	var p models.Product
	if err := r.db.QueryRowxContext(
		ctx,
		createProduct,
		&product.Name,
		&product.Color,
		&product.Description,
		&product.Factory,
		&product.Cost,
	).StructScan(&p); err != nil {
		return nil, errors.Wrap(err, "productRepo.Create.QueryRowxContext")
	}

	return &p, nil
}

func (r *productRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepo.Update")
	defer span.Finish()

	var p models.Product
	if err := r.db.QueryRowxContext(
		ctx,
		updateProduct,
		&product.Name,
		&product.Color,
		&product.Description,
		&product.Factory,
		&product.Cost,
		&product.ProductID,
	).StructScan(&p); err != nil {
		return nil, errors.Wrap(err, "productRepo.Update.QueryRowxContext")
	}

	return &p, nil
}

func (r *productRepo) GetProductByID(ctx context.Context, productID uuid.UUID) (*models.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepo.GetProductByID")
	defer span.Finish()

	p := &models.Product{}
	if err := r.db.GetContext(ctx, p, getProductByID, productID); err != nil {
		return nil, errors.Wrap(err, "productRepo.GetProductByID.GetContext")
	}

	return p, nil
}

func (r *productRepo) Delete(ctx context.Context, productID uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepo.Delete")
	defer span.Finish()

	result, err := r.db.ExecContext(ctx, deleteProduct, productID)
	if err != nil {
		return errors.Wrap(err, "productRepo.Delete.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "productRepo.Delete.RowsAffected")
	}

	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "productRepo.Delete.RowsAffected")

	}

	return nil
}

func (r *productRepo) GetProducts(ctx context.Context, pq *utils.PaginationQuery) (*models.ProductList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepo.GetProducts")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount); err != nil {
		return nil, errors.Wrap(err, "productRepo.GetProducts.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.ProductList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Products:   make([]*models.Product, 0),
		}, nil
	}

	var productList = make([]*models.Product, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, getProducts, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "productRepo.GetProducts.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		p := &models.Product{}
		if err = rows.StructScan(p); err != nil {
			return nil, errors.Wrap(err, "productRepo.GetProducts.StructScan")
		}
		productList = append(productList, p)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "productRepo.GetProducts.rows.Err")
	}

	return &models.ProductList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Products:   productList,
	}, nil
}

func (r *productRepo) SearchByName(ctx context.Context, name string, pq *utils.PaginationQuery) (*models.ProductList, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "productRepo.SearchByName")
	defer span.Finish()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, findByNameCount, name); err != nil {
		return nil, errors.Wrap(err, "productRepo.SearchByName.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.ProductList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Products:   make([]*models.Product, 0),
		}, nil
	}

	var productList = make([]*models.Product, pq.GetSize())
	rows, err := r.db.QueryxContext(ctx, findByName, name, pq.GetOffset(), pq.GetLimit())
	if err != nil {
		return nil, errors.Wrap(err, "productRepo.SearchByName.QueryxContext")
	}
	defer rows.Close()

	for rows.Next() {
		p := &models.Product{}
		if err = rows.StructScan(p); err != nil {
			return nil, errors.Wrap(err, "productRepo.SearchByName.StructScan")
		}
		productList = append(productList, p)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "productRepo.SearchByName.rows.Err")
	}

	return &models.ProductList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Products:   productList,
	}, nil
}
