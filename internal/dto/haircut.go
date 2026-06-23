package dto

type CreateHaircutRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Style       string   `json:"style"`
	Length      string   `json:"length"`
	FaceShapes  []string `json:"face_shapes"`
	HairTypes   []string `json:"hair_types"`
}

type UpdateHaircutRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Style       string   `json:"style"`
	Length      string   `json:"length"`
	FaceShapes  []string `json:"face_shapes"`
	HairTypes   []string `json:"hair_types"`
}

type HaircutResponse struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Style       string          `json:"style"`
	Length      string          `json:"length"`
	Popularity  int             `json:"popularity"`
	FaceShapes  []string        `json:"face_shapes"`
	HairTypes   []string        `json:"hair_types"`
	Images      []ImageResponse `json:"images"`
}

type ImageResponse struct {
	ID     uint   `json:"id"`
	URL    string `json:"url"`
	Angle  string `json:"angle"`
	IsMain bool   `json:"is_main"`
}

type HaircutListResponse struct {
	Data  []HaircutResponse `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

type HaircutFilterRequest struct {
	FaceShape string `form:"face_shape"`
	HairType  string `form:"hair_type"`
	Style     string `form:"style"`
	Length    string `form:"length"`
	SortBy    string `form:"sort_by"`
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=10"`
}
