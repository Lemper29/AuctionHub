package main

import (
	"context"
	"fmt"
	pb "github/auction/auction-service/gen/proto"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuctionServiceClient(conn)

	for {
		// Меню выбора
		fmt.Println("\n=== АУКЦИОН ===")
		fmt.Println("1 - Создать лот")
		fmt.Println("2 - Подписаться на лот")
		fmt.Println("3 - Сделать ставку")
		fmt.Println("0 - Выход")
		fmt.Print("Введите номер: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			createLotInteractive(client, ctx)
		case 2:
			subscribeToLotInteractive(client, ctx)
		case 3:
			placeBidInteractive(client, ctx)
		case 0:
			fmt.Println("До свидания!")
			return
		default:
			fmt.Println("Неверный выбор")
		}
	}
}

func createLotInteractive(client pb.AuctionServiceClient, ctx context.Context) {
	var name, description string
	var startPrice float64
	var durationMinute int64

	fmt.Print("Введите название лота: ")
	fmt.Scanln(&name)

	fmt.Print("Введите описание лота: ")
	fmt.Scanln(&description)

	fmt.Print("Введите стартовую цену: ")
	fmt.Scanln(&startPrice)

	fmt.Print("Введите длительность аукциона (минуты): ")
	fmt.Scanln(&durationMinute)

	createLot, err := client.CreateLot(ctx, &pb.CreateLotRequest{
		Name:           name,
		Description:    description,
		StartPrice:     startPrice,
		DurationMinute: durationMinute,
	})
	if err != nil {
		log.Printf("Ошибка создания лота: %v", err)
		return
	}

	// ВЫВОДИМ ID ДЛЯ ПОЛЬЗОВАТЕЛЯ!
	fmt.Printf("✅ Лот создан! ID: %s\n", createLot.Lot.Id)
	fmt.Printf("   Название: %s\n", createLot.Lot.Name)
	fmt.Printf("   Стартовая цена: %.2f\n", createLot.Lot.StartPrice)
}

func subscribeToLotInteractive(client pb.AuctionServiceClient, ctx context.Context) {
	var lotID string

	fmt.Print("Введите ID лота для подписки: ")
	fmt.Scanln(&lotID)

	// Сначала проверяем существование лота
	_, err := client.GetLot(ctx, &pb.GetLotRequest{LotId: lotID})
	if err != nil {
		log.Printf("❌ Лот с ID %s не найден: %v", lotID, err)
		return
	}

	stream, err := client.SubscribeToLot(ctx, &pb.SubscribeToLotRequest{
		LotId: lotID,
	})
	if err != nil {
		log.Printf("Ошибка подписки: %v", err)
		return
	}

	fmt.Printf("🔔 Подписка на лот %s активна...\n", lotID)
	fmt.Println("Нажмите Ctrl+C для выхода")

	for {
		lot, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Подписка завершена")
			break
		}
		if err != nil {
			log.Printf("Ошибка получения: %v", err)
			break
		}

		fmt.Printf("📢 Обновление: %s - цена: %.2f, победитель: %s\n",
			lot.Lot.Name,
			lot.Lot.CurrentPrice,
			lot.Lot.CurrentWinner)
	}
}

func placeBidInteractive(client pb.AuctionServiceClient, ctx context.Context) {
	var lotID, userID string
	var amount float64

	fmt.Print("Введите ID лота: ")
	fmt.Scanln(&lotID)

	fmt.Print("Введите ваш ID пользователя: ")
	fmt.Scanln(&userID)

	fmt.Print("Введите сумму ставки: ")
	fmt.Scanln(&amount)

	response, err := client.PlaceBid(ctx, &pb.PlaceBidRequest{
		LotId:  lotID,
		UserId: userID,
		Amount: amount,
	})
	if err != nil {
		log.Printf("Ошибка ставки: %v", err)
		return
	}

	if response.Success {
		fmt.Println("✅ Ставка принята!")
		fmt.Printf("   Текущая цена: %.2f\n", response.UpdatedLot.CurrentPrice)
	} else {
		fmt.Printf("❌ Ставка отклонена: %s\n", response.Message)
	}
}
