-- Добавляем поле для длительности аукциона (в секундах)
ALTER TABLE lots ADD COLUMN duration_seconds INTEGER;

-- Обновляем существующие записи
UPDATE lots SET duration_seconds = EXTRACT(EPOCH FROM (end_time_unix - created_at)) 
WHERE duration_seconds IS NULL;