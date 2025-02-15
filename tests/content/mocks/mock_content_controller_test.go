package mock_interfaces

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/cmd/microservices/content/api"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/content/controller"
	models2 "github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/content/controller/models"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/content/controller/models/mapper"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/content/pkg/models"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/global"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/logger"
	bModels "github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/service/models"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/tests/content/mocks/mock_interfaces"
	"github.com/gofrs/uuid"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// Test for the "PostLikePost" route
func TestPostLikePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBehavior := mock_interfaces.NewMockContentBehavior(ctrl)
	handler := controller.New(mockBehavior)

	userId, _ := uuid.NewV4()
	userIdStr := userId.String()

	postId, _ := uuid.NewV4()
	postIdStr := postId.String()

	// Define mock behavior for liking a post
	mockBehavior.EXPECT().LikePost(gomock.Any(), userIdStr, postIdStr).Return(10, nil)

	req := httptest.NewRequest(http.MethodPost, "/post/like", strings.NewReader(`{"userId": "`+userIdStr+`", "postId": "`+postIdStr+`"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), global.UserKey, bModels.User{UserID: bModels.UserID(userIdStr)})

	// Новый запрос с обновленным контекстом
	req = req.WithContext(ctx)

	req.Context()
	//router.ServeHTTP(w, req)
	handler.PostLikePost(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"count":10`)
}

