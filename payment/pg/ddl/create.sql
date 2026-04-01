CREATE SEQUENCE pay_sequence;

CREATE TABLE pay_payment_intent (
  id BIGINT PRIMARY KEY DEFAULT nextval('pay_sequence'),
  reference TEXT NOT NULL UNIQUE,
  stripe_payment_intent_id TEXT UNIQUE,
  stripe_client_secret TEXT,
  stripe_session_url TEXT,
  job_id BIGINT NOT NULL,
  payer_id BIGINT NOT NULL,
  payee_id BIGINT NOT NULL CHECK (
    payer_id <> payee_id
  ),
  amount BIGINT NOT NULL CHECK (amount > 0),
  currency TEXT NOT NULL DEFAULT 'NOK',
  status TEXT NOT NULL CHECK (
    status IN ('CREATED', 'INITIATED', 'SUCCEEDED', 'FAILED')
  ),
  action_created_by_id BIGINT NOT NULL,
  action_created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  action_updated_by_id BIGINT,
  action_updated_at TIMESTAMPTZ
);

CREATE INDEX pay_payment_intent_job_id_idx ON pay_payment_intent(job_id);
CREATE INDEX pay_payment_intent_payer_id_idx ON pay_payment_intent(payer_id);
CREATE INDEX pay_payment_intent_payee_id_idx ON pay_payment_intent(payee_id);