# 05. Workflow, tests и сдача

TODO-07: завершите Make targets worker/client для exercise profile. Пока targets завершаются
с ошибкой, чтобы незавершённая команда не выглядела успешным запуском.
Client — одноразовый GetDocument с DOCUMENT_ID, worker — долгоживущий отдельный process.
Не публикуйте credentials и не удаляйте volumes при обычном down.

TODO-08: расширьте UseCase tests: successful Create, validation before dependencies,
repository/publisher failure, context propagation. Добавьте consumer tests и failure cases worker.
Основной make ci должен проходить. make acceptance и acceptance-redis после решения тоже зелёные.
Сейчас baseline компилируется, функциональные acceptance tests ожидаемо красные из-за TODO.

Сдача: GitHub PR, команды запуска, результаты tests/vet/race, короткая схема responsibility.
Покажите HTTP Create → pending → Kafka → worker → processed → owner API → ready → gRPC Get.
Отдельно объясните временные receipts, retention, duplicate window и dual write.
Outbox, DLQ и production service authentication — необязательное углубление, не скрытый критерий.
