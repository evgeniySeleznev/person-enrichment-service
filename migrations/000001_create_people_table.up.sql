CREATE TABLE IF NOT EXISTS people (
                                      person_id BIGSERIAL PRIMARY KEY,
                                      name TEXT NOT NULL,
                                      surname TEXT NOT NULL,
                                      patronymic TEXT,
                                      age INTEGER,
                                      gender TEXT,
                                      nationality TEXT
);