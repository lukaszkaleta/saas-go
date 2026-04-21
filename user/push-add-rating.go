package user

import (
	"context"
	"log/slog"

	"github.com/lukaszkaleta/saas-go/universal"
)

type PushAddRating struct {
	ratings universal.Adder[*universal.RatingModel, universal.Rating]
	users   Users
	sender  *universal.PushSender
}

func NewPushAddRating(ratings universal.Adder[*universal.RatingModel, universal.Rating], users Users) *PushAddRating {
	return &PushAddRating{
		ratings: ratings,
		users:   users,
		sender:  &universal.PushSender{},
	}
}

func (p *PushAddRating) Add(ctx context.Context, r *universal.RatingModel) (universal.Rating, error) {
	rating, err := p.ratings.Add(ctx, r)
	if err != nil {
		return nil, err
	}

	pushErr := p.sendPush(ctx, r)
	if pushErr != nil {
		slog.Error("Can not send push", "Error", pushErr.Error())
	}

	return rating, nil
}

func (p *PushAddRating) sendPush(ctx context.Context, r *universal.RatingModel) error {
	reviewee, err := p.users.ById(ctx, r.RevieweeId)
	if err != nil {
		return err
	}

	account := reviewee.Account()
	accountModel, err := account.Model(ctx)
	if err != nil || accountModel.FirebaseToken == "" {
		return nil
	}

	reviewer := CurrentUser(ctx)
	if reviewer == nil {
		return nil
	}

	pushMsg := universal.PushMessage{
		Title: reviewer.Person.FirstName,
		Body:  "Added a rating for you",
		Link:  "https://naborly.no/profile",
	}

	p.sender.SendAsync(ctx, accountModel.FirebaseToken, pushMsg)

	return nil
}
