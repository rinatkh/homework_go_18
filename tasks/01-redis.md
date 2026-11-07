# 01. Временные receipts в Redis

Реализуйте TODO-01 в internal/repository/redis/receipt.go: Exists(ctx,id) (bool,error)
и Save(ctx,id,ttl) error. Ключ doc:receipt:<event_id>, TTL положительный.
Receipt означает, что result успешно опубликован, а не что он уже применён PostgreSQL.

Контракт: неизвестный ID → false,nil, сохранённый → true,nil до expiry, после expiry → false.
Пустой ID и неположительный TTL нужно отклонить. Не теряйте caller context и Redis errors.
Пример: Save("event-42",10min), затем Exists("event-42") возвращает true.

Проверьте отсутствие, сохранение, TTL, expiry, отменённый context, ошибку подключения.
Готовый real-Redis контракт запускается через make acceptance-redis с TEST_REDIS_URL.
До реализации он красный. Подсказка: сравните семантику EXISTS и SET с expiration,
не используйте бессрочную запись с отдельной необязательной EXPIRE.
