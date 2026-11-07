-- +goose Up
-- +goose StatementBegin
CREATE TABLE documents (
 id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
 title TEXT NOT NULL CHECK (btrim(title)<>''),
 status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending','ready'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE documents;
-- +goose StatementEnd
