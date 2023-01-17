install:
	@echo "Start installation..."
	@(cd src; go mod tidy; sudo go build -a -o /usr/local/bin/popularity_survey)
	@echo "Installation is completed."

uninstall:
	@echo "Start uninstallation..."
	@which popularity_survey | sudo xargs rm
	@echo "Uninstallation is completed."

test:
	@echo "Start Test..."
	@(cd src; go test ./...)
