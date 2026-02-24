package users

import "context"

type Service struct {
	repo RepositoryGorm
}

func NewService(repo RepositoryGorm) *Service {
	return &Service{repo: repo}
}

func (s *Service) Get(ctx context.Context, id uint64) (UserView, error) {
	if id == 0 {
		return UserView{}, ErrInvalidArgument
	}
	user, err := s.repo.Get(ctx, id)
	if err != nil {
		return UserView{}, err
	}
	return toView(user), nil
}

func (s *Service) List(ctx context.Context, limit, offset int) (ListUsersResponse, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	list, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return ListUsersResponse{}, err
	}

	items := make([]UserView, 0, len(list))
	for _, user := range list {
		items = append(items, toView(user))
	}

	return ListUsersResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}, nil

}

func toView(u User) UserView {
	return UserView{
		ID:        u.ID,
		TgID:      u.TgID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		PhotoURL:  u.PhotoURL,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
