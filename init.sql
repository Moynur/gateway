CREATE TABLE IF NOT EXISTS transactions
(
    id SERIAL PRIMARY KEY,
    transaction_id UUID NOT NULL,
    operation_id UUID NOT NULL,
    amount INT NOT NULL,
    amount_available INT NOT NULL,
    currency TEXT NOT NULL,
    operation_type TEXT NOT NULL,
    -- In reality pan should be stored in a PCI compliant way
    pan TEXT NOT NULL
)