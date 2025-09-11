-- Дополнительные индексы для часто используемых запросов
CREATE INDEX idx_lots_current_price ON lots(current_price);
CREATE INDEX idx_bids_amount ON bids(amount);
CREATE INDEX idx_lots_created_at ON lots(created_at);