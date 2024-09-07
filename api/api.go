package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/orenoid/telegram-account-bot/service/bill"
	"github.com/orenoid/telegram-account-bot/service/user"
)

func GetEcho(hub *ControllersHub) *echo.Echo {
	e := echo.New()

	// register controllers
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/openapi/bills", hub.importBills)

	return e
}

type ControllersHub struct {
	userService *user.Service
	billService *bill.Service
}

func NewControllersHub(userService *user.Service, billService *bill.Service) *ControllersHub {
	return &ControllersHub{userService: userService, billService: billService}
}

type ImportBillsRequest struct {
	Bills []struct {
		Amount    float64    `json:"amount"`
		Category  string     `json:"category"`
		Name      *string    `json:"name,omitempty"`      // optional
		CreatedAt *time.Time `json:"createdAt,omitempty"` // if not provided, then use current time as default
	}
}

func (hub *ControllersHub) importBills(c echo.Context) error {
	req := &ImportBillsRequest{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// parse token
	if len(c.Request().Header["Authorization"]) != 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid header: Authorization")
	}
	token, startWithBearer := strings.CutPrefix(c.Request().Header["Authorization"][0], "Bearer ")
	if !startWithBearer {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid header: Authorization")
	}
	// query base user id
	userID, err := hub.userService.MustGetUserIDByToken(token)
	if err != nil {
		return err
	}
	createBillDTOs := make([]bill.CreateBillDTO, 0, len(req.Bills))
	for _, reqBill := range req.Bills {
		billDTO := bill.CreateBillDTO{
			Name:      reqBill.Name,
			Category:  reqBill.Category,
			Amount:    reqBill.Amount,
			CreatedAt: reqBill.CreatedAt,
		}
		createBillDTOs = append(createBillDTOs, billDTO)
	}
	// create bills
	err = hub.billService.CreateNewBills(userID, createBillDTOs)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, nil)
}
