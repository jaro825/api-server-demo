//go:build !integration

package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jaro825/api-server-demo/cache"
	cacheMocks "github.com/jaro825/api-server-demo/mocks/github.com/jaro825/api-server-demo/cache"
	storeMocks "github.com/jaro825/api-server-demo/mocks/github.com/jaro825/api-server-demo/store"
	"github.com/jaro825/api-server-demo/store"
	"github.com/jaro825/api-server-demo/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserAPI_CreateUser_ok(t *testing.T) {
	t.Parallel()
	reqBody := strings.NewReader("{\"name\":\"jaro\"}")

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/api/v1/user", reqBody)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockStore := storeMocks.NewMockStore(t)
	mockStore.On("CreateUser", mock.Anything, mock.Anything).Return(nil).Once()

	mockCache := cacheMocks.NewMockCache(t)

	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	userAPI := &UserAPI{
		store:  mockStore,
		cache:  mockCache,
		logger: zerolog.Logger{},
	}

	userAPI.CreateUser(w, req)

	if want, got := http.StatusCreated, w.Result().StatusCode; want != got {
		t.Fatalf("expected %d; got %d", want, got)
	}
}

func TestUserAPI_CreateUser_db_error(t *testing.T) {
	t.Parallel()
	reqBody := strings.NewReader("{\"name\":\"jaro\"}")

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/api/v1/user", reqBody)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockStore := storeMocks.NewMockStore(t)
	mockStore.On("CreateUser", mock.Anything, mock.Anything).Return(store.ErrCouldNotCreateUser).Once()

	mockCache := cacheMocks.NewMockCache(t)

	mockCache.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	userAPI := &UserAPI{
		store:  mockStore,
		cache:  mockCache,
		logger: zerolog.Logger{},
	}

	userAPI.CreateUser(w, req)

	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected %d; got %d", want, got)
	}
}

func TestUserAPI_CreateUser_bad_payload(t *testing.T) {
	t.Parallel()
	reqBody := strings.NewReader("{\"unknown_field\":\"jaro\"}")

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/api/v1/user", reqBody)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	userAPI := &UserAPI{
		store:  storeMocks.NewMockStore(t),
		cache:  cacheMocks.NewMockCache(t),
		logger: zerolog.Logger{},
	}

	userAPI.CreateUser(w, req)

	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected %d; got %d", want, got)
	}
}

func TestUserAPI_GetUser_user_not_in_cache_db_OK(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/api/v1/user/3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1", nil)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := req.Context()

	mockStore := storeMocks.NewMockStore(t)
	uid, _ := uuid.Parse("3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1")
	mockStore.EXPECT().GetUser(ctx, "3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1").Return(&types.User{
		ID:   uid,
		Name: "jaro",
	}, nil).Once()

	mockCache := cacheMocks.NewMockCache(t)
	var cachedUser types.User
	mockCache.EXPECT().Get(ctx, "3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1", &cachedUser).Return(cache.ErrKeyNotFound).Once()

	userAPI := &UserAPI{
		store:  mockStore,
		cache:  mockCache,
		logger: zerolog.Logger{},
	}

	userAPI.GetUser(w, req)

	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected %d; got %d", want, got)
	}
}

func TestUserAPI_GetUser_user_not_found(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/api/v1/user/3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1", nil)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	ctx := req.Context()

	mockStore := storeMocks.NewMockStore(t)
	mockStore.EXPECT().GetUser(ctx, "3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1").Return(nil, store.ErrUserNotFound).Once()

	mockCache := cacheMocks.NewMockCache(t)
	var cachedUser types.User
	mockCache.EXPECT().Get(ctx, "3c5bed0e-8dfb-4ac2-aff1-8a76f2c821a1", &cachedUser).Return(cache.ErrKeyNotFound).Once()

	userAPI := &UserAPI{
		store:  mockStore,
		cache:  mockCache,
		logger: zerolog.Logger{},
	}

	userAPI.GetUser(w, req)

	if want, got := http.StatusNotFound, w.Result().StatusCode; want != got {
		t.Fatalf("expected %d; got %d", want, got)
	}
}

func TestUserAPI_GetUser_invalid_user_id(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:8080/api/v1/user/xyz", nil)

	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "xyz")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	userAPI := &UserAPI{
		store:  storeMocks.NewMockStore(t),
		cache:  cacheMocks.NewMockCache(t),
		logger: zerolog.Logger{},
	}

	userAPI.GetUser(w, req)

	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected %d; got %d", want, got)
	}
}
