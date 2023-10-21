package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/validator"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	SignUp(user model.User) (model.UserResponse, error)
	Login(user model.User) (string, error)
}

// userUsecaseのソースコードは、repositoryのInterfaceだけに依存させる
type userUsecase struct {
	ur repository.IUserRepository
	uv validator.IUserValidator
}

// インスタンスを生成することで、そのインスタンスの関数を実行できる
// コンストラクタの作成
// コンストラクタとは、「クラス」（設計図）の処理を実行する「インスタンス」（実際に作ったモノ）が生成される際に実行されるメソッドで
// 主にクラスのメンバ変数を初期化する際に使われる。
// IUserUsecase型がIUserUsecaseのインタフェースを満たすためには、インタフェースで定義されているメソッドを全て実装する必要がある。
func NewUserUsecase(ur repository.IUserRepository, uv validator.IUserValidator) IUserUsecase {
	return &userUsecase{ur, uv}
}

// UserのPasswordをhash化する
func (uu *userUsecase) SignUp(user model.User) (model.UserResponse, error) {
	//validationの実行
	if err := uu.uv.UserValidate(user); err != nil {
		return model.UserResponse{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return model.UserResponse{}, err
	}
	newUser := model.User{Email: user.Email, Password: string(hash)}
	//CreateUserは、userオブジェクトのポインタを引数で受け取るため、&newUserとする
	if err := uu.ur.CreateUser(&newUser); err != nil {
		return model.UserResponse{}, err
	}
	resUser := model.UserResponse{
		ID:    newUser.ID,
		Email: newUser.Email,
	}

	return resUser, nil
}

func (uu *userUsecase) Login(user model.User) (string, error) {
	//validationの実行
	if err := uu.uv.UserValidate(user); err != nil {
		return "", err
	}
	// DB内に存在するのか確認
	// 空のmode.userを作り判定する
	storedUser := model.User{}
	//入力されたEmailの確認
	// 何で空のstoredUserをGetUserByEmailへ渡しているのかわからん
	//-> GetUserByEmailでemailの情報を元にuserの値を取得するため
	//EmailからUser.IDを特定するため？
	if err := uu.ur.GetUserByEmail(&storedUser, user.Email); err != nil {
		return "", err
	}
	//入力されたパスワードの確認
	err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", err
	}
	//JWTトークンの生成
	//JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": storedUser.ID,
		//JWTの有効期限を12時間にしている
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})
	//SECRETはJWTのシークレットキー
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
