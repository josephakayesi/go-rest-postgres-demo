package controller

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/josephakayesi/go-cerbos-abac/application/dto"
	vo "github.com/josephakayesi/go-cerbos-abac/domain/core/valueobject"
	"github.com/josephakayesi/go-cerbos-abac/domain/usecase"
	"github.com/josephakayesi/go-cerbos-abac/infra/config"
	internal "github.com/josephakayesi/go-cerbos-abac/internal"
	"golang.org/x/exp/slog"
)

type AuthController struct {
	Queue       *config.NatsQueue
	AuthUsecase usecase.AuthUsecase
	Logger      slog.Logger
}

func (tc *AuthController) Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.UserContext())
	ctx = internal.SetBrowserFingerprintInContext(ctx, c)
	logId, ctx := internal.SetLogIdInContext(ctx)

	defer cancel()

	registerUserDTO := &dto.RegisterUserDTO{}

	if err := c.BodyParser(&registerUserDTO); err != nil {
		tc.Logger.Error("unable to parse CreateAuthDTO", "error", err)
		return err
	}

	tc.Logger.Info("User attempting to register", "email", registerUserDTO.Email)

	userCredentials := vo.UserCredentials{Email: registerUserDTO.Email, Username: registerUserDTO.Username}

	existingUser, _ := tc.AuthUsecase.FindUserByEmailOrUsername(ctx, userCredentials)

	if existingUser != nil {
		tc.Logger.Info("User with email or username already exists", "log_id", logId, "email", registerUserDTO.Email, "username", registerUserDTO.Username)
		return c.Status(400).JSON(internal.NewErrorResponse("username or email already exists"))
	}

	_, err := tc.AuthUsecase.RegisterUser(ctx, *registerUserDTO)

	if err != nil {
		tc.Logger.Error("AuthUsecase unable to register user", "log_id", logId, "err", err)
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	tc.Logger.Info("User with email and username registered successfully", "log_id", logId, "email", registerUserDTO.Email, "username", registerUserDTO.Username)

	return c.Status(201).JSON(*internal.NewSuccessResponse("registered user successfully. please proceed to login"))
}

func (tc *AuthController) Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.UserContext())
	logId, ctx := internal.SetLogIdInContext(ctx)

	defer cancel()

	loginUserDto := &dto.LoginUserDto{}

	if err := c.BodyParser(&loginUserDto); err != nil {
		tc.Logger.Error("unable to parse LoginUserDto", slog.String("log_id", logId), slog.String("error", err.Error()))
		return err
	}

	tc.Logger.Info("user attempting to login", "log_id", logId, "username_or_email", loginUserDto.UsernameOrEmail)

	user, err := tc.AuthUsecase.LoginUser(ctx, *loginUserDto)

	if err != nil {
		tc.Logger.Error("AuthUsecase.Login unable to login user", slog.String("log_id", logId), slog.String("error", err.Error()))
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	accessToken, _, err := user.CreateAccessToken()

	if err != nil {
		tc.Logger.Error("cannot create access token", slog.String("log_id", logId), slog.String("error", err.Error()))
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	refreshToken, _, err := user.CreateRefreshToken()

	if err != nil {
		tc.Logger.Error("cannot create access token", slog.String("log_id", logId), slog.String("error", err.Error()))
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	tokens := dto.CreateGetTokensDTO(accessToken, refreshToken)

	bf := internal.BrowserFingerprint{
		ClientIP:  c.Context().RemoteIP().String(),
		UserAgent: string(c.Context().UserAgent()),
	}

	key := internal.ContextWithValueKey("BrowserFingerprint")

	_ = context.WithValue(ctx, key, bf)

	tc.Logger.Info("user with email or username logged in successfully", slog.String("log_id", logId), slog.String("username_or_email", loginUserDto.UsernameOrEmail))

	return c.Status(201).JSON(internal.NewSuccessResponse("logged in successfully", internal.WithData(tokens)))
}

func (tc *AuthController) RefreshToken(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.UserContext())
	logId, ctx := internal.SetLogIdInContext(ctx)

	defer cancel()

	refreshTokenDto := &dto.RefreshTokenDto{}

	if err := c.BodyParser(&refreshTokenDto); err != nil {
		tc.Logger.Error("unable to parse RefreshTokenDto", slog.String("log_id", logId), slog.String("error", err.Error()))
		return err
	}

	tc.Logger.Info("User attempting to refresh token", slog.String("log_id", logId))

	payload, err := internal.NewPaseto().VerifyRefreshToken(refreshTokenDto.RefreshToken)
	if err != nil {
		tc.Logger.Error("refresh token verification failed", slog.String("log_id", logId), slog.String("error", err.Error()))
		return c.Status(400).JSON(internal.NewErrorResponse("we couldn't verify your token. please try again"))
	}

	userCredentials := vo.UserCredentials{Email: payload.Email, Username: payload.Username}

	user, err := tc.AuthUsecase.FindUserByEmailOrUsername(ctx, userCredentials)
	if err != nil {
		tc.Logger.Error("could not find user with credentials",
			slog.Group("credentials", slog.String("email", userCredentials.Email), slog.String("username", userCredentials.Username)),
			slog.String("log_id", logId), slog.String("error", err.Error()))

		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	accessToken, _, err := user.CreateAccessToken()
	if err != nil {
		tc.Logger.Error("cannot create access token", slog.String("log_id", logId))
		return c.Status(400).JSON(internal.NewErrorResponse(err.Error()))
	}

	bf := internal.BrowserFingerprint{
		ClientIP:  c.Context().RemoteIP().String(),
		UserAgent: string(c.Context().UserAgent()),
	}

	key := internal.ContextWithValueKey("BrowserFingerprint")

	_ = context.WithValue(ctx, key, bf)

	token := dto.CreateGetAccessTokenDTO(accessToken)

	tc.Logger.Info("user with email and username refreshed token successfully", slog.String("log_id", logId), slog.Group("credentials", slog.String("email", userCredentials.Email), slog.String("username", userCredentials.Username)))

	return c.Status(201).JSON(internal.NewSuccessResponse("refreshed token successfully", internal.WithData(token)))
}

func (tc *AuthController) GetVerificationPublicKeys(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.UserContext())
	logId, _ := internal.SetLogIdInContext(ctx)

	defer cancel()

	bf := internal.BrowserFingerprint{
		ClientIP:  c.Context().RemoteIP().String(),
		UserAgent: string(c.Context().UserAgent()),
	}

	tc.Logger.Info("attempting to get verifcation public keys", slog.String("log_id", logId), slog.String("user_agent", bf.UserAgent), slog.String("client_ip", bf.ClientIP))

	pub := dto.CreateGetPublicKeyDTO(config.Get("PASETO_PRIVATE_KEY_SECRET", ""))

	tc.Logger.Info("got paseto verification key successfully", slog.String("log_id", logId), slog.String("user_agent", bf.UserAgent), slog.String("client_ip", bf.ClientIP))

	return c.Status(201).JSON(internal.NewSuccessResponse("got verification public keys sucessfully", internal.WithData(pub)))
}
