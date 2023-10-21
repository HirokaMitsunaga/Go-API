package repository

import (
	"go-rest-api/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUserByEmail(user *model.User, email string) error
	CreateUser(user *model.User) error
}

// この宣言は、dbという名前のフィールドが、*gorm.DB型のポインタを保持することを示しています。
// つまり、dbフィールドはgorm.DB型のオブジェクトへのポインタを指すポインタ型の変数です。
// このようなフィールド宣言は、通常の方法でポインタを持つフィールドを宣言するために使用されます。
type userRepository struct {
	db *gorm.DB
}

// repositoryにDBのインスタンスをDeoendancy Ingectionさせるためコンストラクタを作成する
// 引数のdbをインスタンスの要素にして、userRepositoryを作成しそのポインタをreturnで返す
// userRepository型がIUserRepositoryインタフェースを満たすためにはGetUserByEmailとCreateUserを実装する必要がある
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (ur *userRepository) GetUserByEmail(user *model.User, email string) error {
	//引数で受け取ったemailのuserが存在する場合は、引数のuserが指し示すアドレスの内容を検索したuserの内容で書き換える（user_usecase.goのLoginメソッドで使われている）
	if err := ur.db.Where("email=?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) CreateUser(user *model.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
