# 02. Producer и consumer

TODO-02/03 в internal/messaging/kafka.go. Publisher.Publish(ctx,event) сериализует JSON,
выбирает key из DocumentID и возвращает ошибку синхронной публикации в пределах deadline.
Consume(ctx,brokers,topic,group,eventType,handle) читает нужную group и фильтрует типы.

Порядок: Fetch → decode/filter → handle → Commit. Временная ошибка handler или output
publish не означает успешную обработку. Cancellation завершает loop и закрывает reader.
Выберите и опишите политику invalid payload: без бесконечного скрытого цикла.

Пример: document.created с document_id=7 читает worker, document.processed — owner API.
Не используйте одну application group для независимого console observer.
Проверьте JSON/key, failure before commit, retry того же record, shutdown, type filtering.
Подсказка: ReadMessage и FetchMessage имеют разные commit consequences при group reader.
Сравните classroom consumer, но сохраняйте Document contract.
