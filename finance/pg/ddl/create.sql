CREATE TABLE financial_ledger (
    id BIGSERIAL PRIMARY KEY,
    -- seller = who earns money (DPI anchor)
    seller_id int64 NOT NULL REFERENCES users (id),
    -- buyer = who pays money
    buyer_id int64 REFERENCES users (id),
    -- related job
    job_id int64 references job(id),
    -- event type (defines meaning)
    type TEXT NOT NULL,
    /*
        escrow_hold
        payout_release
        platform_fee
        payout
        refund
        chargeback
        adjustment
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