// Test for the "PostPost" route (Create Post)
func TestPostPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBehavior := mock_interfaces.NewMockContentBehavior(ctrl)
	//router := api.NewRouter(mockBehavior)
	handler := controller.New(mockBehavior)

	userId, _ := uuid.NewV4()
	userIdStr := userId.String()

	postId, _ := uuid.NewV4()
	postIdStr := postId.String()

	mockBehavior.EXPECT().CreatePost(gomock.Any(), userIdStr, "New Title", "New Content", 0).Return(postIdStr, nil)

	req := httptest.NewRequest(http.MethodPost, "/post", strings.NewReader(
		`{ "title": "New Title", "content": "New Content", "layer": 0}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Новый контекст с добавленным значением
	ctx := context.WithValue(req.Context(), global.UserKey, bModels.User{UserID: bModels.UserID(userIdStr)})

	// Новый запрос с обновленным контекстом
	req = req.WithContext(ctx)

	req.Context()
	//router.ServeHTTP(w, req)
	handler.PostPost(w, req)

	fmt.Println(w.Body)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "postId")
}

// Test for the "FeedPopularGet" route
func TestFeedPopularGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBehavior := mock_interfaces.NewMockContentBehavior(ctrl)
	router := api.NewRouter(mockBehavior)

	userId, _ := uuid.NewV4()
	userIdStr := userId.String()

	postId, _ := uuid.NewV4()
	postId1Str := postId.String()

	postId, _ = uuid.NewV4()
	postId2Str := postId.String()

	mockBehavior.EXPECT().GetPopularPosts(gomock.Any(), userIdStr, gomock.Any()).Return([]*models.Post{
		{PostID: postId1Str, Title: "Popular Post 1"},
		{PostID: postId2Str, Title: "Popular Post 2"},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/feed/popular", nil)
	w := httptest.NewRecorder()

	// Новый контекст с добавленным значением
	ctx := context.WithValue(req.Context(), global.UserKey, bModels.User{UserID: bModels.UserID(userIdStr)})

	// Новый запрос с обновленным контекстом
	req = req.WithContext(ctx)

	req.Context()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Popular Post 1")
	assert.Contains(t, w.Body.String(), "Popular Post 2")
}

// Test for the "PostsPostIdDelete" route
func TestPostsPostIdDelete(t *testing.T) {
	logger.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//mockBehavior := mock_interfaces.NewMockContentBehavior(ctrl)
	//
	////handler := controller.New(mockBehavior)
	//router := api.NewRouter(mockBehavior)

	mockBehavior := mock_interfaces.NewMockContentBehavior(ctrl)
	//router := api.NewRouter(mockBehavior)
	handler := controller.New(mockBehavior)

	userId, _ := uuid.NewV4()
	userIdStr := userId.String()

	postId, _ := uuid.NewV4()
	postIdStr := postId.String()

	mockBehavior.EXPECT().DeletePost(gomock.Any(), userIdStr, postIdStr).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/delete/post/"+postIdStr, nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	SetUserCookie(userIdStr, req)

	//handler.PostsPostIdDelete(w, req)
	req = AddValueToPath(req, controller.PathPostID, postIdStr)

	// Новый контекст с добавленным значением
	ctx := context.WithValue(req.Context(), global.UserKey, bModels.User{UserID: bModels.UserID(userIdStr)})

	// Новый запрос с обновленным контекстом
	req = req.WithContext(ctx)

	handler.PostsPostIdDelete(w, req)

	fmt.Println(w.Body)
	fmt.Println(w.Code)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// Test for the "AuthorPostAuthorIdGet" route
func TestAuthorPostAuthorIdGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBehavior := mock_interfaces.NewMockContentBehavior(ctrl)
	router := api.NewRouter(mockBehavior)

	authorId, _ := uuid.NewV4()
	authorIdStr := authorId.String()

	postId, _ := uuid.NewV4()
	postId1Str := postId.String()

	postId, _ = uuid.NewV4()
	postId2Str := postId.String()

	mockBehavior.EXPECT().GetAuthorPosts(gomock.Any(), gomock.Any(), authorIdStr, gomock.Any()).Return([]*models.Post{
		{PostID: postId1Str, Title: "Author Post 1"},
		{PostID: postId2Str, Title: "Author Post 2"},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/author/post/"+authorIdStr, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Author Post 1")
	assert.Contains(t, w.Body.String(), "Author Post 2")
}

func TestPostUpdatePost(t *testing.T) {
	logger.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBehavior := mock_interfaces.NewMockContentBehavior(ctrl)

	handler := controller.New(mockBehavior)
	//router := api.NewRouter(mockBehavior)

	userID := GenerateUUID()

	postId, _ := uuid.NewV4()
	postIdStr := postId.String()

	title := "New Title"
	content := "New Content"

	mockBehavior.EXPECT().UpdatePost(gomock.Any(), userID, postIdStr, title, content).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/post/update", strings.NewReader(
		`{"postId": "`+postIdStr+`", "title": "`+title+`", "content": "`+content+`"}`))

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	req = AddUserIDInReq(req, userID)

	//SetUserCookie(userID, req)

	handler.PostUpdatePost(w, req)
	//router.ServeHTTP(w, req)

	fmt.Println(w.Body)
	fmt.Println(w.Code)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestFeedSubscriptionsGet(t *testing.T) {
	logger.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBehavior := mock_interfaces.NewMockContentBehavior(ctrl)

	//handler := controller.New(mockBehavior)
	router := api.NewRouter(mockBehavior)

	userId, _ := uuid.NewV4()
	userIdStr := userId.String()

	postId, _ := uuid.NewV4()
	postId1Str := postId.String()

	postId, _ = uuid.NewV4()
	postId2Str := postId.String()

	offset := "1"
	limit := "7"

	opt := bModels.NewFeedOpt(offset, limit)

	mockBehavior.EXPECT().GetFeedSubscription(gomock.Any(), userIdStr, opt).Return([]*models.Post{
		{PostID: postId1Str, Title: "Sub Post 1"},
		{PostID: postId2Str, Title: "Sub Post 2"},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/feed/subscriptions", nil)
	req.Header.Set("Content-Type", "application/json")

	query := req.URL.Query()

	// Добавляем параметры
	query.Add("offset", offset)
	query.Add("limit", limit)

	// Присваиваем обновленные параметры обратно в URL
	req.URL.RawQuery = query.Encode()

	w := httptest.NewRecorder()

	SetUserCookie(userIdStr, req)

	//handler.PostsPostIdDelete(w, req)
	router.ServeHTTP(w, req)

	fmt.Println(w.Body)
	fmt.Println(w.Code)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Sub Post 1")
	assert.Contains(t, w.Body.String(), "Sub Post 2")
}

func TestPostMediaGet_Success(t *testing.T) {
	logger.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContentBehavior := mock_interfaces.NewMockContentBehavior(ctrl)

	handler := controller.New(mockContentBehavior)

	// Настраиваем мок для возврата тестовых данных
	expectedPostID := GenerateUUID()
	expectedUserID := GenerateUUID()
	expectedMedia := []*models.Media{{MediaID: GenerateUUID(), MediaURL: "path/to/media", MediaType: "img"}}
	mockContentBehavior.EXPECT().GetFile(gomock.Any(), expectedUserID, expectedPostID).Return(expectedMedia, nil)

	// Создаем поддельный запрос с параметрами
	req, err := http.NewRequest("GET", "/some-path?postID="+expectedPostID, nil)
	assert.NoError(t, err)

	// Добавляем пользователя в контекст
	req = AddUserIDInReq(req, expectedUserID)

	// Создаем ResponseRecorder для захвата ответа
	rr := httptest.NewRecorder()

	// Вызываем обработчик
	handler.PostMediaGet(rr, req)

	// Проверяем результаты
	assert.Equal(t, http.StatusOK, rr.Code)

	exceptedTransportMedias := mapper.MapCommonMediaSToControllerMedias(expectedMedia)

	ExceptedResponse := models2.MediaResponse{
		PostID:       expectedPostID,
		MediaContent: exceptedTransportMedias,
	}

	// Проверяем тело ответа
	var response models2.MediaResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedPostID, response.PostID)
	assert.Equal(t, ExceptedResponse, response)
}
