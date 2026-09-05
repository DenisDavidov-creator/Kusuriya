ALTER TABLE remedies
ADD CONSTRAINT uq_remedy_name UNIQUE (name);