# System notes

## vulnerability

If system is abruptly shut down some queued tasks in go channels can be lost.

## how to solve

1. blue/green deployment: it can help to gracefully stop the process, while giving it less and less work, until no more items are present in channels.

2. listen to exit/kill signals: this can be used to close the chanel, so no more items can be added to it, and to process remaining items in channels before shut down.