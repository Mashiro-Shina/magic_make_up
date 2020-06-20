package handlers

import (
	"log"
	"magicMakeup/repositories"
)

// 手动模拟递增的 comment id
var commentIDChannel <-chan int
var commentIDKey string
var commentID int

func init() {
	commentIDChannel = GetCommentID()
	commentIDKey = "comment_id"

	ok, err := repositories.HasKey(commentIDKey)
	if err != nil {
		log.Fatalf("initialize comment id failed: %v\n", err)
		return
	}

	if !ok {
		err = repositories.InitCommentID(commentIDKey)
		if err != nil {
			log.Fatalf("initialize comment id failed: %v\n", err)
			return
		}
	}

	commentID, err = repositories.GetCommentIDFromRedis(commentIDKey)
	if err != nil {
		log.Fatalf("get comment id failed: %v\n", err)
		return
	}
}

func GetCommentID() <-chan int {
	ch := make(chan int)

	go func() {
		for {
			ch <- commentID
			commentID++
			_ = repositories.SetCommentID(commentIDKey, commentID)
		}
	}()

	return ch
}
