package event

import (
	"gorm.io/gorm"

	"github.com/chat-merger/merger/server/internal/callback"
)

type Context interface {
	CBClient() callback.Client
	DB() *gorm.DB
}
