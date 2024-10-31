import React, { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';

import './Account-Style.css';

import saskatoonCity from '../images/saskatoon-city.jpg';

export default function Account() {
	const { account_id } = useParams();
	const [account, setAccount] = useState({
		AccountID: '',
		Owner: '',
		Balance: 0
	});

	// Inside your React component or service file
	const fetchData = async () => {
		try {
			const response = await fetch(`http://localhost:8080/account/${account_id}`, {
				method: 'GET',
				headers: {
					'Content-Type': 'application/json',
				},
			});
			const data = await response.json();
			console.log(data.data);
			setAccount({ AccountID: data.data.AccountID, Owner: data.data.Owner, Balance: data.data.Balance });
		} catch (error) {
			console.error('Error fetching data:', error);
		}
	};

	// Call the function inside a useEffect or based on user interaction
	useEffect(() => {
		fetchData();
	}, []);  
	

	return (
		<section className="account-page">
			<div className='header'>
				<img src={saskatoonCity} alt="background-head" />
				<div className='display-name'>{account.Owner}</div>	
			</div>
			<div className='summary'>
				<div className='banking-accounts'>
					<div className='section-account'>
						<h2>Bank Accounts</h2>
						<div className='open-account'>
							<div className='transaction-account'>
								<span className='account'><Link to={`/cheq/${account.AccountID}`} className='link' state={{balance: account.Balance}}>Day to Day Chequing Transaction:</Link></span>
								<span className='balance'>$ {account.Balance.toLocaleString()}</span>
							</div>
							<div className='detail-account'>
								<span className='accn-num'>Account Num: {account.AccountID}</span>
								<span className='currency'>CAD</span>
							</div>
						</div>
						<div className='apply-account'>
							<div className='add-account'>+</div>
							<span className='apply'>Apply for a Bank Account</span>
						</div>
					</div>
					<div className='section-account'>
						<h2>Credit Cards</h2>
						<div className='apply-account'>
							<div className='add-account'>+</div>
							<span className='apply'>Apply for a Credit</span>
						</div>
					</div>
					<div className='section-account'>
						<h2>Loans</h2>
						<div className='apply-account'>
							<div className='add-account'>+</div>
							<span className='apply'>Apply for a Loans</span>
						</div>
					</div>
					<div className='section-account'>
						<h2>Direct Investments</h2>
						<div className='apply-account'>
							<div className='add-account'>+</div>
							<span className='apply'>Apply for a Investment</span>
						</div>
					</div>
					<div className='section-account'>
						<h2>Mortgages</h2>
						<div className='apply-account'>
							<div className='add-account'>+</div>
							<span className='apply'>Apply for a Mortgage</span>
						</div>
					</div>
				</div>
				<div className='banking-actions'>
					<h2>Transfer & Payment Process</h2>
					<span><Link to={`/payment`} className='link' state={{account_id: account.AccountID}}>Send Payment</Link></span>
					<span><Link to={`/deposit`} className='link' state={{account_id: account.AccountID}}>Deposit Money</Link></span>
					<span><Link to={`/withdraw`} className='link' state={{account_id: account.AccountID}}>Withdraw Money</Link></span>
				</div>
			</div>
		</section>
	);
}