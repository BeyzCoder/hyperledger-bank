import React, { useEffect, useState } from 'react';
import { useParams, useLocation } from 'react-router-dom';

import './Chequing-Style.css';

export default function Chequing() {
	const { account_id } = useParams();
	const { balance } = useLocation().state || {};
	const [transactions, setTransactions] = useState([]);

	// Inside your React component or service file
	const fetchData = async () => {
		try {
			const response = await fetch(`http://localhost:8080/cheq/${account_id}`, {
				method: 'GET',
				headers: {
					'Content-Type': 'application/json',
				},
			});
			const data = await response.json();
			setTransactions(data.data);
		} catch (error) {
			console.error('Error fetching data:', error);
		}
	};

	// Call the function inside a useEffect or based on user interaction
	useEffect(() => {
		fetchData();
	}, []);

	return (
		<div className="view-cheq">
			<div className='transactions'>
				<h2>Day to Day Chequing Transaction:</h2>
				<h1>Balance: {balance.toLocaleString()}</h1>
				{transactions.map((transaction, index) => (
					<div className="row">
						<div className="col">
							<div className="hash">{transaction.TransactionHash}</div>
							<div className="time">{transaction.Timestamp}</div>
						</div>
						<div className='col'>
							<div className="amount"
									style={{
										color: transaction.Activity === 'deposit' || transaction.Activity === 'received' ? 'green' : 'red',
									}}
							>
								{transaction.Activity === 'deposit' || transaction.Activity === 'received' ? '+' : '-'}
								{transaction.Amount.toLocaleString()}
							</div>
							<div className="activity">{transaction.Activity}</div>
						</div>
					</div>
				))}
			</div>
		</div>
	);
}