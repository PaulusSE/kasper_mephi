-- +goose Up
-- +goose StatementBegin
-- 1. Новый тип статуса заявки
CREATE TYPE request_status AS ENUM ('pending','approved','rejected');

-- 2. Таблица для хранения «заявок на регистрацию»
CREATE TABLE registration_requests (
                                       request_id     UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
                                       email          VARCHAR(128)   NOT NULL UNIQUE,
                                       password_hash  TEXT           NOT NULL,
                                       user_type      user_type      NOT NULL,       -- 'student' или 'supervisor'
                                       payload        JSONB          NOT NULL,       -- сюда можно класть остальные поля (full_name, group_id и т.д.)
                                       status         request_status NOT NULL DEFAULT 'pending',
                                       created_at     TIMESTAMPTZ    NOT NULL DEFAULT now(),
                                       processed_by   UUID,                            -- кто подтвердил/отклонил (user_id администратора/руководителя)
                                       processed_at   TIMESTAMPTZ                     -- когда подтвердил/отклонил
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS registration_requests;
DROP TYPE IF EXISTS request_status;
-- +goose StatementEnd
