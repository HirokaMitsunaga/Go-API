package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	SignUp(c echo.Context) error
	LogIn(c echo.Context) error
	LogOut(c echo.Context) error
	CsrfToken(e echo.Context) error
}

type userController struct {
	uu usecase.IUserUsecase
}

func NewUserController(uu usecase.IUserUsecase) IUserController {
	return &userController{uu}
}

func (uc *userController) SignUp(c echo.Context) error {
	//userから受け取るbodyを構造体へ変換する処理
	user := model.User{}
	//echoのBindを使うとclientからbodyの情報をuserが指し示すポインタへ格納してくれる(参照渡し)
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	//値渡し
	userRes, err := uc.uu.SignUp(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)

}

func (uc *userController) LogIn(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	tokenString, err := uc.uu.Login(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	//JWTトークンをサーバーサイドでcookieに設定していく
	//httpパッケージのCookie構造体を生成する（構造体が大きすぎるためnewを使っていると思われる）
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = tokenString
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	//cookie.Secure = true
	//HttpOnlyは、clientのJavaScriptからtokeの値が読み取れない用にするためtrueにする
	cookie.HttpOnly = true
	//フロントエンドとバックエンドのドメインが違うクロスドメイン間でのcookieの送受信のため、SameSiteNoneModeとしている
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)
}

func (uc *userController) LogOut(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "token"
	//値をクリアするために空の文字列を代入
	cookie.Value = ""
	//有効期限がすぐ切れるようにtime.Now()にしておく
	cookie.Expires = time.Now()
	cookie.Path = "/"
	cookie.Domain = os.Getenv("API_DOMAIN")
	//cookie.Secure = true
	//HttpOnlyは、clientのJavaScriptからtokeの値が読み取れない用にするためtrueにする
	cookie.HttpOnly = true
	//フロントエンドとバックエンドのドメインが違うクロスドメイン間でのcookieの送受信のため、SameSiteNoneModeとしている
	cookie.SameSite = http.SameSiteLaxMode
	c.SetCookie(cookie)
	return c.NoContent(http.StatusOK)

}

// 「/csrf」にアクセスがあった場合に、JSONでcstfトークンを返す
func (uc *userController) CsrfToken(c echo.Context) error {
	token := c.Get("csrf").(string)
	return c.JSON(http.StatusOK, echo.Map{
		"csrf_token": token,
	})
}
