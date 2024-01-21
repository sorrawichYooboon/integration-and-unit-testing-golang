package accountservice

import (
	"errors"
	"events"
	"log"
	"producer/commands"
	services "producer/services/producer"

	"github.com/google/uuid"
)

type IAccountService interface {
	OpenAccount(command commands.OpenAccountCommand) (id string, err error)
	DepositFund(command commands.DepositFundCommand) error
	WithdrawFund(command commands.WithdrawFundCommand) error
	CloseAccount(command commands.CloseAccountCommand) error
}

type accountService struct {
	eventProducer services.IEventProducer
}

func NewAccountService(eventProducer services.IEventProducer) IAccountService {
	return accountService{eventProducer}
}

func (sv accountService) OpenAccount(command commands.OpenAccountCommand) (id string, err error) {

	if command.AccountHolder == "" || command.AccountType == 0 || command.OpeningBalance == 0 {
		return "", errors.New("bad request")
	}

	event := events.OpenAccountEvent{
		ID:             uuid.NewString(),
		AccountHolder:  command.AccountHolder,
		AccountType:    command.AccountType,
		OpeningBalance: command.OpeningBalance,
	}

	log.Printf("%#v", event)
	return event.AccountHolder, sv.eventProducer.Produce(event)
}

func (sv accountService) DepositFund(command commands.DepositFundCommand) error {
	if command.ID == "" || command.Amount == 0 {
		return errors.New("bad request")
	}

	event := events.DepositFundEvent{
		ID:     command.ID,
		Amount: command.Amount,
	}

	log.Printf("%#v", event)
	return sv.eventProducer.Produce(event)
}

func (sv accountService) WithdrawFund(command commands.WithdrawFundCommand) error {
	if command.ID == "" || command.Amount == 0 {
		return errors.New("bad request")
	}

	event := events.WithdrawFundEvent{
		ID:     command.ID,
		Amount: command.Amount,
	}

	log.Printf("%#v", event)
	return sv.eventProducer.Produce(event)
}

func (sv accountService) CloseAccount(command commands.CloseAccountCommand) error {
	if command.ID == "" {
		return errors.New("bad request")
	}

	event := events.CloseAccountEvent{
		ID: command.ID,
	}

	log.Printf("%#v", event)
	return sv.eventProducer.Produce(event)
}
