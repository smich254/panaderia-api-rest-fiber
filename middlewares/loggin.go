package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Logging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		
		// Si hay un error, registra el error además del tiempo de procesamiento
		if err != nil {
			log.Printf(
				"[ERROR] [%s] %s %s - %v - Processed request in %v",
				c.Method(),
				c.Path(),
				c.IP(),
				err,
				duration,
			)
			// No olvides devolver el error para que el siguiente middleware o el manejador pueda manejarlo
			return err
		}

		log.Printf(
			"[%s] %s %s - Processed request in %v",
			c.Method(),
			c.Path(),
			c.IP(),
			duration,
		)

		return nil
	}
}
