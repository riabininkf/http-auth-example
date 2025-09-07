package handlers

//go:generate mockery --name TokenIssuer --output ./mocks --outpkg mocks --filename token_issuer.go --structname TokenIssuer

// TokenIssuer provides methods to issue access and refresh tokens for a specified user.
type TokenIssuer interface {
	IssueAccessToken(userID string) (string, error)
	IssueRefreshToken(userID string) (string, error)
}
