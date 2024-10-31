import React, { useRef, useState } from 'react';

import './RegisterForm-Style.css';

export default function RegisterForm() {
	const formRef = useRef();

	const [formData, setFormData] = useState({
		first: '',
		last: '',
		phone: '',
		email: '',
		message: '',
	});

	const handleChange = (e) => {
		const { name, value } = e.target;
		setFormData({
			...formData,
			[name]: value,
		});
	};

	const handleSubmit = (e) => {
		e.preventDefault();
	}

  return (
		<div className='container-form'>
			<h2>Register Account</h2>
			<form ref={formRef} onSubmit={handleSubmit}>
				<div className="full-name">
					<input type="text" id='first' name='first' value={formData.first} onChange={handleChange} placeholder='First Name' />
					<input type="text" id='last' name='last' value={formData.last} onChange={handleChange} placeholder='Last Name' />
				</div>
				<div className='contact-info'>
					<input type="tel" id='phone' name='phone' value={formData.phone} onChange={handleChange} placeholder='Phone Number' />
					<input type="email" id='email' name='email' value={formData.email} onChange={handleChange} placeholder='Email' />
				</div>
				<div className='personal-detail'>
					<input type="text" id='address' name='address' value={formData.address} onChange={handleChange} placeholder='Address' />
					<input type="text" id='sin' name='sin' value={formData.sin} onChange={handleChange} placeholder='S I N' />
				</div>
				<button type='submit' value='Send'>Submit</button>
			</form>
		</div>
	);
}