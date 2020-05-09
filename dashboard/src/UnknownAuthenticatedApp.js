import React, {Component} from 'react';
import './App.css';
import Navbar from 'react-bootstrap/Navbar'
import Spinner from "react-bootstrap/Spinner";
import logo from "./img/logo.png"

class UnknownAuthenticatedApp extends Component {

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
                            <Spinner animation="border" variant="primary"/>
                        </Navbar.Collapse>
                    </Navbar>
                </div>
                <div className="App-divider"/>
            </div>
        );
    }
}

export default UnknownAuthenticatedApp;