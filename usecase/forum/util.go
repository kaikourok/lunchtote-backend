package forum

import (
	"fmt"
	"time"

	"github.com/kaikourok/lunchtote-backend/library/secure"
)

const identifierRuneCandidates = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!#$%&()*+,-./:;=?@^"

func (s *ForumUsecase) generateIdentifier(ip *string) *string {
	if ip == nil {
		return nil
	}

	config := s.registry.GetConfig()
	t := time.Now()

	base := *ip
	base += fmt.Sprintf("%d/%s/%d", t.Year(), t.Month().String(), t.Day())
	base += config.GetString("secure.forum-identifier-secret")
	identifier := secure.GenerateShortHash(base, identifierRuneCandidates)
	return &identifier
}
