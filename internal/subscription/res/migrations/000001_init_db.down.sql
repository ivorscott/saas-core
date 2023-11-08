REASSIGN OWNED BY user_a TO postgres;
DROP OWNED BY user_a;
DROP USER user_a;

DROP POLICY customers_isolation_policy ON customers;
DROP POLICY subscriptions_isolation_policy ON subscriptions;
DROP POLICY transactions_isolation_policy ON transactions;

ALTER TABLE customers DISABLE ROW LEVEL SECURITY;
ALTER TABLE subscriptions DISABLE ROW LEVEL SECURITY;
ALTER TABLE transactions DISABLE ROW LEVEL SECURITY;

DROP TABLE subscription_statuses;

DROP TABLE transaction_statuses;

DROP INDEX idx_transaction_status_id;

DROP TABLE transactions;

DROP INDEX idx_subscription_transaction, idx_subscription_status_id, idx_subscription_customer;

DROP TABLE IF EXISTS subscriptions;

DROP INDEX idx_customer_surname, idx_customer_email;

DROP TABLE IF EXISTS customers;