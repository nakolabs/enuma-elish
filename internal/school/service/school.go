package service

import (
	"context"
	"enuma-elish/internal/school/repository"
	"enuma-elish/internal/school/service/data/request"
	"enuma-elish/internal/school/service/data/response"
	"enuma-elish/pkg/jwt"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"time"
)

func (s *service) CreatSchool(ctx context.Context, data request.CreateSchoolRequest) error {

	now := time.Now().UnixMilli()
	school := repository.School{
		ID:        uuid.New(),
		Name:      data.Name,
		Level:     data.Level,
		CreatedAt: now,
		UpdatedAt: 0,
	}

	claims, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return err
	}
	fmt.Println(claims)

	userID, ok := claims["id"].(string)
	if !ok {
		log.Error().Msg("Failed to extract user id")
		return errors.New("error extracting user id")
	}

	err = s.repository.CreateSchool(ctx, uuid.MustParse(userID), school)
	if err != nil {
		log.Err(err).Msg("error creating school")
		return err
	}

	return nil
}

func (s *service) GetDetailSchool(ctx context.Context, schoolID uuid.UUID) (response.DetailSchool, error) {
	var detailSchool response.DetailSchool
	school, err := s.repository.GetSchoolByID(ctx, schoolID)
	if err != nil {
		log.Err(err).Msg("error getting school")
		return detailSchool, err
	}

	detailSchool = response.DetailSchool{
		ID:        school.ID,
		Name:      school.Name,
		Level:     school.Level,
		CreatedAt: school.CreatedAt,
		UpdatedAt: school.UpdatedAt,
	}

	return detailSchool, nil
}

func (s *service) GetListSchool(ctx context.Context) (response.ListSchool, error) {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return nil, err
	}

	userIdString, ok := claim["user_id"].(string)
	if !ok {
		log.Error().Msg("Failed to extract user id")
		return nil, errors.New("error extracting user id")
	}

	userID, err := uuid.Parse(userIdString)
	if err != nil {
		log.Err(err).Msg("Failed to extract user id")
		return nil, err
	}

	listSchool, err := s.repository.GetListSchool(ctx, userID)
	if err != nil {
		log.Err(err).Msg("error getting list school")
		return nil, err
	}

	res := response.ListSchool{}
	for _, school := range listSchool {
		res = append(res, response.DetailSchool{
			ID:        school.ID,
			Name:      school.Name,
			Level:     school.Level,
			CreatedAt: school.CreatedAt,
			UpdatedAt: school.UpdatedAt,
		})
	}

	return res, nil
}

func (s *service) SwitchSchool(ctx context.Context, schoolID uuid.UUID) (string, error) {
	claim, err := jwt.ExtractContext(ctx)
	if err != nil {
		log.Err(err).Msg("Failed to extract claims")
		return "", err
	}

	userIdString, ok := claim["user_id"].(string)
	if !ok {
		log.Error().Msg("Failed to extract user id")
		return "", errors.New("error extracting user id")
	}

	userID, err := uuid.Parse(userIdString)
	if err != nil {
		log.Err(err).Msg("error parsing user id")
		return "", err
	}

	email, ok := claim["email"].(string)
	if !ok {
		log.Error().Msg("Failed to extract email")
		return "", errors.New("error extracting email")
	}

	userSchoolRole, err := s.repository.GetSchoolRoleByUserIDAndSchoolID(ctx, schoolID, userID)
	if err != nil {
		log.Err(err).Msg("error getting school role")
		return "", err
	}

	jwtPayload := map[string]interface{}{
		"email":     email,
		"user_id":   userID,
		"role_id":   userSchoolRole.RoleID,
		"school_id": userSchoolRole.SchoolID,
	}

	token, err := jwt.GenerateToken(time.Duration(s.config.JWT.Duration)*time.Hour, jwtPayload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Msg("error generating token")
		return "", err
	}

	return token, nil
}
