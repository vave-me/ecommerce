package domain

import (
	"context"
)

type FakePostRepository struct {
	posts map[string]*Post
}

func NewFakePostRepository() *FakePostRepository {
	return &FakePostRepository{posts: map[string]*Post{}}
}

var _ PostRepository = (*FakePostRepository)(nil)

func (r *FakePostRepository) Load(ctx context.Context, postID string) (*Post, error) {
	if post, exists := r.posts[postID]; exists {
		return post, nil
	}

	return NewPost(postID), nil
}

func (r *FakePostRepository) Save(ctx context.Context, post *Post) error {
	for _, event := range post.Events() {
		if err := post.ApplyEvent(event); err != nil {
			return err
		}
	}

	r.posts[post.ID()] = post

	return nil
}

func (r *FakePostRepository) Reset(posts ...*Post) {
	r.posts = make(map[string]*Post)

	for _, post := range posts {
		r.posts[post.ID()] = post
	}
}
