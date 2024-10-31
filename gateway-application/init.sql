-- Create 'accounts' table
CREATE TABLE IF NOT EXISTS accounts (
    account_id INT PRIMARY KEY, -- Manually managed account ID
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    address TEXT NOT NULL,
    sin VARCHAR(20) NOT NULL, -- Social Insurance Number
    email VARCHAR(100) NOT NULL, -- Email address
    password VARCHAR(255) NOT NULL -- Password (consider hashing this before storing)
);

-- Create 'transactionsHash' table
CREATE TABLE IF NOT EXISTS transaction_hashes (
    transaction_hash VARCHAR(255), 
    activity VARCHAR(20) CHECK (activity IN ('deposit', 'withdraw', 'payment', 'received')), -- Limit type to specific values
    amount NUMERIC(15, 2) NOT NULL, -- Amount involved in the transaction, allowing up to 15 digits, with 2 decimals
    account_id INT REFERENCES accounts(account_id), -- Foreign key to the 'accounts' table
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Timestamp of the transaction
);

INSERT INTO accounts (account_id, name, phone, address, sin, email, password)
VALUES (123456789, 'Steven Kline Baes', '333-777-9999', '123 Maple St', 'SIN123456', 'baes2374@saskpolytech.ca', 'psswd');
INSERT INTO accounts (account_id, name, phone, address, sin, email, password)
VALUES (987654321, 'National Student Loans Service Centre', '333-999-1111', '124 Maple St', 'SIN654321' , 'nslsc@workemail.ca', 'psswd');
INSERT INTO accounts (account_id, name, phone, address, sin, email, password)
VALUES (098765432, 'SaskEnergy', '333-888-2222', '125 Maple St', 'SIN765432', 'saskenergy@workemail.ca', 'psswd');
