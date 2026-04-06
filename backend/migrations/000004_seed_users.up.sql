INSERT INTO users (id, email, password_hash, role, status, created_at, updated_at)
VALUES
  (gen_random_uuid(), 'admin@synth.cmd', '$2a$10$Q9XfE5sOnr4veNs/kk5nh.NaCcC3b.QsuElp0WCkz4ZbM2q6J92nq', 'admin', 'active', now(), now()),
  (gen_random_uuid(), 'agent@synth.cmd', '$2a$10$CkybTT8Jm7QQmOrVru/HyuE0yFBUdAb24MjDlBEAv5nAdCKIio0y6', 'agent', 'active', now(), now()),
  (gen_random_uuid(), 'client@synth.cmd', '$2a$10$aKg9Hl.DuiFZV2yP9og5DetJxLeF2nR4A4DzNEX8ZRl/VLgTfa3R.', 'client', 'active', now(), now())
ON CONFLICT (email) DO NOTHING;
