import React, {Component} from 'react';
import './App.css';
import UnknownAuthenticatedApp from "./UnknownAuthenticatedApp";
import AuthenticatedApp from "./AuthenticatedApp";
import NonAuthenticatedApp from "./NonAuthenticatedApp";
import ReactGA from 'react-ga';

if (process.env.NODE_ENV === "production") {
    ReactGA.initialize(process.env.REACT_APP_GA_TRACKER);
}

class App extends Component {

    state = {
        is_authenticated: null,
        config: null
    };

    render() {
        if (this.state.is_authenticated == null || this.state.config == null) {
            return (<UnknownAuthenticatedApp/>)
        } else if (this.state.is_authenticated === false) {
            return (<NonAuthenticatedApp/>);
        } else {
            return (<AuthenticatedApp config={this.state.config}/>);
        }
    }

    componentDidMount() {
        fetch(process.env.REACT_APP_API_URL + "/auth/isAuthenticated", {
            method: "GET",
            credentials: 'include'
        })
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        is_authenticated: result.is_authenticated,
                        config: this.state.config
                    })
                },
                (error) => {
                    console.log(error)
                }
            );

        fetch(process.env.REACT_APP_API_URL + "/dashboardConfig", {
            method: "GET",
            credentials: 'include'
        })
            .then(res => res.json())
            .then(
                (result) => {
                    if (result.debug === false) {
                        ReactGA.pageview("/");
                    }

                    this.setState({
                        is_authenticated: this.state.is_authenticated,
                        config: result
                    });
                },
                (error) => {
                    console.log(error)
                }
            )
    }
}

export default App;