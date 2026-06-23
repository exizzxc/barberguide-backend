package repository

import (
	"context"

	"github.com/exizzxc/barberguide-backend/internal/dto"
	"github.com/exizzxc/barberguide-backend/internal/model"
	"gorm.io/gorm"
)

type HaircutRepository interface {
	Create(ctx context.Context, haircut *model.Haircut) error
	FindAll(ctx context.Context, filter dto.HaircutFilterRequest) ([]model.Haircut, int64, error)
	FindByID(ctx context.Context, id uint) (*model.Haircut, error)
	Update(ctx context.Context, haircut *model.Haircut) error
	Delete(ctx context.Context, id uint) error
}

type haircutRepository struct {
	db *gorm.DB
}

func NewHaircutRepository(db *gorm.DB) HaircutRepository {
	return &haircutRepository{db: db}
}

func (r *haircutRepository) Create(ctx context.Context, haircut *model.Haircut) error {
	return r.db.WithContext(ctx).Create(haircut).Error
}

func (r *haircutRepository) FindAll(ctx context.Context, filter dto.HaircutFilterRequest) ([]model.Haircut, int64, error) {
	var haircuts []model.Haircut
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Haircut{}).
		Where("deleted_at IS NULL")

	// Фильтр по форме лица
	if filter.FaceShape != "" {
		query = query.Joins("JOIN haircut_face_shapes ON haircut_face_shapes.haircut_id = haircuts.id").
			Where("haircut_face_shapes.face_shape = ?", filter.FaceShape)
	}

	// Фильтр по типу волос
	if filter.HairType != "" {
		query = query.Joins("JOIN haircut_hair_types ON haircut_hair_types.haircut_id = haircuts.id").
			Where("haircut_hair_types.hair_type = ?", filter.HairType)
	}

	// Фильтр по стилю
	if filter.Style != "" {
		query = query.Where("style = ?", filter.Style)
	}

	// Фильтр по длине
	if filter.Length != "" {
		query = query.Where("length = ?", filter.Length)
	}

	// Считаем total до пагинации
	query.Count(&total)

	// Сортировка
	switch filter.SortBy {
	case "popularity":
		query = query.Order("popularity DESC")
	case "newest":
		query = query.Order("created_at DESC")
	default:
		query = query.Order("name ASC")
	}

	// Пагинация
	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Preload("FaceShapes").
		Preload("HairTypes").
		Preload("Images").
		Offset(offset).
		Limit(filter.Limit).
		Find(&haircuts).Error

	return haircuts, total, err
}

func (r *haircutRepository) FindByID(ctx context.Context, id uint) (*model.Haircut, error) {
	var haircut model.Haircut
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("FaceShapes").
		Preload("HairTypes").
		Preload("Images").
		First(&haircut).Error
	if err != nil {
		return nil, err
	}
	return &haircut, nil
}

func (r *haircutRepository) Update(ctx context.Context, haircut *model.Haircut) error {
	return r.db.WithContext(ctx).Save(haircut).Error
}

func (r *haircutRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Haircut{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
