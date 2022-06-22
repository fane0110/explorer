GO = go

TARGET = explorer-task

all:$(TARGET)

explorer-task: task.go
	$(GO) mod tidy
	$(GO) build -o $@ $^

clean:
	rm -rf $(TARGET)
