import React, {Component} from 'react';
import './App.css';
import Navbar from 'react-bootstrap/Navbar'
import logo from "./img/logo.png"
import Button from "react-bootstrap/Button";

class NonAuthenticatedApp extends Component {

    render() {
        return (
            <div className="App">
                <div className="App-navbar-wrapper">
                    <Navbar className="App-navbar">
                        <Navbar.Brand href="/">
                            <img src={logo} alt="logo"/>
                            Pull Assistant
                        </Navbar.Brand>
                        <Navbar.Collapse className="justify-content-end">
                            <Button href={process.env.REACT_APP_API_URL + "/auth/github"}
                                    variant="outline-secondary">Sign in</Button>
                        </Navbar.Collapse>
                    </Navbar>
                </div>
                <div className="App-divider"/>
                <div className="App-content">
                    <h2>Pull Assistant installations</h2>
                    <div className="App-content-nav">
                        <Button href={process.env.REACT_APP_API_URL + "/auth/github"} variant="primary">
                            Sign in
                        </Button>
                    </div>
                </div>
            </div>
        );
    }
}

export default NonAuthenticatedApp;