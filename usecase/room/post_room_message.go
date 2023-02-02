package room

import (
	"errors"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/entity/service"
	"github.com/kaikourok/lunchtote-backend/entity/validator"
	usecaseErrors "github.com/kaikourok/lunchtote-backend/usecase/errors"
)

func (s *RoomUsecase) PostRoomMessage(characterId int, message *model.RoomPostMessage) error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()
	notificator := s.registry.GetNotificator()

	if message == nil {
		return usecaseErrors.ErrValidate
	}

	err := validation.ValidateStruct(message,
		validation.Field(&message.Name, validator.IsNotContainSpecialRune),
		validation.Field(&message.Message, validator.IsNotOnlySpace, validator.IsNotContainSpecialRune),
		validation.Field(&message.ReplyPermission, validation.In("DISALLOW", "FOLLOW", "FOLLOWED", "MUTUAL_FOLLOW", "ALL")),
	)
	if err != nil {
		return usecaseErrors.ErrValidate
	}

	permissions, roleType, banned, err := repository.RetrieveRoomOwnPermissions(characterId, message.Room)
	if err != nil {
		return err
	}

	if banned || !permissions.Write {
		return errors.New("指定のルームで発言を行う権限がありません")
	}

	if (message.Refer != nil || message.DirectReply != nil) && !permissions.UseReply {
		return errors.New("指定のルームで返信を行う権限がありません")
	}

	if message.Secret && !permissions.UseSecret {
		return errors.New("指定のルームでは秘話を使用できません")
	}

	if roleType != "MASTER" && roleType != "MEMBER" {
		room, targetName, newMemberWebhooks, err := repository.JoinToRoom(characterId, message.Room)
		if err != nil {
			logger.Error(err)
			return err
		}

		if 0 < len(newMemberWebhooks) {
			go func() {
				newMemberReplacer := strings.NewReplacer(
					"{base-path}", config.GetString("general.client-host"),
					"{room-title}", room.Title,
					"{room-id}", strconv.Itoa(room.Id),
					"{entry-number-text}", service.ConvertCharacterIdToText(characterId),
					"{name}", targetName,
				)
				newMemberMessage := newMemberReplacer.Replace(config.GetString("notification.new-member-template"))

				for _, webhook := range newMemberWebhooks {
					notificator.SendWebhook(webhook, newMemberMessage)
				}
			}()
		}
	}

	messageId, err := repository.PostRoomMessage(characterId, message, config.GetString("general.upload-path"))
	if err != nil {
		logger.Error(err)
		return err
	}

	go func() {
		dto, err := repository.NotificateRoomMessage(messageId)
		if err != nil {
			logger.Error(err)
			return
		}

		if 0 < len(dto.RepliedWebhooks) {
			repliedReplacer := strings.NewReplacer(
				"{base-path}", config.GetString("general.client-host"),
				"{entry-number-text}", service.ConvertCharacterIdToText(dto.UserId),
				"{name}", dto.UserName,
				"{refer-root}", strconv.Itoa(dto.ReferRoot),
			)
			repliedMessage := repliedReplacer.Replace(config.GetString("notification.replied-template"))

			for _, webhook := range dto.RepliedWebhooks {
				notificator.SendWebhook(webhook, repliedMessage)
			}
		}

		if 0 < len(dto.SubscribeWebhooks) {
			subscribeReplacer := strings.NewReplacer(
				"{base-path}", config.GetString("general.client-host"),
				"{room-title}", dto.RoomTitle,
				"{room-id}", strconv.Itoa(dto.RoomId),
				"{entry-number-text}", service.ConvertCharacterIdToText(dto.UserId),
				"{name}", dto.UserName,
			)
			subscribeMessage := subscribeReplacer.Replace(config.GetString("notification.subscribe-template"))

			for _, webhook := range dto.SubscribeWebhooks {
				notificator.SendWebhook(webhook, subscribeMessage)
			}
		}
	}()

	return nil
}
