-- Удаляем дополнительные индексы
DROP INDEX IF EXISTS idx_lots_current_price;
DROP INDEX IF EXISTS idx_bids_amount;
DROP INDEX IF EXISTS idx_lots_created_at;