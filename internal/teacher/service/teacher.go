package service

import (
	"context"
	"enuma-elish/internal/teacher/repository"
	"enuma-elish/internal/teacher/service/data/request"
	"enuma-elish/internal/teacher/service/data/response"
	commonError "enuma-elish/pkg/error"
	commonHttp "enuma-elish/pkg/http"
	"fmt"
	"net/smtp"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) InviteTeacher(ctx context.Context, data request.InviteTeacherRequest) error {

	var teachers []repository.User
	now := time.Now().UnixMilli()
	for _, email := range data.Emails {
		teachers = append(teachers, repository.User{
			Email:      email,
			ID:         uuid.New(),
			Name:       "",
			Password:   "",
			IsVerified: false,
			CreatedAt:  now,
			UpdatedAt:  0,
		})
	}

	err := s.repository.CreateTeachers(ctx, teachers, data.SchoolID)
	if err != nil {
		log.Error().Err(err).Msg("create teachers error")
		return err
	}

	go func([]repository.User) {
		for _, v := range teachers {
			token, err := s.repository.CreateTeacherVerifyToken(context.Background(), v.Email)
			if err != nil {
				log.Err(err).Str("email", v.Email).Msg("create teacher verify token")
			}

			link := fmt.Sprintf("%s/email-verification?email=%s&token=%s", s.config.Http.FrontendHost, v.Email, token)
			err = s.sendEmail(v.Email, link, "Email Verification")
			if err != nil {
				log.Err(err).Str("email", v.Email).Msg("send email")
			}
		}
	}(teachers)

	return nil
}

func (s *service) sendEmail(to string, msg, subject string) error {
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)
	message := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, msg))

	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)
	err := smtp.SendMail(addr, auth, s.config.SMTP.Username, []string{to}, message)
	if err != nil {
		log.Err(err).Msg("failed to send email")
		return err
	}

	return nil
}

func (s *service) VerifyTeacherEmail(ctx context.Context, data request.VerifyTeacherEmailRequest) error {
	token, err := s.repository.VerifyEmailToken(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	if token != data.Token {
		log.Err(commonError.ErrInvalidToken).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	return nil
}

func (s *service) UpdateTeacherAfterVerifyEmail(ctx context.Context, data request.UpdateTeacherAfterVerifyEmailRequest) error {

	token, err := s.repository.VerifyEmailToken(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	if token != data.Token {
		log.Err(commonError.ErrInvalidToken).Msg("Failed to verify email token")
		return commonError.ErrInvalidToken
	}

	teacher, err := s.repository.GetTeacherByEmail(ctx, data.Email)
	if err != nil {
		log.Err(err).Msg("Failed to get teacher")
		return commonError.ErrUserNotFound
	}

	hashPass, err := hashPassword(data.Password)
	if err != nil {
		log.Err(err).Msg("Failed to hash password")
		return err
	}

	teacher = &repository.User{
		ID:         teacher.ID,
		Name:       data.Name,
		Email:      data.Email,
		Password:   hashPass,
		IsVerified: true,
		UpdatedAt:  time.Now().UnixMilli(),
	}

	err = s.repository.UpdateTeacher(ctx, *teacher)
	if err != nil {
		log.Err(err).Msg("Failed to update teacher")
		return commonError.ErrInternal
	}

	s.repository.Redis().Del(ctx, repository.TeacherVerifyEmailTokenKey+":"+data.Email)
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *service) ListTeachers(ctx context.Context, httpQuery request.GetListTeacherQuery) (response.GetListTeacherResponse, *commonHttp.Meta, error) {

	listTeacher, total, err := s.repository.GetListTeachers(ctx, httpQuery)
	if err != nil {
		log.Err(err).Msg("list teachers")
		return nil, nil, err
	}

	meta := commonHttp.NewMetaFromQuery(httpQuery, total)
	res := make(response.GetListTeacherResponse, len(listTeacher))
	for i, teacher := range listTeacher {
		res[i] = response.GetTeacherResponse{
			ID:         teacher.ID,
			Name:       teacher.Name,
			Email:      teacher.Email,
			IsVerified: teacher.IsVerified,
			CreateAt:   teacher.CreatedAt,
			UpdateAt:   teacher.UpdatedAt,
		}
	}

	return res, meta, nil
}

func (s *service) DeleteTeacher(ctx context.Context, teacherID uuid.UUID, schoolID uuid.UUID) error {
	err := s.repository.DeleteTeacher(ctx, teacherID, schoolID)
	if err != nil {
		log.Err(err).Msg("Failed to delete teacher")
		return err
	}
	return nil
}

func (s *service) GetDetailTeacher(ctx context.Context, teacherID uuid.UUID) (response.GetDetailTeacherResponse, error) {
	teacher, err := s.repository.GetTeacherByID(ctx, teacherID)
	if err != nil {
		log.Err(err).Msg("Failed to get teacher")
		return response.GetDetailTeacherResponse{}, err
	}

	res := response.GetDetailTeacherResponse{
		ID:         teacher.ID,
		Name:       teacher.Name,
		Email:      teacher.Email,
		IsVerified: teacher.IsVerified,
		CreateAt:   teacher.CreatedAt,
		UpdateAt:   teacher.UpdatedAt,
	}

	return res, nil
}

func (s *service) UpdateTeacherClass(ctx context.Context, data request.UpdateTeacherClassRequest) error {
	err := s.repository.UpdateTeacherClass(ctx, data.TeacherID, data.OldClassID, data.NewClassID)
	if err != nil {
		log.Err(err).Msg("Failed to update teacher class assignment")
		return err
	}
	return nil
}

func (s *service) GetTeacherSubjects(ctx context.Context, teacherID uuid.UUID) ([]response.Subject, error) {
	subjects, err := s.repository.GetTeacherSubjects(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	res := make([]response.Subject, len(subjects))
	for i, subject := range subjects {
		res[i] = response.Subject{
			ID:        subject.ID,
			Name:      subject.Name,
			SchoolID:  subject.SchoolID,
			CreatedAt: subject.CreatedAt,
			UpdatedAt: subject.UpdatedAt,
		}
	}
	return res, nil
}
