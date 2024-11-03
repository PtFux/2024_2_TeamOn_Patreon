package api

import (
	api "github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/content/controller"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/content/controller/interfaces"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/logger"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/middlewares"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(behavior interfaces.ContentBehavior) *mux.Router {
	mainRouter := mux.NewRouter().StrictSlash(true)

	authRouter := mainRouter.PathPrefix("/").Subrouter()
	router := mainRouter.PathPrefix("/").Subrouter()

	handleAuth(authRouter, behavior)
	handleOther(router, behavior)

	authRouter.Use(middlewares.HandlerAuth)
	router.Use(middlewares.AuthMiddleware)
	return mainRouter
}

func handleAuth(router *mux.Router, behavior interfaces.ContentBehavior) *mux.Router {
	op := "content.api.routers.NewRouterWithAuth"

	handler := api.New(behavior)

	var routes = Routes{
		Route{
			"PostLikePost",
			strings.ToUpper("Post"),
			"/post/like",
			handler.PostLikePost,
		},

		Route{
			"PostPost",
			strings.ToUpper("Post"),
			"/post",
			handler.PostPost,
		},

		Route{
			"PostUpdatePost",
			strings.ToUpper("Post"),
			"/post/update",
			handler.PostUpdatePost,
		},

		Route{
			"PostUploadContentPost",
			strings.ToUpper("Post"),
			"/post/upload/content",
			handler.PostUploadContentPost,
		},

		Route{
			"PostsPostIdDelete",
			strings.ToUpper("Delete"),
			"/delete/post/{postId}",
			handler.PostsPostIdDelete,
		},
	}

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		logger.StandardInfoF(op, "Registered: %s %s", route.Method, route.Pattern)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func handleOther(router *mux.Router, behavior interfaces.ContentBehavior) {
	op := "content.api.routers.NewRouterWithAuth"

	handler := api.New(behavior)

	var routes = Routes{
		Route{
			"FeedPopularGet",
			strings.ToUpper("Get"),
			"/feed/popular",
			handler.FeedPopularGet,
		},

		Route{
			"FeedSubscriptionsGet",
			strings.ToUpper("Get"),
			"/feed/subscriptions",
			handler.FeedSubscriptionsGet,
		},

		Route{
			"AuthorPostAuthorIdGet",
			strings.ToUpper("Get"),
			"/author/post/{authorId}",
			handler.AuthorPostAuthorIdGet,
		},

		Route{
			"PostMediaGet",
			strings.ToUpper("Get"),
			"/post/media",
			handler.PostMediaGet,
		},
	}

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		logger.StandardInfoF(op, "Registered: %s %s", route.Method, route.Pattern)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
}