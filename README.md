# Homework Go 18. Document processing starter

Перенесите инфраструктурные идеи урока на небольшой **Document API**. Это новый Homework
starter, предыдущие Homework10–15 не меняются. PostgreSQL хранит documents, Redis —
временные receipts обработанных event IDs, Kafka — document.created/processed,
gRPC — GetDocument. Сложная auth здесь не повторяется, локальные ports loopback-only.

Дано: application skeleton, ports, рабочий PostgreSQL repository и SQL schema,
Compose с PG/Redis/Kafka, multi-stage Dockerfile, отдельные cmd API/worker/client,
исходный .proto и **реально сгенерированный** Go contract в pkg/proto.
Kafka writer creation/result wiring подготовлены. Задание — закончить adapters и flow.

## Запуск baseline

Нужны Docker/Compose, Go1.22, make. `make help`, `make env`, `make up`, затем
`curl -i http://localhost:8082/health`. Одноразовый migrate job применяет goose 001_documents.sql (Up/Down в одном файле).
PostgreSQL5433, Redis6380, Kafka19094, HTTP8082 и gRPC50052 не конфликтуют с classroom.
Redis временный, PG/Kafka volumes сохраняются. down -v уничтожает данные, обычный down — нет.

`make ci` и `make check` проходят **до решения**, проверяет компиляцию и базовую validation.
Это не означает, что упражнения выполнены. POST /documents до TODO-02 сохраняет Document,
но возвращает 501: publisher ещё не готов, это тот же dual-write gap. GetDocument gRPC
возвращает Unimplemented до TODO-05. Не повторяйте POST бесконечно без проверки SQL.
Worker/client не запускаются автоматически, их exercise profile включается после TODO.

## TODO checklist

| ID | Файл | Ожидаемая работа |
| --- | --- | --- |
| 01 | internal/repository/redis/receipt.go | Exists/Save c doc:receipt:key и positive TTL, передавать ctx |
| 02 | internal/messaging/kafka.go | Producer JSON, key=document_id, bounded publish + ack |
| 03 | internal/messaging/kafka.go | Group reader Fetch→handler→Commit, filtering/retry/cancel |
| 04 | internal/worker/processor.go | Created→анализ title→Processed, receipt после успешного publish |
| 05 | internal/grpcapi/server.go | GetDocument, DTO mapping, InvalidArgument/NotFound/Internal |
| 06 | internal/automation/client.go | Typed gRPC Get с 3s deadline и caller context |
| 07 | Makefile worker/client | Реальные commands для exercise profile, document ID configurable |
| 08 | internal/document/usecase_test.go и новые tests | Successful create/failures/worker через fakes |

Receipt фиксируется после успешного output publish. Crash между publish и Save может
породить duplicate output, TTL тоже ограничивает dedup window. MarkReady в owner service
идемпотентен. Не обещайте exactly once и не пишите SQL из worker.
Проверьте, что raw event JSON валиден и event IDs не пустые, transient error не коммитит
входной offset. Producer остаётся за Publisher interface, UseCase не импортирует Kafka.

## Acceptance

Создать Document через HTTP. Увидеть pending в PG и created в Kafka. Отдельный worker
анализирует title, публикует processed. API consumer меняет status на ready. Отдельный
client вызывает gRPC GetDocument и получает ready. Redis receipt имеет TTL10min.
Повтор события безопасен, publish failure не считается обработанным входом.

Для ручного Create: `curl -X POST localhost:8082/documents -H 'Content-Type: application/json' -d '{"title":"Read distributed systems"}'`.
Проверьте также invalid title →400 и отсутствие публикации события при validation error.

`make tools`, compiler protoc3.21.12, `make proto` восстанавливают оба pb.go из source.
Generated output не редактировать. Source располагается api/proto, shared output pkg/proto.
Для HTTP/gRPC errors/TTL/failure paths напишите meaningful tests, дополнительно real Redis
integration с изолированным key prefix. Все cmd должны завершаться при cancellation.

В сдаче: patch своего решения, короткий architecture sketch, make check results,
команды воспроизведения и объяснение dual write/duplicates. Бонус: retry budget/DLQ или
outbox design на бумаге, без обязательного нового infrastructure компонента.


## Проверки задания

`make acceptance` запускает functional tests worker, gRPC server и typed client.
До реализации TODO эти tests обязаны падать, без t.Skip и без ожидания ErrNotImplemented.
`make acceptance-redis` проверяет настоящий Redis receipt/TTL (TEST_REDIS_URL).
Компиляция всех тестов: go test -tags=acceptance,integration -run '^$' ./...
Основной CI запускает make ci и остаётся зелёным на starter.

Миграции: make migrate-up/down/status для Compose, make migrate-local-up/down/status
для DATABASE_URL. Повторный Up не повторяет DDL, один Down удаляет documents только в test DB.
Реализации exercises не включены. Проект не импортирует classroom18.

Условия: [Redis](tasks/01-redis.md), [Kafka](tasks/02-kafka.md),
[worker](tasks/03-worker.md), [gRPC](tasks/04-grpc.md), [workflow/tests](tasks/05-workflow-tests.md).
