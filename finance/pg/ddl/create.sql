CREATE TABLE financial_ledger (
    id BIGSERIAL PRIMARY KEY,
    -- seller = who earns money (DPI anchor)
    seller_id BIGINT NOT NULL REFERENCES users (id),
    -- buyer = who pays money
    buyer_id BIGINT REFERENCES users (id),
    -- related job
    job_id BIGINT references job(id),
    -- event type (defines meaning)
    type SMALLINT NOT NULL,
    /*
        1 = escrow_hold        -> customer paid
        2 = payout_release     -> SELLER EARNED (DAC7 SOURCE OF TRUTH)
        3 = platform_fee       -> revenue
        4 = payout             -> Stripe transfer only (ignore for DAC7)
        5 = refund
        6 = chargeback
        7 = adjustment
    */

    -- MONEY (minor units, e.g. 100 NOK = 10000)
    amount BIGINT NOT NULL,

    -- optional breakdown for reporting
    gross_amount BIGINT,
    fee_amount BIGINT,
    net_amount BIGINT,

    currency CHAR(3) NOT NULL,

    -- Stripe reconciliation
    stripe_payment_intent_id TEXT,
    stripe_transfer_id TEXT,
    stripe_payout_id TEXT,
    stripe_refund_id TEXT,

    -- when it actually happened
    occurred_at TIMESTAMP NOT NULL,

    -- DB insert time
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    -- flexible metadata
    metadata JSONB
);

CREATE UNIQUE INDEX financial_ledger_payment_intent_unique_index
    ON financial_ledger(stripe_payment_intent_id, type)
    WHERE stripe_payment_intent_id IS NOT NULL;

CREATE UNIQUE INDEX financial_ledger_transfer_unique_index
    ON financial_ledger(stripe_transfer_id, type)
    WHERE stripe_transfer_id IS NOT NULL;

CREATE UNIQUE INDEX financial_ledger_payout_unique_index
    ON financial_ledger(stripe_payout_id, type)
    WHERE stripe_payout_id IS NOT NULL;

CREATE UNIQUE INDEX financial_ledger_refund_unique_index
    ON financial_ledger(stripe_refund_id, type)
    WHERE stripe_refund_id IS NOT NULL;

CREATE INDEX financial_ledger_seller_occurred_index
    ON financial_ledger(seller_id, occurred_at);

CREATE INDEX financial_ledger_job_index
    ON financial_ledger(job_id, occurred_at);

CREATE INDEX financial_ledger_buyer_occurred_index
    ON financial_ledger(buyer_id, occurred_at);

CREATE INDEX financial_ledger_type_occurred_index
    ON financial_ledger(type, occurred_at);

CREATE INDEX financial_ledger_occurred_index
    ON financial_ledger(occurred_at);