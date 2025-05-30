package service

import (
	"context"
	"enuma-elish/internal/school/repository"
	"enuma-elish/internal/school/service/data/request"
	"enuma-elish/internal/school/service/data/response"
	"enuma-elish/pkg/jwt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

	err = s.repository.CreateSchool(ctx, claims.User.ID, school)
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

	listSchool, err := s.repository.GetListSchool(ctx, claim.User.ID)
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

	userSchoolRole, err := s.repository.GetSchoolRoleByUserIDAndSchoolID(ctx, schoolID, claim.User.ID)
	if err != nil {
		log.Err(err).Msg("error getting school role")
		return "", err
	}

	now := time.Now()
	exp := now.Add(time.Hour * 2).Unix()
	nbf := now.Unix()
	iat := now.Unix()
	payload := jwt.Payload{
		Exp: exp,
		Iat: iat,
		Nbf: nbf,
		Iss: "genesis",
		Sub: claim.User.ID.String(),
		Aud: "genesis",
		User: jwt.User{
			ID:       claim.User.ID,
			Email:    claim.User.Email,
			SchoolID: schoolID,
			RoleID:   userSchoolRole.RoleID,
		},
	}

	token, err := jwt.GenerateToken(payload, s.config.JWT.Secret)
	if err != nil {
		log.Err(err).Msg("error generating token")
		return "", err
	}

	return token, nil
}

func (s *service) DeleteSchool(ctx context.Context, schoolID uuid.UUID) error {
	err := s.repository.DeleteSchool(ctx, schoolID)
	if err != nil {
		log.Err(err).Msg("Failed to delete school")
		return err
	}
	return nil
}
