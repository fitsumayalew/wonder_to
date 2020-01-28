package repository


import (
	"github.com/jinzhu/gorm"
	"xCut/mock_entity"
	"xCut/user"
)

type MockSessionRepo struct {
	conn *gorm.DB
}


// NewShoptGormRepo returns new object of MockShopRepo
func NewMockSessionRepo(db *gorm.DB) Session.SessionRepository {
	return &MockSessionRepo{conn: db}
}

func (mSessionRepo *MockSessionRepo) Session(sessionId string) (*entity.Session, []error) {
Sess := entity.MockSession
if id != 1 {
return nil, []error{errors.New("Not found")}
}
return &Sess, nil
}

//// Returns all the sessions
func (mSessionRepo *MockSessionRepo) Sessions() ([]*entity.Session, []error) {
	Sess := []entity.Session{entity.MockSession}
	return Sess, nil
}



func (mSessionRepo *MockSessionRepo) DeleteSession(id uint) (*entity.Session, []error) {
Sess := entity.MockSession
if id != 1 {
return nil, []error{errors.New("Not found")}
}
return &Sess, nil
}

// StoreSession stores a given Session in the database
func (mSessionRepo *MockSessionRepo) StoreSession(session *entity.Session) (*entity.Session, []error) {
Sess := session
return Sess,nil
}