package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	e := echo.New()
	//CORSのMWの内容を追加
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FR_URL")},
		//HeaderXCSRFTokenを含めることでheader経由でCsrfトークンを取得する事ができる
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods: []string{"GET", "PUT", "POST", "DELETE"},
		//AllowCredentialsをtrueにする事でcookieの送受信を可能にしている
		AllowCredentials: true,
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		//postmanで動作確認をできるようにするために、samesiteをデフォルトモードにする
		//CookieSameSite: http.SameSiteDefaultMode,
		//postmanの動作確認後NoneModeへ変更する
		CookieSameSite: http.SameSiteNoneMode,
	}))
	//パスによって呼び出すコントローラーをそれぞれ定義
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)
	t := e.Group("/tasks")
	//Useを使うことでエンドポイントにミドルウェアを追加する事ができる。
	//今回はechojwtというミドルウェアを適用
	//TokenLookupはどこにjwtトークンが格納されているのかを指し示す
	t.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)
	return e
}
