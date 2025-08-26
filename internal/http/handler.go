package http

const TagHandlerName = "http.handler"

type TokenIssuer interface {
	IssueAccessToken(userID string) (string, error)
	IssueRefreshToken(userID string) (string, error)
}
