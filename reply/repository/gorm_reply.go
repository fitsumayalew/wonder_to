package repository

import (
	"github.com/jinzhu/gorm"
	"xcut/entity"
	"xcut/reply"
)

type ReplyGormRepo struct {
	conn *gorm.DB
}




// NewReviewtGormRepo returns new object of ReviewGormRepo
func NewReplyGormRepo(db *gorm.DB) reply.ReplyRepository {
	return &ReplyGormRepo{conn: db}
}

//  GetReviews returns all reviews stored in the database
func (replyRepo *ReplyGormRepo) GetReplies() ([]entity.Reply, []error) {
	repl := []entity.Reply{}
	errs := replyRepo.conn.Find(&repl).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return repl, errs
}

func (replyRepo *ReplyGormRepo) GetReplyByReviewID(reviewID uint) (*entity.Reply, []error) {
	revo := entity.Reply{}
	errs := replyRepo.conn.Find(&revo, "shop_id=?",reviewID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &revo, errs
}



func (replyRepo *ReplyGormRepo) GetReply(id uint) (*entity.Reply, []error) {
	reply := entity.Reply{}
	errs := replyRepo.conn.First(&reply,id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &reply, errs
}




// GetReviewByShopID retrieve a review from the database by its shopid

// UpdateReview updates a given review in the database
func (replyRepo *ReplyGormRepo) UpdateReply(reply *entity.Reply) (*entity.Reply , []error){
	repl := reply
	errs := replyRepo.conn.Save(repl).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return repl, errs
}

// DeleteReview deletes a given review from the database
func (replyRepo *ReplyGormRepo) DeleteReply(id uint) (*entity.Reply, []error){
	repl, errs := replyRepo.GetReply(id)
	if len(errs) > 0{
		return nil, errs
	}
	errs = replyRepo.conn.Delete(repl, repl.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return repl, errs
}

// StoreReview stores a given review in the database
func (replyRepo *ReplyGormRepo)StoreReply(reply *entity.Reply) (*entity.Reply, []error){
	repl := reply
	errs := replyRepo.conn.Create(repl).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return repl, errs
}