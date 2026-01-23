package portfolio

import (
	"go-boilerplate/internal/dto"
	"go-boilerplate/pkg/response"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(
	g *echo.Group,
	usecase Usecase,
	bearerMiddleware echo.MiddlewareFunc,
) {
	handler := &Handler{
		usecase: usecase,
	}

	v1 := g.Group("/v1")
	portfolios := v1.Group("/portfolios")

	portfolios.POST("", handler.CreatePortfolio, bearerMiddleware)
	portfolios.GET("", handler.GetPortfolios, bearerMiddleware)
	portfolios.GET("/:id", handler.GetPortfolio, bearerMiddleware)
	portfolios.PUT("/:id", handler.UpdatePortfolio, bearerMiddleware)
	portfolios.DELETE("/:id", handler.DeletePortfolio, bearerMiddleware)
	portfolios.GET("/:id/summary", handler.GetPortfolioSummary, bearerMiddleware)

	// Holdings endpoints
	holdings := portfolios.Group("/:id/holdings")
	holdings.POST("", handler.AddHolding)
	holdings.DELETE("/:holdingId", handler.RemoveHolding)
}

func (h *Handler) CreatePortfolio(c *echo.Context) error {
	var req dto.CreatePortfolioRequest

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id format")
	}

	portfolio, err := h.usecase.CreatePortfolio(c.Request().Context(), userID, req)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	responseData := ToPortfolioResponse(portfolio)

	return response.Created(c, "portfolio created successfully", responseData)
}

func (h *Handler) GetPortfolios(c *echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id format")
	}

	portfolios, err := h.usecase.GetUserPortfolios(c.Request().Context(), userID)
	if err != nil {
		return response.InternalServerError(c, err.Error())
	}

	responseData := ToPortfolioListResponse(portfolios)

	return response.Success(c, "success get portfolios", responseData)
}

func (h *Handler) GetPortfolio(c *echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id format")
	}

	portfolioID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "invalid portfolio id")
	}

	portfolio, err := h.usecase.GetPortfolio(c.Request().Context(), userID, portfolioID)
	if err != nil {
		if err.Error() == "portfolio not found" || err.Error() == "unauthorized access to portfolio" {
			return response.NotFound(c, err.Error())
		}

		return response.InternalServerError(c, err.Error())
	}

	responseData := ToPortfolioResponse(portfolio)

	return response.Success(c, "success get portfolio", responseData)
}

func (h *Handler) UpdatePortfolio(c *echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id format")
	}

	portfolioID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "invalid portfolio id")
	}

	var req dto.UpdatePortfolioRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	portfolio, err := h.usecase.UpdatePortfolio(c.Request().Context(), userID, portfolioID, req)
	if err != nil {
		if err.Error() == "portfolio not found" || err.Error() == "unauthorized access to portfolio" {
			return response.NotFound(c, err.Error())
		}

		return response.InternalServerError(c, err.Error())
	}

	responseData := ToPortfolioResponse(portfolio)

	return response.Success(c, "success update portfolio", responseData)
}

func (h *Handler) DeletePortfolio(c *echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id format")
	}

	portfolioID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "invalid portfolio id")
	}

	if err := h.usecase.DeletePortfolio(c.Request().Context(), userID, portfolioID); err != nil {
		if err.Error() == "portfolio not found" || err.Error() == "unauthorized access to portfolio" {
			return response.NotFound(c, err.Error())
		}

		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "success delete portfolio", nil)
}

func (h *Handler) GetPortfolioSummary(c *echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id format")
	}

	portfolioID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "invalid portfolio id")
	}

	summary, err := h.usecase.GetPortfolioSummary(c.Request().Context(), userID, portfolioID)
	if err != nil {
		if err.Error() == "portfolio not found" || err.Error() == "unauthorized access to portfolio" {
			return response.NotFound(c, err.Error())
		}

		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "success get portfolio summary", summary)
}

func (h *Handler) AddHolding(c *echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id format")
	}

	portfolioID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "invalid portfolio id")
	}

	var req dto.AddHoldingRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "invalid request payload")
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	holding, err := h.usecase.AddHolding(c.Request().Context(), userID, portfolioID, req)
	if err != nil {
		if err.Error() == "portfolio not found" || err.Error() == "unauthorized access to portfolio" {
			return response.NotFound(c, err.Error())
		}

		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "success add holding", holding)
}

func (h *Handler) RemoveHolding(c *echo.Context) error {
	userIDStr := c.Get("user_id").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.BadRequest(c, "invalid user id format")
	}

	portfolioID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "invalid portfolio id")
	}

	holdingID, err := uuid.Parse(c.Param("holdingId"))
	if err != nil {
		return response.BadRequest(c, "invalid holding id")
	}

	if err := h.usecase.RemoveHolding(c.Request().Context(), userID, portfolioID, holdingID); err != nil {
		if err.Error() == "portfolio not found" || err.Error() == "unauthorized access to portfolio" {
			return response.NotFound(c, err.Error())
		}

		return response.InternalServerError(c, err.Error())
	}

	return response.Success(c, "success remove holding", nil)
}
