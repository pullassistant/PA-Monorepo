import React, {Component} from 'react';
import './App.css';
import Navbar from 'react-bootstrap/Navbar'
import Button from 'react-bootstrap/Button'
import Installation from "./Installation";
import logo from "./img/logo.png"

class AuthenticatedApp extends Component {

    state = {
        installations: []
    };

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
                            <Button href={process.env.REACT_APP_API_URL + "/auth/signOut"}
                                    variant="outline-secondary">Sign out</Button>
                        </Navbar.Collapse>
                    </Navbar>
                </div>
                <div className="App-divider"/>
                <div className="App-content">
                    <h2>Pull Assistant installations</h2>
                    <div className="App-content-nav">
                        <Button href={this.props.config.marketplace_url} variant="primary" target="_blank">New
                            installation</Button>
                    </div>
                    {
                        this.state.installations && this.state.installations.map((item, key) =>
                            <Installation installation={item} config={this.props.config}/>)
                    }
                </div>
            </div>
        );
    }

    componentDidMount() {
        fetch(process.env.REACT_APP_API_URL + "/installations", {
            method: "GET",
            credentials: 'include'
        })
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState(result)
                },
                (error) => {
                    console.error(error)
                }
            )
    }
}

export default AuthenticatedApp;