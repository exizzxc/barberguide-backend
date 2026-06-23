package service

import (
	"context"
	"errors"

	"github.com/exizzxc/barberguide-backend/internal/dto"
	"github.com/exizzxc/barberguide-backend/internal/model"
	"github.com/exizzxc/barberguide-backend/internal/repository"
	"gorm.io/gorm"
)

type HaircutService interface {
	Create(ctx context.Context, req dto.CreateHaircutRequest) (*dto.HaircutResponse, error)
	GetAll(ctx context.Context, filter dto.HaircutFilterRequest) (*dto.HaircutListResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.HaircutResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateHaircutRequest) (*dto.HaircutResponse, error)
	Delete(ctx context.Context, id uint) error
}

type haircutService struct {
	haircutRepo repository.HaircutRepository
}

func NewHaircutService(haircutRepo repository.HaircutRepository) HaircutService {
	return &haircutService{haircutRepo: haircutRepo}
}

func (s *haircutService) Create(ctx context.Context, req dto.CreateHaircutRequest) (*dto.HaircutResponse, error) {
	haircut := &model.Haircut{
		Name:        req.Name,
		Description: req.Description,
		Style:       req.Style,
		Length:      req.Length,
	}

	for _, fs := range req.FaceShapes {
		haircut.FaceShapes = append(haircut.FaceShapes, model.HaircutFaceShape{FaceShape: fs})
	}

	for _, ht := range req.HairTypes {
		haircut.HairTypes = append(haircut.HairTypes, model.HaircutHairType{HairType: ht})
	}

	if err := s.haircutRepo.Create(ctx, haircut); err != nil {
		return nil, err
	}

	return toHaircutResponse(haircut), nil
}

func (s *haircutService) GetAll(ctx context.Context, filter dto.HaircutFilterRequest) (*dto.HaircutListResponse, error) {
	haircuts, total, err := s.haircutRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	var responses []dto.HaircutResponse
	for _, h := range haircuts {
		responses = append(responses, *toHaircutResponse(&h))
	}

	return &dto.HaircutListResponse{
		Data:  responses,
		Total: total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}, nil
}

func (s *haircutService) GetByID(ctx context.Context, id uint) (*dto.HaircutResponse, error) {
	haircut, err := s.haircutRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("haircut not found")
		}
		return nil, err
	}

	return toHaircutResponse(haircut), nil
}

func (s *haircutService) Update(ctx context.Context, id uint, req dto.UpdateHaircutRequest) (*dto.HaircutResponse, error) {
	haircut, err := s.haircutRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("haircut not found")
		}
		return nil, err
	}

	if req.Name != "" {
		haircut.Name = req.Name
	}
	if req.Description != "" {
		haircut.Description = req.Description
	}
	if req.Style != "" {
		haircut.Style = req.Style
	}
	if req.Length != "" {
		haircut.Length = req.Length
	}

	if err := s.haircutRepo.Update(ctx, haircut); err != nil {
		return nil, err
	}

	return toHaircutResponse(haircut), nil
}

func (s *haircutService) Delete(ctx context.Context, id uint) error {
	_, err := s.haircutRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("haircut not found")
		}
		return err
	}

	return s.haircutRepo.Delete(ctx, id)
}

func toHaircutResponse(h *model.Haircut) *dto.HaircutResponse {
	resp := &dto.HaircutResponse{
		ID:          h.ID,
		Name:        h.Name,
		Description: h.Description,
		Style:       h.Style,
		Length:      h.Length,
		Popularity:  h.Popularity,
		FaceShapes:  []string{},
		HairTypes:   []string{},
		Images:      []dto.ImageResponse{},
	}

	for _, fs := range h.FaceShapes {
		resp.FaceShapes = append(resp.FaceShapes, fs.FaceShape)
	}

	for _, ht := range h.HairTypes {
		resp.HairTypes = append(resp.HairTypes, ht.HairType)
	}

	for _, img := range h.Images {
		resp.Images = append(resp.Images, dto.ImageResponse{
			ID:     img.ID,
			URL:    img.URL,
			Angle:  img.Angle,
			IsMain: img.IsMain,
		})
	}

	return resp
}
