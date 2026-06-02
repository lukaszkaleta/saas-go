CREATE TABLE financial_ledger (
    id UUID PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users(id),
    job_id bigint REFERENCES job(id),
    type TEXT NOT NULL,
    /*
        TYPES:
        - job_payment
        - platform_fee
        - stripe_fee
        - refund
        - payout
        - adjustment
    */
    direction TEXT NOT NULL,
    /*
        credit = money coming to seller
        debit = money leaving seller
    */
    gross_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    fee_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    net_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL,
    stripe_payment_intent_id TEXT,
    stripe_transfer_id TEXT,
    stripe_refund_id TEXT,
    occurred_at TIMESTAMP NOT NULL,
    action_created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    action_created_by_id bigint references users(id),
    metadata JSONB
);