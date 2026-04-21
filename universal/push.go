package universal

import (
	"context"
	"errors"
	"log/slog"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type PushMessage struct {
	Title string
	Body  string
	Link  string
}

type PushSender struct {
}

func (s *PushSender) SendAsync(ctx context.Context, token string, msg PushMessage) {
	jsonPath := os.Getenv("FIREBASE_JSON_PATH")
	go func() {
		// Wait for the original context to be done.
		select {
		case <-ctx.Done():
		}

		// If the context was canceled, it might be due to an error.
		// However, in many cases, successful requests also end with context cancellation.
		// If we want to be sure it wasn't an error, we'd need to check the cause
		// if WithCancelCause was used.
		if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) && !errors.Is(ctx.Err(), context.DeadlineExceeded) {
			// This would be some other error if any existed, but usually it's one of these two.
			return
		}

		// Now we use a background context to ensure the push notification is sent
		// regardless of the original context being done.
		asyncCtx := context.Background()

		var opt option.ClientOption
		if jsonPath != "" {
			opt = option.WithCredentialsFile(jsonPath)
		} else {
			opt = option.WithoutAuthentication()
		}

		config := &firebase.Config{ProjectID: "naborly-9f7dd"}
		app, err := firebase.NewApp(asyncCtx, config, opt)
		if err != nil {
			slog.Error("Failed to create firebase app", "error", err)
			return
		}

		client, err := app.Messaging(asyncCtx)
		if err != nil {
			slog.Error("Failed to create messaging client", "error", err)
			return
		}

		dataMap := make(map[string]string)
		if msg.Link != "" {
			dataMap["link"] = msg.Link
		}

		resultName, sendError := client.Send(asyncCtx, &messaging.Message{
			Token: token,
			Notification: &messaging.Notification{
				Title: msg.Title,
				Body:  msg.Body,
			},
			Data: dataMap,
		})

		if sendError != nil {
			slog.Error("Failed to send push notification", "error", sendError)
			return
		}

		slog.Info("Push notification sent", "result", resultName)
	}()
}
