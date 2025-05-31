package repo

import (
	"context"
	"errors"
	"online-judge/internal/models"
	"sync"
)

type DiscussionRepo struct {
	mu                  sync.Mutex
	db                  []models.DiscussionDB
	next                int
	discussionReactions map[int]map[int]int // userID -> discussionID -> reaction
	commentReactions    map[int]map[int]int // userID -> commentID -> reaction
	nextCommentID       int
}

func NewDiscussionRepo() *DiscussionRepo {
	return &DiscussionRepo{
		db:                  make([]models.DiscussionDB, 0),
		next:                1,
		nextCommentID:       1,
		discussionReactions: make(map[int]map[int]int),
		commentReactions:    make(map[int]map[int]int),
	}
}

func (r *DiscussionRepo) CreateDiscussion(ctx context.Context, userID int, problemID int, title string, content string, tagIDs []int) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := models.UserInfo{
		ID:       userID,
		Username: "User" + string(rune(userID)),
	}

	var tags []models.Tag
	for _, tagID := range tagIDs {
		tags = append(tags, models.Tag{
			ID:   tagID,
			Name: "Tag" + string(rune(tagID)),
		})
	}

	discussion := models.DiscussionDB{
		ID:        r.next,
		ProblemID: problemID,
		CreatedBy: user,
		Title:     title,
		Content:   content,
		Tags:      tags,
	}

	r.db = append(r.db, discussion)
	r.next++

	return discussion.ID, nil
}

func (r *DiscussionRepo) UpdateDiscussion(ctx context.Context, userID int, problemID int, title string, content string, tagIDs []int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, d := range r.db {
		if d.ProblemID == problemID && d.CreatedBy.ID == userID {
			r.db[i].Title = title
			r.db[i].Content = content

			var tags []models.Tag
			for _, tagID := range tagIDs {
				tags = append(tags, models.Tag{
					ID:   tagID,
					Name: "Tag" + string(rune(tagID)),
				})
			}
			r.db[i].Tags = tags
			return nil
		}
	}
	return errors.New("discussion not found or permission denied")
}

func (r *DiscussionRepo) GetDiscussion(ctx context.Context, problemID int, userID int) ([]models.DiscussionDB, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []models.DiscussionDB
	for _, d := range r.db {
		if d.ProblemID == problemID {
			// Set reaction from current user
			if r.discussionReactions[userID] != nil {
				d.CurrentUserReaction = r.discussionReactions[userID][d.ID]
			}
			// Set reactions for each comment
			for i := range d.Comments {
				if r.commentReactions[userID] != nil {
					d.Comments[i].CurrentUserReaction = r.commentReactions[userID][d.Comments[i].ID]
				}
			}
			result = append(result, d)
		}
	}
	return result, nil
}

func (r *DiscussionRepo) ReactToDiscussion(ctx context.Context, userID int, reaction int, discussionID int) error {
	if reaction != -1 && reaction != 0 && reaction != 1 {
		return errors.New("invalid reaction")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.discussionReactions[userID] == nil {
		r.discussionReactions[userID] = make(map[int]int)
	}

	for i, d := range r.db {
		if d.ID == discussionID {
			prev := r.discussionReactions[userID][discussionID]
			r.db[i].Likes -= boolToInt(prev == 1)
			r.db[i].Dislikes -= boolToInt(prev == -1)

			r.discussionReactions[userID][discussionID] = reaction

			r.db[i].Likes += boolToInt(reaction == 1)
			r.db[i].Dislikes += boolToInt(reaction == -1)
			return nil
		}
	}
	return errors.New("discussion not found")
}

func (r *DiscussionRepo) AddCommentToDiscussion(ctx context.Context, userID int, discussionID int, content string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user := models.UserInfo{
		ID:       userID,
		Username: "User" + string(rune(userID)),
	}

	comment := models.CommentDB{
		ID:           r.nextCommentID,
		DiscussionID: discussionID,
		CreatedBy:    user,
		Content:      content,
	}
	r.nextCommentID++

	for i := range r.db {
		if r.db[i].ID == discussionID {
			r.db[i].Comments = append(r.db[i].Comments, comment)
			return 0, nil
		}
	}

	return comment.ID, errors.New("discussion not found")
}

func (r *DiscussionRepo) ReactToComment(ctx context.Context, userID int, reaction int, commentID int) error {
	if reaction != -1 && reaction != 0 && reaction != 1 {
		return errors.New("invalid reaction")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.commentReactions[userID] == nil {
		r.commentReactions[userID] = make(map[int]int)
	}

	for i := range r.db {
		for j := range r.db[i].Comments {
			if r.db[i].Comments[j].ID == commentID {
				prev := r.commentReactions[userID][commentID]
				r.db[i].Comments[j].Likes -= boolToInt(prev == 1)
				r.db[i].Comments[j].Dislikes -= boolToInt(prev == -1)

				r.commentReactions[userID][commentID] = reaction

				r.db[i].Comments[j].Likes += boolToInt(reaction == 1)
				r.db[i].Comments[j].Dislikes += boolToInt(reaction == -1)
				return nil
			}
		}
	}
	return errors.New("comment not found")
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
