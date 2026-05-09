package crontab

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/samber/do"

	"codebase/internal/bob"
	"codebase/internal/service/alert"
	"codebase/internal/service/stock"
	"codebase/internal/service/utils"
	"codebase/pkg/telegram"
)

func RegisterPriceMonitorJob(container *do.Injector, cronManager interface{
	RegisterCrontab(desc string, runOnStart bool, d time.Duration, f func())
}) {
	stockSvc := do.MustInvoke[stock.ServiceStock](container)
	utilsSvc := do.MustInvoke[*utils.ServiceUtils](container)
	alertSvc := do.MustInvoke[alert.ServiceAlert](container)

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	tgClient := telegram.NewClient(token)

	// Run every 1 minute
	cronManager.RegisterCrontab("PriceMonitor", true, 1*time.Minute, func() {
		ctx := context.Background()

		alerts, err := alertSvc.ListAllActiveStockAlerts(ctx)
		if err != nil {
			log.Printf("PriceMonitor: failed to list alerts: %v\n", err)
			return
		}

		if len(alerts) == 0 {
			return
		}

		// Group alerts by symbol to avoid fetching stock info multiple times
		alertsBySymbol := make(map[string][]int)
		for i, a := range alerts {
			alertsBySymbol[a.Symbol] = append(alertsBySymbol[a.Symbol], i)
		}

		for symbol, indices := range alertsBySymbol {
			info, err := stockSvc.GetStockInfo(ctx, symbol, "MAIN")
			if err != nil {
				log.Printf("PriceMonitor: failed to get stock info for %s: %v\n", symbol, err)
				continue
			}

			currentPrice := info.MatchedPrice

			for _, idx := range indices {
				a := alerts[idx]

				var crossedHigh, crossedLow bool

				if a.HighThreshold != nil && currentPrice >= *a.HighThreshold {
					crossedHigh = true
				}
				if a.LowThreshold != nil && currentPrice <= *a.LowThreshold {
					crossedLow = true
				}

				if crossedHigh || crossedLow {
					// Let's query user by ID directly to get telegram_chat_id
					rawUser, err := utilsSvc.Datastore.FindRawUserByID(ctx, a.UserID)
					if err != nil || rawUser == nil || !rawUser.TelegramChatID.IsSet() || rawUser.TelegramChatID.MustGet() == "" {
						log.Printf("PriceMonitor: user %s has no telegram chat id, skipping alert for %s\n", a.UserID, symbol)
						continue
					}

					chatID := rawUser.TelegramChatID.MustGet()
					msg := fmt.Sprintf("⚠️ STOCK ALERT: %s\nCurrent Price: %.2f", symbol, currentPrice)
					if crossedHigh {
						msg += fmt.Sprintf("\nCrossed high threshold: %.2f", *a.HighThreshold)
					}
					if crossedLow {
						msg += fmt.Sprintf("\nCrossed low threshold: %.2f", *a.LowThreshold)
					}

					if err := tgClient.SendMessage(ctx, chatID, msg); err != nil {
						log.Printf("PriceMonitor: failed to send telegram message to %s: %v\n", chatID, err)
					} else {
						log.Printf("PriceMonitor: sent telegram alert to %s for %s\n", chatID, symbol)

						setter := &bob.UserStockAlertSetter{
							LastNotifiedAt: omitnull.From(time.Now()),
							IsActive:       omit.From(false), // Optional: disable it after firing
						}
						_ = utilsSvc.Datastore.UpdateStockAlert(ctx, a.ID, setter)
					}
				}
			}
		}
	})
}
