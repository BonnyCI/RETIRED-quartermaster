package bot

type Bot interface {
	Configure()
	Connect()
	Disconnect()
	Send()
	Sendf()
}
