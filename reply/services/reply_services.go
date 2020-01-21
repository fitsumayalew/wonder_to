package services

import (
	"xCut/entity"
	"xCut/reply"
)

type ReplyService struct {
	replyRepo reply.ReplyRepository
}

func (rp *ReplyService) GetReplies() ([]entity.Reply, []error) {
	rpl, errs := rp.replyRepo.GetReplies()
	if len(errs) > 0 {
		return nil, errs
	}
	return rpl, errs
}

func (rp *ReplyService) GetReplyByReviewID(reviewID uint) (*entity.Reply, []error) {
	rpl, errs := rp.replyRepo.GetReplyByReviewID(reviewID)
	if len(errs) > 0 {
		return nil, errs
	}
	return rpl, errs
}

func NewReplyService(rplRepo reply.ReplyRepository) reply.ReplyService {
	return &ReplyService{replyRepo: rplRepo}
}

// Getreplys returns all stored replys
func (rp *ReplyService) GetReply(id uint) (*entity.Reply, []error) {
	rpl, errs := rp.replyRepo.GetReply(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return rpl, errs
}

// Updatereply updates a given reply in a database
func (rp *ReplyService) UpdateReply(reply *entity.Reply) (*entity.Reply, []error) {
	rpl, errs := rp.replyRepo.UpdateReply(reply)
	if len(errs) > 0 {
		return nil, errs
	}
	return rpl, errs
}

// Deletereply deletes a given reply
func (rp *ReplyService) DeleteReply(id uint) (*entity.Reply, []error) {
	rpl, errs := rp.replyRepo.DeleteReply(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return rpl, errs
}

// Storereply stores a given reply
func (rp *ReplyService) StoreReply(reply *entity.Reply) (*entity.Reply, []error) {
	rpl, errs := rp.replyRepo.StoreReply(reply)
	if len(errs) > 0 {
		return nil, errs
	}
	return rpl, errs
}
