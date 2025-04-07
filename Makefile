# AGENT
BINARY_NAME_AGENT=agent
BASE_PATH_AGENT=./cmd/agent/
OUTPUT_PATH_AGENT=$(BASE_PATH_AGENT)$(BINARY_NAME_AGENT)
LOG_PATH_AGENT=$(BASE_PATH_AGENT)$(BINARY_NAME_AGENT).log
START_CMD_AGENT=$(BASE_PATH_AGENT)$(BINARY_NAME_AGENT)
PID_FILE_AGENT=$(BASE_PATH_AGENT)$(BINARY_NAME_AGENT).pid

# SERVER
BINARY_NAME_SERVER=server
BASE_PATH_SERVER=./cmd/server/
OUTPUT_PATH_SERVER=$(BASE_PATH_SERVER)$(BINARY_NAME_SERVER)
LOG_PATH_SERVER=$(BASE_PATH_SERVER)$(BINARY_NAME_SERVER).log
START_CMD_SERVER=$(BASE_PATH_SERVER)$(BINARY_NAME_SERVER)
PID_FILE_SERVER=$(BASE_PATH_SERVER)$(BINARY_NAME_SERVER).pid

.PHONY: build-agent
build-agent:
	go build -o $(OUTPUT_PATH_AGENT) $(BASE_PATH_AGENT)

.PHONY: build-server
build-server:
	go build -o $(OUTPUT_PATH_SERVER) $(BASE_PATH_SERVER)

.PHONY: start-agent
start-agent: build-agent
	@echo "Starting $(BINARY_NAME_AGENT)..."
	@nohup $(START_CMD_AGENT) > $(LOG_PATH_AGENT) 2>&1 & echo $$! > $(PID_FILE_AGENT)
	@echo "$(BINARY_NAME_AGENT) started with PID `cat $(PID_FILE_AGENT)`"

.PHONY: start-server
start-server: build-server
	@echo "Starting $(BINARY_NAME_SERVER)..."
	@nohup $(START_CMD_SERVER) > $(LOG_PATH_SERVER) 2>&1 & echo $$! > $(PID_FILE_SERVER)
	@echo "$(BINARY_NAME_SERVER) started with PID `cat $(PID_FILE_SERVER)`"

.PHONY: start
start: start-server start-agent

.PHONY: stop-agent
stop-agent:
	@if [ -f $(PID_FILE_AGENT) ]; then \
		echo "Stopping $(BINARY_NAME_AGENT)..."; \
		kill -9 `cat $(PID_FILE_AGENT)` && rm -f $(PID_FILE_AGENT); \
		echo "$(BINARY_NAME_AGENT) stopped."; \
	else \
		echo "No PID (agent) file found. Is the service running?"; \
	fi

.PHONY: stop-server
stop-server:
	@if [ -f $(PID_FILE_SERVER) ]; then \
		echo "Stopping $(BINARY_NAME_SERVER)..."; \
		kill -9 `cat $(PID_FILE_SERVER)` && rm -f $(PID_FILE_SERVER); \
		echo "$(BINARY_NAME_SERVER) stopped."; \
	else \
		echo "No PID (server) file found. Is the service running?"; \
	fi

.PHONY: stop
stop: stop-agent stop-server

.PHONY: restart
restart: stop start

.PHONY: clean
clean:
	rm -f $(OUTPUT_PATH_AGENT) $(PID_FILE_AGENT) $(LOG_PATH_AGENT)
	rm -f $(OUTPUT_PATH_SERVER) $(PID_FILE_SERVER) $(LOG_PATH_SERVER)