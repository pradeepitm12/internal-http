CREATE TABLE IF NOT EXISTS accounts (
                                        account_id INTEGER PRIMARY KEY,
                                        balance NUMERIC(20, 2) NOT NULL
    );

CREATE TABLE IF NOT EXISTS transactions (
                                            transaction_id SERIAL PRIMARY KEY,
                                            source_account_id INTEGER NOT NULL,
                                            destination_account_id INTEGER NOT NULL,
                                            amount NUMERIC(20, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );