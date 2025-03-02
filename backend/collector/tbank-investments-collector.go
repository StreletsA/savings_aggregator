package collector

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type TBankInvestmentsSavingsCollector struct {
}

func (TBankInvestmentsSavingsCollector) GetSourceType() SavingsCollectorSourceType {
	return T_BANK_INVESTMENTS_SOURCE
}

func (TBankInvestmentsSavingsCollector) Collect() (SavingsCollectionInfo, error) {
	// загружаем конфигурацию для сдк из .yaml файла
	config, err := investgo.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("config loading error %v", err.Error())
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer cancel()
	// сдк использует для внутреннего логирования investgo.Logger
	// для примера передадим uber.zap
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}

	logger := l.Sugar()

	// создаем клиента для investAPI, он позволяет создавать нужные сервисы и уже
	// через них вызывать нужные методы
	client, err := investgo.NewClient(ctx, config, logger)
	if err != nil {
		logger.Fatalf("client creating error %v", err.Error())
	}
	defer func() {
		err := client.Stop()
		if err != nil {
			logger.Errorf("client shutdown error %v", err.Error())
		}
	}()

	// создаем клиента для сервиса счетов
	usersService := client.NewUsersServiceClient()

	var accId string
	status := pb.AccountStatus_ACCOUNT_STATUS_ALL
	accsResp, err := usersService.GetAccounts(&status)
	if err != nil {
		return SavingsCollectionInfo{}, err
	} else {
		accs := accsResp.GetAccounts()
		for _, acc := range accs {
			// берем последний счет
			accId = acc.GetId()
		}
	}

	operationsService := client.NewOperationsServiceClient()
	portfolio, err := operationsService.GetPortfolio(accId, 0)
	if err != nil {
		return SavingsCollectionInfo{}, err
	}

	total := portfolio.GetTotalAmountPortfolio().ToFloat()

	return SavingsCollectionInfo{SourceType: T_BANK_INVESTMENTS_SOURCE, Total: float32(total)}, nil
}
