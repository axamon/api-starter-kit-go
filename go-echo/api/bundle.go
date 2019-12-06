package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
)



type Problem struct {
	Title string
	Status int32
	Detail string
}

type Timestamps struct {
	Ts time.Time
}

// This function wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func sendPetstoreError(ctx echo.Context, code int, message string) error {
	err := ctx.JSON(code, message)
	return err
}

/**
 * This structure contains all shared data for this app,
 * in our case it's a random number generator.
 */
type MyApplication struct {
	r *rand.Rand
}

func CreateApplication() *MyApplication {
	return &MyApplication{
		r: rand.New(rand.NewSource(99)),
	}
}

//
// Implement the methods declared in the generated interface.
//

func (app *MyApplication) GetEcho(ctx echo.Context) error {

	var ts = time.Now()
	result := Timestamps{
		Ts: ts,
	}
	return ctx.JSON(http.StatusOK, result)
}

func (app *MyApplication) GetStatus(ctx echo.Context) error {

	ctx.Response().Header().Set("Cache-Control", "no-store")

	var f = app.r.Int() % 5
	fmt.Println("string %d", f)
	if f < 3 {
		var result = Problem{
			Status: int32(200),
			Title:  "ok",
		}
		return ctx.JSON(http.StatusOK, result)
	}
	var result = Problem{
		Status: int32(503),
		Title:  "ko",
	}

	return ctx.JSON(http.StatusServiceUnavailable, result)
}

func CORSFilter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		return next(c)
	}
}

func ProblemErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)

	httpError, ok := err.(*echo.HTTPError)
	fmt.Println("problemErrorHandler: ", err, ok, httpError)
	if !ok {
		return
	}

	var problem Problem
	problem.Status = int32(httpError.Code)
	problem.Title, ok = httpError.Message.(string)
	if !ok {
		c.Logger().Error("Message is not a string", httpError)
		return
	}

	if httpError.Message == "Path was not found" {
		problem.Status = 404
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead { // Issue #608
			err = c.NoContent(httpError.Code)
		} else {
			c.Response().Header().Set("Content-Type", "application/problem+json")
			err = c.JSON(int(problem.Status), problem)
		}
		if err != nil {
			c.Logger().Error(err)
		}
	}

}
