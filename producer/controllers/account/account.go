package accountcontrollers

import (
	"encoding/json"
	"log"
	"producer/commands"
	services "producer/services/account"

	"github.com/gofiber/fiber/v2"
)

type IAccountController interface {
	OpenAccount(c *fiber.Ctx) error
	DepositFund(c *fiber.Ctx) error
	WithdrawFund(c *fiber.Ctx) error
	CloseAccount(c *fiber.Ctx) error
}

type accountController struct {
	accountService services.IAccountService
}

func NewAccountController(accountService services.IAccountService) IAccountController {
	return accountController{accountService}
}

func (obj accountController) OpenAccount(c *fiber.Ctx) error {
	command := commands.OpenAccountCommand{}
	err := json.Unmarshal(c.Body(), &command)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.SendString(err.Error())
	}

	id, err := obj.accountService.OpenAccount(command)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return err
	}

	c.Status(fiber.StatusCreated)
	return c.JSON(fiber.Map{
		"message": "open account success",
		"id":      id,
	})
}

func (obj accountController) DepositFund(c *fiber.Ctx) error {
	command := commands.DepositFundCommand{}
	err := json.Unmarshal(c.Body(), &command)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.SendString(err.Error())
	}

	err = obj.accountService.DepositFund(command)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"message": "deposit fund success",
	})
}

func (obj accountController) WithdrawFund(c *fiber.Ctx) error {
	command := commands.WithdrawFundCommand{}
	err := c.BodyParser(&command)
	if err != nil {
		return err
	}

	err = obj.accountService.WithdrawFund(command)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(fiber.Map{
		"message": "withdraw fund success",
	})
}

func (obj accountController) CloseAccount(c *fiber.Ctx) error {
	command := commands.CloseAccountCommand{}
	err := c.BodyParser(&command)
	if err != nil {
		return err
	}

	err = obj.accountService.CloseAccount(command)
	if err != nil {
		log.Println(err)
		return err
	}

	return c.JSON(fiber.Map{
		"message": "close account success",
	})
}
