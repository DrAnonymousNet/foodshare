run:
	go build && ./foodshare

kill:
	kill $(lsof -t -i :8000)

