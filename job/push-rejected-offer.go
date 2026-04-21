package job

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/lukaszkaleta/saas-go/universal"
	"github.com/lukaszkaleta/saas-go/user"
)

type PushRejectedOffer struct {
	inner universal.Rejecter
	users user.Users
	offer Offer
}

func (m *PushRejectedOffer) Reject(ctx context.Context) error {
	err := m.inner.Reject(ctx)
	if err != nil {
		return err
	}

	userId, err := universal.CreatedById[OfferModel](ctx, m.offer)
	if err != nil {
		// Even if we fail to get userId for push, the main operation succeeded
		slog.Error("Failed to get creator ID for push notification", "error", err)
		return nil
	}
	u, err := m.users.ById(ctx, userId)
	if err != nil {
		slog.Error("Failed to get user for push notification", "error", err)
		return nil
	}

	model, err := u.Account().Model(ctx)
	if err != nil {
		slog.Error("Failed to get user model for push notification", "error", err)
		return nil
	}
	token := model.FirebaseToken
	if token != "" {
		offerModel, err := m.offer.Model(ctx)
		if err == nil {
			jobId := strconv.FormatInt(offerModel.JobId, 10)
			offerId := strconv.FormatInt(offerModel.Id, 10)
			sender := universal.PushSender{}
			sender.SendAsync(ctx, token, universal.PushMessage{
				Title: "Offer rejected",
				Body:  "Your offer has been rejected.",
				Link:  "https://naborly.no/offer/" + string(jobId) + "/" + offerId,
			})
		}
	}

	return nil
}

func NewPushRejectedOffer(users user.Users, offer Offer, inner universal.Rejecter) universal.Rejecter {
	return &PushRejectedOffer{
		inner: inner,
		users: users,
		offer: offer,
	}
}
