GO = go

TARGET = explorer-task

TARGET2 = explorer-task4

TARGET3 = explorer-task5


all:$(TARGET) $(TARGET2) $(TARGET3) 


explorer-task: task.go
	$(GO) mod tidy
	$(GO) build -o $@ $^

explorer-task4: task4.go
	$(GO) mod tidy
	$(GO) build -o $@ $^

explorer-task5: task5.go
	$(GO) mod tidy
	$(GO) build -o $@ $^


clean:
	rm -rf $(TARGET)
