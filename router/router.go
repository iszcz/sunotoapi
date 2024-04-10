package router

import (
	"fksunoapi/cfg"
	"fksunoapi/serve"
	"fmt"
    "encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func CreateTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		ck := c.Get("Authorization")
		if ck == "" {
			ck = cfg.Config.App.Client
		} else {
			ck = serve.ParseToken(ck)
		}
		serve.Session = serve.GetSession(ck)
		var body []byte
		var errResp *serve.ErrorResponse
		if c.Path() == "/v2/generate" {
			body, errResp = serve.V2Generate(data, ck)
		} else if c.Path() == "/v2/lyrics/create" {
			body, errResp = serve.GenerateLyrics(data, ck)
		}
		if errResp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}

		return c.Status(fiber.StatusOK).Send(body)
	}
}

func GetTask() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		if len(data["ids"]) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeRequestFailed, "Cannot find ids"))
		}
		ck := c.Get("Authorization")
		if ck == "" {
			ck = cfg.Config.App.Client
		} else {
			ck = serve.ParseToken(ck)
		}
		serve.Session = serve.GetSession(ck)
		var body []byte
		var errResp *serve.ErrorResponse
		if c.Path() == "/v2/feed" {
			body, errResp = serve.V2GetFeedTask(data["ids"], ck)
		} else if c.Path() == "/v2/lyrics/task" {
			body, errResp = serve.GetLyricsTask(data["ids"], ck)
		}
		if errResp != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(errResp)
		}
		return c.Status(fiber.StatusOK).Send(body)
	}
}

func SunoChat() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(serve.NewErrorResponse(serve.ErrCodeJsonFailed, "Cannot parse JSON"))
		}
		ck := c.Get("Authorization")
		if ck == "" {
			ck = cfg.Config.App.Client
		} else {
			ck = serve.ParseToken(ck)
		}
		//serve.Session = serve.GetSession(ck)
		fmt.Println(data)
		res, errResp := serve.SunoChat(data, ck)

		// 如果res不是[]byte类型，将其转换为JSON字符串以便打印
        resJSON, err := json.Marshal(res)
        if err != nil {
            // 处理可能的错误
            fmt.Println("Error marshalling response to JSON:", err)
        } else {
            // 打印JSON字符串
            fmt.Println(string(resJSON))
        }

        if errResp != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(errResp)
        }

        c.Set("Content-Type", "application/json")

        return c.Status(fiber.StatusOK).JSON(res)
	}
}

func SetupRoutes(app *fiber.App) {
	app.Use(logger.New(logger.ConfigDefault))
	app.Use(cors.New(cors.ConfigDefault))
	app.Post("/v2/generate", CreateTask())
	app.Post("/v2/feed", GetTask())
	app.Post("/v2/lyrics/create", CreateTask())
	app.Post("/v2/lyrics/task", GetTask())
	app.Post("/v1/chat/completions", SunoChat())
}
