import React, { useRef, useState } from 'react';
import { useLocation } from 'react-router-dom';

import './Payment-Style.css';

export default function Payment() {
	const formRef = useRef();
	const { account_id } = useLocation().state || {};

	const [formData, setFormData] = useState({
		From: '',
		To: '',
		Amount: '',
	});

	const handleChange = (e) => {
		const { name, value } = e.target;
		setFormData({
			...formData,
			[name]: value,
		});
	};

	const handleSubmit = async (e) => {
		e.preventDefault();
		try {
			alert(`transaction processing!`);
			const response = await fetch(`http://localhost:8080/payment`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify(formData),
			});

			if (response.ok) {
				const data = await response.json();
				alert(data.data);
			} else {
				console.log("Error:", response.statusText);
			}
		} catch (error) {
			console.error('Error fetching data:', error);
		}
	}

	return (
		<div className="view-pay">
			<div className="payment-container">
				<h1>Send Payment:</h1>
				<form ref={formRef} onSubmit={handleSubmit}>
					<select name="From" id="From" value={formData.From} onChange={handleChange}>
						<option value="" disabled>Select an account</option>
						<option value={account_id}>Chequing</option>
					</select>
					<input type="To" id='To' name='To' value={formData.To} onChange={handleChange} placeholder='To Whom' />
					<input type="Amount" id='Amount' name='Amount' value={formData.Amount} onChange={handleChange} placeholder='Amount to send' />
					<button type='submit' value='Send'>Send</button>
				</form>	
			</div>
		</div>
	);
}