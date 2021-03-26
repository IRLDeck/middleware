package middleware

import (
	"context"

	ut "github.com/go-playground/universal-translator"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type LocaleContextKey string

const (
	Translator LocaleContextKey = "locale.translator"
	Locale     LocaleContextKey = "locale.locale"
)

type LocaleServer interface {
	UniversalTranslator() *ut.UniversalTranslator
	LanguageMatcher() language.Matcher
}

func LocaleMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		server, ok := info.Server.(LocaleServer)
		if !ok {
			return handler(ctx, req)
		}

		headers, _ := metadata.FromIncomingContext(ctx)

		locale := "en"
		acceptLanguage := headers["grpcgateway-accept-language"]
		if len(acceptLanguage) > 0 {
			t, _ := language.MatchStrings(server.LanguageMatcher(), acceptLanguage[0])
			locale = t.Parent().String()
		}

		t, _ := server.UniversalTranslator().GetTranslator(locale)
		ctx = context.WithValue(ctx, Translator, t)
		ctx = context.WithValue(ctx, Locale, locale)
		return handler(ctx, req)
	}
}
