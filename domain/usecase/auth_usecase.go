package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/josephakayesi/go-cerbos-abac/application/dto"
	"github.com/josephakayesi/go-cerbos-abac/application/validation"

	entity "github.com/josephakayesi/go-cerbos-abac/domain/core/entity"
	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
	"github.com/josephakayesi/go-cerbos-abac/domain/repository"
	"github.com/josephakayesi/go-cerbos-abac/internal"
)

type AuthUsecase interface {
	RegisterUser(c context.Context, r dto.RegisterUserDTO) (*entity.User, error)
	LoginUser(c context.Context, l dto.LoginUserDto) (*entity.User, error)
	FindUserByEmailOrUsername(c context.Context, u vo.UserCredentials) (*entity.User, error)
}

type authUsecase struct {
	userRepository repository.UserRepository
	contextTimeout time.Duration
}

func NewAuthUsecase(userRepository repository.UserRepository, timeout time.Duration) AuthUsecase {
	return &authUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (uu *authUsecase) FindUserByEmailOrUsername(c context.Context, uc vo.UserCredentials) (*entity.User, error) {
	_, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	user, err := uu.userRepository.FindByEmailOrUsername(c, uc)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *authUsecase) RegisterUser(c context.Context, r dto.RegisterUserDTO) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	user := &entity.User{}

	err := validation.RegisterUserDTOValidation(&r)

	if err != nil {
		return nil, err
	}

	user.ID = internal.GenerateUserId()
	user.FirstName = r.FirstName
	user.LastName = r.LastName
	user.Username = r.Username
	user.Email = r.Email
	user.Password = r.Password
	user.Status = internal.Unverified
	user.Role = vo.UserRole

	_, err = uu.userRepository.Create(ctx, *user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uu *authUsecase) LoginUser(c context.Context, l dto.LoginUserDto) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
	defer cancel()

	err := validation.LoginUserDTOValidation(&l)

	if err != nil {
		return nil, err
	}

	userCredentials := vo.UserCredentials{Email: l.UsernameOrEmail, Username: l.UsernameOrEmail}

	user, _ := uu.FindUserByEmailOrUsername(ctx, userCredentials)

	if user == nil {
		return nil, errors.New("sorry. your login details are incorrect")
	}

	doesPasswordMatch := user.Password.DoesPasswordMatch(l.Password.String())

	if !doesPasswordMatch {
		return nil, errors.New("sorry. your login details are incorrect")
	}

	return user, nil
}
