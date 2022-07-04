GO = go

TARGET = explorer-task

TARGET2 = explorer-task4

TARGET3 = task-export-ex

TARGET4 = task-import


all:$(TARGET) $(TARGET2) $(TARGET3) $(TARGET4) 


explorer-task: task.go
	$(GO) mod tidy
	$(GO) build -o $@ $^

task-export: task-export.go
	$(GO) mod tidy
	$(GO) build -o $@ $^

task-export-ex: task-export-ex.go
	$(GO) mod tidy
	$(GO) build -o $@ $^

task-import: task-import.go
	$(GO) mod tidy
	$(GO) build -o $@ $^


clean:
	rm -rf $(TARGET)
