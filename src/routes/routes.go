package routes

import (
	"net/http"
	"os"

	"git.xenonstack.com/xs-onboarding/document-manage/config"
	"git.xenonstack.com/xs-onboarding/document-manage/src/api"
	"git.xenonstack.com/xs-onboarding/document-manage/src/token"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// Routes is a method in which all the service endpoints are defined
func Routes(router *gin.Engine) {

	//Healthz Check
	router.GET("/healthz", api.Healthz)

	router.StaticFile("/openapi.yaml", "./openapi.yaml")
	// developers help endpoint
	if config.Conf.Service.Environment != "production" {
		// endpoint to read variables
		router.GET("/end", checkToken, readEnv)
		router.GET("/logs", checkToken, readLogs)
	}
	authMiddleware := token.MwInitializer()
	router.Use(authMiddleware.MiddlewareFunc())
	{ //filetype can be of two types either uploads or policy
		router.POST("/user/upload", api.UploadDocuments)
		router.GET("/getpolicy", api.GetPolicy)
		router.DELETE("/delete/:foldername/:filename", api.DeleteDocument)
		//router.Use(checkAdmin)
		{
			router.POST("/uploadpolicy", api.UploadPolicy)
			router.GET("/getfile/:foldername/:filename", api.GetFile)
		}
	}
}

//func checkAdmin is used to check user is admin
func checkAdmin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	if claims["role"].(string) != "admin" {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	c.Next()
}

// readLogs is a api handler for reading logs
func readLogs(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "info.txt")
}

// readEnv is api handler for reading configuration variables data
func readEnv(c *gin.Context) {
	if config.TomlFile == "" {
		// if configuration is done using environment variables
		env := make([]string, 0)
		for _, pair := range os.Environ() {
			env = append(env, pair)
		}
		c.JSON(200, gin.H{
			"environments": env,
		})
	} else {
		// if configuration is done using toml file
		http.ServeFile(c.Writer, c.Request, config.TomlFile)
	}
}

// checkToken is a middleware to check header is set or not for secured api
func checkToken(c *gin.Context) {
	xt := c.Request.Header.Get("XSOnboarding-token")
	if xt != "XSOnboarding" {
		c.Abort()
		c.JSON(404, gin.H{})
		return
	}
	c.Next()
}
