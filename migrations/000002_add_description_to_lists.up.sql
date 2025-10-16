-- Добавляем новое поле description в таблицу lists
ALTER TABLE lists ADD COLUMN description TEXT;

-- Обновление таблицы: для всех существующих записей
UPDATE lists SET description = '';

-- Комментарии для документации
COMMENT ON COLUMN lists.description IS 'Описание списка (необязательное поле)';