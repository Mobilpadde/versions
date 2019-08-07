PORT ?= 5000

kill:
	kill $$(lsof -t -i:$(PORT))