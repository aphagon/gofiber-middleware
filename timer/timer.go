package timer

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

// New will measure how long it takes before a response is returned
func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// start timer
		start := time.Now()

		// next routes
		err := c.Next()

		// stop timer
		stop := time.Now()

		// Do something with response
		c.Append("Server-Timing", fmt.Sprintf("app;dur=%v", stop.Sub(start).String()))

		// return stack error if exist
		return err
	}
}
