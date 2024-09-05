package pipes_source

import (
	"time"

	"github.com/turbot/tailpipe-plugin-sdk/collection_state"
)

type AuditLogAPICollectionState struct {
	collection_state.CollectionStateBase

	StartTime time.Time `json:"start_time,omitempty"` // oldest record timestamp
	StartID   string    `json:"start_id,omitempty"`   // oldest record id
	EndTime   time.Time `json:"end_time,omitempty"`   // newest record timestamp
	EndID     string    `json:"end_id,omitempty"`     // newest record id

	prevTime time.Time `json:"-"`
	prevId   string    `json:"-"`

	// TODO #error we may need to add these fields in future if we need to capture gaps in data
	//ResumeToken string `json:"resume_token"`
	//Offset      int    `json:"offset"`
}

func NewAuditLogAPICollectionState() collection_state.CollectionState[*AuditLogAPISourceConfig] {
	return &AuditLogAPICollectionState{}
}

func (s *AuditLogAPICollectionState) Init(*AuditLogAPISourceConfig) error {
	return nil
}

func (s *AuditLogAPICollectionState) IsEmpty() bool {
	return s.StartTime.IsZero() // && s.EndTime.IsZero()
}

func (s *AuditLogAPICollectionState) Upsert(createdAt time.Time, id string) {
	if s.StartTime.IsZero() || createdAt.Before(s.StartTime) {
		s.StartTime = createdAt
		s.StartID = id
	}

	if s.EndTime.IsZero() || createdAt.After(s.EndTime) {
		s.EndTime = createdAt
		s.EndID = id
	}

	if createdAt.Equal(s.StartTime) && id < s.StartID {
		s.StartID = id
	}
	if createdAt.Equal(s.EndTime) && id > s.EndID {
		s.EndID = id
	}
}

// StartCollection stores the current state as previous state
func (s *AuditLogAPICollectionState) StartCollection() {
	s.prevTime = s.EndTime
	s.prevId = s.EndID
}

func (s *AuditLogAPICollectionState) ShouldCollectRow(createdAt time.Time, id string) bool {
	if !s.prevTime.IsZero() && createdAt.Equal(s.prevTime) && id == s.prevId {
		return false
	}

	return true
}
