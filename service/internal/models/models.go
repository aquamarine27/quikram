package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionOption struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct,omitempty"`
}

type QuestionAnswer struct {
	QuestionID       uuid.UUID `json:"question_id"`
	SelectedOptionIDs []int     `json:"selected_option_ids"`
	IsCorrect        bool      `json:"is_correct"`
}

type WeakTopic struct {
	Topic string  `json:"topic"`
	Score float64 `json:"score"`
}

type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;" json:"id"`
	Email           string     `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash    string     `gorm:"not null" json:"-"`
	Name            string     `gorm:"size:100;default:''" json:"name,omitempty"`
	Plan            string     `gorm:"size:20;default:free" json:"plan"`
	UploadsThisMonth int       `gorm:"default:0" json:"uploads_this_month"`
	UploadsResetAt  *time.Time `json:"uploads_reset_at,omitempty"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type Subject struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title    string    `gorm:"not null;size:200" json:"title"`
	Category string    `gorm:"size:100;default:''" json:"category,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (s *Subject) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type Document struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;" json:"id"`
	SubjectID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"subject_id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Filename      string     `gorm:"not null;size:255" json:"filename"`
	FileKey       string     `gorm:"not null" json:"file_key"`
	FileSize      int        `gorm:"default:0" json:"file_size,omitempty"`
	MimeType      string     `gorm:"size:100" json:"mime_type"`
	Status        string     `gorm:"size:30;default:uploaded" json:"status"`
	ExtractedText string     `gorm:"type:text;default:''" json:"extracted_text,omitempty"`
	DeleteAt      *time.Time `json:"delete_at,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type Summary struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	DocumentID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"document_id"`
	SubjectID     uuid.UUID `gorm:"type:uuid;not null;index" json:"subject_id"`
	UserID        uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	ContentShort  string    `gorm:"type:text" json:"content_short"`
	ContentMedium string    `gorm:"type:text" json:"content_medium"`
	ContentLong   string    `gorm:"type:text" json:"content_long"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (s *Summary) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type Quiz struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	SummaryID      uuid.UUID `gorm:"type:uuid;not null;index" json:"summary_id"`
	SubjectID      uuid.UUID `gorm:"type:uuid;not null;index" json:"subject_id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Title          string    `gorm:"size:200;default:''" json:"title,omitempty"`
	QuestionsCount int       `gorm:"not null" json:"questions_count"`
	Difficulty     string    `gorm:"size:20;default:medium" json:"difficulty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

type Question struct {
	ID            uuid.UUID       `gorm:"type:uuid;primaryKey;" json:"id"`
	QuizID        uuid.UUID       `gorm:"type:uuid;not null;index" json:"quiz_id"`
	QuestionText  string          `gorm:"type:text;not null" json:"question_text"`
	QuestionType  string          `gorm:"size:30;not null" json:"question_type"`
	Options       []QuestionOption `gorm:"serializer:json;type:jsonb;not null" json:"options"`
	Explanation   string          `gorm:"type:text;default:''" json:"explanation,omitempty"`
	Topic         string          `gorm:"size:200;default:''" json:"topic,omitempty"`
	OrderIndex    int             `gorm:"default:0" json:"order_index"`
}

func (q *Question) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

type QuizAttempt struct {
	ID          uuid.UUID          `gorm:"type:uuid;primaryKey;" json:"id"`
	QuizID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"quiz_id"`
	UserID      uuid.UUID          `gorm:"type:uuid;not null;index" json:"user_id"`
	Answers     []QuestionAnswer   `gorm:"serializer:json;type:jsonb;default:'[]'" json:"answers"`
	Score       float64            `gorm:"default:0" json:"score,omitempty"`
	WeakTopics  []WeakTopic        `gorm:"serializer:json;type:jsonb;default:'[]'" json:"weak_topics,omitempty"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
	StartedAt   time.Time          `gorm:"autoCreateTime" json:"started_at"`
}

func (a *QuizAttempt) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Role      string    `gorm:"size:20;not null" json:"role"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (m *ChatMessage) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	TokenHash string    `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (t *RefreshToken) BeforeCreate(db *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
