# 03. Обработка Document через application ports

TODO-04: Processor.Handle(ctx,document.Event) error. На document.created посчитайте слова
в title, создайте document.processed и опубликуйте результат. Другие типы не должны
создавать цикл событий. Worker не обращается к DocumentRepository или PostgreSQL.

До работы проверьте receipt по event ID. После успешного Publish сохраните receipt на
10 минут. При ошибке Publish receipt сохранять нельзя. Ошибки зависимостей возвращаются
consumer. Пример: title "Read Go docs" → Words=3, тот же DocumentID, непустой result ID.

Проверьте один/несколько/Unicode words, duplicate receipt, publisher failure, receipt error,
caller cancellation. Acceptance tests используют fakes и не требуют Kafka/Redis.
Подсказка: receipt и publication не составляют общей transaction. Объясните окно duplicate
при падении после Publish до Save и после expiry. Exactly-once обещать нельзя.
