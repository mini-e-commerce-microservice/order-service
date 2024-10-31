package handler

import (
	"errors"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"net/http"
	"order-service/generated/proto/jwt_claims_proto"
	"order-service/internal/util"
	"strings"
)

func (h *handler) getUserFromBearerAuth(w http.ResponseWriter, r *http.Request, mustMerchantUser bool) (*jwt_claims_proto.JwtAuthAccessTokenClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, collection.Err(errors.New("authorization header is missing")))
		return nil, false
	}

	bearerSplit := strings.Split(authHeader, " ")
	if len(bearerSplit) != 2 {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, collection.Err(errors.New("invalid authorization header format")))
		return nil, false
	}

	if bearerSplit[0] != "Bearer" {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, collection.Err(errors.New("authorization scheme must be Bearer")))
		return nil, false
	}

	authAccessTokenJWT := &util.AuthAccessTokenClaims{}
	err := authAccessTokenJWT.ClaimsHS256(bearerSplit[1], h.jwtAccessTokenConf.Key)
	if err != nil {
		h.httpOtel.Err(w, r, http.StatusUnauthorized, collection.Err(err))
		return nil, false
	}

	if !authAccessTokenJWT.IsEmailVerified {
		h.httpOtel.Err(w, r, http.StatusForbidden, errors.New("email user must be verified"), "You must activate your email first")
		return nil, false
	}

	return authAccessTokenJWT.JwtAuthAccessTokenClaims, true
}
