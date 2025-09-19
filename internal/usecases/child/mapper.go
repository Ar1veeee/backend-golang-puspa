package child

import (
	"backend-golang/internal/adapters/http/dto"
	"backend-golang/internal/domain/entities"
	"backend-golang/internal/helpers"
	"backend-golang/internal/infrastructure/config"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Mapper interface {
	ChildResponse(parentDetail *entities.ParentDetail, child *entities.Children) (*dto.ChildResponse, error)
}

type childMapper struct {
	encryptionKey string
}

func NewChildMapper() Mapper {
	key := config.GetEnv("ENCRYPTION_KEY", "")
	if key == "" {
		log.Fatal().Err(fmt.Errorf("missing encrypted key"))
	}

	return &childMapper{
		encryptionKey: key,
	}
}

func (m *childMapper) ChildResponse(parentDetail *entities.ParentDetail, child *entities.Children) (*dto.ChildResponse, error) {
	var parentPhone, parentName string

	if parentDetail != nil {
		parentName = parentDetail.ParentName
		if len(parentDetail.ParentPhone) > 0 {
			if decryptedPhone, err := helpers.DecryptData(parentDetail.ParentPhone, m.encryptionKey); err != nil {
				log.Warn().Err(err).Msg("Failed to decrypt parent phone")
				parentPhone = "[Encrypted]"
			} else {
				parentPhone = string(decryptedPhone)
			}
		}
	}

	return &dto.ChildResponse{
		ChildId:        child.Id,
		ChildName:      child.ChildName,
		ChildBirthDate: child.ChildBirthDate,
		ChildGender:    child.ChildGender,
		ParentName:     parentName,
		ParentPhone:    parentPhone,
		CreatedAt:      child.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      child.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
