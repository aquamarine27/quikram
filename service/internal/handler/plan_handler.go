package handler

import "github.com/gofiber/fiber/v2"

type PlanFeature struct {
	Subjects       int  `json:"subjects"`
	UploadsPerMonth int `json:"uploads_per_month"`
	AISummary       bool `json:"ai_summary"`
	Compression     bool `json:"compression"`
	BasicTests      bool `json:"basic_tests"`
	AdvancedTests   bool `json:"advanced_tests"`
	Difficulty      bool `json:"difficulty"`
	Analytics       bool `json:"analytics"`
	WeakSpots       bool `json:"weak_spots"`
	Export          bool `json:"export"`
	AIChat          bool `json:"ai_chat"`
}

type Plan struct {
	ID         string      `json:"id"`
	Price      int         `json:"price"`
	Period     string      `json:"period"`
	Badge      string      `json:"badge"`
	Highlighted bool       `json:"highlighted"`
	Features   PlanFeature `json:"features"`
}

type PlanHandler struct{}

func NewPlanHandler() *PlanHandler {
	return &PlanHandler{}
}

func (h *PlanHandler) List(c *fiber.Ctx) error {
	plans := []Plan{
		{
			ID:          "free",
			Price:       0,
			Period:      "forever",
			Badge:       "",
			Highlighted: false,
			Features: PlanFeature{
				Subjects:        3,
				UploadsPerMonth: 15,
				AISummary:       true,
				Compression:     false,
				BasicTests:      true,
				AdvancedTests:   false,
				Difficulty:      false,
				Analytics:       false,
				WeakSpots:       false,
				Export:          false,
				AIChat:          false,
			},
		},
		{
			ID:          "pro",
			Price:       49900,
			Period:      "month",
			Badge:       "popular",
			Highlighted: true,
			Features: PlanFeature{
				Subjects:        -1,
				UploadsPerMonth: -1,
				AISummary:       true,
				Compression:     true,
				BasicTests:      true,
				AdvancedTests:   true,
				Difficulty:      true,
				Analytics:       true,
				WeakSpots:       true,
				Export:          true,
				AIChat:          false,
			},
		},
		{
			ID:          "proai",
			Price:       59900,
			Period:      "month",
			Badge:       "ai_chat",
			Highlighted: false,
			Features: PlanFeature{
				Subjects:        -1,
				UploadsPerMonth: -1,
				AISummary:       true,
				Compression:     true,
				BasicTests:      true,
				AdvancedTests:   true,
				Difficulty:      true,
				Analytics:       true,
				WeakSpots:       true,
				Export:          true,
				AIChat:          true,
			},
		},
	}

	return c.JSON(plans)
}
