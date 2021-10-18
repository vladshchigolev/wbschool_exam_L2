package main

import "fmt"
// Фасад — это структурный паттерн, который предоставляет простой (но урезанный) интерфейс
// к сложной системе объектов, библиотеке или фреймворку.
// ---------------------------------------------------------------------------------------
// Процессы, происходящие за кулисами, когда производятся операции с балансом банковского счёта (кредитной карты),
// очень сложны. В этом процесе участвуют десятки подсистем. Вот их малая часть:
// Проверка аккаунта
// Проверка PIN-кода
// Баланс дебет/кредит
// Запись в бухгалтерской книге
// Отправка оповещения
// В такой сложной системе легко потеряться или что-то сломать, если обращаться с ней неправильно.
// Для таких случаев и существует паттерн Фасад — он позволяет клиенту работать с десятками компонентов,
// используя при этом простой интерфейс. Клиенту необходимо лишь ввести реквизиты карты, код безопасности,
// стоимость оплаты и тип операции. Фасад управляет дальнейшей коммуникацией между различными компонентами
// без контакта клиента со сложными внутренними механизмами.
type walletFacade struct {
	account      *account // Внутри фасада находятся различные компоненты, к которым у клиента нет прямого доступа,
	wallet       *wallet // вместо этого - упрощенный интерфейс взаимодействия с ними
	securityCode *securityCode
	notification *notification
	ledger       *ledger
}

func newWalletFacade(accountID string, code int) *walletFacade {
	fmt.Println("Starting create account")
	walletFacacde := &walletFacade{
		account:      newAccount(accountID),
		securityCode: newSecurityCode(code),
		wallet:       newWallet(),
		notification: &notification{},
		ledger:       &ledger{},
	}
	fmt.Println("Account created")
	return walletFacacde
}
// У кошелька простой интерфейс - 1. метод для пополнения счёта, 2. метод для списания суммы со счёта
func (w *walletFacade) addMoneyToWallet(accountID string, securityCode int, amount int) error {
	fmt.Println("Starting add money to wallet")
	err := w.account.checkAccount(accountID)
	if err != nil {
		return err
	}
	err = w.securityCode.checkCode(securityCode)
	if err != nil {
		return err
	}
	w.wallet.creditBalance(amount)
	w.notification.sendWalletCreditNotification()
	w.ledger.makeEntry(accountID, "credit", amount)
	return nil
}

func (w *walletFacade) deductMoneyFromWallet(accountID string, securityCode int, amount int) error {
	fmt.Println("Starting debit money from wallet")
	err := w.account.checkAccount(accountID)
	if err != nil {
		return err
	}

	err = w.securityCode.checkCode(securityCode)
	if err != nil {
		return err
	}
	err = w.wallet.debitBalance(amount)
	if err != nil {
		return err
	}
	w.notification.sendWalletDebitNotification()
	w.ledger.makeEntry(accountID, "credit", amount)
	return nil
}
