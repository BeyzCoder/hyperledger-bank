import React, { useState } from 'react';
import { Link } from 'react-router-dom';

import RegisterForm from '../components/RegisterForm';

import './Home-Style.css';

import hyperledgerLogo from '../images/hyperledger-logo.png';
import bankLogo from '../images/bank-logo.png';
import closeIcon from '../assets/close-icon.svg';

export default function Home() {

	const [isPanelOpen, setIsPanelOpen] = useState(false);

  const handleButtonClick = () => {
    setIsPanelOpen(!isPanelOpen); // Toggle the panel state
  };

	return (
		<section className="home-page">
			<img src={bankLogo} alt="logo" className='home-logo'/>
			<div className='login-container'><Link to={`/account/123456789`} className='link'>Login</Link></div>
			<div className={`wrapper ${isPanelOpen ? 'slide-out' : ''}`}>
				<div className="headline">
					<h1>Hyperledger Fabric Banking</h1>
					<p>
						Lorem ipsum dolor sit, amet consectetur adipisicing elit. Esse quam fugit dicta. 
						Non debitis repellendus delectus placeat magnam, repudiandae obcaecati voluptatibus! 
						Laborum pariatur neque sint laudantium ea dicta excepturi repudiandae.
					</p>
					<button onClick={handleButtonClick}>Open an Account</button>
				</div>
				<div className="display-img">
					<img src={hyperledgerLogo} alt="logo" />
				</div>
			</div>

			{/* Panel Slides */}
			<div className={`side-panel ${isPanelOpen ? 'open' : ''}`}>
				<img src={closeIcon} alt="icon" onClick={handleButtonClick} className="close-panel"/>
				<RegisterForm />
			</div>
		</section>
	);
}