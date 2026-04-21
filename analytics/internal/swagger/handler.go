package swagger

import (
	"os"

	"github.com/gofiber/fiber/v3"
)

var openAPIYAML []byte

func init() {
	var err error
	openAPIYAML, err = os.ReadFile("docs/openapi.yaml")
	if err != nil {
		panic(err)
	}
}

func Handler(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html")
	return c.Send(swaggerHTML)
}

func OpenAPI(c fiber.Ctx) error {
	c.Set("Content-Type", "application/vnd.yaml")
	return c.Send(openAPIYAML)
}

var swaggerHTML = []byte(`<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui.css">
    <script src="https://cdn.jsdelivr.net/npm/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <title>API Documentation</title>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>
        window.onload = () => {
            SwaggerUIBundle({
                url: "/api/v1/swagger/openapi.yaml",
                dom_id: "#swagger-ui",
            });
        };
    </script>
</body>
</html>`)
