package control

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/kaikourok/lunchtote-backend/entity/model"
	"github.com/kaikourok/lunchtote-backend/library/secure"
)

func (s *ControlUsecase) Initialize() error {
	logger := s.registry.GetLogger()
	config := s.registry.GetConfig()
	repository := s.registry.GetRepository()

	err := repository.Initialize()
	if err != nil {
		logger.Error(err)
	}

	redis := redis.NewClient(&redis.Options{
		Addr:     config.GetString("session.host") + ":" + config.GetString("session.port"),
		Password: config.GetString("session.password"),
		DB:       0,
	})

	keys, err := redis.Keys(config.GetString("session.redis-prefix") + "*").Result()
	if err != nil {
		log.Fatal(err)
	}

	if len(keys) != 0 {
		err = redis.Del(keys...).Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	var firstAdministratorCharacterId *int

	initialAdministrators := config.Get("administrator.initial-administrators")
	for _, administratorPropsInterface := range initialAdministrators.([]any) {
		administratorProps := administratorPropsInterface.([]any)

		characterId, err := strconv.Atoi(administratorProps[0].(string))
		if err != nil {
			log.Fatal(err)
		}

		if firstAdministratorCharacterId == nil {
			firstAdministratorCharacterId = &characterId
		}

		builder := administratorProps[1].(string) + config.GetString("secure.frontend-hash-salt")
		for i := 0; i < config.GetInt("secure.frontend-hash-stretch"); i++ {
			b := sha256.Sum256([]byte(builder))
			builder = hex.EncodeToString(b[:])
		}

		hashedPassword, err := secure.HashPassword(builder, config.GetInt("secure.bcrypt-cost"))
		if err != nil {
			log.Fatal(err)
		}

		err = repository.CreateAdministrator(
			characterId,
			hashedPassword,
			administratorProps[2].(string),
			administratorProps[3].(string),
			administratorProps[4].(string),
			secure.GenerateSecureRandomHex(config.GetInt("secure.notification-token-length")),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = os.MkdirAll(config.GetString("general.upload-directory"), 0777)
	if err != nil {
		log.Fatal(err)
	}

	initialForumGroups := config.GetStringSlice("forum.initial-forum-groups")
	initialForumGroupIDsMap := make(map[string]int, len(initialForumGroups))
	for _, group := range initialForumGroups {
		id, err := repository.CreateForumGroup(group)
		if err != nil {
			log.Fatal(err)
		}

		initialForumGroupIDsMap[group] = id
	}

	initialForums := config.Get("forum.initial-forums").([]interface{})
	for _, forumInterface := range initialForums {
		forum := forumInterface.([]interface{})
		forumGroupId := initialForumGroupIDsMap[forum[0].(string)]

		var forceType *string
		s := forum[4].(string)
		if s != "" {
			forceType = &s
		}

		_, err := repository.CreateForum(forumGroupId, &model.ForumCreateData{
			Title:         forum[1].(string),
			Summary:       forum[2].(string),
			Guide:         forum[3].(string),
			ForcePostType: forceType,
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	return err
}
