package reply

import "xCut/entity"

type ReplyService interface {
	GetReplies() ([]entity.Reply, []error)
	GetReply(id uint) (*entity.Reply, []error)
	GetReplyByReviewID(reviewID uint) (*entity.Reply, []error)
	//getRecentReview
	StoreReply(reply *entity.Reply) (*entity.Reply, []error)
	UpdateReply(reply *entity.Reply) (*entity.Reply, []error)
	DeleteReply(id uint) (*entity.Reply, []error)
}
