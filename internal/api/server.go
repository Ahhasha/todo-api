package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

type ServerInterface interface {
	ListLists(ctx echo.Context, params ListListsParams) error

	CreateList(ctx echo.Context) error

	DeleteList(ctx echo.Context, id Id) error

	GetList(ctx echo.Context, id Id) error

	UpdateList(ctx echo.Context, id Id) error

	GetTasks(ctx echo.Context, listID string, params GetTasksParams) error

	CreateTask(ctx echo.Context, listID string) error

	DeleteTask(ctx echo.Context, taskID string) error

	GetTask(ctx echo.Context, taskID string) error

	UpdateTask(ctx echo.Context, taskID string) error
}

type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

func (w *ServerInterfaceWrapper) ListLists(ctx echo.Context) error {
	var err error

	var params ListListsParams

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	err = w.Handler.ListLists(ctx, params)
	return err
}

func (w *ServerInterfaceWrapper) CreateList(ctx echo.Context) error {
	var err error

	err = w.Handler.CreateList(ctx)
	return err
}

func (w *ServerInterfaceWrapper) DeleteList(ctx echo.Context) error {
	var err error
	var id Id

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	err = w.Handler.DeleteList(ctx, id)
	return err
}

func (w *ServerInterfaceWrapper) GetList(ctx echo.Context) error {
	var err error
	var id Id

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	err = w.Handler.GetList(ctx, id)
	return err
}

func (w *ServerInterfaceWrapper) UpdateList(ctx echo.Context) error {
	var err error
	var id Id

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	err = w.Handler.UpdateList(ctx, id)
	return err
}

func (w *ServerInterfaceWrapper) GetTasks(ctx echo.Context) error {
	var err error
	var listID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "listID", runtime.ParamLocationPath, ctx.Param("listID"), &listID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter listID: %s", err))
	}

	var params GetTasksParams

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	err = w.Handler.GetTasks(ctx, listID, params)
	return err
}

func (w *ServerInterfaceWrapper) CreateTask(ctx echo.Context) error {
	var err error
	var listID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "listID", runtime.ParamLocationPath, ctx.Param("listID"), &listID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter listID: %s", err))
	}

	err = w.Handler.CreateTask(ctx, listID)
	return err
}

func (w *ServerInterfaceWrapper) DeleteTask(ctx echo.Context) error {
	var err error
	var taskID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "taskID", runtime.ParamLocationPath, ctx.Param("taskID"), &taskID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter taskID: %s", err))
	}

	err = w.Handler.DeleteTask(ctx, taskID)
	return err
}

func (w *ServerInterfaceWrapper) GetTask(ctx echo.Context) error {
	var err error
	var taskID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "taskID", runtime.ParamLocationPath, ctx.Param("taskID"), &taskID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter taskID: %s", err))
	}

	err = w.Handler.GetTask(ctx, taskID)
	return err
}

func (w *ServerInterfaceWrapper) UpdateTask(ctx echo.Context) error {
	var err error
	var taskID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "taskID", runtime.ParamLocationPath, ctx.Param("taskID"), &taskID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter taskID: %s", err))
	}

	err = w.Handler.UpdateTask(ctx, taskID)
	return err
}

type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/v1/lists", wrapper.ListLists)
	router.POST(baseURL+"/api/v1/lists", wrapper.CreateList)
	router.DELETE(baseURL+"/api/v1/lists/:id", wrapper.DeleteList)
	router.GET(baseURL+"/api/v1/lists/:id", wrapper.GetList)
	router.PATCH(baseURL+"/api/v1/lists/:id", wrapper.UpdateList)
	router.GET(baseURL+"/api/v1/lists/:listID/tasks", wrapper.GetTasks)
	router.POST(baseURL+"/api/v1/lists/:listID/tasks", wrapper.CreateTask)
	router.DELETE(baseURL+"/api/v1/tasks/:taskID", wrapper.DeleteTask)
	router.GET(baseURL+"/api/v1/tasks/:taskID", wrapper.GetTask)
	router.PATCH(baseURL+"/api/v1/tasks/:taskID", wrapper.UpdateTask)

}
