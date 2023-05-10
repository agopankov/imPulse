package monitor

import (
	"context"
	"fmt"
	"github.com/agopankov/binance/client/pkg/telegram"
	"github.com/agopankov/binance/client/pkg/tracker"
	"github.com/agopankov/binance/server/pkg/grpcbinance/proto"
	tele "gopkg.in/telebot.v3"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Monitor struct {
	TelegramClient       *telegram.Client
	SecondTelegramClient *telegram.Client
	BinanceClient        proto.BinanceServiceClient
	ChatID               int64
	SecondChatID         int64
	TrackerInstance      *tracker.Tracker
}

func getPriceForSymbol(symbol string, prices []*proto.USDTPrice) string {
	for _, price := range prices {
		if price.Symbol == symbol {
			return fmt.Sprintf("%.8f", price.Price)
		}
	}
	return ""
}

func PriceChanges(ctx context.Context, client *telegram.Client, secondTelegramClient *telegram.Client, binanceClient proto.BinanceServiceClient, chatID int64, secondChatID int64, trackerInstance *tracker.Tracker) {
	ticker := time.NewTicker(5 * time.Second)
	notifyTicker := time.NewTicker(1 * time.Minute)
	logTicker := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-logTicker.C:
			processLogTicker(trackerInstance)
		case <-ticker.C:
			processTicker(client, binanceClient, chatID, trackerInstance)
		case <-notifyTicker.C:
			processNotifyTicker(client, secondTelegramClient, binanceClient, chatID, secondChatID, trackerInstance)
		}
	}
}

func processTicker(telegramClient *telegram.Client, binanceClient proto.BinanceServiceClient, chatID int64, trackerInstance *tracker.Tracker) {
	ctx := context.Background()
	usdtPrices, err := binanceClient.GetUSDTPrices(ctx, &proto.Empty{})
	if err != nil {
		log.Printf("Error getting USDT prices: %v", err)
		return
	}

	changePercent, err := binanceClient.Get24HChangePercent(ctx, &proto.Empty{})
	if err != nil {
		log.Printf("Error getting 24h change percent: %v", err)
		return
	}

	var newTrackedSymbols []tracker.SymbolChange
	for _, price := range usdtPrices.Prices {
		change := 0.0
		for _, changePercent := range changePercent.ChangePercents {
			if price.Symbol == changePercent.Symbol {
				change = changePercent.ChangePercent
				break
			}
		}

		if change >= 20 && !trackerInstance.IsTracked(price.Symbol) {
			newSymbol := tracker.SymbolChange{
				Symbol:         price.Symbol,
				PriceChange:    fmt.Sprintf("%.8f", price.Price),
				PriceChangePct: change,
				AddedAt:        time.Now(),
			}
			trackerInstance.UpdateTrackedSymbol(newSymbol)
			newTrackedSymbols = append(newTrackedSymbols, newSymbol)
		}
	}

	sort.Slice(newTrackedSymbols, func(i, j int) bool {
		return newTrackedSymbols[i].PriceChangePct > newTrackedSymbols[j].PriceChangePct
	})

	var messageBuilder strings.Builder
	for _, symbolChange := range newTrackedSymbols {
		emoji := "✅"
		price := strings.TrimRight(strings.TrimRight(symbolChange.PriceChange, "0"), ".")
		message := fmt.Sprintf("%s %s / USDT P: %s Ch24h: %.2f%% \n", emoji, symbolChange.Symbol[:len(symbolChange.Symbol)-4], price, symbolChange.PriceChangePct)
		messageBuilder.WriteString(message)
	}

	if messageBuilder.Len() > 0 {
		message := messageBuilder.String()
		recipient := &tele.User{ID: chatID}
		_, err := telegramClient.SendMessage(recipient, message)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	}

	for symbol, _ := range trackerInstance.GetTrackedSymbols() {
		change24h := 0.0
		for _, changePercent := range changePercent.ChangePercents {
			if symbol == changePercent.Symbol {
				change24h = changePercent.ChangePercent
				break
			}
		}

		if change24h <= 20 {
			trackerInstance.RemoveTrackedSymbol(symbol)
			continue
		}
	}
}

