package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jaro825/api-server-demo/cache"
	"github.com/jaro825/api-server-demo/store"
	"github.com/jaro825/api-server-demo/types"
	"github.com/rs/zerolog"
	"net/http"
)

var ErrUserNotFound = errors.New("user not found")

type UserAPI struct {
	store  store.Store
	cache  cache.Cache
	logger zerolog.Logger
}

func NewUserAPI(l zerolog.Logger, s store.Store, c cache.Cache) *UserAPI {
	return &UserAPI{s, c, l}
}

func (u *UserAPI) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")

		// check if id is a valid UUID
		id, err := uuid.Parse(userID)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, Error{Error: "user ID is not a valid UUID"})

			return
		}

		ctx := r.Context()

		// check local cache
		var cachedUser types.User
		err = u.cache.Get(ctx, id.String(), &cachedUser)
		if err == nil {
			u.logger.Debug().Msgf("UserAPI.UserCTX user found in cache: %+v", cachedUser)
			ctx = context.WithValue(ctx, contextKeyUser, &cachedUser)
			next.ServeHTTP(w, r.WithContext(ctx))

			return
		}

		// query db
		user, err := u.store.GetUser(ctx, id.String())
		if err != nil {
			WriteJSON(w, http.StatusNotFound, Error{Error: err.Error()})

			return
		}

		ctx = context.WithValue(ctx, contextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type key int

const (
	contextKeyUser key = iota
)

func (u *UserAPI) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, ok := ctx.Value(contextKeyUser).(*types.User)
	if !ok {
		WriteJSON(w, http.StatusNotFound, Error{Error: ErrUserNotFound.Error()})

		return
	}

	WriteJSON(w, http.StatusOK, user)
}

func (u *UserAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, Error{Error: "could not decode request payload into a user"})

		return
	}

	user.ID = uuid.New()

	ctx := r.Context()

	err = u.cache.Set(ctx, user.ID.String(), &user)
	if err != nil {
		u.logger.Warn().Msgf("UserAPI.CreateUser could not save user in cache: %v", err)
	} else {
		u.logger.Debug().Msgf("UserAPI.CreateUser saved user in cache: %+v", user)
	}

	err = u.store.CreateUser(ctx, &user)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, Error{Error: err.Error()})

		return
	}

	WriteJSON(w, http.StatusCreated, &user)
}

func (u *UserAPI) UpdateUser(_ http.ResponseWriter, _ *http.Request) {
	panic("UserAPI.UpdateUser not implemented") // todo implement
}

func (u *UserAPI) DeleteUser(_ http.ResponseWriter, _ *http.Request) {
	panic("UserAPI.DeleteUser not implemented") // todo implement
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(v) //nolint: errchkjson
}

type Error struct {
	Error string `json:"error"`
}
