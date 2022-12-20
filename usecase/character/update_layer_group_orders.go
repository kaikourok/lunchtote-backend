package character

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/library/slice"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *CharacterUsecase) UpdateLayerGroupOrders(characterId, groupId int, orders *[]model.CharacterIconLayerGroupOrder) error {
	logger := s.registry.GetLogger()
	repository := s.registry.GetRepository()

	err := validation.Validate(groupId, validation.Min(1))
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	if orders == nil {
		return usecaseErrors.ErrValidate
	}

	for _, order := range *orders {
		err := validation.ValidateStruct(&order,
			validation.Field(&order.LayerGroup, validation.Min(1)),
			validation.Field(&order.Order, validation.Min(0), validation.Max(10000)),
		)
		if err != nil {
			return usecaseErrors.ErrValidate
		}
	}

	payloadLayerGroups := make([]int, len(*orders))
	payloadOrders := make([]int, len(*orders))

	for i, v := range *orders {
		payloadLayerGroups[i] = v.LayerGroup
		payloadOrders[i] = v.Order
	}

	if slice.IsContainsDuplicateValue(&payloadLayerGroups) {
		return usecaseErrors.ErrValidate
	}

	if slice.IsContainsDuplicateValue(&payloadOrders) {
		return usecaseErrors.ErrValidate
	}

	err = repository.UpdateLayerGroupOrders(characterId, groupId, orders)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
