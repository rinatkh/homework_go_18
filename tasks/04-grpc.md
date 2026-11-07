# 04. gRPC GetDocument и typed client

TODO-05 в internal/grpcapi/server.go: GetDocument(ctx,request) вызывает прежний UseCase.Get,
проверяет положительный id и преобразует Document в DocumentResponse.
Ошибки: invalid id → InvalidArgument, отсутствующий Document → NotFound,
неожиданная repository error → нейтральный Internal, без SQL details.

TODO-06 в internal/automation/client.go: Get(ctx,client,id) выполняет typed RPC с верхним
пределом три секунды, сохраняет более короткий parent deadline/cancellation и вызывает cancel.
Пример: id=5 → тот же ID/title/status из PostgreSQL через API, без прямого DB access client.

Сначала прочитайте api/proto/document/v1/document.proto, затем pkg/proto/document/v1.
Не редактируйте generated code. make tools и make proto воспроизводят output.
Подсказка: protobuf DTO не становится application model. Для tests достаточно fake repository
и fake generated client; отдельный bufconn experiment можно добавить после unit tests.
