INSERT INTO stages(name, sort) VALUES
 ('New', 10),
 ('Contacted', 20),
 ('Appointment Set', 30),
 ('Active Client', 40),
 ('Under Contract', 50),
 ('Closed', 60)
ON CONFLICT DO NOTHING;
