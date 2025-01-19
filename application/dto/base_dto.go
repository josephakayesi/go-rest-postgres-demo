package dto

type GetTokenDTO struct {
	Token string `json:"token"`
}

type GetAccessTokenDTO struct {
	AccessToken string `json:"access_token"`
}

type GetRefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type GetTokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type GetPublicKeyDTO struct {
	PublicKey string `json:"public_key"`
}

func CreateGetTokenDTO(token string) *GetTokenDTO {
	return &GetTokenDTO{
		Token: token,
	}
}

func CreateGetAccessTokenDTO(token string) *GetAccessTokenDTO {
	return &GetAccessTokenDTO{
		AccessToken: token,
	}
}

func CreateGetRefreshTokenDTO(token string) *GetRefreshTokenDTO {
	return &GetRefreshTokenDTO{
		RefreshToken: token,
	}
}

func CreateGetTokensDTO(accessToken string, refreshToken string) *GetTokensDTO {
	return &GetTokensDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func CreateGetPublicKeyDTO(publicKey string) *GetPublicKeyDTO {
	return &GetPublicKeyDTO{
		PublicKey: publicKey,
	}
}

type EventConvertible interface {
	ToString() string
	FromString(event ...string) error
}
