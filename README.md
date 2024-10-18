# Picpay Backend Challenge

Application that implements picpay (brazilian company) software engineering interview challenge

pt-br system specification can be checked at [specs.md](spec.md)

## Run app

```sh
make run
```

### Todos

- Include transfer datetime in database


## System notes

### Banking Race conditions problem

The system is using a goroutine and a channel to process each transfer request atomically, to avoid race conditions that could cause the double spending problem.

An account could spend money more than once while the system don't persist the first debit.

The problem with this aproach is that, no metter the account, ALL transactions are processed in a queue: one after the other.

In a real world scenario, user base growth could make transaction processing extremely slow.

**How I woult try to solve**

1. Use a message broker like kafka to process parallel transfer if this 'same time' transfers are related to different accounts. Transfer events could be published to a customer specific partition, so 'same customer' transfers would be processed in a queue fifo way and different customer transactions could be processed in the same time 


### Goroutines and channels problem

If system is abruptly shut down some queued tasks in go channels can be lost.

**How I woult try to solve**

1. blue/green deployment: it can help to gracefully stop the process, while giving it less and less work, until no more items are present in channels.

2. listen to exit/kill signals: this can be used to close the chanel, so no more items can be added to it, and to process remaining items in channels before shut down.

### Database ATOMIC transactions problem

System is storing two transaction registers per transfer: a debit to the payer and a credit to the peyee

It is recommend to do it in a sigle DATABASE TRANSACTION wich will persist both registers and rollback all when error occurs

The system is currently vulnerable to data corruption in error scenarios

**How I woult try to solve**

1. create a MONGO DATABASE TRANSACTION to persist all data required to register a transfer
2. change database strategy and use a relational DB that is ACID compliant