func processNotifyTicker(telegramClient *telegram.Client, secondTelegramClient *telegram.Client, binanceClient proto.BinanceServiceClient, chatID int64, secondChatID int64, trackerInstance *tracker.Tracker) {
	ctx := context.Background()
	usdtPrices, err := binanceClient.GetUSDTPrices(ctx, &proto.Empty{})
	if err != nil {
		log.Printf("Error getting USDT prices: %v", err)
		return
	}

	changePercent, err := binanceClient.Get24HChangePercent(ctx, &proto.Empty{})
	if err != nil {
		log.Printf("Error getting 24h change percent: %v", err)
		return
	}

	trackedSymbols := trackerInstance.GetTrackedSymbols()

	var sortedSymbols []tracker.SymbolChange
	for _, symbolChange := range trackedSymbols {
		sortedSymbols = append(sortedSymbols, symbolChange)
	}
	sort.Slice(sortedSymbols, func(i, j int) bool {
		return sortedSymbols[i].PriceChangePct > sortedSymbols[j].PriceChangePct
	})

	var messageBuilder strings.Builder
	for _, symbolChange := range sortedSymbols {
		currentPrice := getPriceForSymbol(symbolChange.Symbol, usdtPrices.Prices)

		var change24h float64
		for _, changePercentData := range changePercent.ChangePercents {
			if symbolChange.Symbol == changePercentData.Symbol {
				change24h = changePercentData.ChangePercent
				break
			}
		}

		price := strings.TrimRight(strings.TrimRight(currentPrice, "0"), ".")
		emoji := ""
		currentPriceFloat, _ := strconv.ParseFloat(currentPrice, 64)
		previousPriceFloat, _ := strconv.ParseFloat(symbolChange.PriceChange, 64)
		if currentPriceFloat > previousPriceFloat {
			emoji = "📈"
		} else if currentPriceFloat < previousPriceFloat {
			emoji = "📉"
		} else {
			emoji = "🔹"
		}

		message := fmt.Sprintf("%s %s / USDT P: %s Ch24h: %.2f%% \n", emoji, symbolChange.Symbol[:len(symbolChange.Symbol)-4], price, change24h)
		messageBuilder.WriteString(message)

		//if symbolChange.IsNew {
		//	recipient := &tele.User{ID: secondChatID}
		//	_, err := secondTelegramClient.SendMessage(recipient, message)
		//	if err != nil {
		//		log.Printf("Error sending message to the second chat: %v\n", err)
		//	}
		//	symbolChange.IsNew = false
		//}

		symbolChange.PriceChange = currentPrice
		trackerInstance.UpdateTrackedSymbol(symbolChange)
	}

	if messageBuilder.Len() > 0 {
		message := messageBuilder.String()
		recipient := &tele.User{ID: chatID}
		_, err := telegramClient.SendMessage(recipient, message)
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
	}

	for _, symbolChange := range sortedSymbols {
		currentPrice := getPriceForSymbol(symbolChange.Symbol, usdtPrices.Prices)

		currentPriceFloat, _ := strconv.ParseFloat(currentPrice, 64)
		previousPriceFloat, _ := strconv.ParseFloat(symbolChange.PriceChange, 64)

		if !symbolChange.IsNew && time.Since(symbolChange.AddedAt) <= 15*time.Minute && (currentPriceFloat/previousPriceFloat)-1 >= 0.05 {
			message := fmt.Sprintf("🚀 %s / USDT P: %.3f Ch24h: %.2f%% (PrP: %.3f) \n",
				symbolChange.Symbol[:len(symbolChange.Symbol)-4],
				currentPriceFloat,
				symbolChange.PriceChangePct,
				previousPriceFloat,
			)
			messageBuilder.WriteString(message)

			recipient := &tele.User{ID: secondChatID}
			_, err := secondTelegramClient.SendMessage(recipient, message)
			if err != nil {
				log.Printf("Error sending message to the second chat: %v\n", err)
			}
		}

		symbolChange.PriceChange = currentPrice
		trackerInstance.UpdateTrackedSymbol(symbolChange)
	}

}

func processLogTicker(trackerInstance *tracker.Tracker) {
	trackedSymbols := trackerInstance.GetTrackedSymbols()
	if len(trackedSymbols) == 0 {
		log.Printf("TrackedSymbols is empty")
	}
	for symbol, symbolChange := range trackedSymbols {
		log.Printf("Symbol: %s, PriceChange: %s, PriceChangePct: %.2f%%, AddedAt: %v\n",
			symbol, symbolChange.PriceChange, symbolChange.PriceChangePct, symbolChange.AddedAt)
	}
}
